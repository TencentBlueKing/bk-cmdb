### 描述

更新业务自定义模型属性(权限：业务自定义字段编辑权限)

### 输入参数

| 参数名称              | 参数类型                  | 必选 | 描述                                                                                                                                                                         |
|-------------------|-----------------------|----|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| id                | int                   | 是  | 目标数据的记录ID                                                                                                                                                                  |
| bk_biz_id         | int                   | 是  | 业务id                                                                                                                                                                       |
| description       | string                | 否  | 数据的描述信息                                                                                                                                                                    |
| isonly            | bool                  | 否  | 表明唯一性                                                                                                                                                                      |
| isreadonly        | bool                  | 否  | 表明是否只读                                                                                                                                                                     |
| isrequired        | bool                  | 否  | 表明是否必填                                                                                                                                                                     |
| bk_property_group | string                | 否  | 字段分栏的名字                                                                                                                                                                    |
| option            | object                | 否  | 用户自定义内容，存储的内容及格式由调用方决定, 以数字内容为例（{"min":1,"max":2}）                                                                                                                         |
| bk_property_name  | string                | 否  | 模型属性名，用于展示                                                                                                                                                                 |
| unit              | string                | 否  | 单位                                                                                                                                                                         |
| bk_property_type  | string                | 否  | 定义的属性字段用于存储数据的数据类型 （singlechar(短字符),longchar(长字符),int(整形),enum(枚举类型),date(日期),time(时间),objuser(用户),enummulti(枚举多选),enumquote(枚举引用),timezone(时区),bool(布尔),organization(组织)) |
| placeholder       | string                | 否  | 占位符                                                                                                                                                                        |
| bk_asst_obj_id    | string                | 否  | 如果有关联其它的模型，那么就必需设置此字段，否则就不需要设置                                                                                                                                             |
| default           | 随bk_property_type字段类型 | 否  | 默认值                                                                                                                                                                        |

### 调用示例

```json
{

    "id":1,
    "bk_biz_id": 2,
    "description":"test",
    "placeholder":"test",
    "unit":"1",
    "isonly":false,
    "isreadonly":false,
    "isrequired":false,
    "bk_property_group":"default",
    "option":{"min":1,"max":4},
    "bk_property_name":"aaa",
    "bk_property_type":"int",
    "bk_asst_obj_id":"0"
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
