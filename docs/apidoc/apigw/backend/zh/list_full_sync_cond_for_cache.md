### 描述

查询全量同步缓存条件(版本: v3.14.1+，权限: 全量同步缓存条件的查询权限)

### 输入参数

| 参数名称         | 参数类型      | 必选 | 描述                                      |
|--------------|-----------|----|-----------------------------------------|
| resource     | string    | 否  | 需要查询的全量同步数据缓存的资源类型。resource和ids至少需要选择一个 |
| sub_resource | string    | 否  | 需要查询的全量同步的下级数据类型                        |
| ids          | int array | 否  | 需要查询的全量同步缓存条件的ID列表。resource和ids至少需要选择一个 |

### 调用示例

```json
{
  "ids": [
    123,
    456
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
        "id": 123,
        "resource": "object_instance",
        "sub_resource": "bk_switch",
        "is_all": true,
        "interval": 24
      },
      {
        "id": 456,
        "resource": "host",
        "is_all": false,
        "interval": 6,
        "condition": {
          "condition": "AND",
          "rules": [
            {
              "field": "bk_host_innerip",
              "operator": "not_equal",
              "value": "127.0.0.1"
            },
            {
              "condition": "OR",
              "rules": [
                {
                  "field": "bk_os_type",
                  "operator": "in",
                  "value": [
                    "3"
                  ]
                },
                {
                  "field": "bk_cloud_id",
                  "operator": "equal",
                  "value": 0
                }
              ]
            }
          ]
        }
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

| 参数名称 | 参数类型  | 描述           |
|------|-------|--------------|
| info | array | 全量同步缓存条件数据列表 |

#### info[x]

| 参数名称         | 参数类型   | 描述                                                                                                                                                                                                                                                       |
|--------------|--------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id           | int    | 全量同步缓存条件的自增ID                                                                                                                                                                                                                                            |
| resource     | string | 全量同步数据缓存的资源类型，枚举值为：host, biz, set, module, process, object_instance, mainline_instance, biz_set, plat, project。其中host代表主机，biz代表业务，set代表集群，module代表模块，process代表进程，object_instance代表通用模型实例，mainline_instance代表主线模型实例，biz_set代表业务集，plat代表管控区域, project代表项目。 |
| sub_resource | string | 全量同步数据缓存的下级数据类型。resource为object_instance或mainline_instance时代表需要同步的模型的bk_obj_id                                                                                                                                                                           |
| is_all       | bool   | 是否同步全量数据                                                                                                                                                                                                                                                 |
| condition    | object | is_all为false时用于指定同步条件，组装规则可参考: https://github.com/TencentBlueKing/bk-cmdb/blob/master/pkg/filter/README.md                                                                                                                                               |
| interval     | int    | 同步周期，单位为小时                                                                                                                                                                                                                                               |
