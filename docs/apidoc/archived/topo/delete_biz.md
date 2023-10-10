### 批量删除业务

- API: POST /api/v3/deletemany/biz
- API 名称：delete_biz
- 功能说明：
	- 中文：批量删除业务
	- English：batch delete business

- input body

```json
{
    "bk_biz_id": [3,4]
}
```

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
| bk_biz_id | array | 是 | 无 | 用于更新的业务ID列表，最多不能超过20个 | business id list. the max length is 20 |

- output

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": null
}
```

注意：当删除的业务列表中，只要有一个没有权限，则会全部删除失败。
