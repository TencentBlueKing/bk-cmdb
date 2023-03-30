### 功能描述

更新对象模型属性

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段              | 类型   | 必选 | 描述                                                         |
| ----------------- | ------ | ---- | ------------------------------------------------------------ |
| id                | int    | 是   | 目标数据的记录ID                                             |
| description       | string | 否   | 数据的描述信息                                               |
| isonly            | bool   | 否   | 表明唯一性                                                   |
| isreadonly        | bool   | 否   | 表明是否只读                                                 |
| isrequired        | bool   | 否   | 表明是否必填                                                 |
| bk_property_group | string | 否   | 字段分栏的名字                                               |
| option            | string | 否   | 用户自定义内容，存储的内容及格式由调用方决定, 以数字内容为例（{"min":"1","max":"2"}） |
| bk_property_name  | string | 否   | 模型属性名，用于展示                                         |
| bk_property_type  | string | 否   | 定义的属性字段用于存储数据的数据类型（singlechar,longchar,int,enum,date,time,objuser,enummulti,enumquote,timezone,bool,organization) |
| unit              | string | 否   | 单位                                                         |
| placeholder       | string | 否   | 占位符                                                       |
| ismultiple        | bool   | 否   | 是否可多选，其中字段类型为短字符，长字符，数字，浮点，枚举，日期，时间，时区，布尔，列表暂时不支持可多选，在更新属性时，字段类型为上述类型时，不能将ismultiple更新为true，如果更新为true则会提示该类型暂不支持可多选。枚举多选，枚举引用，用户，组织字段支持可多选。 |
| default           | object | 否   | 给属性添加默认值，更新的时候，default的值根据字段的实际类型进行传递，如果想要置空字段的默认值，需要传递default:null |

#### bk_property_type

| 标识         | 名字       |
| ------------ | ---------- |
| singlechar   | 短字符     |
| longchar     | 长字符     |
| int          | 整形       |
| enum         | 枚举类型   |
| date         | 日期       |
| time         | 时间       |
| objuser      | 用户       |
| enummulti    | 枚举(多选) |
| enumquote    | 枚举(引用) |
| timezone     | 时区       |
| bool         | 布尔       |
| organization | 组织       |

### 请求参数示例

更新默认值场景

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "id":1,
    "description":"test",
    "placeholder":"test",
    "unit":"1",
    "isonly":false,
    "isreadonly":false,
    "isrequired":false,
    "bk_property_group":"default",
    "option":"{\"min\":\"1\",\"max\":\"4\"}",
    "bk_property_name":"aaa",
    "bk_property_type":"int",
    "bk_asst_obj_id":"0",
    "default":3
}
```

不更新默认值场景

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "id":1,
    "description":"test",
    "placeholder":"test",
    "unit":"1",
    "isonly":false,
    "isreadonly":false,
    "isrequired":false,
    "bk_property_group":"default",
    "option":"{\"min\":\"1\",\"max\":\"4\"}",
    "bk_property_name":"aaa",
    "bk_property_type":"int",
    "bk_asst_obj_id":"0"
}
```

置空默认值场景

```python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "id":1,
    "description":"test",
    "placeholder":"test",
    "unit":"1",
    "isonly":false,
    "isreadonly":false,
    "isrequired":false,
    "bk_property_group":"default",
    "option":"{\"min\":\"1\",\"max\":\"4\"}",
    "bk_property_name":"aaa",
    "bk_property_type":"int",
    "bk_asst_obj_id":"0",
    "default":null
}
```



### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "e43da4ef221746868dc4c837d36f3807",
    "data": "success"
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
| data       | object | 无数据返回                                 |