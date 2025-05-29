### Description

This interface is used to query the path information from a node instance in the business topology hierarchy to the
business vertex based on a certain business topology level (including custom topology level). (v3.9.1)

**Note** This interface has a cache, and the longest cache update time is 5 minutes.

### Parameters

| Name      | Type  | Required | Description                                                                                              |
|-----------|-------|----------|----------------------------------------------------------------------------------------------------------|
| bk_biz_id | int   | Yes      | Business ID                                                                                              |
| bk_nodes  | array | Yes      | List of business topology instance node information to be queried, with a maximum query quantity of 1000 |

#### Explanation of bk_nodes Fields

| Name       | Type   | Required | Description                                                                                     |
|------------|--------|----------|-------------------------------------------------------------------------------------------------|
| bk_obj_id  | string | Yes      | Business topology node model name, such as biz, set, module, and the model name of custom level |
| bk_inst_id | int    | Yes      | The instance ID of the business topology node                                                   |

### Request Example

```json
{
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

### Response Example

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

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error    |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Data returned by the request                                       |

#### Explanation of data

| Name       | Type   | Description                                                                                                                      |
|------------|--------|----------------------------------------------------------------------------------------------------------------------------------|
| bk_obj_id  | string | Business topology node model name, such as biz, set, module, and the model name of custom level                                  |
| bk_inst_id | int    | The instance ID of the business topology node                                                                                    |
| bk_paths   | array  | The hierarchical information of the node, that is, the hierarchical information from the business to the parent node of the node |

#### Explanation of bk_paths

| Name         | Type   | Description        |
|--------------|--------|--------------------|
| bk_obj_id    | string | Node type          |
| bk_inst_id   | int    | Node instance ID   |
| bk_inst_name | string | Node instance name |
