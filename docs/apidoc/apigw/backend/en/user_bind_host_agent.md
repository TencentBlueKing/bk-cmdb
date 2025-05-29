### Description

Bind Agent to Host (Version: v3.10.25+, Permission: Host AgentID Management Permission)

### Parameters

| Name        | Type   | Required | Description                            |
|-------------|--------|----------|----------------------------------------|
| bk_host_id  | int    | Yes      | Host ID to bind the agent to           |
| bk_agent_id | string | Yes      | Agent ID to bind to the specified host |

### Request Example

```json
{
    "list": [
        {
            "bk_host_id": 1,
            "bk_agent_id": "xxxxxxxxxx"
        },
        {
            "bk_host_id": 2,
            "bk_agent_id": "yyyyyyyyyy"
        }
    ]
}
```

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
