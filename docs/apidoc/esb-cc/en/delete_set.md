### Function Description

Delete Cluster (Permission: Business Topology Deletion Permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type   | Required | Description       |
| ------------------- | ------ | -------- | ----------------- |
| bk_supplier_account | string | No       | Developer account |
| bk_biz_id           | int    | Yes      | Business ID       |
| bk_set_id           | int    | Yes      | Cluster ID        |

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "123456789",
    "bk_biz_id": 1,
    "bk_set_id": 10
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

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure        |
| message    | string | Error message returned in case of request failure            |
| data       | object | Request returned data                                        |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |