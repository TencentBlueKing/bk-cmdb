### Function Description

Query Business Instance Topology

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type   | Required | Description                                                  |
| ------------------- | ------ | -------- | ------------------------------------------------------------ |
| bk_biz_id           | int    | Yes      | Business ID                                                  |
| level               | int    | No       | Topology level index, index value starts from 0, default is 2. When set to -1, the complete business instance topology will be read. |

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request was successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failure        |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Request returned data                                        |

#### data

| Field        | Type   | Description                                                  |
| ------------ | ------ | ------------------------------------------------------------ |
| bk_inst_id   | int    | Instance ID                                                  |
| bk_inst_name | string | Name used for displaying the instance                        |
| bk_obj_icon  | string | Model icon name                                              |
| bk_obj_id    | string | Model ID                                                     |
| bk_obj_name  | string | Name used for displaying the model                           |
| child        | array  | Collection of all instances under the current node           |
| default      | int    | 0-ordinary cluster, 1-built-in module collection, default is 0 |

#### child

| Field        | Type   | Description                                                  |
| ------------ | ------ | ------------------------------------------------------------ |
| bk_inst_id   | int    | Instance ID                                                  |
| bk_inst_name | string | Name used for displaying the instance                        |
| bk_obj_icon  | string | Model icon name                                              |
| bk_obj_id    | string | Model ID                                                     |
| bk_obj_name  | string | Name used for displaying the model                           |
| child        | array  | Collection of all instances under the current node           |
| default      | int    | 0-ordinary cluster, 1-built-in module collection, default is 0 |