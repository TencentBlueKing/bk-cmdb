### Functional Description

Query the host and its corresponding topology.

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                     | Type   | Required | Description                                                  |
| ------------------------- | ------ | -------- | ------------------------------------------------------------ |
| bk_biz_id                 | int    | Yes      | Business ID                                                  |
| mainline_property_filters | array  | No       | Custom level query filter                                    |
| set_property_filter       | object | No       | Set query filter                                             |
| module_property_filter    | object | No       | Module query filter                                          |
| host_property_filter      | object | No       | Host query filter                                            |
| fields                    | array  | No       | Host property list, which controls which fields are in the host that returns the results |
| page                      | object | Yes      | Paging query condition, returning host data according to bk_ host_ id sorting |

#### mainline_property_filters

| Field     | Type   | Required | Description                    |
| --------- | ------ | -------- | ------------------------------ |
| bk_obj_id | string | Yes      | bk_obj_id of the filter target |
| filter    | object | Yes      | Target object's query filter   |

##### filter

- This parameter is a combination of business attribute field filtering rules, which is used to search for business based on business attribute fields. Combinations only support AND operations and can be nested up to 2 levels.

| Field     | Type   | Required | Description                                     |
| --------- | ------ | -------- | ----------------------------------------------- |
| condition | string | Yes      | Rule operator                                   |
| rules     | array  | Yes      | Scope condition rules for the selected business |

#### rules

| Field    | Type   | Required | Description                                                  |
| -------- | ------ | -------- | ------------------------------------------------------------ |
| field    | string | No       | Field name                                                   |
| operator | string | No       | Available values: equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Values format depends on the operator                        |

