package ksyun

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/ks3sdklib/aws-sdk-go/aws"
	"github.com/ks3sdklib/aws-sdk-go/service/s3"
	"log"
	"time"
)

func resourceKsyunKs3Bucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunKs3BucketCreate,
		Update: resourceKsyunKs3BucketUpdate,
		Read:   resourceKsyunKs3BucketRead,
		Delete: resourceKsyunKs3BucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 63),
			},
			"type": {
				Type:     schema.TypeString,
				Default:  "NORMAL", //NORMAL 普通 ARCHIVE 归档
				Optional: true,
			},
			"acl": {
				Type:     schema.TypeString,
				Default:  "private", //private public-read public-read-write
				Optional: true,
			},
			"logging": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"affect": {
							Type:         schema.TypeString,
							Default:      "disable", //enable disable
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"enable", "disable"}, true),
						},
						"target_bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"target_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: func(v interface{}) int {
					var buf bytes.Buffer
					m := v.(map[string]interface{})
					buf.WriteString(fmt.Sprintf("%s-", m["target_bucket"]))
					return hashcode.String(buf.String())
				},
			},
			"cors_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_origins": { //指定哪个域名可访问 www.example.com
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_headers": { //指定head enable disable
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"allowed_methods": { //访问域名 GET PUT POST DELETE HEAD
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"expose_headers": { //在响应中暴露的头
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"max_age_seconds": { //浏览器缓存时间
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"referer_config": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"affect": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"enable", "disable"}, true),
						},
						"allow_empty": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"referers": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"lifecycle_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(0, 255),
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"status": {
							Type:         schema.TypeString,
							Required:     true, //enable disable
							ValidateFunc: validation.StringInSlice([]string{"enable", "disable"}, true),
						},
						"expiration": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      expirationHash,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateKs3BucketLifecycleTimestamp,
									},
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"expired_object_delete_marker": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"transition": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      transitionHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateKs3BucketLifecycleTimestamp,
									},
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"storage_class": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateKs3BucketLifecycleTransitionStorageClass(),
									},
								},
							},
						},
					},
				},
			},
			"policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateJsonString,
				//DiffSuppressFunc: suppressEquivalentKsyunPolicyDiffs,
			},
		},
	}
}

func resourceKsyunKs3BucketCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.ks3conn

	var bucket string
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		bucket = resource.UniqueId()
	}
	d.Set("bucket", bucket)
	acl := d.Get("acl").(string)

	log.Printf("[DEBUG] KS3 bucket create: %s, ACL: %s", bucket, acl)

	req := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String(acl),
	}

	if err := validateKs3BucketName(bucket); err != nil {
		return fmt.Errorf("Error validating KS3 bucket name: %s", err)
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] Trying to create new S3 bucket: %q", bucket)
		_, err := conn.CreateBucket(req)
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "OperationAborted" {
				log.Printf("[WARN] Got an error while trying to create S3 bucket %s: %s", bucket, err)
				return resource.RetryableError(
					fmt.Errorf("Error creating S3 bucket %s, retrying: %s", bucket, err))
			}
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if isResourceTimeoutError(err) {
		_, err = conn.CreateBucket(req)
	}
	if err != nil {
		return fmt.Errorf("Error creating S3 bucket: %s", err)
	}

	// Assign the bucket name as the resource ID
	d.SetId(bucket)
	return resourceKsyunKs3BucketUpdate(d, meta)
}

func resourceKsyunKs3BucketUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.ks3conn
	if d.HasChange("cors_rule") {
		if err := resourceKsyunKs3BucketCorsUpdate(conn, d); err != nil {
			return err
		}
	}
	if d.HasChange("acl") && !d.IsNewResource() {
		if err := resourceKsyunKs3BucketAclUpdate(conn, d); err != nil {
			return err
		}
	}
	if d.HasChange("logging") {
		if err := resourceKsyunKs3BucketLoggingUpdate(conn, d); err != nil {
			return err
		}
	}
	return resourceKsyunKs3BucketRead(d, meta)
}

