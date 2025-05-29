### Description

Get hosts under a service template (Version: v3.8.6, Permission: Business Access Permission)

### Parameters

| Name                    | Type   | Required | Description                                                                           |
|-------------------------|--------|----------|---------------------------------------------------------------------------------------|
| bk_biz_id               | int    | Yes      | Business ID                                                                           |
| bk_service_template_ids | array  | Yes      | List of service template IDs, up to 500                                               |
| bk_module_ids           | array  | No       | List of module IDs, up to 500                                                         |
| fields                  | array  | Yes      | List of host attributes, controls which fields are included in the module information |
| page                    | object | Yes      | Pagination information                                                                |

#### page Field Explanation

| Name  | Type | Required | Description                             |
|-------|------|----------|-----------------------------------------|
| start | int  | Yes      | Record start position                   |
| limit | int  | Yes      | Number of records per page, maximum 500 |

### Request Example

```json
{
    "bk_biz_id": 5,
    "bk_service_template_ids": [
        48,
        49
    ],
    "bk_module_ids": [
        65,
        68
    ],
    "fields": [
        "bk_host_id",
        "bk_cloud_id"
    ],
    "page": {
        "start": 0,
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
    "data": {
        "count": 6,
        "info": [
            {
                "bk_cloud_id": 0,
                "bk_host_id": 1
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 2
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 3
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 4
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 7
            },
            {
                "bk_cloud_id": 0,
                "bk_host_id": 8
            }
        ]
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### data

| Name  | Type  | Description       |
|-------|-------|-------------------|
| count | int   | Number of records |
| info  | array | Actual host data  |

#### data.info

| Name                 | Type   | Description                          |
|----------------------|--------|--------------------------------------|
| bk_host_name         | string | Host name                            |
| bk_host_innerip      | string | Inner IP                             |
| bk_host_id           | int    | Host ID                              |
| bk_cloud_id          | int    | Control area                         |
| import_from          | string | Host import source, 3 for API import |
| bk_asset_id          | string | Fixed asset number                   |
| bk_cloud_inst_id     | string | Cloud host instance ID               |
| bk_cloud_vendor      | string | Cloud vendor                         |
| bk_cloud_host_status | string | Cloud host status                    |
| bk_comment           | string | Comment                              |
| bk_cpu               | int    | Logical core count of CPU            |
| bk_cpu_architecture  | string | CPU architecture                     |
| bk_cpu_module        | string | CPU model                            |
| bk_disk              | int    | Disk capacity (GB)                   |
| bk_host_outerip      | string | Host public IP                       |
| bk_host_innerip_v6   | string | Host inner IPv6                      |
| bk_host_outerip_v6   | string | Host public IPv6                     |
| bk_isp_name          | string | Internet service provider            |
| bk_mac               | string | Host inner MAC address               |
| bk_mem               | int    | Host memory capacity (MB)            |
| bk_os_bit            | string | Operating system bit number          |
| bk_os_name           | string | Operating system name                |
| bk_os_type           | string | Operating system type                |
| bk_os_version        | string | Operating system version             |
| bk_outer_mac         | string | Host public MAC address              |
| bk_province_name     | string | Province where the host is located   |
| bk_service_term      | int    | Warranty period                      |
| bk_sla               | string | SLA level                            |
| bk_sn                | string | Device SN                            |
| bk_state             | string | Current status                       |
| bk_state_name        | string | Country where the host is located    |
| operator             | string | Main maintainer                      |
| bk_bak_operator      | string | Backup maintainer                    |

**Note: The explanation of the return values here only covers the system's built-in attribute fields. The rest of the
return values depend on the user's own defined attribute fields** 
