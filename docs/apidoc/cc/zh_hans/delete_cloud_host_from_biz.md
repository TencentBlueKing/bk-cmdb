### 功能描述

从业务空闲机集群删除云主机 (云主机管理专用接口, 版本：v3.10.19+，权限：业务主机编辑权限)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段          | 类型        | 必选  | 描述                                    |
|-------------|-----------|-----|---------------------------------------|
| bk_biz_id   | int       | 是   | 业务ID                                  |
| bk_host_ids | array | 是   | 删除的云主机ID数组，数组长度最多为200，一批主机仅可同时成功或同时失败 |

**注意：通过该接口只能删除云主机，填写其它普通主机的id会报主机不存在，云主机的属性中bk_cloud_host_identifier字段为true，其它普通主机的为false，可通过添加云主机相关的接口如（add_cloud_host_to_biz 新增云主机到业务的空闲机模块）添加云主机**

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 123,
    "bk_host_ids": [
        1,
        2
    ]
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "data": null,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807"
}
```

### 返回结果参数说明

#### response

| 名称         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |
