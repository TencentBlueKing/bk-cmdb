### 功能描述

查询主机及其对应topo

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 参数                      | 类型   | 必选 | 描述                                                         |
| ------------------------- | ------ | ---- | ------------------------------------------------------------ |
| bk_biz_id                 | int    | 是   | 业务id                                                       |
| mainline_property_filters | array  | 否   | 自定义层级模型组合查询条件                                   |
| set_property_filter       | object | 否   | 集群属性组合查询条件                                         |
| module_property_filter    | object | 否   | 模块属性组合查询条件                                         |
| host_property_filter      | object | 否   | 主机属性组合查询条件                                         |
| fields                    | array  | 否   | 主机属性列表，控制返回结果的主机里有哪些字段，能够加速接口请求和减少网络流量传输 |
| page                      | object | 是   | 分页查询条件，返回的主机数据按照bk_host_id排序               |

#### mainline_property_filters

| 名称      | 类型   | 必填 | 默认值           |
| --------- | ------ | ---- | ---------------- |
| bk_obj_id | string | 是   | 自定义层级模型id |
| filter    | object | 是   | 目标的筛选条件   |

##### filter

该参数为主机属性字段过滤规则的组合，用于根据主机属性字段搜索主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。 过滤规则为四元组 `field`, `operator`, `value`

| 名称     | 类型   | 必填 | 默认值 | 说明   | Description                                                  |
| :------- | :----- | :--- | :----- | :----- | :----------------------------------------------------------- |
| field    | string | 是   | 无     | 字段名 |                                                              |
| operator | string | 是   | 无     | 操作符 | 可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否   | 无     | 操作数 | 不同的operator对应不同的value格式                            |

组装规则可参考: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### set_property_filter

该参数为集群属性字段过滤规则的组合，用于根据集群属性字段搜索集群下的主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。 过滤规则为四元组 `field`, `operator`, `value`

| 名称     | 类型   | 必填 | 默认值 | 说明   | Description                                                  |
| :------- | :----- | :--- | :----- | :----- | :----------------------------------------------------------- |
| field    | string | 是   | 无     | 字段名 |                                                              |
| operator | string | 是   | 无     | 操作符 | 可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否   | 无     | 操作数 | 不同的operator对应不同的value格式                            |

组装规则可参考: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### module_property_filter

该参数为模块属性字段过滤规则的组合，用于根据模块属性字段搜索模块下的主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。 过滤规则为四元组 `field`, `operator`, `value`

| 名称     | 类型   | 必填 | 默认值 | 说明   | Description                                                  |
| :------- | :----- | :--- | :----- | :----- | :----------------------------------------------------------- |
| field    | string | 是   | 无     | 字段名 |                                                              |
| operator | string | 是   | 无     | 操作符 | 可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否   | 无     | 操作数 | 不同的operator对应不同的value格式                            |

组装规则可参考: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### host_property_filter

该参数为主机属性字段过滤规则的组合，用于根据主机属性字段搜索主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。 过滤规则为四元组 `field`, `operator`, `value`

| 名称     | 类型   | 必填 | 默认值 | 说明   | Description                                                  |
| :------- | :----- | :--- | :----- | :----- | :----------------------------------------------------------- |
| field    | string | 是   | 无     | 字段名 |                                                              |
| operator | string | 是   | 无     | 操作符 | 可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否   | 无     | 操作数 | 不同的operator对应不同的value格式                            |

组装规则可参考: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### page

| 字段  | 类型 | 必选 | 描述                 |
| :---- | :--- | :--- | :------------------- |
| start | int  | 是   | 记录开始位置         |
| limit | int  | 是   | 每页限制条数,最大500 |

#### 请求参数示例

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

### 返回结果示例

```json
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

### 返回结果参数说明

#### response

| 字段       | 类型   | 描述                                          |
| ---------- | ------ | --------------------------------------------- |
| result     | bool   | 请求成功与否。true：请求成功；false：请求失败 |
| code       | int    | 错误编吗。0表示success，>0表示失败错误        |
| message    | string | 请求失败返回的错误信息                        |
| permission | object | 权限信息                                      |
| request_id | string | 请求链id                                      |
| data       | array  | 请求结果                                      |

#### data

| 字段  | 类型  | 描述               |
| :---- | :---- | :----------------- |
| count | int   | 记录条数           |
| info  | array | 主机数据和拓扑信息 |

#### data.info

| 字段 | 类型  | 描述         |
| :--- | :---- | :----------- |
| host | dict  | 主机实际数据 |
| topo | array | 主机拓扑信息 |

#### data.info.host

| 名称             | 类型   | 说明          |                                 |
| :--------------- | :----- | :------------ | :------------------------------ |
| bk_isp_name      | string | 所属运营商    | 0:其它；1:电信；2:联通；3:移动  |
| bk_sn            | string | 设备SN        |                                 |
| operator         | string | 主要维护人    |                                 |
| bk_outer_mac     | string | 外网MAC       |                                 |
| bk_state_name    | string | 所在国家      | CN:中国，详细值，请参考CMDB页面 |
| bk_province_name | string | 所在省份      |                                 |
| import_from      | string | 录入方式      | 1:excel;2:agent;3:api           |
| bk_sla           | string | SLA级别       | 1:L1;2:L2;3:L3                  |
| bk_service_term  | int    | 质保年限      | 1-10                            |
| bk_os_type       | string | 操作系统类型  | 1:Linux;2:Windows;3:AIX         |
| bk_os_version    | string | 操作系统版本  |                                 |
| bk_os_bit        | int    | 操作系统位数  |                                 |
| bk_mem           | string | 内存容量      |                                 |
| bk_mac           | string | 内网MAC地址   |                                 |
| bk_host_outerip  | string | 外网IP        |                                 |
| bk_host_name     | string | 主机名称      |                                 |
| bk_host_innerip  | string | 内网IP        |                                 |
| bk_host_id       | int    | 主机ID        |                                 |
| bk_disk          | int    | 磁盘容量      |                                 |
| bk_cpu_module    | string | CPU型号       |                                 |
| bk_cpu_mhz       | int    | CPU频率       |                                 |
| bk_cpu           | int    | CPU逻辑核心数 | 1-1000000                       |
| bk_comment       | string | 备注          |                                 |
| bk_cloud_id      | int    | 云区域        |                                 |
| bk_bak_operator  | string | 备份维护人    |                                 |
| bk_asset_id      | string | 固资编号      |                                 |

#### data.info.topo

| 字段     | 类型   | 描述           |
| :------- | :----- | :------------- |
| inst     | Object | 拓扑实例信息   |
| children | array  | 拓扑实例的子集 |

#### data.info.topo.inst

| 字段 | 类型   | 描述   |
| :--- | :----- | :----- |
| obj  | string | 模型id |
| name | string | 实例名 |
| id   | int    | 实例id |

 