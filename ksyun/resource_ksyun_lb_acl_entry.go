package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunLoadBalancerAclEntry() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunLoadBalancerAclEntryCreate,
		Delete: resourceKsyunLoadBalancerAclEntryDelete,
		Read:   resourceKsyunLoadBalancerAclEntryRead,
		Schema: map[string]*schema.Schema{
			"load_balancer_acl_entry_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer_acl_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"rule_number": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"rule_action": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}
func resourceKsyunLoadBalancerAclEntryRead(d *schema.ResourceData, m interface{}) error {
	return nil
}
func resourceKsyunLoadBalancerAclEntryCreate(d *schema.ResourceData, m interface{}) error {
	Slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	creates := []string{
		"load_balancer_acl_id",
		"cidr_block",
		"rule_number",
		"rule_action",
		"protocol",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "CreateLoadBalancerAclEntry"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := Slbconn.CreateLoadBalancerAclEntry(&req)
	if err != nil {
		return fmt.Errorf("create LoadBalancerAclEntry : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)

	loadBalancerAclEntry, ok := (*resp)["LoadBalancerAclEntry"]
	if !ok {
		return fmt.Errorf("create LoadBalancerAclEntry : no LoadBalancerAclEntry found")
	}

	lbae, ok := loadBalancerAclEntry.(map[string]interface{})
	if !ok {
		return fmt.Errorf("create LoadBalancerAclEntry : no LoadBalancerAclEntry data found")
	}

	id, ok := lbae["LoadBalancerAclEntryId"]
	if !ok {
		return fmt.Errorf("create LoadBalancerAclEntry : no LoadBalancerAclEntry id found")
	}
	d.Set("load_balancer_acl_entry_id", id)
	ids, ok := id.(string)
	if !ok {
		return fmt.Errorf("create LoadBalancerAclEntry : no LoadBalancerAclEntry id found")
	}
	SetDByResp(d, lbae, lbAclEntryKeys, map[string]bool{})
	d.SetId(ids)
	return nil
}

func resourceKsyunLoadBalancerAclEntryDelete(d *schema.ResourceData, m interface{}) error {
	Slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["LoadBalancerAclEntryId"] = d.Id()
	req["LoadBalancerAclId"] = d.Get("load_balancer_acl_id")
	/*
		_, err := Slbconn.DeregisterInstancesFromListener(&req)
		if err != nil {
			return fmt.Errorf("delete LoadBalancerAclEntry error:%v", err)
		}
		return nil
	*/
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		action := "DeleteLoadBalancerAclEntry"
		logger.Debug(logger.ReqFormat, action, req)

		if resp, err := Slbconn.DeleteLoadBalancerAclEntry(&req); err != nil {
			if strings.Contains(err.Error(), "NotFound") {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("error on deleting lbacl %q, %s", d.Id(), err))
		} else {
			logger.Debug(logger.RespFormat, action, req, *resp)
		}
		return nil
	})
}
