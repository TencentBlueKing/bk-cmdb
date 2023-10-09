### Functional description

Modeling

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required   | Description                                                    |
|----------------------|------------|--------|----------------------------------------------------------|
| creator              | string      | no     | Creator of this data                                           |
| bk_classification_id | string     | yes     | The classification ID of the object model, which can only be named by English letter sequence                 |
| bk_obj_id            |  string     | yes     | The ID of the object model, which can only be named in English letter sequence                     |
| bk_obj_name          |  string     | yes     | The name of the object model, for presentation, can be used in any language that humans can read|
| bk_obj_icon          |  string     | no     | ICON information for the object model for front-end display|


### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "creator": "admin",
    "bk_classification_id": "test",
    "bk_obj_name": "test",
    "bk_obj_icon": "icon-cc-business",
    "bk_obj_id": "test"
}
```


### Return Result Example

```python

{
    "code": 0,
    "permission": null,
    "result": true,
    "request_id": "b529879b85c74e3c91b3d8119df8dbc7",
    "message": "success",
    "data": {
        "description": "",
        "bk_ishidden": false,
        "bk_classification_id": "test",
        "creator": "admin",
        "bk_obj_name": "test",
        "bk_ispaused": false,
        "last_time": null,
        "bk_obj_id": "test",
        "create_time": null,
        "bk_supplier_account": "0",
        "position": "",
        "bk_obj_icon": "icon-cc-business",
        "modifier": "",
        "id": 2000002118,
        "ispre": false
    }
}

```

### Return Result Parameters Description
#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |

#### data

| Field      | Type      | Description               |
|-----------|-----------|--------------------|
| id        |  int       | ID of the new data record|
| bk_classification_id | int    | Classification ID of the object model   |
| creator             |  string |Creator       |
| modifier            |  string |Last modified by|
| create_time         |  string |Settling time     |
| last_time           |  string |Update time     |
| bk_supplier_account | string |Developer account number   |
| bk_obj_id | string |Model type   |
| bk_obj_name | string |Model name   |
| bk_obj_icon          |  string             | ICON information for the object model for front-end display|
| position             |  json object string |Coordinates for front-end presentation   |
| ispre                |  bool               | Predefined, true or false   |