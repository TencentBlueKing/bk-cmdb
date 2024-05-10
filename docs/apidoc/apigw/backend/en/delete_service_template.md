### Description

Delete service template based on service template ID (Permission: Service template deletion permission)

### Parameters

| Name                | Type | Required | Description         |
|---------------------|------|----------|---------------------|
| service_template_id | int  | Yes      | Service template ID |
| bk_biz_id           | int  | Yes      | Business ID         |

### Request Example

```python
{
  "bk_biz_id": 1,
  "service_template_id": 1
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
