### Function Description

Query host locks based on host ID list (Version: v3.8.6, Permission: Business host edit permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field   | Type  | Required | Description  |
| ------- | ----- | -------- | ------------ |
| id_list | array | Yes      | Host ID list |

### Request Parameter Example

```python
{
   "bk_app_code": "esb_test",
   "bk_app_secret": "xxx",
   "bk_username": "xxx",
   "bk_token": "xxx",
   "id_list":[1, 2]
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
    "data": {
        1: true,
        2: false
    }
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

#### data

| Field | Type   | Description                                                  |
| ----- | ------ | ------------------------------------------------------------ |
| data  | object | Data returned by the request, where the key is ID, and the value is whether it is locked |