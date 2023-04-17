### Functional description

Query the service template ID corresponding to the host. This interface is dedicated to node management and may be adjusted at any time. Do not use other services (v3.10.11+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Parameter       | Type| Required| Description                |
| ---------- | ----- | ---- | ------------------- |
| bk_host_id | array |yes   | Host IDs, up to 200|

#### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_host_id": [
        258,
        259
    ]
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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

### Return Result Parameters Description

#### response

| Field                | Type| Description       |
| ------------------- | ----- | ---------- |
| result     |  bool   | Whether the request was successful or not. True: request succeeded;false: Request failed|
| code       |  int    | Wrong. 0 indicates success,>0 indicates failure error        |
| message    |  string |Error message returned by request failure                        |
| permission | object |Permission information                                      |
| request_id | string |Request chain id                                      |
| data       |  array  |Request result                                      |

#### data

| Field                | Type| Description       |
| ------------------- | ----- | ---------- |
| bk_host_id          |  int   | Host id     |
| service_template_id | array |Service template id|
