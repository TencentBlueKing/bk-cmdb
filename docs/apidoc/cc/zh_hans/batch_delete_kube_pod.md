### 功能描述

批量删除Pod (版本：v3.10.23+，权限：容器Pod删除)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段   | 类型           | 必选  | 描述                             |
|------|--------------|-----|--------------------------------|
| data | object array | 是   | 要删除的pod信息数组，data里所有pod总和最多200条 |

#### data 字段说明

| 字段        | 类型        | 必选  | 描述                                    |
|-----------|-----------|-----|---------------------------------------|
| bk_biz_id | int       | 是   | 业务ID                                  |
| ids       | int array | 是   | 要删除的pod的cc ID数组，data里所有pod的ID总和最多200条 |

### 请求参数示例

```json
{
    "bk_app_code": "code",
    "bk_app_secret": "secret",
    "bk_username": "xxx",
    "bk_token": "xxxx",
    "data": [
        {
            "bk_biz_id": 123,
            "ids": [
                5,
                6
            ]
        }
    ]
}
```

### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### 返回结果参数

#### response

| 名称         | 类型     | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| request_id | string | 请求链id                      |
