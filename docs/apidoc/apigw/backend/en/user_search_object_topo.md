### Description

Query the topology of a common model through the classification ID of the object model (Permission: Model Topology View
Edit Permission)

### Parameters

| Name                 | Type   | Required | Description                                                                         |
|----------------------|--------|----------|-------------------------------------------------------------------------------------|
| bk_classification_id | string | Yes      | Classification ID of the object model, can only be named with alphabetical sequence |

### Request Example

```python
{
    "bk_classification_id": "test"
}
```

### Response Example

```python
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": [
        {
           "arrows": "to",
           "from": {
               "bk_classification_id": "bk_host_manage",
               "bk_obj_id": "host",
               "bk_obj_name": "主机",
               "position": "{\"bk_host_manage\":{\"x\":-357,\"y\":-344},\"lhmtest\":{\"x\":163,\"y\":75}}",
               "bk_supplier_account": "0"
           },
           "label": "switch_to_host",
           "label_name": "",
           "label_type": "",
           "to": {
               "bk_classification_id": "bk_network",
               "bk_obj_id": "bk_switch",
               "bk_obj_name": "交换机",
               "position": "{\"bk_network\":{\"x\":-172,\"y\":-160}}",
               "bk_supplier_account": "0"
           }
        }
   ]
}
```

### Response Parameters

| Name       | Type   | Description                                                        |
|------------|--------|--------------------------------------------------------------------|
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 indicates success, >0 indicates failed error         |
| message    | string | Error message returned in case of failure                          |
| permission | object | Permission information                                             |
| data       | object | Request returned data                                              |

#### data

| Name       | Type   | Description                                                                    |
|------------|--------|--------------------------------------------------------------------------------|
| arrows     | string | Value can be "to" (unidirectional) or "to,from" (bidirectional)                |
| label_name | string | Name of the relationship                                                       |
| label      | string | Indicates through which field From is related to To                            |
| from       | string | English id of the object model, the initiator of the topological relationship  |
| to         | string | English ID of the object model, the terminator of the topological relationship |

#### from、to

| Name                 | Type               | Description                            |
|----------------------|--------------------|----------------------------------------|
| bk_classification_id | string             | Classification ID                      |
| bk_obj_id            | string             | Model ID                               |
| bk_obj_name          | string             | Model name                             |
| bk_supplier_account  | string             | Vendor account                         |
| position             | json object string | Coordinates used for front-end display |
