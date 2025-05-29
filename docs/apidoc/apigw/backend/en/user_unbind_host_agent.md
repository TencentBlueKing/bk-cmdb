### Description

Unbind agent from host (Version: v3.10.25+, Permission: Host AgentID Management Permission)

### Parameters

| Name        | Type   | Required | Description                    |
|-------------|--------|----------|--------------------------------|
| bk_host_id  | int    | Yes      | Host ID of the agent to unbind |
| bk_agent_id | string | Yes      | Agent ID to unbind             |

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

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
