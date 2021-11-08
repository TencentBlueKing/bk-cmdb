### Functional description

delete object attribute, can delete business custom field

### Request Parameters

{{ common_args_desc }}

#### Request Parameters Example

| Field  |  Type       | Required	   |  Description                         |
|-------|-------------|--------|-------------------------------|
| id    | int         | No     | The unique identifier ID of the deleted data record  |


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