func resourceKsyunKs3BucketRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.ks3conn

	var err error

	_, err = retryOnAwsCode("NoSuchBucket", func() (interface{}, error) {
		return conn.HeadBucket(&s3.HeadBucketInput{
			Bucket: aws.String(d.Id()),
		})
	})

	if err != nil {
		if awsError, ok := err.(awserr.RequestFailure); ok && awsError.StatusCode() == 404 {
			log.Printf("[WARN] KS3 Bucket (%s) not found, error code (404)", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error reading KS3 Bucket (%s): %s", d.Id(), err)
	}

	// In the import case, we won't have this
	if _, ok := d.GetOk("bucket"); !ok {
		d.Set("bucket", d.Id())
	}

	// Read the CORS
	corsResponse, err := retryOnAwsCode("NoSuchBucket", func() (interface{}, error) {
		return conn.GetBucketCORS(&s3.GetBucketCORSInput{
			Bucket: aws.String(d.Id()),
		})
	})
	if err != nil && !isAWSErr(err, "NoSuchCORSConfiguration", "") {
		return fmt.Errorf("error getting S3 Bucket CORS configuration: %s", err)
	}

	corsRules := make([]map[string]interface{}, 0)
	if cors, ok := corsResponse.(*s3.GetBucketCORSOutput); ok && len(cors.CORSRules) > 0 {
		corsRules = make([]map[string]interface{}, 0, len(cors.CORSRules))
		for _, ruleObject := range cors.CORSRules {
			rule := make(map[string]interface{})
			rule["allowed_headers"] = flattenStringList(ruleObject.AllowedHeaders)
			rule["allowed_methods"] = flattenStringList(ruleObject.AllowedMethods)
			rule["allowed_origins"] = flattenStringList(ruleObject.AllowedOrigins)
			// Both the "ExposeHeaders" and "MaxAgeSeconds" might not be set.
			if ruleObject.AllowedOrigins != nil {
				rule["expose_headers"] = flattenStringList(ruleObject.ExposeHeaders)
			}
			if ruleObject.MaxAgeSeconds != nil {
				rule["max_age_seconds"] = int(*ruleObject.MaxAgeSeconds)
			}
			corsRules = append(corsRules, rule)
		}
	}
	if err := d.Set("cors_rule", corsRules); err != nil {
		return fmt.Errorf("error setting cors_rule: %s", err)
	}

	// Read the logging configuration
	loggingResponse, err := retryOnAwsCode("NoSuchBucket", func() (interface{}, error) {
		return conn.GetBucketLogging(&s3.GetBucketLoggingInput{
			Bucket: aws.String(d.Id()),
		})
	})

	if err != nil {
		return fmt.Errorf("error getting KS3 Bucket logging: %s", err)
	}

	lcl := make([]map[string]interface{}, 0, 1)
	if logging, ok := loggingResponse.(*s3.GetBucketLoggingOutput); ok && logging.LoggingEnabled != nil {
		v := logging.LoggingEnabled
		lc := make(map[string]interface{})
		if *v.TargetBucket != "" {
			lc["target_bucket"] = *v.TargetBucket
		}
		if *v.TargetPrefix != "" {
			lc["target_prefix"] = *v.TargetPrefix
		}
		lcl = append(lcl, lc)
	}
	if err := d.Set("logging", lcl); err != nil {
		return fmt.Errorf("error setting logging: %s", err)
	}

	return nil
}

func resourceKsyunKs3BucketDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*KsyunClient)
	conn := client.ks3conn

	log.Printf("[DEBUG] KS3 Delete Bucket: %s", d.Id())
	_, err := conn.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(d.Id()),
	})

	if isAWSErr(err, "NoSuchBucket", "") {
		return nil
	}

	if err != nil {
		return fmt.Errorf("error deleting S3 Bucket (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceKsyunKs3BucketLoggingUpdate(s3conn *s3.S3, d *schema.ResourceData) error {
	logging := d.Get("logging").(*schema.Set).List()
	bucket := d.Get("bucket").(string)
	loggingStatus := &s3.BucketLoggingStatus{}

	if len(logging) > 0 {
		c := logging[0].(map[string]interface{})

		loggingEnabled := &s3.LoggingEnabled{}
		if val, ok := c["target_bucket"]; ok {
			loggingEnabled.TargetBucket = aws.String(val.(string))
		}
		if val, ok := c["target_prefix"]; ok {
			loggingEnabled.TargetPrefix = aws.String(val.(string))
		}

		loggingStatus.LoggingEnabled = loggingEnabled
	}

	contentType := "application/xml"
	i := &s3.PutBucketLoggingInput{
		Bucket:              aws.String(bucket),
		BucketLoggingStatus: loggingStatus,
		ContentType:         &contentType,
	}
	log.Printf("[DEBUG] KS3 put bucket logging: %#v", i)

	_, err := retryOnAwsCode("NoSuchBucket", func() (interface{}, error) {
		return s3conn.PutBucketLogging(i)
	})
	if err != nil {
		return fmt.Errorf("Error putting KS3 logging: %s", err)
	}

	return nil
}

func resourceKsyunKs3BucketAclUpdate(s3conn *s3.S3, d *schema.ResourceData) error {
	acl := d.Get("acl").(string)
	bucket := d.Get("bucket").(string)

	i := &s3.PutBucketACLInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String(acl),
	}
	log.Printf("[DEBUG] S3 put bucket ACL: %#v", i)

	_, err := retryOnAwsCode("NoSuchBucket", func() (interface{}, error) {
		return s3conn.PutBucketACL(i)
	})
	if err != nil {
		return fmt.Errorf("Error putting S3 ACL: %s", err)
	}

	return nil
}

