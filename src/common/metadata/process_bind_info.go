/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */
package metadata

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

/*
version: 1.0 test
description:  由于在不同的运行版本中,进程的绑定信息的二维结构中，数据的列是不一样的。 所以又下面的实现

通过定义数据反序列化的方法来实现struct 的同一个属性在不同运行版本的环境上，实现进程绑定信息的多态。
主要是利用interface 的特性来实现，

defaultPropertyBindInfoHandle，defaultProcBindInfoHandle 是在使用反序列化 进程，进程模板中进程
绑定信息实际结构的对象

下面是defaultPropertyBindInfoHandle，defaultProcBindInfoHandle中UJSON和UBSON 含义的介绍
UJSON json反序列的方法，用于HTTP的消息处理,将数据解析到不同的struct上。 这个结构需要是ProcPropertyBindInfo，ProcBindInfoInterface interface 的实现
UBSON  bson 反序列化的方法， 用于数据库存储,将数据解析到不同的struct上。这个结构需要是ProcPropertyBindInfo，ProcBindInfoInterface interface 的实现

*/

var (
	//  内部变量，不允许改变，改变值请用对应的Register 方案
	defaultPropertyBindInfoHandle ProcPropertyBindInfoInterface = &openVersionPropertyBindInfo{}
	//  内部变量，不允许改变，改变值请用对应的Register 方案
	defaultProcBindInfoHandle ProcBindInfoInterface = &openVersionProcBindInfo{}
)

// TODO 等待实现， 替换已有的进程，进程模板中进程绑定信息实际结构的处理对象
func Register(propertyBindInfo ProcPropertyBindInfoInterface, procBindInfo ProcBindInfoInterface) {
	defaultPropertyBindInfoHandle = propertyBindInfo
	defaultProcBindInfoHandle = procBindInfo
}

// ProcPropertyBindInfoInterface 用来处理进程模板中bind info 数据的反序列化，
// 序列号使用默认的方法，目前只支持json, bson, 如果需要其他请新加
type ProcPropertyBindInfoInterface interface {
	UJSON(data []byte) (*ProcPropertyBindInfoValue, error)
	UBSON(data []byte) (*ProcPropertyBindInfoValue, error)
	/* 	JSON() ([]byte, error)
	   	BSON() ([]byte, error) */
}

// ProcBindInfoInterface 用来处理进程中bind info 数据的序反序列化，
// 序列号使用默认的方法，目前只支持json, bson, 如果需要其他请新加
type ProcBindInfoInterface interface {
	UJSON(data []byte) (*ProcBindInfo, error)
	UBSON(data []byte) (*ProcBindInfo, error)
	/* 	JSON() ([]byte, error)
	   	BSON() ([]byte, error) */
}

// ProcPropertyBindInfoRaw 用来限定进程模板中的校验和通过模板来限定进程实例的数据内容
type ProcPropertyBindInfoRaw interface {
	ExtractChangeInfoBindInfo(i *Process) (*ProcBindInfo, bool, bool)
	ExtractInstanceUpdateData(i *Process) *ProcBindInfo
	Validate() error
	// 留做扩招使用
	//Update(input ProcessProperty, rawProperty map[string]interface{})
	Datas() []interface{}
}

// ProcPropertyBindInfo 给服务模板使用的，来存储，校验服务模板中进程绑定的信息
type ProcPropertyBindInfo struct {
	// 通过Unmarshal 方法实现不同的数据类型
	Value ProcPropertyBindInfoValue
	// 给前端做兼容
	AsDefaultValue *bool `field:"as_default_value" json:"as_default_value" bson:"as_default_value"`
}

// ProcPropertyBindInfoValue 给服务模板使用的，来存储，校验服务模板中进程绑定的信息
type ProcPropertyBindInfoValue struct {
	// 通过Unmarshal 方法实现不同的数据类型
	raw ProcPropertyBindInfoRaw `field:"value" json:"value" bson:"value"`
}

// ProcBindInfo 给服务模板使用的，来存储，校验服务实例中进程绑定的信息
type ProcBindInfo struct {
	// 通过Unmarshal 方法实现不同的数据类型
	raw interface{}
}

