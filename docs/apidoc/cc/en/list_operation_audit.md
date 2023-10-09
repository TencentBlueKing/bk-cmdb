### Functional description

 Get action Audit log based on condition

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type      | Required   | Description                       |
|---------------------|------------|--------|-----------------------------|
| page                |  object     | yes  | Paging parameter                    |
| condition           |  object     | no     | Operation audit log query criteria                   |

#### page

| Field      | Type      | Required   | Description                |
|-----------|------------|--------|----------------------|
| start     |   int       | no     | Record start position         |
| limit     |   int       | yes  | Limit bars per page, Max. 200|
| sort      |   string    | no     | Sort field             |

#### condition

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id     | int      | no      | Business ID                                    |
| resource_type  |string      | no      | Specific resource type of operation|
| action     |    array  |no    | Operation type|
|   operation_time   |    object  |yes    | Operating time|
|   user   |    string  |no    | Operator|
|    resource_name  |    string  |  no   | Resource name |
|    category  |    string  |  no  | Type of query |
| fuzzy_query    |  bool         | no   | Whether to use fuzzy query to query the resource name is **inefficient and poor in performance**. This field only affects resource_name. This field will be ignored when using condition to perform fuzzy query. Please choose one of the two. |
| condition | array |no| Specify query criteria, which can not be provided at the same time as user and resource_name|

##### condition.condition

| Field     | Type         | Required| Description                                                         |
| -------- | ------------ | ---- | ------------------------------------------------------------ |
| field    |  string       | yes | Object,"user" only,"resource_name"                      |
| operator | string       | yes | Operator: in is belongs to, not_in is does not belong to, contains is contains, and field is resource_name. Contains can be used for fuzzy query|
| value    |  string/array |yes   | The value corresponding to the field, in and not_in require array type, and contexts require String type|

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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
            "operator":"in",
            "value":["admin"]
          },
          {
            "field":"resource_name",
            "operator":"in",
            "value":["1.1.1.1"]
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
    "code":0,
    "message":"success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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

| Field      | Type      | Description         |
|-----------|-----------|--------------|
| count     |  int       | Number of records     |
| info      |  array     | Record information of operation audit|

#### info
| Field      | Type      | Description         |
|-----------|-----------|--------------|
|    id |      int  | Audit ID  |
|   audit_type  |     string   |   Operational audit type   |
|   bk_supplier_account  |    string    | Developer account number     |
|   user  |      string  |    Operator|
|   resource_type  |    string    |   Resource type   |
|  action   |    string    |    Operation type|
|    operate_from |    string    | Source platform          |
|  operation_detail   |     object     | Operational details    |
| operation_time    |     string   |    Operating time|
|  bk_biz_id   |       int | Business ID |
| resource_id    |     int   |    Resource id|
|   resource_name  |     string   | Resource Name    |