### Functional description

delete object instances in batches

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type       | Required	   | Description                            |
|---------------------|-------------|--------|----------------------------------|
| bk_supplier_account | string      | Yes     | Supplier account                       |
| bk_obj_id           | string      | Yes     | Object ID |
| inst_ids            | int array   |Yes      | Instance ID group                       |


### Request Parameters Example

```python
{
    "bk_supplier_account": "0",
    "bk_obj_id": "test",
    "delete":{
    "inst_ids":[123]
    }
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
