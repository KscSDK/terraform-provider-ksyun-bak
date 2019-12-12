package ksyun

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/ks3sdklib/aws-sdk-go/aws"
	"github.com/ks3sdklib/aws-sdk-go/aws/awserr"
	"github.com/ks3sdklib/aws-sdk-go/service/s3"
)

func TestAccKsyunKs3Bucket_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "ksyun_ks3.bucket-create"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName:   resourceName,
		IDRefreshIgnore: []string{"force_destroy"},
		Providers:       testAccProviders,
		CheckDestroy:    testAccCheckKsyunKs3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKsyunKs3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKsyunKs3BucketExists(resourceName),
				),
			},
		},
	})
}

func TestAccKsyunKs3Bucket_Cors_Update(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "ksyun_ks3.bucket-cors"

	updateBucketCors := func(n string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources[n]
			if !ok {
				return fmt.Errorf("Not found: %s", n)
			}

			conn := testAccProvider.Meta().(*KsyunClient).ks3conn
			_, err := conn.PutBucketCORS(&s3.PutBucketCORSInput{
				Bucket: aws.String(rs.Primary.ID),
				CORSConfiguration: &s3.CORSConfiguration{
					CORSRules: []*s3.CORSRule{
						{
							AllowedHeaders: []*string{aws.String("*")},
							AllowedMethods: []*string{aws.String("GET")},
							AllowedOrigins: []*string{aws.String("https://www.example.com")},
						},
					},
				},
			})
			if err != nil && !isAWSErr(err, "NoSuchCORSConfiguration", "") {
				return err
			}
			return nil
		}
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKsyunKs3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKsyunKs3BucketConfigWithCORS(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKsyunKs3BucketExists(resourceName),
					testAccCheckKsyunKs3BucketCors(
						resourceName,
						[]*s3.CORSRule{
							{
								AllowedHeaders: []*string{aws.String("*")},
								AllowedMethods: []*string{aws.String("PUT"), aws.String("POST")},
								AllowedOrigins: []*string{aws.String("https://www.example.com")},
								ExposeHeaders:  []*string{aws.String("x-amz-server-side-encryption"), aws.String("ETag")},
								MaxAgeSeconds:  Int64(3000),
							},
						},
					),
					updateBucketCors(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_destroy", "acl"},
			},
			{
				Config: testAccKsyunKs3BucketConfigWithCORS(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKsyunKs3BucketExists(resourceName),
					testAccCheckKsyunKs3BucketCors(
						resourceName,
						[]*s3.CORSRule{
							{
								AllowedHeaders: []*string{aws.String("*")},
								AllowedMethods: []*string{aws.String("PUT"), aws.String("POST")},
								AllowedOrigins: []*string{aws.String("https://www.example.com")},
								ExposeHeaders:  []*string{aws.String("x-amz-server-side-encryption"), aws.String("ETag")},
								MaxAgeSeconds:  Int64(3000),
							},
						},
					),
				),
			},
		},
	})
}

func TestAccKsyunKs3Bucket_Logging(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "ksyun_ks3.bucket-logging"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKsyunKs3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKsyunKs3BucketConfigWithLogging(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKsyunKs3BucketExists(resourceName),
					testAccCheckKsyunKs3BucketLogging(
						resourceName, "ksyun_ks3.bucket-target", "log/"),
				),
			},
		},
	})
}

func testAccCheckKsyunKs3BucketDestroy(s *terraform.State) error {
	return testAccCheckKsyunKs3BucketDestroyWithProvider(s, testAccProvider)
}

