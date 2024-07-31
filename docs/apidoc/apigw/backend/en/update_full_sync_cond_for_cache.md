### Description

Update full synchronization cache condition (version: v3.14.1+, permission: update permission for full sync cache cond)

### Parameters

| Name | Type   | Required | Description                                             |
|------|--------|----------|---------------------------------------------------------|
| id   | int    | yes      | ID of the full sync cache cond that needs to be updated |
| data | object | Yes      | Data that needs to be updated                           |

#### data

| Name     | Type | Required                                                                                                        | Description |
|----------|------|-----------------------------------------------------------------------------------------------------------------|-------------|
| interval | int  | Sync period, in hours, used to specify the cache expiration time, the minimum is 6 hours, the maximum is 7 days |

### Request Example

```json
{
   "id": 123,
   "data": {
     "interval": 24
   }
}
```

### Response Example

```json
{
   "result": true,
   "code": 0,
   "message": "success",
   "permission": null,
   "data": null
}
```

### Response Parameters

| Name       | Type   | Description                                                      |
|------------|--------|------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error    |
| message    | string | Error message returned in case of request failure                |
| permission | object | Permission information                                           |
| data       | object | Data returned in the request                                     |
