### Functional description

Host query without business information

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id | int        | no     | Business ID |
| page       |   object    | yes  | Query criteria|
| host_property_filter|  object| no | Host attribute combination query criteria|
| fields  |  array   | yes  | Host attribute list, which controls which fields are in the host that returns the result, can speed up interface requests and reduce network traffic   |

#### host_property_filter
This parameter is a combination of filtering rules for the host attribute field and is used to search for hosts based on the host attribute field. The combination supports AND and OR, and can be nested, with a maximum of 2 layers.
The filtering rule is a quadruple`field`,`operator`,`value`

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| condition       |   string    | no     | Combined query criteria|
| rules      |   array    | no     | Rule|


#### rules
| Name| Type| Required| Default value|  Description|
| ---  | ---  | --- |---  | ---|
| field| string| yes | None| Field name| Field name|
| operator| string| yes | None| Operator| Optional values equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between|
| value| - |no| None| Operand| Different values correspond to different value formats|

Assembly rules can be found at: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md



#### page

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| start    |   int    | yes  | Record start position|
| limit    |   int    | yes  | Limit bars per page, Max. 500|



### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_supplier_account": "0",
    "page": {
        "start": 0,
        "limit": 3
    },
    "fields": [
        "bk_host_id",
        "bk_cloud_id",
        "bk_host_innerip",
        "bk_os_type",
        "bk_mac"
    ],
    "host_property_filter": {
        "condition": "AND",
        "rules": [
            {
                "field": "bk_host_innerip",
                "operator": "equal",
                "value": "127.0.0.1"
            },
            {
                "condition": "OR",
                "rules": [
                    {
                        "field": "bk_os_type",
                        "operator": "not_in",
                        "value": [
                            "3"
                        ]
                    },
                    {
                        "field": "bk_cloud_id",
                        "operator": "equal",
                        "value": 0
                    }
                ]
            }
        ]
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
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": {
    "count": 30,
    "info": [
      {
        "bk_cloud_id": 0,
        "bk_host_id": 1,
        "bk_host_innerip": "192.168.15.18",
        "bk_mac": "",
        "bk_os_type": null
      },
      {
        "bk_cloud_id": 0,
        "bk_host_id": 2,
        "bk_host_innerip": "192.168.15.4",
        "bk_mac": "",
        "bk_os_type": null
      },
      {
        "bk_cloud_id": 0,
        "bk_host_id": 3,
        "bk_host_innerip": "192.168.15.12",
        "bk_mac": "",
        "bk_os_type": null
      }
    ]
  }
}
```

### Return Result Parameters Description
#### response

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request was successful or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | object |Data returned by request|

#### data

| Field      | Type      | Description      |
|-----------|-----------|-----------|
| count     |  int       | Number of records|
| info      |  array     | Host actual data|

#### data.info
| Name| Type| Description|
| ---------------- | ------ | -------------------------------  |
| bk_os_type       |  string |Operating system type| 1:Linux;2:Windows; 3:AIX         |
| bk_mac           |  string |Intranet MAC address   |                                 |
| bk_host_innerip  | string |Intranet IP        |                                 |
| bk_host_id       |  int    | Host ID        |                                 |
| bk_cloud_id      |  int    | Cloud area    |                                 |
