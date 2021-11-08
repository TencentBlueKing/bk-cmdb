### Functional description

delete object

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field  |  Type       | Required	   |  Description                 |
|-------|-------------|--------|-----------------------|
| id    | int         | No     | ID of the deleted data record   |


### Request Parameters Example

```python

{
    "id" : 0
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