Assembly rules are available for reference: [QueryBuilder Rules](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

#### set_property_filter

- This parameter is a combination of business attribute field filtering rules, which is used to search for business based on business attribute fields. Combinations only support AND operations and can be nested up to 2 levels.

| Field     | Type   | Required | Description                                     |
| --------- | ------ | -------- | ----------------------------------------------- |
| condition | string | Yes      | Rule operator                                   |
| rules     | array  | Yes      | Scope condition rules for the selected business |

#### rules

| Field    | Type   | Required | Description                                                  |
| -------- | ------ | -------- | ------------------------------------------------------------ |
| field    | string | No       | Field name                                                   |
| operator | string | No       | Available values: equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Values format depends on the operator                        |

Assembly rules are available for reference: [QueryBuilder Rules](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

#### module_property_filter

- This parameter is a combination of business attribute field filtering rules, which is used to search for business based on business attribute fields. Combinations only support AND operations and can be nested up to 2 levels.

| Field     | Type   | Required | Description                                     |
| --------- | ------ | -------- | ----------------------------------------------- |
| condition | string | Yes      | Rule operator                                   |
| rules     | array  | Yes      | Scope condition rules for the selected business |

#### rules

| Field    | Type   | Required | Description                                                  |
| -------- | ------ | -------- | ------------------------------------------------------------ |
| field    | string | No       | Field name                                                   |
| operator | string | No       | Available values: equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Values format depends on the operator                        |

Assembly rules are available for reference: [QueryBuilder Rules](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

#### host_property_filter

- This parameter is a combination of business attribute field filtering rules, which is used to search for business based on business attribute fields. Combinations only support AND operations and can be nested up to 2 levels.

| Field     | Type   | Required | Description                                     |
| --------- | ------ | -------- | ----------------------------------------------- |
| condition | string | Yes      | Rule operator                                   |
| rules     | array  | Yes      | Scope condition rules for the selected business |

#### rules

| Field    | Type   | Required | Description                                                  |
| -------- | ------ | -------- | ------------------------------------------------------------ |
| field    | string | No       | Field name                                                   |
| operator | string | No       | Available values: equal, not_equal, in, not_in, less, less_or_equal, greater, greater_or_equal, between, not_between |
| value    | -      | No       | Values format depends on the operator                        |

Assembly rules are available for reference: [QueryBuilder Rules](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

#### page

| Field | Type | Required | Description                                     |
| :---- | :--- | :------- | :---------------------------------------------- |
| start | int  | No       | Record start position                           |
| limit | int  | Yes      | Limit the number of entries per page, up to 500 |

### Request Parameters Example

```json{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "mainline_property_filters": [
        {
            "bk_obj_id": "test_custom_level_1",
            "filter": {
                "condition": "AND",
                "rules": [
                    {
                        "field": "bk_inst_name",
                        "operator": "equal",
                        "value": "test2"
                    }
                ]
            }
        }
    ],
    "set_property_filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "bk_set_name",
                "operator": "equal",
                "value": "set1"
            }
        ]
    },
    "module_property_filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "id",
                "operator": "not_equal",
                "value": "2"
            }
        ]
    },
    "host_property_filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "id",
                "operator": "not_equal",
                "value": "17"
            }
        ]
    },
    "page": {
        "start": 0,
        "limit": 500
    }
}
```

### Return Result Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "count": 5,
        "info": [
            {
                "host": {
                    "bk_host_id": 1
                },
                "topo": [
                    {
                        "inst": {
                            "obj": "test_custom_level_2",
                            "name": "test1",
                            "id": 8
                        },
                        "children": [
                            {
                                "inst": {
                                    "obj": "test_custom_level_1",
                                    "name": "test2",
                                    "id": 9
                                },
                                "children": [
                                    {
                                        "inst": {
                                            "obj": "set",
                                            "name": "set1",
                                            "id": 5
                                        },
                                        "children": [
                                            {
                                                "inst": {
                                                    "obj": "module",
                                                    "name": "model1",
                                                    "id": 11
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
                    "bk_host_id": 2
                },
                "topo": [
                    {
                        "inst": {
                            "obj": "test_custom_level_2",
                            "name": "test1",
                            "id": 8
                        },
                        "children": [
                            {
                                "inst": {
                                    "obj": "test_custom_level_1",
                                    "name": "test2",
                                    "id": 9
                                },
                                "children": [
                                    {
                                        "inst": {
                                            "obj": "set",
                                            "name": "set1",
                                            "id": 5
                                        },
                                        "children": [
                                            {
                                                "inst": {
                                                    "obj": "module",
                                                    "name": "model1",
                                                    "id": 11
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
                    "bk_host_id": 3
                },
                "topo": [
                    {
                        "inst": {
                            "obj": "test_custom_level_2",
                            "name": "test1",
                            "id": 8
                        },
                        "children": [
                            {
                                "inst": {
                                    "obj": "test_custom_level_1",
                                    "name": "test2",
                                    "id": 9
                                },
                                "children": [
                                    {
                                        "inst": {
                                            "obj": "set",
                                            "name": "set1",
                                            "id": 5
                                        },
                                        "children": [
                                            {
                                                "inst": {
                                                    "obj": "module",
                                                    "name": "model1",
                                                    "id": 11
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
                    "bk_host_id": 4
                },
                "topo": [
                    {
                        "inst": {
                            "obj": "test_custom_level_2",
                            "name": "test3",
                            "id": 10
                        },
                        "children": [
                            {
                                "inst": {
                                    "obj": "test_custom_level_1",
                                    "name": "test4",
                                    "id": 11
                                },
                                "children": [
                                    {
                                        "inst": {
                                            "obj": "set",
                                            "name": "ttt",
                                            "id": 6
                                        },
                                        "children": [
                                            {
                                                "inst": {
                                                    "obj": "module",
                                                    "name": "sss",
                                                    "id": 12
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
                            "obj": "test_custom_level_2",
                            "name": "test1",
                            "id": 8
                        },
                        "children": [
                            {
                                "inst": {
                                    "obj": "test_custom_level_1",
                                    "name": "test2",
                                    "id": 9
                                },
                                "children": [
                                    {
                                        "inst": {
                                            "obj": "set",
                                            "name": "set1",
                                            "id": 5
                                        },
                                        "children": [
                                            {
                                                "inst": {
                                                    "obj": "module",
                                                    "name": "model1",
                                                    "id": 11
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
                    "bk_host_id": 5
                },
                "topo": [
                    {
                        "inst": {
                            "obj": "test_custom_level_2",
                            "name": "test1",
                            "id": 8
                        },
                        "children": [
                            {
                                "inst": {
                                    "obj": "test_custom_level_1",
                                    "name": "test2",
                                    "id": 9
                                },
                                "children": [
                                    {
                                        "inst": {
                                            "obj": "set",
                                            "name": "set1",
                                            "id": 5
                                        },
                                        "children": [
                                            {
                                                "inst": {
                                                    "obj": "module",
                                                    "name": "model1",
                                                    "id": 11
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
                            "obj": "test_custom_level_2",
                            "name": "test3",
                            "id": 10
                        },
                        "children": [
                            {
                                "inst": {
                                    "obj": "test_custom_level_1",
                                    "name": "test4",
                                    "id": 11
                                },
                                "children": [
                                    {
                                        "inst": {
                                            "obj": "set",
                                            "name": "ttt",
                                            "id": 6
                                        },
                                        "children": [
                                            {
                                                "inst": {
                                                    "obj": "module",
                                                    "name": "sss",
                                                    "id": 12
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
            }
        ]
    }
}
```

### Return Result Parameters Description

response

| Field      | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Success or failure of the request. true: success; false: failure |
| code       | int    | Error code. 0 represents success, >0 represents failure error |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | array  | Request result                                               |


#### data

| Field  | Type  | Description        |
| :----- | :---- | :------------------ |
| count  | int   | Number of records   |
| info   | array | Host data and topology information |

#### data.info

| Field | Type  | Description       |
| :---- | :---- | :----------------- |
| host  | dict  | Actual host data   |
| topo  | array | Host topology information |

#### data.info.host

| Field                | Type   | Description          |
| -------------------- | ------ | -------------------- |
| bk_host_name         | string | Host name             |
| bk_host_innerip      | string | Inner IP              |
| bk_host_id           | int    | Host ID               |
| bk_cloud_id          | int    | Cloud control area   |
| import_from          | string | Host import source, 3 for API import |
| bk_asset_id          | string | Fixed asset number    |
| bk_cloud_inst_id     | string | Cloud host instance ID |
| bk_cloud_vendor      | string | Cloud vendor          |
| bk_cloud_host_status | string | Cloud host status     |
| bk_comment           | string | Remark                |
| bk_cpu               | int    | CPU logical cores     |
| bk_cpu_architecture  | string | CPU architecture      |
| bk_cpu_module        | string | CPU model             |
| bk_disk              | int    | Disk capacity (GB)   |
| bk_host_outerip      | string | Host outer IP         |
| bk_host_innerip_v6   | string | Host inner IPv6       |
| bk_host_outerip_v6   | string | Host outer IPv6       |
| bk_isp_name          | string | ISP name              |
| bk_mac               | string | Host inner MAC address |
| bk_mem               | int    | Host memory capacity (MB) |
| bk_os_bit            | string | Operating system bit  |
| bk_os_name           | string | Operating system name |
| bk_os_type           | string | Operating system type |
| bk_os_version        | string | Operating system version |
| bk_outer_mac         | string | Host outer MAC address |
| bk_province_name     | string | Province name         |
| bk_service_term      | int    | Warranty period       |
| bk_sla               | string | SLA level             |
| bk_sn                | string | Device SN             |
| bk_state             | string | Current state         |
| bk_state_name        | string | Country               |
| operator             | string | Main maintainer       |
| bk_bak_operator      | string | Backup maintainer     |

**Note: The returned values here only provide explanations for the system's built-in attribute fields. Other returned values depend on user-defined attribute fields.**

#### data.info.topo

| Field     | Type   | Description       |
| :-------- | :----- | :----------------- |
| inst      | Object | Topology instance information |
| children  | array  | Subset of topology instances |

#### data.info.topo.inst

| Field | Type   | Description |
| :---- | :----- | :---------- |
| obj   | string | Model ID    |
| name  | string | Instance name |
| id    | int    | Instance ID  |
