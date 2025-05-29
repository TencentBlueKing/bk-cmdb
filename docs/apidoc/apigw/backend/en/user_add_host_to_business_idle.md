### Description

Add hosts to the business idle hosts. This interface ensures that hosts are either added successfully together or fail
together (v3.10.25+, Permission: Host pool host allocation to business permission)

### Parameters

| Name         | Type  | Required | Description                                    |
|--------------|-------|----------|------------------------------------------------|
| bk_host_list | array | Yes      | Host information (array length limited to 200) |
| bk_biz_id    | int   | Yes      | Business ID                                    |

#### bk_host_list (Fields related to hosts)

| Name               | Type   | Required | Description                                                                    |
|--------------------|--------|----------|--------------------------------------------------------------------------------|
| bk_host_innerip    | string | No       | Host's internal IPv4, one of bk_host_innerip or bk_host_innerip_v6 is required |
| bk_host_innerip_v6 | string | No       | Host's internal IPv6, one of bk_host_innerip or bk_host_innerip_v6 is required |
| bk_cloud_id        | int    | Yes      | Control area ID                                                                |
| bk_addressing      | string | Yes      | Addressing method, "static" or "dynamic"                                       |
| operator           | string | No       | Main maintainer                                                                |
| ...                |        |          |                                                                                |

### Request Example

```python
{
    "bk_biz_id": 3,
    "bk_host_list": [
        {
            "bk_host_innerip": "10.0.0.1",
            "bk_cloud_id": 0,
            "bk_addressing": "dynamic",
            "operator": "admin"
        },
        {
            "bk_host_innerip": "10.0.0.2",
            "bk_cloud_id": 0,
            "bk_addressing": "dynamic",
            "operator": "admin"
        }
    ]
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": {
        "bk_host_ids": [
            1,
            2
        ]
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                       |
|------------|--------|-------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error     |
| message    | string | Error message returned for a failed request                       |
| data       | object | Data returned by the request                                      |
| permission | object | Permission information                                            |

#### data

| Name        | Type  | Description                   |
|-------------|-------|-------------------------------|
| bk_host_ids | array | Host IDs of the created hosts |