func testAccCheckKsyunKs3BucketDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	conn := provider.Meta().(*KsyunClient).ks3conn

	fmt.Println(len(s.RootModule().Resources))
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ksyun_ks3" {
			continue
		}

		input := &s3.HeadBucketInput{
			Bucket: aws.String(rs.Primary.ID),
		}

		// Retry for S3 eventual consistency
		err := resource.Retry(1*time.Minute, func() *resource.RetryError {
			_, err := conn.HeadBucket(input)

			if err.(awserr.RequestFailure).StatusCode() == 404 {
				return nil
			}

			if err != nil {
				return resource.NonRetryableError(err)
			}

			return resource.RetryableError(fmt.Errorf("KS3 Bucket still exists: %s", rs.Primary.ID))
		})

		if isResourceTimeoutError(err) {
			_, err = conn.HeadBucket(input)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckKsyunKs3BucketExists(n string) resource.TestCheckFunc {
	return testAccCheckKsyunKs3BucketExistsWithProvider(n, func() *schema.Provider { return testAccProvider })
}

func testAccCheckKsyunKs3BucketExistsWithProvider(n string, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		provider := providerF()

		conn := provider.Meta().(*KsyunClient).ks3conn
		_, err := conn.HeadBucket(&s3.HeadBucketInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			if isAWSErr(err, "NoSuchBucket", "") {
				return fmt.Errorf("KS3 bucket not found")
			}
			return err
		}
		return nil

	}
}

func testAccCheckKsyunKs3BucketCors(n string, corsRules []*s3.CORSRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		conn := testAccProvider.Meta().(*KsyunClient).ks3conn

		out, err := conn.GetBucketCORS(&s3.GetBucketCORSInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() != "NoSuchCORSConfiguration" {
				return fmt.Errorf("GetBucketCors error: %v", err)
			}
		}

		if !reflect.DeepEqual(out.CORSRules, corsRules) {
			return fmt.Errorf("bad error cors rule, expected: %v, got %v", corsRules, out.CORSRules)
		}

		return nil
	}
}

func testAccCheckKsyunKs3BucketLogging(n, b, p string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]
		conn := testAccProvider.Meta().(*KsyunClient).ks3conn

		out, err := conn.GetBucketLogging(&s3.GetBucketLoggingInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("GetBucketLogging error: %v", err)
		}

		if out.LoggingEnabled == nil {
			return fmt.Errorf("logging not enabled for bucket: %s", rs.Primary.ID)
		}

		tb := s.RootModule().Resources[b]

		if v := out.LoggingEnabled.TargetBucket; v == nil {
			if tb.Primary.ID != "" {
				return fmt.Errorf("bad target bucket, found nil, expected: %s", tb.Primary.ID)
			}
		} else {
			if *v != tb.Primary.ID {
				return fmt.Errorf("bad target bucket, expected: %s, got %s", tb.Primary.ID, *v)
			}
		}

		if v := out.LoggingEnabled.TargetPrefix; v == nil {
			if p != "" {
				return fmt.Errorf("bad target prefix, found nil, expected: %s", p)
			}
		} else {
			if *v != p {
				return fmt.Errorf("bad target prefix, expected: %s, got %s", p, *v)
			}
		}

		return nil
	}
}

func testAccKsyunKs3BucketConfig(randInt int) string {
	return fmt.Sprintf(`
resource "ksyun_ks3" "bucket-create" {
  bucket = "tf-test-bucket-create-%d"
  acl    = "public-read"
}
`, randInt)
}

func testAccKsyunKs3BucketConfigWithCORS(randInt int) string {
	return fmt.Sprintf(`
resource "ksyun_ks3" "bucket-cors" {
  bucket = "tf-test-bucket-cors-%d"
  acl    = "public-read"

  cors_rule {
    allowed_header = ["*"]
    allowed_method = ["PUT", "POST"]
    allowed_origin = ["https://www.example.com"]
    expose_header  = ["x-amz-server-side-encryption", "ETag"]
    max_age_seconds = 3000
  }
}
`, randInt)
}

func testAccKsyunKs3BucketConfigWithLogging(randInt int) string {
	return fmt.Sprintf(`
resource "ksyun_ks3" "bucket-target" {
  bucket = "tf-test-bucket-target-%d"
  acl    = "public-read"
}

resource "ksyun_ks3" "bucket-logging" {
  bucket = "tf-test-bucket-log-%d"
  acl    = "private"

  logging {
    target_bucket = "${ksyun_ks3.bucket-target.id}"
    target_prefix = "log/"
  }
}
`, randInt, randInt)
}
