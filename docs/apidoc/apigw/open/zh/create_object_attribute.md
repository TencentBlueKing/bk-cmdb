### 描述

创建模型属性(权限：模型编辑权限)

### 输入参数

| 参数名称              | 参数类型   | 必选 | 描述                                                                                                                                                                                             |
|-------------------|--------|----|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| creator           | string | 否  | 数据的创建者                                                                                                                                                                                         |
| description       | string | 否  | 数据的描述信息                                                                                                                                                                                        |
| editable          | bool   | 否  | 表明数据是否可编辑                                                                                                                                                                                      |
| isonly            | bool   | 否  | 表明唯一性                                                                                                                                                                                          |
| ispre             | bool   | 否  | true:预置字段,false:非内置字段                                                                                                                                                                          |
| isreadonly        | bool   | 否  | true:只读，false:非只读                                                                                                                                                                              |
| isrequired        | bool   | 否  | true:必填，false:可选                                                                                                                                                                               |
| option            | string | 否  | 用户自定义内容，存储的内容及格式由调用方决定，以数字类型为例（{"min":"1","max":"2"}）                                                                                                                                          |
| unit              | string | 否  | 单位                                                                                                                                                                                             |
| placeholder       | string | 否  | 占位符                                                                                                                                                                                            |
| bk_property_group | string | 否  | 字段分栏的名字                                                                                                                                                                                        |
| bk_obj_id         | string | 是  | 模型ID                                                                                                                                                                                           |
| bk_property_id    | string | 是  | 模型的属性ID                                                                                                                                                                                        |
| bk_property_name  | string | 是  | 模型属性名，用于展示                                                                                                                                                                                     |
| bk_property_type  | string | 是  | 定义的属性字段用于存储数据的数据类型,可取值范围 （singlechar(短字符),longchar(长字符),int(整形),enum(枚举类型),date(日期),time(时间),objuser(用户),enummulti(枚举多选),enumquote(枚举引用),timezone(时区),bool(布尔),organization(组织),id_rule(id规则)) |
| ismultiple        | bool   | 否  | 是否可多选，其中字段类型为短字符，长字符，数字，浮点，枚举，日期，时间，时区，布尔，列表暂时不支持可多选，在创建属性时，字段类型为上述类型可以不传ismultiple参数，默认为false，如果传true则会提示该类型暂不支持可多选。枚举多选，枚举引用，用户，组织字段支持可多选，其中用户字段，组织字段默认为true                                 |
| default           | object | 否  | 给属性字段添加默认值，default的值根据字段的实际类型进行传递，比如创建int类型字段，如果想要给该字段设置默认值，可以传default:5，如果是短字符类型，那么default:"aaa"，不想设置默认值则不传该字段                                                                                |

### 调用示例

```json
{
    "creator": "user",
    "description": "test",
    "editable": true,
    "isonly": false,
    "ispre": false,
    "isreadonly": false,
    "isrequired": false,
    "option": "^[0-9a-zA-Z_]{1,}$",
    "unit": "1",
    "placeholder": "test",
    "bk_property_group": "default",
    "bk_obj_id": "cc_test_inst",
    "bk_property_id": "cc_test",
    "bk_property_name": "cc_test",
    "bk_property_type": "singlechar",
    "bk_asst_obj_id": "test",
    "ismultiple": false,
    "default":"aaaa"
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
	"data": {
		"id": 7,
		"bk_supplier_account": "0",
		"bk_obj_id": "cc_test_inst",
		"bk_property_id": "cc_test",
		"bk_property_name": "cc_test",
		"bk_property_group": "default",
		"bk_property_index": 4,
		"unit": "1",
		"placeholder": "test",
		"editable": true,
		"ispre": false,
		"isrequired": false,
		"isreadonly": false,
		"isonly": false,
		"bk_issystem": false,
		"bk_isapi": false,
		"bk_property_type": "singlechar",
		"option": "",
		"description": "test",
		"creator": "user",
		"create_time": "2020-03-25 17:12:08",
		"last_time": "2020-03-25 17:12:08",
		"bk_property_group_name": "default",
        	"ismultiple": false,
        	"default":"aaaa"
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

| 参数名称                   | 参数类型   | 描述                                                                                                                                                                                       |
|------------------------|--------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| creator                | string | 数据的创建者                                                                                                                                                                                   |
| description            | string | 数据的描述信息                                                                                                                                                                                  |
| editable               | bool   | 表明数据是否可编辑                                                                                                                                                                                |
| isonly                 | bool   | 表明唯一性                                                                                                                                                                                    |
| ispre                  | bool   | true:预置字段,false:非内置字段                                                                                                                                                                    |
| isreadonly             | bool   | true:只读，false:非只读                                                                                                                                                                        |
| isrequired             | bool   | true:必填，false:可选                                                                                                                                                                         |
| option                 | string | 用户自定义内容，存储的内容及格式由调用方决定                                                                                                                                                                   |
| unit                   | string | 单位                                                                                                                                                                                       |
| placeholder            | string | 占位符                                                                                                                                                                                      |
| bk_property_group      | string | 字段分栏的名字                                                                                                                                                                                  |
| bk_obj_id              | string | 模型ID                                                                                                                                                                                     |
| bk_supplier_account    | string | 开发商账号                                                                                                                                                                                    |
| bk_property_id         | string | 模型的属性ID                                                                                                                                                                                  |
| bk_property_name       | string | 模型属性名，用于展示                                                                                                                                                                               |
| bk_property_type       | string | 定义的属性字段用于存储数据的数据类型 （singlechar(短字符),longchar(长字符),int(整形),enum(枚举类型),date(日期),time(时间),objuser(用户),enummulti(枚举多选),enumquote(枚举引用),timezone(时区),bool(布尔),organization(组织),id_rule(id规则)) |
| bk_biz_id              | int    | 业务自定义字段的业务id                                                                                                                                                                             |
| bk_property_group_name | string | 字段分栏的名字                                                                                                                                                                                  |
| ismultiple             | bool   | 字段是否支持可多选                                                                                                                                                                                |
| default                | object | 属性默认值                                                                                                                                                                                    |
