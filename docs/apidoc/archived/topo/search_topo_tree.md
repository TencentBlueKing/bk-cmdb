### 功能描述

根据业务名、自定义层级名、集群、模块名模糊搜索业务的拓扑树。

主机规则包括：
- 在多业务搜索下，如果扫描的业务数量超过一定数量，可以直接拒绝搜索，直接返回，报查询数据过多的错误，暂定为20个业务。
- 可以支持业务名、自定义层级名、集群、模块组合的查询方式进行模糊搜索。
- 除主机、模块外，其它的这些可以提供当前层级或者下一层级的拓扑信息，但是总的搜索的节点数量不能超过50个。如果超过50个，则直接拒绝搜索，报查询数据过多的错误。 

注意： 
- 该接口仅用于web页面搜索使用，不建议后台使用。
- 该接口有5分钟的缓存，在最极端的情况下，5分钟内的数据可能会不对。
- 每次搜索会自动触发缓存更新，所以在数据不对的情况下，第二次搜索数据即准确。

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选   |  描述                                                    |
|----------------------|------------|--------|--------------------------------------------------|
| bk_biz_id              | int     | 是     | 业务ID，其中-1代表所有业务，不能为0。正整数值代表对应的业务ID           |
| bk_biz_name              | string     | 是     | 业务名，业务名和业务ID可以选择只用一个，但不能同时为空。如果bk_biz_id > 0, 则bk_biz_name会被忽略,因为已经指定了具体的业务。      |
| bk_set_name              | string     | 否     | 集群名   |
| bk_module_name              | string     | 否     | 模块名   |
| bk_level              | object     | 否     | 自定义层级描述   |

注意： bk_set_name, bk_module_name, bk_level只能使用期中一个。 

- bk_level

| 字段                 |  类型      | 必选   |  描述                       |
|----------------------|------------|--------|--------------------------|
| bk_obj_id            | string     | 是     | 自定义层级的模型名称         |
| bk_inst_name         | string     | 是     | 自定义层级的实例名称         |



### 请求参数示例

```json
{
    "bk_biz_id": 2,
    "bk_biz_name": "",
    "bk_set_name": "paas",
    "bk_module_name": "",
    "bk_level": {
        "bk_obj_id": "",
        "bk_inst_name": ""
    }
}
```

### 返回结果示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": [
    {
      "bk_biz_id": 2,
      "bk_biz_name": "蓝鲸",
      "bk_topo_tree": [
        {
          "bk_obj_id": "set",
          "bk_inst_name": "PaaS平台",
          "bk_inst_id": 4,
          "children": [
            {
              "bk_obj_id": "module",
              "bk_inst_name": "nginx",
              "bk_inst_id": 15,
              "children": null
            },
            {
              "bk_obj_id": "module",
              "bk_inst_name": "elasticsearch",
              "bk_inst_id": 14,
              "children": null
            },
            {
              "bk_obj_id": "module",
              "bk_inst_name": "nfs",
              "bk_inst_id": 5000083,
              "children": null
            }
          ]
        }
      ]
    }
  ]
}
```

### 返回结果参数说明

#### data array

| 字段         | 类型         | 描述     |
| ------------ | ------------ | -------- |
| bk_biz_id    | int          | 业务ID   |
| bk_biz_name  | string       | 业务名   |
| bk_topo_tree | Object Array | 根据查询条件生成的拓扑树层级信息，该树包含查询对象到业务根节点的所有层级信息，及下一层级信息，但模块例外，没有下一层级。 |

#### bk_obj_tree array

| 字段         | 类型         | 描述                                                     |
| ------------ | ------------ | -------------------------------------------------------- |
| bk_obj_id    | string       | 对象类型，可以为自定义层级对名称，set,module名称         |
| bk_inst_name | string       | 对象名，如集群名、模块名等                               |
| bk_inst_id   | Int          | 对象实例身份id，唯一标识这个实例ID                       |
| children     | object array | 表示该对象实例的子节点信息，可以是自定义层级集群、模块等 |

#### children array 

| 字段         | 类型   | 描述                                             |
| ------------ | ------ | ------------------------------------------------ |
| bk_obj_id    | string | 对象类型，可以为自定义层级对名称，set,module名称 |
| bk_inst_name | string | 对象名，如集群名、模块名等                       |
| bk_inst_id   | Int    | 对象实例身份id，唯一标识这个实例ID               |

注意： 模块没有子层级，即没有children.