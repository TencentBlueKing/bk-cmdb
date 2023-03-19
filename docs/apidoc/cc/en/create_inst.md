### Functional description

Create instance

- This interface only applies to custom hierarchical models and generic model instances, not to business, set, module, host and other model instances

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                       | Type      | Required   | Description                                      |
|----------------------------|------------|--------|--------------------------------------------|
| bk_obj_id                  |  string     | yes     | Model ID                 |
| bk_inst_name | string     | yes     | Instance name|
| bk_biz_id                  |  int        | no     | Business ID, which must be transferred when creating a custom mainline level model instance|
| bk_parent_id                  |  int        | no     | It must be passed when creating a custom mainline level model instance, representing the parent level instance ID|

 Note: If the operation is a user-defined mainline hierarchy model instance and permission Center is used, for the version with CMDB less than 3.9, the metadata parameter containing the service id of the instance needs to be transferred. Otherwise, the permission Center authentication will fail. The format is
"metadata": {
    "label": {
        "bk_biz_id": "64"
    }
}

Other fields that belong to instance properties can also be input parameters.


### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id":"test3",
    "bk_inst_name":"example18",
    "bk_biz_id":0
}
```

### Return Result Example

```json

{
    "result": true,
    "code": 0,
    "data": {
        "bk_biz_id": 0,
        "bk_inst_id": 1177099,
        "bk_inst_name": "example18",
        "bk_obj_id": "test3",
        "bk_supplier_account": "0",
        "create_time": "2022-01-05T17:28:27.069+08:00",
        "last_time": "2022-01-05T17:28:27.069+08:00",
        "test4": ""
    },
    "message": "success",
    "permission": null,
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### Return Result Parameters Description
#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |

#### data

| Field       | Type      | Description     |
|----------- |-----------|----------|
| bk_inst_id | int       | Instance id   |
| bk_biz_id |     int   | Business ID |
| bk_inst_name |   string     | Instance name   |
| bk_obj_id |      string  |   Model id|
| bk_supplier_account|  string       | Developer account number                                                 |
| create_time         |  string |Settling time     |
| last_time           |  string |Update time     |

