### 功能描述

根据服务实例id和设置的标签为服务实例添加标签

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
|instance_ids|array|是|服务实例ID,最大值为100|
|labels|object|是|添加的Label|
| bk_biz_id  | int     |是     | 业务ID       |

#### labels 字段说明
- key 校验规则: `^[a-zA-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`
- value 校验规则: `^[a-z0-9A-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`

### 请求参数示例

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "instance_ids": [59, 62],
  "labels": {
    "key1": "value1",
    "key2": "value2"
  }
}
```

### 返回结果示例

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": null

}
```

### 返回结果说明

#### response

| 名称  | 类型  | 描述 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                           |


