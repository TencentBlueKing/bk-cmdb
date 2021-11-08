### Functional description

get mainline object's business topology

### Request Parameters

{{ common_args_desc }}

#### Request Parameters Example

| Field                 |  Type      | Required	   |  Description                 |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     | Yes     | Supplier account            |

### Request Parameters Example

``` python
{
    "bk_supplier_account":"0"
}
```

### Return Result Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "data": [
    {
      "bk_obj_id": "biz",
      "bk_obj_name": "business",
      "bk_supplier_account": "0",
      "bk_next_obj": "set",
      "bk_next_name": "set",
      "bk_pre_obj_id": "",
      "bk_pre_obj_name": ""
    },
    {
      "bk_obj_id": "set",
      "bk_obj_name": "set",
      "bk_supplier_account": "0",
      "bk_next_obj": "module",
      "bk_next_name": "module",
      "bk_pre_obj_id": "biz",
      "bk_pre_obj_name": "business"
    },
    {
      "bk_obj_id": "module",
      "bk_obj_name": "module",
      "bk_supplier_account": "0",
      "bk_next_obj": "host",
      "bk_next_name": "host",
      "bk_pre_obj_id": "set",
      "bk_pre_obj_name": "set"
    },
    {
      "bk_obj_id": "host",
      "bk_obj_name": "host",
      "bk_supplier_account": "0",
      "bk_next_obj": "",
      "bk_next_name": "",
      "bk_pre_obj_id": "module",
      "bk_pre_obj_name": "module"
    }
  ]
}
```

### Return Result Parameters Description

#### data

| Field       | Type     | Description         |
|------------|----------|--------------|
|bk_obj_id | string | object's unique id |
|bk_obj_name | string | object's name |
|bk_supplier_account | string | supplier's account |
|bk_next_obj | string | the next object's unique id |
|bk_next_name | string | the next object's name |
|bk_pre_obj_id | string | the previous object's unique id |
|bk_pre_obj_name | string | the previous object's name |
