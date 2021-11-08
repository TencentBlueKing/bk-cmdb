### Functional description

list set templates

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field               | Type   | Required | Description           |
| ------------------- | ------ | -------- | --------------------- |
| bk_supplier_account | string | Yes      | Supplier Account Code |
| bk_biz_id           | int    | Yes      | Business ID           |
| set_template_ids    | array  | No       | Set Template ID Array |
| page                | object | No       | page info             |

#### page

| Field | Type   | Required | Description                                       |
| ----- | ------ | -------- | ------------------------------------------------- |
| start | int    | No       | start record                                      |
| limit | int    | No       | page limit, maximum value is 1000                 |
| sort  | string | No       | the field for sort, '-' represent decending order |

### Request Parameters Example

```json
{
  "bk_supplier_account": "0",
  "bk_biz_id": 10,
  "set_template_ids":[1, 11],
  "page": {
    "start": 0,
    "limit": 10,
    "sort": "-name"
  }
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
    "count": 2,
    "info": [
      {
        "id": 1,
        "name": "zk1",
        "bk_biz_id": 10,
        "version": 0,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2020-03-16T15:09:23.859+08:00",
        "last_time": "2020-03-25T18:59:00.167+08:00",
        "bk_supplier_account": "0"
      },
      {
        "id": 11,
        "name": "q",
        "bk_biz_id": 10,
        "version": 0,
        "creator": "admin",
        "modifier": "admin",
        "create_time": "2020-03-16T15:10:05.176+08:00",
        "last_time": "2020-03-16T15:10:05.176+08:00",
        "bk_supplier_account": "0"
      }
    ]
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

#### Data field description

| Field | Type  | Description   |
| ----- | ----- | ------------- |
| count | int   | total count   |
| info  | array | response data |

#### Info field description

#### data description

| Field               | Type    | Description           |
| ------------------- | ------- | --------------------- |
| id                  | int     | set template ID       |
| name                | array   | set template name     |
| bk_biz_id           | int     | business ID           |
| version             | int     | set template version  |
| creator             | string  | creator               |
| modifier            | string  | last modifier         |
| create_time         | string  | creation time         |
| last_time           | string  | last modify time      |
| bk_supplier_account | string  | supplier account code |
