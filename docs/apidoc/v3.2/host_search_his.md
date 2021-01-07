
### 新加主机查询历史
* API: POST /api/{version}/hosts/history
* API名称： create_host_search_history
* 功能说明：
	* 中文：主机查询历史
	* English ：create host search history
* input body：
```
{
    "content":"{"bk_host_id":"10"}"
}
```


* input参数说明


| 名称    | 类型   | 必填 | 默认值 | 说明         | Description                   |
| ------- | ------ | ---- | ------ | ------------ | ----------------------------- |
| content | object | 是   | 无     | 主机查询条件 | host search condition content |
content为对象

* output
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":{
        "id":"b8b4bglmjrcav115n31g"
    }
}
```
*  output字段说明：


| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result                                             |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | object | 请求返回的数据                             | return data                                                |

data 字段说明：

| 名称 | 类型   | 说明             | Description                  |
| ---- | ------ | ---------------- | ---------------------------- |
| id   | string | 查询记录的主键ID | search record primary key ID |

### 获取主机查询历史
* API: GET /api/{version}/hosts/history/{start}/{limit} 
* API名称： search_host_search_history
* 功能说明：
	* 中文：主机查询历史
	* English ：search host search history
* input body：
无

* input参数说明

| 名称  | 类型 | 必填 | 默认值 | 说明                 | Description            |
| ----- | ---- | ---- | ------ | -------------------- | ---------------------- |
| start | int  | 是   | 无     | 记录开始位置         | start record           |
| limit | int  | 是   | 无     | 每页限制条数,最大200 | page limit, max is 200 |



* output
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":{
        "count":5,
        "info":[
            {
                "content":"test",
                "create_time":"2017-11-23T10:42:07.754+08:00",
                "id":"b8b3avtmjrcav115n2vg",
                "user":"test default"
            },
            {
                "content":"test",
                "create_time":"2017-11-23T10:42:08.89+08:00",
                "id":"b8b3b05mjrcav115n300",
                "user":"test default"
            }
        ]
    }
}
```
*  output字段说明

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result                                             |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | object | 请求返回的数据                             | return data                                                |

data 字段说明：

| 名称  | 类型         | 说明         | Description     |
| ----- | ------------ | ------------ | --------------- |
| count | int          | 返回数据条数 | record num      |
| info  | object array | 返回数据     | the data return |

info 字段说明：

| 名称        | 类型         | 说明     | Description       |
| ----------- | ------------ | -------- | ----------------- |
| content     | object array | 查询内容 | search content    |
| create_time | 时间格式     | 创建时间 | create time       |
| id          | string       | 主键ID   | primary key ID    |
| user        | string       | 创建人   | the man create it |
