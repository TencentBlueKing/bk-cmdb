### Description

Batch delete instances of referenced models (v3.10.30+, Permission: Edit permission of source model instances)

### Parameters

| Name           | Type        | Required | Description                                                       |
|----------------|-------------|----------|-------------------------------------------------------------------|
| bk_obj_id      | string      | Yes      | Source model ID                                                   |
| bk_property_id | string      | Yes      | ID of the property in the source model that references this model |
| ids            | int64 array | Yes      | Array of instance IDs to be deleted, with a maximum of 500        |

### Request Example

```json
{
  "bk_obj_id": "host",
  "bk_property_id": "disk",
  "ids": [
    1,
    2
  ]
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": null,
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Data returned in the request                                                |
