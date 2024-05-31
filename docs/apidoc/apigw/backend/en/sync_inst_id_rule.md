### Description

When the id rule field of the model instance is empty, the id rule field value is asynchronously refreshed to the model instance field through this interface (Version: v3.14.1, Permission: edit permission of the corresponding instance)

### Parameters

| Name      | Type   | Required | Description |
|-----------|--------|----------|-------------|
| bk_obj_id | string | Yes      | model id    |

### Request Example

```json
{
  "bk_obj_id": "host"
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "",
  "permission": null,
  "data": {
    "task_id": "111"
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| permission | object | Permission information                                              |
| data       | object | Data returned by the request                                        |

#### data

| Name       | Type   | Description |
|------------|--------|------------|
| task_id     | string | task id    |