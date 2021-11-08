### Functional description

delete dynamic group (V3.9.6)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field               | Type    | Required  | Description                            |
|---------------------|---------|-----------|----------------------------------------|
| bk_biz_id           | int     | Yes       | Business ID                            |
| id                  | string  | Yes       | Primary key ID of target dynamic group |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
    "id": "XXXXXXXX"
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "data": {}
}
```
