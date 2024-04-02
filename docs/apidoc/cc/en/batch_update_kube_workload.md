### Function Description

Batch Update Workloads (Version: v3.12.1+, Permission: Edit Container Workload Permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description                                                  |
| --------- | ------ | -------- | ------------------------------------------------------------ |
| bk_biz_id | int    | Yes      | Business ID                                                  |
| kind      | string | Yes      | Workload type, currently supported workload types include deployment, daemonSet, statefulSet, gameStatefulSet, gameDeployment, cronJob, job, pods (directly create Pods without going through workload) |
| data      | object | Yes      | Fields to be updated                                         |
| ids       | array  | Yes      | Unique ID array in cc                                        |

#### data

| Field                   | Type   | Required | Description                                                  |
| ----------------------- | ------ | -------- | ------------------------------------------------------------ |
| labels                  | map    | No       | Labels                                                       |
| selector                | object | No       | Workload selector                                            |
| replicas                | int    | No       | Number of workload instances                                 |
| strategy_type           | string | No       | Workload update mechanism                                    |
| min_ready_seconds       | int    | No       | Specifies the minimum readiness time for newly created Pods, only Pods that exceed this time are considered available |
| rolling_update_strategy | object | No       | Rolling update strategy                                      |

#### selector

| Field             | Type  | Required | Description           |
| ----------------- | ----- | -------- | --------------------- |
| match_labels      | map   | No       | Match based on labels |
| match_expressions | array | No       | Match expressions     |

#### match_expressions[x]

| Field    | Type   | Required | Description                                                  |
| -------- | ------ | -------- | ------------------------------------------------------------ |
| key      | string | Yes      | Label key                                                    |
| operator | string | Yes      | Operator, optional values: "In", "NotIn", "Exists", "DoesNotExist" |
| values   | array  | No       | String array, must be provided if the operator is "In" or "NotIn", must be empty if the operator is "Exists" or "DoesNotExist" |

#### rolling_update_strategy

Only applicable when strategy_type is RollingUpdate, otherwise it is empty

| Field           | Type   | Required | Description         |
| --------------- | ------ | -------- | ------------------- |
| max_unavailable | object | No       | Maximum unavailable |
| max_surge       | object | No       | Maximum surge       |

#### max_unavailable

| Field   | Type   | Required | Description                                                  |
| ------- | ------ | -------- | ------------------------------------------------------------ |
| type    | int    | Yes      | Optional values are 0 (representing int type) or 1 (representing string type) |
| int_val | int    | No       | Must be provided if type is 0 (representing int type), corresponding int value |
| str_val | string | No       | Must be provided if type is 1 (representing string type), corresponding string value |

#### max_surge

| Field   | Type   | Required | Description                                                  |
| ------- | ------ | -------- | ------------------------------------------------------------ |
| type    | int    | Yes      | Optional values are 0 (representing int type) or 1 (representing string type) |
| int_val | int    | No       | Must be provided if type is 0 (representing int type), corresponding int value |
| str_val | string | No       | Must be provided if type is 1 (representing string type), corresponding string value |

Note: Use either the unique identifier of k8s or cc to pass in the association information. These two methods can only be used separately and cannot be mixed.

### Request Parameter Example

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 3,
  "kind": "deployment",
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
}
```

### Response Example

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

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error  |
| message    | string | Error message returned in case of request failure            |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Data returned in the request                                 |