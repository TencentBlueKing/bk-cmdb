
### 新加自定义API接口

*  API: POST /api/{version}/userapi
* API名称： create_custom_query
* 功能说明：
	* 中文： 添加自定义api
	* English ：create customize query

*  input body:
```
{
    "bk_biz_id":12,
    "info":"{\"condition\":[{\"bk_obj_id\":\"biz\",\"condition\":[{\"field\":\"default\",\"operator\":\"$ne\",\"value\":1}],\"fields\":[]},{\"bk_obj_id\":\"set\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"module\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"host\",\"condition\":[{\"field\":\"bk_host_innerip\",\"operator\":\"$eq\",\"value\":\"127.0.0.1\"}],\"fields\":[\"bk_host_innerip\",\"bk_host_outerip\",\"bk_agent_status\"]}]}",
    "name":"api1"
}
``` 

* input参数说明


| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | --- |---| --- | --- | ---|
| bk_biz_id|int|是|无| 业务ID |business ID|
| info|json string|是|无|通用查询条件 | common search query parameters|
| name|string|是|无|收藏的名称|the name of user api|

info 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_obj_id| string| 否| 无|对象名,可以为biz,set,module,host,object|object name, it can be biz,set,module,host,object|
| fields| string数组| 否| 无|查询输出字段|fields output|
| condition| object array| 否| 无|查询条件|search condition|

二级condition 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| field| string| 否| 无|对象的字段|field of object|
| operator| string| 否| 无|操作符, $eq为相等，$neq为不等，$in为属于，$nin为不属于|$eq is equal,$in is belongs, $nin is not belong,$neq is not equal|
| value| string| 否| 无|字段对应的值|the value of field|



* output 

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":{
        "id":"b80nu3dmjrccd9i5r1eg"
    }
}
```
*  output字段说明

| 名称  | 类型  | 说明 | Description|
|---|---|---|---|
| id| string| 自定义api主键ID|Primary key ID|


### 更新自定义API接口

*  API: PUT /api/{version}/userapi/{bk_biz_id}/{id}

* API名称： update_custom_query
* 功能说明：
	* 中文：更新自定义api
	* English ：update customize query

*  input body:
```
{
    "info":"{\"condition\":[{\"bk_obj_id\":\"biz\",\"condition\":[{\"field\":\"default\",\"operator\":\"$ne\",\"value\":1}],\"fields\":[]},{\"bk_obj_id\":\"set\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"module\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"host\",\"condition\":[{\"field\":\"bk_host_innerip\",\"operator\":\"$eq\",\"value\":\"127.0.0.1\"}],\"fields\":[\"bk_host_innerip\",\"bk_host_outerip\",\"bk_agent_status\"]}]}",
    "name":"api1"
}
```


* input参数说明

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | --- |---| --- | --- | ---|
| bk_biz_id|int|是|无| 业务ID |business ID|
| id|string|是|无| 主键ID |Primary key ID|
| info|json string|否|无|通用查询条件 | common search query parameters|
| name|string|否|无|收藏的名称|the name of user api|

info 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_obj_id| string| 否| 无|对象名,可以为biz,set,module,host,object|object name, it can be biz,set,module,host,object|
| fields| string数组| 否| 无|查询输出字段|fields output|
| condition| object array| 否| 无|查询条件|search condition|

condition 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| field| string| 否| 无|对象的字段|field of object|
| operator| string| 否| 无|操作符, $eq为相等，$neq为不等，$in为属于，$nin为不属于|$eq is equal,$in is belongs, $nin is not belong,$neq is not equal|
| value| string| 否| 无|字段对应的值|the value of field|


* output 
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":null
}
```
* output 参数说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null| 请求返回的数据 |return data|

### 删除自定义API接口

*  API:  DELETE /api/{version}/userapi/{bk_biz_id}/{id}

* API名称：  delete_custom_query
* 功能说明：
	* 中文： 删除自定义api
	* English ：delete customize query
