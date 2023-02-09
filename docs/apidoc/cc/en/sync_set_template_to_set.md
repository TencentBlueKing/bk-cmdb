### Functional description

Synchronize that set template unde the specified service to the set according to the service id, the clust template id and the set id list to be synchronized

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                  | Type   | Required| Description           |
| -------------------- | ------ | ---- | ------------- |
| bk_biz_id            |  int    | yes      | Business ID |
| set_template_id      |  int    | yes | Set template ID   |
| bk_set_ids           |  array  |yes   | List of set IDs to be synchronized |


### Request Parameters Example

```json
{

    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": null
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
