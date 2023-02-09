### Functional description

batch update container cluster attribute fields (v3.10.23+, permission: kube cluster editing permissions)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|---------------------|------------|--------|------------|
| bk_biz_id    |  int  | yes     | business ID|
| ids           | array        | no     |unique IDs of the cluster in cmdb|
| data         | object     | yes     | data that needs to be updated|

#### data

| 字段      |  类型      | 必选   |  描述      |
|-----------|------------|--------|------------|
| scheduling_engine |  string  | no  | scheduling engine|
| xid |  string  | no   | associated cluster ID|
| version   |  string  | no   | cluster version |
| network_type   |  string  | no   | network type|
| region |  string  | no    | the region where the cluster is located|
| vpc |  string  | no    | vpc network|
| network |  array  | no    | cluster network|
| type |  string  | no     | cluster type |

**Note:**
- the number of clusters to be updated at one time does not exceed 100

### Request Parameters Example

```json
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "bk_biz_id":2,
    "ids":[
        1
    ],
    "data":{
        "scheduling_engine":"engine1",
        "version":"1.20.6",
        "network_type":"underlay",
        "region":"xxx",
        "vpc":"xxx",
        "network":"127.0.0.0/21",
        "type":"public-cluster"
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
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7",
    "data": null
}
```
### Return Result Parameters Description

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                           |
