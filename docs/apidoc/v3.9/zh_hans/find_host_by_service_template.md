### 功能描述

获取服务模板下的主机 (v3.8.6)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| bk_biz_id  | int  | 是     | 业务ID |
| bk_service_template_ids  | int array  | 是     | 服务模板ID列表，最多可填500个 |
| bk_module_ids  | int array  | 否     | 模块ID列表, 最多可填500个 |
| fields  |  string array   | 是     | 主机属性列表，控制返回结果的模块信息里有哪些字段 |
| page       |  object    | 是     | 分页信息 |

#### page 字段说明

| 字段  | 类型   | 必选 | 描述                  |
| ----- | ------ | ---- | --------------------- |
| start | int    | 是   | 记录开始位置          |
| limit | int    | 是   | 每页限制条数,最大500 |

### 请求参数示例

```json
{
    "bk_biz_id": 5,
    "bk_service_template_ids": [
        48,
        49
    ],
    "bk_module_ids": [
        65,
        68
    ],
    "fields": [
        "bk_host_id",
        "bk_cloud_id"
    ],
    "page": {
        "start": 0,
        "limit": 10
    }
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": {
        "count": 6,
        "info": [
            {
                "bk_cloud_id": 0,
                "bk_host_id": 1
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 2
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 3
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 4
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 7
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 8
            }
        ]
    }
}
```
