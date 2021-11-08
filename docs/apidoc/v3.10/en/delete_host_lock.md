### Functional description

delete host lock (v3.8.6)

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type       | Required	   | Description                            |
|---------------------|-------------|--------|----------------------------------|
|id_list| string| yes| host innerip|

### Request Parameters Example

```python
{
   "id_list":[1, 2, 3]
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "success",
    "data": null
}
```
