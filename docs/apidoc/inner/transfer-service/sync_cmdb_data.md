### 请求方式

POST /transfer/v3/sync/cmdb/data

### 描述

触发同步cmdb数据任务

注意：

- 该接口仅用于在源环境触发同步指定的cmdb数据任务
- 该接口是异步接口，触发后就会返回，可以通过rid追踪触发的同步任务
- 同一时间只允许一个同步任务

### 输入参数

| 参数名称          | 参数类型             | 必选 | 描述                                                                                                 |
|---------------|------------------|----|----------------------------------------------------------------------------------------------------|
| resource_type | string           | 是  | 同步数据的资源类型，枚举值：biz,set,module,host,host_relation,object_instance,inst_asst,service_instance,process |
| sub_resource  | string           | 否  | 下级数据类型。resource为object_instance和inst_asst时代表需要同步的模型的bk_obj_id                                      |
| is_all        | bool             | 否  | 是否同步全量数据，is_all为false时必须传start或end                                                                 |
| start         | map[string]int64 | 否  | 同步的起始区间ID信息，同步的数据不包含该ID代表的数据                                                                       |
| end           | map[string]int64 | 否  | 同步的结束区间ID信息，同步的数据包含该ID代表的数据                                                                        |

### 调用示例

```json
{
  "resource_type": "object_instance",
  "sub_resource": "bk_switch",
  "is_all": false,
  "start": {
    "bk_inst_id": 10
  },
  "end": {
    "bk_inst_id": 100
  }
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "",
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
