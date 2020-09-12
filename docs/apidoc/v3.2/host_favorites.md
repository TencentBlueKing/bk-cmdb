
### 新加收藏
* API: POST /api/{version}/hosts/favorites
* API名称： create_favorites
* 功能说明：
	* 中文：添加收藏
	* English ：create favorites
* input body：
```
{
    "info":"{\"biz\":{\"bk_biz_id\":12},\"exact_search\":false,\"bk_host_innerip\":true,\"bk_host_outerip\":true,\"ip_list\":[]}",
    "query_params":"[{\"bk_obj_id\":\"host\",\"field\":\"operator\",\"operator\":\"$eq\",\"value\":\"admin\"}]",
    "name":"my5"
}
```

* input参数说明

| 名称         | 类型        | 必填 | 默认值 | 说明         | Description                    |
| ------------ | ----------- | ---- | ------ | ------------ | ------------------------------ |
| info         | json string | 否   | 无     | ip查询条件   | ip search parameters           |
| query_params | json string | 否   | 无     | 通用查询条件 | common search query parameters |
| name         | string      | 是   | 无     | 收藏的名称   | the name of favorites          |

info 参数说明：

| 名称            | 类型       | 必填 | 默认值 | 说明             | Description              |
| --------------- | ---------- | ---- | ------ | ---------------- | ------------------------ |
| biz             | object     | 是   | 无     | 业务信息查询条件 | business info for search |
| exact_search    | bool       | 是   | 无     | 是否精确查询     | is exeact search         |
| bk_host_innerip | bool       | 是   | 无     | true 或者false   | true or false            |
| bk_host_outerip | bool       | 是   | 无     | true 或者false   | true or false            |
| ip_list         | string数组 | 是   | 无     | ip地址列表       | ip address list          |

biz 参数信息：

| 名称      | 类型 | 必填 | 默认值 | 说明   | Description |
| --------- | ---- | ---- | ------ | ------ | ----------- |
| bk_biz_id | int  | 是   | 无     | 业务ID | business ID |

query_params 参数说明:

| 名称      | 类型   | 必填 | 默认值 | 说明                                                   | Description                                                       |
| --------- | ------ | ---- | ------ | ------------------------------------------------------ | ----------------------------------------------------------------- |
| bk_obj_id | string | 是   | 无     | 对象ID                                                 | object ID                                                         |
| field     | string | 否   | 无     | 对象的字段                                             | field of object                                                   |
| operator  | string | 否   | 无     | 操作符, $eq为相等，$neq为不等，$in为属于，$nin为不属于 | $eq is equal,$in is belongs, $nin is not belong,$neq is not equal |
| value     | object | 否   | 无     | 字段对应的值                                           | the value of field                                                |


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

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result                                             |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | object | 请求返回的数据                             | return data                                                |

data 结构说明：

| 名称 | 类型   | 说明         | Description              |
| ---- | ------ | ------------ | ------------------------ |
| id   | string | 收藏的主键ID | favorites primary key ID |



### 编辑收藏
* API:  PUT /api/{version}/hosts/favorites/{id} 
* API名称： update_favorites
* 功能说明：
	* 中文：编辑收藏
	* English ：update favorites
* input body：
```
{
    "count":6,
    "id":"bacb3j4kd42325venmag",
    "info":"{\"biz\":{\"id\":bk_biz_id},\"exact\":0,\"bk_host_innerip\":true,\"bk_host_outerip\":true,\"ip\":[]}",
    "is_default":2,
    "name":"my211",
    "query_params":"[{\"bk_biz_id\":12,\"bk_obj_id\":\"biz\",\"field\":\"default\",\"operator\":\"$ne\",\"value\":1}]"
}
```

* input参数说明


| 名称         | 类型        | 必填 | 默认值 | 说明         | Description                    |
| ------------ | ----------- | ---- | ------ | ------------ | ------------------------------ |
| info         | json string | 否   | 无     | ip查询条件   | ip search parameters           |
| query_params | json string | 否   | 无     | 通用查询条件 | common search query parameters |
| name         | string      | 是   | 无     | 收藏的名称   | the name of favorites          |
| id           | string      | 是   | 无     | 收藏的主键   | favorites primary key ID       |
| count        | int         | 否   | 无     | 收藏次数     | the time of favorites          |

info 参数说明：

