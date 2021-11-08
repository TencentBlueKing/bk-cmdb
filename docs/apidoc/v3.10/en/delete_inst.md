### Functional description

delete instance

- the api is just suitable for instances of self-defined mainline model and common model, not suitable for instances of business, set, module, host model, etc.

### Request Parameters

{{ common_args_desc }}

#### Request Parameters Example

| Field                |  Type       | Required	   |  Description                            |
|---------------------|-------------|--------|----------------------------------|
| bk_supplier_account | string      | Yes     | Supplier account                       |
| bk_obj_id           | string      | Yes     | Object ID |
| bk_inst_id          | int         | Yes     | instance ID  |
| bk_biz_id                  | int        | No     | Business ID, when the obj is self-defined mainline model，it must be set |

Note: when the obj is self-defined mainline model with using IAM and cmdb version < 3.9，it must have another param metadata，else it will cause auth fail error，the metadata param format is
"metadata": {
  "label": {
      "bk_biz_id": "64"
  }
}


### Request Parameters Example

```json

{
    "bk_supplier_account": "0",
    "bk_obj_id": "test",
    "bk_inst_id": 0
}
```


### Return Result Example

```json

{
    "result": true,
    "code": 0,
    "message": "",
    "data": "success"
}
```
