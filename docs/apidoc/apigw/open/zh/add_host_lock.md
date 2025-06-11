### 描述

根据主机的id列表对主机加锁，新加主机锁，如果主机已经加过锁，同样提示加锁成功(版本：v3.8.6，权限：业务主机编辑权限)

### 输入参数

| 参数名称    | 参数类型      | 必选 | 描述     |
|---------|-----------|----|--------|
| id_list | int array | 是  | 主机ID列表 |

### 调用示例

```json
{
   "id_list":[1, 2, 3]
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": null,
    "permission": null,
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| data       | object | 请求返回的数据                    |
| permission | object | 权限信息                       |
