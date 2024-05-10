### 描述

获取动态分组详情 (版本：v3.9.6，权限：业务访问权限)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述         |
|-----------|--------|----|------------|
| bk_biz_id | int    | 是  | 业务ID       |
| id        | string | 是  | 目标动态分组主键ID |

### 调用示例

```json
{
    "bk_biz_id": 1,
    "id": "XXXXXXXX"
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": {
    	"bk_biz_id": 1,
    	"name": "my-dynamic-group",
    	"id": "XXXXXXXX",
    	"bk_obj_id": "host",
    	"info": {
    		"condition":[
    			{
    				"bk_obj_id":"set",
    				"condition":[
    					{
    						"field":"default",
    						"operator":"$ne",
    						"value":1
    					}
    				]
    			},
    			{
    				"bk_obj_id":"module",
    				"condition":[
    					{
    						"field":"default",
    						"operator":"$ne",
    						"value":1
    					}
    				]
    			},
    			{
    				"bk_obj_id":"host",
    				"condition":[
    					{
    						"field":"bk_host_innerip",
    						"operator":"$eq",
    						"value":"127.0.0.1"
    					}
    				]
    			}
    		]
    	},
       "create_user": "admin",
       "create_time": "2018-03-27T16:22:43.271+08:00",
       "modify_user": "admin",
       "last_time": "2018-03-27T16:29:26.428+08:00"
    },
    "permission": null,
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |

#### data

| 参数名称        | 参数类型   | 描述                          |
|-------------|--------|-----------------------------|
| bk_biz_id   | int    | 业务ID                        |
| id          | string | 动态分组主键ID                    |
| bk_obj_id   | string | 动态分组的目标资源对象类型,目前可以为host,set |
| name        | string | 动态分组命名                      |
| info        | object | 动态分组规则信息                    |
| last_time   | string | 更新时间                        |
| modify_user | string | 修改者                         |
| create_time | string | 创建时间                        |
| create_user | string | 创建者                         |

#### data.info.condition

| 参数名称      | 参数类型   | 描述                                                                                    |
|-----------|--------|---------------------------------------------------------------------------------------|
| bk_obj_id | string | 条件对象资源类型, host类型的动态分组支持的info.conditon:set,module,host；set类型的动态分组支持的info.condition:set |
| condition | array  | 查询条件                                                                                  |

#### data.info.condition.condition

| 参数名称     | 参数类型   | 描述                                     |
|----------|--------|----------------------------------------|
| field    | string | 对象的字段                                  |
| operator | string | 操作符, op值为eq(相等)/ne(不等)/in(属于)/nin(不属于) |
| value    | object | 字段对应的值                                 |
