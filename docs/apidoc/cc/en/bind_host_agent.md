### Function Description

Bind Agent to Host (Version: v3.10.25+, Permission: Host AgentID Management Permission)

### Request Parameters

{{ common_args_desc }}

### Request Parameters

| Field | Type  | Required | Description                                                  |
| ----- | ----- | -------- | ------------------------------------------------------------ |
| list  | array | Yes      | List of host IDs and agent IDs to be bound, up to 200 entries |

### list

| Field       | Type   | Required | Description                            |
| ----------- | ------ | -------- | -------------------------------------- |
| bk_host_id  | int    | Yes      | Host ID to bind the agent to           |
| bk_agent_id | string | Yes      | Agent ID to bind to the specified host |

### Request Parameter Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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
    "request_id": "e43da4ef221746868dc4c837d36f3807"
}
```

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error  |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |