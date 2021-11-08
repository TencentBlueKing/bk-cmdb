### Functional description

search host lock. (v3.8.6)

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type       | Required	   | Description                            |
|---------------------|-------------|--------|----------------------------------|
|id_list| int array| yes| host id list|

### Request Parameters Example

```python
{
   "id_list":[1, 2]
}
```

### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "success",
    "data": {
        1: true,
        2: false
    }
}
```

### Return Result Parameters Description
#### data

| Field      | Type         | Description                 |
|-----------|--------------|----------------------|
| data | map[int]bool |the data response,Key is the ID, value is locked status|

