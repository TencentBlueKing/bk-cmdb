### Description

Create Containers, Pods, and Containers (v3.12.1+, Permission: Container pods creation permission)

### Parameters

| Name | Type  | Required | Description                            |
|------|-------|----------|----------------------------------------|
| data | array | Yes      | Detailed information for creating pods |

#### data[x]

| Name      | Type  | Required | Description                                             |
|-----------|-------|----------|---------------------------------------------------------|
| bk_biz_id | int   | Yes      | Business ID                                             |
| pods      | array | Yes      | Detailed information for creating pods in this business |

#### pods[x]

| Name           | Type         | Required | Description                    |
|----------------|--------------|----------|--------------------------------|
| spec           | object       | Yes      | Associated pod information     |
| bk_host_id     | int          | Yes      | Associated host ID             |
| name           | string       | Yes      | Pod name                       |
| operator       | string array | Yes      | Person in charge of the pod    |
| priority       | object       | No       | Priority                       |
| labels         | object       | No       | Labels                         |
| ip             | string       | No       | Container network IP           |
| ips            | array        | No       | Array of container network IPs |
| volumes        | object       | No       | Volume information             |
| qos_class      | string       | No       | Quality of service             |
| node_selectors | object       | No       | Node label selector            |
| tolerations    | object       | No       | Tolerations                    |
| containers     | array        | No       | Container information          |

#### spec

| Name            | Type   | Required | Description                                  |
|-----------------|--------|----------|----------------------------------------------|
| bk_cluster_id   | int    | Yes      | ID of the cluster where the pod is located   |
| bk_namespace_id | int    | Yes      | ID of the namespace to which the pod belongs |
| bk_node_id      | int    | Yes      | ID of the node where the pod is located      |
| ref             | object | Yes      | Relevant information about the pod           |
| bk_pod_id       | int    | No       | ID of the pod (optional)                     |

#### ref

| Name | Type | Required | Description                                                                    |
|------|------|----------|--------------------------------------------------------------------------------|
| kind | int  | Yes      | Category of the workload related to the pod, see notes for specific categories |
| id   | int  | Yes      | ID of the workload related to the pod                                          |

#### containers[x]

| Name          | Type   | Required | Description             |
|---------------|--------|----------|-------------------------|
| name          | string | Yes      | Container name          |
| container_uid | string | Yes      | Container ID            |
| image         | string | No       | Image information       |
| ports         | array  | No       | Container ports         |
| host_ports    | array  | No       | Host port mapping       |
| args          | array  | No       | Startup parameters      |
| started       | int    | No       | Startup time            |
| limits        | object | No       | Resource limits         |
| requests      | object | No       | Requested resource size |
| liveness      | object | No       | Liveness probe          |
| environment   | array  | No       | Environment variables   |
| mounts        | array  | No       | Mounted volumes         |

#### ports[x]

| Name          | Type   | Required | Description    |
|---------------|--------|----------|----------------|
| name          | string | Yes      | Port name      |
| hostPort      | int    | No       | Host port      |
| containerPort | int    | No       | Container port |
| protocol      | string | No       | Protocol name  |
| hostIP        | string | No       | Host IP        |

#### liveness

| Name      | Type   | Required | Description      |
|-----------|--------|----------|------------------|
| exec      | object | Yes      | Execution action |
| httpGet   | object | No       | Http Get action  |
| tcpSocket | object | No       | tcp socket       |
| grpc      | object | No       | grpc protocol    |

**Note:**

- The number of pods created at once should not exceed 200.
- Specific workload categories: deployment, statefulSet, daemonSet, gameStatefulSet, gameDeployment, cronJob, job, pods.
- This interface will synchronously create pods and their corresponding containers.

### Request Example

```json
{
  "data": [
    {
      "bk_biz_id": 1,
      "pods": [
        {
          "spec": {
            "bk_cluster_id": 1,
            "bk_namespace_id": 1,
            "ref": {
              "kind": "deployment",
              "id": 1
            },
            "bk_node_id": 1
          },
          "name": "name",
          "operator": [
            "user1",
            "user2"
          ],
          "bk_host_id": 1,
          "priority": 1,
          "labels": {
            "env": "test"
          },
          "ip": "127.0.0.1",
          "ips": [
            {
              "ip": "127.0.0.1"
            },
            {
              "ip": "127.0.0.2"
            }
          ],
          "containers": [
            {
              "name": "name",
              "container_uid": "uid",
              "image": "xxx",
              "started": 1
            }
          ]
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
  "message": "success",
  "permission": null,
  "data": {
    "ids": [
      1,
      2
    ]
  },
}
```

**Note:**

- The order of the pod ID array returned in the data field corresponds to the order of the array data in the parameters.

### Response Parameters

| Name       | Type   | Description                                                       |
|------------|--------|-------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error     |
| message    | string | Error message returned for a failed request                       |
| permission | object | Permission information                                            |
| data       | object | Data returned by the request                                      |
