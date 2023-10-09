### Functional description

According to the service topology instance node, querying the host relationship information under the instance node

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type      | Required   | Description                       |
|---------------------|------------|--------|-----------------------------|
| page       |   dict    | yes  | Query criteria|
| fields    |   array   | yes  | Host attribute list, which controls the fields in the host that returns the result. Please fill them in as required. They can be bk_biz_id,bk_host_id,bk_module_id,bk_set_id,bk_supplier_account|
| bk_obj_id | string |yes| The model ID of the topology node, which can be a user-defined hierarchical model ID, set, module, etc., But can not be a business|
| bk_inst_ids | array |yes| The instance ID of the topology node, supporting up to 50 instance nodes|
| bk_biz_id | int |yes| Business ID |

#### page

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| start    |   int    | yes  | Record start position|
| limit    |   int    | yes  | Limit bars per page, maximum 500|
| sort     |   string |no     | Sort field|

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 1,
    "page": {
        "start": 0,
        "limit": 10
    },
    "fields": [
        "bk_module_id",
        "bk_host_id"
    ],
    "bk_obj_id": "province",
    "bk_inst_ids": [10,11]
}
```

### Return Result Example

```json
{
  "result":true,
  "code":0,
  "message":"success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data":  {
      "count": 1,
      "info": [
          {
              "bk_host_id": 2,
              "bk_module_id": 51
          }
      ]
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

#### data

| Field      | Type      | Description      |
|-----------|-----------|-----------|
| count     |  int       | Number of records|
| info      |  array     | Host relationship information|

#### info 
| Field      | Type      | Description      |
|-----------|-----------|-----------|
| bk_host_id     |  int       | Host id|
| bk_module_id      |  int     | Module id|