*  input body
无

* input参数说明

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | --- |---| --- | --- | ---|
| bk_biz_id|int|是|无| 业务ID |business ID|
| id|string|是|无| 主键ID |Primary key ID|


* output 

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":null
}
```

* output 参数说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null| 请求返回的数据 |return data|


### 查询自定义API

*  API: POST /api/{version}/userapi/search/{bk_biz_id}
* API名称：  search_custom_query
* 功能说明：
	* 中文： 查询自定义api
	* English ：search customize query 
*  input body

```
{
     "condition":{
          "name": "aa",
    },
    "start": 0,
    "limit": 20,

}
```


* input参数说明

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | --- |---| --- | --- | --- |
| bk_biz_id|int|是|无|业务ID | business ID |
| condition|object对象|是|无|查询条件 | search condition|
| start|int|是|无|记录开始位置 |start record|
| limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|


condition 参数说明： condition 字段为自定义api的属性字段, 可以是create_user,modify_user, name



* ouput 
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":null,
    "data":{
        "count":1,
        "info":[
            {
                "bk_biz_id":12,
                "create_time":"2018-03-02T15:04:20.117+08:00",
                "create_user":"admin_default",
                "id":"bacfet4kd42325venmcg",
                "info":"{\"condition\":[{\"bk_obj_id\":\"biz\",\"condition\":[{\"field\":\"default\",\"operator\":\"$ne\",\"value\":1}],\"fields\":[]},{\"bk_obj_id\":\"set\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"module\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"host\",\"condition\":[{\"field\":\"bk_host_innerip\",\"operator\":\"$eq\",\"value\":\"127.0.0.1\"}],\"fields\":[\"bk_host_innerip\",\"bk_host_outerip\",\"bk_agent_status\"]}]}",
                "last_time":"",
                "modify_user":"",
                "name":"api1"
            }
        ]
    }
}
```

*  output字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object| 请求返回的数据 |return data|

data字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|--- |
| count| int| 记录条数| count of record|
| info| object array| 自定义api数据| detail of record|

info 字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|--- |
| bk_biz_id|int| 业务ID| business ID|
| create_time|时间格式| 创建时间|create time |
| create_user|string| 创建者|create user|
| id|string| 自定义api主键ID|primary key ID|
| info|json string| 自定义api信息| the query info|
| last_time|string| 更新时间|last update time |
| modify_user|string| 修改者| modify user|
| name|string| 自定义api命名|the name of api|

info 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_obj_id| string| 否| 无|对象名,可以为biz,set,module,host,object|object name, it can be biz,set,module,host,object|
| fields| string数组| 否| 无|查询输出字段|fields output|
| condition| object array| 否| 无|查询条件|search condition|

condition 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| field| string| 否| 无|对象的字段|field of object|
| operator| string| 否| 无|操作符, $eq为相等，$neq为不等，$in为属于，$nin为不属于|$eq is equal,$in is belongs, $nin is not belong,$neq is not equal|
| value| string| 否| 无|字段对应的值|the value of field|


### 获取自定义API详情

*  API: GET /api/{version}/userapi/detail/{bk_biz_id}/{id}
* API名称：  get_custom_query_detail
* 功能说明：
	* 中文： 获取自定义api详情
	* English ：get customize query detail
*  input body
无

* input参数说明

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | --- |---| --- | --- | ---|
| bk_biz_id|int|是|无|业务ID | business ID|
| id|string|是|无|主键ID | pripary key ID|


