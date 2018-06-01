# 添加对象模型

- API: POST /api/{version}/object
- API 名称: create_object
- 功能说明：
    - 中文：定义模型
    - English：create a object

- input body

``` json
{
    "creator": "admin",
    "modifier": "admin",
    "bk_classification_id": "cc_test",
    "bk_obj_name": "cc_test_inst",
    "bk_supplier_account": "0",
    "bk_obj_icon": "icon-cc-business",
    "bk_obj_id": "cc_test_inst"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|creator|string|否|无|本条数据创建者|creator|
|modifier|string|否|无|本条数据的最后修改人员|modifier|
|description|string|否|无|对象模型的描述信息|description|
|bk_classification_id|string|是|无|对象模型的分类ID，只能用英文字母序列命名|the classification identifier|
|bk_obj_id|string|是|无|对象模型的ID，只能用英文字母序列命名|the object identifier|
|bk_obj_name|string|是|无|对象模型的名字，用于展示，可以使用人类可以阅读的任何语言|the object name ,it will be used to shown|
| bk_supplier_account| string| 是| 无|开发商账号|supplier account code|
|bk_obj_icon|string|否|无|对象模型的ICON信息，用于前端显示|the icon of the object|
|position|string|否|无|用于存储前端在页面上显示的模型的位置信息，数据格式自己定义，长度最多 1024个字符|	the position ,it will be show in the page |


- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": {
        "id": 1038
    }
}
```


**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- output 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|操作结果|the result|

data字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|id|int|新增的数据记录的ID|the data record identifier|

# 删除对象模型

- API: DELETE /api/{version}/object/{id}
- API 名称: delete_object
- 功能说明：
    - 中文：删除模型
    - English：delete a object

- input body

    无

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|id|int|否|无|被删除的数据记录的ID|the id of the target data record|


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
|data|string|结果信息|the result |


# 更新对象模型

- API: PUT /api/{version}/object/{id}
- API 名称: update_object
- 功能说明：
    - 中文：更新模型定义
    - English：update a object

- input body

``` json
{
    "creator": "admin",
    "modifier": "admin",
    "bk_classification_id": "cc_test",
    "bk_obj_name": "cc2_test_inst",
    "bk_supplier_account": "0",
    "bk_obj_icon": "icon-cc-business",
}
```


**注:以上 JSON 数据中各字段的取值仅为示例数据。**

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|id|int|否|无|目标数据的记录ID|the record id|
|creator|string|否|无|本条数据创建者|the creator|
|modifier|string|否|无|本条数据的最后修改人员|the modifier|
|description|string|否|无|对象模型的描述信息|description|
|bk_classification_id|string|是|无|对象模型的分类ID，只能用英文字母序列命名|the classification identifier|
|bk_obj_name|string|否|无|对象模型的名字|the name of the object|
|bk_supplier_account| string| 是| 无|开发商账号|supplier account code|
|bk_obj_icon|string|否|无|对象模型的ICON信息，用于前端显示|the icon of the object|


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
|data|string|结果数据|the result|


#  查询对象模型

- API: POST /api/{version}/objects
- API 名称: search_objects
- 功能说明：
    - 中文：查询模型
    - English：search a object

- input body

``` json
{
    "bk_obj_id": "biz",
    "bk_supplier_account":"0"
}
```


**注:以上 JSON 数据中各字段的取值仅为示例数据。实际使用中可以按照查询的需求填写多个字段。**

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|id|string|否|无|数据记录的ID|the data record id|
|creator|string|否|无|本条数据创建者|creator|
|modifier|string|否|无|本条数据的最后修改人员|modifier|
|description|string|否|无|对象模型的描述信息|description|
|bk_classification_id|string|否|无|对象模型的分类ID，只能用英文字母序列命名|the classifition identifier|
|bk_obj_id|string|否|无|对象模型的ID，只能用英文字母序列命名|the object identifier|
|bk_obj_name|string|否|无|对象模型的名字，用于展示，可以使用人类可以阅读的任何语言|the name of the object, it will be used to shown|
|bk_supplier_account| string| 否| 无|开发商账号|supplier account code|
|bk_obj_icon|string|否|无|对象模型的ICON信息，用于前端显示|the icon of the object|
|position|string|否|无|用于存储前端在页面上现实的模型的位置信息，数据格式自己定义，长度最多 1024个字符|the position ,it will be show in the page|


- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": [
        {
            "bk_classification_id": "bk_organization",
            "create_time": "2018-03-08T11:30:28.005+08:00",
            "creator": "cc_system",
            "description": "",
            "id": 4,
            "ispaused": false,
            "ispre": true,
            "last_time": null,
            "modifier": "",
            "bk_obj_icon": "icon-XXX",
            "bk_obj_id": "XX",
            "bk_obj_name": "XXX",
            "position": "{\"test_obj\":{\"x\":-253,\"y\":137}}",
            "bk_supplier_account": "0"
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
|data|object|结果数据|the result|

data 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|id|int|数据记录的ID|the record identifier|
|creator|string|本条数据创建者|creator|
|modifier|string|本条数据的最后修改人员|modifier|
|description|string|对象模型的描述信息|description|
|bk_classification_id|string|对象模型的分类ID，只能用英文字母序列命名|the classifition identifier|
|bk_obj_id|string|对象模型的ID，只能用英文字母序列命名|the object identifier|
|bk_obj_name|string|对象模型的名字，用于展示|the name of the object, it will be used to shown|
| bk_supplier_account| string|开发商账号|supplier account code|
|bk_obj_icon|string|对象模型的ICON信息，用于前端显示|the icon of the object|
|position|string|用于存储前端在页面上显示的模型的位置信息，数据格式自己定义，长度最多 1024个字符|the position ,it will be show in the page|

#  查询普通对象模型的拓扑结构

- API: POST  /api/{version}/objects/topo
- API 名称: search_object_topo
- 功能说明：
    - 中文：查询普通模型拓扑
    - english：search a object topo

- input body

``` json
{
    "bk_classification_id": "bk_host_manage"
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。实际使用中可以按照查询的需求填写多个字段。**

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---| ---|
|id|string|否|无|数据记录的ID|the the record identifier|
|creator|string|否|无|本条数据创建者|creator|
|modifier|string|否|无|本条数据的最后修改人员|modifier|
|description|string|否|无|对象模型的描述信息|description|
|bk_classification_id|string|否|无|对象模型的分类ID，只能用英文字母序列命名|the classification identifier|
|bk_obj_id|string|否|无|对象模型的ID，只能用英文字母序列命名|the object identifier|
|bk_obj_name|string|否|无|对象模型的名字，用于展示，可以使用人类可以阅读的任何语言|the object name|
| bk_supplier_account| string| 否| 无|开发商账号|supplier account code|
|bk_obj_icon|string|否|无|对象模型的ICON信息，用于前端显示|the icon of the object|
|position|string|否|无|用于存储前端在页面上现实的模型的位置信息，数据格式自己定义，长度最多 1024个字符|the coordinate of the object|


- output

``` json
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": null,
    "data": [
        {
            "arrows": "to",
            "from": {
                "bk_classification_id": "bk_host_manage",
                "bk_obj_id": "host",
                "bk_obj_name": "主机",
                "position": "{\"bk_host_manage\":{\"x\":-357,\"y\":-344},\"lhmtest\":{\"x\":163,\"y\":75}}",
                "bk_supplier_account": "0"
            },
            "label": "bk_cloud_id",
            "label_name": "",
            "label_type": "",
            "to": {
                "bk_classification_id": "bk_host_manage",
                "bk_obj_id": "plat",
                "bk_obj_name": "云区域",
                "position": "{\"bk_host_manage\":{\"x\":-172,\"y\":-160}}",
                "bk_supplier_account": "0"
            }
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
|data|array|结果数据|the result|

data 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|arrows|string|取值 to（单向） 或 to,from（双向）|direction|
|label_name|string|关联关系的名字|the associated name|
|label|string|表明From通过哪个字段关联到To的|the associated attribute|
|from|string|对象模型的英文id，拓扑关系的发起方|the starting point of the association|
|to|string|对象模型的英文ID，拓扑关系的终止方|the associated end point|



