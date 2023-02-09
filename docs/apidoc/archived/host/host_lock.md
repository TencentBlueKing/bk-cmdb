
###  新加主机锁

* API: POST /api/v3/host/lock
* API名称： add_host_lock
* 功能说明：
	* 中文：新加主机锁，如果主机已经加过锁，同样提示加锁成功
	* English：add host lock. If the host has been locked, the same prompt is successful.
* input body：
```
{
   "id_list":[1, 2, 3]
}
```
* input字段说明

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
|id_list| string| 是|无| 主机ID列表| host id list|


* output:
```
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "success",
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


###  查询主机锁

* API: POST /api/v3/host/lock/search
* API名称： search_host_lock
* 功能说明：
	* 中文：查询主机锁
	* English: search host lock. 
* input body：
```
{
   "id_list":[1, 2]
}
```
* input字段说明

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
|id_list| string| 是|无| 主机ID列表| host id list|




* output:
```
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "success",
    "data": {
        1: true,
        2: false
    }
}
```


* output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result true or false|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | map | 请求返回的数据, key 是 ID，value 是否上锁 |the data response,Key is the ID, value is locked status|




###  删除主机锁

* API: DELETE /api/v3/host/lock
* API名称： delete_host_lock
* 功能说明：
	* 中文：删除主机锁
	* English：delete host lock
* input body：
```
{
   "id_list":[1, 2, 3]
}
```
* input字段说明

| 名称  | 类型 |必填| 默认值 | 说明 |Description|
| ---  | ---  | --- |---  | --- | ---|
|id_list| string| 是|无| 主机ID列表| host id list|


* output:
```
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "success",
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


