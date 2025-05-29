### Description

Delete process templates based on a list of process template IDs (Permission: Service template editing permission)

### Parameters

| Name              | Type  | Required | Description                                               |
|-------------------|-------|----------|-----------------------------------------------------------|
| bk_biz_id         | int   | Yes      | Business ID                                               |
| process_templates | array | Yes      | List of process template IDs, with a maximum value of 500 |

### Request Example

```python
{
  "bk_biz_id": 1,
  "process_templates": [50]
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "data": null,
    "message": "success",
    "permission": null,
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
