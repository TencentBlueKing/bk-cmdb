###  新增主机到主机资源池
* API: POST /api/v3/hosts/add
* API名称： add_host_to_resource
* 功能说明：
	* 中文：新增主机到资源池
	* English ：add host to resource
* input body：
```
{
　　"host_info":{
　　　　"0":{
　　　　　　"bk_host_innerip":"127.0.0.1",
　　　　　　"import_from":"3",
　　　　　　"bk_cloud_id":1
　　　　}
　　},
　　"bk_biz_id":3
}
```
* input字段说明:

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| host_info| object array| 是|无| 主机信息 | host info|
| bk_biz_id| int| 否| 无| 业务ID |business ID|

host_info object 说明：


| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_host_innerip| string| 是|无| 主机内网ip | host inner ip|
| import_from| string| 是|api| 主机导入来源,以api方式导入为3 | host import source|
| bk_cloud_id| int| 是| 无| 云区域ID|cloud area ID|



* output：
```
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "",
    "data": null
}
```

* output字段说明:

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | json string | 请求返回的数据,示例如下（{"error":["6行字段正则校验失败","7行字段正则校验失败"],"success":["4","5"],"update_error":null}） |the data response|

###  主机转移到业务内模块
* API: POST /api/v3/hosts/modules
* API名称： transfer_host_module
* 功能说明：
	* 中文：业务内主机转移模块
	* English ：transfer host to module in biz
* input body：
```
{
    "bk_biz_id":151,
    "bk_host_id":[
        10,
        9
    ],
    "bk_module_id":[
        170
    ],
    "is_increment":true
}
```
* input字段说明:

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_biz_id| int| 是|无|业务ID |  business ID|
| bk_host_id| int数组| 是| 无|主机 ID|host ID|
| bk_module_id| int数组| 是| 无|模块 id| module ID |
| is_increment| bool| 是| 无|覆盖或者追加,会删除原有关系. true是更新，false是覆盖|cover or pursue ,true will cover |


