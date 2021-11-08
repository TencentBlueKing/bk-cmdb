### Functional description

update business enable status

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type      | Required	   |  Description      |
|---------------------|------------|--------|------------|
| bk_biz_id           | int        | Yes     | Business ID     |
| bk_supplier_account | string     | Yes     | Supplier ID   |
| flag                | string     | Yes     |  Enable status, disabled or enable  |

### Request Parameters Example

```python
{
    "bk_biz_id": "3",
    "bk_supplier_account": "0",
    "flag": "enable"
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
