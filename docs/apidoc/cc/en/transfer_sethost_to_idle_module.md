### Function Description

Transfer hosts to the business's idle module under the specified business cluster and module (Permission: Service Instance Edit Permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field        | Type | Required | Description |
| ------------ | ---- | -------- | ----------- |
| bk_biz_id    | int  | Yes      | Business ID |
| bk_set_id    | int  | Yes      | Cluster ID  |
| bk_module_id | int  | Yes      | Module ID   |

### Request Parameter Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id":10,
    "bk_module_id":58,
    "bk_set_id":1
}
```

### Return Result Example

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