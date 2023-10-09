### Functional description

Query the host details and the topology information it belongs to according to the host condition information

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| page       |   dict    | yes  | Query criteria|
| host_property_filter|  object| no | Host attribute combination query criteria|
| fields    |   array   | yes  | Host attribute list, which controls the fields in the host that returns the result. Please fill in as required|

#### host_property_filter
This parameter is a combination of filtering rules for the host attribute field and is used to search for hosts based on the host attribute field. The combination supports AND and OR, and can be nested, with a maximum of 2 layers.
The filtering rules are quaternions`field`,`operator`,`value`

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| condition       |   string    | no     | Combined query criteria|
| rules      |   array    | no     | Rule|


#### rules

| Name| Type| Required| Default value|  Description|
| ---  | ---  | --- |---  | ---|
| field| string| yes | None| Field name| Field name|
| operator| string| yes | None| Operator| Optional values equal,not_equal,in, not_in, less,less_or_equal,greater,greater_or_equal,between,not_between, contexts, exists,not_exists; Where contexts is regular matching and case insensitive, exists is the condition for filtering the existence of a field, and not_exists is the condition for filtering the non-existence of a field|
| value| - |no| None| Operand| Different values correspond to different value formats|

Assembly rules can be found at https: <//github.com/Lucent/bk CMDB/blob/master/src/common/QueryBuilder/README.md>



#### page

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| start    |   int    | yes  | Record start position|
| limit    |   int    | yes  | Limit bars per page, maximum 500|
| sort     |   string |no     | Sort field|



### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
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
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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
                            "name": "China",
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
                            "name": "Country",
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
                            "name": "Idle pool",
                            "id": 2
                        },
                        "children": [
                            {
                                "inst": {
                                    "obj": "module",
                                    "name": "Faulty machine",
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
#### response

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request was successful or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | object |Data returned by request|

#### data

| Field      | Type      | Description      |
|-----------|-----------|-----------|
| count     |  int       | Number of records|
| info      |  array     | Host data and topology information|

#### data.info
| Field      | Type      | Description      |
|-----------|-----------|-----------|
| host      |  dict      | Host actual data|
| topo      |  array     | Host topology information|

#### data.info.host
| Name| Type| Description|
|---|---|---|
| bk_isp_name|  string |Operator| 0: other;1: telecommunications;2: Unicom;3: mobile|
| bk_sn | string |Equipment SN||
| operator | string |Main maintainer||
| bk_outer_mac | string |Extranet MAC||
| bk_state_name | string |Country| CN: China please refer to CMDB page for detailed values|
| bk_province_name | string |Province||
| import_from | string |Entry method| 1:excel;2:agent; 3:api|
| bk_sla | string |SLA level| 1:L1;2:L2;3: L3|
| bk_service_term | int |Warranty period|  1-10 |
| bk_os_type | string |Operating system type| 1:Linux;2:Windows; 3:AIX|
| bk_os_version | string |Operating system version||
| bk_os_bit | int |Operating system bits||
| bk_mem | string |Memory capacity||
| bk_mac | string |Intranet MAC address||
| bk_host_outerip | string |Extranet IP||
| bk_host_name | string |Host name||
| bk_host_innerip | string |Intranet IP||
| bk_host_id | int |Host ID||
| bk_disk | int |Disk capacity||
| bk_cpu_module | string |CPU model||
| bk_cpu_mhz | int |CPU frequency||
| bk_cpu | int |Number of CPU logical cores|  1-1000000 |
| bk_comment | string |Remarks||
| bk_cloud_id | int |Cloud area||
| bk_bak_operator | string |Backup maintainer||
| bk_asset_id | string |Fixed capital no||

#### data.info.topo
| Field      | Type      | Description      |
|-----------|-----------|-----------|
| inst | object       | Details describing the node instance|
| inst.obj | string       | Model type of this node, such as set/module and user-defined hierarchical model type|
| inst.name | string  |The instance name of the node|
| inst.id    |  int     | The instance ID of this node|
| children | object array       | Details of child nodes describing the current instance node, possibly multiple|
| children.inst | object       | Instance detail information of the child node|
| children.children | string  |Detail information describing the child nodes of the node|


