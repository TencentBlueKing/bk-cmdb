# 添加模型分类

- API:   POST  /api/{version}/object/classification
- API 名称：create_classification
- 功能说明：
	- 中文： 新建模型分类
	- English：create a object classification

- input  body

``` json
{
    "bk_classification_id": "cs_test",
    "bk_classification_name": "test_name",
    "bk_classification_icon": "icon-cc-business"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

-  input 字段说明参数

|字段|类型|是否必须|默认值|说明|Description|
|---|---|---|---|---| ---|
|bk_classification_id|string|是|无|分类ID，英文描述用于系统内部使用|the classification identifier|
|bk_classification_name|string|是|无|分类名|the name of the classification  |
|bk_classification_icon|string|否|无|模型分类的图标|the icon of the classification|


- output 

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": {
        "id": 18
    }
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Engilish Describe|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|操作结果.|the result|

data 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|id|int|新增数据记录的ID|the data record id|



# 删除模型分类

- API: DELETE   /api/{version}/object/classification/{id}
- API 名称：delete_classification
- 功能说明：
	- 中文： 根据创建分类的时候返回的id删除目标模型分类。
	- English：delete the classification by the id


- input body

    无

- input 字段说明

|字段|类型|是否必须|默认值|说明|Description|
|---|---|---|---|---|---|
|id|int|是|无|分类数据记录ID|the id of the classification data record|


- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": "success"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明


| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|操作结果|the result information|


# 更新模型分类数据


- API: PUT /api/{version}/object/classification/{id}
- API 名称：update_classification
- 功能说明：
	- 中文： 更新模型分类
	- English：update the  object classification

- input body

``` json
{
    "bk_classification_name": "cc_test_new",
    "bk_classification_icon": "icon-cc-business"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**


- input 字段说明

|字段|类型|是否必须|默认值|说明|Description|
|---|---|---|---|---|---|
|id|int|是|无|数据记录的ID|the id of the classification data record|
|bk_classification_name|string|否|无|分类名|the classification name of a object |
|bk_classification_icon|string|否|无|模型分类的图标|the classification icon of a object|


-  output

```  json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": "success"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|string|操作结果|the result information|


# 查询模型分类列表

- API: POST /api/{version}/object/classifications
- API 名称：search_classifications
- 功能说明：
	- 中文： 查询所有模型分类
	- English：get all classifications

- input body

    无

- output
``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": [
        {
            "bk_classification_icon": "icon-cc-business",
            "bk_classification_id": "bk_host_manage",
            "bk_classification_name": "主机管理",
            "bk_classification_type": "inner",
            "id": 1
        }
    ]
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|array|操作结果|the result information|

data 字段说明：

|字段|类型|说明|Description|
|---|---|---|---|
|bk_classification_id|string|分类ID，英文描述用于系统内部使用|the classification identifier|
|bk_classification_name|string|分类名|the classification name|
|bk_classification_type|string|用于对分类进行分类（如：inner）|system inner classification or customize|
|bk_classification_icon|string|模型分类的图标|the classification icon|
|id|int|数据记录ID|the data record id|



# 查询模型分类及附属模型信息

- API: POST /api/{version}/object/classification/{bk_supplier_account}/objects
- API 名称：search_classifications_objects
- 功能说明：
	- 中文： 查询所有模型分类下的模型集合
	- English：get all classification objects

- input body


``` json
{"bk_classification_id":"cc_test"}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。根据实际查询需要，可以在示例的JSON基础上增加或修改为其它字段。**

-  input 字段说明

|字段|类型|是否必须|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_classification_id|string|否|无|分类的标识符|the classification identifier|


- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": [
        {
            "bk_classification_icon": "icon-cc-business",
            "bk_classification_id": "bk_host_manage",
            "bk_classification_name": "主机管理",
            "bk_classification_type": "inner",
            "id": 1,
            "bk_asst_objects": {
                "host": [
                    {
                        "bk_classification_id": "bk_host_manage",
                        "create_time": "2018-03-08T11:30:28.005+08:00",
                        "creator": "cc_system",
                        "description": "",
                        "id": 6,
                        "ispaused": false,
                        "ispre": true,
                        "last_time": null,
                        "modifier": "",
                        "bk_obj_icon": "icon-cc-subnet",
                        "bk_obj_id": "plat",
                        "bk_obj_name": "plat-XXX",
                        "position": "{\"bk_host_manage\":{\"x\":-172,\"y\":-160}}",
                        "bk_supplier_account": "0"
                    }
                ]
            },
            "bk_objects": [
                {
                    "bk_classification_id": "bk_host_manage",
                    "create_time": "2018-03-08T11:30:28.005+08:00",
                    "creator": "cc_system",
                    "description": "",
                    "id": 6,
                    "ispaused": false,
                    "ispre": true,
                    "last_time": null,
                    "modifier": "",
                    "bk_obj_icon": "icon-cc-subnet",
                    "bk_obj_id": "plat",
                    "bk_obj_name": "plat-XXX",
                    "position": "{\"bk_host_manage\":{\"x\":-172,\"y\":-160}}",
                    "bk_supplier_account": "0"
                }
            ]
        }
    ]
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Dscription|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|array|操作结果|the result|

data字段说明

|字段|类型|说明|Description|
|---|---|---|---|
|bk_classification_id|string|分类ID，英文描述用于系统内部使用|the classification identifier|
|bk_classification_name|string|分类名|the classification name it will be shown|
|bk_classification_type|string|用于对分类进行分类（如：inner）|system inner classification or customize|
|bk_classification_icon|string|模型分类的图标|the classification icon|
|bk_objects|array|当前分类下的所有模型|the objects of the classification|
|bk_asst_objects|map[string]array|当前分类下的模型关联的其他模型|the association map|

bk_objects 字段说明

|字段|类型|说明|Description|
|---|---|---|---|
|bk_classification_id|string|分类ID|the classification indentifier|
|create_time|string|创建时间|the creation time|
|creator|string|创建者|the creator|
|description|string|模型描述|the object describtion|
|id|int|模型数据记录的ID|the object data record id|
|ispaused|bool|是否停用|if it is paused|
|ispre|bool|是否内置|if it is the system inner|
|last_time|string|更新时间|the last updated time|
|modifier|string|最后修改人员|the last modifier|
|bk_obj_icon|string|图标|the object icon|
|bk_obj_id|string|模型标识符|the object indentify id|
|bk_obj_name|string|模型的名字，用于展示|the object name ,it will be used to shown|
|bk_position|string|模型在图上的位置|the position ,it will be show in the page|
| bk_supplier_account| string| 开发商账号|supplier account code|


bk_asst_objects 字段说明

|字段|类型|说明|Description|
|---|---|---|---|
|bk_classification_id|string|分类ID|the classification indentifier|
|create_time|string|创建时间|the creation time|
|creator|string|创建者|the creator|
|description|string|模型描述|the object describtion|
|id|int|模型数据记录的ID|the object data record id|
|ispaused|bool|是否停用|if it is paused|
|ispre|bool|是否内置|if it is the system inner|
|last_time|string|更新时间|the last updated time|
|modifier|string|最后修改人员|the last modifier|
|bk_obj_icon|string|图标|the object icon|
|bk_obj_id|string|模型标识符|the object indentify id|
|bk_obj_name|string|模型的名字，用于展示|the object name ,it will be used to shown|
|bk_position|string|模型在图上的位置|the position ,it will be show in the page|
| bk_supplier_account| string| 开发商账号|supplier account code|
