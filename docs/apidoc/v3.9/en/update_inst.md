### Functional description

update instance

- the api is just suitable for instances of self-defined mainline model and common model, not suitable for instances of business, set, module, host model, etc.

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type      | Required	   |  Description                            |
|---------------------|------------|--------|----------------------------------|
| bk_supplier_account | string     | Yes     | Supplier account                       |
| bk_obj_id           | string     | Yes     | Object ID       |
| bk_inst_id          | int        | Yes     | Instance ID |
| bk_inst_name        | string     | No     | Field instance ID,also it can be used for custom   |
| bk_biz_id                  | int        | No     | Business ID, when the obj is self-defined mainline model，it must be set |

Note: when the obj is self-defined mainline model with using IAM and cmdb version < 3.9，it must have another param metadata，else it will cause auth fail error，the metadata param format is
"metadata": {
  "label": {
      "bk_biz_id": "64"
  }
}


### Request Parameters Example(General instance example)

```json
{
    "bk_supplier_account": "0",
    "bk_obj_id": "1",
    "bk_inst_id": 0,
    "bk_inst_name": "test"
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
