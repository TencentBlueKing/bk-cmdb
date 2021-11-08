### Functional description

update object instances in batches

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type       | Required	   |  Description                            |
|---------------------|-------------|--------|----------------------------------|
| bk_supplier_account | string      | Yes     | Supplier account                       |
| bk_obj_id           | string      | Yes     | Object ID                           |
| update              | object array| Yes     | The updated fields and values for the instance             |

#### update
| Field         | Type   | Required	  | Description                           |
|--------------|--------|-------|--------------------------------|
| bk_inst_name | string | No    | Instance name, or custom field |
| datas        | object | Yes    | The updated fields for the instance           |
| inst_id      | int    | Yes    | Point out datas, for instance update    |

#### datas

**datas is an object of map type，key is a field defined by an instance of the model, value is a field **


### Request Parameters Example

```python
{
    "bk_supplier_account":"0",
    "bk_obj_id":"test",
    "update":[
        {
          "datas":{
            "bk_inst_name":"batch_update"
          },
          "inst_id":46
         }
        ]
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
