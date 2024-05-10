### Description

Query the service template ID corresponding to the host. This interface is dedicated to node management and may be
adjusted at any time. Please do not use it for other services (Version: v3.10.11+, Permission: Host pool host view
permission)

### Parameters

| Name       | Type  | Required | Description              |
|------------|-------|----------|--------------------------|
| bk_host_id | array | Yes      | Host ID, up to 200 hosts |

### Request Example

```json
{
    "bk_host_id": [
        258,
        259
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
    "data": [
        {
            "bk_host_id": 258,
            "service_template_id": [
                3
            ]
        },
        {
            "bk_host_id": 259,
            "service_template_id": [
                1,
                2
            ]
        }
    ]
}
```

### Response Parameters

| Name       | Type   | Description                                                      |
|------------|--------|------------------------------------------------------------------|
| result     | bool   | Success or failure of the request. true: success; false: failure |
| code       | int    | Error code. 0 represents success, >0 represents failure error    |
| message    | string | Error message returned in case of failure                        |
| permission | object | Permission information                                           |
| data       | array  | Request result                                                   |

#### data

| Name                | Type  | Description         |
|---------------------|-------|---------------------|
| bk_host_id          | int   | Host ID             |
| service_template_id | array | Service template ID |
