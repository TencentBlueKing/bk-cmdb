### Description

Batch delete service instances based on service instance IDs (Permission: Service instance deletion permission)

### Parameters

| Name                 | Type  | Required | Description                                    |
|----------------------|-------|----------|------------------------------------------------|
| service_instance_ids | array | Yes      | Service instance ID list, maximum value is 500 |
| bk_biz_id            | int   | Yes      | Business ID                                    |

### Request Example

```python
{
  "bk_biz_id": 1,
  "service_instance_ids": [48]
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
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |
