### Functional description

find biz info by host id

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                 |  Type      | Required	   |  Description          |
|----------------------|------------|--------|-----------------------------|
| bk_host_id            | int array  | Yes    | host ID array ,the array length need be less than 500     |
| bk_biz_id             | int     | No    | business ID  |

### Request Parameters Example

```json
{
    "bk_host_id": [
        3,
        4
    ]
}
```

### Return Result Example

```json
{
  "result":true,
  "code":0,
  "message":"success",
  "data": [
    {
      "bk_biz_id": 3,
      "bk_host_id": 3,
      "bk_module_id": 59,
      "bk_set_id": 11,
      "bk_supplier_account": "0"
    },
    {
      "bk_biz_id": 3,
      "bk_host_id": 3,
      "bk_module_id": 60,
      "bk_set_id": 11,
      "bk_supplier_account": "0"
    },
    {
      "bk_biz_id": 3,
      "bk_host_id": 3,
      "bk_module_id": 61,
      "bk_set_id": 12,
      "bk_supplier_account": "0"
    },
    {
      "bk_biz_id": 3,
      "bk_host_id": 4,
      "bk_module_id": 60,
      "bk_set_id": 11,
      "bk_supplier_account": "0"
    }
  ]
}
```


### Return Result Parameters Description

#### data ：

| Field       | Type     | Description         |
|------------|----------|--------------|
| bk_biz_id| int| business ID |
| bk_host_id| int | host ID |
| bk_module_id| int| module ID |
| bk_set_id| int | set ID |
| bk_supplier_account| string| supplier account |