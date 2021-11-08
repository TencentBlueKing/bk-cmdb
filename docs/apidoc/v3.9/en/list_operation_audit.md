### Functional description

 Obtain operation audit logs based on conditions

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field | Type | Required | Description |
|---------------------|------------|--------|-----------------------------|
| page | object | Yes | paging parameters |
| condition | object | No | Operation audit log query condition |

#### page

| Field | Type | Required | Description |
|-----------|------------|--------|----------------------|
| start | int | No | record start position |
| limit | int | Yes | limit the number of entries per page, maximum 200 |
| sort | string | No | Sort field |

#### condition

| Field | Type | Required | Description |
|-----------|------------|--------|------------|
| bk_biz_id | int | No | business id |
| resource_type |string | No | The specific resource type of the operation |
| action | array | No | Action type |
| operation_time | object | Yes | Operation time |
| user | string | No | Operator |
| resource_name | string | No | Resource name |
| category | string | No | type of query |
| fuzzy_query    | bool         | No       | use fuzzy query or not when searching with resource name, **fuzzy query has low efficiency and poor performance** |
| condition | array | No | condition query's condition, with user and resource_name cannot be provided at the same time |

##### condition.condition

| Field    | Type         | Required | Description                                                  |
| -------- | ------------ | -------- | ------------------------------------------------------------ |
| field    | string       | Yes      | fields of query, only can be "user" and "resource_name"      |
| operator | string       | Yes      | operator，$in means "belong to"，$nin means "does not belong  to"，$regex means "contain"，$regex only can be used by resource_name |
| value    | string/array | Yes      | value of query，$in and $nin need array type，$regex need string type |

### Request Parameters Example

```json
{
    "condition":{
        "bk_biz_id":2,
        "resource_type":"host",
        "action":[
            "create",
            "delete"
        ],
        "operation_time":{
            "start":"2020-09-23 00:00:00",
            "end":"2020-11-01 23:59:59"
        },
        "user":"admin",
        "resource_name":"1.1.1.1",
        "category":"host",
        "fuzzy_query": false
    },
    "page":{
        "start":0,
        "limit":10,
        "sort":"-operation_time"
    }
}
```

```json
{
    "condition":{
        "bk_biz_id":2,
        "resource_type":"host",
        "action":[
            "create",
            "delete"
        ],
        "operation_time":{
            "start":"2020-09-23 00:00:00",
            "end":"2020-11-01 23:59:59"
        },
      	"condition":[
          {
            "field":"user",
            "operatior":"$in",
            "value":"admin"
          },
          {
            "field":"resource_name",
            "operatior":"$in",
            "value":"1.1.1.1"
          }
        ],
        "category":"host"
    },
    "page":{
        "start":0,
        "limit":10,
        "sort":"-operation_time"
    }
}
```

### Return Result Example

```json
{
    "result":true,
    "bk_error_code":0,
    "bk_error_msg":"success",
    "permission":null,
    "data":{
        "count":2,
        "info":[
            {
                "id":7,
                "audit_type":"",
                "bk_supplier_account":"",
                "user":"admin",
                "resource_type":"host",
                "action":"delete",
                "operate_from":"",
                "operation_detail":null,
                "operation_time":"2020-10-09 21:30:51",
                "bk_biz_id":1,
                "resource_id":4,
                "resource_name":"2.2.2.2"
            },
            {
                "id":2,
                "audit_type":"",
                "bk_supplier_account":"",
                "user":"admin",
                "resource_type":"host",
                "action":"delete",
                "operate_from":"",
                "operation_detail":null,
                "operation_time":"2020-10-09 17:13:55",
                "bk_biz_id":1,
                "resource_id":1,
                "resource_name":"1.1.1.1"
            }
        ]
    }
}
```

### Return Result Parameters Description

#### data

| Field | Type | Description |
|-----------|-----------|--------------|
| count | int | Number of records |
| info | array | Operation audit record information |
