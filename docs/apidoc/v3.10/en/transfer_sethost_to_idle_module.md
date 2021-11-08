### Functional description

 transfer sethost to idle module

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field          |  Type      | Required	     |  Description    |
|---------------|------------|----------|----------|
| bk_biz_id     | int        | Yes       | Business ID   |
| bk_set_id     | int        | Yes       | Set ID   |
| bk_module_id  | int        | Yes       | Module ID   |


### Request Parameters Example

```python
{
    "bk_biz_id":10,
    "bk_module_id":58,
    "bk_set_id":1
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "data": "sucess"
}
```
