### Functional description

create instance

- the api is just suitable for instances of self-defined mainline model and common model, not suitable for instances of business, set, module, host model, etc.

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                       |  Type      | Required	   |  Description                                      |
|----------------------------|------------|--------|--------------------------------------------|
| bk_obj_id                  | string     | Yes     | Object ID                 |
| bk_supplier_account        | string     | No     | Supplier account, please fill '0' by independent deployment                |
| bk_inst_name | string     | Yes     | Instance ID |
| bk_biz_id                  | int        | No     | Business ID, when the obj is self-defined mainline model，it must be set |
| bk_parent_id                  | int        | No     | when the obj is self-defined mainline model，it must be set, which reprent its parent instance id|

Note: when the obj is self-defined mainline model with using IAM and cmdb version < 3.9，it must have another param metadata，else it will cause auth fail error，the metadata param format is
"metadata": {
  "label": {
      "bk_biz_id": "64"
  }
}


other object's attribute filed can also be the request parameters.

### Request Parameters Example

```json
{
    "bk_obj_id: "switch",
    "bk_inst_name": "example18",
    "bk_supplier_account": "0",
    "bk_biz_id": 0
}
```

### Return Result Example

```json

{
    "result": true,
    "code": 0,
    "message": "",
    "data": {
        "bk_inst_id": 67
    }
}
```

### Return Result Parameters Description

#### data

| Field       | Type      | Description     |
|----------- |-----------|----------|
| bk_inst_id | int       | Instance ID   |
