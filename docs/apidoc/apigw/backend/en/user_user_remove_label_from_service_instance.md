### Description

Remove tags from specified service instances under the specified business based on the business ID, service instance ID,
and the tags to be removed. (Permission: Service instance deletion permission)

### Parameters

| Name         | Type  | Required | Description                                         |
|--------------|-------|----------|-----------------------------------------------------|
| bk_biz_id    | int   | Yes      | Business ID                                         |
| instance_ids | array | Yes      | List of service instance IDs, with a maximum of 500 |
| keys         | array | Yes      | List of tag keys to be removed                      |

### Request Example

```python
{
  "bk_biz_id": 1,
  "instance_ids": [60, 62],
  "keys": ["value1", "value3"]
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

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |
