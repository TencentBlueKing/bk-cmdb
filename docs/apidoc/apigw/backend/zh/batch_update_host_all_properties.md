### 描述

根据主机id和属性批量更新主机属性（版本：v3.13.6+，权限：业务主机编辑权限）

### 输入参数

| 参数名称   | 参数类型  | 必选 | 描述                      |
|--------|-------|----|-------------------------|
| update | array | 是  | 主机被更新的属性和值，最多同时更新500台主机 |

#### update

| 参数名称        | 参数类型   | 必选 | 描述                           |
|-------------|--------|----|------------------------------|
| properties  | object | 是  | 主机被更新的属性和值 |
| bk_host_ids | array  | 是  | 用于更新的主机ID    |

#### properties

| 参数名称         | 参数类型   | 必选 | 描述                                  |
|--------------|--------|----|-------------------------------------|
| bk_host_name | string | 否  | 主机名，也可以为其它属性                        |
| operator     | string | 否  | 主要维护人，也可以为其它属性                      |
| bk_comment   | string | 否  | 备注，也可以为其它属性                         |
| bk_isp_name  | string | 否  | 所属运营商，也可以为其它属性                      |
| bk_cloud_id  | int    | 否  | 管控区域id，只能更新管控区域为:“未分配[90000001]”的主机 |

### 调用示例

```json
{
    "update":[
      {
        "properties":{
          "bk_host_name":"batch_update",
          "operator": "admin",
          "bk_comment": "test",
          "bk_isp_name": "1",
          "bk_cloud_id": 0
        },
        "bk_host_ids":[46]
      }
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
    "data": null
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
