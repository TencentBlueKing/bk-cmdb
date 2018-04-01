
### 添加对象实例

- API POST /api/{version}/inst/{bk_supplier_account}/{bk_obj_id}
- API 名称：create_inst
- 功能说明：
	- 中文：添加实例
	- English：create a new inst 

- input body
``` json
{
	"bk_inst_name": "example",
	"bk_parent_id": 96,
	"bk_biz_id": 96
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input字段说明

|字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_obj_id|string|是|无|模型ID|the object id|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_parent_id|int|是|无|主线模型的父实例ID,在拓扑结构中，当前实例节点的上一级实例节点|the main line parent id|
|bk_inst_name|string|是|无|实例名|the inst name|
|bk_biz_id|int|是|无|业务ID|the application id|

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


### 删除对象实例

- API: DELETE  /api/{version}/inst/{bk_supplier_account}/{bk_biz_id}/{bk_inst_id}
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
|bk_biz_id|int|是|无|业务ID|the application id|
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
