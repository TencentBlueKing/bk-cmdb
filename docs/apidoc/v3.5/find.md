
### 全文检索
* API:  POST /api/{version}/find/full_text
* API名称： full_text_find
* 功能说明：
	* 中文：全文检索
	* English ：full text find
* input body:
```
{
    "page": {
        "start": 0,
        "limit": 10
    },
    "query_string": "test",
    "bk_obj_id": "host",
    "bk_biz_id": "2",
    "filter": ["model"]
}
```


* input字段说明：

| 名称  | 类型 |必填| 默认值|说明 | Description|
|---|---|---|---|---|---|
| start| int|否|-1|分页开始| 分页开始|
| limit| int|否|-1|分页大小| 分页大小|
| query_string| string|否|无|检索内容 | query string content|
| bk_obj_id| string|否|无|对象实例id| 指定对象|
| bk_biz_id| string|否|无|业务id|指定业务|
| filter| string list|否|空|过滤搜索表| 取值范围["model", "object", "host", "process", "application"]|


* output:

```
{
    "result": true,
    "bk_error_code": 0,
    "bk_error_msg": "",
    "data": [
        {
            "type": "object",
            "score": 3.566052,
            "source": {
                  "jw_test_4": 1,
                  "bk_inst_id": 5,
                  "bk_supplier_account": "0",
                  "metadata": {
                      "label": {
                          "bk_biz_id": "2"
                      }
                  },
                  "bk_obj_id": "test",
                  "bk_inst_name": "1",
                  "jw_test_1": "1",
                  "jw_test_2": 12,
                  "jw_test_3": "2019-03-06"
            },
            "highlight": {
                "bk_obj_id": [
                    "<em>test</em>"
                ],
                "bk_obj_id.keyword": [
                    "<em>test</em>"
                 ]
            }
        },
        {
            "type": "host",
            "score": 2.2986379,
            "source": {
                "bk_bak_operator" : null,
                "bk_supplier_account" : "0",
                "bk_disk" : 50,
                "bk_host_innerip" : "10.0.0.6",
                "bk_os_name" : "windows",
                "import_from" : "1",
                "bk_state_name" : null,
                "bk_cloud_id" : 0,
                "bk_cpu_mhz" : 2,
                "bk_mac" : "aa:aa:aa:aa:aa:aa",
                "bk_asset_id" : "",
                "bk_comment" : "this is test host",
                "bk_host_name" : "",
                "bk_host_outerip" : "175.0.0.6",
                "bk_outer_mac" : "aa:aa:aa:aa:aa:aa",
                "operator" : null,
                "bk_isp_name" : null,
                "bk_os_version" : "",
                "bk_service_term" : null,
                "bk_sla" : null,
                "bk_os_type" : null,
                "bk_cpu_module" : "",
                "bk_mem" : null,
                "bk_os_bit" : "32",
                "bk_sn" : "",
                "bk_province_name" : null,
                "bk_cpu" : null,
                "create_time" : null,
                "bk_host_id" : 2
            },
            "highlight": {
                "bk_comment": [
                    "this is <em>test</em> host"
                ],
            }
        },
    ]
}
```
*  output字段说明

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |request result|
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |error code. 0 represent success, >0 represent failure code |
| bk_error_msg | string | 请求失败返回的错误信息 |error message from failed request|
| data | object list | 请求返回的数据 |return data|

data说明：

| 名称  | 类型  | 说明 |Description|
|---|---|---|---|
| type | string | 类型 | 查询结果所属数据类型 |
| score | float | 最佳匹配得分 | 搜索结果匹配程度 |
| source | object | 属性值 | 搜索结果的具体内容 |
| highlight| object| 高亮字段| 匹配高亮显示的字段 |
