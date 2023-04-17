### Functional description

Specify model ID and instance iddelete object instances under the specified model

-  This interface only applies to custom hierarchical models and generic model instances, not to business, set, module, host and other model instances 

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type       | Required   | Description                            |
|---------------------|-------------|--------|----------------------------------                     |
| bk_obj_id           |  string      | yes     | Model ID|
| bk_inst_id          |  int         | yes     | Instance ID   |
| bk_biz_id                  |  int        | no     | Business ID, which must be transferred when deleting a user-defined mainline level model instance|

 Note: If the operation is a user-defined mainline hierarchy model instance and permission Center is used, for the version with CMDB less than 3.9, the metadata parameter containing the service id of the instance needs to be transferred. Otherwise, the permission Center authentication will fail. The format is
"metadata": {
    "label": {
        "bk_biz_id": "64"
    }
}

### Request Parameters Example

```json

{ 
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id": "test",
    "bk_inst_id": 0
}
```


### Return Result Example

```json

{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": "success"
}
```
#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |

