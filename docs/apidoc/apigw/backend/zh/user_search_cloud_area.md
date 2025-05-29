### 描述

查询管控区域(权限：管控区域查看权限)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述   |
|-----------|--------|----|------|
| condition | object | 否  | 查询条件 |
| page      | object | 是  | 分页信息 |

#### condition

| 参数名称          | 参数类型   | 必选 | 描述     |
|---------------|--------|----|--------|
| bk_cloud_id   | int    | 否  | 管控区域ID |
| bk_cloud_name | string | 否  | 管控区域名称 |

#### page 字段说明

| 参数名称  | 参数类型 | 必选 | 描述               |
|-------|------|----|------------------|
| start | int  | 否  | 获取数据偏移位置         |
| limit | int  | 是  | 过去数据条数限制，建议 为200 |

### 调用示例

```json
{

    "condition": {
        "bk_cloud_id": 12,
        "bk_cloud_name": "aws"
    },
    "page":{
        "start":0,
        "limit":10
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
    "count": 10,
    "info": [
         {
            "bk_cloud_id": 0,
            "bk_cloud_name": "aws",
            "bk_supplier_account": "0",
            "create_time": "2019-05-20T14:59:48.354+08:00",
            "last_time": "2019-05-20T14:59:48.354+08:00"
        },
        .....
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

| 参数名称  | 参数类型  | 描述           |
|-------|-------|--------------|
| count | int   | 记录条数         |
| info  | array | 查询到的管控区域列表信息 |

#### data.info 字段说明：

| 参数名称          | 参数类型   | 描述     |
|---------------|--------|--------|
| bk_cloud_id   | int    | 管控区域ID |
| bk_cloud_name | string | 管控区域名字 |
| create_time   | string | 创建时间   |
| last_time     | string | 最后修改时间 |
