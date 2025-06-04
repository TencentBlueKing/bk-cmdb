### 描述

新增主机到业务空闲机

- 此接口保证主机要么同时添加成功，要么同时失败(v3.10.25+，权限：主机池主机分配到业务权限)

### 输入参数

| 参数名称         | 参数类型  | 必选 | 描述                |
|--------------|-------|----|-------------------|
| bk_host_list | array | 是  | 主机信息(数组长度一次限制200) |
| bk_biz_id    | int   | 是  | 业务ID              |

#### bk_host_list(主机相关的字段)

| 参数名称               | 参数类型   | 必选 | 描述                                    |
|--------------------|--------|----|---------------------------------------|
| bk_host_innerip    | string | 否  | 主机内网ipv4, 与bk_host_innerip_v6两者其中一个必传 |
| bk_host_innerip_v6 | string | 否  | 主机内网ipv6, 与bk_host_innerip两者其中一个必传    |
| bk_cloud_id        | int    | 是  | 管控区域ID                                |
| bk_addressing      | string | 是  | 寻址方式， "static"、"dynamic"              |
| operator           | string | 否  | 主要维护人                                 |

...

### 调用示例

```json
{
    "bk_biz_id": 3,
    "bk_host_list": [
        {
            "bk_host_innerip": "10.0.0.1",
            "bk_cloud_id": 0,
            "bk_addressing": "dynamic",
            "operator": "admin"
        },
        {
            "bk_host_innerip": "10.0.0.2",
            "bk_cloud_id": 0,
            "bk_addressing": "dynamic",
            "operator": "admin"
        }
    ]
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
        "bk_host_ids": [
            1,
            2
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
| data       | object | 请求返回的数据                    |
| permission | object | 权限信息                       |

#### data

| 参数名称        | 参数类型  | 描述        |
|-------------|-------|-----------|
| bk_host_ids | array | 主机的hostID |
