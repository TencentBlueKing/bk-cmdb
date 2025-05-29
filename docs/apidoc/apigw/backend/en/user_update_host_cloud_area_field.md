### Description

Update the host's cloud area field based on the host id list and cloud area id

### Parameters

| Name        | Type  | Required | Description          |
|-------------|-------|----------|----------------------|
| bk_biz_id   | int   | no       | Business ID          |
| bk_cloud_id | int   | yes      | Cloud area ID        |
| bk_host_ids | array | yes      | Host IDs, up to 2000 |

### Request Example

```python
{
    "bk_host_ids": [43, 44], 
    "bk_cloud_id": 27,
    "bk_biz_id": 1
}
```

### Response Example

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": ""
}
```

### Response Parameters
