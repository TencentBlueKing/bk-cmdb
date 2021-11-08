### 功能描述

根据业务id和集群实例id列表，以及想要获取的属性列表，批量获取指定业务下集群的属性详情 (v3.8.6)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id  | int  | 是     | 业务ID |
| bk_ids  | int array  | 是     | 集群实例ID列表, 即bk_set_id列表，最多可填500个 |
| fields  |  string array   | 是     | 集群属性列表，控制返回结果的集群信息里有哪些字段 |

### 请求参数示例

```json
{
    "bk_biz_id": 3,
    "bk_ids": [
        11,
        12
    ],
    "fields": [
        "bk_set_id",
        "bk_set_name",
        "create_time"
    ]
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": [
        {
            "bk_set_id": 12,
            "bk_set_name": "ss1",
            "create_time": "2020-05-15T22:15:51.67+08:00",
            "default": 0
        },
        {
            "bk_set_id": 11,
            "bk_set_name": "set1",
            "create_time": "2020-05-12T21:04:36.644+08:00",
            "default": 0
        }
    ]
}
```
