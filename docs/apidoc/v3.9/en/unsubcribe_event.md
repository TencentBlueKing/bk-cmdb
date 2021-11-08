### Functional description

unsubscribe event

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field               |  Type      | Required	   |  Description      |
|--------------------|------------|--------|------------|
|bk_supplier_account | string     | Yes     | Supplier account |
|subscription_id     | int        | Yes     | Subscription ID     |

### Request Parameters Example

```python
{
    "bk_supplier_account":"0",
    "subscription_id":1
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "data":"success"
}
```
