### 描述

创建全量同步缓存条件(版本: v3.14.1+，权限: 全量同步缓存条件的创建权限)

### 输入参数

| 参数名称         | 参数类型   | 必选 | 描述                                                                                                                                                                                                                                                          |
|--------------|--------|----|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| resource     | string | 是  | 要进行全量同步数据缓存的资源类型，枚举值为：host, biz, set, module, process, object_instance, mainline_instance, biz_set, plat, project。其中host代表主机，biz代表业务，set代表集群，module代表模块，process代表进程，object_instance代表通用模型实例，mainline_instance代表主线模型实例，biz_set代表业务集，plat代表管控区域, project代表项目。 |
| sub_resource | string | 否  | 要全量同步的下级数据类型。resource为object_instance或mainline_instance时需要指定，代表需要同步的模型的bk_obj_id                                                                                                                                                                            |
| is_all       | bool   | 否  | 是否同步全量数据，每种资源仅允许存在一条is_all为true的条件                                                                                                                                                                                                                          |
| condition    | object | 否  | is_all为false时用于指定同步条件，组装规则可参考: https://github.com/TencentBlueKing/bk-cmdb/blob/master/pkg/filter/README.md                                                                                                                                                  |
| interval     | int    | 是  | 同步周期，单位为小时，用于指定缓存的过期时间，最短为6小时，最长为7天                                                                                                                                                                                                                         |

**注意：**

- 非同步全量数据的定制化条件的数量上限为100个
- 仅创建了对应的全量同步缓存条件的资源才会进行全量数据缓存

### 调用示例

```json
{
  "resource": "object_instance",
  "sub_resource": "bk_switch",
  "is_all": true,
  "interval": 24
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
    "id": 123
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

| 参数名称 | 参数类型 | 描述             |
|------|------|----------------|
| id   | int  | 创建的全量同步缓存条件的ID |
