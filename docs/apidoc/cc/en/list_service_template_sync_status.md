### Function Description

Query the synchronization status of service templates (Version: v3.12.3+, Permission: Business Access).

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type  | Required | Description                                        |
| ------------------- | ----- | -------- | -------------------------------------------------- |
| bk_biz_id           | int   | Yes      | Business ID                                        |
| service_template_id | int   | Yes      | Service template ID                                |
| bk_module_ids       | array | Yes      | List of module IDs to query synchronization status |

### Request Parameter Example

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

### Response Example

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

### Response Result Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | array  | Data returned by the request                                 |

#### data

| Field       | Type   | Description                                  |
| ----------- | ------ | -------------------------------------------- |
| bk_inst_id  | int    | Instance ID, in this case, the module ID     |
| status      | string | Synchronization status                       |
| creator     | string | Creator of the synchronization task          |
| create_time | string | Creation time of the synchronization task    |
| last_time   | string | Last update time of the synchronization task |

**Explanation of Synchronization Status**: There are 6 statuses for instances, including need_sync, new, waiting, executing, finished, and failure. Among them:

- **need_sync**: Awaiting synchronization
- **new/waiting/executing**: Synchronization in progress
- **finished**: Synchronization completed
- **failure**: Synchronization failed