func (ppb *ProcPropertyBindInfoValue) UnmarshalJSON(data []byte) error {
	val, err := defaultPropertyBindInfoHandle.UJSON(data)
	if err != nil {
		return err
	}
	ppb.raw = val.raw

	return nil
}

func (ppb *ProcPropertyBindInfoValue) UnmarshalBSON(data []byte) error {
	val, err := defaultPropertyBindInfoHandle.UBSON(data)
	if err != nil {
		return err
	}
	ppb.raw = val.raw
	return nil
}

func (ppb *ProcPropertyBindInfoValue) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(ppb.raw)

	return b, err
}

func (ppb ProcPropertyBindInfo) MarshalBSON() ([]byte, error) {
	data := make(map[string]interface{}, 0)
	data["as_default_value"] = ppb.AsDefaultValue
	if ppb.Value.raw != nil && len(ppb.Value.raw.Datas()) > 0 {
		data["value"] = ppb.Value.raw.Datas()

	}
	return bson.Marshal(data)
}

func (ppb ProcPropertyBindInfo) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{}, 0)
	data["as_default_value"] = ppb.AsDefaultValue
	if ppb.Value.raw != nil && len(ppb.Value.raw.Datas()) > 0 {
		data["value"] = ppb.Value.raw.Datas()

	}
	return json.Marshal(data)

}

// Validate 校验数据是否合法
func (pbi *ProcPropertyBindInfo) Validate() error {
	if pbi == nil || pbi.Value.raw == nil {
		return nil
	}
	return pbi.Value.Validate()
}

func (pbi *ProcPropertyBindInfo) ExtractInstanceUpdateData(input *Process) *ProcBindInfo {
	return pbi.Value.raw.ExtractInstanceUpdateData(input)
}

// Validate 校验数据是否合法
func (pbi *ProcPropertyBindInfoValue) Validate() error {
	if pbi == nil || pbi.raw == nil {
		return nil
	}
	return pbi.raw.Validate()
}

// ExtractChangeInfoBindInfo 生成对应环境中存储数据的对象
func (pbi *ProcPropertyBindInfoValue) ExtractChangeInfoBindInfo(i *Process) (*ProcBindInfo, bool, bool) {
	return pbi.raw.ExtractChangeInfoBindInfo(i)
}

func (pbi *ProcBindInfo) UnmarshalJSON(data []byte) error {
	val, err := defaultProcBindInfoHandle.UJSON(data)
	if err != nil {
		return err
	}
	pbi.raw = val.raw
	return nil
}

func (pbi *ProcBindInfo) UnmarshalBSON(data []byte) error {
	val, err := defaultProcBindInfoHandle.UBSON(data)
	if err != nil {
		return err
	}
	pbi.raw = val.raw
	return nil
}

func (pbi ProcBindInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(pbi.raw)
}

func (pbi ProcBindInfo) MarshalBSON() ([]byte, error) {
	if pbi.raw == nil {
		// 避免bson marshal nil 出现问题
		return bson.Marshal([]struct{}{})
	}
	return bson.Marshal(pbi.raw)
}

// Update  bind info 每次更新采用的是全量更新
func (pbi *ProcPropertyBindInfo) Update(input ProcessProperty, rawProperty map[string]interface{}) {
	pbi.AsDefaultValue = input.BindInfo.AsDefaultValue
	pbi.Value = input.BindInfo.Value
	return
}

/* 公开版本的进程bind 信息处理的方法 */

type openVersionPropertyBindInfo struct {
}

func (ov *openVersionPropertyBindInfo) UJSON(data []byte) (*ProcPropertyBindInfoValue, error) {
	if data == nil || len(data) == 0 {
		return nil, nil
	}
	bindInfo := &processPropertyBindInfo{}
	err := json.Unmarshal(data, bindInfo)
	return &ProcPropertyBindInfoValue{raw: bindInfo}, err
}

