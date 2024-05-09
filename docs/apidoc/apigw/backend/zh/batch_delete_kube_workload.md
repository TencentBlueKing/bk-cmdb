### 描述

批量删除workload (版本：v3.12.1+，权限：容器工作负载删除)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述                                                                                                                                |
|-----------|--------|----|-----------------------------------------------------------------------------------------------------------------------------------|
| bk_biz_id | int    | 是  | 业务id                                                                                                                              |
| kind      | string | 是  | workload类型，目前支持的workload类型有deployment、daemonSet、statefulSet、gameStatefulSet、gameDeployment、cronJob、job、pods(放不通过workload而直接创建Pod) |
| ids       | array  | 是  | 要删除的workload在cc中的id唯一标识数组, 一次限制大小为200                                                                                             |

### 调用示例

```json
{
  "bk_biz_id": 3,
  "kind": "deployment",
  "ids": [
    1
  ]
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "data": null,
  "message": "success",
  "permission": null,
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
