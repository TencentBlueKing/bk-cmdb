### 描述

可通过可选参数根据模型id或业务id查询对象模型属性(权限：模型查看权限)

### 输入参数

| 参数名称      | 参数类型   | 必选 | 描述                    |
|-----------|--------|----|-----------------------|
| bk_obj_id | string | 是  | 模型ID                  |
| bk_biz_id | int    | 否  | 业务id，设置后查询结果包含业务自定义字段 |

### 调用示例

```json
{
    "bk_obj_id": "test",
    "bk_biz_id": 2
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "data": [
       {
           "bk_biz_id": 0,
           "bk_asst_obj_id": "",
           "bk_asst_type": 0,
           "create_time": "2018-03-08T11:30:27.898+08:00",
           "creator": "cc_system",
           "description": "",
           "editable": false,
           "id": 51,
           "isapi": false,
           "isonly": true,
           "ispre": true,
           "isreadonly": false,
           "isrequired": true,
           "last_time": "2018-03-08T11:30:27.898+08:00",
           "bk_obj_id": "process",
           "option": "",
           "placeholder": "",
           "bk_property_group": "default",
           "bk_property_group_name": "基础信息",
           "bk_property_id": "bk_process_name",
           "bk_property_index": 0,
           "bk_property_name": "进程名称",
           "bk_property_type": "singlechar",
           "bk_supplier_account": "0",
           "unit": ""
       },
       {
            "bk_biz_id": 2,
            "id": 7,
            "bk_supplier_account": "0",
            "bk_obj_id": "process",
            "bk_property_id": "biz_custom_field",
            "bk_property_name": "业务自定义字段",
            "bk_property_group": "biz_custom_group",
            "bk_property_index": 4,
            "unit": "",
            "placeholder": "",
            "editable": true,
            "ispre": true,
            "isrequired": false,
            "isreadonly": false,
            "isonly": false,
            "bk_issystem": false,
            "bk_isapi": false,
            "bk_property_type": "singlechar",
            "option": "",
            "description": "",
            "creator": "admin",
            "create_time": "2020-03-25 17:12:08",
            "last_time": "2020-03-25 17:12:08",
            "bk_property_group_name": "业务自定义分组"
       }
   ]
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

| 参数名称                | 参数类型   | 描述                                                                                                                                                                         |
|---------------------|--------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| creator             | string | 数据的创建者                                                                                                                                                                     |
| description         | string | 数据的描述信息                                                                                                                                                                    |
| editable            | bool   | 表明数据是否可编辑                                                                                                                                                                  |
| isonly              | bool   | 表明唯一性                                                                                                                                                                      |
| ispre               | bool   | true:预置字段,false:非内置字段                                                                                                                                                      |
| isreadonly          | bool   | true:只读，false:非只读                                                                                                                                                          |
| isrequired          | bool   | true:必填，false:可选                                                                                                                                                           |
| option              | string | 用户自定义内容，存储的内容及格式由调用方决定                                                                                                                                                     |
| unit                | string | 单位                                                                                                                                                                         |
| placeholder         | string | 占位符                                                                                                                                                                        |
| bk_property_group   | string | 字段分栏的名字                                                                                                                                                                    |
| bk_obj_id           | string | 模型ID                                                                                                                                                                       |
| bk_supplier_account | string | 开发商账号                                                                                                                                                                      |
| bk_property_id      | string | 模型的属性ID                                                                                                                                                                    |
| bk_property_name    | string | 模型属性名，用于展示                                                                                                                                                                 |
| bk_property_type    | string | 定义的属性字段用于存储数据的数据类型 （singlechar(短字符),longchar(长字符),int(整形),enum(枚举类型),date(日期),time(时间),objuser(用户),enummulti(枚举多选),enumquote(枚举引用),timezone(时区),bool(布尔),organization(组织)) |
| bk_asst_obj_id      | string | 如果有关联其它的模型，那么就必需设置此字段，否则就不需要设置                                                                                                                                             |
| bk_biz_id           | int    | 业务自定义字段的业务id                                                                                                                                                               |
| create_time         | string | 创建时间                                                                                                                                                                       |
| last_time           | string | 更新时间                                                                                                                                                                       |
| id                  | int    | 查询对象的id值                                                                                                                                                                   |
