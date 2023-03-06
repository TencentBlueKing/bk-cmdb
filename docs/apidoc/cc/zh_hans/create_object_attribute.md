### 功能描述

创建模型属性

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段              | 类型   | 必选 | 描述                                                         |
| ----------------- | ------ | ---- | ------------------------------------------------------------ |
| creator           | string | 否   | 数据的创建者                                                 |
| description       | string | 否   | 数据的描述信息                                               |
| editable          | bool   | 否   | 表明数据是否可编辑                                           |
| isonly            | bool   | 否   | 表明唯一性                                                   |
| ispre             | bool   | 否   | true:预置字段,false:非内置字段                               |
| isreadonly        | bool   | 否   | true:只读，false:非只读                                      |
| isrequired        | bool   | 否   | true:必填，false:可选                                        |
| option            | string | 否   | 用户自定义内容，存储的内容及格式由调用方决定，以数字类型为例（{"min":"1","max":"2"}） |
| unit              | string | 否   | 单位                                                         |
| placeholder       | string | 否   | 占位符                                                       |
| bk_property_group | string | 否   | 字段分栏的名字                                               |
| bk_obj_id         | string | 是   | 模型ID                                                       |
| bk_property_id    | string | 是   | 模型的属性ID                                                 |
| bk_property_name  | string | 是   | 模型属性名，用于展示                                         |
| bk_property_type  | string | 是   | 定义的属性字段用于存储数据的数据类型,可取值范围（singlechar,longchar,int,enum,date,time,objuser,enummulti,enumquote,timezone,bool,organization) |
| ismultiple        | bool   | 否   | 是否可多选，其中字段类型为短字符，长字符，数字，浮点，枚举，日期，时间，时区，布尔，列表暂时不支持可多选，在创建属性时，字段类型为上述类型可以不传ismultiple参数，默认为false，如果传true则会提示该类型暂不支持可多选。枚举多选，枚举引用，用户，组织字段支持可多选，其中用户字段，组织字段默认为true |
| default           | object | 否   | 给属性字段添加默认值，default的值根据字段的实际类型进行传递，比如创建int类型字段，如果想要给该字段设置默认值，可以传default:5，如果是短字符类型，那么default:"aaa"，不想设置默认值则不传该字段 |

#### bk_property_type

| 标识         | 名字                                                         |
| ------------ | ------------------------------------------------------------ |
| singlechar   | 短字符(不支持可多选，ismultiple参数必须为false，默认为false) |
| longchar     | 长字符(不支持可多选，ismultiple参数必须为false，默认为false) |
| int          | 整形(不支持可多选，ismultiple参数必须为false，默认为false)   |
| enum         | 枚举类型(不支持可多选，ismultiple参数必须为false，默认为false) |
| date         | 日期(不支持可多选，ismultiple参数必须为false，默认为false)   |
| time         | 时间(不支持可多选，ismultiple参数必须为false，默认为false)   |
| objuser      | 用户(支持可多选，ismultiple参数必须为true，默认为true，可在更新时修改为false) |
| timezone     | 时区(不支持可多选，ismultiple参数必须为false，默认为false)   |
| bool         | 布尔(不支持可多选，ismultiple参数必须为false，默认为false)   |
| enummulti    | 枚举(多选) (支持可多选，ismultiple默认为false)               |
| enumquote    | 枚举(引用) (支持可多选，ismultiple默认为false)               |
| organization | 组织(支持可多选，ismultiple默认为true)                       |

### 请求参数示例

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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


### 返回结果示例

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
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

### 返回结果参数说明
#### response

| 名称       | 类型   | 描述                                       |
| ---------- | ------ | ------------------------------------------ |
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误    |
| message    | string | 请求失败返回的错误信息                     |
| permission | object | 权限信息                                   |
| request_id | string | 请求链id                                   |
| data       | object | 请求返回的数据                             |

#### data

| 字段                   | 类型   | 描述                                                         |
| ---------------------- | ------ | ------------------------------------------------------------ |
| creator                | string | 数据的创建者                                                 |
| description            | string | 数据的描述信息                                               |
| editable               | bool   | 表明数据是否可编辑                                           |
| isonly                 | bool   | 表明唯一性                                                   |
| ispre                  | bool   | true:预置字段,false:非内置字段                               |
| isreadonly             | bool   | true:只读，false:非只读                                      |
| isrequired             | bool   | true:必填，false:可选                                        |
| option                 | string | 用户自定义内容，存储的内容及格式由调用方决定                 |
| unit                   | string | 单位                                                         |
| placeholder            | string | 占位符                                                       |
| bk_property_group      | string | 字段分栏的名字                                               |
| bk_obj_id              | string | 模型ID                                                       |
| bk_supplier_account    | string | 开发商账号                                                   |
| bk_property_id         | string | 模型的属性ID                                                 |
| bk_property_name       | string | 模型属性名，用于展示                                         |
| bk_property_type       | string | 定义的属性字段用于存储数据的数据类型 （singlechar,longchar,int,enum,date,time,objuser,singleasst,multiasst,timezone,bool) |
| bk_biz_id              | int    | 业务自定义字段的业务id                                       |
| bk_property_group_name | string | 字段分栏的名字                                               |
| ismultiple             | bool   | 字段是否支持可多选                                           |
| default                | object | 属性默认值                                                   |

#### bk_property_type

| 标识         | 名字     |
| ------------ | -------- |
| singlechar   | 短字符   |
| longchar     | 长字符   |
| int          | 整形     |
| enum         | 枚举类型 |
| date         | 日期     |
| time         | 时间     |
| objuser      | 用户     |
| timezone     | 时区     |
| bool         | 布尔     |
| enummulti    | 枚举多选 |
| enumquote    | 枚举引用 |
| organization | 组织     |