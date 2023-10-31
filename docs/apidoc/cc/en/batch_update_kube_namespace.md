### Function Description

Batch update namespace (version: v3.12.1+, auth: edit container namespace)

### Request parameters

{{ common_args_desc }}

#### Interface Parameters

| field     | type   | required | description                             |
|-----------|--------|----------|-----------------------------------------|
| bk_biz_id | int    | yes      | business id                             |
| data      | object | yes      | contains the fields to be updated       |
| ids       | array  | yes      | an array of id unique identifiers in cc |

#### data

| field           | type  | required | description                                  |
|-----------------|-------|----------|----------------------------------------------|
| labels          | map   | no       | labels                                       |
| resource_quotas | array | no       | namespace CPU and memory requests and limits |

### Example request parameters

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 1,
  "ids": [
    1,
    2,
    3
  ],
  "data": {
    "labels": {
      "test": "test",
      "test2": "test2"
    },
    "resource_quotas": [
      {
        "hard": {
          "memory": "20000Gi",
          "pods": "100",
          "cpu": "10k"
        },
        "scope_selector": {
          "match_expressions": [
            {
              "values": [
                "high"
              ],
              "operator": "In",
              "scope_name": "PriorityClass"
            }
          ]
        }
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
  "data": null,
  "message": "success",
  "permission": null,
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### Return result parameter description

#### response

| name       | type   | description                                                                               |
|------------|--------|-------------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. true:request successful; false request failed. |
| code       | int    | The error code. 0 means success, >0 means failure error.                                  |
| message    | string | The error message returned by the failed request.                                         |
| permission | object | Permission information                                                                    |
| request_id | string | request_chain_id                                                                          |
| data       | object | The data returned by the request.                                                         |
