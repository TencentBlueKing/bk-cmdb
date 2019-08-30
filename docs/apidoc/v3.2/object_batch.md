 
### 导出模型属性

- API: POST /object/owner/{bk_supplier_account}/object/{bk_obj_id}/export
- API 名称：export_object_attribute
- 功能说明：
	- 中文：导出模型属性
	- English：export object attribute

- input body

    无

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_obj_id|string|是|无|模型ID|the object id|

- output

    输出的结果是文件


### 导入模型属性

- API: POST /object/owner/{bk_supplier_account}/object/{bk_obj_id}/import
- API 名称：import_object_attribute
- 功能说明：
	- 中文：导入模型属性
	- English：import object attribute

- input body

    输入的是文件

- input 字段说明

| 字段|类型|必填|默认值|说明|Description|
|---|---|---|---|---|---|
|bk_supplier_account|string|是|无|开发商账号|supplier account code|
|bk_obj_id|string|是|无|模型ID|the object id|

- output

``` json
{
    "result": true,
    "code": 0,
    "message": null,
    "data": {
        "obj_id": {
            "insert_failed": [
                "line:3 msg: type assertion to bool failed"
            ],
            "success": [
                "4",
                "5"
            ],
            "update_failed": [
                "line:4 msg: type assertion to bool failed",
                "line:5 msg: type assertion to bool failed"
            ]
        }
    }
}
```

**注:以上 JSON 数据中各字段的取值仅为示例数据。obj_id 是模型的标识符，其值与输入的 bk_obj_id 保持一致。**

- output字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|result|bool|ture：成功，false：失败 |true:success, false: failure|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
|data|object|操作结果|the result|

data 字段说明

| 字段|类型|说明|Description|
|---|---|---|---|
|obj_id|string|具体的模型ID|the real object identifier|
|insert_failed|array|创建记录失败的错误信息集合，如果没有数据此字段会被省略|the inserted errors array,if there is no data, it will be ignore|
|update_failed|array|更新数据失败的错误信息集合，如果没有数据此字段会被省略|the updated errors array,if there is no data, it will be ignore|
|success|array|操作成功的行号集合，如果没有数据此字段会被省略|the success line num array, if there is no data, it will be ignore|

