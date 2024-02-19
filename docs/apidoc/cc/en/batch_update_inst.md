### Function Description

Batch Update Object Instances (Permission: Model Instance Editing Permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description                                   |
| --------- | ------ | -------- | --------------------------------------------- |
| bk_obj_id | string | Yes      | Model ID                                      |
| update    | array  | Yes      | Fields and values to be updated for instances |

#### update

| Field   | Type   | Required | Description                                            |
| ------- | ------ | -------- | ------------------------------------------------------ |
| datas   | object | Yes      | Fields and values to be updated for instances          |
| inst_id | int    | Yes      | Specific instance for which datas is used for updating |

#### datas

| Field        | Type   | Required | Description                                    |
| ------------ | ------ | -------- | ---------------------------------------------- |
| bk_inst_name | string | No       | Instance name, can also be other custom fields |

**datas is a map-type object, where the key is the field defined in the model for the instance, and the value is the value of the field**

### Request Parameter Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id": "test",
    "update": [
        {
            "datas": {
                "bk_inst_name": "batch_update"
            },
            "inst_id": 46
        }
    ]
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": "success"
}
```

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error  |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Data returned in the request                                 |