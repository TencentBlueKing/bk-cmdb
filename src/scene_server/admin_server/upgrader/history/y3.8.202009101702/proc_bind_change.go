/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package y3_8_202009101702

import (
	"context"
	"fmt"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

func addProcBindInfo(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	bindIPAttrFilter := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: "bind_info",
	}
	cnt, err := db.Table(common.BKTableNameObjAttDes).Find(bindIPAttrFilter).Count(ctx)
	if err != nil {
		return err
	}
	if cnt != 0 {
		blog.Infof("add process attribute bind_info error.")
		return nil
	}

	// add module attribute field
	newAttributeID, err := db.NextSequence(ctx, common.BKTableNameObjAttDes)
	if err != nil {
		blog.Errorf("addProcBindIP failed, NextSequence failed, err: %s", err.Error())
		return fmt.Errorf("NextSequence failed, err: %s", err.Error())
	}

	bindIPAttrIdxFilter := map[string]interface{}{
		common.BKObjIDField: common.BKInnerObjIDProc,
	}
	sort := common.BKPropertyIndexField + ":-1"
	procAttr := &Attribute{}
	if err := db.Table(common.BKTableNameObjAttDes).Find(bindIPAttrIdxFilter).Sort(sort).One(ctx, procAttr); err != nil {
		blog.Errorf("addProcBindIP failed, find proc max property index id failed, filter: %s err: %s", bindIPAttrIdxFilter, err.Error())
		return err
	}

	nowTime := metadata.Now()
	procBindIPAttr := &Attribute{
		ID:            int64(newAttributeID),
		OwnerID:       conf.OwnerID,
		ObjectID:      common.BKInnerObjIDProc,
		PropertyID:    "bind_info",
		PropertyName:  "bind_info",
		PropertyGroup: "proc_port",
		PropertyIndex: procAttr.PropertyIndex + 1,
		Placeholder:   "process bind port information",
		IsEditable:    true,
		IsPre:         true,
		IsRequired:    false,
		PropertyType:  common.FieldTypeTable,
		Option:        getSubAttr(),
		Creator:       conf.User,
		CreateTime:    &nowTime,
		LastTime:      &nowTime,
	}

	return db.Table(common.BKTableNameObjAttDes).Insert(ctx, procBindIPAttr)
}

func migrateProcBindInfo(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	// 已经处理过最大的进程id
	maxProcID := int64(0)
	for {
		procs := make([]dbProcess, 0)
		filter := map[string]interface{}{
			common.BKProcIDField:  map[string]interface{}{common.BKDBGT: maxProcID},
			common.BKProcBindInfo: map[string]interface{}{common.BKDBExists: false},
		}
		if err := db.Table(common.BKTableNameBaseProcess).Find(filter).Limit(500).All(ctx, &procs); err != nil {
			return err
		}
		blog.InfoJSON("start process  bind info. process id start: %d", maxProcID)

		for _, proc := range procs {
			if proc.ProcessID > maxProcID {
				maxProcID = proc.ProcessID
			}
			bindInfoArr := make([]map[string]interface{}, 0)
			if proc.Port != nil {
				portArr := strings.Split(*proc.Port, ",")
				for idx, port := range portArr {
					bindInfoArr = append(bindInfoArr, map[string]interface{}{
						"template_row_id": idx + 1,
						"ip":              proc.BindIP,
						"port":            port,
						"protocol":        proc.Protocol,
						"enable":          proc.PortEnable,
					})
				}
			} else {
				bindInfoArr = append(bindInfoArr, map[string]interface{}{
					"template_row_id": 1,
					"ip":              proc.BindIP,
					"port":            "",
					"protocol":        proc.Protocol,
					"enable":          proc.PortEnable,
				})
			}

			updateFilter := map[string]interface{}{
				common.BKProcIDField: proc.ProcessID,
				//common.BKProcBindInfo: map[string]interface{}{common.BKDBExists: false},
			}
			doc := map[string]interface{}{
				common.BKProcBindInfo: bindInfoArr,
			}
			if err := db.Table(common.BKTableNameBaseProcess).Update(ctx, updateFilter, doc); err != nil {
				return err
			}

		}
		if len(procs) == 0 {
			break
		}
	}

	return nil
}

