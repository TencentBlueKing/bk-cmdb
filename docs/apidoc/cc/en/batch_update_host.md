### Functional description

Bulk update host attributes based on host id and attributes (can not be used to update cloud area field in host attributes)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type         | Required   | Description                           |
|---------------------|--------------|--------|---------------------------------|
| update              |  array |yes     | Host updated attributes and values, up to 500   |

#### update
| Field        | Type    | Required   | Description                                                |
|-------------|--------|--------|----------------------------------------------------|
| properties  | object |yes     | The updated properties and values of the host can not be used to update the cloud area field in the host properties |
| bk_host_id  | int    | yes     | Host ID for update                                     |

#### properties
| Field         | Type   | Required   | Description                                                      |
|--------------|--------|-------|-----------------------------------------------------------|
| bk_host_name | string |no    | The host name, which can also be another attribute, can not be used to update the cloud area field in the host properties |
| operator     |  string |no    | The primary maintainer, which can also be another attribute, can not be used to update the cloud area field in the host properties|
| bk_comment   |  string |no    | Note, which can be other properties, can not be used to update the cloud area field in host properties  |
| bk_isp_name  | string |no    | The operator, or other attributes, can not be used to update the cloud area field in the host attribute|



### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "update":[
      {
        "properties":{
          "bk_host_name":"batch_update",
          "operator": "admin",
          "bk_comment": "test",
          "bk_isp_name": "1"
        },
        "bk_host_id":46
      }
    ]
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

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |
