### Functional description

find a business topology node's path which is from the top biz instance node directly to this node's
 parent instance node.(v3.9.1)

**Note**
this api has cache service, the longest cache ttl is 5 minutes. 

### Request Parameters

{{ common_args_desc }}

#### Parameters description

| Field               | Type   | Required | Description  
|-----------|------------|--------|------------|
| bk_biz_id  | int  | Yes     | business id |
| bk_nodes  | object array  | Yes    | the node's basic info to be queried, max length 1000 |


#### bk_nodes description

| Field               | Type   | Required | Description  
| ----- | ------ | ---- | --------------------- |
| bk_obj_id | string    | Yes   | the node's object name，such as biz,set,module and user custom object        |
| bk_inst_id | int64    | Yes   | the object's instance id |

### Request Parameters Example

```json
{
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

#### data field description
| Field       | Type     | Description       |
|-----------|------------|------------|
| bk_obj_id | string     | the node's object name，such as biz,set,module and user custom object        |
| bk_inst_id | int64     | the object's instance id |
| bk_paths | object array| the node's path info |

#### bk_paths field description
| Field       | Type     | Description       |
|-----------|------------|------------|
| bk_obj_id | string     | the node's object name，such as biz,set,module and user custom object        |
| bk_inst_id | int64     | the object's instance id |
| bk_inst_name | string   |the object's instance name |
