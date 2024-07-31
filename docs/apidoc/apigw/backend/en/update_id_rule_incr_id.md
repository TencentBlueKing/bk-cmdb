### Description

Update id rule self-incremented id (Version: v3.14.1, Permission: id self-increment id update permission)

### Parameters

| Name        | Type   | Required | Description                                                                                                                                                     |
|-------------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| type        | string | Yes      | if it is the self-increasing id of the updated model, it is the bk_obj_id of the corresponding model, and if it is a global self-increasing id, it is "global". |
| sequence_id | int    | Yes      | self-increment id                                                                                                                                               |

### Request Example

```json
{
  "type": "host",
  "sequence_id": 1000000
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "",
  "permission": null,
  "data": null
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
