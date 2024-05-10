### 描述

统计每个业务下主机CPU数量 (成本管理专用接口，v3.8.17+/v3.10.18+，权限：全局设置权限)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述   |
|-----------|--------|----|------|
| bk_biz_id | int    | 否  | 业务ID |
| page      | object | 否  | 分页信息 |

**注：bk_biz_id和page参数必须且只能传其中一个**

#### page 字段说明

| 参数名称  | 参数类型 | 必选 | 描述           |
|-------|------|----|--------------|
| start | int  | 是  | 记录开始位置       |
| limit | int  | 是  | 每页限制条数，最多10条 |

### 调用示例

```json
{
    "page": {
        "start": 10,
        "limit": 10
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
    "data": [
        {
            "bk_biz_id": 5,
            "host_count": 100,
            "cpu_count": 192,
            "no_cpu_host_count": 5
        },
        {
            "bk_biz_id": 7,
            "host_count": 40,
            "cpu_count": 58,
            "no_cpu_host_count": 11
        }
    ]
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

| 参数名称              | 参数类型 | 描述             |
|-------------------|------|----------------|
| bk_biz_id         | int  | 业务ID           |
| host_count        | int  | 主机数量           |
| cpu_count         | int  | CPU数量          |
| no_cpu_host_count | int  | 没有CPU数量字段的主机数量 |
