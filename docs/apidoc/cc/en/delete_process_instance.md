### Function Description

Delete process instances based on a list of process instance IDs (Permission: Service instance editing permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type | Required | Description                                               |
| -------------------- | ---- | -------- | --------------------------------------------------------- |
| process_instance_ids | int  | Yes      | List of process instance IDs, with a maximum value of 500 |
| bk_biz_id            | int  | Yes      | Business ID                                               |

### Request Parameter Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "process_instance_ids": [54]
}
```

### Response Example

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