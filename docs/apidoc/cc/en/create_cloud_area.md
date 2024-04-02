### Function Description

Create a control area based on the control area name (Permission: Control Area Creation Permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field         | Type   | Required | Description       |
| ------------- | ------ | -------- | ----------------- |
| bk_cloud_name | string | Yes      | Control area name |

### Request Parameter Example

```python
{
    
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_cloud_name": "test1"
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
    "data": {
        "created": {
            "origin_index": 0,
            "id": 6
        }
    }
}
```

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error  |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Request return data                                          |

#### data

| Field   | Type   | Description                              |
| ------- | ------ | ---------------------------------------- |
| created | object | Created successfully, return information |

#### data.created

| Field         | Type | Description                                       |
| ------------ | ---- | ------------------------------------------------- |
| origin_index | int  | Corresponding to the order of the request results |
| id           | int  | Control area id, bk_cloud_id                      |