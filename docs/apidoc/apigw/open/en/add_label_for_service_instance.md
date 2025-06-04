### Description

Add labels to service instances based on service instance ID and set labels. (Permission: Service instance editing
permission)

### Parameters

| Name         | Type   | Required | Description                                            |
|--------------|--------|----------|--------------------------------------------------------|
| instance_ids | array  | Yes      | Service instance IDs, supports up to 100 IDs at a time |
| labels       | object | Yes      | Labels to be added                                     |
| bk_biz_id    | int    | Yes      | Business ID                                            |

#### labels Field Description

- key Validation Rule: `^[a-zA-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`
- value Validation Rule: `^[a-z0-9A-Z]([a-z0-9A-Z\-_.]*[a-z0-9A-Z])?$`

### Request Example

```python
{
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
  "data": null

}
```

### Response Parameters

| Name       | Type   | Description                                                       |
|------------|--------|-------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error     |
| message    | string | Error message returned for a failed request                       |
| permission | object | Permission information                                            |
| data       | object | Data returned by the request                                      |
