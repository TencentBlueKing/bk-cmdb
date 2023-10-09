### Functional description

query host and its corresponding Topo


### Request Parameters

{{ common_args_desc }}


#### General Parameters

| Field                     | Type   | Required | Description                                                  |
| ------------------------- | ------ | -------- | ------------------------------------------------------------ |
| bk_biz_id                 | int    | Yes      | business id                                                  |
| mainline_property_filters | array  | No       | custom level query filter                                    |
| set_property_filter       | object | No       | set query filter                                             |
| module_property_filter    | object | No       | module query filter                                          |
| host_property_filter      | object | No       | host query filter                                            |
| fields                    | array  | No       | host property list, which controls which fields are in the host that returns the results |
| page                      | object | Yes      | paging query condition, returning host data according to bk_ host_ id sorting |

#### mainline_property_filters

| Field     | Type   | Required | Description                    |
| --------- | ------ | -------- | ------------------------------ |
| bk_obj_id | string | Yes      | bk_obj_id of the filter target |
| filter    | object | Yes      | target object's query filter   |

##### filter

- This parameter is a combination of business attribute field filtering rules, which is used to search for business based on business attribute fields. Combinations only support AND operations and can be nested up to 2 levels.

| Field     | Type   | Required | Description                                     |
| --------- | ------ | -------- | ----------------------------------------------- |
| condition | string | Yes      | rule operator                                   |
| rules     | array  | Yes      | scope condition rules for the selected business |


#### rules

| Field    | Type   | Required | Description                                                  |
| -------- | ------ | -------- | ------------------------------------------------------------ |
| field    | string | No       | field name                                                   |
| operator | string | No       | available values: equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | No       | values's format depend on operator                           |

Assembly rules are available for reference:https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### set_property_filter

- This parameter is a combination of business attribute field filtering rules, which is used to search for business based on business attribute fields. Combinations only support AND operations and can be nested up to 2 levels.

| Field     | Type   | Required | Description                                     |
| --------- | ------ | -------- | ----------------------------------------------- |
| condition | string | Yes      | rule operator                                   |
| rules     | array  | Yes      | scope condition rules for the selected business |


#### rules

| Field    | Type   | Required | Description                                                  |
| -------- | ------ | -------- | ------------------------------------------------------------ |
| field    | string | No       | field name                                                   |
| operator | string | No       | available values: equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | No       | values's format depend on operator                           |

Assembly rules are available for reference:https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### module_property_filter

- This parameter is a combination of business attribute field filtering rules, which is used to search for business based on business attribute fields. Combinations only support AND operations and can be nested up to 2 levels.

| Field     | Type   | Required | Description                                     |
| --------- | ------ | -------- | ----------------------------------------------- |
| condition | string | Yes      | rule operator                                   |
| rules     | array  | Yes      | scope condition rules for the selected business |


#### rules

| Field    | Type   | Required | Description                                                  |
| -------- | ------ | -------- | ------------------------------------------------------------ |
| field    | string | No       | field name                                                   |
| operator | string | No       | available values: equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | No       | values's format depend on operator                           |

Assembly rules are available for reference:https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### host_property_filter

- This parameter is a combination of business attribute field filtering rules, which is used to search for business based on business attribute fields. Combinations only support AND operations and can be nested up to 2 levels.

| Field     | Type   | Required | Description                                     |
| --------- | ------ | -------- | ----------------------------------------------- |
| condition | string | Yes      | rule operator                                   |
| rules     | array  | Yes      | scope condition rules for the selected business |


#### rules

| Field    | Type   | Required | Description                                                  |
| -------- | ------ | -------- | ------------------------------------------------------------ |
| field    | string | No       | field name                                                   |
| operator | string | No       | available values: equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | No       | values's format depend on operator                           |

Assembly rules are available for reference:https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### page

| Field | Type | Required | Description                                     |
| :---- | :--- | :------- | :---------------------------------------------- |
| start | int  | No       | Record start position                           |
| limit | int  | Yes      | Limit the number of entries per page, up to 500 |

### Request Parameters Example

```json
{
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

| Field      | Type   | Description                                   |
| ---------- | ------ | --------------------------------------------- |
| result     | bool   | 请求成功与否。true：请求成功；false：请求失败 |
| code       | int    | 错误编吗。0表示success，>0表示失败错误        |
| message    | string | 请求失败返回的错误信息                        |
| permission | object | 权限信息                                      |
| request_id | string | 请求链id                                      |
| data       | array  | 请求结果                                      |

#### data

| Field | Type  | Description        |
| :---- | :---- | :----------------- |
| count | int   | 记录条数           |
| info  | array | 主机数据和拓扑信息 |

#### data.info

| Field | Type  | Description  |
| :---- | :---- | :----------- |
| host  | dict  | 主机实际数据 |
| topo  | array | 主机拓扑信息 |

#### data.info.host

| Field            | Type   | Description                 |                                        |
| :--------------- | :----- | :-------------------------- | :------------------------------------- |
| bk_isp_name      | string | Internet Service Provider   | 0:other；1:telecom；2:unicom；3:mobile |
| bk_sn            | string | Device SN                   |                                        |
| operator         | string | Main maintainer             |                                        |
| bk_outer_mac     | string | Internet mac                |                                        |
| bk_state_name    | string | Country                     | CN:China                               |
| bk_province_name | string | Province                    |                                        |
| import_from      | string | Entry method                | 1:excel;2:agent;3:api                  |
| bk_sla           | string | SLA level                   | 1:L1;2:L2;3:L3                         |
| bk_service_term  | int    | Warranty period             | 1-10                                   |
| bk_os_type       | string | Operating system type       | 1:Linux;2:Windows;3:AIX                |
| bk_os_version    | string | Operating system version    |                                        |
| bk_os_bit        | int    | Operating system bits       |                                        |
| bk_mem           | string | Memory capacity             |                                        |
| bk_mac           | string | Intranet MAC address        |                                        |
| bk_host_outerip  | string | Internet IP                 |                                        |
| bk_host_name     | string | Host name                   |                                        |
| bk_host_innerip  | string | Intranet IP                 |                                        |
| bk_host_id       | int    | Host ID                     |                                        |
| bk_disk          | int    | Disk capacity               |                                        |
| bk_cpu_module    | string | CPU model                   |                                        |
| bk_cpu_mhz       | int    | CPU frequency               |                                        |
| bk_cpu           | int    | Number of CPU logical cores | 1-1000000                              |
| bk_comment       | string | remarks                     |                                        |
| bk_cloud_id      | int    | Cloud region                |                                        |
| bk_bak_operator  | string | Backup maintainer           |                                        |
| bk_asset_id      | string | Fixed assets No             |                                        |

#### data.info.topo

| Field    | Type   | Description                   |
| :------- | :----- | :---------------------------- |
| inst     | object | Topology instance information |
| children | array  | Subset of topology instances  |

#### data.info.topo.inst

| Field | Type   | Description   |
| :---- | :----- | :------------ |
| obj   | string | Object ID     |
| name  | string | Instance name |
| id    | int    | Instance ID   |

