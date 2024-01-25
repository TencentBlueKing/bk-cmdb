### Function Description

Synchronize service template information to the corresponding service instances (Version: v3.12.3+, Permission: Create, edit, delete permissions for service instances)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type  | Required | Description                           |
| ------------------- | ----- | -------- | ------------------------------------- |
| bk_biz_id           | int   | Yes      | Business ID                           |
| service_template_id | int   | Yes      | Service template ID                   |
| bk_module_ids       | array | Yes      | List of module IDs to be synchronized |

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
    28
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
  "data": null
}
```

### Return Result Parameter Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | object | Request returned data                                        |