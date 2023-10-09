### Functional description

Update model definition

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type              | Required   | Description                                   |
|---------------------|--------------------|--------|-----------------------------------------|
| id                  |  int                | no     | The ID of the object model as a condition for the update operation    |
| modifier            |  string             | no     | The last person to modify this piece of data    |
| bk_classification_id|  string             | yes  | The classification ID of the object model, which can only be named by English letter sequence|
| bk_obj_name         |  string             | no     | The name of the object model                          |
| bk_obj_icon         |  string             | no     | ICON information of object model, used for front-end display, value can be referred to [(modleIcon.json)](/static/esb/api_docs/res/cc/modleIcon.json)|
| position            |  json object string |no     | Coordinates for front-end presentation                      |



### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "id": 1,
    "modifier": "admin",
    "bk_classification_id": "cc_test",
    "bk_obj_name": "cc2_test_inst",
    "bk_obj_icon": "icon-cc-business",
    "position":"{\"ff\":{\"x\":-863,\"y\":1}}"
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
    "data": "success"
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
| data | object |No data return|
