### 功能描述

创建动态分组 (V3.9.6)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id |  int     | 是     | 业务ID |
| bk_obj_id |  string  | 是     | 动态分组的目标资源对象类型,目前可以为host,set |
| info      |  object  | 是     | 通用查询条件 |
| name      |  string  | 是     | 动态分组名称 |

#### info.condition

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_obj_id |  string   | 是     | 条件对象资源类型, host类型的动态分组支持的info.conditon:set,module,host；set类型的动态分组支持的info.condition:set |
| condition |  array    | 是     | 查询条件 |

#### info.condition.condition

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| field     |  string    | 是     | 对象的字段 |
| operator  |  string    | 是     | 操作符, op值为eq(相等)/ne(不等)/in(属于)/nin(不属于) |
| value     |  object    | 是     | 字段对应的值 |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
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
    	]
    }
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "data": {
        "id": "XXXXXXXX"
    }
}
```

### 返回结果参数

#### data

| 字段    | 类型  | 描述      |
|--------|-------|-----------|
| id     | string | 创建成功后返回新的动态分组主键ID |
