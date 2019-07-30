package ksyun

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-ksyun/logger"
	"strings"
)

//将得到的map[string]interface{}类型的值转换为可赋值给d的类型。只支持[k,v]中的v为基本数据类型
func GetSubDByRep(data interface{}, include, exclude map[string]bool) []interface{} {
	ma, ok := data.(map[string]interface{})
	if !ok {
		return nil
	}
	subD := make(map[string]interface{})
	for k, v := range ma {
		if exclude[k] || !include[k] {
			continue
		}
		subD[Hump2Downline(k)] = v
	}
	return []interface{}{subD}
}

//sdk resp->terraform d
//将获取的list类型的值转换为可赋值给d的类型。
//遍历[]interface{}数据（其中interface{}可以转换为map[string]interface{}),将驼峰型字段名转为下划线。
//用于 d.Set（ key，value）中的value，且value为TypeList类型
func GetSubSliceDByRep(items []interface{}, include /*,exclude*/ map[string]bool) []map[string]interface{} {
	datas := []map[string]interface{}{}
	for _, v := range items {
		data := map[string]interface{}{}
		vv, _ := v.(map[string]interface{})
		for key, value := range vv {
			//此处不判断，需在后面对非基本类型单独处理
			if /*exclude[key]||*/ !include[key] {
				continue //不判断，当有新字段加入时无法解析
			}
			data[Hump2Downline(key)] = value
		}
		datas = append(datas, data)
	}
	return datas
}

//将获取的map类型的值转换为可赋值给d的类型。只支持Map中的elem为map[string]interface{}，且[k,v]中的v为基本数据类型
func GetSubStructDByRep(datas interface{}, exclude map[string]bool) map[string]interface{} {

	subStruct := map[string]interface{}{}
	items, ok := datas.(map[string]interface{})
	if !ok {
		return subStruct
	}
	for k, v := range items {
		if exclude[k] {
			continue
		}
		subStruct[Hump2Downline(k)] = v
	}
	return subStruct
}

//将得到的map[string]interface{}类型的值赋值给d。只支持[k,v]中的v为基本数据类型
func SetDByRespV1(d *schema.ResourceData, m interface{}, exclud map[string]bool) map[string]interface{} {
	ma, ok := m.(map[string]interface{})
	fmt.Println("ok:", ok)
	mre := make(map[string]interface{}, 0)
	if !ok {
		return mre
	}
	for k, v := range ma {
		if exclud[k] {
			if mm, ok := v.(map[string]interface{}); ok {
				mre[k] = mm
			} else {
				mre[k] = v
			}
			continue
		}
		d.Set(Hump2Downline(k), v)
	}
	return mre
}

//将得到的map[string]interface{}类型的值赋值给d。只支持[k,v]中的v为基本数据类型
//Include 是已有字段，防止日后新加字段无法解析。exclude是排除字段，非基本类型字段需排除做特殊处理。
func SetDByResp(d *schema.ResourceData, m interface{}, includ, exclude map[string]bool) map[string]interface{} {
	mre := make(map[string]interface{}, 0)
	ma, ok := m.(map[string]interface{})
	if !ok {
		return mre
	}
	for k, v := range ma {
		if !includ[k] || exclude[k] {
			if mm, ok := v.(map[string]interface{}); ok {
				mre[k] = mm
			} else {
				mre[k] = v
			}
			continue
		}

		d.Set(Hump2Downline(k), v)
	}
	return mre
}

//简单驼峰转下划线 未对连写大写字母做判断进行特殊处理。
//即aDDCC 转为a_d_d_c_c 而非a_ddc_c
func Hump2Downline(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return s
	}
	var s1 string
	if len(s) == 1 {
		s1 = strings.ToLower(s[:1])
		return s1
	}
	for k, v := range s {
		if k == 0 {
			s1 = strings.ToLower(s[0:1])
			continue
		}
		if v >= 65 && v <= 90 {
			v1 := "_" + strings.ToLower(s[k:k+1])
			s1 = s1 + v1
		} else {
			s1 = s1 + s[k:k+1]
		}
	}
	return s1
}

//简单下划线转驼峰
func Downline2Hump(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return s
	}
	var s1 []string
	ss := strings.Split(s, "_")
	for _, v := range ss {
		vv := strings.ToUpper(v[:1]) + v[1:]
		s1 = append(s1, vv)
	}
	return strings.Join(s1, "")
}

//flattern struct 用于创建时，结构体类型的入参转换为map型 ,不带前缀拼接
func FlatternStruct(v interface{}, req *map[string]interface{}) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				vv := Downline2Hump(k2)
				(*req)[vv] = fmt.Sprintf("%v", v2)
			}
		}
	}
}

//flattern struct 用于创建时，结构体类型的入参转换为map型,带前缀拼接
func FlatternStructPrefix(v interface{}, req *map[string]interface{}, prex string) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				vv := Downline2Hump(k2)
				(*req)[fmt.Sprintf("%s.%s", prex, vv)] = fmt.Sprintf("%v", v2)
			}
		}
	}
}