func (ov *openVersionPropertyBindInfo) UBSON(data []byte) (*ProcPropertyBindInfoValue, error) {
	if data == nil || len(data) == 0 {
		return nil, nil
	}
	var ret interface{}
	bson.Unmarshal(data, &ret)
	bindInfo1 := make(map[string]processPropertyBindInfoRow, 0)
	err := bson.Unmarshal(data, &bindInfo1)
	if err != nil {
		return nil, err
	}
	// TODO 找出为什么是map[string]processPropertyBindInfoRow 原因
	// bson 怀疑是因为bson 不支持数组做为顶级节点，
	var bindInfo []processPropertyBindInfoRow
	for _, item := range bindInfo1 {
		bindInfo = append(bindInfo, item)
	}
	tmp := processPropertyBindInfo(bindInfo)
	/*bindInfo := &processPropertyBindInfo{}
	err := bson.Unmarshal(data, bindInfo)*/

	return &ProcPropertyBindInfoValue{raw: &tmp}, err
}

type openVersionProcBindInfo struct {
}

func (ov *openVersionProcBindInfo) UJSON(data []byte) (*ProcBindInfo, error) {
	if data == nil || len(data) == 0 {
		return nil, nil
	}
	bindInfo := &processBindInfo{}
	err := json.Unmarshal(data, bindInfo)
	return &ProcBindInfo{raw: bindInfo}, err
}

func (ov *openVersionProcBindInfo) UBSON(data []byte) (*ProcBindInfo, error) {
	if data == nil || len(data) == 0 {
		return nil, nil
	}
	bindInfo := &processBindInfo{}
	err := bson.Unmarshal(data, bindInfo)
	return &ProcBindInfo{raw: bindInfo}, err
}

type processBindInfoRow struct {
	IP       *string `field:"ip" json:"ip" bson:"ip"`
	Port     *string `field:"port" json:"port" bson:"port"`
	Protocol *string `field:"protocol" json:"protocol" bson:"protocol"`
	Enable   *bool   `field:"enable" json:"enable" bson:"enable"`
}

type processPropertyBindInfoRow struct {
	IP       PropertyBindIP   `field:"ip" json:"ip" bson:"ip"`
	Port     PropertyPort     `field:"port" json:"port" bson:"port"`
	Protocol PropertyProtocol `field:"protocol" json:"protocol" bson:"protocol"`
	Enable   PropertyBool     `field:"enable" json:"enable" bson:"enable"`
}

type processBindInfo []processBindInfoRow
type processPropertyBindInfo []processPropertyBindInfoRow

