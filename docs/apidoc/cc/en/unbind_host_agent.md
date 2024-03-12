### Function Description

Unbind agent from host (Version: v3.10.25+, Permission: Host AgentID Management Permission)

### Request Parameters

{{ common_args_desc }}

### Request Parameters

| Field | Type  | Required | Description                                          |
| ----- | ----- | -------- | ---------------------------------------------------- |
| list  | array | Yes      | List of host IDs and agent IDs to unbind (up to 200) |

### list

| Field       | Type   | Required | Description                    |
| ----------- | ------ | -------- | ------------------------------ |
| bk_host_id  | int    | Yes      | Host ID of the agent to unbind |
| bk_agent_id | string | Yes      | Agent ID to unbind             |

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

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807"
}
```

### Return Result Parameter Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |