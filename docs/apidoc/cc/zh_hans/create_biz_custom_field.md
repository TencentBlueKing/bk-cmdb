### 功能描述

创建业务自定义模型属性

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                  |  类型      | 必选   |  描述                                                    |
|-----------------------|------------|--------|----------------------------------------------------------|
| bk_biz_id             | int        | 是     | 业务id                                                 |
| creator               | string     | 否     | 数据的创建者                                             |
| description           | string     | 否     | 数据的描述信息                                           |
| editable              | bool       | 否     | 表明数据是否可编辑                                       |
| isonly                | bool       | 否     | 表明唯一性                                               |
| ispre                 | bool       | 否     | true:预置字段,false:非内置字段                           |
| isreadonly            | bool       | 否     | true:只读，false:非只读                                  |
| isrequired            | bool       | 否     | true:必填，false:可选                                    |
| option                | string     | 否     |用户自定义内容，存储的内容及格式由调用方决定，以数字类型为例（{"min":"1","max":"2"}）|
| unit                  | string     | 否     | 单位                                                     |
| placeholder           | string     | 否     | 占位符                                                   |
| bk_property_group     | string     | 否     | 字段分栏的名字                                           |
| bk_obj_id             | string     | 是     | 模型ID                                                   |                                              |
| bk_property_id        | string     | 是     | 模型的属性ID                                             |
| bk_property_name      | string     | 是     | 模型属性名，用于展示                                     |
| bk_property_type      | string     | 是     | 定义的属性字段用于存储数据的数据类型,可取值范围（singlechar,longchar,int,enum,date,time,objuser,singleasst,multiasst,timezone,bool）|
| bk_asst_obj_id        | string     | 否     | 如果有关联其它的模型，那么就必需设置此字段，否则就不需要设置                                                                        |

#### bk_property_type

| 标识       | 名字     |
|------------|----------|
| singlechar | 短字符   |
| longchar   | 长字符   |
| int        | 整形     |
| enum       | 枚举类型 |
| date       | 日期     |
| time       | 时间     |
| objuser    | 用户     |
| singleasst | 单关联   |
| multiasst  | 多关联   |
| timezone   | 时区     |
| bool       | 布尔     |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_biz_id": 2,
    "creator": "user",
    "description": "test",
    "editable": true,
    "isonly": false,
    "ispre": false,
    "isreadonly": false,
    "isrequired": false,
    "option": {"min":1,"max":2},
    "unit": "1",
    "placeholder": "test",
    "bk_property_group": "default",
    "bk_obj_id": "cc_test_inst",
    "bk_property_id": "cc_test",
    "bk_property_name": "cc_test",
    "bk_property_type": "singlechar",
    "bk_asst_obj_id": "test"
}
```


### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
	"data": {
		"bk_biz_id": 2,
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
		"option": {"min":1,"max":2},
		"description": "test",
		"creator": "user",
		"create_time": "2020-03-25 17:12:08",
		"last_time": "2020-03-25 17:12:08"
	}
}
```

### 返回结果参数说明
#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                           |

#### data

| 字段                | 类型         | 描述                                                       |
|---------------------|--------------|------------------------------------------------------------|
| bk_biz_id           | int          | 业务自定义字段的业务id                                       |
| creator             | string       | 数据的创建者                                               |
| description         | string       | 数据的描述信息                                             |
| editable            | bool         | 表明数据是否可编辑                                         |
| isonly              | bool         | 表明唯一性                                                 |
| ispre               | bool         | true:预置字段,false:非内置字段                             |
| isreadonly          | bool         | true:只读，false:非只读                                    |
| isrequired          | bool         | true:必填，false:可选                                      |
| option              | string       | 用户自定义内容，存储的内容及格式由调用方决定               |
| unit                | string       | 单位                                                       |
| placeholder         | string       | 占位符                                                     |
| bk_property_group   | string       | 字段分栏的名字                                             |
| bk_property_group_name | string    | 字段分栏的名字 |
| bk_obj_id           | string       | 模型ID                                                     |
| bk_supplier_account | string       | 开发商账号                                                 |
| bk_property_id      | string       | 模型的属性ID                                               |
| bk_property_name    | string       | 模型属性名，用于展示                                       |
| bk_property_type    | string       | 定义的属性字段用于存储数据的数据类型 （singlechar,longchar,int,enum,date,time,objuser,singleasst,multiasst,timezone,bool)|
| bk_asst_obj_id      | string       | 如果有关联其它的模型，那么就必需设置此字段，否则就不需要设置|
| create_time         | string | 创建时间     |
| last_time           | string | 更新时间     |
| id| int | 主键id |

#### bk_property_type

| 标识       | 名字     |
|------------|----------|
| singlechar | 短字符   |
| longchar   | 长字符   |
| int        | 整形     |
| enum       | 枚举类型 |
| date       | 日期     |
| time       | 时间     |
| objuser    | 用户     |
| singleasst | 单关联   |
| multiasst  | 多关联   |
| timezone   | 时区     |
| bool       | 布尔     |
