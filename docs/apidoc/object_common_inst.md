
### 添加对象实例

- API POST /api/{version}/inst/{bk_supplier_account}/{bk_obj_id}
- API 名称：create_inst
- 功能说明：
	- 中文：添加实例
	- English：create a new inst 

- input body
``` json
{
	"bk_inst_name": "example"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input字段说明

|字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_obj_id|string|是|无|模型ID|the object id|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_inst_name|string|是|无|实例名|the inst name|

注：此处的输入参数仅对必填参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段。

- output

``` json
{
	"result": true,
	"bk_error_code": 0,
	"bk_error_msg": null,
	"data": {
		"bk_inst_id": 67
	}
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object| 请求返回的数据 |the data response|

data 字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
|bk_inst_id|int|新增数据记录的ID|the id of the new inst |


### 查询实例

- API POST /api/{version}/inst/association/search/owner/{bk_supplier_account}/object/{bk_obj_id}
- API 名称：search_inst
- 功能说明：
	- 中文：查询实例
	- English：search insts by condition 

- input body
``` json
{
	"page": {
		"start": 0,
		"limit": 10,
		"sort": ""
	},
	"condition": {
		"test_obj": [{
			"field": "bk_inst_name",
			"operator": "$eq",
			"value": "test"
		}]
	},
	"fields": {
		"test":[]
	}
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input字段说明

|字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_obj_id|string|是|无|模型ID|the object id|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
| page| object| 是|无|分页参数 |page parameter|
|condition| object | 否|无|查询条件|the search condition|
|fields|object|否|无|查询的字段|the search fields|


page 参数说明：

|名称|类型|必填| 默认值 | 说明 | Description|
|---|---| --- |---  | --- | ---|
| start|int|是|无|记录开始位置 |start record|
| limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
| sort| string| 否| 无|排序字段|the field for sort|

condition 参数说明：

|名称|类型|必填| 默认值 | 说明 | Description|
|---|---| --- |---  | --- | ---|
|test_obj|string|是|无|此处仅为示例数据，需要被设置为模型的标识符，在页面上配置的英文名|the engilish name, the object indentifier|
|field|string|是|无|取值为模型的字段名|the field name of a object|
|operator|string|是|无|取值为：$regex $eq $ne|the available value: $regex $eq $ne|
|value|string|是|无|field配置的模型字段名所对应的值|the value of the field|

fields 参数说明：

|名称|类型|必填| 默认值 | 说明 | Description|
|---|---| --- |---  | --- | ---|
|test|string|无|此处仅为示例数据，需要被设置为模型的标识符，在页面上配置的英文名，此字段所取得值为改模型所定义的模型的字段的集合|the engilish name, the object indentifier. The value is the collection of all the fields of the object.|

注：此处的输入参数仅对必填参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段。

- output

``` json
  {
	"result": true,
	"bk_error_code": 0,
	"bk_error_msg": "success",
	"data": {
		"count": 1,
		"info": [{
			"bk_inst_id": 1,
			"bk_inst_name": "test",
			"bk_obj_id": "test",
			"bk_supplier_account": "0",
			"create_time": "2018-04-17T14:50:15.993+08:00",
			"last_time": "2018-04-17T15:00:49.274+08:00",
			"test_asst": [{
				"bk_inst_id": 2,
				"bk_inst_name": "test2",
				"bk_obj_id": "test_obj",
				"id": "2"
			}]
		}]
	}
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object| 请求返回的数据 |the data response|

data 字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
|id|string|已存储的关联实例的id|The id of the associated instance that has been stored.|
|bk_inst_id|int|新增数据记录的ID|the id of the new inst |
|bk_supplier_account|string|开发商账号|supplier account code|
|bk_obj_id|string|模型ID|the object id|
|create_time|string|数据创建的时间|the creation date time|
|last_time|string|最后修改时间|the last modify date time|
|test_asst|string|test 的关联字段，用户自定义的名字，实际应用中不会与此相同|Test's associated field, user-defined name, will not be the same in practice.|

### 删除对象实例

- API: DELETE  /api/{version}/inst/{bk_supplier_account}/{bk_obj_id}/{bk_inst_id}
- API 名称：delete_inst
- 功能说明：
	- 中文： 删除实例
	- English：delete a inst

- input body

	无

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_obj_id|int|是|无|对象ID|the object id|
|bk_inst_id|int|是|无|实例ID|the inst id|


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

- output 字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | string| 操作结果 |the result|
