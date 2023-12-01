Functional description

Example Query the synchronization status of a service template(v3.12.3+, permission: biz access)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type  | Required | Description                                                       |
|---------------------|-------|----------|-------------------------------------------------------------------|
| bk_biz_id           | int   | yes      | Business ID                                                       |
| service_template_id | int   | yes      | Service template ID                                               |
| bk_module_ids       | array | yes      | List of module ids whose synchronization status you want to query |

### Request Parameters Example

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

### Return Result Example

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

### Return Result Parameters Description

#### response

| Name       | Type   | Description                                                                             |
|------------|--------|-----------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. True: request succeeded;false request failed |
| code       | int    | Wrong code. 0 indicates success,>0 indicates failure error                              |
| message    | string | Error message returned by request failure                                               |
| permission | object | Permission information                                                                  |
| request_id | string | Request chain id                                                                        |
| data       | array  | Data returned by request                                                                |

#### data

| Name        | Type   | Description                                      |
|-------------|--------|--------------------------------------------------|
| bk_inst_id  | int    | Instance id, where is the module ID              |
| status      | string | Synchronous status                               |
| creator     | string | The creator of a synchronization task            |
| create_time | string | The create time of the synchronization task      |
| last_time   | string | The last update time of the synchronization task |

**Synchronization status declaration**： An instance has six states: need_sync, new, waiting, executing, finished, and
failure，among：

- **need_sync** Indicates to be synchronized
- **new/waiting/executing** Indicates synchronization
- **finished** Indicates synchronization complete
- **failure** Indicates synchronization failure
