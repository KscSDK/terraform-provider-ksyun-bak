package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"log"
	"strings"
)

func resourceKsyunNetworkInterface() *schema.Resource {
	return &schema.Resource{
		//	Create: resourceKsyunNetworkInterfaceCreate,
		Read:   resourceKsyunNetworkInterfaceRead,
		Update: resourceKsyunNetworkInterfaceUpdate,
		Delete: resourceKsyunNetworkInterfaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_group_id": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"d_n_s1": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"d_n_s2": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vpc_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceKsyunNetworkInterfaceRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).vpcconn
	readReq := make(map[string]interface{})
	id := d.Id()
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return fmt.Errorf("id shoule like instance_id:networkinterface_id")
	}
	readReq["NetworkInterfaceId.1"] = ids[1]
	action := "DescribeNetworkInterfaces"
	logger.Debug(logger.ReqFormat, action, readReq)
	resp, err := conn.DescribeNetworkInterfaces(&readReq)
	if err != nil {
		return fmt.Errorf("error on reading Instance %q, %s", d.Id(), err)
	}
	logger.Debug(logger.AllFormat, action, readReq, *resp, err)
	itemset, ok := (*resp)["NetworkInterfaceSet"]
	items, ok := itemset.([]interface{})
	if !ok || len(items) == 0 {
		d.SetId("")
		return nil
	}
	excludesKeys := map[string]bool{
		"SecurityGroupSet": true,
	}
	excludes := SetDByResp(d, items[0], networkInterfaceKeys, excludesKeys)
	log.Printf("excludes:%v", excludes)
	if sg, ok := excludes["SecurityGroupSet"]; ok {
		itemSetSub := GetSubSliceDByRep(sg.([]interface{}), kecSecurityGroupKeys)
		if len(itemSetSub) != 0 {
			var itemSetSlice []string
			for _, v := range itemSetSub {
				for k1, v1 := range v {
					log.Printf("k1:%v;v1:%v", k1, v1)
					if k1 == "security_group_id" {
						itemSetSlice = append(itemSetSlice, fmt.Sprintf("%v", v1))
					}
				}
			}
			log.Printf("itemSetSlice:%v", itemSetSlice)
			d.Set("security_group_id", itemSetSlice)
		}
	}
	return nil
}

func resourceKsyunNetworkInterfaceUpdate(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return fmt.Errorf("id shoule like instance_id:networkinterface_id")
	}
	conn := m.(*KsyunClient).kecconn
	// Enable partial attribute modification
	d.Partial(true)
	// Whether the representative has any modifications
	attributeUpdate := false
	updateReq := make(map[string]interface{})
	updateReq["InstanceId"] = ids[0]
	updateReq["NetworkInterfaceId"] = ids[1]
	updateReq["SubnetId"] = fmt.Sprintf("%v", d.Get("subnet_id"))
	allAttributes := []string{
		"private_ip_address",
		"d_n_s1",
		"d_n_s2",
	}
	var updates []string
	for _, v := range allAttributes {
		if d.HasChange(v) {
			attributeUpdate = true
			updates = append(updates, v)
		}
	}
	if d.HasChange("subnet_id") && !d.IsNewResource() {
		attributeUpdate = true
	}
	if !attributeUpdate {
		return nil
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		securityGroupIds := SchemaSetToStringSlice(v)
		for k, v := range securityGroupIds {
			updateReq[fmt.Sprintf("SecurityGroupId.%v", k+1)] = v
		}
	}
	attributeUpdate = true
	for _, v := range updates {
		if v1, ok := d.GetOk(v); ok {
			updateReq[Downline2Hump(v)] = fmt.Sprintf("%v", v1)
		}
	}
	action := "ModifyNetworkInterfaceAttribute"
	logger.Debug(logger.ReqFormat, action, updateReq)
	resp, err := conn.ModifyNetworkInterfaceAttribute(&updateReq)
	logger.Debug(logger.AllFormat, action, updateReq, *resp, err)
	if err != nil {
		return fmt.Errorf("update NetworkInterface (%v)error:%v", updateReq, err)
	}
	result, ok := (*resp)["Return"]
	if !ok || fmt.Sprintf("%v", result) != "true" {
		return fmt.Errorf("update NetworkInterface (%v)error:%v", updateReq, result)
	}
	d.SetPartial("subnet_id")
	d.SetPartial("security_group_id")
	for _, v := range updates {
		d.SetPartial(v)
	}
	d.Partial(false)
	return resourceKsyunNetworkInterfaceRead(d, m)
}
func resourceKsyunNetworkInterfaceDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()
	ids := strings.Split(id, ":")
	if len(ids) != 2 {
		return fmt.Errorf("id shoule like instance_id:networkinterface_id")
	}
	deleteReq := make(map[string]interface{})
	deleteReq["InstanceId"] = ids[0]
	deleteReq["NetworkInterfaceId"] = ids[1]
	conn := m.(*KsyunClient).kecconn
	action := "DetachNetworkInterface"
	logger.Debug(logger.ReqFormat, action, deleteReq)
	resp, err2 := conn.DetachNetworkInterface(&deleteReq)
	logger.Debug(logger.AllFormat, action, deleteReq, *resp, err2)
	if err2 == nil || notFoundError(err2) {
		return nil
	}
	return fmt.Errorf("fail when delete instance :%v", err2)
}
