### Description

Add the module for idle hosts in the business related to cloud hosts (Dedicated interface for cloud host management,
Version: v3.10.19+, Permission: Business host editing permission)

### Parameters

| Name      | Type  | Required | Description                                                                                                                    |
|-----------|-------|----------|--------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id | int   | Yes      | Business ID                                                                                                                    |
| host_info | array | Yes      | Information of newly added cloud hosts, array length can be up to 200, success or failure occurs for the entire batch of hosts |

#### host_info

Host information, where bk_cloud_id, bk_host_innerip, cloud vendor, and cloud host instance ID fields are required.
Other fields are attributes defined in the host model. Only a subset of fields is shown here, please fill in other
fields as needed.

| Name                 | Type   | Required | Description                                                             |
|----------------------|--------|----------|-------------------------------------------------------------------------|
| bk_cloud_id          | int    | Yes      | Control area ID                                                         |
| bk_host_innerip      | string | Yes      | IPv4 format of host's internal IP, separated by commas for multiple IPs |
| bk_cloud_vendor      | array  | Yes      | Cloud vendor                                                            |
| bk_cloud_inst_id     | array  | Yes      | Cloud host instance ID                                                  |
| bk_addressing        | string | No       | Addressing method, static for cloud hosts                               |
| bk_host_name         | string | No       | Hostname or other attributes                                            |
| operator             | string | No       | Main maintainer or other attributes                                     |
| bk_comment           | string | No       | Remark or other attributes                                              |
| import_from          | string | No       | Host import source, 3 for API import                                    |
| bk_asset_id          | string | No       | Fixed asset number                                                      |
| bk_created_at        | string | No       | Creation time                                                           |
| bk_updated_at        | string | No       | Update time                                                             |
| bk_created_by        | string | No       | Creator                                                                 |
| bk_updated_by        | string | No       | Updater                                                                 |
| bk_cloud_host_status | string | No       | Cloud host status                                                       |
| bk_cpu               | int    | No       | CPU logical cores                                                       |
| bk_cpu_architecture  | string | No       | CPU architecture                                                        |
| bk_cpu_module        | string | No       | CPU model                                                               |
| bk_disk              | int    | No       | Disk capacity (GB)                                                      |
| bk_host_outerip      | string | No       | Host's external IP                                                      |
| bk_host_innerip_v6   | string | No       | Host's internal IPv6                                                    |
| bk_host_outerip_v6   | string | No       | Host's external IPv6                                                    |
| bk_isp_name          | string | No       | Affiliated ISP                                                          |
| bk_mac               | string | No       | Host's internal MAC address                                             |
| bk_mem               | int    | No       | Host's RAM capacity (MB)                                                |
| bk_os_bit            | string | No       | Operating system bit                                                    |
| bk_os_name           | string | No       | Operating system name                                                   |
| bk_os_type           | string | No       | Operating system type                                                   |
| bk_os_version        | string | No       | Operating system version                                                |
| bk_outer_mac         | string | No       | Host's external MAC address                                             |
| bk_province_name     | string | No       | Province where the host is located                                      |
| bk_service_term      | int    | No       | Warranty period                                                         |
| bk_sla               | string | No       | SLA level                                                               |
| bk_sn                | string | No       | Device SN                                                               |
| bk_state             | string | No       | Current state                                                           |
| bk_state_name        | string | No       | Country where the host is located                                       |
| bk_bak_operator      | string | No       | Backup maintainer                                                       |

### Request Example

```json
{
    "bk_biz_id": 123,
    "host_info": [
        {
            "bk_cloud_id": 0,
            "bk_host_innerip": "127.0.0.1",
            "bk_cloud_vendor": "2",
            "bk_cloud_inst_id": "45515",
            "bk_host_name": "host1",
            "operator": "admin",
            "bk_comment": "comment"
        },
        {
            "bk_cloud_id": 0,
            "bk_host_innerip": "127.0.0.2",
            "bk_cloud_vendor": "2",
            "bk_cloud_inst_id": "45656",
            "bk_host_name": "host2",
            "operator": "admin",
            "bk_comment": "comment"
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
    "data": {
        "ids": [
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

| Name | Type  | Description                                 |
|------|-------|---------------------------------------------|
| ids  | array | Array of IDs for successfully created hosts |
