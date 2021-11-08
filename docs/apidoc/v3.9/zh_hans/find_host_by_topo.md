### 功能描述

查询拓扑节点下的主机 (v3.8.13)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段        |  类型  | 必选   |  描述      |
|------------|--------|--------|------------|
| bk_biz_id  | int    | 是     | 业务ID |
| bk_obj_id  | string | 是     | 拓扑节点模型ID，不可为biz |
| bk_inst_id | int    | 是     | 拓扑节点实例ID |
| fields     | array  | 是     | 主机属性列表，控制返回结果的主机里有哪些字段，能够加速接口请求和减少网络流量传输   |
| page       | object | 是     | 分页信息 |

#### page 字段说明

| 字段  | 类型   | 必选 | 描述                  |
| ----- | ------ | ---- | --------------------- |
| start | int    | 是   | 记录开始位置          |
| limit | int    | 是   | 每页限制条数,最大500 |

### 请求参数示例

```json
{
    "bk_biz_id": 5,
    "bk_obj_id": "xxx",
    "bk_inst_id": 10,
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
        "count": 2,
        "info": [
            {
                "bk_cloud_id": 0,
                "bk_host_id": 1
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 2
            }
        ]
    }
}
```

### 返回结果参数说明

#### data

| 字段      | 类型      | 描述      |
|-----------|-----------|-----------|
| count     | int       | 记录条数 |
| info      | array     | 主机实际数据 |
