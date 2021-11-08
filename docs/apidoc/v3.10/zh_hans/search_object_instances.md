### 功能描述

通用模型实例查询 (v3.10.1+)

### 请求参数

{{ common_args_desc }}

#### 接口参数

|    字段    |  类型  | 必选 | 描述                                                                                                            |
|------------|--------|------|-----------------------------------------------------------------------------------------------------------------|
| bk_obj_id  | string |  是  | 模型ID                                                                                                          |
| conditions | object |  否  | 组合查询条件,  组合支持AND和OR两种方式，可以嵌套，最多嵌套3层, 每层OR条件最大支持20个, 不指定该参数表示匹配全部(即conditions为null) |
| fields     | array  |  否  | 指定需要返回的字段, 不具备的字段将被忽略, 不指定则返回全部字段（返回全部字段会对性能产生影响，建议按需返回）    |
| page       | object |  是  | 分页设置                                                                                                        |

#### conditions

|   字段   |  类型  | 必选 |  描述                                                                                                     |
|----------|--------|------|-----------------------------------------------------------------------------------------------------------|
| field    | string |  是  | 条件字段                                                                                                  |
| operator | string |  是  | 操作符, 可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between等|
| value    |   -    |  否  | 条件字段期望的值, 不同的operator对应不同的value格式, 数组类型值最大支持500个元素                          |

组装规则详细可参考: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### page

|  字段 |  类型  | 必选 |  描述                                                            |
|-------|--------|------|------------------------------------------------------------------|
| start | int    |  是  | 记录开始位置                                                     |
| limit | int    |  是  | 每页限制条数, 最大500                                            |
| sort  | string |  否  | 检索排序，遵循MongoDB语义格式{KEY}:{ORDER}，默认按照创建时间排序 |

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
    },
    "fields":[
        "bk_inst_id",
        "bk_inst_name"
    ],
    "page":{
        "start":0,
        "limit":500
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
        "info": [
            {
                "bk_inst_id": 1,
                "bk_inst_name": "switch-instance"
            }
        ]
    }
}
```

### 返回结果参数

#### data

| 字段 | 类型  | 描述                                |
|------|-------|-------------------------------------|
| info | array | map数组格式, 返回满足条件的实例数据 |
