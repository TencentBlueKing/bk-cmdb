### Functional description

delete set template

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                | Type   | Required | Description              |
| -------------------- | ------ | -------- | ------------------------ |
| bk_supplier_account  | string | Yes      | Supplier Account Code    |
| bk_biz_id            | int    | Yes      | Business ID              |
| set_template_ids     | array  | Yes      | Set Template ID List     |


### Request Parameters Example

```json
{
    "bk_supplier_account": "0",
    "bk_biz_id": 20,
    "set_template_ids": [59]
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
