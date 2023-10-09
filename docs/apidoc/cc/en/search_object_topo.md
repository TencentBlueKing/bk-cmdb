### Functional description

Query general model topology by classification ID of object model

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                  | Type      | Required   | Description                                    |
|----------------------|------------|--------|------------------------------------------|
| bk_classification_id |string      | yes   | The classification ID of the object model, which can only be named by English letter sequence|


### Request Parameters Example

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_classification_id": "test"
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
           "arrows": "to",
           "from": {
               "bk_classification_id": "bk_host_manage",
               "bk_obj_id": "host",
               "bk_obj_name": "Host",
               "position": "{\"bk_host_manage\":{\"x\":-357,\"y\":-344},\"lhmtest\":{\"x\":163,\"y\":75}}",
               "bk_supplier_account": "0"
           },
           "label": "switch_to_host",
           "label_name": "",
           "label_type": "",
           "to": {
               "bk_classification_id": "bk_network",
               "bk_obj_id": "bk_switch",
               "bk_obj_name": "Switch",
               "position": "{\"bk_network\":{\"x\":-172,\"y\":-160}}",
               "bk_supplier_account": "0"
           }
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

| Field       | Type      | Description                               |
|------------|-----------|------------------------------------|
| arrows     |  string    | Take to (one-way) or to,from (two-way)|
| label_name | string    | The name of the Association                     |
| label      |  string    | Indicates by which field From is associated with To     |
| from       |  string    | The English id of the object model, the initiator of the topological relationship|
| to         |  string    | The English ID of the object model, the termination party of the topological relationship|

#### from、to
| Field       | Type      | Description                               |
|------------|-----------|------------------------------------|
|bk_classification_id| string| Class ID|
|  bk_obj_id    | string     | Model id|
|  bk_obj_name    | string     | Model name|
| bk_supplier_account | string |Developer account number   |
| position             |  json object string |Coordinates for front-end presentation   |