//FlatternStructSlicePrefix 用于创建时，结构体切片类型的入参转换为map型 ,【
func FlatternStructSlicePrefix(values interface{}, req *map[string]interface{}, prex string) {
	v, _ := values.([]interface{})
	for k1, v1 := range v {
		vv := v1.(map[string]interface{})
		for k2, v2 := range vv {
			vv := Downline2Hump(k2)
			(*req)[fmt.Sprintf("%s.%d.%s", prex, k1+1, vv)] = fmt.Sprintf("%v", v2)
		}
	}
}

//ConvertFilterStruct  用于创建时，结构体类型的入参转换为map型 ,不带前缀拼接
func ConvertFilterStruct(v interface{}, req *map[string]interface{}) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				vv := strings.ReplaceAll(k2, "_", "-")
				(*req)[vv] = fmt.Sprintf("%v", v2)
			}
		}
	}
}

//ConvertFilterStruct  用于创建时，结构体类型的入参转换为map型,带前缀拼接
func ConvertFilterStructPrefix(v interface{}, req *map[string]interface{}, prex string) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			if v1[0] == nil {
				return
			}
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				vv := strings.ReplaceAll(k2, "_", "-")
				(*req)[fmt.Sprintf("%s.%s", prex, vv)] = v2
			}
		}
	}
}

/*
func ConvertFilterStructStructPrefix(v interface{}, req *map[string]interface{}, prex string) {
	if v1, ok1 := v.([]interface{}); ok1 {
		if len(v1) > 0 {
			if v1[0] == nil {
				return
			}
			vv := v1[0].(map[string]interface{})
			for k2, v2 := range vv {
				vv := strings.ReplaceAll(k2, "-", "_")
				v3, ok3 := v2.([]string)
				if !ok3 || len(v3) == 0 {
					(*req)[fmt.Sprintf("%s.%s", prex, vv)] = fmt.Sprintf("%v", v2)
				}
				(*req)[fmt.Sprintf("%s.%s", prex, vv)] = fmt.Sprintf("%v", v3[0])

			}
		}
	}
}

*/
func dataSourceKscSave(d *schema.ResourceData, dataKey string, ids []string, datas []map[string]interface{}) error {

	d.SetId(hashStringArray(ids))
	d.Set("total_count", len(datas))

	if err := d.Set(dataKey, datas); err != nil {
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		writeToFile(outputFile.(string), datas)
	}

	return nil
}
func dataSourceKscSaveSlice(d *schema.ResourceData, dataKey string, ids []string, datas []string) error {

	d.SetId(hashStringArray(ids))
	d.Set("total_count", len(datas))

	if err := d.Set(dataKey, datas); err != nil {
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		writeToFile(outputFile.(string), datas)
	}

	return nil
}

func dataSourceSqlserverSave(d *schema.ResourceData, dataKey string, ids []string, datas []map[string]interface{}) error {

	if len(ids) == 1 {
		d.SetId(ids[0])
	} else {
		d.SetId(strings.Join(ids,","))
	}

	d.Set("total_count", len(datas))

	logger.DebugInfo("$$$$$$$$$datasdatasdatas$$$$$$$$ %+v",datas)
	if err := d.Set(dataKey, datas); err != nil {
		logger.DebugInfo("$$$$$$$$$omg$$$$$$$$ %+v",err)
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}
	logger.DebugInfo("$$$$$$$$$fuckfuckfuck$$$$$$$$ %+v",datas)
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		logger.DebugInfo(" ------------ %+v", outputFile)
		writeToFile(outputFile.(string), datas)
	} else {
		logger.DebugInfo(" !!!!!!!!!!! %+v",  outputFile)
	}

	return nil
}

func dataSourceSqlserverDataSave(d *schema.ResourceData, dataKey string, ids []string, datas []map[string]interface{}) error {

	if len(ids) == 1 {
		d.SetId(ids[0])
	} else {
		d.SetId(strings.Join(ids,","))
	}

	d.Set("total_count", len(datas))

	logger.DebugInfo("$$$$$$$$$datasdatasdatas$$$$$$$$ %+v",datas)
	if err := d.Set(dataKey, datas); err != nil {
		logger.DebugInfo("$$$$$$$$$omg$$$$$$$$ %+v",err)
		return fmt.Errorf("error set datas %v :%v", datas, err)
	}
	logger.DebugInfo("$$$$$$$$$fuckfuckfuck$$$$$$$$ %+v",datas)
	if outputFile, ok := d.GetOk("output_file"); ok && outputFile.(string) != "" {
		logger.DebugInfo(" ------------ %+v", outputFile)
		writeToFile(outputFile.(string)+"_data", datas)
	} else {
		logger.DebugInfo(" !!!!!!!!!!! %+v",  outputFile)
	}

	return nil
}