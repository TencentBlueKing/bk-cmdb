### 新增进程
* API： POST /api/{version}/proc/{bk_supplier_account}/{bk_biz_id}
* API名称： create_process
* 功能说明：
	* 中文：创建进程
	* English ：create process
* input body：
```
{
    "bk_process_name":"nginx",
    "port":80,
    "bind_ip":"1",
    "protocol":"1",
    "bk_func_name":"nginx",
    "work_path":"/data/cc/running",
    "user":"cc"
}
```

* input字段说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_process_name| string| 是|无|进程名 |process name|
| port|  string| 是| 无|主机端口|host port|
|protocol|string|协议:1/2(1:tcp, 2:udp)|protocol:1/2(1:tcp, 2:udp)|
|bind_ip|string|绑定IP:1/2/3/4(1:127.0.0.1,2:0.0.0.0,3:第一内网IP,4:第一外网IP)|1/2/3/4(1:127.0.0.1,2:0.0.0.0,3:first intranet IP,4:first extranet IP)|

 其它字段依赖变量的定义

* output：
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":"success"
}
```

* output 字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | string | 请求返回的数据 |return data|


### 查询进程

* API:  POST /api/{version}/proc/search/{bk_supplier_account}/{bk_biz_id}
* API名称： search_process
* 功能说明：
	* 中文：查询进程
	* English ：search process

* input body :

```
{
    "page":{
        "start":0,
        "limit":10,
        "sort":"bk_process_name"
    },
    "fields":[
        "bk_process_id",
        "bk_process_name"
    ],
    "condition":{
        "bk_biz_id":"12233",
        "bk_process_name":"nginx"
    }
}

```


* input字段说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| page| object| 是|无|分页参数 |page parameter|
| fields| array | 是| 无|查询字段|search fields|
| condition|  object| 是| 无|查询条件|search condition|

page 参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| start|int|是|无|记录开始位置 |start record|
| limit|int|是|无|每页限制条数,最大200 |page limit, max is 200|
| sort| string| 否| 无|排序字段|the field for sort|

fields参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_process_id| int| 否|无|进程ID |process id|
| bk_process_name| string| 否|无|进程名称 |process name|

参数为进程的任意属性

condition 参数说明：condition 参数为进程的属性

* output:
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":{
        "count":5,
        "info":[
            {
                "bk_process_name":"nginx",
                "port":80,
                "bind_ip":"1",
                "protocol":"1",
                "bk_func_name":"nginx",
                "work_path":"/data/cc/running",
                "user":"cc"
            },
            {
                "bk_process_name":"apache",
                "port":8080,
                "bind_ip":"1",
                "protocol":"1",
                "bk_funcName":"nginx",
                "work_path":"/data/cc/running",
                "bk_user":"cc"
            }
        ]
    }
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object | 请求返回的数据 |the data response|



data 数据信息：

| 名称  | 类型  | 说明 |request result true or false|
|---|---|---|---|
| count| int | 请求失败返回的错误信息 |the count of data|
| info| object | 请求返回的数据 |list of process|

info字段说明：
### 获取进程详情

* API: GET    /api/{version}/proc/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}
* API名称： get_process_detail
* 功能说明：
	* 中文：获取进程详情
	* English ：get process detail
* input body:
不需要
* input字段说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_supplier_account| string| 是|无|开发商 code |supplier account code|
| bk_biz_id| int | 是| 无|业务 id|business id |
| bk_process_id|  int| 是| 无|进程 id |process id|


* output:
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":[
        {
            "bk_property_id":"bk_process_name",
            "bk_property_name":"进程名",
            "bk_property_value":"nginx"
        },
        {
            "bk_property_id":"bk_process_name",
            "bk_property_name":"功能名",
            "bk_property_value":"nginx"
        }
    ]
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object | 请求返回的数据 |the data response|

data 数据说明： 进程属性的具体数据

### 删除进程

* API: DELETE    /api/{version}/proc/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}
* API名称： delete_process
* 功能说明：
	* 中文：删除进程
	* English ：delete process
* input body:
不需要

* output:
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":"success"
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | string | 请求返回的数据 |the data response|

### 更新进程
* API:  PUT  /api/{version}/proc/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}
* API名称： update_process
* 功能说明：
	* 中文：更新进程
	* English ：update process
* input body:
```
{
    "bk_process_name":"nginx"
}
```

* input字段说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_supplier_account| string| 是|无|开发商 code |supplier account code|
| bk_biz_id| int | 是| 无|业务 id|business id |
| bk_process_id|  int| 是| 无|进程 id |process id|
body 字段为进程属性


* output:

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":"success"
}
```

* output字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | string | 请求返回的数据 |the data response|


### 批量更新进程
* API:  PUT  /api/{version}/proc/{bk_supplier_account}/{bk_biz_id}
* API名称： batch_update_process
* 功能说明：
	* 中文：批量更新进程
	* English ：batch update process
* input body:
```
{
    "bk_process_id" : "44,45,46,47,48",
	"start_cmd": "./start.sh 8080",
	"port": "1000"
}
```

* input字段说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_supplier_account| string| 是|无|开发商 code |supplier account code|
| bk_biz_id| int | 是| 无|业务 id|business id |
| bk_process_id|  string| 是| 无|进程id,int类型的bk_process_id,分割|process ids joined by ','|
body 字段为进程属性，可指定除`bk_func_id`和`bk_process_name`以外的属性


* output:

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":"success"
}
```

* output字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | string | 请求返回的数据 |the data response|


### 获取进程绑定模块
* API: GET    /api/{version}/proc/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}
* API名称： get_process_bind_module
* 功能说明：
	* 中文：获取进程绑定的模块
	* English ：get the process bind module
* input body:
无

* input字段说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_supplier_account| string| 是|无|开发商 code |supplier account code|
| bk_biz_id| int | 是| 无|业务 id|business id |
| bk_process_id|  int| 是| 无|进程 id |process id|


* output：
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":[
        {
            "bk_module_name":"db",
            "set_num":10,
            "is_bind":0
        },
        {
            "bk_module_name":"gs",
            "set_num":5,
            "is_bind":1
        }
    ]
}
```

* output字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object| 请求返回的数据 |the data response|

data 数据结构

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| bk_module_name| string| 模块名 |module name|
| set_num| int | 属于几个集群 | bind set num |
| is_bind| int| 是否绑定模块 |is bind to module|

### 绑定进程到模块
* API: PUT   /api/{version}/proc/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}/{bk_module_name}
* API名称： bind_process_module
* 功能说明：
	* 中文：绑定进程到模块
	* English ：bind process to module
* input body :
* 
无

* input字段说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_supplier_account| string| 是|无|开发商 code |supplier account code|
| bk_biz_id| int | 是| 无|业务 id|business id |
| bk_process_id|  int| 是| 无|进程 id |process id|
| bk_module_name|  string| 是| 无|模块名称 |module name|


* output:
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":"success"
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | string| 请求返回的数据 |the data response|


### 解绑进程模块
* API: DELETE   /api/{version}/proc/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}/{bk_module_name}
* API名称： delete_process_module_binding
* 功能说明：
	* 中文： 删除进程模块绑定关系
	* English ：delete process  module binding relationship
* input body:
无
* input字段说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_supplier_account| string| 是|无|开发商 code |supplier account code|
| bk_biz_id| int | 是| 无|业务 ID|business id |
| bk_process_id|  int| 是| 无|进程 ID |process id|
| bk_module_name|  string| 是| 无|模块名称 |module name|


* output:
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":"success"
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | string| 请求返回的数据 |the data response|


