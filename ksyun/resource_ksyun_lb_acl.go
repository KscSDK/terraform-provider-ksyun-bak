package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunLoadBalancerAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunLoadBalancerAclCreate,
		Read:   resourceKsyunLoadBalancerAclRead,
		Update: resourceKsyunLoadBalancerAclUpdate,
		Delete: resourceKsyunLoadBalancerAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_acl_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"load_balancer_acl_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"load_balancer_acl_entry_set": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"load_balancer_acl_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_acl_entry_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rule_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}
func resourceKsyunLoadBalancerAclCreate(d *schema.ResourceData, m interface{}) error {
	Slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	creates := []string{
		"load_balancer_acl_name",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "CreateLoadBalancerAcl"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := Slbconn.CreateLoadBalancerAcl(&req)
	if err != nil {
		return fmt.Errorf("create LoadBalancerAcl : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)

	loadBalancerAcl, ok := (*resp)["LoadBalancerAcl"]
	if !ok {
		return fmt.Errorf("create LoadBalancerAcl : no LoadBalancerAcl found")
	}

	p, ok := loadBalancerAcl.(map[string]interface{})
	if !ok {
		return fmt.Errorf("create LoadBalancerAcl : no LoadBalancerAcl data found")
	}
	id, ok := p["LoadBalancerAclId"]
	if !ok {
		return fmt.Errorf("create LoadBalancerAcl : no LoadBalancerAcl id found")
	}
	ids, ok := id.(string)
	if !ok {
		return fmt.Errorf("create LoadBalancerAcl : no LoadBalancerAcl id found")
	}
	d.Set("load_balancer_acl_id", id)
	d.SetId(ids)
	return resourceKsyunLoadBalancerAclRead(d, m)
}

func resourceKsyunLoadBalancerAclRead(d *schema.ResourceData, m interface{}) error {
	Slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["LoadBalancerAclId.1"] = d.Id()
	action := "DescribeLoadBalancerAcls"
	logger.Debug(logger.ReqFormat, action, req)

	resp, err := Slbconn.DescribeLoadBalancerAcls(&req)
	if err != nil {
		return fmt.Errorf(" read LoadBalancerAcls : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)

	resSet := (*resp)["LoadBalancerAclSet"]
	res, ok := resSet.([]interface{})
	if !ok || len(res) == 0 {
		d.SetId("")
		return nil
	}
	subPara := SetDByResp(d, res[0], lbAclKeys, map[string]bool{"LoadBalancerAclEntrySet": true})
	lbes, ok := subPara["LoadBalancerAclEntrySet"].([]interface{})
	if ok {
		subRes := GetSubSliceDByRep(lbes, lbAclEntryKeys)
		d.Set("load_balancer_acl_entry_set", subRes)
	}
	return nil
}

func resourceKsyunLoadBalancerAclUpdate(d *schema.ResourceData, m interface{}) error {
	Slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["LoadBalancerAclId"] = d.Id()
	allAttributes := []string{
		"load_balancer_acl_name",
	}
	attributeUpdate := false
	var updates []string
	//获取修改属性
	for _, v := range allAttributes {
		if d.HasChange(v) {
			attributeUpdate = true
			updates = append(updates, v)
		}
	}
	if !attributeUpdate {
		return nil
	}
	//创建修改请求
	for _, v := range allAttributes {
		if v1, ok := d.GetOk(v); ok {
			req[Downline2Hump(v)] = fmt.Sprintf("%v", v1)
		}
	}
	// 开启 允许部分属性修改 功能
	d.Partial(true)
	action := "ModifyLoadBalancerAcl"
	logger.Debug(logger.ReqFormat, action, req)

	resp, err := Slbconn.ModifyLoadBalancerAcl(&req)
	if err != nil {
		return fmt.Errorf("update LoadBalancerAcl (%v)error:%v", req, err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)

	// 设置部分修改属性
	for _, v := range updates {
		d.SetPartial(v)
	}
	d.Partial(false)
	return resourceKsyunLoadBalancerAclRead(d, m)
}

func resourceKsyunLoadBalancerAclDelete(d *schema.ResourceData, m interface{}) error {
	Slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["LoadBalancerAclId"] = d.Id()
	/*
		_, err := Slbconn.DeleteLoadBalancerAcl(&req)
		if err != nil {
			return fmt.Errorf("delete LoadBalancerAcl error:%v", err)
		}
		return nil
	*/
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		action := "DeleteLoadBalancerAcl"
		logger.Debug(logger.ReqFormat, action, req)

		if resp, err := Slbconn.DeleteLoadBalancerAcl(&req); err != nil {
			if strings.Contains(err.Error(), "NotFound") {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("error on deleting lbacl %q, %s", d.Id(), err))
		} else {
			logger.Debug(logger.RespFormat, action, req, *resp)
		}
		req := make(map[string]interface{})
		req["LoadBalancerAclId"] = d.Id()
		action = "DescribeLoadBalancerAcls"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := Slbconn.DescribeLoadBalancerAcls(&req)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on reading lbacl when deleting %q, %s", d.Id(), err))
		}
		logger.Debug(logger.RespFormat, action, req, *resp)

		itemSet, ok := (*resp)["LoadBalancerAclSet"]
		if !ok {
			return nil
		}
		items, ok := itemSet.([]interface{})
		if !ok || len(items) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf(" the specified lbacl %q has not been deleted due to unknown error", d.Id()))
	})
}
