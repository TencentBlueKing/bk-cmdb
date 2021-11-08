### Functional description
get business's idle, fault and recycle modules.

### Request Parameters

{{ common_args_desc }}

### Request Parameters Example

``` python

```

### Return Result Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "data": {
    "bk_set_id": 2,
    "bk_set_name": "idle pool",
    "module": [
      {
        "bk_module_id": 3,
        "bk_module_name": "idle machine"
      },
      {
        "bk_module_id": 4,
        "bk_module_name": "fault machine"
      },
      {
        "bk_module_id": 5,
        "bk_module_name": "recycle machine"
      }
    ]
  }
}
```

### Return Result Parameters Description

#### data description

| Field       | Type     | Description         |
|------------|----------|--------------|
|bk_set_id | int64 | the set id that idle, fault and recycle module belongs to  |
|bk_set_name | string |the set name that idle, fault and recycle module belongs to |

#### module description
| Field       | Type     | Description         |
|------------|----------|--------------|
|bk_module_id | int64 | module's id |
|bk_module_name | string |module's name|
