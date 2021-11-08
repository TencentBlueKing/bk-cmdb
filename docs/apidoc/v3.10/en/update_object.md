### Functional description

update object

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type              | Required	   |  Description                                   |
|---------------------|--------------------|--------|-----------------------------------------|
| id                  | int                | No     | ID of target object data record, as a condition for update     |
| modifier            | string             | No     | Last editor of data    |
| bk_classification_id| string             | Yes     | Classification ID, can be named in English alphabet sequence only|
| bk_obj_name         | string             | No     | Object name                           |
| bk_supplier_account | string             | Yes     | Supplier account                              |
| bk_obj_icon         | string             | No     | Object icon information, display in front, that can refer to [(modleIcon.json)](/static/esb/api_docs/res/cc/modleIcon.json)|
| position            | json object string | No     |  Position of front display                     |



### Request Parameters Example

```python
{
    "id": 1,
    "modifier": "admin",
    "bk_classification_id": "cc_test",
    "bk_obj_name": "cc2_test_inst",
    "bk_supplier_account": "0",
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
    "data": "success"
}
```
