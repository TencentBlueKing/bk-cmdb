### Function Description

Query instance association topology

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type   | Required | Description |
| ---------- | ------ | -------- | ----------- |
| bk_obj_id  | string | Yes      | Model ID    |
| bk_inst_id | int    | Yes      | Instance ID |

### Request Parameter Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id":"test",
    "bk_inst_id":1
}
```

### Return Result Example

```python
{
    "result": true,
    "code": 0,
    "data": [
        {
            "id": "",
            "bk_obj_id": "biz",
            "bk_obj_icon": "icon-cc-business",
            "bk_inst_id": 0,
            "bk_obj_name": "business",
            "bk_inst_name": "",
            "asso_id": 0,
            "count": 1,
            "children": [
                {
                    "id": "6",
                    "bk_obj_id": "biz",
                    "bk_obj_icon": "icon-cc-business",
                    "bk_inst_id": 6,
                    "bk_obj_name": "business",
                    "bk_inst_name": "",
                    "asso_id": 558
                }
            ]
        }
    ],
    "message": "success",
    "permission": null,
    "request_id": "94c85fdf6a9341e18750a44d6e18c127"
}
```

### Return Result Parameter Explanation

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error   |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain id                                             |
| data       | object | Request returned data                                        |

#### data

| Field        | Type         | Description                                            |
| ------------ | ------------ | ------------------------------------------------------ |
| bk_inst_id   | int          | Instance ID                                            |
| bk_inst_name | string       | Display name of the instance                           |
| bk_obj_icon  | string       | Icon name of the model                                 |
| bk_obj_id    | string       | Model ID                                               |
| bk_obj_name  | string       | Display name of the model                              |
| children     | object array | Collection of all instances associated with this model |
| count        | int          | Number of nodes in children                            |

#### children

| Field        | Type   | Description                  |
| ------------ | ------ | ---------------------------- |
| bk_inst_id   | int    | Instance ID                  |
| bk_inst_name | string | Display name of the instance |
| bk_obj_icon  | string | Icon name of the model       |
| bk_obj_id    | string | Model ID                     |
| bk_obj_name  | string | Display name of the model    |
| asso_id      | string | Association ID               |