### Functional description

create a new container pod and container(v3.10.23+，permission: kube pod creation permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| data    |  array  | yes     | Details of the pod to be created|

#### data

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_biz_id    |  int  | yes     | business ID|
| pods    |  array  | yes     | Details of the pod to be created under this business|

#### pods

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| spec   |  object  | yes   | pod association information|
| bk_host_id   |  int  | yes   | pod associated host id|
| name |  string  | yes     | pod name |
| priority  |  object  | no     | priority |
| labels  |  object  | no     | labels |
| ip  |  string  | no     | Container network IP|
| ips  |  array  | no     | Container network IP array|
| volumes  |  object  | no     | Volume information|
| qos_class  |  string  | no     | service quality|
| node_selectors  |  object  | no | Node label selector|
| tolerations  |  object  | no     | tolerance |
| containers  |  array  | no     | container information|

#### spec

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| bk_cluster_id |  int  | yes     | ID of the cluster where the pod is located|
| bk_namespace_id  |  int  | yes     | The ID of the namespace to which the pod belongs|
| bk_node_id  |  int  | yes     | ID of the node where the pod is located|
| ref  |  object  | yes     | Information about the workload corresponding to the pod|

#### ref

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| kind |  int  | yes     | the workload category associated with the pod. For specific categories, see Note|
| id  |  int  | yes     | the ID of the workload associated with the pod|

#### containers

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| name |  string  | yes     | container name|
| container_uid  |  string  | yes     | container ID|
| image  |  string  | no     | mirror information|
| ports  |  array  | no     | container port|
| host_ports  |  array  | no     | host port mapping|
| args  |  array  | no     | startup parameters|
| started  |  int  | no     | start time|
| limits  |  object  | no     | resource constraints|
| requests  |  object  | no     | application resource size|
| liveness  |  object  | no     | survival probe|
| environment  |  array  | no     | environment variable|
| mounts  |  array  | no     | mount volume|

#### ports

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| name |  string  | yes     | port name|
| hostPort  |  int  | no     | host port |
| containerPort  |  int  | no     | container port|
| protocol  |  string  | no     | protocol name|
| hostIP  |  string  | no     | host IP |

#### liveness

| Field      | Type      | Required   | Description      |
|-----------|------------|--------|------------|
| exec |  object  | yes     | perform action|
| httpGet  |  object  | no     | Http Get action |
| tcpSocket  |  object  | no     | tcp socket |
| grpc  |  object  | no     | grpc protocol |

**注意：**
- create no more than 200 pods at one time .
- specific workload category: deployment、statefulSet、daemonSet、gameStatefulSet、gameDeployment、cronJob、job、pods.
- this interface will create pods and corresponding containers synchronously.

### Request Parameters Example

```json
 {
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "bk_biz_id":2,
    "data":[
        {
            "bk_biz_id":1,
            "pods":[
                {
                    "spec":{
                        "bk_cluster_id":1,
                        "bk_namespace_id":1,
                        "ref":{
                            "kind":"deployment",
                            "id":1
                        },
                        "bk_node_id":1
                    },
                    "name":"name",
                    "bk_host_id":1,
                    "priority":1,
                    "labels":{
                        "env":"test"
                    },
                    "ip":"127.0.0.1",
                    "ips":[
                        {
                            "ip":"127.0.0.1"
                        },
                        {
                            "ip":"127.0.0.2"
                        }
                    ],
                    "containers":[
                        {
                            "name":"name",
                            "container_uid":"uid",
                            "image":"xxx",
                            "started":1
                        }
                    ]
                }
            ]
        }
    ]
}
```

### Return Result Example

```json
{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data":{
        "ids":[
            1,
            2
        ]
    },
    "request_id":"87de106ab55549bfbcc46e47ecf5bcc7"
}
```
**注意：**
- the order of the node ID array in the returned data is consistent with the order of the array data in the parameter.


### Return Result Parameters Description
#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request succeeded or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| permission    |  object |Permission information    |
| data    |  object |Data returned by request                           |
| request_id    |  string |Request chain id    |

### data

| Name    | Type   | Description                   |
| ------- | ------ | ------------------------------- |
| ids  | array   |list of kube pod IDs created |
