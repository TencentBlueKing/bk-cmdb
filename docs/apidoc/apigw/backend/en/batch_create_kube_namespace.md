### Description

Batch create namespaces (Version: v3.12.1+, Permission: Container namespace creation permission)

### Parameters

| Name      | Type  | Required | Description                                             |
|-----------|-------|----------|---------------------------------------------------------|
| bk_biz_id | int   | Yes      | Business ID                                             |
| data      | array | Yes      | Array of namespaces, up to 200 can be created at a time |

#### data[x]

| Name            | Type   | Required | Description                                          |
|-----------------|--------|----------|------------------------------------------------------|
| bk_cluster_id   | int    | Yes      | Unique ID that identifies the cluster in CMDB        |
| name            | string | Yes      | Namespace name                                       |
| labels          | map    | No       | Labels                                               |
| resource_quotas | array  | No       | CPU and memory requests and limits for the namespace |

#### resource_quotas[x]

| Name           | Type   | Required | Description                                                                                                                                    |
|----------------|--------|----------|------------------------------------------------------------------------------------------------------------------------------------------------|
| hard           | object | No       | Hard limits for each named resource                                                                                                            |
| scopes         | array  | No       | Quota scope, optional values are: "Terminating", "NotTerminating", "BestEffort", "NotBestEffort", "PriorityClass", "CrossNamespacePodAffinity" |
| scope_selector | object | No       | Scope selector                                                                                                                                 |

#### scope_selector

| Name              | Type  | Required | Description       |
|-------------------|-------|----------|-------------------|
| match_expressions | array | No       | Match expressions |

#### match_expressions[x]

| Name       | Type   | Required | Description                                                                                                                                    |
|------------|--------|----------|------------------------------------------------------------------------------------------------------------------------------------------------|
| scope_name | array  | Yes      | Quota scope, optional values are: "Terminating", "NotTerminating", "BestEffort", "NotBestEffort", "PriorityClass", "CrossNamespacePodAffinity" |
| operator   | string | Yes      | Selector operator, optional values are: "In", "NotIn", "Exists", "DoesNotExist"                                                                |
| values     | array  | No       | String array, must not be empty if the operator is "In" or "NotIn", must be empty if the operator is "Exists" or "DoesNotExist"                |

### Request Example

```json
{
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

### Response Example

```json
{
  "result": true,
  "code": 0,
  "data": {
    "ids": [
      1
    ]
  },
  "message": "success",
  "permission": null,
}
```

**Note:**

- The order of the namespace ID array returned in the data field corresponds to the order of the array data in the
  parameters.

### Response Parameters

| Name       | Type   | Description                                                       |
|------------|--------|-------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error     |
| message    | string | Error message returned for a failed request                       |
| permission | object | Permission information                                            |
| data       | object | Data returned by the request                                      |

#### data

| Name | Type  | Description                                   |
|------|-------|-----------------------------------------------|
| ids  | array | Unique identifiers for namespaces in CC array |
