### 描述

查询主机及其对应topo

### 输入参数

| 参数名称                      | 参数类型   | 必选 | 描述                                       |
|---------------------------|--------|----|------------------------------------------|
| bk_biz_id                 | int    | 是  | 业务id                                     |
| mainline_property_filters | array  | 否  | 自定义层级模型组合查询条件                            |
| set_property_filter       | object | 否  | 集群属性组合查询条件                               |
| module_property_filter    | object | 否  | 模块属性组合查询条件                               |
| host_property_filter      | object | 否  | 主机属性组合查询条件                               |
| fields                    | array  | 否  | 主机属性列表，控制返回结果的主机里有哪些字段，能够加速接口请求和减少网络流量传输 |
| page                      | object | 是  | 分页查询条件，返回的主机数据按照bk_host_id排序             |

#### mainline_property_filters

| 参数名称      | 参数类型   | 必选 | 描述        |
|-----------|--------|----|-----------|
| bk_obj_id | string | 是  | 自定义层级模型id |
| filter    | object | 是  | 目标的筛选条件   |

##### filter

该参数为主机属性字段过滤规则的组合，用于根据主机属性字段搜索主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。
过滤规则为四元组 `field`, `operator`, `value`

| 参数名称     | 参数类型   | 必选 | 描述                                                                                                |
|----------|--------|----|---------------------------------------------------------------------------------------------------|
| field    | string | 是  | 字段名                                                                                               |
| operator | string | 是  | 操作符,可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否  | 操作数,不同的operator对应不同的value格式                                                                       |

组装规则可参考: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### set_property_filter

该参数为集群属性字段过滤规则的组合，用于根据集群属性字段搜索集群下的主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。
过滤规则为四元组 `field`, `operator`, `value`

| 参数名称     | 参数类型   | 必选 | 描述                                                                                                |
|----------|--------|----|---------------------------------------------------------------------------------------------------|
| field    | string | 是  | 字段名                                                                                               |
| operator | string | 是  | 操作符,可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否  | 操作数,不同的operator对应不同的value格式                                                                       |

组装规则可参考: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### module_property_filter

该参数为模块属性字段过滤规则的组合，用于根据模块属性字段搜索模块下的主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。
过滤规则为四元组 `field`, `operator`, `value`

| 参数名称     | 参数类型   | 必选 | 描述                                                                                                |
|----------|--------|----|---------------------------------------------------------------------------------------------------|
| field    | string | 是  | 字段名                                                                                               |
| operator | string | 是  | 操作符,可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否  | 操作数,不同的operator对应不同的value格式                                                                       |

组装规则可参考: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### host_property_filter

该参数为主机属性字段过滤规则的组合，用于根据主机属性字段搜索主机。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。
过滤规则为四元组 `field`, `operator`, `value`

| 参数名称     | 参数类型   | 必选 | 描述                                                                                                |
|----------|--------|----|---------------------------------------------------------------------------------------------------|
| field    | string | 是  | 字段名                                                                                               |
| operator | string | 是  | 操作符,可选值 equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between |
| value    | -      | 否  | 操作数,不同的operator对应不同的value格式                                                                       |

组装规则可参考: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

#### page

| 参数名称  | 参数类型 | 必选 | 描述           |
|-------|------|----|--------------|
| start | int  | 是  | 记录开始位置       |
| limit | int  | 是  | 每页限制条数,最大500 |

### 调用示例

```json
{
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
    "fields": [
        "bk_host_id",
        "bk_host_innerip"
    ],
    "page": {
        "start": 0,
        "limit": 500
    }
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
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

### 响应参数说明

| 参数名称       | 参数类型   | 描述                          |
|------------|--------|-----------------------------|
| result     | bool   | 请求成功与否。true：请求成功；false：请求失败 |
| code       | int    | 错误编吗。0表示success，>0表示失败错误    |
| message    | string | 请求失败返回的错误信息                 |
| permission | object | 权限信息                        |
| data       | array  | 请求结果                        |

#### data

| 参数名称  | 参数类型  | 描述        |
|-------|-------|-----------|
| count | int   | 记录条数      |
| info  | array | 主机数据和拓扑信息 |

#### data.info

| 参数名称 | 参数类型  | 描述     |
|------|-------|--------|
| host | dict  | 主机实际数据 |
| topo | array | 主机拓扑信息 |

#### data.info.host

| 参数名称                 | 参数类型   | 描述                |
|----------------------|--------|-------------------|
| bk_host_name         | string | 主机名               |
| bk_host_innerip      | string | 内网IP              |
| bk_host_id           | int    | 主机ID              |
| bk_cloud_id          | int    | 管控区域              |
| import_from          | string | 主机导入来源,以api方式导入为3 |
| bk_asset_id          | string | 固资编号              |
| bk_cloud_inst_id     | string | 云主机实例ID           |
| bk_cloud_vendor      | string | 云厂商               |
| bk_cloud_host_status | string | 云主机状态             |
| bk_comment           | string | 备注                |
| bk_cpu               | int    | CPU逻辑核心数          |
| bk_cpu_architecture  | string | CPU架构             |
| bk_cpu_module        | string | CPU型号             |
| bk_disk              | int    | 磁盘容量（GB）          |
| bk_host_outerip      | string | 主机外网IP            |
| bk_host_innerip_v6   | string | 主机内网IPv6          |
| bk_host_outerip_v6   | string | 主机外网IPv6          |
| bk_isp_name          | string | 所属运营商             |
| bk_mac               | string | 主机内网MAC地址         |
| bk_mem               | int    | 主机名内存容量（MB）       |
| bk_os_bit            | string | 操作系统位数            |
| bk_os_name           | string | 操作系统名称            |
| bk_os_type           | string | 操作系统类型            |
| bk_os_version        | string | 操作系统版本            |
| bk_outer_mac         | string | 主机外网MAC地址         |
| bk_province_name     | string | 所在省份              |
| bk_service_term      | int    | 质保年限              |
| bk_sla               | string | SLA级别             |
| bk_sn                | string | 设备SN              |
| bk_state             | string | 当前状态              |
| bk_state_name        | string | 所在国家              |
| operator             | string | 主要维护人             |
| bk_bak_operator      | string | 备份维护人             |

**注意：此处的返回值仅对系统内置的属性字段做了说明，其余返回值取决于用户自己定义的属性字段**

#### data.info.topo

| 参数名称     | 参数类型   | 描述      |
|----------|--------|---------|
| inst     | Object | 拓扑实例信息  |
| children | array  | 拓扑实例的子集 |

#### data.info.topo.inst

| 参数名称 | 参数类型   | 描述   |
|------|--------|------|
| obj  | string | 模型id |
| name | string | 实例名  |
| id   | int    | 实例id |
 
 