func resourceKsyunKs3BucketCorsUpdate(s3conn *s3.S3, d *schema.ResourceData) error {
	bucket := d.Get("bucket").(string)
	rawCors := d.Get("cors_rule").([]interface{})

	if len(rawCors) == 0 {
		// Delete CORS
		log.Printf("[DEBUG] KS3 bucket: %s, delete CORS", bucket)

		_, err := retryOnAwsCode("NoSuchBucket", func() (interface{}, error) {
			return s3conn.DeleteBucketCORS(&s3.DeleteBucketCORSInput{
				Bucket: aws.String(bucket),
			})
		})
		if err != nil {
			return fmt.Errorf("Error deleting KS3 CORS: %s", err)
		}
	} else {
		// Put CORS
		rules := make([]*s3.CORSRule, 0, len(rawCors))
		for _, cors := range rawCors {
			corsMap := cors.(map[string]interface{})
			r := &s3.CORSRule{}
			for k, v := range corsMap {
				log.Printf("[DEBUG] KS3 bucket: %s, put CORS: %#v, %#v", bucket, k, v)
				if k == "max_age_seconds" {
					r.MaxAgeSeconds = Int64(int64(v.(int)))
				} else {
					vMap := make([]*string, len(v.([]interface{})))
					for i, vv := range v.([]interface{}) {
						if str, ok := vv.(string); ok {
							vMap[i] = aws.String(str)
						}
					}
					switch k {
					case "allowed_headers":
						r.AllowedHeaders = vMap
					case "allowed_methods":
						r.AllowedMethods = vMap
					case "allowed_origins":
						r.AllowedOrigins = vMap
					case "expose_headers":
						r.ExposeHeaders = vMap
					}
				}
			}
			rules = append(rules, r)
		}
		corsInput := &s3.PutBucketCORSInput{
			Bucket: aws.String(bucket),
			CORSConfiguration: &s3.CORSConfiguration{
				CORSRules: rules,
			},
		}
		log.Printf("[DEBUG] KS3 bucket: %s, put CORS: %#v", bucket, corsInput)

		_, err := retryOnAwsCode("NoSuchBucket", func() (interface{}, error) {
			return s3conn.PutBucketCORS(corsInput)
		})
		if err != nil {
			return fmt.Errorf("Error putting KS3 CORS: %s", err)
		}
	}

	return nil
}

// Int64 returns a pointer to the int64 value passed in.
func Int64(v int64) *int64 {
	return &v
}

func expirationHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	if v, ok := m["expired_object_delete_marker"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
	}
	return hashcode.String(buf.String())
}

func transitionHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	if v, ok := m["storage_class"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	return hashcode.String(buf.String())
}
