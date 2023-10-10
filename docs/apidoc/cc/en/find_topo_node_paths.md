### Functional description

This interface is used to query the path information from the parent level of a node to the service vertex according to a node instance (including a custom node level instance) in the service topology level. (v3.9.1)

**Attention.**
The interface has a cache, and the maximum time for cache updates is 5 minutes.

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id  | int  |yes     | Business ID |
| bk_nodes  | array  |yes     | List of service topology instance node information to query. The maximum number of queries is 1000|


#### bk_nodes Field Description

| Field| Type   | Required| Description                  |
| ----- | ------ | ---- | --------------------- |
| bk_obj_id | string    | yes | Business topology node model name, such as biz,set,module and user-defined hierarchy model name       |
| bk_inst_id | int    | yes | The instance ID of the service topology node|

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 3,
    "bk_nodes": [
        {
            "bk_obj_id": "set",
            "bk_inst_id": 11
        },
        {
            "bk_obj_id": "module",
            "bk_inst_id": 60
        },
        {
            "bk_obj_id": "province",
            "bk_inst_id": 3
        }
    ]
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
    "data": [
        {
            "bk_obj_id": "set",
            "bk_inst_id": 11,
            "bk_inst_name": "gz",
            "bk_paths": [
                [
                    {
                        "bk_obj_id": "biz",
                        "bk_inst_id": 3,
                        "bk_inst_name": "demo"
                    },
                    {
                        "bk_obj_id": "province",
                        "bk_inst_id": 3,
                        "bk_inst_name": "sz"
                    }
                ]
            ]
        },
        {
            "bk_obj_id": "module",
            "bk_inst_id": 60,
            "bk_inst_name": "m2",
            "bk_paths": [
                [
                    {
                        "bk_obj_id": "biz",
                        "bk_inst_id": 3,
                        "bk_inst_name": "demo"
                    },
                    {
                        "bk_obj_id": "province",
                        "bk_inst_id": 3,
                        "bk_inst_name": "sz"
                    },
                    {
                        "bk_obj_id": "set",
                        "bk_inst_id": 12,
                        "bk_inst_name": "set1"
                    }
                ]
            ]
        },
        {
            "bk_obj_id": "province",
            "bk_inst_id": 3,
            "bk_inst_name": "sz",
            "bk_paths": [
                [
                    {
                        "bk_obj_id": "biz",
                        "bk_inst_id": 3,
                        "bk_inst_name": "demo"
                    }
                ]
            ]
        }
    ]
}
```

### Return Result Parameters Description
#### response
| Name    | Type   | Description                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                     |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                             |


#### Data description
| Field      | Type      | Description      |
|-----------|------------|------------|
| bk_obj_id | string   | Business topology node model name, such as biz,set,module and user-defined hierarchy model name       |
| bk_inst_id | int   | The instance ID of the service topology node|
| bk_paths | array| Hierarchy information of the node, i.e. Hierarchy information from the service to the parent node of the node|

#### bk_paths description
| Field      | Type      | Description      |
|-----------|------------|------------|
| bk_obj_id | string   | Node type       |
| bk_inst_id | int   | Node instance ID|
| bk_inst_name | string   | Node instance name|
