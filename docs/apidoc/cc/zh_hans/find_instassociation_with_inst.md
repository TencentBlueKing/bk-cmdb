### 功能描述

查询模型实例关联关系，可选择返回源模型实例与目标模型实例的详情(v3.10.11+)

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 参数      | 类型 | 必选 | 描述     |
| --------- | ---- | ---- | -------- |
| condition | map  | 是   | 查询参数 |
| page      | map  | 是   | 分页条件 |

**condition**

| 参数        | 类型  | 必选 | 描述                                      |
| :---------- | ----- | ---- | ----------------------------------------- |
| asst_filter | map   | 是   | 查询关联关系的filter                      |
| asst_fields | array | 否   | 关联关系需要返回的内容，不填返回全部      |
| src_fields  | array | 否   | 源模型需要返回的属性，不填返回全部        |
| dst_fields  | array | 否   | 目标模型需要返回的属性，不填返回全部      |
| src_detail  | bool  | 否   | 不填默认为false，不返回源模型的实例详情   |
| dst_detail  | bool  | 否   | 不填默认为false，不返回目标模型的实例详情 |

**asst_filter**

该参数为关联关系属性字段过滤规则的组合，用于根据关联关系属性搜索关联关系。组合支持AND 和 OR 两种方式，可以嵌套，最多嵌套2层。 过滤规则为四元组 `field`, `operator`, `value`

| 参数      | 类型   | 必选 | 描述                          |
| --------- | ------ | ---- | ----------------------------- |
| condition | string | 是   | 查询条件的组合方式，AND或者OR |
| rule      | array  | 是   | 包含所有查询条件的集合        |

**rule**

| 参数     | 类型   | 必选 | 描述                                                         |
| -------- | ------ | ---- | ------------------------------------------------------------ |
| field    | string | 是   | 查询条件中的字段，例如：bk_obj_id，bk_asst_obj_id，bk_inst_id |
| operator | string | 是   | 查询条件中的查询方式，equal、in、nin等                       |
| value    | string | 是   | 查询条件对应的值                                             |

组装规则可参考: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

**page**

| 参数  | 类型   | 必选 | 描述                 |
| ----- | ------ | ---- | -------------------- |
| start | int    | 否   | 记录开始位置         |
| limit | int    | 是   | 每页限制条数,最大200 |
| sort  | string | 否   | 排序字段             |

**分页对象为关联关系**

#### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "condition": {
        "asst_filter": {
            "condition": "AND",
            "rules": [
                {
                    "field": "bk_obj_id",
                    "operator": "equal",
                    "value": "bk_switch"
                },
                {
                    "field": "bk_inst_id",
                    "operator": "equal",
                    "value": 1
                },
                {
                    "field": "bk_asst_obj_id",
                    "operator": "equal",
                    "value": "host"
                }
            ]
        },
        "src_fields": [
            "bk_inst_id",
            "bk_inst_name"
        ],
        "dst_fields": [
            "bk_host_innerip"
        ],
        "src_detail": true,
        "dst_detail": true
    },
    "page": {
        "start": 0,
        "limit": 20,
        "sort": "-bk_asst_inst_id"
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
        "association": [
            {
                "id": 3,
                "bk_inst_id": 1,
                "bk_obj_id": "bk_switch",
                "bk_asst_inst_id": 3,
                "bk_asst_obj_id": "host",
                "bk_obj_asst_id": "bk_switch_connect_host",
                "bk_asst_id": "connect"
            },
            {
                "id": 2,
                "bk_inst_id": 1,
                "bk_obj_id": "bk_switch",
                "bk_asst_inst_id": 2,
                "bk_asst_obj_id": "host",
                "bk_obj_asst_id": "bk_switch_connect_host",
                "bk_asst_id": "connect"
            },
            {
                "id": 1,
                "bk_inst_id": 1,
                "bk_obj_id": "bk_switch",
                "bk_asst_inst_id": 1,
                "bk_asst_obj_id": "host",
                "bk_obj_asst_id": "bk_switch_connect_host",
                "bk_asst_id": "connect"
            }
        ],
        "src": [
            {
                "bk_inst_id": 1,
                "bk_inst_name": "s1"
            }
        ],
        "dst": [
            {
                "bk_host_innerip": "10.11.11.1"
            },
            {
                "bk_host_innerip": "10.11.11.2"
            },
            {
                "bk_host_innerip": "10.11.11.3"
            }
        ]
    }
}
```

### 返回结果参数说明

#### response

| 字段                | 类型  | 描述       |
| ------------------- | ----- | ---------- |
| result     | bool   | 请求成功与否。true：请求成功；false：请求失败 |
| code       | int    | 错误编吗。0表示success，>0表示失败错误        |
| message    | string | 请求失败返回的错误信息                        |
| permission | object | 权限信息                                      |
| request_id | string | 请求链id                                      |
| data       | object | 请求结果                                      |

#### data

| 字段        | 类型  | 描述                                     |
| ----------- | ----- | ---------------------------------------- |
| association | array | 查询到的关联关系详情，按分页排序参数排序 |
| src         | array | 源模型实例的详情                         |
| dst         | array | 目标模型实例的详情                       |

##### association

| 名称            | 类型   | 说明                     |
| --------------- | ------ | ------------------------ |
| id              | int64  | 关联id                   |
| bk_inst_id      | int64  | 源模型实例id             |
| bk_obj_id       | string | 关联关系源模型id         |
| bk_asst_inst_id | int64  | 关联关系目标模型id       |
| bk_asst_obj_id  | string | 目标模型实例id           |
| bk_obj_asst_id  | string | 自动生成的模型关联关系id |
| bk_asst_id      | string | 关系名称                 |

##### src

| 名称         | 类型   | 说明   |
| ------------ | ------ | ------ |
| bk_inst_name | string | 实例名 |
| bk_inst_id   | int    | 实例id |

##### dst

| 名称             | 类型   | 说明       |
| ---------------- | ------ | ---------- |
| bk_host_inner_ip | string | 主机内网ip |

