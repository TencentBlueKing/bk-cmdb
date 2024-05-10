### Description

Get the result of pushing host identity to machines (can only get tasks pushed within 30 minutes) (Version: v3.10.23+,
Permission: When the hosts included in the task belong to a business, the corresponding business access permission is
required; when the hosts belong to a host pool, host update permission is required)

### Parameters

| Name    | Type   | Required | Description |
|---------|--------|----------|-------------|
| task_id | string | Yes      | Task ID     |

### Request Example

```json
{
    "task_id": "GSETASK:F:202201251046313618521052:198"
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "msg": "success",
    "permission": null,
    "data": {
            "success_list": [
                1,
                2
            ],
            "pending_list": [
                3,
                4
            ],
            "failed_list": [
                5,
                6
            ]
        }
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

#### data Field Explanation

| Name         | Type  | Description                                                                                              |
|--------------|-------|----------------------------------------------------------------------------------------------------------|
| success_list | array | List of host IDs that executed successfully                                                              |
| failed_list  | array | List of host IDs that failed to execute                                                                  |
| pending_list | array | List of host IDs for which the host identity was called by GSE, but the result has not been obtained yet |
