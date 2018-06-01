### 创建分组基本信息
- API: POST /api/{version}/objectatt/group/new
- API 名称：create_group
- 功能说明: 
	- 中文: 创建分组
	- English: create a group for a object

- input body

``` json

{
	"bk_group_id": "znv94e4aa0l",
	"bk_group_name": "test_group_name",
	"bk_group_index": 1,
	"bk_obj_id": "cc_test_inst",
	"bk_supplier_account": "0"
}

```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

-  intput字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---| ---|
|bk_group_id|string| 是|无|分组ID，纯英文字符序列，不允许修改|the group identifier|
|bk_group_name|string|是|无|分组名字，用于展示|the group name|
|bk_group_index| int|是|0|分组排序|the group index|
|bk_obj_id|string|是|无|模型ID，用于指明该分组的所属|the object identifier|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|ispre|bool|否|false|true-内置分组，false- 非内置分组|true-inner, false- no inner|
|isdefault|bool|否|false|true-默认分组，false-用户自定义分组|true - the default, false-the customized|


- output

``` json
{
	"result": true,
	"bk_error_code": 0,
	"bk_error_msg": null,
	"data": {
		"id": 1046
	}
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object| 请求返回的数据 |return data|

data 说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
|id|int|新增数据记录的ID|the id of the new data record|

### 查询分组基本信息
- API: POST /api/{version}/objectatt/group/property/owner/{bk_supplier_account}/object/{bk_obj_id}
- API 名称：search_group
- 功能说明：
	- 中文：查询分组信息
	- English：query the grouping of models

- input body

	无

- input 字段说明

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
|bk_obj_id|string|是|无|模型ID|the object id|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|

- output


``` json
{
    "result": true,
    "code": 0,
    "message": null,
    "data": [
        {
            "bk_group_id": "default",
            "bk_group_index": 1,
            "bk_group_name": "基础信息",
            "bk_isdefault": true,
            "bk_obj_id": "host",
            "bk_supplier_account": "0",
            "id": 5,
            "ispre": false
        }
    ]
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | array| 请求返回的数据 |return data|

data 说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
|bk_group_id|string|分组标识|the group identifier|
|bk_group_index|int|分组排序|the group sort index|
|bk_group_name|string|分组名|the group name|
|bk_isdefault|bool|true-默认分组,false-普通分组|true - the defualt group, false - the common group|
|bk_obj_id|string|模型标识|the object identifier|
|bk_supplier_account|string|开发商账号|supplier account code|
|id|int|数据记录ID|the data record id|
|ispre|bool|true - 内置分组, false - 自定义定义分组|true - is inner, false - is customer|


### 修改分组基本信息
- API: PUT  /api/{version}/objectatt/group/update
- API 名称：update_group
- input body
``` json
{
	"condition": {
		"id": 5
	},
	"data": {
		"bk_group_index": 2
	}
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input参数说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_group_name|string|否|无|分组名字，用于展示|the group name|
|bk_group_index| int|否|0|分组排序|the group index|
|bk_obj_id|string|否|无|模型ID，用于指明该分组的所属|the object identifier|
|bk_supplier_account|string|否|无|开发商账号|supplier account code|
|ispre|bool|否|false|true内置分组，false 非内置分组|the group is inner|
|isdefault|bool|否|false|用于指明是否默认分组|true-the default group, false-custum|

condition 字段说明：

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_group_id|string| 否|无|分组ID，纯英文字符序列，不允许修改|the group identifier|
|bk_group_name|string|否|无|分组名字，用于展示|the group name|
|bk_group_index|int|否|0|分组排序|the group index|
|bk_obj_id|string|否|是|模型ID，用于指明该分组的所属|the object identifier|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_obj_id|string|是|无|模型ID，用于指明该分组的所属|the object identifier|

data 字段说明：

|字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_group_name|string|否|无|分组名字，用于展示|the group name|
|bk_group_index| int|否|0|分组排序|the group index|


- output

``` json
{
	"result": true,
	"bk_error_code": 0,
	"bk_error_msg": null,
	"data": "success"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | string| 操作结果 |the result|

### 删除分组基本信息
- API: DELETE /api/{version}/objectatt/group/groupid/{id}
- API 名称：delete_group
- 功能说明：
	- 中文：删除模型分组
	- English：delete the group of the object

- input body 

	无

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|id|int|是|无|分组记录标识|the group record id|


- output

``` json
{
	"result": true,
	"bk_error_code": 0,
	"bk_error_msg": null,
	"data": "success"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | string| 操作结果 |the result|

### 更新模型属性分组

- API: PUT /api/{version}/objectatt/group/property
- API 名称：update_property_group
- 功能说明：
	- 中文：更新某个属性所属的分组
	- English：update the grouping of a property

* input body
``` json
{
  "condition": {
      "":"" 
   },
  "data":{
      "":"" 
   }
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input 字段说明

condition 字段说明:

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_property_id|string|是|无|属性ID|the object attribute's identifier|
|bk_obj_id|string|是|无|模型ID|the object identifier|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|


data 字段说明：

|字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_property_group|string|是|无|分组ID|the group id|
|bk_property_index|int|否|0|定义属性在同一分组下的顺序|the property index|


- output

``` json
{
	"result": true,
	"bk_error_code": 0,
	"bk_error_msg": null,
	"data": "success"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | string| 操作结果|the result|


### 删除模型属性分组
- API: DELETE  /api/{version}/objectatt/group/owner/{bk_supplier_account}/object/{bk_obj_id}/propertyids/{bk_property_id}/groupids/{bk_group_id}
- API 名称：delete_object_property_group
- 功能说明：
	- 中文：删除模型属性分组
	- English：delete the group of the object property

- input body

	无

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_group_id|string|是|无|分组ID|the group record  id|
|bk_property_id|string|是|无|属性ID|the property identifier|
|bk_obj_id|string|是|无|模型ID|the object identifier|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|


- output

``` json
{
	"result": true,
	"bk_error_code": 0,
	"bk_error_msg": null,
	"data": "success"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | string| 操作结果 |the result|