func migrateProcTempBindInfo(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	// 已经处理过最大的id
	maxID := int64(0)
	for {
		procTemps := make([]dbProcessTemplate, 0)
		filter := map[string]interface{}{
			common.BKFieldID:                    map[string]interface{}{common.BKDBGT: maxID},
			"property." + common.BKProcBindInfo: map[string]interface{}{common.BKDBExists: false},
		}
		if err := db.Table(common.BKTableNameProcessTemplate).Find(filter).Limit(500).All(ctx, &procTemps); err != nil {
			return err
		}
		blog.InfoJSON("start process template bind info. id start: %d", maxID)

		for _, procTemp := range procTemps {
			if procTemp.ID > maxID {
				maxID = procTemp.ID
			}

			// 保证数据库的数据不是null
			bindInfoValue := make([]map[string]interface{}, 0)
			if procTemp.Property != nil {
				if procTemp.Property.Port != nil && procTemp.Property.Port.Value != nil {
					portArr := strings.Split(*procTemp.Property.Port.Value, ",")
					for idx := range portArr {
						bindInfoValue = append(bindInfoValue, map[string]interface{}{
							// 老版本只能有一行数据
							"row_id": idx + 1,
							"ip":     procTemp.Property.BindIP,
							"port": PropertyString{
								AsDefaultValue: procTemp.Property.Port.AsDefaultValue,
								Value:          &portArr[idx],
							},
							"protocol": procTemp.Property.Protocol,
							"enable":   procTemp.Property.Enable,
						})
					}

				} else {
					bindInfoValue = append(bindInfoValue, map[string]interface{}{
						// 老版本只能有一行数据
						"row_id":   1,
						"ip":       procTemp.Property.BindIP,
						"port":     procTemp.Property.Port,
						"protocol": procTemp.Property.Protocol,
						"enable":   procTemp.Property.Enable,
					})
				}

			}
			updateFilter := map[string]interface{}{
				common.BKFieldID: procTemp.ID,
			}
			doc := map[string]interface{}{
				"property." + common.BKProcBindInfo: map[string]interface{}{
					"as_default_value": true,
					"value":            bindInfoValue,
				},
			}
			if err := db.Table(common.BKTableNameProcessTemplate).Update(ctx, updateFilter, doc); err != nil {
				return err
			}

		}
		if len(procTemps) == 0 {
			break
		}

	}

	return nil
}

