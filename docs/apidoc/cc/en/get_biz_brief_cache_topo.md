### Function Description

This interface is used to query the full simplified topology tree information of a business based on the business ID. (v3.9.14) The full information of the business topology contains all the topology tree data from the root node of the business, to custom level instances (if included in the main topology level), to clusters, modules, and other intermediate topology levels.

Note:

- This interface is a cache interface, and the default full cache refresh time is 15 minutes.
- If the topology information of the business changes, the cache of the business topology data will be refreshed in real-time through the event mechanism.

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type | Required | Description                                                  |
| --------- | ---- | -------- | ------------------------------------------------------------ |
| bk_biz_id | int  | Yes      | ID of the business to which the business topology to be queried belongs |

### Request Parameter Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 2
}
```

### Response Example

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
        "nm": "Idle Host Pool",
        "dft": 1,
        "nds": [
          {
            "obj": "module",
            "id": 7,
            "nm": "Idle Host",
            "dft": 1,
            "nds": null
          },
          {
            "obj": "module",
            "id": 8,
            "nm": "Fault Host",
            "dft": 2,
            "nds": null
          },
          {
            "obj": "module",
            "id": 9,
            "nm": "To Be Recycled",
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

### Response Parameters Description

#### response

| Field       | Type   | Description                                                  |
| ---------- | ------ | ------------------------------------------------------------ |
| result     | bool   | Whether the request is successful. true: successful; false: failed |
| code       | int    | Error code. 0 represents success, >0 represents a failure error |
| message    | string | Error message returned in case of failure                    |
| permission | object | Permission information                                       |
| request_id | string | Request chain ID                                             |
| data       | object | Data returned by the request                                 |

#### Explanation of data.biz Parameters

| Field               | Type   | Description                                                  |
| ------------------- | ------ | ------------------------------------------------------------ |
| id                  | int    | Business ID                                                  |
| nm                  | string | Business name                                                |
| dft                 | int    | Business type, the value is >=0, 0: indicates that the business is a regular business. 1: indicates that the business is a resource pool business |
| bk_supplier_account | string | Supplier account                                             |

#### Explanation of data.idle Object Parameters

The data in the idle object represents the data in the idle set of the business. Currently, there is only one idle set, and there may be more sets in the future. Do not rely on this quantity.

| Field | Type   | Description                                                  |
| ----- | ------ | ------------------------------------------------------------ |
| obj   | string | Object of this resource, can be the module ID corresponding to the business custom level (value of the bk_obj_id field), set, module, etc. |
| id    | int    | ID of this instance                                          |
| nm    | string | Name of this instance                                        |
| dft   | int    | The value is >=0. Only set and module have this field. 0: indicates a regular cluster or module, >1: indicates a set or module for idle hosts. |
| nds   | object | Subnode information of this node                             |

#### Explanation of data.nds Object Parameters

Describes the topology data of other topology nodes except the idle set in the business. This object is an array object. If there are no other nodes, it is empty. Each node object is described as follows, nesting one by one with its corresponding child nodes. It should be noted that the "nds" node of the module must be empty. The module is the bottommost node in the entire business topology tree.

| Field | Type   | Description                                                  |
| ----- | ------ | ------------------------------------------------------------ |
| obj   | string | Object of this resource, can be the module ID corresponding to the business custom level (value of the bk_obj_id field), set, module, etc. |
| id    | int    | ID of this instance                                          |
| nm    | string | Name of this instance                                        |
| dft   | int    | The value is >=0. Only set and module have this field. 0: indicates a regular cluster or module, >1: indicates a set or module for idle hosts. |
| nds   | object | Subnode information of this node, nested one by one according to the topology level. |