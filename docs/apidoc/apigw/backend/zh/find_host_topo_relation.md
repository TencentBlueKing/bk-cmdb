### 描述

获取主机与拓扑的关系(权限：业务访问权限)

### 输入参数

| 参数名称          | 参数类型   | 必选 | 描述            |
|---------------|--------|----|---------------|
| bk_biz_id     | int    | 是  | 业务ID          |
| bk_set_ids    | array  | 否  | 集群ID列表，最多200条 |
| bk_module_ids | array  | 否  | 模块ID列表，最多500条 |
| bk_host_ids   | array  | 否  | 主机ID列表，最多500条 |
| page          | object | 是  | 分页信息          |

#### page 字段说明

| 参数名称  | 参数类型 | 必选 | 描述               |
|-------|------|----|------------------|
| start | int  | 否  | 获取数据偏移位置         |
| limit | int  | 是  | 过去数据条数限制，建议 为200 |

### 调用示例

```json
{
    "page":{
        "start":0,
        "limit":10
    },
    "bk_biz_id":2,
    "bk_set_ids": [1, 2],
    "bk_module_ids": [23, 24],
    "bk_host_ids": [25, 26]
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "data": {
        "count": 2,
        "data": [
            {
                "bk_biz_id": 2,
                "bk_host_id": 2,
                "bk_module_id": 2,
                "bk_set_id": 2,
                "bk_supplier_account": "0"
            },
            {
                "bk_biz_id": 1,
                "bk_host_id": 1,
                "bk_module_id": 1,
                "bk_set_id": 1,
                "bk_supplier_account": "0"
            }
        ],
        "page": {
            "limit": 10,
            "start": 0
        }
    },
    "message": "success",
    "permission": null,
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

#### data 字段说明：

| 参数名称  | 参数类型         | 描述                    |
|-------|--------------|-----------------------|
| count | int          | 记录条数                  |
| data  | object array | 业务下主机与集群，模块，集群的数据详情列表 |
| page  | object       | 页                     |

#### data.data 字段说明：

| 参数名称                | 参数类型   | 描述    |
|---------------------|--------|-------|
| bk_biz_id           | int    | 业务ID  |
| bk_set_id           | int    | 集群ID  |
| bk_module_id        | int    | 模块ID  |
| bk_host_id          | int    | 主机ID  |
| bk_supplier_account | string | 开发商账号 |

#### data.page 字段说明:

| 参数名称  | 参数类型 | 描述       |
|-------|------|----------|
| start | int  | 数据偏移位置   |
| limit | int  | 过去数据条数限制 |
