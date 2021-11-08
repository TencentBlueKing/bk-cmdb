### Functional description

Newly added host lock, if the host has been locked, it also prompts the lock to succeed. (v3.8.6)

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type       | Required	   | Description                            |
|---------------------|-------------|--------|----------------------------------|
|id_list| int array| yes| host id list|

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
    "message": "",
    "data": null
}
```
