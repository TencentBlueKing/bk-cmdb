### Functional description

list host's detail info and its topology it belongs to with host's attributes

### Request Parameters


#### General Parameters

| Field | Type | Required |  Description |
|-----------|------------|--------|------------|
| bk_app_code  |  string    | Yes | APP ID     |
| bk_app_secret|  string    | Yes | APP Secret(APP TOKEN), which can be got via BlueKing Developer Center -&gt; Click APP ID -&gt; Basic Info  |
| bk_token     |  string    | No | Current user login token, bk_token or bk_username must be valid, bk_token can be got by Cookie |
| bk_username  |  string    | No | Current user username, APP in the white list, can use this field to specify the current user |
| fields  |  array   | Yes     | host property list, the specified host property feilds will be returned <br>it can speed up the request and reduce the network payload  |

#### Interface Parameters

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| page       |  dict    | Yes     | search condition |
| fields       |  array string    | Yes     | host's attribute list need to return |
| host_property_filter    |  dict  | No     | host property filter |


#### host_property_filter
host property filter is a combined of atom filter rule, combine operator could be `AND` or `OR`, nested up to 2 levels。
atom rule has three fields: `field`, `operator`, `value`

| Field      |  Type      | Required   |  Description      |
| ---  | ---  | --- |---  |
| field|string|Yes|field |
| operator|string|No|operator |available values: equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between,contains,exists,not_exists; contains is a regular pattern operation which is case in-sensitive; exists/not_exists operator is to filter out a attribute is exists in a instances or not. |
| value| - | No| value|values's format depend on operator|

reference: <https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md>

#### page

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| start    |  int    | Yes     | start record |
| limit    |  int    | Yes     | page limit, maximum value is 500 |
| sort     |  string | No     | the field for sort |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_token": "xxx",
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
### Return Result Example

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
                            "name": "china",
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
                            "name": "cn",
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
                            "name": "IdleSet",
                            "id": 2
                        },
                        "children": [
                            {
                                "inst": {
                                    "obj": "module",
                                    "name": "FaultModule",
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

### Return Result Parameters Description

#### data

| Field      | Type      | Description      |
|-----------|-----------|-----------|
| count     | int       | the num of record |
| info      | array     | host data and topology information |

#### data.info
| Field      | Type      | Description      |
|-----------|-----------|-----------|
| host     | dict       | the num of record |
| topo      | array     | host's topology info |

#### data.info.host
| Field      | Type      | Description      |
|---|---|---|
| bk_isp_name| string | telecom operators | 0: Others; 1: China Telecom; 2: China Unicom; 3: China Mobile |
| bk_sn | string | device SN |
| operator | string | maintainer |
| bk_outer_mac | string | outer MAC |
| bk_state_name | string | country | CN: China, please refer to CMDB web page for detailed value |
| bk_province_name | string | province |  |
| import_from | string | import from | 1:excel;2:agent;3:api |
| bk_sla | string | SLA level | 1:L1;2:L2;3:L3 |
| bk_service_term | int | warranty | 1-10 |
| bk_os_type | string | os type | 1:Linux;2:Windows;3:AIX |
| bk_os_version | string | os version |
| bk_os_bit | int | os bits |
| bk_mem | string | memory capacity |
| bk_mac | string | mac address |
| bk_host_outerip | string | outer ip |
| bk_host_name | string | hostname |
| bk_host_innerip | string | inner ip |
| bk_host_id | int | host id |
| bk_disk | int | disk capacity | |
| bk_cpu_module | string | CPU module |
| bk_cpu_mhz | int | CPU hz | |
| bk_cpu | int | CPU count | 1-1000000 |
| bk_comment | string | comment |
| bk_cloud_id | int | cloud area id |
| bk_bak_operator | string | backup maintainer |
| bk_asset_id | string | device id |

#### data.info.topo
| Field      | Type      | Description      |
|-----------|-----------|-----------|
| inst | object       | this node's instance's details information |
| inst.obj | string       | this node's object type, can be set/module/custom_level etc. |
| inst.name | string  | node's name |
| inst.id    | int     | node's instance id |
| children | object array       | node's children details topology information |
| children.inst | object       | this child's instance's information |
| children.children | string  | this instance's children details information |