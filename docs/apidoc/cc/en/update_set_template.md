### Functional description

According to the service id and the set template id, edit the set template under the specified service

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                  | Type   | Required   | Description           |
| -------------------- | ------ | ----- | -------------- |
| bk_biz_id            |  int    | yes | Business ID |
| set_template_id      |  int    | yes | Set template ID |
| name                 |  string |Either service_template_ids or service_template_ids is required, or both| Set template name |
| service_template_ids | array  |And name are required, or both| Service template ID list|


### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error   |
| message | string |Error message returned by request failure                   |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                          |

#### Data field Description

| Field                | Type   | Description         |
| ------------------- | ------ | ------------ |
| id                  |  int    | Set template ID |
| name                |  string  |Set template name|
| bk_biz_id           |  int    | Business ID |
| version             |  int    | Set template version |
| creator             |  string |Creator       |
| modifier            |  string |Last modified by|
| create_time         |  string |Settling time     |
| last_time           |  string |Update time     |
| bk_supplier_account | string |Developer account number   |
