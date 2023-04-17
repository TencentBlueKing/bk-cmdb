### Functional description

Unbind agent to host (v3.10.25+).

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field | Type         | Required | Description                                                   |
| ----- | ------------ | -------- | ------------------------------------------------------------- |
| list  | object array | yes      | list of host IDs and agent IDs to bind, maximum length is 200 |

#### list

| Field       | Type   | Required | Description           |
| ----------- | ------ | -------- | --------------------- |
| bk_host_id  | int    | yes      | host ID to bind agent |
| bk_agent_id | string | yes      | agent ID to bind host |

### Request Parameters Example

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

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| data    |  object |Data returned by request                           |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |