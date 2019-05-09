### 创建集群

- API: POST  /api/{version}/set/{bk_biz_id}   
- API 名称：create_set
- 功能说明：
	- 中文： 新建集群
	- English：create set

- input body
``` json
{

    "bk_set_name":"",
    "bk_parent_id":0,
    "bk_supplier_account":"",
    "bk_biz_id":1,
    "default":0
}
```

**注:以上 JSON 数据中各字段为必填字段或内置字段，它们在示例中的值仅为示例数据。**

- input 参数说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_set_id|int|是|无|集群ID|the set id|
|bk_parent_id|int|是|无|父节点的ID|the parent inst identifier|
|bk_biz_id|int|是|无|业务ID|business ID|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_set_name|string|是|无|集群名字 |set name|
|bk_capacity|int|否|无|设计容量|the design the capacity|
|description|string|否|无|备注|the remark|
|bk_service_status|enum|否|开放|服务状态:1/2(1:开放,2:关闭)|the service status:1/2 (1:open,2:close)|
|bk_set_env|enum|否|正式|环境类型：1/2/3(1:测试,2:体验,3:正式)|environment type:1/2/3(1:test,2:experience,3:formal)|
|bk_set_desc|string|否|无|集群描述|the set description|

 **注: 用户自定义的字段也可以作为参数传入。**




### 批量删除集群

- API: DELETE  /api/{version}/set/{bk_biz_id}/batch
- API 名称：batch_delete_set
- 功能说明：
	- 中文：批量删除集群
	- English：batch to delete set

- input body

``` json

{
    "delete":{
    "inst_ids":[]
    }
}

```

- input参数说明：

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_biz_id|int|是|无|业务ID|business ID|
|inst_ids|int array|是|无|集群ID集合|the set id collection|


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

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|操作结果|the result|



### 删除集群

- API: DELETE  /api/{version}/set/{bk_biz_id}/{bk_set_id}   
- API 名称：delete_set
- 功能说明：
	- 中文： 删除集群
	- English：delete set

- input body

    无

- input参数说明：

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_biz_id|int|是|无|业务ID|business ID|
|bk_set_id|int|是|无|集群ID|the set id|


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

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|操作结果|the result|

### 更新集群
- API： PUT /api/{version}/set/{bk_biz_id}/{bk_set_id}   
- API 名称：update_set
- 功能说明：
	- 中文： 更新集群
	- English：update set

- input body

``` json
{
    "bk_biz_id":2,
    "bk_set_name":"公共组件",
    "bk_set_desc":"",
    "bk_set_env":"3",
    "bk_service_status":"1",
    "description":"",
    "bk_capacity":null,
    "bk_supplier_account":"0"
}
```

**注:以上 JSON 数据中各字段为必填字段或内置字段，它们在示例中的值仅为示例数据。**

- input 参数说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_set_id|int|是|无|集群ID|the set id|
|bk_biz_id|int|是|无|业务ID|business ID|
|bk_set_name|string|否|无|集群名字 |set name|
|bk_capacity|int|否|无|设计容量|the design the capacity|
|description|string|否|无|备注|the remark|
|bk_service_status|enum|否|开放|服务状态:1/2(1:开放,2:关闭)|the service status:1/2 (1:open,2:close)|
|bk_set_env|enum|否|正式|环境类型：1/2/3(1:测试,2:体验,3:正式)|environment type:1/2/3(1:test,2:experience,3:formal)|
|bk_set_desc|string|否|无|集群描述|the set description|


 **注: 用户在使用的时候可以为每个Set增加字段的数量，这些自定义的字段也可以作为参数传入。**

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

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|操作结果|the result|

### 查询集群

- API： POST /api/{version}/set/search/{bk_supplier_account}/{bk_biz_id}   
- API 名称：search_set
- 功能说明：
	- 中文： 查询集群
	- English：search set

-  input body:
``` json
{
    "fields":[
        "bk_set_name"
    ],
    "page":{
        "start":0,
        "limit":100,
        "sort":"bk_set_name"
    },
    "condition":{
        "bk_set_name":"set_new"
    }
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input参数说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
| bk_supplier_account| string| 是| 无|开发商账号|supplier account code|
| bk_biz_id| int| 是|无|业务ID |  business ID|
| page| object| 是|无|分页参数 |page parameter|
| fields| array | 是| 无|查询字段|search fields|
| condition|  object| 是| 无|查询条件|search condition|

page 参数说明：

|名称|类型|必填| 默认值 | 说明 | Description|
|---|---| --- |---  | --- | ---|
| start|int|是|无|记录开始位置 |start record|
| limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
| sort| string| 否| 无|排序字段|the field for sort|

fields参数说明：

|名称|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_parent_id|int|否|无|父节点的ID|the parent inst identifier|
|bk_set_id|int|是|无|集群ID|the set id|
|bk_set_name|string|否|无|集群名字 |set name|
|bk_capacity|int|否|无|设计容量|the design the capacity|
|description|string|否|无|备注|the remark|
|bk_service_status|enum|否|开放|服务状态:1/2(1:开放,2:关闭)|the service status:1/2 (1:open,2:close)|
|bk_set_env|enum|否|正式|环境类型：1/2/3(1:测试,2:体验,3:正式)|environment type:1/2/3(1:test,2:experience,3:formal)|
|bk_set_desc|string|否|无|集群描述|the set description|

**注:所有字段均为Set定义的字段，这些字段包括预置字段，也包括用户自定义字段。**

condition 参数说明：

|名称|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_parent_id|int|否|无|父节点的ID|the parent inst identifier|
|bk_set_id|int|是|无|集群ID|the set id|
|bk_set_name|string|否|无|集群名字 |set name|
|bk_capacity|int|否|无|设计容量|the design the capacity|
|description|string|否|无|备注|the remark|
|bk_service_status|enum|否|开放|服务状态:1/2(1:开放,2:关闭)|the service status:1/2 (1:open,2:close)|
|bk_set_env|enum|否|正式|环境类型：1/2/3(1:测试,2:体验,3:正式)|environment type:1/2/3(1:test,2:experience,3:formal)|
|bk_set_desc|string|否|无|集群描述|the set description|

**注:所有字段均为Set定义的字段，这些字段包括预置字段，也包括用户自定义字段。**

- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": {
        "count": 1,
        "info": [
            {
                "bk_set_name": "内置模块集"
            }
        ]
    }
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|操作结果|the result|

data 说明

| 字段|类型|说明|Description|
|---|---|---|---|
|count|int|数据条数|the data item count|
|info|array|数据集合|the data array|

info 说明

| 字段|类型|说明|Description|
|---|---|---|---|
|bk_parent_id|int|父节点的ID|the parent inst identifier|
|bk_set_id|int|集群ID|the set id|
|bk_set_name|string|集群名字 |set name|
|bk_capacity|int|设计容量|the design the capacity|
|description|string|备注|the remark|
|bk_service_status|enum|服务状态:1/2(1:开放,2:关闭)|the service status:1/2 (1:open,2:close)|
|bk_set_env|enum|环境类型：1/2/3(1:测试,2:体验,3:正式)|environment type:1/2/3(1:test,2:experience,3:formal)|
|bk_set_desc|string|集群描述|the set description|

**注：此处按照fields所指定的字段进行配置，所有字段均为Set定义的字段，这些字段包括预置字段，也包括用户自定义字段。**

