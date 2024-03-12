### Function Description

Delete service categories based on service category IDs (Permission: Service category deletion permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type | Required | Description         |
| --------- | ---- | -------- | ------------------- |
| id        | int  | Yes      | Service category ID |
| bk_biz_id | int  | Yes      | Business ID         |

### Request Parameter Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "id": 6
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
    "request_id": "28807929b7af4fcd9b834fd200ceb2ad"
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