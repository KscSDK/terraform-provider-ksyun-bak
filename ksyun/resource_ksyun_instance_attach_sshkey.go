package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunAttachKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunAttachKeyCreate,
		Read:   resourceKsyunAttachKeyRead,
		Delete: resourceKsyunAttachKeyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
func resourceKsyunAttachKeyCreate(d *schema.ResourceData, m interface{}) error {
	kecConn := m.(*KsyunClient).kecconn
	req := make(map[string]interface{})
	creates := []string{
		"key_id",
		"instance_id",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "AttachKey"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := kecConn.AttachKey(&req)
	logger.Debug(logger.AllFormat, action, req, *resp, err)
	if err != nil {
		return fmt.Errorf("Error AttachKey : %s", err)
	}
	status, ok := (*resp)["Return"]
	if !ok {
		return fmt.Errorf("Error AttachKey ")
	}
	status1, ok := status.(bool)
	if !ok || !status1 {
		return fmt.Errorf("Error AttachKey:fail ")
	}
	d.SetId(fmt.Sprintf("%s:%s", d.Get("instance_id"), d.Get("key_id")))
	return resourceKsyunAttachKeyRead(d, m)
}

func resourceKsyunAttachKeyRead(d *schema.ResourceData, m interface{}) error {
	kecConn := m.(*KsyunClient).kecconn
	p := strings.Split(d.Id(), ":")
	if len(p) != 2 {
		return fmt.Errorf("error attachkeyid:%v", d.Id())
	}
	req := make(map[string]interface{})
	req["InstanceId.1"] = p[0]
	if pd, ok := d.GetOk("project_id"); ok {
		req["project_id"] = fmt.Sprintf("%v", pd)
	}
	action := "DescribeInstances"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := kecConn.DescribeInstances(&req)
	if err != nil {
		return fmt.Errorf("error on reading Instance %q, %s", d.Id(), err)
	}
	logger.Debug(logger.AllFormat, action, req, *resp, err)
	itemset, ok := (*resp)["InstancesSet"]
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	item, ok := items[0].(map[string]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	keySet, ok := item["KeySet"]
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	keys, ok := keySet.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	for k, v := range keys {
		if fmt.Sprintf("%v", v) == p[1] {
			d.Set("key_id", v)
			break
		}
		if k == len(keys)-1 {
			d.SetId("")
			return nil
		}
	}
	d.Set("instance_id", p[0])
	return nil
}

func resourceKsyunAttachKeyDelete(d *schema.ResourceData, m interface{}) error {
	kecConn := m.(*KsyunClient).kecconn
	deleteReq := make(map[string]interface{})
	p := strings.Split(d.Id(), ":")
	if len(p) != 2 {
		return fmt.Errorf("error attachkeyid:%v", d.Id())
	}
	deleteReq["InstanceId.1"] = p[0]
	deleteReq["KeyId.1"] = p[1]
	action := "DisAttachKey"
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		logger.Debug(logger.ReqFormat, action, deleteReq)
		resp, err1 := kecConn.DetachKey(&deleteReq)
		logger.Debug(logger.AllFormat, action, deleteReq, *resp, err1)
		if err1 == nil || (err1 != nil && notFoundError(err1)) {
			return nil
		}
		if err1 != nil && inUseError(err1) {
			return resource.RetryableError(err1)
		}
		req := make(map[string]interface{})
		req["AllocationId.1"] = p[0]
		action = "DescribeAddresses"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := kecConn.DescribeInstances(&req)
		logger.Debug(logger.AllFormat, action, req, *resp, err)
		if err != nil && notFoundError(err1) {
			return nil
		}
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on  reading eip when delete %q, %s", d.Id(), err))
		}
		addressesSets, ok := (*resp)["AddressesSet"]
		if !ok {
			return nil
		}
		addsets, ok := addressesSets.([]interface{})
		if !ok || len(addsets) == 0 {
			return nil
		}
		addset, ok := addsets[0].(map[string]interface{})
		if !ok {
			return nil
		}
		if instanceId, ok := addset["InstanceId"]; ok {
			if instanceId == p[1] {
				return resource.NonRetryableError(fmt.Errorf("the specified DisAttachKey %q has not been deleted", d.Id()))
			}
		}
		return nil
	})
}
