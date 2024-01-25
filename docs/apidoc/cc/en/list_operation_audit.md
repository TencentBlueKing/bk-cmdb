### Function Description

Retrieve operation audit logs based on conditions (Permission: Operation audit query permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description                          |
| --------- | ------ | -------- | ------------------------------------ |
| page      | object | Yes      | Pagination parameters                |
| condition | object | No       | Operation audit log query conditions |

#### page

| Field | Type   | Required | Description                             |
| ----- | ------ | -------- | --------------------------------------- |
| start | int    | No       | Record start position                   |
| limit | int    | Yes      | Number of records per page, maximum 200 |
| sort  | string | No       | Sorting field                           |

#### condition

| Field          | Type   | Required | Description                                                  |
| -------------- | ------ | -------- | ------------------------------------------------------------ |
| bk_biz_id      | int    | No       | Business ID                                                  |
| resource_type  | string | No       | Specific resource type of the operation                      |
| action         | array  | No       | Operation types                                              |
| operation_time | object | Yes      | Operation time                                               |
| user           | string | No       | Operator                                                     |
| resource_name  | string | No       | Resource name                                                |
| category       | string | No       | Query type                                                   |
| fuzzy_query    | bool   | No       | Whether to use fuzzy query on resource name. **Fuzzy queries are inefficient and have poor performance. This field only affects resource_name, and when using the condition method for fuzzy queries, this field will be ignored. Please choose one of the two methods to use.** |
| condition      | array  | No       | Specify query conditions, cannot be provided at the same time as user and resource_name |

##### condition.condition

| Field    | Type         | Required | Description                                                  |
| -------- | ------------ | -------- | ------------------------------------------------------------ |
| field    | string       | Yes      | Field of the object, only "user" and "resource_name" are supported |
| operator | string       | Yes      | Operator, "in" means belonging to, "not_in" means not belonging to, "contains" means contains. When field is "resource_name," you can use "contains" for fuzzy queries |
| value    | string/array | Yes      | Value corresponding to the field, array type is required for "in" and "not_in," string type is required for "contains" |

### Request Parameter Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "condition": {
        "bk_biz_id": 2,
        "resource_type": "host",
        "action": [
            "create",
            "delete"
        ],
        "operation_time": {
            "start": "2020-09-23 00:00:00",
            "end": "2020-11-01 23:59:59"
        },
        "user": "admin",
        "resource_name": "1.1.1.1",
        "category": "host",
        "fuzzy_query": false
    },
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "-operation_time"
    }
}
```
```json
{
    "condition": {
        "bk_biz_id": 2,
        "resource_type": "host",
        "action": [
            "create",
            "delete"
        ],
        "operation_time": {
            "start": "2020-09-23 00:00:00",
            "end": "2020-11-01 23:59:59"
        },
        "condition": [
            {
                "field": "user",
                "operator": "in",
                "value": ["admin"]
            },
            {
                "field": "resource_name",
                "operator": "in",
                "value": ["1.1.1.1"]
            }
        ],
        "category": "host"
    },
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "-operation_time"
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
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "count": 2,
        "info": [
            {
                "id": 7,
                "audit_type": "",
                "bk_supplier_account": "",
                "user": "admin",
                "resource_type": "host",
                "action": "delete",
                "operate_from": "",
                "operation_detail": null,
                "operation_time": "2020-10-09 21:30:51",
                "bk_biz_id": 1,
                "resource_id": 4,
                "resource_name": "2.2.2.2"
            },
            {
                "id": 2,
                "audit_type": "",
                "bk_supplier_account": "",
                "user": "admin",
                "resource_type": "host",
                "action": "delete",
                "operate_from": "",
                "operation_detail": null,
                "operation_time": "2020-10-09 17:13:55",
                "bk_biz_id": 1,
                "resource_id": 1,
                "resource_name": "1.1.1.1"
            }
        ]
    }
}
```

### Response Result Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Data returned by the request                                 |

#### data

| Field | Type  | Description                     |
| ----- | ----- | ------------------------------- |
| count | int   | Number of records               |
| info  | array | Operation audit log information |

#### info

| Field               | Type   | Description          |
| ------------------- | ------ | -------------------- |
| id                  | int    | Audit ID             |
| audit_type          | string | Operation audit type |
| bk_supplier_account | string | Supplier account     |
| user                | string | Operator             |
| resource_type       | string | Resource type        |
| action              | string | Operation type       |
| operate_from        | string | Source platform      |
| operation_detail    | object | Operation details    |
| operation_time      | string | Operation time       |
| bk_biz_id           | int    | Business ID          |
| resource_id         | int    | Resource ID          |
| resource_name       | string | Resource name        |