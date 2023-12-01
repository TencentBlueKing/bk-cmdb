### 功能描述

查询服务模板的同步状态(版本：v3.12.3+，权限：业务访问)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                  | 类型    | 必选 | 描述             |
|---------------------|-------|----|----------------|
| bk_biz_id           | int   | 是  | 业务ID           |
| service_template_id | int   | 是  | 服务模板ID         |
| bk_module_ids       | array | 是  | 要查询同步状态的模块ID列表 |

### 请求参数示例

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 3,
  "service_template_id": 1,
  "bk_module_ids": [
    28,
    29,
    30
  ]
}
```

### 返回结果示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": [
    {
      "bk_inst_id": 30,
      "status": "finished",
      "creator": "admin",
      "create_time": "2023-10-07T12:43:22.795Z",
      "last_time": "2023-11-10T03:37:31.009Z"
    },
    {
      "bk_inst_id": 29,
      "status": "finished",
      "creator": "admin",
      "create_time": "2023-10-07T07:22:43.167Z",
      "last_time": "2023-11-10T03:37:31.005Z"
    },
    {
      "bk_inst_id": 28,
      "status": "new",
      "creator": "admin",
      "create_time": "2023-11-30T09:52:13.706Z",
      "last_time": "2023-11-30T09:52:13.706Z"
    }
  ]
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
| data       | array  | 请求返回的数据                    |

#### data

| 字段          | 类型     | 描述            |
|-------------|--------|---------------|
| bk_inst_id  | int    | 实例id，此处为模块ID  |
| status      | string | 同步状态          |
| creator     | string | 同步任务的创建者      |
| create_time | string | 同步任务的创建时间     |
| last_time   | string | 同步任务的最后一次更新时间 |

**同步状态说明**： 实例状态共有need_sync、new、waiting、executing、finished、failure 6种状态，其中：

- **need_sync** 为待同步
- **new(新建)/waiting(等待中)/executing(执行中)** 为同步中
- **finished** 为已同步
- **failure** 为同步失败
