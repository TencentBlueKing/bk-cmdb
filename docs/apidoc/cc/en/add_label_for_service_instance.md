### Function Description

Add labels to service instances based on service instance ID and set labels. (Permission: Service instance editing permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field        | Type   | Required | Description                                            |
| ------------ | ------ | -------- | ------------------------------------------------------ |
| instance_ids | array  | Yes      | Service instance IDs, supports up to 100 IDs at a time |
| labels       | object | Yes      | Labels to be added                                     |
| bk_biz_id    | int    | Yes      | Business ID                                            |

#### labels Field Description

- key Validation Rule: `^[a-zA-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`
- value Validation Rule: `^[a-z0-9A-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`

### Request Parameters Example

```python
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "instance_ids": [59, 62],
  "labels": {
    "key1": "value1",
    "key2": "value2"
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
  "data": null

}
```

### Response Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error |
| message    | string | Error message returned for a failed request                  |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Data returned by the request                                 |