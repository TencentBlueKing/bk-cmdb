### Function Description

Batch create namespace (version: v3.12.1+, auth: create container namespace)

### Request parameters

{{ common_args_desc }}

#### Interface Parameters

| field     | type  | required | description                               |
|-----------|-------|----------|-------------------------------------------|
| bk_biz_id | int   | yes      | business_id                               |
| data      | array | Yes      | namespace array, limited to 200 at a time |

#### data[x]

| field           | type   | required | description                                                                  |
|-----------------|--------|----------|------------------------------------------------------------------------------|
| bk_cluster_id   | int    | yes      | The unique id of the cluster identified in cmdb, passed in with cluster_uid. |
| name            | string | yes      | namespace name                                                               |
| labels          | map    | no       | labels                                                                       |
| resource_quotas | array  | no       | namespace CPU and memory requests and limits                                 |

#### resource_quotas[x]

| field          | type   | required | description                                                                                                                                    |
|----------------|--------|----------|------------------------------------------------------------------------------------------------------------------------------------------------|
| hard           | object | no       | hard limits required per named resource                                                                                                        |
| scopes         | array  | no       | Quota scopes,optional values are: "Terminating", "NotTerminating", "BestEffort", "NotBestEffort", "PriorityClass", "CrossNamespacePodAffinity" |scope_selector
| scope_selector | no     | object   | scope selector                                                                                                                                 |

#### scope_selector

| field             | type | required | description       |
|-------------------|------|----------|-------------------|
| match_expressions | no   | array    | match_expressions |

#### match_expressions[x]

| field      | type   | required | description                                                                                                                                    |
|------------|--------|----------|------------------------------------------------------------------------------------------------------------------------------------------------|
| scope_name | array  | is       | quota scope,optional values are: "Terminating", "NotTerminating", "BestEffort", "NotBestEffort", "PriorityClass", " CrossNamespacePodAffinity" |
| operator   | string | Yes      | selector operator, with optional values "In", "NotIn", "Exists", "DoesNotExist"                                                                |
| values     | array  | no       | Array of strings, cannot be empty if the operator is "In" or "NotIn", must be empty if it is "Exists" or "DoesNotExist"                        |

### Request parameter examples

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 3,
  "data": [
    {
      "bk_cluster_id": 1,
      "name": "test",
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
  ]
}
```

### Return Result Example

```json

{
    "result": true,
    "code": 0,
    "data": {
      "ids": [1]
    },
    "message": "success",
    "permission": null,
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

**Note:**

- The order of the namespace ID array in the returned data is consistent with the order of the array data in the
  parameter.

### Return result parameter description

#### response

| name       | type   | description                                                                               |
|------------|--------|-------------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. true:request successful; false request failed. |
| code       | int    | The error code. 0 means success, >0 means failure error.                                  |
| message    | string | The error message returned by the failed request.                                         |
| permission | object | Permission information                                                                    |
| request_id | string | request_chain_id                                                                          |
| data       | object | data returned by the request                                                              |

#### data

| field | type  | description                                     |
|-------|-------|-------------------------------------------------|
| ids   | array | array of unique identifiers for namespace in cc |
