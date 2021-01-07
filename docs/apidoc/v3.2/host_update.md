


###  更新主机属性

* API: PUT /api/{version}/hosts/batch
* API名称： update_host
* 功能说明：
	* 中文：更新主机属性
	* English ：update host attributes
* input body：
```
{
    "bk_host_id":"1,2,3"   ,
    "bk_host_name":"test"
}
```

* input参数说明：

| 名称       | 类型   | 必填 | 默认值 | 说明                            | Description         |
| ---------- | ------ | ---- | ------ | ------------------------------- | ------------------- |
| bk_host_id | string | 是   | 无     | 主机id,int类型的bk_host_id,分割 | host id join by "," |
需要更新的主机属性见主机属性列表

* output:
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":"success"
}
```

* 字段说明

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result                                             |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | string | 请求返回的数据                             | return data                                                |

