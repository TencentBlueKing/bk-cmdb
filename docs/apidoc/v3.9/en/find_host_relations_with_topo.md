### Functional description

find host relations with topology instance node

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                 |  Type      | Required	   |  Description          |
|----------------------|------------|--------|-----------------------------|
| page       |  dict    | Yes     | search condition |
| fields       |  array string    | Yes     | relation's attribute list need to return, can be: bk_biz_id,bk_host_id,bk_module_id,bk_set_id,bk_supplier_account |
| bk_obj_id | string | Yes | this object's model |
| bk_inst_ids | int array | Yes | this object's instance id list, max size is 50. |

#### page

| Field      |  Type      | Required   |  Description      |
|-----------|------------|--------|------------|
| start    |  int    | Yes     | start record |
| limit    |  int    | Yes     | page limit, maximum value is 500 |


### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_token": "xxx",
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

#### data:

| Field       | Type     | Description         |
|------------|----------|--------------|
| count     | int       | the num of record |
| info      | array     | host data and topology information |
