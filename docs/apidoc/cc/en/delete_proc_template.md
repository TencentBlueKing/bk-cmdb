### Function Description

Delete process templates based on a list of process template IDs (Permission: Service template editing permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field             | Type  | Required | Description                                               |
| ----------------- | ----- | -------- | --------------------------------------------------------- |
| bk_biz_id         | int   | Yes      | Business ID                                               |
| process_templates | array | Yes      | List of process template IDs, with a maximum value of 500 |

### Request Parameter Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "process_templates": [50]
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "data": null,
    "message": "success",
    "permission": null,
    "request_id": "069cbd4eed2846a0b4c995f3d040e2a5"
}
```

### Return Result Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Data returned by the request                                 |