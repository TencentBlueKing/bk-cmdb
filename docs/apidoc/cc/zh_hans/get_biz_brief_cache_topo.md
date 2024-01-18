### 功能描述

根据业务ID,查询该业务的全量简明拓扑树信息。（v3.9.14）
该业务拓扑的全量信息，包含了从业务这个根节点开始，到自定义层级实例(如果主线的拓扑层级中包含)，到集群、模块等中间的所有拓扑层级树数据。

注意： 
- 该接口为缓存接口，默认全量缓存刷新时间为15分钟。
- 如果业务的拓扑信息发生变化，会通过事件机制实时刷新缓存该业务的拓扑数据。

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选   |  描述                                                    |
|----------------------|------------|--------|--------------------------------------------------|
| bk_biz_id              | int     | 是     | 要查询的业务拓扑所属的业务的ID          |


### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 2
}
```

### 返回结果示例

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
        "nm": "空闲机池",
        "dft": 1,
        "nds": [
          {
            "obj": "module",
            "id": 7,
            "nm": "空闲机",
            "dft": 1,
            "nds": null
          },
          {
            "obj": "module",
            "id": 8,
            "nm": "故障机",
            "dft": 2,
            "nds": null
          },
          {
            "obj": "module",
            "id": 9,
            "nm": "待回收",
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

### 返回结果参数说明
#### response
| 名称    | 类型   | 说明                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                     |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                             |

#### data.biz 参数说明

| 字段         | 类型         | 描述     |
| ------------ | ------------ | -------- |
| id    | int          | 业务ID   |
| nm  | string       | 业务名   |
| dft | int | 业务类型，该值>=0，0: 表示该业务为普通业务。1: 表示该业务为资源池业务 |
| bk_supplier_account | string       | 开发商账号    |
#### data.idle 对象参数说明
idle对象中的数据表示该业务的空闲set中的数据，目前只有一个空闲set，后续可能有多个set，请勿依赖此数量。

| 字段         | 类型         | 描述     |
| ------------ | ------------ | -------- |
| obj    | string| 该资源的对象，可以是业务自定义层级对应的模块id(bk_obj_id字段值)，set, module等。|
| id    | int          | 该实例的ID   |
| nm  | string       | 该实例的名称  |
| dft  | int       | 该值>=0，只有set和module有该字段，0:表示普通的集群或者模块，>1:表示为空闲机类的set或module。  |
| nds  | object       | 该节点所属的子节点信息 |

#### data.nds 对象参数说明 
描述该业务下除空闲set外的其它拓扑节点的拓扑数据。该对象是一个数组对象，若无其它节点，则为空。
每个节点的对象描述如下，按照拓扑层级，各节点和其对应的子节点逐个嵌套。
需要注意的是，module的"nds"节点一定为空，module是整个业务拓扑树中最底层的节点。

| 字段         | 类型         | 描述     |
| ------------ | ------------ | -------- |
| obj    | string| 该资源的对象，可以是业务自定义层级对应的模块id(bk_obj_id字段值)，set, module等。|
| id    | int          | 该实例的ID   |
| nm  | string       | 该实例的名称  |
| dft  | int       | 该值>=0，只有set和module有该字段，0:表示普通的集群或者模块，>1:表示为空闲机类的set或module。  |
| nds  | object       | 该节点所属的子节点信息，按照拓扑层级逐级循环嵌套。 |

