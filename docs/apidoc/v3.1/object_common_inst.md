### 添加对象实例

- API POST /api/{version}/inst/{bk_supplier_account}/{bk_obj_id}
- API 名称：create_inst
- 功能说明：
	- 中文：创建实例
	- English：create a new inst 

- input body  (通用实例示例)
``` json
{
	"bk_inst_name": "example"
}
```

- input body  (云区域示例)
``` json
{
    "bk_cloud_name":"example18",
    "bk_supplier_account":"0",
    "bk_biz_id":0
}
```


**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input字段说明

|字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_obj_id|string|是|无|模型ID，新建云区域时为plat|the object id，when to create a new bk cloud it should be plat |
|bk_supplier_account|string|是|无|开发商账号,独立部署请填"0"|supplier account code,enterprise version is "0"|
|bk_inst_name/bk_cloud_name|string|是|无|实例名,当创建对象为云区域时为bk_cloud_name|the inst name, when the object is plat, it should be bk_cloud_name|
| bk_biz_id|int|否|无| 业务ID |business ID|

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
    "page":{
        "start":0,
        "limit":10,
        "sort":"-bk_inst_id"
    },
    "fields":{

    },
    "condition":{
        "bk_weblogic":[
            {
                "field":"bk_inst_name",
                "operator":"$regex",
                "value":"qq"
            }
        ]
    }
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input字段说明

|字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_obj_id|string|是|无|模型ID|the object id|
|bk_supplier_account|string|是|无|开发商账号,独立部署请填"0"|supplier account code,enterprise version is "0"|
|page| object| 是|无|分页参数 |page parameter|
|condition| object | 否|无|查询条件|the search condition|
|fields|string array|否|无|查询的字段|the search fields|


page 参数说明：

|名称|类型|必填| 默认值 | 说明 | Description|
|---|---| --- |---  | --- | ---|
| start|int|是|无|记录开始位置 |start record|
| limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
| sort| string| 否| 无|排序字段|the field for sort|

condition 参数说明：

|名称|类型|必填| 默认值 | 说明 | Description|
|---|---| --- |---  | --- | ---|
|bk_weblogic|string|是|无|此处仅为示例数据，需要被设置为模型的标识符，在页面上配置的英文名|the engilish name, the object indentifier|
|field|string|是|无|取值为模型的字段名|the field name of a object|
|operator|string|是|无|取值为：$regex $eq $ne|the available value: $regex $eq $ne|
|value|string|是|无|field配置的模型字段名所对应的值|the value of the filed|

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
|test_asst|string|test_asst为此实例的关联字段，返回关联模型对应的实例。|Test's associated field, user-defined name.|


### 更新对象实例（包含云区域）

- API: PUT  /api/{version}/inst/{bk_supplier_account}/{bk_obj_id}/{bk_inst_id}
- API 名称：update_inst
- 功能说明：
	- 中文： 更新对象实例
	- English：update a inst

- input body (通用实例示例)

``` json
  {
	"bk_inst_name": "aaaaaa"
}
```

- input body (云区域示例)

``` json
  {
	"bk_cloud_name": "cloud1"
}
```

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_obj_id|string|是|无|模型ID，更新云区域时为plat|the object id, when update bk cloud it should be plat|
|bk_inst_id|int|是|无|实例ID,更新云区域是为bk_cloud_id|the inst id, when update bk cloud it should be cloud area ID|
|bk_inst_name|string|否|无|实例名，也可以为其它自定义字段|the inst name, can be other field|
|bk_cloud_name|string|否|无|云区域名，更新云区域名时需要|the cloud area name, it is in need where update plat|


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



### 批量更新通用对象实例

- API: PUT  /api/{version}/inst/{bk_supplier_account}/{bk_obj_id}/batch
- API 名称：batch_update_inst
- 功能说明：
	- 中文： 更新对象实例
	- English：update a inst

- input body (通用实例示例)

``` json
{
"update":[
	{
	  "datas":{
	  	"bk_inst_name":"batch_update"
	  },
      "inst_id":46
	 }
    ]
}
```


- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_obj_id|string|是|无|模型ID|the object id|
|update|object array|是|无|实例被更新的字段及值|the inst value|

- update 字段说明
| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_inst_name|string|否|无|实例名，也可以为其它自定义字段|the inst name, can be other field|
|datas|object|是|无|实例被更新的字段取值|the inst value|
|inst_id|int|是|无|指明datas 用于更新的具体实例|set the datas owner|

- datas 字段说明

**datas 是map类型的对象，key 是实例对应的模型定义的字段，value是字段的取值**



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


### 批量删除对象实例

