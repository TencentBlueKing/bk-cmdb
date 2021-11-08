### Functional description

get a business's brief topology data, which contains each nodes's brief information from biz to the bottom modules. （v3.9.14）

Note： 
- this is a cache api, which has a default 15Min ttl。
- the cache will also be refreshed by event which is trigged by add mainline topology instance, such as add a new set, module etc.

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field                |  Type       | Required	   | Description                            |
|----------------------|------------|--------|--------------------------------------------------|
| bk_biz_id              | int     | Yes    | the business's id to be searched          |


### Request Parameters Example

```json
{
    "bk_biz_id": 2
}
```

### Return Result Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "data": {
    "biz": {
      "id": 3,
      "nm": "lee",
      "dft": 0
    },
    "idle": [
      {
        "obj": "set",
        "id": 3,
        "nm": "idle set",
        "dft": 1,
        "nds": [
          {
            "obj": "module",
            "id": 7,
            "nm": "idle module",
            "dft": 1,
            "nds": null
          },
          {
            "obj": "module",
            "id": 8,
            "nm": "fault module",
            "dft": 2,
            "nds": null
          },
          {
            "obj": "module",
            "id": 9,
            "nm": "recycle module",
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
        "nm": "广东",
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


#### data.biz object description

| Field                 |  Type    	   |  Description       |
| ------------ | ------------ | -------- |
| id    | int          | business id   |
| nm  | string       | business name  |
| dft | int | the type of business，dft >=0，0: a common business. 1: a resource pool business |

#### data.idle object description

idle describe the special idle set's topology nodess. for now, only have a idle set. it may have more in the future.

| Field                 |  Type    	   |  Description       |
| ------------ | ------------ | -------- |
| obj    | string| describe which kind of resource is. such as module, set and custom level object's id, as is a model's bk_object_id value.   |
| id    | int          | instance's identity id   |
| nm  | string       | instance's name |
| dft  | int       | dft>=0，only set or module have this field. 0: means this is a common set or module，>1: a special set or module, such as idle set, fault module etc.  |
| nds  | object       | this instance's all the sub nodes. |

#### data.nds description
nds describe all the topology instance tree, except the former idle set's data. nds may be null if this nodes has no sub-resources.

| Field                 |  Type    	   |  Description       |
| ------------ | ------------ | -------- |
| obj    | string| describe which kind of resource is. such as module, set and custom level object's id, as is a model's bk_object_id value.   |
| id    | int          | instance's identity id   |
| nm  | string       | instance's name |
| dft  | int       | dft>=0，only set or module have this field. 0: means this is a common set or module，>1: a special set or module, such as idle set, fault module etc.  |
| nds  | object       | this instance's all the sub nodes. |

