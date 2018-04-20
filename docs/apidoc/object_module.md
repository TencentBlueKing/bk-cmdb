
### 创建模块
- API： POST  /api/{version}/module/{bk_biz_id}/{bk_set_id}
- API 名称：create_module
- 功能说明：
	- 中文：创建模块	
	- English：create a module

- input body:
``` json
{
    "default":0,
    "bk_module_name":"cc_module",
    "bk_supplier_account":"0",
    "bk_parent_id":0
}
``` 

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

-  input字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_set_id|int|是|无|集群id|the set id|
|bk_parent_id|int|是|无|父实例节点的ID，当前实例节点的上一级实例节点，在拓扑结构中对于Module一般指的是Set的bk_set_id|the parent inst id|
|bk_module_id|string|是|无|模块标识|the module indentifier|
|bk_module_name|string|是|无|模块名|the module name|
|bk_supplier_account|string|否|无|开发商账号|supplier account code|
|bk_module_type|enum|否|普通|模块类型：普通/数据库|the module type: common/database|
|operator|string|否|无|主要维护人|the main maintainer|
|bk_bak_operator|string|否|无|备份维护人|the backup maintainer|


- output:
```
{
	"result": true,
	"bk_error_code": 0,
	"bk_error_msg": null,
	"data": {
		"id": 11142
	}
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object| 请求返回的数据 |the data response|

data 字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
|id|int|新增数据记录的ID|the id of the new module |


### 删除模块

- API： DELETE /api/{version}/module/{bk_biz_id}/{bk_set_id}/{bk_module_id}  
- API 名称：delete_module
- 功能说明：
	- 中文：删除模块
	- English：delete the module

- input body

    无

- input参数说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_biz_id|int|是|无|业务id|the application id|
|bk_set_id|int|是|无|集群id|the set id|
|bk_module_id|int|是|无|模块id|the module identifier|

- output:
```
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
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | string| 操作结果数据 |the result|

### 更新模块
- API： PUT /api/{version}/module/{bk_biz_id}/{bk_set_id}/{bk_module_id} 
- API 名称：update_module
- 功能说明：
	- 中文：更新模块
	- English：update the module

- input body
``` json
{
    "bk_module_name":"module_new"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

-  input参数说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_set_id|int|是|无|集群id|the set id|
|bk_parent_id|int|是|无|父实例节点的ID，当前实例节点的上一级实例节点，在拓扑结构中对于Module一般指的是Set的bk_set_id|the parent inst id|
|bk_module_id|string|是|无|模块标识|the module indentifier|
|bk_module_name|string|否|无|模块名|the module name|
|bk_supplier_account|string|否|无|开发商账号|supplier account code|
|bk_module_type|enum|否|普通|模块类型：普通/数据库|the module type: common/database|
|operator|string|否|无|主要维护人|the main maintainer|
|bk_bak_operator|string|否|无|备份维护人|the backup maintainer|

**注：以上字段仅为内置或必填参数，用户自定义的模块字段也可以作为输入参数。**

- output:
```
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
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | string| 操作结果数据 |the result|


### 查询模块
- API： POST /api/{version}/module/search/{bk_supplier_account}/{bk_biz_id}/{bk_set_id}           
- API 名称：search_module
- 功能说明：
	- 中文：查询模块	
	- English：search some modules

- input body

``` json
{
    "fields":[
        "bk_module_name"
    ],
    "page":{
        "start":0,
        "limit":100,
        "sort":"bk_module_name"
    },
    "condition":{
        "bk_module_name":"module_new"
    }
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input参数说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_biz_id|int|是|无|业务id|the application id|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_set_id|int|是|无|集群ID|the set identifier|
|bk_module_name|string|否|无|模块名字|the module name|
| page| object| 是|无|分页参数 |page parameter|
| fields| array | 是| 无|查询字段|search fields|
| condition|  object| 是| 无|查询条件|search condition|

page 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| start|int|是|无|记录开始位置 |start record|
| limit| int | 是| 无|每页限制条数|page limit |
| sort| string| 否| 无|排序字段|the field for sort|

fields参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
|bk_set_id|int|否|无|集群id|the set id|
|bk_parent_id|int|否|无|父实例节点的ID，当前实例节点的上一级实例节点，在拓扑结构中对于Module一般指的是Set的bk_set_id|the parent inst id|
|bk_module_id|string|否|无|模块标识|the module indentifier|
|bk_module_name|string|否|无|模块名|the module name|
|bk_supplier_account|string|否|无|开发商账号|supplier account code|
|bk_module_type|enum|否|普通|模块类型：普通/数据库|the module type: common/database|
|operator|string|否|无|主要维护人|the main maintainer|
|bk_bak_operator|string|否|无|备份维护人|the backup maintainer|

condition 说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|default|int|否|false|内置模块:0-普通模块,1-空闲机,2-故障机|inner field. 0-common module.1-idle module, 2- fault module|
|bk_set_id|int|否|无|集群id|the set id|
|bk_parent_id|int|否|无|父实例节点的ID，当前实例节点的上一级实例节点，在拓扑结构中对于Module一般指的是Set的bk_set_id|the parent inst id|
|bk_module_id|string|否|无|模块标识|the module indentifier|
|bk_module_name|string|否|无|模块名|the module name|
|bk_supplier_account|string|否|无|开发商账号|supplier account code|
|bk_module_type|enum|否|普通|模块类型：普通/数据库|the module type: common/database|
|operator|string|否|无|主要维护人|the main maintainer|
|bk_bak_operator|string|否|无|备份维护人|the backup maintainer|

- output
``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": {
        "count": 2,
        "info": [
            {
                "bk_module_name": "cc_service"
            },
            {
                "bk_module_name": "cmdb"
            }
        ]
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
| data | object| 操作结果数据 |the result|

data字段说明

|名称| 类型|说明|Description|
|---|---|---|---|
|count|int|数据数量|the data item count|
|info|array|结果集|the data result array|

info 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|default|int|内置模块:0-普通模块,1-空闲机,2-故障机|inner field. 0-common module.1-idle module, 2- fault module|
|bk_set_id|int|集群id|the set id|
|bk_parent_id|int|父实例节点的ID，当前实例节点的上一级实例节点，在拓扑结构中对于Module一般指的是Set的bk_set_id|the parent inst id|
|bk_module_id|string|模块标识|the module indentifier|
|bk_module_name|string|模块名|the module name|
|bk_supplier_account|string|开发商账号|supplier account code|
|bk_module_type|enum|模块类型：普通/数据库|the module type: common/database|
|operator|string|主要维护人|the main maintainer|
|bk_bak_operator|string|备份维护人|the backup maintainer|

**注:以上 字段仅为预置字段，不包含用户自定义字段。**
