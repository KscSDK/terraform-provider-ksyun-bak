package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunHealthCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunHealthCheckCreate,
		Read:   resourceKsyunHealthCheckRead,
		Update: resourceKsyunHealthCheckUpdate,
		Delete: resourceKsyunHealthCheckDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"health_check_state": {
				Type:     schema.TypeString,
				Required: true,
			},
			"healthy_threshold": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"unhealthy_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"url_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_default_host_name": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"health_check_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}
func resourceKsyunHealthCheckCreate(d *schema.ResourceData, m interface{}) error {
	Slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	/*
		if v, ok := d.GetOk("listener_id"); ok {
			req["ListenerId"] = v
		}
		if v, ok := d.GetOk("healthy_threshold"); ok {
			req["HealthyThreshold"] = v
		}
		if v, ok := d.GetOk("interval"); ok {
			req["Interval"] = v
		}
		if v, ok := d.GetOk("timeout"); ok {
			req["Timeout"] = v
		}
		if v, ok := d.GetOk("unhealthy_threshold"); ok {
			req["UnhealthyThreshold"] = v
		}
		if v, ok := d.GetOk("health_check_state"); ok {
			req["HealthCheckState"] = v
		}
		if v, ok := d.GetOk("url_path"); ok {
			req["UrlPath"] = v
		}
		if v, ok := d.GetOk("is_default_host_name"); ok {
			req["IsDefaultHostName"] = v
		}
		if v, ok := d.GetOk("host_name"); ok {
			req["HostName"] = v
		}
	*/
	creates := []string{
		"listener_id",
		"healthy_threshold",
		"interval",
		"timeout",
		"unhealthy_threshold",
		"health_check_state",
		"url_path",
		"is_default_host_name",
		"host_name",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			req[vv] = fmt.Sprintf("%v", v1)
		}
	}
	action := "ConfigureHealthCheck"
	logger.Debug(logger.ReqFormat, action, req)

	resp, err := Slbconn.ConfigureHealthCheck(&req)
	if err != nil {
		return fmt.Errorf("create HealthCheck : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	id, ok := (*resp)["HealthCheckId"]
	if !ok {
		return fmt.Errorf("create HealthCheck : no HealthCheckId found")
	}
	idres, ok := id.(string)
	if !ok {
		return fmt.Errorf("create HealthCheck : no HealthCheckId found")
	}
	d.Set("health_check_id", idres)
	d.SetId(idres)
	return resourceKsyunHealthCheckRead(d, m)
}

func resourceKsyunHealthCheckRead(d *schema.ResourceData, m interface{}) error {
	Slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["HealthCheckId.1"] = d.Id()
	action := "DescribeHealthChecks"
	logger.Debug(logger.ReqFormat, action, req)

	resp, err := Slbconn.DescribeHealthChecks(&req)
	if err != nil {
		return fmt.Errorf(" read HealthChecks : %s", err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)

	healthCheckSet := (*resp)["HealthCheckSet"]
	health, ok := healthCheckSet.([]interface{})
	if !ok || len(health) == 0 {
		d.SetId("")
		return nil
	}
	SetDByResp(d, health[0], healthCheckKeys, map[string]bool{})

	return nil
}

func resourceKsyunHealthCheckUpdate(d *schema.ResourceData, m interface{}) error {
	Slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["HealthCheckId"] = d.Id()
	if _, ok := d.GetOk("health_check_state"); !ok {
		return fmt.Errorf("cann't change health_check_state to empty string")
	}
	allAttributes := []string{
		"health_check_state",
		"healthy_threshold",
		"interval",
		"timeout",
		"unhealthy_threshold",
		"is_default_host_name",
		"host_name",
		"health_check_state",
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
	//必填字端，未修改也要传
	if _, ok := req["HealthCheckState"]; !ok {
		req["HealthCheckState"] = d.Get("health_check_state")
	}
	// 开启 允许部分属性修改 功能
	d.Partial(true)
	action := "ModifyHealthCheck"
	logger.Debug(logger.ReqFormat, action, req)

	resp, err := Slbconn.ModifyHealthCheck(&req)
	if err != nil {
		return fmt.Errorf("update HealthCheck (%v)error:%v", req, err)
	}
	logger.Debug(logger.RespFormat, action, req, *resp)
	// 设置部分修改属性
	for _, v := range updates {
		d.SetPartial(v)
	}
	d.Partial(false)
	return resourceKsyunHealthCheckRead(d, m)
}

func resourceKsyunHealthCheckDelete(d *schema.ResourceData, m interface{}) error {
	Slbconn := m.(*KsyunClient).slbconn
	req := make(map[string]interface{})
	req["HealthCheckId"] = d.Id()
	/*
		_, err := Slbconn.DeleteHealthCheck(&req)
		if err != nil {
			return fmt.Errorf("delete HealthCheck error:%v", err)
		}
		return nil
	*/
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		action := "DeleteHealthCheck"
		logger.Debug(logger.ReqFormat, action, req)
		if resp, err := Slbconn.DeleteHealthCheck(&req); err != nil {
			if strings.Contains(err.Error(), "NotFound") {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("error on deleting healthcheck %q, %s", d.Id(), err))
		} else {
			logger.Debug(logger.RespFormat, action, req, *resp)

		}
		req := make(map[string]interface{})
		req["HealthCheckId.1"] = d.Id()
		action = "DescribeHealthChecks"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := Slbconn.DescribeHealthChecks(&req)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on reading healthcheck when deleting %q, %s", d.Id(), err))
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		items, ok := (*resp)["HealthCheckSet"]
		if !ok {
			return nil
		}
		itemsspe, ok := items.([]interface{})
		if !ok || len(itemsspe) == 0 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf(" the specified healthcheck %q has not been deleted due to unknown error", d.Id()))
	})
}
