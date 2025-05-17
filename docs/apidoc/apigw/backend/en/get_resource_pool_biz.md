### Description

Query resource pool business info (version: v3.15.1+, permission: View resource pool hosts permission)

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": {
    "bk_biz_id": 1
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| data       | object | Request returned data                                               |
| permission | object | Permission information                                              |

#### data

| Name      | Type | Description |
|-----------|------|-------------|
| bk_biz_id | int  | Business ID |