| 名称            | 类型       | 必填 | 默认值 | 说明             | Description              |
| --------------- | ---------- | ---- | ------ | ---------------- | ------------------------ |
| biz             | object     | 是   | 无     | 业务信息查询条件 | business info for search |
| exact_search    | bool       | 是   | 无     | 是否精确查询     | is exeact search         |
| bk_host_innerip | bool       | 是   | 无     | true 或者false   | true or false            |
| bk_host_outerip | bool       | 是   | 无     | true 或者false   | true or false            |
| ip_list         | string数组 | 是   | 无     | ip地址列表       | ip address list          |

biz 参数信息：

| 名称      | 类型 | 必填 | 默认值 | 说明   | Description |
| --------- | ---- | ---- | ------ | ------ | ----------- |
| bk_biz_id | int  | 是   | 无     | 业务ID | business ID |

query_params 参数说明:

| 名称      | 类型   | 必填 | 默认值 | 说明                                                   | Description                                                       |
| --------- | ------ | ---- | ------ | ------------------------------------------------------ | ----------------------------------------------------------------- |
| bk_obj_id | string | 是   | 无     | 对象ID                                                 | object ID                                                         |
| field     | string | 否   | 无     | 对象的字段                                             | field of object                                                   |
| operator  | string | 否   | 无     | 操作符, $eq为相等，$neq为不等，$in为属于，$nin为不属于 | $eq is equal,$in is belongs, $nin is not belong,$neq is not equal |
| value     | object | 否   | 无     | 字段对应的值                                           | the value of field                                                |


* output
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


### 删除收藏
* API: DELETE /api/{version}/hosts/favorites/{id} 
* API名称： delete_favorites
* 功能说明：
	* 中文：删除收藏
	* English ：delete favorites
* input body：
无

* output
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"",
    "data":null
}
```

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result                                             |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | null   | 请求返回的数据                             | return data                                                |


### 获取收藏列表
* API: POST /api/{version}/hosts/favorites/search 
* API名称： search_favorites
* 功能说明：
	* 中文：获取收藏列表
	* English ：search favorites
* input body：

```
{
    "condition": {
        "is_default": 1,
        "name": "保存名称" 
    },
    "limit": 10,
    "start": 0
}
```
* input参数说明

| 名称      | 类型   | 必填 | 默认值 | 说明                 | Description            |
| --------- | ------ | ---- | ------ | -------------------- | ---------------------- |
| condition | object | 是   | 无     | 查询条件对象         | search condition       |
| start     | int    | 是   | 无     | 记录开始位置         | start record           |
| limit     | int    | 是   | 无     | 每页限制条数,最大200 | page limit, max is 200 |



* output
```
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":null,
    "data":{
        "count":1,
        "info":[
{"count":1,"id":"bacb3j4kd42325venmag","info":"{\"biz\":{\"bk_biz_id\":12},\"exact\":0,\"bk_host_innerip\":true,\"bk_host_outerip\":true,\"ip\":[]}","is_default":2,"name":"my2","bk_query_params":"[{\"bk_biz_id\":12,\"bk_obj_id\":\"biz\",\"field\":\"Default\",\"operator\":\"$ne\",\"value\":1}]"}
        ]
    }
}
```

* output 字段说明：

| 名称          | 类型   | 说明                                       | Description                                                |
| ------------- | ------ | ------------------------------------------ | ---------------------------------------------------------- |
| result        | bool   | 请求成功与否。true:请求成功；false请求失败 | request result                                             |
| bk_error_code | int    | 错误编码。 0表示success，>0表示失败错误    | error code. 0 represent success, >0 represent failure code |
| bk_error_msg  | string | 请求失败返回的错误信息                     | error message from failed request                          |
| data          | object | 请求返回的数据                             | return data                                                |

 data：
 
| 名称  | 类型   | 说明         | Description           |
| ----- | ------ | ------------ | --------------------- |
| count | int    | 请求记录条数 | num of record         |
| info  | object | 请求记录信息 | the info of favorites |
info object说明：为添加收藏的存储

### 收藏使用次数加一
* API: PUT /api/{version}/hosts/favorites/{id}/incr  
* API名称： incr_favorites
* 功能说明：
	* 中文：收藏使用次数自增长
	* English ：add favorites use times
* input body：
无

* input参数说明：


| 名称 | 类型   | 必填 | 默认值 | 说明         | Description    |
| ---- | ------ | ---- | ------ | ------------ | -------------- |
| id   | string | 是   | 无     | 收藏的主键ID | primary key ID |


* output
```
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "",
    "data": {
        "count": 3, 
        "id": "b81gpe04m7vhbr71qlk0" 
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

| 名称  | 类型   | 说明         | Description    |
| ----- | ------ | ------------ | -------------- |
| id    | string | 收藏的主键ID | primary key ID |
| count | int    | 收藏使用次数 | used times     |



