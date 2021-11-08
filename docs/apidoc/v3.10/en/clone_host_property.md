### Functional description

clone host property

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field        |  Type   | Required	   |  Description                       |
|-------------|---------|--------|-----------------------------|
| bk_org_ip   | string  | Yes     | source host inner ip    |
| bk_dst_ip   | string  | Yes     | dest host inner ip |
| bk_org_id   | int  | Yes    | souce host identify id  |
| bk_dst_id   | int  | Yes    | destination host identify id |
| bk_biz_id   | int     | Yes     | Business ID                      |
| bk_cloud_id | int     | No     | Cloud ID                    |


Note: use host inner ip clone or use host id to clone can only choose one, can not be used both at the same time.

### Request Parameters Example

```json
{
    "bk_biz_id":2,
    "bk_org_ip":"127.0.0.1",
    "bk_dst_ip":"127.0.0.2",
    "bk_cloud_id":0
}
```

or

```json
{
    "bk_biz_id":2,
    "bk_org_id": 10,
    "bk_dst_id": 11,
    "bk_cloud_id":0
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
