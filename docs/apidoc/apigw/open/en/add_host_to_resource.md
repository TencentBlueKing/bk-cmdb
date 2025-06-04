### Description

Add hosts to the resource pool (Permission: Host Pool Host Creation Permission)

### Parameters

| Name                | Type   | Required | Description       |
|---------------------|--------|----------|-------------------|
| bk_supplier_account | string | No       | Developer account |
| host_info           | dict   | Yes      | Host information  |
| bk_biz_id           | int    | No       | Business ID       |

#### host_info

| Name                 | Type   | Required | Description                                                                  |
|----------------------|--------|----------|------------------------------------------------------------------------------|
| bk_host_innerip      | string | Yes      | Host's inner IP address                                                      |
| import_from          | string | No       | Host import source, 3 for API import                                         |
| bk_cloud_id          | int    | No       | Managed area ID, not filled in to add to the default area 0                  |
| bk_addressing        | string | No       | Addressing mode, default to static addressing mode if not filled in (static) |
| bk_host_name         | string | No       | Host name                                                                    |
| bk_asset_id          | string | No       | Fixed asset number                                                           |
| bk_created_at        | string | No       | Creation time                                                                |
| bk_updated_at        | string | No       | Update time                                                                  |
| bk_created_by        | string | No       | Creator                                                                      |
| bk_updated_by        | string | No       | Last updater                                                                 |
| bk_cloud_inst_id     | string | No       | Cloud host instance ID                                                       |
| bk_cloud_vendor      | string | No       | Cloud vendor                                                                 |
| bk_cloud_host_status | string | No       | Cloud host status                                                            |
| bk_comment           | string | No       | Comment                                                                      |
| bk_cpu               | int    | No       | Number of CPU logical cores                                                  |
| bk_cpu_architecture  | string | No       | CPU architecture                                                             |
| bk_cpu_module        | string | No       | CPU model                                                                    |
| bk_disk              | int    | No       | Disk capacity (GB)                                                           |
| bk_host_outerip      | string | No       | Host's outer IP address                                                      |
| bk_host_innerip_v6   | string | No       | Host's inner IPv6 address                                                    |
| bk_host_outerip_v6   | string | No       | Host's outer IPv6 address                                                    |
| bk_isp_name          | string | No       | ISP name                                                                     |
| bk_mac               | string | No       | Host's inner MAC address                                                     |
| bk_mem               | int    | No       | Host's memory capacity (MB)                                                  |
| bk_os_bit            | string | No       | Operating system bit number                                                  |
| bk_os_name           | string | No       | Operating system name                                                        |
| bk_os_type           | string | No       | Operating system type                                                        |
| bk_os_version        | string | No       | Operating system version                                                     |
| bk_outer_mac         | string | No       | Host's outer MAC address                                                     |
| bk_province_name     | string | No       | Province where the host is located                                           |
| bk_service_term      | int    | No       | Warranty period                                                              |
| bk_sla               | string | No       | SLA level                                                                    |
| bk_sn                | string | No       | Device SN                                                                    |
| bk_state             | string | No       | Current status                                                               |
| bk_state_name        | string | No       | Country where the host is located                                            |
| operator             | string | No       | Main maintainer                                                              |
| bk_bak_operator      | string | No       | Backup maintainer                                                            |

**Note: The input parameters here only explain the required and system-built parameters. The rest of the parameters to
be filled in depend on the host's property fields defined by the user. Refer to the configuration of the host's property
fields for the setting of parameter values.**

### Request Example

```json
{
    "bk_biz_id": 3,
    "host_info": {
        "0": {
            "bk_host_innerip": "127.0.0.1",
            "bk_host_name": "host02",
            "bk_cloud_id": 0,
            "import_from": "3",
            "bk_addressing": "dynamic",
            "bk_asset_id":"udschdfhebv",
            "bk_created_at": "",
            "bk_updated_at": "",
            "bk_created_by": "admin",
            "bk_updated_by": "admin",
            "bk_cloud_inst_id": "1",
            "bk_cloud_vendor": "15",
            "bk_cloud_host_status":"2",
            "bk_comment": "canway-host",
            "bk_cpu": 8,
            "bk_cpu_architecture": "x86",
            "bk_cpu_module": "Intel(R) X87",
            "bk_disk": 195,
            "bk_host_outerip": "12.0.0.1",
            "bk_host_innerip_v6": "0000:0000:0000:0000:0000:0000:0000:0234",
            "bk_host_outerip_v6": "0000:0000:0000:0000:0000:0000:0000:0345",
            "bk_isp_name": "1",
            "bk_mac": "00:00:00:00:00:02",
            "bk_mem": 32155,
            "bk_os_bit": "64-bit",
            "bk_os_name": "linux redhat",
            "bk_os_type": "1",
            "bk_os_version": "7.8",
            "bk_outer_mac": "00:00:00:00:00:02",
            "bk_province_name": "110000",
            "bk_service_term": 6,
            "bk_sla": "1",
            "bk_sn": "abcdsd3252425",
            "bk_state": "测试中",
            "bk_state_name": "CN",
            "operator": "admin",
            "bk_bak_operator": "admin"
        }
    }
}
```

In the example, "0" in host_info indicates the row number, which can be increased sequentially.

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": {
        "success": [
            "0"
        ]
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| data       | object | Request returned data                                               |
| permission | object | Permission information                                              |

#### data

| Name    | Type  | Description            |
|---------|-------|------------------------|
| success | array | Successful row numbers |
