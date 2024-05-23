### Description

Query host details and its topology information based on host condition information (Permission: Host pool host view
permission)

### Parameters

| Name                 | Type   | Required | Description                                                                                                 |
|----------------------|--------|----------|-------------------------------------------------------------------------------------------------------------|
| page                 | dict   | Yes      | Query conditions                                                                                            |
| host_property_filter | object | No       | Combined query conditions for host properties                                                               |
| fields               | array  | Yes      | List of host properties, control which fields are in the returned result, fill in according to requirements |

#### host_property_filter

This parameter is a combination of filtering rules for host property fields, used to search for hosts based on host
property fields. The combination supports AND and OR, and can be nested up to two levels. Filtering rules are
quadruples `field`, `operator`, `value`

| Name      | Type   | Required | Description               |
|-----------|--------|----------|---------------------------|
| condition | string | No       | Combined query conditions |
| rules     | array  | No       | Rules                     |

#### rules

| Name     | Type   | Required | Description                                                                                                                      |
|----------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------|
| field    | string | Yes      | Field name                                                                                                                       |
| operator | string | Yes      | Operator, optional values are equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Operand, different operators correspond to different value formats                                                               |

Assembly rules can refer to: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### page

| Name  | Type   | Required | Description                                      |
|-------|--------|----------|--------------------------------------------------|
| start | int    | Yes      | Record start position                            |
| limit | int    | Yes      | Number of records per page, maximum value is 500 |
| sort  | string | No       | Sorting field                                    |

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
        "bk_host_innerip"
    ],
    "host_property_filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "bk_host_innerip",
                "operator": "equal",
                "value": "192.168.1.1"
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
                "host": {
                    "bk_host_id": 2,
                    "bk_host_innerip": "192.168.1.1"
                },
                "topo": [
                    {
                        "inst": {
                            "obj": "nation",
                            "name": "中国",
                            "id": 30
                        },
                        "children": [
                            {
                                "inst": {
                                    "obj": "province",
                                    "name": "prov-xxx",
                                    "id": 31
                                },
                                "children": [
                                    {
                                        "inst": {
                                            "obj": "set",
                                            "name": "set-xxx",
                                            "id": 20
                                        },
                                        "children": [
                                            {
                                                "inst": {
                                                    "obj": "module",
                                                    "name": "mod-xxx",
                                                    "id": 52
                                                },
                                                "children": null
                                            },
                                            {
                                                "inst": {
                                                    "obj": "module",
                                                    "name": "mod-yy",
                                                    "id": 53
                                                },
                                                "children": null
                                            }
                                        ]
                                    }
                                ]
                            }
                        ]
                    },
                    {
                        "inst": {
                            "obj": "nation",
                            "name": "国家",
                            "id": 29
                        },
                        "children": [
                            {
                                "inst": {
                                    "obj": "province",
                                    "name": "prv1",
                                    "id": 26
                                },
                                "children": [
                                    {
                                        "inst": {
                                            "obj": "set",
                                            "name": "set11",
                                            "id": 19
                                        },
                                        "children": [
                                            {
                                                "inst": {
                                                    "obj": "module",
                                                    "name": "m22",
                                                    "id": 51
                                                },
                                                "children": null
                                            }
                                        ]
                                    }
                                ]
                            }
                        ]
                    }
                ]
            },
            {
                "host": {
                    "bk_host_id": 4,
                    "bk_host_innerip": "192.168.1.2"
                },
                "topo": [
                    {
                        "inst": {
                            "obj": "set",
                            "name": "空闲机池",
                            "id": 2
                        },
                        "children": [
                            {
                                "inst": {
                                    "obj": "module",
                                    "name": "故障机",
                                    "id": 4
                                },
                                "children": null
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

| Name       | Type   | Description                                                      |
|------------|--------|------------------------------------------------------------------|
| result     | bool   | Success or failure of the request. true: success; false: failure |
| code       | int    | Error code. 0 represents success, >0 represents failure error    |
| message    | string | Error message returned in case of failure                        |
| permission | object | Permission information                                           |
| data       | object | Data returned by the request                                     |

#### data

| Name  | Type  | Description                        |
|-------|-------|------------------------------------|
| count | int   | Number of records                  |
| info  | array | Host data and topology information |

#### data.info

| Name | Type  | Description               |
|------|-------|---------------------------|
| host | dict  | Actual host data          |
| topo | array | Host topology information |

#### data.info.host

| Name                 | Type   | Description                          |
|----------------------|--------|--------------------------------------|
| bk_host_name         | string | Host name                            |
| bk_host_innerip      | string | Internal IP of the host              |
| bk_host_id           | int    | Host ID                              |
| bk_cloud_id          | int    | Managed area                         |
| import_from          | string | Host import source, 3 for API import |
| bk_asset_id          | string | Fixed asset number                   |
| bk_cloud_inst_id     | string | Cloud host instance ID               |
| bk_cloud_vendor      | string | Cloud vendor                         |
| bk_cloud_host_status | string | Cloud host status                    |
| bk_comment           | string | Comment                              |
| bk_cpu               | int    | CPU logical cores                    |
| bk_cpu_architecture  | string | CPU architecture                     |
| bk_cpu_module        | string | CPU model                            |
| bk_disk              | int    | Disk capacity (GB)                   |
| bk_host_outerip      | string | Host external IP                     |
| bk_host_innerip_v6   | string | Host internal IPv6                   |
| bk_host_outerip_v6   | string | Host external IPv6                   |
| bk_isp_name          | string | ISP name                             |
| bk_mac               | string | Host internal MAC address            |
| bk_mem               | int    | Host memory capacity (MB)            |
| bk_os_bit            | string | Operating system bit number          |
| bk_os_name           | string | Operating system name                |
| bk_os_type           | string | Operating system type                |
| bk_os_version        | string | Operating system version             |
| bk_outer_mac         | string | Host external MAC address            |
| bk_province_name     | string | Province where the host is located   |
| bk_service_term      | int    | Warranty period                      |
| bk_sla               | string | SLA level                            |
| bk_sn                | string | Device SN                            |
| bk_state             | string | Current status                       |
| bk_state_name        | string | Country where the host is located    |
| operator             | string | Main maintainer                      |
| bk_bak_operator      | string | Backup maintainer                    |

**Note: The return value here only explains the system-built property fields, other return values depend on user-defined
property fields**

#### data.info.topo

| Name              | Type         | Description                                                                   |
|-------------------|--------------|-------------------------------------------------------------------------------|
| inst              | object       | Details of the node instance                                                  |
| inst.obj          | string       | Model type of the node, such as set/module and custom hierarchical model type |
| inst.name         | string       | Instance name of the node                                                     |
| inst.id           | int          | Instance ID of the node                                                       |
| children          | object array | Details of child nodes of the current instance, there may be multiple         |
| children.inst     | object       | Details of the instance of the child node                                     |
| children.children | string       | Details of child nodes of the current instance                                |