* ouput 

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":null,
    "data":{
                "biz_id":12,
                "creat_time":"2018-03-02T15:04:20.117+08:00",
                "create_user":"admin_default",
                "id":"bacfet4kd42325venmcg",
                "info":"{\"condition\":[{\"bk_obj_id\":\"biz\",\"condition\":[{\"field\":\"default\",\"operator\":\"$ne\",\"value\":1}],\"fields\":[]},{\"bk_obj_id\":\"set\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"module\",\"condition\":[],\"fields\":[]},{\"bk_obj_id\":\"host\",\"condition\":[{\"field\":\"bk_host_innerip\",\"operator\":\"$eq\",\"value\":\"127.0.0.1\"}],\"fields\":[\"bk_host_innerip\",\"bk_host_outerip\",\"bk_agent_status\"]}]}",
                "last_time":"",
                "modify_user":"",
                "name":"api1"
    }
}
```
*  output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object| 请求返回的数据 |return data|

data字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|--- |
| bk_biz_id|int| 业务ID| business ID|
| create_time|时间格式| 创建时间|create time |
| create_user|string| 创建者|create user|
| id|string| 自定义api主键ID|primary key ID|
| info|json string| 自定义api信息| the query info|
| last_time|string| 更新时间|last update time |
| modify_user|string| 修改者| modify user|
| name|string| 自定义api命名|the name of api|

info 参数说明：

| 名称  | 类型  | 说明 | Description|
| ---  | ---  | --- |---  | 
| bk_obj_id| string|对象名,可以为biz,set,module,host,object|object name, it can be biz,set,module,host,object|
| fields| string数组|查询输出字段|fields output|
| condition| object array|查询条件|search condition|

condition 参数说明：

| 名称  | 类型  | 说明 | Description|
| ---  | ---  | --- | ---|
| field| string|对象的字段|field of object|
| operator| string|操作符, $eq为相等，$neq为不等，$in为属于，$nin为不属于|$eq is equal,$in is belongs, $nin is not belong,$neq is not equal|
| value| string|字段对应的值|the value of field|

### 根据自定义api获取数据

*  API:
GET /api/{version}/userapi/data/{bk_biz_id}/{id}/{start}/{limit}
* API名称：  get_custom_query_data
* 功能说明：
	* 中文： 根据自定义api获取数据
	* English ：get customize query data
*  input body
无
* input参数说明

| 名称  | 类型 |必填| 默认值 | 说明 | Description |
| ---  | --- |---| --- | --- | ---|
| bk_biz_id|int|是|无|业务ID | business ID|
| id|string|是|无|主键ID | primary key ID|
| start|int|是|无|记录开始位置 |start record|
| limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|


* output 

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":null,
    "data":{
        "count":1,
        "info":[
            {
                "biz":{
                    "bk_biz_id":11,
                    "bk_biz_name":"1111",
                    "create_time":"2017-12-20T14:45:22.04+08:00",
                    "default":0,
                    "last_time":"2017-12-20T14:45:22.04+08:00",
                    "bk_biz_maintainer":"tt",
                    "bk_supplier_account":"0",
                    "bk_biz_productor":"tt",
                    "dg":""
                },
                "host":{
                    "bk_host_assetid":"",
                    "bk_comment":"准备下架设备",
                    "create_time":"2018-01-04T14:41:17.376+08:00",
                    "bk_host_id":187,
                    "bk_host_name":"nginx.27",
                    "bk_host_type":"虚拟机",
                    "import_from":"1",
                    "bk_host_innerip":"10.0.0.0",
                    "bk_cloud_id":0,
                    "aaa":"",
                    "cpu":"",
                    "enum2":""
                },
                "module":{
                    "bk_module_name":"空闲机"
                },
                "set":{
                    "bk_set_name":"内置模块集"
                }
            }
        ]
    }
}
```

*  output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object| 请求返回的数据 |return data|

data 字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| count| int| 记录条数 |the num of record|
| info| object array | 主机实际数据 |host data|

info 字段说明:

| 名称  | 类型  | 说明 |Description|
|---|---|---|---| 
| biz | object | 主机所属的业务信息 |host biz info|
| set| object | 主机所属的集群信息 |host set info|
| module| object | 主机所属的模块信息 |host module info|
| host| object | 主机自身属性|host attr info|
