### Functional description

Query the model instance Association, and optionally return the details of the source model instance and the target model instance (v3.10.11+)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Parameter      | Type| Required| Description     |
| --------- | ---- | ---- | -------- |
| condition | map  |yes   | Query parameter|
| page      |  map  |yes   | Paging condition|

**condition**

| Parameter        | Type| Required| Description                                      |
| :---------- | ----- | ---- | ----------------------------------------- |
| asst_filter | map   | yes | Filter for querying Association relationship                      |
| asst_fields | array |no   | Content to be returned for Association relationship. All are returned without filling in      |
| src_fields  | array |no   | The attributes to be returned by the source model. All are returned without filling in        |
| dst_fields  | array |no   | Attributes to be returned by the target model. All are returned without filling in      |
| src_detail  | bool  |no   | The default value is false, and the instance details of the source model are not returned   |
| dst_detail  | bool  |no   | The default value is false, and the instance details of the target model are not returned|

**asst_filter**

This parameter is a combination of filtering rules for Association attribute fields, and is used to search Association according to Association attribute. The combination supports AND and OR, and can be nested, with a maximum of 2 layers. The filtering rules are quaternions`field`,`operator`,`value`

| Parameter      | Type   | Required| Description                          |
| --------- | ------ | ---- | ----------------------------- |
| condition | string |yes   | Combination of query criteria, AND or OR|
| rule      |  array  |yes   | Collection containing all query criteria        |

**rule**

| Parameter     | Type   | Required| Description                                                         |
| -------- | ------ | ---- | ------------------------------------------------------------ |
| field    |  string |yes   | Fields in query criteria, for example: bk_obj_id，bk_asst_obj_id，bk_inst_id|
| operator | string |yes   | Query method in query criteria, equal, in, nin, etc.                       |
| value    |  string |yes   | Value corresponding to query criteria                                             |

For assembly rules, refer to: https: //github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md

**page**

| Parameter| Type   | Required| Description                 |
| ----- | ------ | ---- | -------------------- |
| start | int    | no   | Record start position         |
| limit | int    | yes | Limit bars per page, Max. 200|
| sort  | string |no   | Sort field             |

**Paging object is associated**

#### Request Parameters Example

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

### Return Result Example

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

### Return Result Parameters Description

#### response

| Field                | Type| Description       |
| ------------------- | ----- | ---------- |
| result     |  bool   | Whether the request was successful or not. True: request succeeded;false: Request failed|
| code       |  int    | Wrong. 0 indicates success,>0 indicates failure error        |
| message    |  string |Error message returned by request failure                        |
| permission | object |Permission information                                      |
| request_id | string |Request chain id                                      |
| data       |  object |Request result                                      |

#### data

| Field        | Type| Description                                     |
| ----------- | ----- | ---------------------------------------- |
| association | array |The queried Association relationship details are sorted by the paging sorting parameter|
| src         |  array |Details of the source model instance                         |
| dst         |  array |Details of the target model instance                       |

##### association

| Name            | Type   | Description                     |
| --------------- | ------ | ------------------------ |
| id              |  int64  |Association id                   |
| bk_inst_id      |  int64  |Source model instance id             |
| bk_obj_id       |  string |Association relationship source model id         |
| bk_asst_inst_id | int64  |Association relation target model id       |
| bk_asst_obj_id  | string |Target model instance id           |
| bk_obj_asst_id  | string |Auto-generated model association id|
| bk_asst_id      |  string |Relationship name                 |

##### src

| Name         | Type   | Description   |
| ------------ | ------ | ------ |
| bk_inst_name | string |Instance name|
| bk_inst_id   |  int    | Instance id|

##### dst

| Name             | Type   | Description       |
| ---------------- | ------ | ---------- |
| bk_host_inner_ip | string |Host intranet ip|

