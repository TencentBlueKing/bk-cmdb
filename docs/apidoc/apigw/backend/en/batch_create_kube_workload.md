### Description

Batch Create Workloads (Version: v3.12.1+, Permission: Container workloads creation permission)

### Parameters

| Name      | Type   | Required | Description                                                                                                                                                                                                |
|-----------|--------|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id | int    | Yes      | Business ID                                                                                                                                                                                                |
| kind      | string | Yes      | Workload type, currently supported workload types include deployment, daemonSet, statefulSet, gameStatefulSet, gameDeployment, cronJob, job, pods (create Pods directly without passing through workloads) |
| data      | array  | Yes      | Array, limited to creating 200 at a time                                                                                                                                                                   |

#### data[x]

| Name                    | Type   | Required | Description                                                                                                                                              |
|-------------------------|--------|----------|----------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_namespace_id         | int    | Yes      | The unique identifier of the namespace in cc                                                                                                             |
| name                    | string | Yes      | Workload name                                                                                                                                            |
| labels                  | map    | No       | Labels                                                                                                                                                   |
| selector                | object | No       | Workload selector                                                                                                                                        |
| replicas                | int    | No       | Number of workload instances                                                                                                                             |
| strategy_type           | string | No       | Workload update mechanism                                                                                                                                |
| min_ready_seconds       | int    | No       | Specifies the minimum ready time for a newly created Pod without any container crashes. Only when this time is exceeded, the Pod is considered available |
| rolling_update_strategy | object | No       | Rolling update strategy                                                                                                                                  |

#### selector

| Name              | Type  | Required | Description          |
|-------------------|-------|----------|----------------------|
| match_labels      | map   | No       | Match labels         |
| match_expressions | array | No       | Matching expressions |

#### match_expressions[x]

| Name     | Type   | Required | Description |
|----------|--------|----------|-------------|
| key      | string | Yes      | Label key   |
| operator | string | Yes      | Operator    |
| values   | array  | No       | Values      |

#### rolling_update_strategy

When strategy_type is RollingUpdate, it cannot be empty. Otherwise, it is empty.

| Name            | Type   | Required | Description         |
|-----------------|--------|----------|---------------------|
| max_unavailable | object | No       | Maximum unavailable |
| max_surge       | object | No       | Maximum surge       |

#### max_unavailable

| Name    | Type   | Required | Description                    |
|---------|--------|----------|--------------------------------|
| type    | int    | Yes      | Type (0 for int, 1 for string) |
| int_val | int    | No       | Integer value                  |
| str_val | string | No       | String value                   |

#### max_surge

| Name    | Type   | Required | Description                    |
|---------|--------|----------|--------------------------------|
| type    | int    | Yes      | Type (0 for int, 1 for string) |
| int_val | int    | No       | Integer value                  |
| str_val | string | No       | String value                   |

### Request Example

```json
{
  "bk_biz_id": 3,
  "kind": "deployment",
  "data": [
    {
      "bk_namespace_id": 1,
      "name": "test",
      "labels": {
        "test": "test",
        "test2": "test2"
      },
      "selector": {
        "match_labels": {
          "test": "test",
          "test2": "test2"
        },
        "match_expressions": [
          {
            "key": "tier",
            "operator": "In",
            "values": [
              "cache"
            ]
          }
        ]
      },
      "replicas": 1,
      "strategy_type": "RollingUpdate",
      "min_ready_seconds": 1,
      "rolling_update_strategy": {
        "max_unavailable": {
          "type": 0,
          "int_val": 1
        },
        "max_surge": {
          "type": 0,
          "int_val": 1
        }
      }
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

- The order of the workload ID array returned in the data field corresponds to the order of the array data in the
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

| Name | Type  | Description                       |
|------|-------|-----------------------------------|
| ids  | array | Array of unique identifiers in cc |
