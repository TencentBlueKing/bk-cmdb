### 批量更新业务属性

- API: PUT /api/v3/updatemany/biz/property
- API 名称：batch_update_biz
- 功能说明：
	- 中文：批量更新业务属性
	- English：batch update business properties

- input body

example 1:
```json
{
    "properties":{
      "bk_biz_developer":"developer",
      "bk_biz_maintainer": "maintainer",
      "bk_biz_name":"biz_test",
      "bk_biz_productor": "productor",
      "bk_biz_tester":"tester",
      "operator": "operator"
    },
    "condition": {
        "bk_biz_id": {"$in": [3,4]}
    }
}
```

example 2:
```json
{
    "properties":{
      "bk_biz_developer":"developer",
      "bk_biz_maintainer": "maintainer",
      "bk_biz_name":"biz_test",
      "bk_biz_productor": "productor",
      "bk_biz_tester":"tester",
      "operator": "operator"
    },
    "condition": {
        "bk_biz_id": {"$nin": [2]},	// exclude bk_biz_id 2
        "life_cycle": {"$in": ["1"]}
    }
}
```

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
| properties  | object | 是 | 无    | 业务被更新的属性和值 | business property keys and values to be updated |
| condition | object | 是 |  无 | 被更新业务的过滤条件 |  business property update condition |

properties 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
| bk_biz_developer | string | 否 | 无    | 开发人员    | developer |
| bk_biz_maintainer     | string | 否 | 无   | 运维人员 | maintainer |
| bk_biz_name   | string | 否 | 无   | 业务名      | business name |
| bk_biz_productor  | string | 否 | 无   | 产品人员 | productor |
| bk_biz_tester | string | 否 | 无 | 测试人员 | tester |
| operator | string | 否 | 无 | 操作人员 | operator |

- output

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": null
}
```

注意：当更新的业务列表中，只要有一个没有权限，则会全部更新失败。
