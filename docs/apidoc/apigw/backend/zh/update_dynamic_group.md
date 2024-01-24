### 描述

更新动态分组 (版本：v3.9.6，权限：动态分组编辑权限)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述                                                  |
|-----------|--------|----|-----------------------------------------------------|
| bk_biz_id | int    | 是  | 业务ID                                                |
| id        | string | 是  | 主键ID                                                |
| bk_obj_id | string | 否  | 动态分组的目标资源对象类型, 目前可以为host,set.更新规则时需同时提供该字段和info两个字段 |
| info      | object | 否  | 通用查询条件                                              |
| name      | string | 否  | 动态分组名称                                              |

#### info.condition

| 参数名称      | 参数类型   | 必选 | 描述                                                                                    |
|-----------|--------|----|---------------------------------------------------------------------------------------|
| bk_obj_id | string | 是  | 条件对象资源类型, host类型的动态分组支持的info.conditon:set,module,host；set类型的动态分组支持的info.condition:set |
| condition | array  | 是  | 查询条件                                                                                  |

#### info.condition.condition

| 参数名称     | 参数类型   | 必选 | 描述                                     |
|----------|--------|----|----------------------------------------|
| field    | string | 是  | 对象的字段                                  |
| operator | string | 是  | 操作符, op值为eq(相等)/ne(不等)/in(属于)/nin(不属于) |
| value    | object | 是  | 字段对应的值                                 |

### 调用示例

```json
{
    "bk_biz_id": 1,
    "id": "XXXXXXXX",
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
    	]
    }
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": {}
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
