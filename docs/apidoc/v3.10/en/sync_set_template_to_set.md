### Functional description

sync set template to set

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field               | Type   | Required | Description              |
| ------------------- | ------ | -------- | ------------------------ |
| bk_supplier_account | string | Yes      | Supplier Account Code    |
| bk_biz_id           | int    | Yes      | Business ID              |
| set_template_id     | int    | Yes      | Set Template ID          |
| bk_set_ids          | array  | Yes      | IDs Of Set To Sync       |


### Request Parameters Example

```json
{
    "bk_supplier_account": "0",
    "bk_biz_id": 20,
    "set_template_id": 6,
    "bk_set_ids": [46]
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": null
}
```

### Return Result Parameters Description

#### response

| Field   | Type   | Description                                            |
| ------- | ------ | ------------------------------------------------------ |
| result  | bool   | request success or failed. true:success；false: failed |
| code    | int    | error code. 0: success, >0: something error            |
| message | string | error info description                                 |
| data    | object | response data                                          |
