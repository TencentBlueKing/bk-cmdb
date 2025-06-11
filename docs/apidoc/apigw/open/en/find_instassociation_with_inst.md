### Description

Query the instance association relationship of the model, and optionally return details of the source model instance and
target model instance. (Version: v3.10.11+, Permission: Model instance query permission)

### Parameters

| Name      | Type   | Required | Description                    |
|-----------|--------|----------|--------------------------------|
| bk_obj_id | string | Yes      | Unique identifier of the model |
| condition | map    | Yes      | Query parameters               |
| page      | map    | Yes      | Pagination conditions          |

**condition**

| Name        | Type  | Required | Description                                                                        |
|-------------|-------|----------|------------------------------------------------------------------------------------|
| asst_filter | map   | Yes      | Filter for querying association relationships                                      |
| asst_fields | array | No       | Content to be returned for association relationships, returns all if not specified |
| src_fields  | array | No       | Properties to be returned for the source model, returns all if not specified       |
| dst_fields  | array | No       | Properties to be returned for the target model, returns all if not specified       |
| src_detail  | bool  | No       | Defaults to false, does not return details of the source model instance            |
| dst_detail  | bool  | No       | Defaults to false, does not return details of the target model instance            |

**asst_filter**

This parameter is a combination of filter rules for association relationship attribute fields, used to search
association relationships based on association relationship attributes. The combination supports both AND and OR, and
can be nested, with a maximum nesting depth of 2. The filter rule is a quadruple `field`, `operator`, `value`.

| Name      | Type   | Required | Description                                       |
|-----------|--------|----------|---------------------------------------------------|
| condition | string | Yes      | Combination method of query conditions, AND or OR |
| rules     | array  | Yes      | Collection containing all query conditions        |

**rules**

| Name     | Type   | Required | Description                                                                 |
|----------|--------|----------|-----------------------------------------------------------------------------|
| field    | string | Yes      | Field in the query condition, such as bk_obj_id, bk_asst_obj_id, bk_inst_id |
| operator | string | Yes      | Query method in the query condition, such as equal, in, nin, etc.           |
| value    | string | Yes      | Value corresponding to the query condition                                  |

Assembly rules can refer
to: [bk-cmdb Query Builder](https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md)

**page**

| Name  | Type   | Required | Description                        |
|-------|--------|----------|------------------------------------|
| start | int    | No       | Record start position              |
| limit | int    | Yes      | Record limit per page, maximum 200 |
| sort  | string | No       | Sorting field                      |

**Pagination object is for association relationships**

### Request Example

```json
{
    "bk_obj_id":"bk_switch",
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

### Response Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
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

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### data

| Name        | Type  | Description                                                                                             |
|-------------|-------|---------------------------------------------------------------------------------------------------------|
| association | array | Details of the queried association relationships, sorted according to the pagination sorting parameters |
| src         | array | Details of the source model instance                                                                    |
| dst         | array | Details of the target model instance                                                                    |

##### association

| Name            | Type   | Description                                               |
|-----------------|--------|-----------------------------------------------------------|
| id              | int64  | Association ID                                            |
| bk_inst_id      | int64  | Source model instance ID                                  |
| bk_obj_id       | string | Source model ID of the association relationship           |
| bk_asst_inst_id | int64  | Target model instance ID                                  |
| bk_asst_obj_id  | string | Target model ID of the association relationship           |
| bk_obj_asst_id  | string | Automatically generated model association relationship ID |
| bk_asst_id      | string | Relationship name                                         |

##### src

| Name         | Type   | Description   |
|--------------|--------|---------------|
| bk_inst_name | string | Instance name |
| bk_inst_id   | int    | Instance ID   |

##### dst

| Name             | Type   | Description   |
|------------------|--------|---------------|
| bk_host_inner_ip | string | Host inner IP |
