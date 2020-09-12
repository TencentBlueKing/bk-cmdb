
### 保存用户字段配置
* API:  POST /api/{version}/usercustom
* API名称： save_user_custom
* 功能说明：
	* 中文：保存用户自定义配置
	* English ：save user custom
* input body:
```
{
    "host_query_column":"["bk_host_innerip", "bk_host_name"]",
    "host_display_column":"["bk_host_innerip", "bk_host_name"]",
    "biz_query_column":"["bk_biz_name", "bk_biz_productor", "bk_biz_maintainer"}",
    "biz_display_column":"{"bk_biz_name", "bk_biz_productor", "bk_biz_maintainer"}",
    "bk_biz_id":123
}
```


* input字段说明：

| 名称                | 类型        | 必填 | 默认值 | 说明         | Description              |
| ------------------- | ----------- | ---- | ------ | ------------ | ------------------------ |
| host_query_column   | string数组  | 否   | 无     | 主机查询字段 | host query fields        |
| host_display_column | string数组  | 否   | 无     | 主机展示字段 | host display fields      |
| biz_query_column    | string 数组 | 否   | 无     | 业务查询字段 | business query fields    |
| biz_display_column  | string数组  | 否   | 无     | 业务展示字段 | business  display fields |
| bk_biz_id           | int         | 否   | 无     | 业务ID       | business ID              |


* output:

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":null
}
```
*  output字段说明

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result                                             |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | null   | 请求返回的数据                             | return data                                                |


### 获取用户字段配置
* API: POST /api/{version}/usercustom/user/search
* API名称： search_user_custom
* 功能说明：
	* 中文：获取用户自定义配置
	* English ：save user custom
* input body:
无
* output:

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":{
        "bk_biz_id":123,
        "host_query_column":"{"bk_host_innerip", "bk_host_name"}",
        "host_display_column":"{"bk_host_innerip", "bk_host_name"}",
        "biz_query_column":"{"bk_biz_name", "bk_bus_productor", "bk_biz_maintainer"}",
        "biz_display_column":"{"bk_biz_name", "bk_bus_productor", "bk_biz_maintainer"}",
        "id":"b81ervtmjrcduf67mm9g",
        "user":"test default"
    }
}

```

* output参数说明：

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result                                             |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | object | 请求返回的数据                             | return data                                                |

data字段说明：

| 名称                | 类型        | 说明         | Description             |
| ------------------- | ----------- | ------------ | ----------------------- |
| host_query_column   | string数组  | 主机查询字段 | host query fields       |
| host_display_column | string数组  | 主机展示字段 | host display fields     |
| biz_query_column    | string 数组 | 业务查询字段 | business query fields   |
| biz_display_column  | string数组  | 业务展示字段 | business display fields |
| bk_biz_id           | int         | 业务ID       | business ID             |
| id                  | int         | 主键ID       | primary key ID          |
| user                | string      | 用户名       | user name               |

### 获取默认字段配置
* API:  POST /api/{version}/usercustom/default/search
* API名称： get_user_default_custom
* 功能说明：
	* 中文：获取默认用户自定义配置
	* English ：save default user custom
* input body:
无

* output:

```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":{
        "host_query_column":"{"bk_host_innerip", "bk_host_name"}",
        "host_display_column":"{"bk_host_innerip", "bk_host_name"}",
        "biz_query_column":"{"bk_biz_name", "bk_bus_productor", "bk_biz_maintainer"}",
        "biz_display_column":"{"bk_biz_name", "bk_bus_productor", "bk_biz_maintainer"}",
        "id":"b81ervtmjrcduf67mm9g",
        "user":"test default",
        "is_default":1,
    }
}
```


* output参数说明：

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result                                             |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | object | 请求返回的数据                             | return data                                                |

data字段说明：

| 名称                | 类型        | 说明         | Description             |
| ------------------- | ----------- | ------------ | ----------------------- |
| host_query_column   | string数组  | 主机查询字段 | host query fields       |
| host_display_column | string数组  | 主机展示字段 | host display fields     |
| biz_query_column    | string 数组 | 业务查询字段 | business query fields   |
| biz_display_column  | string数组  | 业务展示字段 | business display fields |
| bk_biz_id           | int         | 业务ID       | business ID             |
| id                  | int         | 主键ID       | primary key ID          |
| user                | string      | 用户名       | user name               |


