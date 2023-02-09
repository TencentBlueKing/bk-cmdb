### Functional description

According to the service ID, querying the full-volume concise topology tree information of the service.（v3.9.14）
The full information of the service topology includes all topology hierarchy tree data from the root node of the service, to the user-defined hierarchy instance (if included in the topology hierarchy of the main line), to the middle of the  set , module, etc.

Note:
- This interface is a cache interface and the default full cache refresh time is 15 minutes.
- If the topology information of the service changes, the topology data of the service will be refreshed and cached in real time through the event mechanism.

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                 | Type      | Required   | Description                                                    |
|----------------------|------------|--------|--------------------------------------------------|
| bk_biz_id              |  int     | yes  | ID of the business to which the business topology to query belongs |


### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 2
}
```

### Return Result Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807",
  "data": {
    "biz": {
      "id": 3,
      "nm": "lee",
      "dft": 0,
      "bk_supplier_account": "0"
    },
    "idle": [
      {
        "obj": "set",
        "id": 3,
        "nm": "Idle pool",
        "dft": 1,
        "nds": [
          {
            "obj": "module",
            "id": 7,
            "nm": "Idle machine",
            "dft": 1,
            "nds": null
          },
          {
            "obj": "module",
            "id": 8,
            "nm": "Faulty machine",
            "dft": 2,
            "nds": null
          },
          {
            "obj": "module",
            "id": 9,
            "nm": "To be recycled",
            "dft": 3,
            "nds": null
          }
        ]
      }
    ],
    "nds": [
      {
        "obj": "province",
        "id": 22,
        "nm": "Guangdong",
        "nds": [
          {
            "obj": "set",
            "id": 16,
            "nm": "magic-set",
            "dft": 0,
            "nds": [
              {
                "obj": "module",
                "id": 48,
                "nm": "gameserver",
                "dft": 0,
                "nds": null
              },
              {
                "obj": "module",
                "id": 49,
                "nm": "mysql",
                "dft": 0,
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

### Return Result Parameters Description
#### response
| Name    | Type   | Description                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                     |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                             |

#### Data.biz Parameter Description

| Field         | Type         | Description     |
| ------------ | ------------ | -------- |
| id    |  int          | Service ID   |
| nm  | string       | Business name   |
| dft | int |Business type, the value>= 0,0: Indicates that the service is an ordinary service. 1: indicates that the service is a resource pool service|
| bk_supplier_account | string       | Developer account number    |
#### Data.idle object parameter Description
The data in the idle object indicates the data in the idle set of the service. There is only one idle set at present, and there may be multiple sets in the future. Please do not rely on this quantity.

| Field         | Type         | Description     |
| ------------ | ------------ | -------- |
| obj    |  string| The object of the resource may be the module id(bk_obj_id field value), set, module, etc. Corresponding to the business user-defined level.|
| id    |  int          | The ID of the instance   |
| nm  | string       | The name of the instance|
| dft  | int       | This value>=0. Only set and module have this field. 0: Represents a common  set  or module,>1: set or module represented as an idle machine class.  |
| nds  | object       | Child node information to which this node belongs|

#### Data.NDS object parameter Description
Describe the topology data of other topology nodes except idle set under the service. The object is an array object and is empty if there are no other nodes.
The object for each node is described below, and each node and its corresponding child nodes are nested one by one according to the topology level.
It should be noted that the nds node of module must be empty, and module is the lowest node in the whole service topology tree.

| Field         | Type         | Description     |
| ------------ | ------------ | -------- |
| obj    |  string| The object of the resource may be the module id(bk_obj_id field value), set, module, etc. Corresponding to the business user-defined level.|
| id    |  int          | The ID of the instance   |
| nm  | string       | The name of the instance|
| dft  | int       | This value>=0. Only set and module have this field. 0: Represents a common set or module,>1: set or module represented as an idle machine class.  |
| nds  | object       | The information of the child nodes to which the node belongs is circularly nested step by step according to the topology level. |

