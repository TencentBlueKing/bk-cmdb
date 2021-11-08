### Functional description

update set template

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                | Type   | Required   | Description              |
| -------------------- | ------ | ---------- | ------------------------ |
| bk_supplier_account  | string | Yes        | Supplier Account Code    |
| bk_biz_id            | int    | Yes        | Business ID              |
| set_template_id      | int    | Yes        | Set Template ID          |
| name                 | string | Choose One | Set Template Name        |
| service_template_ids | array  | Choose One | Service Template ID List |


### Request Parameters Example

```json
{
    "bk_supplier_account": "0",
    "name": "test",
    "bk_biz_id": 20,
    "set_template_id": 6,
    "service_template_ids": [59]
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "permission": null,
    "data": {
        "id": 6,
        "name": "test",
        "bk_biz_id": 20,
        "version": 0,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2019-11-27T17:24:10.671658+08:00",
        "last_time": "2019-11-27T17:24:10.671658+08:00",
        "bk_supplier_account": "0"
    }
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

#### data description

| Field               | Type    | Description           |
| ------------------- | ------- | --------------------- |
| id                  | integer | set template ID       |
| name                | array   | set template name     |
| bk_biz_id           | int     | business ID           |
| version             | int     | set template version  |
| creator             | string  | creator               |
| modifier            | string  | last modifier         |
| create_time         | string  | creation time         |
| last_time           | string  | last modify time      |
| bk_supplier_account | string  | supplier account code |