- API: DELETE  /api/{version}/inst/{bk_supplier_account}/{bk_obj_id}/batch
- API 名称：batch_delete_inst
- 功能说明：
	- 中文： 批量删除实例
	- English：batch delete a inst

- input body

``` json
{
    "delete":{
    "inst_ids":[]
    }
}
```

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_obj_id|string|是|无|模型ID，删除对象为云区域时为plat|the object id, when delete bk cloud it should be plat|
|inst_ids|int array|是|无|实例ID集合|the inst id collection|


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

### 删除对象实例（包含云区域）

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
|bk_obj_id|string|是|无|模型ID，删除对象为云区域时为plat|the object id, when delete bk cloud it should be plat|
|bk_inst_id|int|是|无|实例ID，删除云区域时为云区域ID|the inst id, when delete bk cloud it should be cloud area ID|


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



### 查询实例关联拓扑

- API: POST  /api/{version}/inst/search/topo/owner/{bk_supplier_account}/object/{bk_obj_id}/inst/{bk_inst_id}
- API 名称：search_inst_association_topo
- 功能说明：
	- 中文： 查询实例关联拓扑
	- English：query the instance association topology

- input body

	无

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_obj_id|string|是|无|模型ID|the object id|
|bk_inst_id|int|是|无|实例ID|the inst id|


- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "success",
    "data": [
        {
            "curr": {
                "bk_inst_id": 17,
                "bk_inst_name": "192.168.1.1",
                "bk_obj_icon": "icon-cc-host",
                "bk_obj_id": "host",
                "bk_obj_name": "主机",
                "children": [],
                "count": 0
            },
            "next": [
                {
                    "bk_inst_id": 0,
                    "bk_inst_name": "",
                    "bk_obj_icon": "icon-cc-subnet",
                    "bk_obj_id": "plat",
                    "bk_obj_name": "云区域",
                    "children": [
                        {
                            "bk_inst_id": 0,
                            "bk_inst_name": "default area",
                            "bk_obj_icon": "",
                            "bk_obj_id": "plat",
                            "bk_obj_name": "",
                            "id": "0"
                        }
                    ],
                    "count": 1
                }
            ],
            "prev": [
                {
                    "bk_inst_id": 0,
                    "bk_inst_name": "",
                    "bk_obj_icon": "icon-cc-business",
                    "bk_obj_id": "rel",
                    "bk_obj_name": "关联",
                    "children": [
                        {
                            "bk_inst_id": 162,
                            "bk_inst_name": "test1",
                            "bk_obj_icon": "",
                            "bk_obj_id": "rel",
                            "bk_obj_name": ""
                        }
                    ],
                    "count": 1
                }
            ]
        }
    ]
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object array | 查询结果 |the result|


- data 字段说明

|名称|类型|说明|Description|
|---|---|---|---|
|curr|object|当前实例节点的信息|the current instance node information |
|next|object array|当前节点的子节点集合|the current node's child node collection|
|prev|object array| 当前节点的父节点结合|the current node's parent node collection|


- curr 字段说明

|名称|类型|说明|Description|
|---|---|---|---|
|bk_inst_id|int|实例ID|the inst ID|
|bk_inst_name|string|实例用于展示的名字|the name of the instance is used to display|
|bk_obj_icon|string|模型图标的名字|the object's icon|
|bk_obj_id|string|模型ID|the object's id|
|bk_obj_name|string|模型用于展示的名字|the name of the object is used to display|
|children|object array|本模型下所有被关联的实例的集合|a collection of all associated instances in this model|
|count|children 包含节点的数量|children contains the number of nodes|


- next 字段说明

|名称|类型|说明|Description|
|---|---|---|---|
|bk_inst_id|int|实例ID|the inst ID|
|bk_inst_name|string|实例用于展示的名字|the name of the instance is used to display|
|bk_obj_icon|string|模型图标的名字|the object's icon|
|bk_obj_id|string|模型ID|the object's id|
|bk_obj_name|string|模型用于展示的名字|the name of the object is used to display|
|children|object array|本模型下所有被关联的实例的集合|a collection of all associated instances in this model|
|count|children 包含节点的数量|children contains the number of nodes|

- next/children 字段说明

|名称|类型|说明|Description|
|---|---|---|---|
|bk_inst_id|int|实例ID|the inst ID|
|bk_inst_name|string|实例用于展示的名字|the name of the instance is used to display|
|bk_obj_icon|string|模型图标的名字|the object's icon|
|bk_obj_id|string|模型ID|the object's id|
|bk_obj_name|string|模型用于展示的名字|the name of the object is used to display|



- prev 字段说明