func clearProcAttrAndGroup(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {
	delPropertyID := []string{common.BKProcGatewayIP, common.BKProcGatewayPort, common.BKProcGatewayProtocol, common.BKProcGatewayCity, common.BKBindIP, common.BKPort, common.BKProtocol, common.BKProcPortEnable}

	delProcAttr := map[string]interface{}{
		common.BKObjIDField:      common.BKInnerObjIDProc,
		common.BKPropertyIDField: map[string]interface{}{common.BKDBIN: delPropertyID},
	}

	if err := db.Table(common.BKTableNameObjAttDes).Delete(ctx, delProcAttr); err != nil {
		blog.ErrorJSON("clearProcAttrAndGroup failed, delete attribute, filter:%s err: %s", delProcAttr, err.Error())
		return err
	}

	proxyGroupAttrFilter := map[string]interface{}{
		common.BKObjIDField:         common.BKInnerObjIDProc,
		common.BKPropertyGroupField: "network_proxy",
	}
	cnt, err := db.Table(common.BKTableNameObjAttDes).Find(proxyGroupAttrFilter).Count(ctx)
	if err != nil {
		blog.ErrorJSON("clearProcAttrAndGroup failed, find network proxy  attribute, filter:%s err: %s", proxyGroupAttrFilter, err.Error())
		return err
	}
	if cnt > 0 {
		return nil
	}
	delProxyGroupAttrFilter := map[string]interface{}{
		common.BKObjIDField:         common.BKInnerObjIDProc,
		common.BKPropertyGroupField: "network_proxy",
	}
	if err := db.Table(common.BKTableNamePropertyGroup).Delete(ctx, delProxyGroupAttrFilter); err != nil {
		blog.ErrorJSON("clearProcAttrAndGroup failed, find network proxy  attribute, filter:%s err: %s", delProxyGroupAttrFilter, err.Error())
		return err
	}

	return nil

}

func getSubAttr() []SubAttriubte {

	return []SubAttriubte{
		SubAttriubte{
			PropertyID:    "ip",
			PropertyName:  "IP",
			Placeholder:   "bind ip",
			IsEditable:    true,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        "^([0-9]{1,3}\\.){3}[0-9]{1,3}$",
			PropertyGroup: common.BKProcBindInfo,
			IsRequired:    true,
		},
		SubAttriubte{
			PropertyID:    "port",
			PropertyName:  "Port",
			Placeholder:   "single port: 8080, </br>multiple consecutive ports: 8080-8089</br> multiple discontinuous ports: 8080-8089.",
			IsEditable:    true,
			PropertyType:  common.FieldTypeSingleChar,
			Option:        "^(((([1-9][0-9]{0,3})|([1-5][0-9]{4})|(6[0-4][0-9]{3})|(65[0-4][0-9]{2})|(655[0-2][0-9])|(6553[0-5]))-(([1-9][0-9]{0,3})|([1-5][0-9]{4})|(6[0-4][0-9]{3})|(65[0-4][0-9]{2})|(655[0-2][0-9])|(6553[0-5])))|((([1-9][0-9]{0,3})|([1-5][0-9]{4})|(6[0-4][0-9]{3})|(65[0-4][0-9]{2})|(655[0-2][0-9])|(6553[0-5]))))$",
			PropertyGroup: common.BKProcBindInfo,
			IsRequired:    true,
		},
		SubAttriubte{
			PropertyID:    "protocol",
			PropertyName:  "Protocol",
			Placeholder:   "service use protocol,",
			IsEditable:    true,
			PropertyType:  common.FieldTypeEnum,
			Option:        []metadata.EnumVal{{ID: "1", Name: "TCP", Type: "text", IsDefault: true}, {ID: "2", Name: "UDP", Type: "text"}},
			PropertyGroup: common.BKProcBindInfo,
			IsRequired:    true,
		},
		SubAttriubte{
			PropertyID:    "enable",
			PropertyName:  "Enable",
			Placeholder:   "enable port",
			IsEditable:    true,
			PropertyType:  common.FieldTypeBool,
			PropertyGroup: common.BKProcBindInfo,
			IsRequired:    true,
		},
	}
}

// Attribute attribute metadata definition
type Attribute metadata.Attribute

type SubAttriubte metadata.SubAttribute

// dbProcess 当前数据库中数据的格式
type dbProcess struct {
	ProcNum         *int64                 `field:"proc_num" json:"proc_num" bson:"proc_num" structs:"proc_num" mapstructure:"proc_num"`
	StopCmd         *string                `field:"stop_cmd" json:"stop_cmd" bson:"stop_cmd" structs:"stop_cmd" mapstructure:"stop_cmd"`
	RestartCmd      *string                `field:"restart_cmd" json:"restart_cmd" bson:"restart_cmd" structs:"restart_cmd" mapstructure:"restart_cmd"`
	ForceStopCmd    *string                `field:"face_stop_cmd" json:"face_stop_cmd" bson:"face_stop_cmd" structs:"face_stop_cmd" mapstructure:"face_stop_cmd"`
	ProcessID       int64                  `field:"bk_process_id" json:"bk_process_id" bson:"bk_process_id" structs:"bk_process_id" mapstructure:"bk_process_id"`
	FuncName        *string                `field:"bk_func_name" json:"bk_func_name" bson:"bk_func_name" structs:"bk_func_name" mapstructure:"bk_func_name"`
	WorkPath        *string                `field:"work_path" json:"work_path" bson:"work_path" structs:"work_path" mapstructure:"work_path"`
	BindIP          *string                `field:"bind_ip" json:"bind_ip" bson:"bind_ip" structs:"bind_ip" mapstructure:"bind_ip"`
	Priority        *int64                 `field:"priority" json:"priority" bson:"priority" structs:"priority" mapstructure:"priority"`
	ReloadCmd       *string                `field:"reload_cmd" json:"reload_cmd" bson:"reload_cmd" structs:"reload_cmd" mapstructure:"reload_cmd"`
	ProcessName     *string                `field:"bk_process_name" json:"bk_process_name" bson:"bk_process_name" structs:"bk_process_name" mapstructure:"bk_process_name"`
	Port            *string                `field:"port" json:"port" bson:"port" structs:"port" mapstructure:"port"`
	PidFile         *string                `field:"pid_file" json:"pid_file" bson:"pid_file" structs:"pid_file" mapstructure:"pid_file"`
	AutoStart       *bool                  `field:"auto_start" json:"auto_start" bson:"auto_start" structs:"auto_start" mapstructure:"auto_start"`
	AutoTimeGap     *int64                 `field:"auto_time_gap" json:"auto_time_gap" bson:"auto_time_gap" structs:"auto_time_gap" mapstructure:"auto_time_gap"`
	LastTime        time.Time              `field:"last_time" json:"last_time" bson:"last_time" structs:"last_time" mapstructure:"last_time"`
	CreateTime      time.Time              `field:"create_time" json:"create_time" bson:"create_time" structs:"create_time" mapstructure:"create_time"`
	BusinessID      int64                  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id" structs:"bk_biz_id" mapstructure:"bk_biz_id"`
	StartCmd        *string                `field:"start_cmd" json:"start_cmd" bson:"start_cmd" structs:"start_cmd" mapstructure:"start_cmd"`
	FuncID          *string                `field:"bk_func_id" json:"bk_func_id" bson:"bk_func_id" structs:"bk_func_id" mapstructure:"bk_func_id"`
	User            *string                `field:"user" json:"user" bson:"user" structs:"user" mapstructure:"user"`
	TimeoutSeconds  *int64                 `field:"timeout" json:"timeout" bson:"timeout" structs:"timeout" mapstructure:"timeout"`
	Protocol        *metadata.ProtocolType `field:"protocol" json:"protocol" bson:"protocol" structs:"protocol" mapstructure:"protocol"`
	Description     *string                `field:"description" json:"description" bson:"description" structs:"description" mapstructure:"description"`
	SupplierAccount string                 `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account" structs:"bk_supplier_account" mapstructure:"bk_supplier_account"`
	StartParamRegex *string                `field:"bk_start_param_regex" json:"bk_start_param_regex" bson:"bk_start_param_regex" structs:"bk_start_param_regex" mapstructure:"bk_start_param_regex"`
	PortEnable      *bool                  `field:"bk_enable_port" json:"bk_enable_port" bson:"bk_enable_port"`
	GatewayIP       *string                `field:"bk_gateway_ip" json:"bk_gateway_ip" bson:"bk_gateway_ip" structs:"bk_gateway_ip" mapstructure:"bk_gateway_ip"`
	GatewayPort     *string                `field:"bk_gateway_port" json:"bk_gateway_port" bson:"bk_gateway_port" structs:"bk_gateway_port" mapstructure:"bk_gateway_port"`
	GatewayProtocol *metadata.ProtocolType `field:"bk_gateway_protocol" json:"bk_gateway_protocol" bson:"bk_gateway_protocol" structs:"bk_gateway_protocol" mapstructure:"bk_gateway_protocol"`
	GatewayCity     *string                `field:"bk_gateway_city" json:"bk_gateway_city" bson:"bk_gateway_city" structs:"bk_gateway_city" mapstructure:"bk_gateway_city"`
}

// this works for the process instance which is used for a template.
type dbProcessTemplate struct {
	ID          int64  `field:"id" json:"id" bson:"id"`
	ProcessName string `field:"bk_process_name" json:"bk_process_name" bson:"bk_process_name"`
	BizID       int64  `field:"bk_biz_id" json:"bk_biz_id" bson:"bk_biz_id"`
	// the service template's, which this process template belongs to.
	ServiceTemplateID int64 `field:"service_template_id" json:"service_template_id" bson:"service_template_id"`

	// stores a process instance's data includes all the process's
	// properties's value.
	Property *dbProcessTemplateProperty `field:"property" json:"property" bson:"property"`

	Creator         string    `field:"creator" json:"creator" bson:"creator"`
	Modifier        string    `field:"modifier" json:"modifier" bson:"modifier"`
	CreateTime      time.Time `field:"create_time" json:"create_time" bson:"create_time"`
	LastTime        time.Time `field:"last_time" json:"last_time" bson:"last_time"`
	SupplierAccount string    `field:"bk_supplier_account" json:"bk_supplier_account" bson:"bk_supplier_account"`
}

type dbProcessTemplateProperty struct {
	BindIP   *PropertyString `field:"bind_ip" json:"bind_ip" bson:"bind_ip"`
	Port     *PropertyString `field:"port" json:"port" bson:"port"`
	Protocol *PropertyString `field:"protocol" json:"protocol" bson:"protocol"`
	Enable   *PropertyBool   `field:"bk_enable_port" json:"bk_enable_port" bson:"bk_enable_port"`
}

type PropertyString struct {
	Value          *string `field:"value" json:"value" bson:"value"`
	AsDefaultValue *bool   `field:"as_default_value" json:"as_default_value" bson:"as_default_value"`
}

type PropertyBool struct {
	Value          *bool `field:"value" json:"value" bson:"value"`
	AsDefaultValue *bool `field:"as_default_value" json:"as_default_value" bson:"as_default_value"`
}
