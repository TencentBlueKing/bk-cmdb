### Functional description

list service templates of a set template

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field           | Type   | Required | Description           |
| ----------------| ------ | -------- | --------------------- |
| set_template_id | int    | Yes      | set template ID |
| bk_biz_id       | int    | Yes      | Business ID           |

### Request Parameters Example

```json
{
  "set_template_id": 1,
  "bk_biz_id": 3
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": [
        {
            "bk_biz_id": 3,
            "id": 48,
            "name": "sm1",
            "service_category_id": 2,
            "creator": "admin",
            "modifier": "admin",
            "create_time": "2020-05-15T14:14:57.691Z",
            "last_time": "2020-05-15T14:14:57.691Z",
            "bk_supplier_account": "0"
        },
        {
            "bk_biz_id": 3,
            "id": 49,
            "name": "sm2",
            "": 16,
            "creator": "admin",
            "modifier": "admin",
            "create_time": "2020-05-15T14:19:09.813Z",
            "last_time": "2020-05-15T14:19:09.813Z",
            "bk_supplier_account": "0"
        }
    ]
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

#### Data field description

| Field | Type  | Description   |
| ----- | ----- | ------------- |
| bk_biz_id           | int    | business ID       |
| id                  | int    | service template ID   |
| name                | array  | service template name |
| service_category_id | int    | service category id   |
| creator             | string | creator       |
| modifier            | string | last modifier |
| create_time         | string  | creation time         |
| last_time           | string  | last modify time      |
| bk_supplier_account | string  | supplier account code |

