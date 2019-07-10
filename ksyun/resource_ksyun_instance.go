package ksyun

//import "github.com/hashicorp/terraform/helper/schema"

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ksc/ksc-sdk-go/service/kec"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceKsyunInstanceCreate,
		Update: resourceKsyunInstanceUpdate,
		Read:   resourceKsyunInstanceRead,
		Delete: resourceKsyunInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"image_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"system_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"disk_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"data_disk_gb": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"data_disk": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"delete_with_instance": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
			/*
				"max_count": {
					Type:     schema.TypeInt,
					Required: true,
				},
				"min_count": {
					Type:     schema.TypeInt,
					Required: true,
				},
			*/
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_password": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"keep_image_login": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"charge_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"purchase_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_name_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sriov_net_support": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"data_guard_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"key_id": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"address_band_width": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"line_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"address_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"address_purchase_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"address_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_configure": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"v_c_p_u": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"g_p_u": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"data_disk_gb": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						/*
							"root_disk_gb": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"data_disk_type": {
								Type:     schema.TypeString,
								Computed: true,
							},
						*/
					},
				},
			},
			"instance_state": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"monitoring": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"network_interface_set": {
				Type:     schema.TypeList,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_interface_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_interface_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_set": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"security_group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"group_set": {
							Type:     schema.TypeList,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"d_n_s1": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"d_n_s2": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"stopped_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"product_what": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"auto_scaling_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_show_sriov_net_support": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceKsyunInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).kecconn
	var resp *map[string]interface{}
	createReq := make(map[string]interface{})
	var err error
	creates := []string{
		"image_id",
		"instance_type",
		//	"system_disk",
		"data_disk_gb",
		//	"data_disk",
		//"max_count",=1
		//"min_count",=1
		"subnet_id",
		"instance_password",
		"keep_image_login",
		"charge_type",
		"purchase_time",
		"security_group_id",
		"private_ip_address",
		"instance_name",
		"instance_name_suffix",
		"sriov_net_support",
		"project_id",
		"data_guard_id",
		"address_band_width",
		"line_id",
		"address_charge_type",
		"address_purchase_time",
		"address_project_id",
	}
	for _, v := range creates {
		if v1, ok := d.GetOk(v); ok {
			vv := Downline2Hump(v)
			createReq[vv] = fmt.Sprintf("%v", v1)
		}
	}
	createReq["MaxCount"] = "1"
	createReq["MinCount"] = "1"
	createStructs := []string{
		"system_disk",
	}
	for _, v := range createStructs {
		if v1, ok := d.GetOk(v); ok {
			FlatternStructPrefix(v1, &createReq, "SystemDisk")
		}
	}
	if v1, ok := d.GetOk("data_disk"); ok {
		FlatternStructSlicePrefix(v1, &createReq, "DataDisk")
	}
	action := "RunInstances"
	logger.Debug(logger.ReqFormat, action, createReq)
	resp, err = conn.RunInstances(&createReq)
	if err != nil {
		return fmt.Errorf("error on creating Instance: %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)
	if resp != nil {
		instances := (*resp)["InstancesSet"].([]interface{})
		if len(instances) == 0 {
			return fmt.Errorf("error on creating Instance")
		}
		Instance := instances[0].(map[string]interface{})
		InstanceId := Instance["InstanceId"].(string)
		d.SetId(InstanceId)
	}
	// after create instance, we need to wait it initialized
	stateConf := &resource.StateChangeConf{
		Pending:    []string{statusPending},
		Target:     []string{"active"},
		Refresh:    instanceStateRefreshForCreateFunc(conn, d.Id(), []string{"active"}),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForState()
	/*
		if err != nil {
			return fmt.Errorf("error on waiting for instance %q complete creating, %s", d.Id(), err)
		}
	*/
	return resourceKsyunInstanceRead(d, meta)
}

func resourceKsyunInstanceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).kecconn
	readReq := make(map[string]interface{})
	readReq["InstanceId.1"] = d.Id()
	if pd, ok := d.GetOk("project_id"); ok {
		readReq["project_id"] = fmt.Sprintf("%v", pd)
	}
	if pd, ok := d.GetOk("project_id"); ok {
		readReq["project_id"] = fmt.Sprintf("%v", pd)
	}
	action := "DescribeInstances"
	logger.Debug(logger.ReqFormat, action, readReq)
	resp, err := conn.DescribeInstances(&readReq)
	if err != nil {
		return fmt.Errorf("error on reading Instance %q, %s", d.Id(), err)
	}
	logger.Debug(logger.AllFormat, action, readReq, *resp, err)
	itemset, ok := (*resp)["InstancesSet"]
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	excludesKeys := map[string]bool{
		"InstanceConfigure":   true,
		"InstanceState":       true,
		"Monitoring":          true,
		"NetworkInterfaceSet": true,
		"SystemDisk":          true,
	}
	excludes := SetDByResp(d, items[0], instanceKeys, excludesKeys)
	if excludes["InstanceConfigure"] != nil {
		itemSet := GetSubDByRep(excludes["InstanceConfigure"], instanceConfigureKeys, map[string]bool{})
		if len(itemSet) > 0 {
			if instanceConfigure, ok := itemSet[0].(map[string]interface{}); ok {
				d.Set("data_disk_gb", instanceConfigure["data_disk_gb"])
			}
		}
		d.Set(Hump2Downline("InstanceConfigure"), itemSet)
	}
	if excludes["InstanceState"] != nil {
		itemSet := GetSubDByRep(excludes["InstanceState"], instanceStateKeys, map[string]bool{})
		d.Set(Hump2Downline("InstanceState"), itemSet)
	}
	if excludes["Monitoring"] != nil {
		itemSet := GetSubDByRep(excludes["Monitoring"], monitoringKeys, map[string]bool{})
		d.Set(Hump2Downline("Monitoring"), itemSet)
	}
	if excludes["NetworkInterfaceSet"] != nil {
		itemSet := GetSubSliceDByRep(excludes["NetworkInterfaceSet"].([]interface{}), kecNetworkInterfaceKeys)
		for k, v := range itemSet {
			if nit, ok := v["network_interface_type"]; ok {
				if nit == "primary" {
					if snId, ok := v["subnet_id"]; ok {
						d.Set("subnet_id", snId)
					}
				}
			}
			if gs, ok := v["group_set"]; ok {
				itemSetSub := GetSubSliceDByRep(gs.([]interface{}), groupSetKeys)
				itemSet[k]["group_set"] = itemSetSub
			}
			if sg, ok := v["security_group_set"]; ok {
				itemSetSub := GetSubSliceDByRep(sg.([]interface{}), kecSecurityGroupKeys)
				itemSet[k]["security_group_set"] = itemSetSub
				if len(itemSetSub) == 1 {
					if sgIdnew, ok := itemSetSub[0]["security_group_id"]; ok {
						d.Set("security_group_id", sgIdnew)
					}
				}
			}
		}
		d.Set(Hump2Downline("NetworkInterfaceSet"), itemSet)
	}
	if excludes["SystemDisk"] != nil {
		itemSet := GetSubDByRep(excludes["SystemDisk"], systemDiskKeys, map[string]bool{})
		d.Set(Hump2Downline("SystemDisk"), itemSet)
	}
	return nil
}

func resourceKsyunInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).kecconn
	d.Partial(true)
	imageUpdate := false
	updateReq := make(map[string]interface{})
	updateReq["InstanceId"] = d.Id()

	//ModifyInstanceAttribute instancename
	attributeNameUpdate := false
	updateInstanceAttributes := []string{
		"instance_name",
	}
	for _, v := range updateInstanceAttributes {
		if d.HasChange(v) && !d.IsNewResource() {
			updateReq[Downline2Hump(v)] = fmt.Sprintf("%v", d.Get(v))
			attributeNameUpdate = true
		}
	}
	if attributeNameUpdate {
		action := "ModifyInstanceAttribute"
		logger.Debug(logger.ReqFormat, action, updateReq)
		resp, err := conn.ModifyInstanceAttribute(&updateReq)
		if err != nil {
			return fmt.Errorf("error on updating  instance name, %s", err)
		}
		logger.Debug(logger.AllFormat, action, updateReq, *resp, err)
		for _, v := range updateInstanceAttributes {
			d.SetPartial(v)
		}
	}

	//ModifyInstanceImage
	updateImages := []string{
		"image_id",
		//	"system_disk",
		//"instance_password",
		"keep_image_login",
	}
	var imageUpdated []string
	for _, v := range updateImages {
		if d.HasChange(v) && !d.IsNewResource() {
			updateReq[Downline2Hump(v)] = fmt.Sprintf("%v", d.Get(v))
			imageUpdated = append(imageUpdated, v)
			imageUpdate = true
		}
	}
	if d.HasChange("system_disk") && !d.IsNewResource() {
		FlatternStructPrefix(d.Get("system_disk"), &updateReq, "SystemDisk")
		imageUpdate = true
		imageUpdated = append(imageUpdated, "system_disk")
	}
	if !imageUpdate {
		d.Partial(false)
		return resourceKsyunInstanceRead(d, meta)
	}
	/*
		passwordUpdate := false
		updatePassword := []string{
			"instance_password",
		}
		var updatedAttributePassword []string
		for _, v := range updatePassword {
			if d.HasChange(v) && !d.IsNewResource() && !imageUpdate {
				updateReq3[Downline2Hump(v)] = fmt.Sprintf("%v", d.Get(v))
				passwordUpdate = true
				updatedAttributePassword = append(updatedAttributePassword, v)
			}
		}
	*/
	var initState string
	//TODO判断初始开机状态
	readReq := make(map[string]interface{})
	readReq["InstanceId.1"] = d.Id()
	action := "DescribeInstances"
	logger.Debug(logger.ReqFormat, action, readReq)
	resp, err := conn.DescribeInstances(&readReq)
	logger.Debug(logger.AllFormat, action, readReq, *resp, err)
	if err != nil {
		return fmt.Errorf("error on reading Instance %q, %s", d.Id(), err)
	}

	itemset, ok := (*resp)["InstancesSet"]
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		return fmt.Errorf("error on reading Instance %q, %s", d.Id(), err)
	}
	state := items[0].(map[string]interface{})["InstanceState"]
	initState = state.(map[string]interface{})["Name"].(string)
	if initState == "error" {
		return fmt.Errorf("instance with error state can't be modify image")
	}
	if initState != "stopped" {
		action = "StopInstances"
		logger.Debug(logger.ReqFormat, action, readReq)
		resp, err = conn.StopInstances(&readReq) //同步
		logger.Debug(logger.AllFormat, action, readReq, *resp, err)
		if err != nil {
			return fmt.Errorf("error on stop  instance %s", err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{statusPending},
			Target:     []string{"stopped"},
			Refresh:    instanceStateRefreshFunc(conn, d.Id(), []string{"stopped"}),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      3 * time.Second,
			MinTimeout: 2 * time.Second,
		}
		if _, err = stateConf.WaitForState(); err != nil {
			return fmt.Errorf("error on waiting for starting instance when stopping %q, %s", d.Id(), err)
		}
	}
	//不支持单独修改密码
	updateReq["InstancePassword"] = fmt.Sprintf("%v", d.Get("instance_password"))
	action = "ModifyInstanceImage"
	logger.Debug(logger.ReqFormat, action, updateReq)
	resp, err = conn.ModifyInstanceImage(&updateReq)
	logger.Debug(logger.AllFormat, action, updateReq, *resp, err)
	if err != nil {
		return fmt.Errorf("error on updating instance image, %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{statusPending},
		Target:     []string{"rebuilding", "overriding"},
		Refresh:    instanceStateRefreshFunc(conn, d.Id(), []string{"rebuilding", "overriding"}),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      3 * time.Second,
		MinTimeout: 2 * time.Second,
	}
	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error on waiting for reinstalling instance when ModifyInstanceImage %q, %s", d.Id(), err)
	}
	stateConf = &resource.StateChangeConf{
		Pending: []string{statusPending},
		Target:  []string{"active"},
		//final state may be "stopped" ,need to return error
		Refresh:    instanceStateRefreshForReinstallFunc(conn, d.Id(), []string{"active"}),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      3 * time.Second,
		MinTimeout: 2 * time.Second,
	}
	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("error on waiting for starting instance when ModifyInstanceImage %q, %s", d.Id(), err)
	}
	for _, v := range imageUpdated {
		d.SetPartial(v)
	}

	//ModifyInstanceType //需要重启
	/*	typeUpdate := false
			updateInstanceTypes := []string{
				"instance_type",
				"data_disk_gb",
			}
			var typeUpdated []string
			for _, v := range updateInstanceTypes {
				if d.HasChange(v) && !d.IsNewResource() {
					updateReq[Downline2Hump(v)] = fmt.Sprintf("%v", d.Get(v))
					typeUpdate = true
					typeUpdated = append(typeUpdated, v)
				}
			}
			if typeUpdate {
				action := "ModifyInstanceType"
				logger.Debug(logger.ReqFormat, action, updateReq)
				resp, err := conn.ModifyInstanceType(&updateReq)
				if err != nil {
					return fmt.Errorf("error on updating  instance type, %s", err)
				}
				logger.Debug(logger.RespFormat, action, updateReq, *resp)
				for _, v := range typeUpdated {
					d.SetPartial(v)
				}
			}


		//后续修改需要关机
		//StopInstances
		  if d.HasChange("force_stop") && !d.IsNewResource() {
		  	if d.Get("force_stop").(bool) {
		  		_, err := conn.StopInstances(&updateReq1)
		  		if err != nil {
		  			return fmt.Errorf("error on stop instance , %s", err)
		  		}
		  	}
		  }


				if passwordUpdate {
					action := "ModifyInstanceType"
					logger.Debug(logger.ReqFormat, action, updateReq3)
					resp, err := conn.ModifyInstanceType(&updateReq3)
					if err != nil {
						return fmt.Errorf("error on updating  instance type, %s", err)
					}
					logger.Debug(logger.RespFormat, action, updateReq3, *resp)
					for _, v := range updatedAttributePassword {
						d.SetPartial(v)
					}
				}

			if initState == "active" {
				//判断初始开机状态,若开机则开机
				action := "RebootInstances"
				logger.Debug(logger.ReqFormat, action, updateReq1)
				resp, err := conn.RebootInstances(&updateReq1) //同步
				if err != nil {
					return fmt.Errorf("error on updating instance type, %s", err)
				}
				logger.Debug(logger.RespFormat, action, updateReq1, *resp)
			}
	*/
	d.Partial(false)
	return resourceKsyunInstanceRead(d, meta)
}

func resourceKsyunInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).kecconn
	//delete
	deleteReq := make(map[string]interface{})
	deleteReq["InstanceId.1"] = d.Id()
	if v, ok := d.GetOk("force_delete"); ok {
		deleteReq["ForceDelete"] = fmt.Sprintf("%v", v)
	}

	return resource.Retry(15*time.Minute, func() *resource.RetryError {
		readReq := make(map[string]interface{})
		readReq["InstanceId.1"] = d.Id()
		action := "DescribeInstances"
		logger.Debug(logger.ReqFormat, action, readReq)
		resp, err1 := conn.DescribeInstances(&readReq)
		logger.Debug(logger.AllFormat, action, readReq, *resp, err1)
		/*
			{
			    "Marker": 0,
			    "InstanceCount": 0,
			    "RequestId": "7c34a18b-b562-44f3-8ea9-a350e3afe649",
			    "InstancesSet": []
			}
		*/
		if err1 != nil && notFoundError(err1) {
			return nil
		}
		if err1 != nil {
			return resource.NonRetryableError(err1)
		}
		itemset, ok := (*resp)["InstancesSet"]
		if !ok {
			return nil
		}
		items, ok := itemset.([]interface{})
		if !ok || len(items) == 0 {
			return nil
		}
		state, ok := items[0].(map[string]interface{})["InstanceState"]
		if !ok {
			return nil
		}
		initState, ok := state.(map[string]interface{})["Name"].(string)
		if !ok {
			return nil
		}
		if strings.ToLower(initState) != "stopped" && strings.ToLower(initState) != "error" {
			action = "StopInstances"
			logger.Debug(logger.ReqFormat, action, readReq)
			resp, err := conn.StopInstances(&readReq) //同步
			logger.Debug(logger.AllFormat, action, readReq, *resp, err)
			if err1 != nil && notFoundError(err1) {
				return nil
			}
			if err != nil {
				return resource.RetryableError(err)
			}
			stateConf := &resource.StateChangeConf{
				Pending:    []string{statusPending},
				Target:     []string{"stopped"},
				Refresh:    instanceStateRefreshFunc(conn, d.Id(), []string{"stopped"}),
				Timeout:    d.Timeout(schema.TimeoutUpdate),
				Delay:      3 * time.Second,
				MinTimeout: 2 * time.Second,
			}
			if _, err = stateConf.WaitForState(); err != nil {
				return resource.RetryableError(err)
			}
		}

		action = "TerminateInstances"
		logger.Debug(logger.ReqFormat, action, deleteReq)
		resp, err2 := conn.TerminateInstances(&deleteReq)
		logger.Debug(logger.AllFormat, action, deleteReq, *resp, err2)
		if err2 == nil || notFoundError(err2) {
			return nil
		}
		if err2 != nil && inUseError(err2) {
			return resource.RetryableError(err2)
		}
		//查询验证
		readReq = make(map[string]interface{})
		readReq["InstanceId.1"] = d.Id()
		action = "DescribeInstances"
		logger.Debug(logger.ReqFormat, action, readReq)
		resp, err := conn.DescribeInstances(&readReq)
		logger.Debug(logger.AllFormat, action, readReq, *resp)
		if err != nil && notFoundError(err) {
			return nil
		}
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error on  reading kec when delete %q, %s", d.Id(), err))
		}

		itemset, ok = (*resp)["InstancesSet"]
		if !ok {
			return nil
		}
		items, ok = itemset.([]interface{})
		if !ok || len(items) == 0 {
			return nil
		}
		return resource.RetryableError(err2)
	})
}
func instanceStateRefreshFunc(client *kec.Kec, instanceId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := map[string]interface{}{"InstanceId.1": instanceId}
		resp, err := client.DescribeInstances(&req)
		if err != nil {
			return nil, "", err
		}
		itemset, ok := (*resp)["InstancesSet"]
		items, ok := itemset.([]interface{})
		if !ok || len(items) == 0 {
			return nil, "", fmt.Errorf("no instance set get")
		}
		item, ok1 := items[0].(map[string]interface{})
		if !ok1 {
			return nil, "", fmt.Errorf("no instance set get")
		}
		instanceState, ok2 := item["InstanceState"]
		if !ok2 {
			return nil, "", fmt.Errorf("no instance state get")
		}
		instancestate, ok3 := instanceState.(map[string]interface{})
		if !ok3 {
			return nil, "", fmt.Errorf("no instance state get")
		}
		state := strings.ToLower(instancestate["Name"].(string))
		for k, v := range target {
			if v == state {
				return resp, state, nil
			}
			if k == len(target)-1 {
				state = statusPending
			}
		}
		return resp, state, nil
	}
}

