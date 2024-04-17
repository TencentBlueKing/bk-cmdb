### 功能描述

创建动态分组 (版本：v3.9.6+，权限：动态分组新建权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段        | 类型     | 必选 | 描述                          |
|-----------|--------|----|-----------------------------|
| bk_biz_id | int    | 是  | 业务ID                        |
| bk_obj_id | string | 是  | 动态分组的目标资源对象类型,目前可以为host,set |
| info      | object | 是  | 通用查询条件                      |
| name      | string | 是  | 动态分组名称                      |

#### info
| 字段        | 类型     | 必选 | 描述                   |
|-----------|--------|----|----------------------|
| condition | object | 否  | 动态分组锁定条件, 和可变条件至少传一个 |
| variable_condition | object | 否  | 动态分组可变条件, 和锁定条件至少传一个 |

#### info.condition

| 字段        | 类型     | 必选 | 描述                                                                                    |
|-----------|--------|----|---------------------------------------------------------------------------------------|
| bk_obj_id | string | 是  | 条件对象资源类型, host类型的动态分组支持的info.conditon:set,module,host；set类型的动态分组支持的info.condition:set |
| condition | array  | 是  | 查询条件                                                                                  |

#### info.condition.condition

| 字段       | 类型     | 必选 | 描述                                     |
|----------|--------|----|----------------------------------------|
| field    | string | 是  | 对象的字段                                  |
| operator | string | 是  | 操作符, op值为$eq(相等)/$ne(不等)/$in(属于)/$nin(不属于))/$like(模糊匹配) |
| value    | object | 是  | 字段对应的值                                 |

#### info.variable_condition

| 字段        | 类型     | 必选 | 描述                                                                                    |
|-----------|--------|----|---------------------------------------------------------------------------------------|
| bk_obj_id | string | 是  | 条件对象资源类型, host类型的动态分组支持的info.conditon:set,module,host；set类型的动态分组支持的info.condition:set |
| condition | array  | 是  | 查询条件                                                                                  |

#### info.variable_condition.condition

| 字段       | 类型     | 必选 | 描述                                     |
|----------|--------|----|----------------------------------------|
| field    | string | 是  | 对象的字段                                  |
| operator | string | 是  | 操作符, op值为$eq(相等)/$ne(不等)/$in(属于)/$nin(不属于))/$like(模糊匹配) |
| value    | object | 是  | 字段对应的值                                 |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
    "bk_obj_id": "host",
    "name": "my-dynamic-group",
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
    	],
        "variable_condition":[
          {
            "bk_obj_id":"set",
            "condition":[
              {
                "field":"bk_parent_id",
                "operator":"$ne",
                "value":1
              }
            ]
          },
          {
            "bk_obj_id":"module",
            "condition":[
              {
                "field":"bk_parent_id",
                "operator":"$ne",
                "value":1
              }
            ]
          },
          {
            "bk_obj_id":"host",
            "condition":[
              {
                "field":"bk_host_outerip",
                "operator":"$eq",
                "value":"127.0.0.1"
              }
            ]
          }
        ]
    }
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "id": "XXXXXXXX"
    }
}
```

### 返回结果参数

#### response

| 字段         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |
| data       | object | 请求返回的数据                    |

#### data

| 字段 | 类型     | 描述                |
|----|--------|-------------------|
| id | string | 创建成功后返回新的动态分组主键ID |
