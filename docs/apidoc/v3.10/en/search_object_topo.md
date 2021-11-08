### Functional description

search object topology

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                  |  Type      | Required	   |  Description                                    |
|----------------------|------------|--------|------------------------------------------|
| bk_classification_id |string      |Yes      | Classification ID, can be named in English alphabet sequence only |


### Request Parameters Example

```python
{
    "bk_classification_id": "test"
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "data": [
        {
           "arrows": "to",
           "from": {
               "bk_classification_id": "bk_host_manage",
               "bk_obj_id": "host",
               "bk_obj_name": "host",
               "position": "{\"bk_host_manage\":{\"x\":-357,\"y\":-344},\"lhmtest\":{\"x\":163,\"y\":75}}",
               "bk_supplier_account": "0"
           },
           "label": "switch_to_host",
           "label_name": "",
           "label_type": "",
           "to": {
               "bk_classification_id": "bk_network",
               "bk_obj_id": "bk_switch",
               "bk_obj_name": "switch",
               "position": "{\"bk_network\":{\"x\":-172,\"y\":-160}}",
               "bk_supplier_account": "0"
           }
        }
   ]
}
```

### Return Result Parameters Description

#### data

| Field       | Type      | Description                               |
|------------|-----------|------------------------------------|
| arrows     | string    | Value to(uniderection) or to,from(bidirectional) |
| label_name | string    | The relationship name                    |
| label      | string    | Indicating which field From is associated with To     |
| from       | string    | Object English ID, initiator of the topological relationship |
| to         | string    | Object English ID,terminate  of topological relationship |
