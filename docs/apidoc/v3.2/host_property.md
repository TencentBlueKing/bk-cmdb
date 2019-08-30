###  克隆主机属性

* API: PUT /api/{version}/hosts/propery/clone
* API名称： clone_host_property
* 功能说明：
	* 中文：克隆主机属性
	* English ：clone host property 
	
* input body：
```
{
    "bk_biz_id":2,
    "bk_org_ip":"127.0.0.1",
    "bk_dst_ip":"127.0.0.2",
    "bk_cloud_id":0
}
```

* input字段说明

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
| bk_org_ip| string| 是|无| 源主机ip, 只支持传入单ip |origin host ip ,only support single|
| bk_dst_ip| string| 是| 无|目标主机ip, 多个ip用","分割|destination host ip,  multiple ip splits with "," |
| bk_biz_id|int|是|无| 业务ID |business ID|
| bk_cloud_id| int| 否| 无| 云区域ID|cloud area ID|

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
