### 描述

模型实例关系数量查询 (v3.10.1+)

### 输入参数

| 参数名称       | 参数类型   | 必选 | 描述                                                                                 |
|------------|--------|----|------------------------------------------------------------------------------------|
| bk_biz_id  | int    | 否  | 业务ID, 针对主线模型查询时需要提供                                                                |
| bk_obj_id  | string | 是  | 模型ID                                                                               |
| conditions | object | 否  | 组合查询条件,  组合支持AND和OR两种方式，可以嵌套，最多嵌套3层, 每层OR条件最大支持20个, 不指定该参数表示匹配全部(即conditions为null) |

#### conditions

| 参数名称      | 参数类型   | 必选 | 描述          |
|-----------|--------|----|-------------|
| condition | string | 是  | 规则操作符       |
| rules     | array  | 是  | 所选业务的范围条件规则 |

#### conditions.rules

| 参数名称     | 参数类型   | 必选 | 描述                                                                                                  |
|----------|--------|----|-----------------------------------------------------------------------------------------------------|
| field    | string | 是  | 条件字段, 可选值 id, bk_inst_id, bk_obj_id, bk_asst_inst_id, bk_asst_obj_id, bk_obj_asst_id, bk_asst_id    |
| operator | string | 是  | 操作符, 可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between等 |
| value    | -      | 否  | 条件字段期望的值, 不同的operator对应不同的value格式, 数组类型值最大支持500个元素                                                  |

组装规则详细可参考: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

### 调用示例

```json
{
    "bk_obj_id":"bk_switch",
    "conditions":{
        "condition": "AND",
        "rules": [
            {
                "field": "bk_obj_asst_id",
                "operator": "equal",
                "value": "bk_switch_connect_host"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                         "field": "bk_inst_id",
                         "operator": "in",
                         "value": [2,4,6]
                    },
                    {
                        "field": "bk_asst_id",
                        "operator": "equal",
                        "value": 3
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
    "message": "success",
    "permission": null,
    "data": {
        "count": 1
    }
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

| 参数名称  | 参数类型 | 描述            |
|-------|------|---------------|
| count | int  | 返回满足条件的实例数据数量 |
