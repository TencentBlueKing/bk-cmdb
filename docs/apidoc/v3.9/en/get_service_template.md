### Functional description

get service template

#### General Parameters

{{ common_args_desc }}

### Request Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|-----------------------|
| bk_supplier_account  | string     |Yes     | Supplier Account ID       |
| service_template_id | int  | Yes   | Service Template ID |

### Request Parameters Example

```json
{
  "bk_supplier_account": "0",
  "service_template_id": 51
}
```

### Return Result Example

```json
{
   "result": true,
   "code": 0,
   "message": "success",
   "data": {
       "bk_biz_id": 3,
       "id": 51,
       "name": "mm2",
       "service_category_id": 12,
       "creator": "admin",
       "modifier": "admin",
       "create_time": "2020-05-26T09:46:15.259Z",
       "last_time": "2020-05-26T09:46:15.259Z",
       "bk_supplier_account": "0"
   }
}
```

### Return Result Parameters Description

#### response

| Field       | Type     | Description         |
|---|---|---|
| result | bool | request success or failed. true:success；false: failed |
| code | int | error code. 0: success, >0: something error |
| message | string | error info description |
| data | object | response data |

#### data description

| Field       | Type     | Description         |
|---|---|---|
|id|integer|Service Template ID|
|name|array|Service Template name|
|service_category_id|integer|Service Category ID|
