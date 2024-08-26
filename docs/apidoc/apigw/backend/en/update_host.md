### Description

Update Host Properties (Permission: For hosts already assigned to a business, business host editing permission is
required. For host pool hosts, host pool host editing permission is required)

### Parameters

| Name                | Type   | Required | Description                  |
|---------------------|--------|----------|------------------------------|
| bk_host_id          | string | Yes      | Host ID, separated by commas |
| bk_host_name        | string | No       | Host name                    |
| bk_comment          | string | No       | Comment                      |
| bk_cpu              | int    | No       | CPU logical cores            |
| bk_cpu_architecture | string | No       | CPU architecture             |
| bk_cpu_module       | string | No       | CPU model                    |
| bk_disk             | int    | No       | Disk capacity (GB)           |
| bk_host_outerip     | string | No       | Host outer IP                |
| bk_host_outerip_v6  | string | No       | Host outer IPv6              |
| bk_isp_name         | string | No       | ISP name                     |
| bk_mac              | string | No       | Host inner MAC address       |
| bk_mem              | int    | No       | Host memory capacity (MB)    |
| bk_os_bit           | string | No       | Operating system bit         |
| bk_os_name          | string | No       | Operating system name        |
| bk_os_type          | string | No       | Operating system type        |
| bk_os_version       | string | No       | Operating system version     |
| bk_outer_mac        | string | No       | Host outer MAC address       |
| bk_province_name    | string | No       | Province name                |
| bk_sla              | string | No       | SLA level                    |
| bk_sn               | string | No       | Device SN                    |
| bk_state            | string | No       | Current state                |
| bk_state_name       | string | No       | Country                      |
| operator            | string | No       | Main maintainer              |
| bk_bak_operator     | string | No       | Backup maintainer            |

**Note: The parameter here only explains the system-built editable parameters, and the rest of the parameters to be
filled depend on the user's own defined attribute fields.**

### Request Example

```json
{
  "bk_host_id": "1,2,3",
  "bk_host_name": "test",
  "bk_comment": "canway-host-101",
  "bk_cpu": 16,
  "bk_cpu_architecture": "arm",
  "bk_cpu_module": "Intel(R) 2.00GHz",
  "bk_disk": 120,
  "bk_host_outerip": "12.0.0.3",
  "bk_host_outerip_v6": "0000:0000:0000:0000:0000:0000:0000:0248",
  "bk_isp_name": "3",
  "bk_mac": "00:00:00:00:00:56",
  "bk_mem": 36666,
  "bk_os_bit": "32-bit",
  "bk_os_name": "ubuntu",
  "bk_os_type": "4",
  "bk_os_version": "7.9.1",
  "bk_outer_mac": "00:00:00:00:00:56",
  "bk_province_name": "440000",
  "bk_sla": "2",
  "bk_sn": "abcd3252425",
  "bk_state": "Backup",
  "bk_state_name": "BE",
  "operator": "admin",
  "bk_bak_operator": "admin"
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": null
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
