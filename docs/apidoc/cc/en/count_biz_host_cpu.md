### Function Description

Count the number of CPUs for each business's hosts (Special interface for cost management, v3.8.17+/v3.10.18+, Permission: Global Settings Permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description        |
| --------- | ------ | -------- | ------------------ |
| bk_biz_id | int    | No       | Business ID        |
| page      | object | No       | Paging information |

**Note: The bk_biz_id and page parameters must be provided, and only one of them can be passed.**

#### Page Field Description

| Field | Type | Required | Description                          |
| ----- | ---- | -------- | ------------------------------------ |
| start | int  | Yes      | Record start position                |
| limit | int  | Yes      | Number of records per page, up to 10 |

### Request Parameter Example

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

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
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

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error  |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Data returned in the request                                 |

#### data

| Field             | Type | Description                                 |
| ----------------- | ---- | ------------------------------------------- |
| bk_biz_id         | int  | Business ID                                 |
| host_count        | int  | Number of hosts                             |
| cpu_count         | int  | Number of CPUs                              |
| no_cpu_host_count | int  | Number of hosts without the CPU count field |