package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ksc/ksc-sdk-go/service/sqlserver"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
	"time"
)

func resourceKsyunSqlServer() *schema.Resource{

	return &schema.Resource{
		Create: resourceKsyunSqlServerCreate,
		Update: resourceKsyunSqlServerUpdate,
		Read:   resourceKsyunSqlServerRead,
		Delete: resourceKsyunSqlServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dbinstanceidentifier":{
				Type:schema.TypeString,
				Optional:true,
			},
			"dbinstanceclass":{
				Type:schema.TypeString,
				Required:true,
			},
			"dbinstancename":{
				Type:schema.TypeString,
				Required:true,
			},
			"dbinstancetype":{
				Type:schema.TypeString,
				Required:true,
			},
			"engine":{
				Type:schema.TypeString,
				Required:true,
			},
			"engineversion":{
				Type:schema.TypeString,
				Required:true,
			},
			"masterusername":{
				Type:schema.TypeString,
				Required:true,
			},
			"masteruserpassword":{
				Type:schema.TypeString,
				Required:true,
			},
			"vpcid":{
				Type:schema.TypeString,
				Required:true,
			},
			"subnetid":{
				Type:schema.TypeString,
				Required:true,
			},
			"billtype":{
				Type:schema.TypeString,
				Required:true,
			},

			// oxoxoxox
			"sqlservers": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema:map[string]*schema.Schema{
						"dbinstanceclass": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema:map[string]*schema.Schema{
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
									"iops":{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"maxconn": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"mem":{
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
						"port":{
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
						"productwhat":{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"servicestarttime" :{
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
func resourceKsyunSqlServerCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).sqlserverconn
	var resp *map[string]interface{}
	createReq := make(map[string]interface{})
	var err error
	creates :=[]string{
		"DBInstanceClass",
		"DBInstanceName",
		"DBInstanceType",
		"Engine",
		"EngineVersion",
		"MasterUserName",
		"MasterUserPassword",
		"VpcId",
		"SubnetId",
		"BillType",
	}
	for _,v :=range creates{
		if v1,ok:=d.GetOk(strings.ToLower(v)); ok{
			createReq[v]=fmt.Sprintf("%v",v1)
		}
	}
	action :="CreateDBInstance"
	logger.Debug(logger.RespFormat,action,createReq)
	resp,err = conn.CreateDBInstance(&createReq)
	logger.Debug(logger.AllFormat,action,createReq,*resp,err)
	if err != nil{
		return fmt.Errorf("error on creating Instance(sqlserver): %s",err)
	}

	if resp != nil {
		bodyData := (*resp)["Data"].(map[string]interface{})
		instances := bodyData["Instances"].([]interface{})
		sqlserverInstance := instances[0].(map[string]interface{})
		instanceId := sqlserverInstance["DBInstanceIdentifier"].(string)
		logger.DebugInfo("~*~*~*~*~ DBInstanceIdentifier : %v",instanceId)
		d.SetId(instanceId)
	}
	stateConf := &resource.StateChangeConf{
		Pending:[]string{tCreatingStatus},
		Target:[]string{tActiveStatus,tFailedStatus,tDeletedStatus,tStopedStatus},
		Timeout:d.Timeout(schema.TimeoutCreate),
		Delay:10*time.Second,
		MinTimeout:10*time.Second,
		Refresh:sqlserverInstanceStateRefreshForCreate(conn,d.Id(),[]string{tCreatingStatus}),
	}
	_,err = stateConf.WaitForState()


	return resourceKsyunSqlServerRead(d,meta)
}

func sqlserverInstanceStateRefreshForCreate(client *sqlserver.Sqlserver, instanceId string, target []string) resource.StateRefreshFunc {
	return func() ( interface{}, string, error) {
		req:= map[string]interface{}{"DBInstanceIdentifier": instanceId}
		action := "DescribeDBInstances"
		logger.Debug(logger.ReqFormat,action,req)
		resp,err:=client.DescribeDBInstances(&req)
		logger.Debug(logger.AllFormat,action,req,*resp,err)
		if err != nil {
			return nil,"",err
		}
		bodyData := (*resp)["Data"].(map[string]interface{})
		instances := bodyData["Instances"].([]interface{})
		sqlserverInstance := instances[0].(map[string]interface{})
		state := sqlserverInstance["DBInstanceStatus"].(string)

		return resp,state,nil

	}
}

func resourceKsyunSqlServerRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).sqlserverconn
	req:= map[string]interface{}{"DBInstanceIdentifier": d.Id()}
	action := "DescribeDBInstances"
	logger.Debug(logger.ReqFormat,action,req)
	resp,err:=conn.DescribeDBInstances(&req)
	logger.Debug(logger.AllFormat,action,req,*resp,err)
	if err != nil {
		return fmt.Errorf("error on reading Instance(sqlserver) %q, %s",d.Id(),err)
	}


	bodyData,dataOk := (*resp)["Data"].(map[string]interface{})
	if !dataOk {
		return fmt.Errorf("error on reading Instance(sqlserver) body %q, %+v",d.Id(), (*resp)["Error"])
	}
	instances := bodyData["Instances"].([]interface{})



	sqlserverIds := make([]string,len(instances))
	sqlserverMap := make([]map[string]interface{}, len(instances))
	for k,instance := range instances {
		instanceInfo,_ := instance.(map[string]interface{})
		for k,v := range instanceInfo  {
			if k == "DBInstanceClass" {
				dbclass := v.(map[string]interface{})
				dbinstanceclass := make(map[string]interface{})
				for j,q := range dbclass {
					dbinstanceclass[strings.ToLower(j)] = q
				}
				wtf := make([]interface{},1)
				wtf[0] = dbinstanceclass
				instanceInfo["dbinstanceclass"] = wtf
				delete(instanceInfo,"DBInstanceClass")
			}else {
				delete(instanceInfo,k)
				instanceInfo[strings.ToLower(k)] = v
			}
		}
		sqlserverMap[k] = instanceInfo
		logger.DebugInfo(" converted ---- %+v ",  instanceInfo)


		sqlserverIds[k] = instanceInfo["dbinstanceidentifier"].(string)
		logger.DebugInfo("sqlserverIds fuck : %v",sqlserverIds)
	}

	logger.DebugInfo(" converted ---- %+v ",  sqlserverMap)
	dataSourceSqlserverSave(d,"sqlservers",sqlserverIds,sqlserverMap)
	//sqlserverInstance := instances[0].(map[string]interface{})
	//DBInstanceClass := sqlserverInstance["DBInstanceClass"].(map[string]interface{})
	//for k,v := range DBInstanceClass{
	//	sqlserverInstance["DBInstanceClass."+k] = v
	//}
	//
	//for k,v := range sqlserverInstance {
	//	if !sqlserverIncludeKeys[k] || sqlserverExcludeKeys[k] {
	//		d.Set(k,v)
	//	}
	//}
	//state := sqlserverInstance["DBInstanceStatus"].(string)
	
	return nil
}

func resourceKsyunSqlServerUpdate(d *schema.ResourceData, meta interface{}) error {
	// 关闭事务，允许部分属性被修改  d.Partial(true) d.Partial(false)

	return fmt.Errorf("error on updating instance , sqlserver is not support update")
}

func resourceKsyunSqlServerDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*KsyunClient).sqlserverconn
	deleteReq := make(map[string]interface{})
	deleteReq["DBInstanceIdentifier"]=d.Id()

	return resource.Retry(15*time.Minute,func() *resource.RetryError{
		readReq:= map[string]interface{}{"DBInstanceIdentifier": d.Id()}
		discribeAction := "DescribeInstances"
		logger.Debug(logger.ReqFormat, discribeAction, readReq)
		desResp, desErr := conn.DescribeDBInstances(&readReq)
		logger.Debug(logger.AllFormat, discribeAction, readReq, *desResp, desErr)

		if desErr != nil {
			if notFoundError(desErr) {
				return nil
			} else {
				return resource.NonRetryableError(desErr)
			}
		}

		bodyData := (*desResp)["Data"].(map[string]interface{})
		instances := bodyData["Instances"].([]interface{})
		sqlserverInstance := instances[0].(map[string]interface{})
		state := sqlserverInstance["DBInstanceStatus"].(string)

		if state != tDeletedStatus{
			deleteAction := "DeleteDBInstance"
			logger.Debug(logger.ReqFormat, deleteAction, deleteReq)
			deleteResp,deleteErr := conn.DeleteDBInstance(&deleteReq)
			logger.Debug(logger.AllFormat, deleteAction, deleteReq, *deleteResp, deleteErr)
			if deleteErr == nil || notFoundError(deleteErr) {
				return nil
			}
			if deleteErr !=nil {
				return resource.RetryableError(deleteErr)
			}

			logger.Debug(logger.ReqFormat, discribeAction, readReq)
			postDesResp, postDesErr := conn.DescribeDBInstances(&readReq)
			logger.Debug(logger.AllFormat, discribeAction, readReq, *postDesResp, postDesErr)

			if desErr != nil {
				if notFoundError(desErr) {
					return nil
				} else {
					return resource.NonRetryableError(fmt.Errorf("error on  reading kec when delete %q, %s", d.Id(), desErr))
				}
			}
		}

		return resource.RetryableError(desErr)
	})
}
