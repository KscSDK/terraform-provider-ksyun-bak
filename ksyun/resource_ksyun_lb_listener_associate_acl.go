package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunListenerLBAcl() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunListenerLBAclCreate,
		Read:   resourceKsyunListenerLBAclRead,
		Delete: resourceKsyunListenerLBAclDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"load_balancer_acl_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
func resourceKsyunListenerLBAclCreate(d *schema.ResourceData, m interface{}) error {
	Slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	creates := []string{
		"load_balancer_acl_id",
		"listener_id",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}

	action := "AssociateLoadBalancerAcl"
	logger.Debug(logger.ReqFormat, action, req)
	resp, err := Slbconn.AssociateLoadBalancerAcl(&req)
	if err != nil {
		return fmt.Errorf("Error CreateListenerLBAcls : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	status, ok := (*resp)["Return"]
	if !ok {
		return fmt.Errorf("Error CreateListenerLBAcls")
	}
	statu, ok := status.(bool)
	if !ok {
		return fmt.Errorf("Error CreateListenerLBAcls ")
	}
	if !statu {
		return fmt.Errorf("Error CreateListenerLBAcls : fail")
	}
	id := fmt.Sprintf("%s:%s", d.Get("listener_id"), d.Get("load_balancer_acl_id"))
	d.SetId(id)
	return resourceKsyunListenerLBAclRead(d, m)
}

func resourceKsyunListenerLBAclRead(d *schema.ResourceData, m interface{}) error {
	Slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	p := strings.Split(d.Id(), ":")
	req["LoadBalancerAclId.1"] = p[1]
	action := "DescribeLoadBalancerAcls"
	logger.Debug(logger.ReqFormat, action, req)

	resp, err := Slbconn.DescribeLoadBalancerAcls(&req)
	if err != nil {
		return fmt.Errorf("Error DescribeListenerLBAcls : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)

	resSet := (*resp)["LoadBalancerAclSet"]
	res, ok := resSet.([]interface{})
	if !ok || len(res) == 0 {
		d.SetId("")
		return nil
	}

	subPara := SetDByResp(d, res[0], lbAclKeys, map[string]bool{"LoadBalancerAclEntrySet": true})
	datas, ok := subPara["LoadBalancerAclEntrySet"].([]interface{})
	if !ok {
		d.SetId("")
		return nil
	}
	SubSlice := GetSubSliceDByRep(datas, lbAclEntryKeys)
	d.Set("load_balancer_acl_entry_set", SubSlice)
	return nil
}

func resourceKsyunListenerLBAclDelete(d *schema.ResourceData, m interface{}) error {
	Slbconn := m.(*KsyunClient).slbconn
	p := strings.Split(d.Id(), ":")
	req := make(map[string]interface{})
	req["ListenerId"] = p[0]
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		action := "DisassociateLoadBalancerAcl"
		logger.Debug(logger.ReqFormat, action, req)

		if resp, err := Slbconn.DisassociateLoadBalancerAcl(&req); err != nil {
			if strings.Contains(err.Error(), "NotFound") {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("error on DisassociateLoadBalancerAcl %q, %s", d.Id(), err))
		} else {
			logger.Debug(logger.RespFormat, action, req, *resp)
		}
		return nil
	})
}
