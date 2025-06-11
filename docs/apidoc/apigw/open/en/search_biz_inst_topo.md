### Description

Query Business Instance Topology

### Parameters

| Name                | Type   | Required | Description                                                                                                                          |
|---------------------|--------|----------|--------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id           | int    | Yes      | Business ID                                                                                                                          |

### Request Example

```json
{
    "bk_biz_id": 1
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
            "bk_inst_id": 2,
            "bk_inst_name": "blueking",
            "bk_obj_id": "biz",
            "bk_obj_name": "business",
            "default": 0,
            "child": [
                {
                    "bk_inst_id": 3,
                    "bk_inst_name": "job",
                    "bk_obj_id": "set",
                    "bk_obj_name": "set",
                    "default": 0,
                    "child": [
                        {
                            "bk_inst_id": 5,
                            "bk_inst_name": "job",
                            "bk_obj_id": "module",
                            "bk_obj_name": "module",
                            "child": []
                        }
                    ]
                }
            ]
        }
    ]
}
```

### Response Parameters

| Name       | Type   | Description                                                         |
|------------|--------|---------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure               |
| message    | string | Error message returned in case of request failure                   |
| permission | object | Permission information                                              |
| data       | object | Request returned data                                               |

#### data

| Name         | Type   | Description                                                    |
|--------------|--------|----------------------------------------------------------------|
| bk_inst_id   | int    | Instance ID                                                    |
| bk_inst_name | string | Name used for displaying the instance                          |
| bk_obj_icon  | string | Model icon name                                                |
| bk_obj_id    | string | Model ID                                                       |
| bk_obj_name  | string | Name used for displaying the model                             |
| child        | array  | Collection of all instances under the current node             |
| default      | int    | 0-ordinary cluster, 1-built-in module collection, default is 0 |

#### child

| Name         | Type   | Description                                                    |
|--------------|--------|----------------------------------------------------------------|
| bk_inst_id   | int    | Instance ID                                                    |
| bk_inst_name | string | Name used for displaying the instance                          |
| bk_obj_icon  | string | Model icon name                                                |
| bk_obj_id    | string | Model ID                                                       |
| bk_obj_name  | string | Name used for displaying the model                             |
| child        | array  | Collection of all instances under the current node             |
| default      | int    | 0-ordinary cluster, 1-built-in module collection, default is 0 |
