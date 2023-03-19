### Functional description

Query hosts in resource pool

### Request Parameters
{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| page       |   dict    | no     | Query criteria|
| host_property_filter|  object| no | Host attribute combination query criteria|
| fields  |  array   | yes  | Host attribute list, which controls which fields are in the host that returns the result, can speed up interface requests and reduce network traffic transmission   |

#### host_property_filter

This parameter is a combination of filtering rules for the host attribute field and is used to search for hosts based on the host attribute field. The combination supports AND and OR, and can be nested, with a maximum of 2 layers.

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| condition       |   string    | no     ||
| rules      |   array    | no     | Rule|

#### rules
The filtering rule is a quadruple`field`,`operator`,`value`

| Name| Type| Required| Default value| Description|  Description|
| ---  | ---  | --- |---  | --- | ---|
| field| string| yes | None| Field name| Field name|
| operator| string| yes | None| Operator| Optional values equal,not_equal,in,not_in,less,less_or_equal,greater,greater_or_equal,between,not_between|
| value|  string |no| None| Operand| Different values correspond to different value formats|

Assembly rules can be found at: https://github.com/Tencent/bk-cmdb/blob/master/src/common/querybuilder/README.md



#### page

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| start    |   int    | yes  | Record start position|
| limit    |   int    | yes  | Limit bars per page, Max. 500|
| sort     |   string |no     | Sort field|



### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "page": {
        "start": 0,
        "limit": 10,
        "sort": "bk_host_id"
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
            "field": "bk_host_outerip",
            "operator": "equal",
            "value": "127.0.0.1"
        }, {
            "condition": "OR",
            "rules": [{
                "field": "bk_os_type",
                "operator": "not_in",
                "value": ["3"]
            }, {
                "field": "bk_sla",
                "operator": "equal",
                "value": "1"
            }]
        }]
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
    "count": 1,
    "info": [
      {
        "bk_cloud_id": "0",
        "bk_host_id": 17,
        "bk_host_innerip": "192.168.1.1",
        "bk_mac": "",
        "bk_os_type": "1"
      }
    ]
  }
}
```

### Return Result Parameters Description
#### response

| Name| Type| Description|
|---|---|---|
| result | bool |Whether the request succeeded or not. True: request succeeded;false request failed|
| code | int |Wrong code. 0 indicates success,>0 indicates failure error|
| message | string |Error message returned by request failure|
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data | array |Data returned by request|

#### data

| Field      | Type      | Description      |
|-----------|-----------|-----------|
| count     |  int       | Number of records|
| info      |  array     | Host actual data|

#### data.info
| Name             | Type   |  Description                     |
| ---------------- | ------ | -------------------------------  |
| bk_os_type       |  string |Operating system type| 1:Linux;2:Windows; 3:AIX         |
| bk_mac           |  string |Intranet MAC address   |                                 |
| bk_host_innerip  | string |Intranet IP        |                                 |
| bk_host_id       |  int    | Host ID        |                                 |
| bk_cloud_id      |  int    | Cloud area    ||
