
###  删除主机

* API: DELETE /api/{version}/hosts/batch
* API名称： delete_host
* 功能说明：
	* 中文：删除主机
	* English ：delete host 
* input body：
```
{
    "bk_host_id":"235,236",
    "bk_supplier_account":"0"
}
```
* input字段说明

| 名称                | 类型   | 必填 | 默认值 | 说明       | Description           |
| ------------------- | ------ | ---- | ------ | ---------- | --------------------- |
| bk_host_id          | string | 是   | 无     | 主机id     | host id join by","    |
| bk_supplier_account | string | 是   | 无     | 开发商账号 | supplier account code |


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

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result true or false                               |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | null   | 请求返回的数据                             | the data response                                          |
