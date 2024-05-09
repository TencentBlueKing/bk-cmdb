### Description

Batch delete workloads (Version: v3.12.1+, Permission: Container workload deletion)

### Parameters

| Name      | Type   | Required | Description                                                                                                                                                                                    |
|-----------|--------|----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id | int    | Yes      | Business ID                                                                                                                                                                                    |
| kind      | string | Yes      | Workload type. Currently supported workload types are deployment, daemonSet, statefulSet, gameStatefulSet, gameDeployment, cronJob, job, pods (creating pods directly without using workloads) |
| ids       | array  | Yes      | Array of unique identifiers of workloads in CC, with a limit of 200 at a time                                                                                                                  |

### Request Example

```json
{
  "bk_biz_id": 3,
  "kind": "deployment",
  "ids": [
    1
  ]
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
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Data returned in the request                                                |
