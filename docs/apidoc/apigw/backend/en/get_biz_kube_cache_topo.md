### Description

Query the full kube topology tree information of the business based on the biz ID (version: v3.12.5+, permission: biz
access)
The full information of the business kube topology tree includes all topology hierarchy tree data starting from the root
node of the business to the Cluster, Namespace, and Workload levels.

Note:

- This interface is a cache interface, and the default full cache refresh time is 15 minutes.
- If the kube topology information of the biz changes, the kube topology data cache of the biz will be refreshed in real
  time through the event mechanism.

### Parameters

| Name      | Type | Required | Description                                                 |
|-----------|------|----------|-------------------------------------------------------------|
| bk_biz_id | int  | yes      | The business ID of the business kube topology to be queried |

### Request Example

```json
{
  "bk_biz_id": 3
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "biz": {
      "id": 3,
      "nm": "biz1",
      "cnt": 100
    },
    "nds": [
      {
        "kind": "cluster",
        "id": 22,
        "nm": "cluster1",
        "cnt": 100,
        "nds": [
          {
            "kind": "namespace",
            "id": 16,
            "nm": "namespace1",
            "cnt": 100,
            "nds": [
              {
                "kind": "deployment",
                "id": 48,
                "nm": "deployment1",
                "cnt": 11,
                "nds": null
              },
              {
                "kind": "daemonSet",
                "id": 49,
                "nm": "daemonSet1",
                "cnt": 89,
                "nds": null
              }
            ]
          }
        ]
      }
    ]
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                                               |
|------------|--------|-------------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. true:request successful; false request failed. |
| code       | int    | The error code. 0 means success, >0 means failure error.                                  |
| message    | string | The error message returned by the failed request.                                         |
| data       | object | The data returned by the request.                                                         |
| permission | object | Permission information                                                                    |

#### data

| Name | Type         | Description                       |
|------|--------------|-----------------------------------|
| biz  | object       | The business information          |
| nds  | object array | Topology node data under business |

#### data.biz

| Name | Type   | Description                             |
|------|--------|-----------------------------------------|
| id   | int    | Service ID                              |
| nm   | string | Business name                           |
| cnt  | int    | Number of Containers under the business |

#### data.nds[n]

| Name | Type   | Description                                                                                                                                                                           |
|------|--------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| kind | string | The resource kind of the topo node. Supports: cluster, namespace, deployment, daemonSet, statefulSet, gameStatefulSet, gameDeployment, cronJob, job, pods.                            |
| id   | int    | The ID of the topo node.                                                                                                                                                              |
| nm   | string | The name of the topo node.                                                                                                                                                            |
| cnt  | int    | The number of Containers under the topo node                                                                                                                                          |
| nds  | object | The child topo node info under the topo node. Topo node info is circularly nested step by step according to the topological level. The value will be empty if there is no child node. |