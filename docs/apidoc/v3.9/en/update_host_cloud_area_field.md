### Functional description

update host's cloud area field

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type      | Required	   |  Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier account       |
| bk_biz_id            | int  | No   | Biz ID |
| bk_cloud_id         | int  | Yes   | cloud area ID |
| bk_host_ids         | array  | Yes   | host IDs, max length 2000 |


### Request Parameters Example(General instance example)

```python
{
	"bk_host_ids": [43, 44], 
	"bk_cloud_id": 27,
	"bk_biz_id": 1
}
```

### Return Result Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": ""
}

```

