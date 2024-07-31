### 描述

更新业务信息(权限：业务编辑权限)

### 输入参数

| 参数名称                | 参数类型   | 必选 | 描述                            |
|---------------------|--------|----|-------------------------------|
| bk_supplier_account | string | 是  | 开发商账号                         |
| bk_biz_id           | int    | 是  | 业务id                          |
| bk_biz_name         | string | 否  | 业务名                           |
| bk_biz_developer    | string | 否  | 开发人员                          |
| bk_biz_maintainer   | string | 否  | 运维人员                          |
| bk_biz_productor    | string | 否  | 产品人员                          |
| bk_biz_tester       | string | 否  | 测试人员                          |
| operator            | string | 否  | 操作人员                          |
| life_cycle          | string | 否  | 生命周期：测试中(1)，已上线(2, 默认值)，停运(3) |
| language            | string | 否  | 语言, "1"代表中文, "2"代表英文          |

**注意：此处仅对系统内置可编辑的参数做了说明，其余需要填写的参数取决于用户自己定义的属性字段**

### 调用示例

```json
{
  "bk_biz_name": "cc_app_test",
  "bk_biz_maintainer": "admin",
  "bk_biz_productor": "admin",
  "bk_biz_developer": "admin",
  "bk_biz_tester": "admin",
  "language": "1",
  "operator": "admin",
  "life_cycle": "2"
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
| data       | object | 请求返回的数据                    |
| permission | object | 权限信息                       |
