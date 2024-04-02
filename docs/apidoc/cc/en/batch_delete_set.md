### Function Description

Batch delete set by set ID under a specified business ID (Permission: Business topology deletion permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description |
| --------- | ------ | -------- | ----------- |
| bk_biz_id | int    | Yes      | Business ID |
| delete    | object | Yes      | Deletion    |

#### delete

| Field    | Type      | Required | Description                        |
| -------- | --------- | -------- | ---------------------------------- |
| inst_ids | int array | Yes      | Array of Cluster IDs to be deleted |

### Request Parameter Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 0,
    "delete": {
        "inst_ids": [123]
    }
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
    "data": "success"
}
```

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error  |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Data returned in the request                                 |