|名称|类型|说明|Description|
|---|---|---|---|
|bk_inst_id|int|实例ID|the inst ID|
|bk_inst_name|string|实例用于展示的名字|the name of the instance is used to display|
|bk_obj_icon|string|模型图标的名字|the object's icon|
|bk_obj_id|string|模型ID|the object's id|
|bk_obj_name|string|模型用于展示的名字|the name of the object is used to display|
|children|object array|本模型下所有被关联的实例的集合|a collection of all associated instances in this model|
|count|children 包含节点的数量|children contains the number of nodes|

- prev/children 字段说明

|名称|类型|说明|Description|
|---|---|---|---|
|bk_inst_id|int|实例ID|the inst ID|
|bk_inst_name|string|实例用于展示的名字|the name of the instance is used to display|
|bk_obj_icon|string|模型图标的名字|the object's icon|
|bk_obj_id|string|模型ID|the object's id|
|bk_obj_name|string|模型用于展示的名字|the name of the object is used to display|


### 查询业务实例拓扑

- API: GET /api/{version}/topo/inst/{bk_supplier_account}/{bk_biz_id}?level={level}
- API 名称：search_biz_inst_topo
- 功能说明：
	- 中文： 查询业务实例拓扑
	- English：query business instance topology

- input body

	无

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_biz_id|int|是|无|业务id|the business id|
|level|int|否|2|拓扑的层级索引，索引取值从0开始，当设置为 -1 的时候会读取完整的业务实例拓扑|the topology level, read full topology when set to -1|


- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "success",
    "data": [
        {
            "bk_inst_id": 2,
            "bk_inst_name": "蓝鲸",
            "bk_obj_id": "biz",
            "bk_obj_name": "业务",
            "child": [
                {
                    "bk_inst_id": 3,
                    "bk_inst_name": "作业平台",
                    "bk_obj_id": "set",
                    "bk_obj_name": "集群",
                    "child": [
                        {
                            "bk_inst_id": 5,
                            "bk_inst_name": "job",
                            "bk_obj_id": "module",
                            "bk_obj_name": "模块",
                            "child": []
                        }
                    ]
                }
            ]
        }
    ]
}

```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object array | 查询结果 |the result|

- data 字段说明

|名称|类型|说明|Description|
|---|---|---|---|
|bk_inst_id|int|实例ID|the inst ID|
|bk_inst_name|string|实例用于展示的名字|the name of the instance is used to display|
|bk_obj_id|string|模型ID|the object's id|
|bk_obj_name|string|模型用于展示的名字|the name of the object is used to display|
|child|object array|当前节点下的所有实例的集合|Collection of all instances under the current node|

- child 字段说明

|名称|类型|说明|Description|
|---|---|---|---|
|bk_inst_id|int|实例ID|the inst ID|
|bk_inst_name|string|实例用于展示的名字|the name of the instance is used to display|
|bk_obj_id|string|模型ID|the object's id|
|bk_obj_name|string|模型用于展示的名字|the name of the object is used to display|
|child|object array|当前节点下的所有实例的集合|Collection of all instances under the current node|





### 查询实例列表

- API: POST /api/{version}/inst/search/owner/{bk_supplier_account}/object/{bk_obj_id}
- API 名称：search_inst_by_object
- 功能说明：
	- 中文： 查询给定模型的实例列表
	- English：query  instance list

- input body

``` json

{
    "page": {
        "start": 0,
        "limit": 10,
        "sort": ""
    },
    "fields": [],
    "condition": {   
    }
}

```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_obj_id|string|是|无|自定义模型ID，查询区域时为plat|the object id， when search bk cloud it should be plat|
|fields|array|否|无|指定查询的字段|need to show|
|condition|object|否|无|查询条件|search condition|
|page|object|否|无|分页条件|page condition|

page 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| start|int|是|无|记录开始位置 |start record|
| limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
| sort| string| 否| 无|排序字段|the field for sort|


fields参数说明：

参数为查询的目标实例对应的模型定义的所有字段。


condition 参数说明：

condition 参数为查询的目标实例对应的模型定义的所有字段。

- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "success",
    "data": {
        "count": 4,
        "info": [
            {
                "bk_cloud_id": 0,
                "bk_cloud_name": "default area",
                "bk_supplier_account": ""
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
|data|string|操作结果|the result|

- data 字段说说明

|名称|类型|说明|Description|
|---|---|---|---|
|count|int|info 集合中元素的数量|the inst ID|
|info|object array |查询的模型的实例集合|the set of instances of the model being queried|

- info 字段说明（此处仅对示例中数据做说明）

|名称|类型|说明|Description|
|---|---|---|---|
|bk_cloud_id|int|云区域ID|the cloud id|
|bk_cloud_name|string|云区域名|the cloud name|
|bk_supplier_account|string|开发商账号|supplier account code|

