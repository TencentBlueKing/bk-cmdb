### Function Description

Delete business set (Version: v3.10.12+, Permission: business set deletion permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field          | Type  | Required | Description              |
| -------------- | ----- | -------- | ------------------------ |
| bk_biz_set_ids | array | Yes      | List of business set IDs |

### Request Parameter Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_set_ids": [
        10,
        12
    ]
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {},
    "request_id": "dsda1122adasadadada2222"
}
```

### Response Parameter Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error  |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| data       | object | Data returned in the request                                 |
| request_id | string | Request chain ID                                             |