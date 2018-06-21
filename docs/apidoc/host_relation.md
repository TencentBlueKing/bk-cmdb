
###  新增主机
* API: POST /api/{version}/hosts/add
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
　　　　　　"import_from":"1",
　　　　　　"bk_cloud_id":1
　　　　}
　　},
　　"bk_supplier_id":0
}
```
* input字段说明:

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| host_info| object array| 是|无| 主机信息 | host info|
| bk_supplier_id| int| 是| 无| 开发商 ID|supplier ID|


host_info object 说明：


| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_host_innerip| string| 是|无| 主机内网ip | host inner ip|
| import_from| string| 是|api| 主机导入来源 | host import source|
| bk_cloud_id| int| 是| 无| 云区域ID|cloud ID|



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
| data | null | 请求返回的数据 |the data response|

###  主机转移到业务内模块
* API: POST /api/{version}/hosts/modules
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
* API: POST /api/{version}/hosts/modules/resource/idle
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
* API: POST /api/{version}/hosts/modules/fault
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
* API: POST /api/{version}/hosts/modules/idle
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
* API:  POST /api/{version}/hosts/modules/resource
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
* API:  POST /api/{version}/hosts/modules/biz/mutilple
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
| bk_cloud_id| int | 是| 无|主机云区域ID| host cloud ID|

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
