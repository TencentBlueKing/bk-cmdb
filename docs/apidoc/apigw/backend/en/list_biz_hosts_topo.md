### Description

Query host and topology information based on the business ID, with optional filtering for clusters, modules, and hosts (
Permission: Business Access Permission).

### Parameters

| Name                   | Type   | Required | Description                                                                                                                         |
|------------------------|--------|----------|-------------------------------------------------------------------------------------------------------------------------------------|
| page                   | object | Yes      | Pagination query conditions, the returned host data is sorted by bk_host_id.                                                        |
| bk_biz_id              | int    | Yes      | Business ID                                                                                                                         |
| set_property_filter    | object | No       | Cluster attribute combination query conditions                                                                                      |
| module_property_filter | object | No       | Module attribute combination query conditions                                                                                       |
| host_property_filter   | object | No       | Host attribute combination query conditions                                                                                         |
| fields                 | array  | Yes      | Host attribute list, control which fields are in the returned result to accelerate the interface request and reduce network traffic |

#### set_property_filter

This parameter is the combination of cluster attribute field filtering rules, used to search for hosts under the cluster
based on cluster attribute fields. The combination supports AND and OR, nested up to 2 layers. The filtering rule is a
quadruple `field`, `operator`, `value`.

| Name      | Type   | Required | Description                 |
|-----------|--------|----------|-----------------------------|
| condition | string | No       | Combination query condition |
| rules     | array  | No       | Rules                       |

#### rules

| Name     | Type   | Required | Description                                                                                                                   |
|----------|--------|----------|-------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Field name                                                                                                                    |
| operator | string | Yes      | Operator, optional values: equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats                                                            |

Assembly rules can refer to: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### module_property_filter

This parameter is the combination of module attribute field filtering rules, used to search for hosts under the module
based on module attribute fields. The combination supports AND and OR, nested up to 2 layers. The filtering rule is a
quadruple `field`, `operator`, `value`.

| Name      | Type   | Required | Description                 |
|-----------|--------|----------|-----------------------------|
| condition | string | No       | Combination query condition |
| rules     | array  | No       | Rules                       |

#### rules

| Name     | Type   | Required | Description                                                                                                                   |
|----------|--------|----------|-------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Field name                                                                                                                    |
| operator | string | Yes      | Operator, optional values: equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats                                                            |

Assembly rules can refer to: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### host_property_filter

This parameter is the combination of host attribute field filtering rules, used to search for hosts based on host
attribute fields. The combination supports AND and OR, nested up to 2 layers. The filtering rule is a
quadruple `field`, `operator`, `value`.

| Name      | Type   | Required | Description                 |
|-----------|--------|----------|-----------------------------|
| condition | string | No       | Combination query condition |
| rules     | array  | No       | Rules                       |

#### rules

| Name     | Type   | Required | Description                                                                                                                   |
|----------|--------|----------|-------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Field name                                                                                                                    |
| operator | string | Yes      | Operator, optional values: equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats                                                            |

Assembly rules can refer to: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### page

| Name  | Type | Required | Description                         |
|-------|------|----------|-------------------------------------|
| start | int  | Yes      | Record start position               |
| limit | int  | Yes      | Records limit per page, maximum 500 |

### Request Example

