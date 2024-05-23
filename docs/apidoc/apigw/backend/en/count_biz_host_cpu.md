### Description

Count the number of CPUs for each business's hosts (Special interface for cost management, v3.8.17+/v3.10.18+,
Permission: Global Settings Permission)

### Parameters

| Name      | Type   | Required | Description        |
|-----------|--------|----------|--------------------|
| bk_biz_id | int    | No       | Business ID        |
| page      | object | No       | Paging information |

**Note: The bk_biz_id and page parameters must be provided, and only one of them can be passed.**

#### Page Field Description

| Name  | Type | Required | Description                          |
|-------|------|----------|--------------------------------------|
| start | int  | Yes      | Record start position                |
| limit | int  | Yes      | Number of records per page, up to 10 |

### Request Example

```json
{
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

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Data returned in the request                                                |

#### data

| Name              | Type | Description                                 |
|-------------------|------|---------------------------------------------|
| bk_biz_id         | int  | Business ID                                 |
| host_count        | int  | Number of hosts                             |
| cpu_count         | int  | Number of CPUs                              |
| no_cpu_host_count | int  | Number of hosts without the CPU count field |
