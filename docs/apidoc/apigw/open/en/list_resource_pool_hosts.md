### Description

Query hosts in the resource pool (Permission: View permission for hosts in the host pool)

### Parameters

| Name                 | Type   | Required | Description                                                                                                                                                                                                |
|----------------------|--------|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| page                 | dict   | No       | Query conditions                                                                                                                                                                                           |
| host_property_filter | object | No       | Combination query conditions for host properties                                                                                                                                                           |
| fields               | array  | No       | List of host properties, controls which fields are returned in the result, speeding up interface requests and reducing network traffic transmission. If not filled in, all fields are returned by default. |

#### host_property_filter

This parameter is a combination of filtering rules for host property fields, used to search hosts based on host property
fields. The combination supports both AND and OR, can be nested, and has a maximum nesting of 2 levels.

| Name      | Type   | Required | Description |
|-----------|--------|----------|-------------|
| condition | string | No       |             |
| rules     | array  | No       | Rules       |

#### rules

Filtering rules are quadruples `field`, `operator`, `value`

| Name     | Type   | Required | Description                                                                                                                      |
|----------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Field name                                                                                                                       |
| operator | string | Yes      | Operator, optional values are equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats                                                               |

Assembly rules can refer to: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### page

| Name  | Type   | Required | Description                                |
|-------|--------|----------|--------------------------------------------|
| start | int    | Yes      | Record start position                      |
| limit | int    | Yes      | Number of records per page, maximum is 500 |
| sort  | string | No       | Sorting field                              |

### Request Example

```json
{
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "bk_host_id"
    },
    "fields": [
        "bk_host_id",
        "bk_cloud_id",
        "bk_host_innerip",
        "bk_os_type",
        "bk_mac"
    ],
    "host_property_filter": {
        "condition": "AND",
        "rules": [
        {
            "field": "bk_host_outerip",
            "operator": "equal",
            "value": "127.0.0.1"
        }, {
            "condition": "OR",
            "rules": [{
                "field": "bk_os_type",
                "operator": "not_in",
                "value": ["3"]
            }, {
                "field": "bk_sla",
                "operator": "equal",
                "value": "1"
            }]
        }]
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
    "count": 1,
    "info": [
      {
        "bk_cloud_id": "0",
        "bk_host_id": 17,
        "bk_host_innerip": "192.168.1.1",
        "bk_mac": "",
        "bk_os_type": "1"
      }
    ]
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | array  | Data returned by the request                                       |

#### data

| Name  | Type  | Description                |
|-------|-------|----------------------------|
| count | int   | Number of records          |
| info  | array | Actual data, list of hosts |

#### data.info

| Name                 | Type   | Description                          |
|----------------------|--------|--------------------------------------|
| bk_host_name         | string | Host name                            |
| bk_host_innerip      | string | Private IP address                   |
| bk_host_id           | int    | Host ID                              |
| bk_cloud_id          | int    | Control area                         |
| import_from          | string | Host import source, 3 for API import |
| bk_asset_id          | string | Fixed asset number                   |
| bk_cloud_inst_id     | string | Cloud host instance ID               |
| bk_cloud_vendor      | string | Cloud vendor                         |
| bk_cloud_host_status | string | Cloud host status                    |
| bk_comment           | string | Comment                              |
| bk_cpu               | int    | Number of CPU logical cores          |
| bk_cpu_architecture  | string | CPU architecture                     |
| bk_cpu_module        | string | CPU model                            |
| bk_disk              | int    | Disk capacity (GB)                   |
| bk_host_outerip      | string | Host public IP address               |
| bk_host_innerip_v6   | string | Host private IPv6 address            |
| bk_host_outerip_v6   | string | Host public IPv6 address             |
| bk_isp_name          | string | ISP name                             |
| bk_mac               | string | Host private MAC address             |
| bk_mem               | int    | Host memory capacity (MB)            |
| bk_os_bit            | string | Operating system bit number          |
| bk_os_name           | string | Operating system name                |
| bk_os_type           | string | Operating system type                |
| bk_os_version        | string | Operating system version             |
| bk_outer_mac         | string | Host public MAC address              |
| bk_province_name     | string | Province                             |
| bk_service_term      | int    | Warranty period (years)              |
| bk_sla               | string | SLA level                            |
| bk_sn                | string | Device SN                            |
| bk_state             | string | Current state                        |
| bk_state_name        | string | Country                              |
| operator             | string | Primary maintainer                   |
| bk_bak_operator      | string | Backup maintainer                    |

**Note: The returned values here only explain the system-built attribute fields, and the rest of the returned values
depend on the user-defined attribute fields.**
