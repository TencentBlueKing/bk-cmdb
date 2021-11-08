### Functional description

delete set in batches

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field      |  Type      | Required	   |  Description      |
|-----------|------------|--------|------------|
| bk_biz_id | int        | Yes     | Business ID     |
| inst_ids  | int array  | Yes     | Set ID group |

### Request Parameters Example

```python
{
    "bk_biz_id":0,
    "delete": {
    "inst_ids": [123]
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
