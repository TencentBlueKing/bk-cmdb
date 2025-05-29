### 描述

通用模型实例关系查询 (版本：v3.10.1+，权限：模型实例查询权限)

### 输入参数

| 参数名称       | 参数类型   | 必选 | 描述                                                                                 |
|------------|--------|----|------------------------------------------------------------------------------------|
| bk_biz_id  | int    | 否  | 业务ID, 针对主线模型查询时需要提供                                                                |
| bk_obj_id  | string | 是  | 模型ID                                                                               |
| conditions | object | 否  | 组合查询条件,  组合支持AND和OR两种方式，可以嵌套，最多嵌套3层, 每层OR条件最大支持20个, 不指定该参数表示匹配全部(即conditions为null) |
| fields     | array  | 否  | 指定需要返回的字段, 不具备的字段将被忽略, 不指定则返回全部字段（返回全部字段会对性能产生影响，建议按需返回）                           |
| page       | object | 是  | 分页设置                                                                               |

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

#### page

| 参数名称  | 参数类型   | 必选 | 描述                                         |
|-------|--------|----|--------------------------------------------|
| start | int    | 是  | 记录开始位置                                     |
| limit | int    | 是  | 每页限制条数, 最大500                              |
| sort  | string | 否  | 检索排序，遵循MongoDB语义格式{KEY}:{ORDER}，默认按照创建时间排序 |

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
    },
    "fields":[
        "bk_inst_id",
        "bk_asst_inst_id",
        "bk_asst_obj_id",
        "bk_asst_id",
        "bk_obj_asst_id"
    ],
    "page":{
        "start":0,
        "limit":500
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
    "data": {
        "info": [
            {
                "bk_inst_id": 2,
                "bk_asst_inst_id": 8,
                "bk_asst_obj_id": "host",
                "bk_asst_id": "connect",
                "bk_obj_asst_id": "bk_switch_connect_host"
            }
        ]
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

| 参数名称 | 参数类型  | 描述                   |
|------|-------|----------------------|
| info | array | map数组格式, 返回满足条件的实例数据 |

#### info

| 参数名称            | 参数类型   | 描述            |
|-----------------|--------|---------------|
| bk_inst_id      | int    | 源模型实例id       |
| bk_asst_inst_id | int    | 目标模型实例id      |
| bk_asst_obj_id  | string | 关联关系目标模型id    |
| bk_asst_id      | string | 关联类型id        |
| bk_obj_asst_id  | string | 自动生成的模型关联关系id |
