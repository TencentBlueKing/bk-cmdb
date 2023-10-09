### Functional description

According to the service id and the set template id, obtaining a service template list of a set template under the specified service

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required   | Description                       |
|---------------------|------------|--------|-----------------------------|
| set_template_id     |  int  |yes   | Set template ID |
| bk_biz_id           |  int        | yes  | Business ID |

### Request Parameters Example

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
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
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request was successful or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | array| Data returned by request|

Data field Description:

| Name| Type| Description|
|---|---|---|
| bk_biz_id           |  int    | Business ID |
| id                  |  int    | Service template ID   |
| name                |  string  |Service template name|
| service_category_id | int    | Service class ID   |
| creator             |  string |Creator       |
| modifier            |  string |Last modified by|
| create_time         |  string |Settling time     |
| last_time           |  string |Update time     |
| bk_supplier_account | string |Developer account number   |
