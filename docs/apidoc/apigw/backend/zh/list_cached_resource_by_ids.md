### 描述

根据ID列表拉取缓存的资源详情(版本: v3.14.1+，权限: 通用缓存查询权限)

### 输入参数

| 参数名称         | 参数类型         | 必选 | 描述                                                                                                                                                                                                                                                   |
|--------------|--------------|----|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| resource     | string       | 是  | 需要查询的资源类型，枚举值为：host, biz, set, module, process, object_instance, mainline_instance, biz_set, plat, project。其中host代表主机，biz代表业务，set代表集群，module代表模块，process代表进程，object_instance代表通用模型实例，mainline_instance代表主线模型实例，biz_set代表业务集，plat代表管控区域, project代表项目。 |
| sub_resource | string       | 否  | 需要查询的下级数据类型。resource为object_instance或mainline_instance时需要指定，代表需要同步的模型的bk_obj_id                                                                                                                                                                      |
| ids          | int array    | 是  | 需要查询的ID列表，最多500个                                                                                                                                                                                                                                     |
| fields       | string array | 否  | 属性列表，控制返回结果的资源详情里有哪些字段                                                                                                                                                                                                                               |

### 调用示例

```json
{
  "resource": "object_instance",
  "sub_resource": "bk_switch",
  "ids": [
    123,
    456
  ],
  "fields": [
    "bk_asset_id",
    "bk_inst_id",
    "bk_inst_name",
    "bk_obj_id"
  ]
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "info": [
      {
        "bk_asset_id": "sw00001",
        "bk_inst_id": 1,
        "bk_inst_name": "sw1",
        "bk_obj_id": "bk_switch"
      },
      {
        "bk_asset_id": "sw00002",
        "bk_inst_id": 2,
        "bk_inst_name": "sw2",
        "bk_obj_id": "bk_switch"
      }
    ]
  }
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |

#### data

| 参数名称 | 参数类型  | 描述     |
|------|-------|--------|
| info | array | 缓存详情数据 |

#### data.info

| 参数名称         | 参数类型   | 描述   |
|--------------|--------|------|
| bk_asset_id  | string | 固资编号 |
| bk_inst_id   | int    | 实例ID |
| bk_inst_name | string | 实例名  |
| bk_obj_id    | string | 模型ID |

**注意：此处的返回值仅以拉取交换机的部分字段的场景为例对其属性字段进行了说明，具体返回值取决于资源类型和用户自定义的属性字段**
