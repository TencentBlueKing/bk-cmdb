### 描述

克隆主机属性(权限：业务主机编辑权限)

### 输入参数

| 参数名称        | 参数类型   | 必选 | 描述       |
|-------------|--------|----|----------|
| bk_org_ip   | string | 是  | 源主机内网ip  |
| bk_dst_ip   | string | 是  | 目标主机内网ip |
| bk_org_id   | int    | 是  | 源主机ID    |
| bk_dst_id   | int    | 是  | 目标主机ID   |
| bk_biz_id   | int    | 是  | 业务ID     |
| bk_cloud_id | int    | 是  | 管控区域ID   |

注： 使用主机内网IP进行克隆与使用主机身份ID进行克隆，这两种方式只能使用期中的一种，不能混用。

### 调用示例

```json
{
    "bk_biz_id":2,
    "bk_org_ip":"127.0.0.1",
    "bk_dst_ip":"127.0.0.2",
    "bk_cloud_id":0
}
```

或

```json
{
    "bk_biz_id":2,
    "bk_org_id": 10,
    "bk_dst_id": 11,
    "bk_cloud_id":0
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": null
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
