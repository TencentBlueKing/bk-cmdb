### Functional description

Query model based on optional criteria

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required   | Description                                                    |
|----------------------|------------|--------|----------------------------------------------------------|
| creator              |  string     | no     | Creator of this data                                           |
| modifier             |  string     | no     | The last person to modify this piece of data                                   |
| bk_classification_id | string     | no     | The classification ID of the object model, which can only be named by English letter sequence                 |
| bk_obj_id            |  string     | no     | The ID of the object model, which can only be named in English letter sequence                     |
| bk_obj_name          |  string     | no     | The name of the object model, for presentation, can be used in any language that humans can read|

### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "creator": "user",
    "modifier": "user",
    "bk_classification_id": "test",
    "bk_obj_id": "biz"
    "bk_obj_name": "aaa"
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": [
        {
            "bk_classification_id": "bk_organization",
            "create_time": "2018-03-08T11:30:28.005+08:00",
            "creator": "cc_system",
            "description": "",
            "id": 4,
            "bk_ispaused": false,
            "ispre": true,
            "last_time": null,
            "modifier": "",
            "bk_obj_icon": "icon-XXX",
            "bk_obj_id": "XX",
            "bk_obj_name": "XXX",
            "position": "{\"test_obj\":{\"x\":-253,\"y\":137}}",
            "bk_supplier_account": "0"
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

#### data

| Field                 | Type               | Description                                                                                           |
|----------------------|--------------------|------------------------------------------------------------------------------------------------|
| id                   |  int                | ID of the data record                                                                                   |
| creator              |  string             | Creator of this data                                                                                 |
| modifier             |  string             | The last person to modify this piece of data                                                                         |
| bk_classification_id | string             | The classification ID of the object model, which can only be named by English letter sequence                                                       |
| bk_obj_id            |  string             | The ID of the object model, which can only be named by English letter sequence                                                           |
| bk_obj_name          |  string             | The name of the object model, used to show                                                                       |
| bk_supplier_account  | string             | Developer account number                                                                                     |
| bk_ispaused          |  bool               | Disable, true or false                                                                        |
| ispre                |  bool               | Predefined, true or false                                                                      |
| bk_obj_icon          |  string             | ICON information of object model, used for front-end display, and the value can be referred to [(modleIcon.json)](/static/esb/api_docs/res/cc/modleIcon.json)|
| position             |  json object string |Coordinates for front-end presentation                                                                             |
| description           |  string     | Description information of data                                           |
