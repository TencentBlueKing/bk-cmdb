### Description

Delete service categories based on service category IDs (Permission: Service category deletion permission)

### Parameters

| Name      | Type | Required | Description         |
|-----------|------|----------|---------------------|
| id        | int  | Yes      | Service category ID |
| bk_biz_id | int  | Yes      | Business ID         |

### Request Example

```python
{
  "bk_biz_id": 1,
  "id": 6
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
