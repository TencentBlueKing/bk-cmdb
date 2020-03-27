# Mainline core service

## What is it ?
提供对模型主线拓扑和业务实例主线拓扑操作接口

## Why we need it ?
对接权限中心时，注册主机资源需要查询主机所在的拓扑主线（权限同步和host server两个模块需要用到），但是现有的拓扑信息实现在场景层，场景层之间相互调用不符合CMDB系统架构
现有的实现方案复杂度高，需要`o(n)`次 mongo 查询，不适合用于基础模块

另外，从系统设计的角度触发，场景层之间不应该相互调用，但是之前的主线拓扑的操作封装在topo_server中，而注册权限中心时又需要拿到实例的所有拓扑层级（比如主机实例的注册及权限信息同步模块都需要拓扑数据）。

## 实现方案

### 模型拓扑查询
模型拓扑属于全局数据，不依赖业务信息存在，主线模型及关联关系可直接通过模型关联关系表 `cc_ObjAsst` 拿到，最终的拓扑模型可以用一个链表表示。

### 业务下的实例拓扑数据
业务下的实例拓扑数据，及具体业务下的主线模型实例组成的拓扑层级，它的实现相对模型拓扑查询更复杂些。由于实例数据只记录了父节点的ID（ID可能存在于多张表结构中），并不记录父节点模型类型，因此，拿到实例父节点信息需要首先根据实例本身的模型类型查找出其父实例的模型类型，然后根据父实例所属的模型类型找出找出父节点ID指向的表结构，最终查询出父实例的详情。根据实例的模型类型查询其父实例模型类型可根据模型拓扑查询实现。


## Features
- 查询模型拓扑
- 查询业务的实例拓扑

## TODO
- 查询模型拓扑支持 `withDetail` 选项

## 输出demo

- 查询模型拓扑

```json
{
  "result": true,
  "bk_error_code": 0,
  "bk_error_msg": "success",
  "data": {
    "Children": [
      {
        "Children": [
          {
            "Children": [
              {
                "Children": [
                  {
                    "Children": [],
                    "ObjectID": "host"
                  }
                ],
                "ObjectID": "module"
              }
            ],
            "ObjectID": "set"
          }
        ],
        "ObjectID": "mainlinelevel1"
      }
    ],
    "ObjectID": "biz"
  }
}
```

- 业务实例拓扑

```json
{
  "result": true,
  "bk_error_code": 0,
  "bk_error_msg": "success",
  "data": {
    "Children": [
      {
        "Children": [
          {
            "Children": [
              {
                "Children": [],
                "ObjectID": "module",
                "InstanceID": 1,
                "Detail": {}
              },
              {
                "Children": [],
                "ObjectID": "module",
                "InstanceID": 2,
                "Detail": {}
              }
            ],
            "ObjectID": "set",
            "InstanceID": 1,
            "Detail": {}
          }
        ],
        "ObjectID": "mainlinelevel1",
        "InstanceID": 1,
        "Detail": {}
      }
    ],
    "ObjectID": "biz",
    "InstanceID": 1,
    "Detail": {}
  }
}
```
