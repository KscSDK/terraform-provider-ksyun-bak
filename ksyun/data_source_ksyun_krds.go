package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
)

func dataSourceKsyunKrds() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceKsyunSqlServerRead,
		Schema: map[string]*schema.Schema{
			//DBInstanceStatus  DBInstanceType DBInstanceIdentifier Keyword ExpiryDateLessThan Marker MaxRecords
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dbinstancestatus": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dbinstanceidentifier": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sqlservers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dbinstanceclass": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"vcpus": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"disk": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"ram": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"iops": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"maxconn": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"mem": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"dbinstanceidentifier": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dbinstancename": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dbinstancestatus": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dbinstancetype": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"groupid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vip": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"engineversion": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"instancecreatetime": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"masterusername": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vpcid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"subnetid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"publiclyaccessible": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"readreplicadbinstanceidentifiers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"billtype": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ordertype": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ordersource": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"masteravailabilityzone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"slaveavailabilityzone": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"multiavailabilityzone": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"productid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"orderuse": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"projectid": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"projectname": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"billtypeid": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"dbparametergroupid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"datastoreversionid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"diskused": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"preferredbackuptime": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"productwhat": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"servicestarttime": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"suborderid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"audit": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}

}

func dataSourceKsyunKrdsRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).sqlserverconn
	desReq := make(map[string]interface{})
	des := []string{
		"DBInstanceStatus",
		"DBInstanceType",
		"DBInstanceIdentifier",
		"Keyword",
		"ExpiryDateLessThan",
		"Marker",
		"MaxRecords",
	}
	for _, v := range des {
		if v1, ok := d.GetOk(strings.ToLower(v)); ok {
			desReq[v] = fmt.Sprintf("%v", v1)
		}
	}
	action := "DescribeDBInstances"
	logger.Debug(logger.ReqFormat, action, desReq)
	resp, err := conn.DescribeDBInstances(&desReq)
	logger.Debug(logger.AllFormat, action, desReq, *resp, err)
	if err != nil {
		return fmt.Errorf("error on reading Instance(sqlserver)  %s", err)
	}

	bodyData, dataOk := (*resp)["Data"].(map[string]interface{})
	if !dataOk {
		return fmt.Errorf("error on reading Instance(sqlserver) body %+v", (*resp)["Error"])
	}
	instances := bodyData["Instances"].([]interface{})
	if len(instances) == 0 {
		return fmt.Errorf("error on reading Instance(sqlserver) body %+v", (*resp))
	}

	logger.Debug("sqlserver start get ids ", action, bodyData)
	sqlserverIds := make([]string, len(instances))
	sqlserverMap := make([]map[string]interface{}, len(instances))
	for k, instance := range instances {
		instanceInfo, _ := instance.(map[string]interface{})
		for k, v := range instanceInfo {
			if k == "DBInstanceClass" {
				dbclass := v.(map[string]interface{})
				dbinstanceclass := make(map[string]interface{})
				for j, q := range dbclass {
					dbinstanceclass[strings.ToLower(j)] = q
				}
				wtf := make([]interface{}, 1)
				wtf[0] = dbinstanceclass
				instanceInfo["dbinstanceclass"] = wtf
				delete(instanceInfo, "DBInstanceClass")
			} else {
				delete(instanceInfo, k)
				instanceInfo[strings.ToLower(k)] = v
			}
		}
		sqlserverMap[k] = instanceInfo
		logger.DebugInfo(" converted ---- %+v ", instanceInfo)

		sqlserverIds[k] = instanceInfo["dbinstanceidentifier"].(string)
	}

	logger.DebugInfo(" converted ---- %+v ", sqlserverMap)
	dataSourceSqlserverDataSave(d, "sqlservers", sqlserverIds, sqlserverMap)

	return nil
}

var krdsIncludeKeys = map[string]bool{
	"DBInstanceIdentifier": true,
	"PreferredBackupTime":  true,
	"DBInstanceName":       true,
	"DBInstanceStatus":     true,
	"DBInstanceType":       true,
	"DBParameterGroupId":   true,
	"GroupId":              true,
	"Vip":                  true,
	"Port":                 true,
	"Engine":               true,
	"EngineVersion":        true,
	"InstanceCreateTime":   true,
	//	"MasterUserName": true,
	//	"DatastoreVersionId": true,
	"VpcId":              true,
	"SubnetId":           true,
	"PubliclyAccessible": true,
	//	"BillType": true,
	//	"OrderType": true,
	"MultiAvailabilityZone": true,
	"DiskUsed":              true,
	//	"ProductId": true,
	//	"ProductWhat": true,
	"ProjectId":        true,
	"ProjectName":      true,
	"Region":           true,
	"ServiceStartTime": true,
	//	"SubOrderId": true,
	"Audit":                            true,
	"ReadReplicaDBInstanceIdentifiers": true,
	//	"BillTypeId": true,
	"DBInstanceClass.Id":      true,
	"DBInstanceClass.Iops":    true,
	"DBInstanceClass.Vcpus":   true,
	"DBInstanceClass.Disk":    true,
	"DBInstanceClass.Ram":     true,
	"DBInstanceClass.Mem":     true,
	"DBInstanceClass.MaxConn": true,
}
var krdsExcludeKeys = map[string]bool{}