* output：
```
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "",
    "data": null
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null | 请求返回的数据 |the data response|

### 资源池主机分配至业务的空闲机模块
* API: POST /api/v3/hosts/modules/resource/idle
* API名称： transfer_resourcehost_to_idlemodule
* 功能说明：
	* 中文：  分配资源池主机到业务的空闲机模块
	* English ：transfer resource host to  idle module
* input body：
```
{
  "bk_biz_id": 115, 
  "bk_host_id": [
    10, 
    9
  ]
}
```
* input字段说明:

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_biz_id| int| 是|无|业务ID | host ID|
| bk_host_id| int数组| 是| 无|主机ID|host ID |


* output:
```
{
  "result": true, 
  "bk_error_code": 0, 
  "bk_error_msg": "", 
  "data": null
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---| 
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null | 请求返回的数据 |the data response|

### 主机上交至业务的故障机模块
* API: POST /api/v3/hosts/modules/fault
* API名称： transfer_host_to_faultmodule
* 功能说明：
	* 中文： 上交主机到业务的故障机模块
	* English ：transfer host to  fault module
* input body:
```
{
  "bk_biz_id": 115, 
  "bk_host_id": [
    10, 
    9
  ]
}
```
* input字段说明:

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_biz_id| int| 是|无|业务id | business ID|
| bk_host_id| int数组| 是| 无|主机id| host ID|


* output:
```
{
  "result": true, 
  "bk_error_code": 0, 
  "bk_error_msg": "", 
  "data": null
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null | 请求返回的数据 |the data response|


### 主机上交至业务的空闲机模块
* API: POST /api/v3/hosts/modules/idle
* API名称：   transfer_host_to_idlemodule
* 功能说明：
	* 中文：上交主机到业务的空闲机模块
	* English ：transfer host to idle module
* input boy:
```
{
  "bk_biz_id": 115, 
  "bk_host_id": [
    10, 
    9
  ]
}
```
* input字段说明

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_biz_id| int| 是|无|业务id | business ID|
| bk_host_id| int数组| 是| 无|主机id| host ID|


* output:
```
{
  "result": true, 
  "bk_error_code": 0, 
  "bk_error_msg": "", 
  "data": null
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null | 请求返回的数据 |the data response|

### 主机回收至资源池
* API:  POST /api/v3/hosts/modules/resource
* API名称： transfer_host_to_resourcemodule
* 功能说明：
	* 中文：上交主机至资源池
	* English ：transfer host to resource module
* input boy:
* input:
```
{
"bk_biz_id":269,
"bk_host_id":[204]
}
```
* input字段说明:

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_biz_id| int| 是|无|业务id | business ID|
| bk_host_id| int数组| 是| 无|主机id| host ID|


* output:
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":null,
    "data":""
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null | 请求返回的数据 |the data response|

### 转移主机至模块
* API:  POST /api/v3/hosts/modules/biz/mutilple
* API名称： transfer_host_to_mutiple_biz_modules
* 功能说明：
	* 中文：转移主机至模块，此模块属于不同业务
	* English ：transfer host to module,this module belongs to different business
* input boy:
* input:
```
{
    "bk_biz_id":10,
    "bk_module_id":58,
    "host_info":[
     {
       "bk_host_innerip":"10.20.30.45",
       "bk_cloud_id":0
    }]
}
```
* input字段说明:

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_biz_id| int| 是|无|业务id | business ID|
| bk_host_id| int array| 是| 无|主机id| host ID|
| host_info| object array| 是| 无|主机信息| 主机信息数组|

host_info说明：

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_host_innerip| string| 是|无|主机内网ip | host inner ip|
| bk_cloud_id| int | 是| 无|云区域ID| cloud area ID|

* output:
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"success",
    "data":"sucess"
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null | 请求返回的数据 |the data response|

### 清空集群/模块下主机
* API:  POST /api/v3/hosts/modules/idle/set
* API名称： transfer_sethost_to_idle_module
* 功能说明：
	* 中文：清空业务下集群/模块中主机
	* English ：transfer host in set or module to idle module
* input boy:
* input:
```
{
    "bk_biz_id":10,
    "bk_module_id":58,
    "bk_set_id":1
}
```
* input字段说明:

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_biz_id| int| 是|无|业务id | business ID|
| bk_set_id| int | 是| 无|集群id| set ID|
| bk_module_id| int| 是| 无|模块id| module ID|


* output:
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"success",
    "data":"sucess"
}
```

* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null | 请求返回的数据 |the data response|




###  同步新增主机到业务
* API: POST /api/v3/hosts/sync/new/host
* API名称： 
* 功能说明：
	* 中文：同步新增主机到业务
	* English ：add new hosts to the service
* input body：
```
{
　　"host_info":{
　　　　"0":{
　　　　　　"bk_host_innerip":"127.0.0.1",
　　　　　　"bk_cloud_id":1
　　　　}
　　},
　　"bk_biz_id": 3,
    "bk_module_id":[7,8],
}
```
* input字段说明:
 
| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| host_info| object array| 是|无| 主机信息 | host info|
| bk_biz_id| int| 是| 无| 业务 ID|bussines ID|
| bk_module_id|int array| 是| 无| 模块 ID|模块 ID|
 	 
 	 
host_info object 说明：
 
 
| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_host_innerip| string| 是|无| 主机内网ip | host inner ip|
| bk_cloud_id| int| 是| 无| 云区域ID|cloud ID|
|...|主机其他属性|
 
 
 
* output：
```
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "success",
    "data": {
        "success": [
            "4"
        ]
    }
}
```
 
* output字段说明:
 
| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | null | 请求返回的数据 |the data response|


###  获取业务下主机与集群，模块的关系
*  API: POST /api/v3/host/topo/relation/read
* API名称： search_app_host_relation_info
* 功能说明：
  * 中文：获取业务下主机与集群，集群的关系
  * English ：Obtain the relationship between the host and the cluster and the cluster under the service
* input body：
```
{
    "page":{
        "start":0,
        "limit":10,
        "sort":"bk_host_id"
    },
    "bk_biz_id":2,
    "bk_set_ids": [1, 2],
    "bk_module_ids": [23, 24],
    "bk_host_ids": [25, 26]
}
```

* input参数说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---| 
| bk_biz_id| int| 是| 无|业务ID|biz condition|
| bk_set_ids|array | 否| 无|集群ID列表| set ID list|
| bk_module_ids|array | 否| 无|模块ID列表| module ID list|
| bk_host_ids|array | 否| 无|主机ID列表| host ID list|

page 字段说明：

| 名称  | 类型 |必填| 默认值 | 说明 | Description|
| ---  | ---  | --- |---  | --- | ---| 
|start|int|否|0|获取数据偏移位置|get the data offset location|
|limit|int|否|200|过去数据条数限制，最大值为500|The number of data points in the past is limited to a maximum of 500|


* output 

```
{
  "result": true,
  "bk_error_code": 0,
  "bk_error_msg": "success",
  "permission": null,
  "data": {
    "count": 10,
    "info": [
        {
        "bk_biz_id": 3,
        "bk_host_id": 5,
        "bk_module_id": 54,
        "bk_set_id": 10,
        "bk_supplier_account": "0"
        },
        .....
    ],
    "page":{
       "start":0,
       "limit":10
    }
   
  }
}
```

* output 参数说明


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
| info| object array | 业务下主机与集群，模块，集群的数据详情列表 |business details list of hosts and clusters, modules, clusters|

data.info[n] 字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| bk_biz_id | int | 业务ID | business ID|
| bk_set_id | int | 集群ID | set ID|
| bk_module_id | int | 模块ID | module ID|
| bk_host_id | int | 主机ID | host ID|



page 字段说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
|start|int|服务端本次获取数据偏移位置|server obtains the data offset position this time|
|limit|int|服务端返回数据条数限制|server client return data limit，|
