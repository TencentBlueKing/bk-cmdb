### Function Description

Delete service template based on service template ID (Permission: Service template deletion permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type | Required | Description         |
| ------------------- | ---- | -------- | ------------------- |
| service_template_id | int  | Yes      | Service template ID |
| bk_biz_id           | int  | Yes      | Business ID         |

### Request Parameter Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "service_template_id": 1
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
    "request_id": "b78feeebd55b4265b463200ab966f506"
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