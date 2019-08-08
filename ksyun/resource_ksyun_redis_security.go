package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
)

// instance security rule
func resourceRedisSecurityRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRedisSecurityRuleCreate,
		Delete: resourceRedisSecurityRuleDelete,
		Update: resourceRedisSecurityRuleUpdate,
		Read:   resourceRedisSecurityRuleRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"available_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cache_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"security_rule_id": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"cidr": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceRedisSecurityRuleCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		rules []interface{}
		resp *map[string]interface{}
		err error
		ok bool
	)

	conn := meta.(*KsyunClient).kcsv1conn
	createReq := make(map[string]interface{})
	createReq["CacheId"] = d.Get("cache_id")
	if az, ok := d.GetOk("available_zone"); ok {
		createReq["AvailableZone"] = az
	}
	if rules, ok = d.Get("rules").([]interface{}); !ok {
		return fmt.Errorf("type of security_rule.rules must be array map")
	}
	var i int
	for _, rule := range rules {
		i = i + 1
		var r map[string]interface{}
		if r, ok = rule.(map[string]interface{}); !ok {
			return fmt.Errorf("type of security_group.rules.[%v] must be map", i)
		}
		createReq[fmt.Sprintf("%v%v", "SecurityRules.Cidr.", i)] = r["cidr"]
	}
	action := "SetCacheSecurityRules"
	logger.Debug(logger.ReqFormat, action, createReq)
	if resp, err = conn.SetCacheSecurityRules(&createReq); err != nil {
		return fmt.Errorf("error on set instance security rule: %s", err)
	}
	logger.Debug(logger.RespFormat, action, createReq, *resp)
	d.SetId(createReq["CacheId"].(string))
	resourceRedisSecurityRuleRead(d, meta)
	return nil
}

func resourceRedisSecurityRuleDelete(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err error
	)
	conn := meta.(*KsyunClient).kcsv1conn
	deleteReq := make(map[string]interface{})
	deleteReq["CacheId"] = d.Get("cache_id")
	if az, ok := d.GetOk("available_zone"); ok {
		deleteReq["AvailableZone"] = az
	}
	rules := d.Get("rules").([]interface{})
	action := "DeleteCacheSecurityRule"
	for _, rule := range rules {
		r := rule.(map[string]interface{})
		deleteReq["SecurityRuleId"] = fmt.Sprintf("%v", r["security_rule_id"].(float64))
		logger.Debug(logger.ReqFormat, action, deleteReq)
		if resp, err = conn.DeleteCacheSecurityRule(&deleteReq); err != nil {
			return fmt.Errorf("error on delete instance security rule: %s", err)
		}
		logger.Debug(logger.RespFormat, action, deleteReq, *resp)
	}
	return nil
}

func resourceRedisSecurityRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	var (
		deleteRulesResults 	[]float64
		addRulesResults 	[]string
		resp *map[string]interface{}
		err error
	)
	d.Partial(true)
	defer d.Partial(false)
	updateReq := make(map[string]interface{})
	updateReq["CacheId"] = d.Get("cache_id")
	if az, ok := d.GetOk("available_zone"); ok {
		updateReq["AvailableZone"] = az
	}
	if d.HasChange("rules") {
		oldMc, newMc := d.GetChange("rules")
		oldRules := oldMc.([]interface{})
		newRules := newMc.([]interface{})
		for _, oldRule := range oldRules {
			oldR := oldRule.(map[string]interface{})
			exist := false
			for _, newRule := range newRules {
				newR := newRule.(map[string]interface{})
				if newR["cidr"] == oldR["cidr"] {
					exist = true
				}
			}
			if !exist {
				deleteRulesResults = append(deleteRulesResults, oldR["security_rule_id"].(float64))
			}
		}

		for _, newRule := range newRules {
			newR := newRule.(map[string]interface{})
			ip := newR["cidr"]
			exist := false
			for _, oldRule  := range oldRules {
				oldR := oldRule.(map[string]interface{})
				if oldR["cidr"] == ip {
					exist = true
				}
			}
			if !exist {
				addRulesResults = append(addRulesResults, ip.(string))
			}
		}
	}
	if len(addRulesResults) > 0 {
		var i int
		for _, rule := range addRulesResults {
			i = i + 1
			updateReq[fmt.Sprintf("%v%v", "SecurityRules.Cidr.", i)] = rule
		}
		conn := meta.(*KsyunClient).kcsv1conn
		action := "SetCacheSecurityRules"
		logger.Debug(logger.ReqFormat, action, updateReq)
		if resp, err = conn.SetCacheSecurityRules(&updateReq); err != nil {
			return fmt.Errorf("error on add instance security rule: %s", err)
		}
		logger.Debug(logger.RespFormat, action, updateReq, *resp)
	}
	if len(deleteRulesResults) > 0 {
		deleteRuleReq := make(map[string]interface{})
		deleteRuleReq["CacheId"] = d.Get("cache_id")
		for _, delRuleId := range deleteRulesResults {
			conn := meta.(*KsyunClient).kcsv1conn
			deleteRuleReq["SecurityRuleId"] = fmt.Sprintf("%v", delRuleId)
			if az, ok := d.GetOk("available_zone"); ok {
				deleteRuleReq["AvailableZone"] = az
			}
			action := "DeleteCacheSecurityRule"
			logger.Debug(logger.ReqFormat, action, deleteRuleReq)
			if resp, err = conn.DeleteCacheSecurityRule(&deleteRuleReq); err != nil {
				return fmt.Errorf("error on delete instance security rule: %s", err)
			}
			logger.Debug(logger.RespFormat, action, deleteRuleReq, *resp)
		}
	}
	resourceRedisSecurityRuleRead(d, meta)
	return nil
}

func resourceRedisSecurityRuleRead(d *schema.ResourceData, meta interface{}) error {
	var (
		resp *map[string]interface{}
		err error
	)
	conn := meta.(*KsyunClient).kcsv1conn
	readReq := make(map[string]interface{})
	readReq["CacheId"] = d.Get("cache_id")
	if az, ok := d.GetOk("available_zone"); ok {
		readReq["AvailableZone"] = az
	}
	action := "DescribeCacheSecurityRules"
	logger.Debug(logger.ReqFormat, action, readReq)
	if resp, err = conn.DescribeCacheSecurityRules(&readReq); err != nil {
		return fmt.Errorf("error on reading instance security rule %q, %s", d.Id(), err)
	}
	logger.Debug(logger.RespFormat, action, readReq, *resp)
	data := (*resp)["Data"].([]interface{})
	if len(data)  == 0 {
		logger.Info("instance security rule result size : 0")
		return nil
	}
	result := make(map[string]interface{})
	var rulesTemp []map[string]interface{}
	for _, v := range data {
		group := v.(map[string]interface{})
		rule := make(map[string]interface{})
		rule[Hump2Downline("securityRuleId")] = group["securityRuleId"]
		rule[Hump2Downline("cidr")] = group["cidr"]
		rulesTemp = append(rulesTemp, rule)
	}
	result["rules"] = rulesTemp
	for k, v := range result  {
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("error set data %v :%v", v, err)
		}
	}
	return nil
}