func (pbi *processPropertyBindInfo) Validate() (err error) {
	// call all field's Validate method one by one
	for _, row := range *pbi {
		if err := row.IP.Validate(); err != nil {
			return err
		}
		if err := row.Port.Validate(); err != nil {
			return err
		}
		if err := row.Protocol.Validate(); err != nil {
			return err
		}
		if err := row.Enable.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// ExtractChangeInfo get bind info changes that will be applied to process instance,
func (pbi *processPropertyBindInfo) ExtractChangeInfoBindInfo(i *Process) (*ProcBindInfo, bool, bool) {
	var changed, isNamePortChanged bool

	preProcBindInfoMap := make(map[int]processBindInfoRow, 0)
	if i.BindInfo != nil {
		preProcBindInfoArr, ok := i.BindInfo.raw.(*processBindInfo)
		if !ok {
			// 这个地方出现panic， 证明在代码处理上面出现问题了
			panic("ExtractChangeInfoBindInfo bind info panic, reason: not open version process bind info struct")
		}
		for idx, row := range *preProcBindInfoArr {
			preProcBindInfoMap[idx] = row
		}
	}

	for idx, row := range *pbi {
		preRow, ok := preProcBindInfoMap[idx]
		if !ok {
			preRow = processBindInfoRow{}
		}
		if IsAsDefaultValue(row.IP.AsDefaultValue) {
			if row.IP.Value == nil && preRow.IP != nil {
				preRow.IP = nil
				changed = true
			} else if row.IP.Value != nil && preRow.IP == nil {
				*preRow.IP = row.IP.Value.IP()
				changed = true
			} else if row.IP.Value != nil && preRow.IP != nil && row.IP.Value.IP() != *preRow.IP {
				*preRow.IP = row.IP.Value.IP()
				changed = true
			}
		}

		if IsAsDefaultValue(row.Port.AsDefaultValue) {
			if row.Port.Value == nil && preRow.Port != nil {
				preRow.Port = nil
				changed = true
				isNamePortChanged = true
			} else if row.Port.Value != nil && preRow.Port == nil {
				*preRow.Port = *row.Port.Value
				isNamePortChanged = true
				changed = true
			} else if row.Port.Value != nil && preRow.Port != nil && *row.Port.Value != *preRow.Port {
				*preRow.Port = *row.Port.Value
				isNamePortChanged = true
				changed = true
			}
		}

		if IsAsDefaultValue(row.Protocol.AsDefaultValue) {
			if row.Protocol.Value == nil && preRow.Protocol != nil {
				preRow.Protocol = nil
				changed = true
			} else if row.Protocol.Value != nil && preRow.Protocol == nil {
				*preRow.Protocol = row.Protocol.Value.String()
				changed = true
			} else if row.Protocol.Value != nil && preRow.Protocol != nil && row.Protocol.Value.String() != *preRow.Protocol {
				*preRow.Protocol = row.Protocol.Value.String()
				changed = true
			}
		}

		if IsAsDefaultValue(row.Enable.AsDefaultValue) {
			if row.Enable.Value == nil && preRow.Enable != nil {
				preRow.Enable = nil
				changed = true
			} else if row.Enable.Value != nil && preRow.Enable == nil {
				*preRow.Enable = *row.Enable.Value
				changed = true
			} else if row.Enable.Value != nil && preRow.Enable != nil && *row.Enable.Value != *preRow.Enable {
				*preRow.Enable = *row.Enable.Value
				changed = true
			}
		}
	}
	newProBindInfo := make([]processBindInfoRow, len(*pbi))
	for idx, row := range preProcBindInfoMap {
		newProBindInfo[idx] = row
	}
	if len(newProBindInfo) == 0 {
		return nil, changed, isNamePortChanged
	}

	return &ProcBindInfo{raw: newProBindInfo}, false, false
}

// ExtractChangeInfo get bind info changes that will be applied to process instance,
func (pbi *processPropertyBindInfo) ExtractInstanceUpdateData(i *Process) *ProcBindInfo {
	// 用户输入的进程绑定的信息
	procBindInfoMap := make(map[int]processBindInfoRow, 0)
	if i.BindInfo != nil && i.BindInfo.raw != nil {
		procBindInfoArr, ok := i.BindInfo.raw.(*processBindInfo)
		if !ok {
			// 这个地方出现panic， 证明在代码处理上面出现问题了
			panic("ExtractChangeInfoBindInfo bind info panic, reason: not open version process bind info struct")
		}
		for idx, row := range *procBindInfoArr {
			procBindInfoMap[idx] = row
		}
	}

	//	注意 process bind info 更新必须是全量
	// 数据已模板为主
	var newProBindInfo []processBindInfoRow
	for idx, row := range *pbi {
		// 用户输入进程绑定信息中的一行
		procBindInfoItem, ok := procBindInfoMap[idx]
		if !ok {
			//
			continue
		}

		// 通过模板处理后，存储到db中process bind info 中的一行
		newbProcBindInfoItem := processBindInfoRow{}

		if IsAsDefaultValue(row.IP.AsDefaultValue) == false {
			if procBindInfoItem.IP != nil {
				newbProcBindInfoItem.IP = procBindInfoItem.IP
			}
		} else {
			if row.IP.Value != nil {
				ip := row.IP.Value.String()
				newbProcBindInfoItem.IP = &ip
			}
		}

		if IsAsDefaultValue(row.Port.AsDefaultValue) == false {
			if procBindInfoItem.Port != nil {
				newbProcBindInfoItem.Port = procBindInfoItem.Port
			}
		} else {
			if row.Port.Value != nil {
				newbProcBindInfoItem.Port = row.Port.Value
			}
		}

		if IsAsDefaultValue(row.Protocol.AsDefaultValue) == false {
			if procBindInfoItem.Protocol != nil {
				newbProcBindInfoItem.Protocol = procBindInfoItem.Protocol
			}
		} else {
			if row.Protocol.Value != nil {
				protocol := row.Protocol.Value.String()
				newbProcBindInfoItem.Protocol = &protocol
			}
		}

		newProBindInfo = append(newProBindInfo, newbProcBindInfoItem)
	}

	return &ProcBindInfo{raw: newProBindInfo}
}

func (ov processPropertyBindInfo) Datas() []interface{} {
	var ret []interface{}
	for _, item := range ov {
		ret = append(ret, item)
	}
	return ret
}
