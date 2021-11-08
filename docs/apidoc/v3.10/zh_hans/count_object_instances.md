### 功能描述

通用模型实例数量查询 (v3.10.1+)

### 请求参数

{{ common_args_desc }}

#### 接口参数

|    字段    |  类型  | 必选 | 描述                                                                                                            |
|------------|--------|------|-----------------------------------------------------------------------------------------------------------------|
| bk_obj_id  | string |  是  | 模型ID                                                                                                          |
| conditions | object |  否  | 组合查询条件,  组合支持AND和OR两种方式，可以嵌套，最多嵌套3层, 每层OR条件最大支持20个, 不指定该参数表示匹配全部(即conditions为null) |

#### conditions

|   字段   |  类型  | 必选 |  描述                                                                                                     |
|----------|--------|------|-----------------------------------------------------------------------------------------------------------|
| field    | string |  是  | 条件字段                                                                                                  |
| operator | string |  是  | 操作符, 可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between等|
| value    |   -    |  否  | 条件字段期望的值, 不同的operator对应不同的value格式, 数组类型值最大支持500个元素                          |

组装规则详细可参考: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

### 请求参数示例

```json
{
    "bk_app_code":"code",
    "bk_app_secret":"secret",
    "bk_token":"xxxx",
    "bk_obj_id":"bk_switch",
    "conditions":{
        "condition": "AND",
        "rules": [
            {
                "field": "bk_inst_name",
                "operator": "equal",
                "value": "switch"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                         "field": "bk_inst_id",
                         "operator": "not_in",
                         "value": [2,4,6]
                    },
                    {
                        "field": "bk_inst_id",
                        "operator": "equal",
                        "value": 3
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
        "count": 1
    }
}
```

### 返回结果参数

#### data

| 字段  |   类型  | 描述                       |
|-------|---------|----------------------------|
| count | integer | 返回满足条件的实例数据数量 |
