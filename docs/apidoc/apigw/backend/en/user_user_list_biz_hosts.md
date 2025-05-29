### Description

Query hosts under a business based on business ID, with optional additional filtering information such as cluster ID,
module ID, etc. (Permission: Business access permission)

### Parameters

| Name                 | Type   | Required | Description                                                                                                                                        |
|----------------------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------|
| page                 | object | Yes      | Query conditions                                                                                                                                   |
| bk_biz_id            | int    | Yes      | Business ID                                                                                                                                        |
| bk_set_ids           | array  | No       | Cluster ID list, up to 200 items **bk_set_ids and set_cond can only use one of them**                                                              |
| set_cond             | array  | No       | Cluster query condition **bk_set_ids and set_cond can only use one of them**                                                                       |
| bk_module_ids        | array  | No       | Module ID list, up to 500 items **bk_module_ids and module_cond can only use one of them**                                                         |
| module_cond          | array  | No       | Module query condition **bk_module_ids and module_cond can only use one of them**                                                                  |
| host_property_filter | object | No       | Combined query conditions for host properties                                                                                                      |
| fields               | array  | Yes      | List of host properties, control which fields are returned in the result to speed up the interface request and reduce network traffic transmission |

#### host_property_filter

This parameter is a combination of filtering rules for host property fields, used to search for hosts based on host
property fields. The combination supports both AND and OR, and can be nested, with a maximum of 2 levels. The filtering
rule is a quadruple `field`, `operator`, `value`

| Name      | Type   | Required | Description              |
|-----------|--------|----------|--------------------------|
| condition | string | No       | Combined query condition |
| rules     | array  | No       | Rules                    |

#### rules

| Name     | Type   | Required | Description                                                                                                                  |
|----------|--------|----------|------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Field name                                                                                                                   |
| operator | string | Yes      | Operator, optional values equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats                                                           |

Assembly rules can refer to: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### set_cond

| Name     | Type   | Required | Description                                                          |
|----------|--------|----------|----------------------------------------------------------------------|
| field    | string | Yes      | Value is the field of the cluster                                    |
| operator | string | Yes      | Optional values are $eq, $ne                                         |
| value    | string | Yes      | The value corresponding to the field configured as the cluster field |

#### module_cond

| Name     | Type   | Required | Description                                                         |
|----------|--------|----------|---------------------------------------------------------------------|
| field    | string | Yes      | Value is the field of the module                                    |
| operator | string | Yes      | Optional values are $eq, $ne                                        |
| value    | string | Yes      | The value corresponding to the field configured as the module field |

#### page

| Name  | Type   | Required | Description                 |
|-------|--------|----------|-----------------------------|
| start | int    | Yes      | Record start position       |
| limit | int    | Yes      | Limit per page, maximum 500 |
| sort  | string | No       | Sorting field               |

### Request Example

```json
{
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "bk_host_id"
    },
    "set_cond": [
        {
            "field": "bk_set_name",
            "operator": "$eq",
            "value": "set1"
        }
    ],
    "bk_biz_id": 3,
    "bk_module_ids": [54,56],
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
                "field": "bk_host_innerip",
                "operator": "equal",
                "value": "127.0.0.1"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                        "field": "bk_os_type",
                        "operator": "not_in",
                        "value": [
                            "3"
                        ]
                    },
                    {
                        "field": "bk_cloud_id",
                        "operator": "equal",
                        "value": 0
                    }
                ]
            }
        ]
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
    "count": 2,
    "info": [
      {
        "bk_cloud_id": 0,
        "bk_host_id": 1,
        "bk_host_innerip": "192.168.15.18",
        "bk_mac": "",
        "bk_os_type": null
      },
      {
        "bk_cloud_id": 0,
        "bk_host_id": 2,
        "bk_host_innerip": "192.168.15.4",
        "bk_mac": "",
        "bk_os_type": null
      }
    ]
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                      |
|------------|--------|------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error    |
| message    | string | Error message returned in case of request failure                |
| permission | object | Permission information                                           |
| data       | object | Data returned in the request                                     |

#### data

| Name  | Type  | Description       |
|-------|-------|-------------------|
| count | int   | Number of records |
| info  | array | Host actual data  |

#### data.info

| Name                 | Type   | Description                               |
|----------------------|--------|-------------------------------------------|
| bk_host_name         | string | Host name                                 |
| bk_host_innerip      | string | Intranet IP                               |
| bk_host_id           | int    | Host ID                                   |
| bk_cloud_id          | int    | Control area                              |
| import_from          | string | Host import source, imported as 3 via API |
| bk_asset_id          | string | Fixed asset number                        |
| bk_cloud_inst_id     | string | Cloud host instance ID                    |
| bk_cloud_vendor      | string | Cloud vendor                              |
| bk_cloud_host_status | string | Cloud host status                         |
| bk_comment           | string | Remarks                                   |
| bk_cpu               | int    | CPU logical cores                         |
| bk_cpu_architecture  | string | CPU architecture                          |
| bk_cpu_module        | string | CPU model                                 |
| bk_disk              | int    | Disk capacity (GB)                        |
| bk_host_outerip      | string | Host public IP                            |
| bk_host_innerip_v6   | string | Host intranet IPv6                        |
| bk_host_outerip_v6   | string | Host public IPv6                          |
| bk_isp_name          | string | ISP name                                  |
| bk_mac               | string | Host intranet MAC address                 |
| bk_mem               | int    | Host memory capacity (MB)                 |
| bk_os_bit            | string | Operating system bit                      |
| bk_os_name           | string | Operating system name                     |
| bk_os_type           | string | Operating system type                     |
| bk_os_version        | string | Operating system version                  |
| bk_outer_mac         | string | Host public MAC address                   |
| bk_province_name     | string | Province where the host is located        |
| bk_service_term      | int    | Warranty period                           |
| bk_sla               | string | SLA level                                 |
| bk_sn                | string | Device SN                                 |
| bk_state             | string | Current state                             |
| bk_state_name        | string | Country where the host is located         |
| operator             | string | Main maintainer                           |
| bk_bak_operator      | string | Backup maintainer                         |

**Note: The returned values here only explain the system's built-in property fields. Other returned values depend on the
user's self-defined property fields.**