```json
{
    "page": {
        "start": 0,
        "limit": 10
    },
    "bk_biz_id": 3,
    "fields": [
        "bk_host_id",
        "bk_cloud_id",
        "bk_host_innerip",
        "bk_os_type",
        "bk_mac"
    ],
    "set_property_filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "bk_set_name",
                "operator": "not_equal",
                "value": "test"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                        "field": "bk_set_id",
                        "operator": "in",
                        "value": [
                            1,
                            2,
                            3
                        ]
                    },
                    {
                        "field": "bk_service_status",
                        "operator": "equal",
                        "value": "1"
                    }
                ]
            }
        ]
    },
    "module_property_filter": {
        "condition": "OR",
        "rules": [
            {
                "field": "bk_module_name",
                "operator": "equal",
                "value": "test"
            },
            {
                "condition": "AND",
                "rules": [
                    {
                        "field": "bk_module_id",
                        "operator": "not_in",
                        "value": [
                            1,
                            2,
                            3
                        ]
                    },
                    {
                        "field": "bk_module_type",
                        "operator": "equal",
                        "value": "1"
                    }
                ]
            }
        ]
    },
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
        "count": 3,
        "info": [
            {
                "host": {
                    "bk_cloud_id": 0,
                    "bk_host_id": 1,
                    "bk_host_innerip": "192.168.15.18",
                    "bk_mac": "",
                    "bk_os_type": null
                },
                "topo": [
                    {
                        "bk_set_id": 11,
                        "bk_set_name": "set1",
                        "module": [
                            {
                                "bk_module_id": 56,
                                "bk_module_name": "m1"
                            }
                        ]
                    }
                ]
            },
            {
                "host": {
                    "bk_cloud_id": 0,
                    "bk_host_id": 2,
                    "bk_host_innerip": "192.168.15.4",
                    "bk_mac": "",
                    "bk_os_type": null
                },
                "topo": [
                    {
                        "bk_set_id": 11,
                        "bk_set_name": "set1",
                        "module": [
                            {
                                "bk_module_id": 56,
                                "bk_module_name": "m1"
                            }
                        ]
                    }
                ]
            },
            {
                "host": {
                    "bk_cloud_id": 0,
                    "bk_host_id": 3,
                    "bk_host_innerip": "192.168.15.12",
                    "bk_mac": "",
                    "bk_os_type": null
                },
                "topo": [
                    {
                        "bk_set_id": 10,
                        "bk_set_name": "Idle Pool",
                        "module": [
                            {
                                "bk_module_id": 54,
                                "bk_module_name": "Idle Host"
                            }
                        ]
                    }
                ]
            }
        ]
    }
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 means success, >0 means failed error                 |
| message    | string | Error message for failed requests                                  |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### data

| Name  | Type  | Description                        |
|-------|-------|------------------------------------|
| count | int   | Number of records                  |
| info  | array | Host data and topology information |

#### data.info

| Name | Type  | Description          |
|------|-------|----------------------|
| host | dict  | Actual host data     |
| topo | array | Topology information |

#### data.info.host

| Name                 | Type   | Description                          |
|----------------------|--------|--------------------------------------|
| bk_host_name         | string | Host name                            |
| bk_host_innerip      | string | Internal IP                          |
| bk_host_id           | int    | Host ID                              |
| bk_cloud_id          | int    | Control area                         |
| import_from          | string | Host import source, 3 for API import |
| bk_asset_id          | string | Fixed asset number                   |
| bk_cloud_inst_id     | string | Cloud host instance ID               |
| bk_cloud_vendor      | string | Cloud vendor                         |
| bk_cloud_host_status | string | Cloud host status                    |
| bk_comment           | string | Comment                              |
| bk_cpu               | int    | CPU logical core count               |
| bk_cpu_architecture  | string | CPU architecture                     |
| bk_cpu_module        | string | CPU model                            |
| bk_disk              | int    | Disk capacity (GB)                   |
| bk_host_outerip      | string | Host external IP                     |
| bk_host_innerip_v6   | string | Host internal IPv6                   |
| bk_host_outerip_v6   | string | Host external IPv6                   |
| bk_isp_name          | string | Operator name                        |
| bk_mac               | string | Host internal MAC address            |
| bk_mem               | int    | Host memory capacity (MB)            |
| bk_os_bit            | string | Operating system bit number          |
| bk_os_name           | string | Operating system name                |
| bk_os_type           | string | Operating system type                |
| bk_os_version        | string | Operating system version             |
| bk_outer_mac         | string | Host external MAC address            |
| bk_province_name     | string | Province name                        |
| bk_service_term      | int    | Warranty period                      |
| bk_sla               | string | SLA level                            |
| bk_sn                | string | Device SN                            |
| bk_state             | string | Current status                       |
| bk_state_name        | string | Country name                         |
| operator             | string | Main maintainer                      |
| bk_bak_operator      | string | Backup maintainer                    |

**Note: The return value here only explains the system's built-in attribute fields. Other return values depend on the
user's own defined attribute fields**

#### data.info.topo

| Name        | Type   | Description        |
|-------------|--------|--------------------|
| bk_set_id   | int    | Cluster ID         |
| bk_set_name | string | Cluster name       |
| module      | array  | Module information |

#### data.info.topo.module

| Name           | Type   | Description |
|----------------|--------|-------------|
| bk_module_id   | int    | Module ID   |
| bk_module_name | string | Module name |
