### Functional description

craate object

### Request Parameters

{{ common_args_desc }}

#### Request Parameters Example

| Field                 |  Type      | Required	   |  Description                                                    |
|----------------------|------------|--------|----------------------------------------------------------|
| creator              |string      | No     | The creator of data                                           |
| bk_classification_id | string     | Yes     | Classification ID of object model, can be named in English alphabet sequence only                  |
| bk_obj_id            | string     | Yes     | Object model ID, can be named in English alphabet sequence only                     |
| bk_obj_name          | string     | Yes     | Object model name, for display, can be named with any language that human can read |
| bk_supplier_account  | string     | Yes     | Supplier account                                               |
| bk_obj_icon          | string     | No     | Icon infomation of object model, dispaly in front, that can refer to[(modleIcon.json)](/static/esb/api_docs/res/cc/modleIcon.json)|


### Request Parameters Example

```python
{
    "creator": "admin",
    "bk_classification_id": "test",
    "bk_obj_name": "test",
    "bk_supplier_account": "0",
    "bk_obj_icon": "icon-cc-business",
    "bk_obj_id": "test"
}
```


### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "data": {
        "id": 1038
    }
}
```

### Return Result Parameters Description

#### data

| Field      | Type      | Description               |
|-----------|-----------|--------------------|
| id        | int       | ID of the new data record |