func instanceStateRefreshForReinstallFunc(client *kec.Kec, instanceId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := map[string]interface{}{"InstanceId.1": instanceId}
		resp, err := client.DescribeInstances(&req)
		if err != nil {
			return nil, "", err
		}
		itemset, ok := (*resp)["InstancesSet"]
		items, ok := itemset.([]interface{})
		if !ok || len(items) == 0 {
			return nil, "", fmt.Errorf("no instance set get")
		}
		item, ok1 := items[0].(map[string]interface{})
		if !ok1 {
			return nil, "", fmt.Errorf("no instance set get")
		}
		instanceState, ok2 := item["InstanceState"]
		if !ok2 {
			return nil, "", fmt.Errorf("no instance state get")
		}
		instancestate, ok3 := instanceState.(map[string]interface{})
		if !ok3 {
			return nil, "", fmt.Errorf("no instance state get")
		}
		state := strings.ToLower(instancestate["Name"].(string))
		if state == "stopped" {
			return nil, "", fmt.Errorf("instance restart error")
		}
		for k, v := range target {
			if v == state {
				return resp, state, nil
			}
			if k == len(target)-1 {
				state = statusPending
			}
		}
		return resp, state, nil
	}
}

func instanceStateRefreshForCreateFunc(client *kec.Kec, instanceId string, target []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		req := map[string]interface{}{"InstanceId.1": instanceId}
		action := "DescribeInstances"
		logger.Debug(logger.ReqFormat, action, req)
		resp, err := client.DescribeInstances(&req)
		if err != nil {
			return nil, "", err
		}
		logger.Debug(logger.RespFormat, action, req, *resp)
		itemset, ok := (*resp)["InstancesSet"]
		items, ok := itemset.([]interface{})
		if !ok || len(items) == 0 {
			return nil, "", fmt.Errorf("no instance set get")
		}
		item, ok1 := items[0].(map[string]interface{})
		if !ok1 {
			return nil, "", fmt.Errorf("no instance set get")
		}
		instanceState, ok2 := item["InstanceState"]
		if !ok2 {
			return nil, "", fmt.Errorf("no instance state get")
		}
		instancestate, ok3 := instanceState.(map[string]interface{})
		if !ok3 {
			return nil, "", fmt.Errorf("no instance state get")
		}
		state := strings.ToLower(instancestate["Name"].(string))
		if state == "error" {
			return nil, "", fmt.Errorf("instance create error")
		}
		for k, v := range target {
			if v == state {
				//resp cant't be null else will try
				return resp, state, nil
			}
			if k == len(target)-1 {
				state = statusPending
			}
		}
		return resp, state, nil
	}
}
