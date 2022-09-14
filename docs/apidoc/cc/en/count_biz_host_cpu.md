### Functional description

count host cpu num in biz (special interface for cost managing, v3.8.17+/v3.10.17+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description        |
| --------- | ------ | -------- | ------------------ |
| bk_biz_id | int    | no       | Business ID        |
| page      | object | no       | Paging information |

**Note: only one of bk_biz_id and page parameters must set**

#### Page field Description

| Field | Type | Required | Description           |
| ----- | ---- | -------- | --------------------- |
| start | int  | yes      | Record start position |
| limit | int  | yes      | Page limit, maxium 10 |


### Request Parameters Example

```json
{
    "bk_app_code": "code",
    "bk_app_secret": "secret",
    "bk_username": "xxx",
    "bk_token": "xxxx",
    "page": {
        "start": 10,
        "limit": 10
    }
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": [
        {
            "bk_biz_id": 5,
            "host_count": 100,
            "cpu_count": 192,
            "no_cpu_host_count": 5
        },
        {
            "bk_biz_id": 7,
            "host_count": 40,
            "cpu_count": 58,
            "no_cpu_host_count": 11
        }
    ]
}
```

### Return result parameter

#### response

| Name       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request was successful or not. True: request succeeded;false request failed |
| code       | int    | Wrong code. 0 indicates success,>0 indicates failure error   |
| message    | string | Error message returned by request failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | object | Data returned by request                                     |

#### data

| Field             | Type | Description                                     |
| ----------------- | ---- | ----------------------------------------------- |
| bk_biz_id         | int  | no                                              |
| host_count        | int  | The number of hosts in the biz                  |
| cpu_count         | int  | The number of host cpus in the biz              |
| no_cpu_host_count | int  | The number of hosts with no cpu info in the biz |
