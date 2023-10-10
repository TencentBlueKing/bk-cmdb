### Functional description

Create a set template with the specified name under the specified service id, and the set template created to contain the service template by the specified service template id

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type   | Required| Description           |
| -------------------- | ------ | ---- | -------------- |
| bk_biz_id            |  int    | yes   | Business ID |
| name                 |  string |yes   | Set template name |
| service_template_ids | array  |yes   | Service template ID list|


### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "0",
    "name": "test",
    "bk_biz_id": 20,
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
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": {
        "id": 6,
        "name": "test",
        "bk_biz_id": 20,
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

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error   |
| message | string |Error message returned by request failure                   |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                          |

#### Data field Description

| Field                | Type   | Description         |
| ------------------- | ------ | ------------ |
| id                  |  int    | Set template ID |
| name                |  array  |Set template name|
| bk_biz_id           |  int    | Business ID |
| creator             |  string |Creator       |
| modifier            |  string |Last modified by|
| create_time         |  string |Settling time     |
| last_time           |  string |Update time     |
| bk_supplier_account | string |Developer account number   |
