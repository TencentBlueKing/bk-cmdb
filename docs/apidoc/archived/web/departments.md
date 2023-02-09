
### 获取部门列表
#
* API:  GET /user/department
* API名称： departments
* 功能说明：
	* 中文： 获取部门列表
	* English ：get department list

#
*  input body：无

#
* input参数说明

| 名称  | 类型 |必填| 默认值 | 说明 | 
| ---  | --- |---| --- | --- | 
| lookup_field|string|否|id|查找字段 |
| exact_lookups|string|否|无|精确查找内容列表，与lookup_field字段一起使用 |
| fuzzy_lookups|string|否|无|模糊查找内容列表，与lookup_field字段一起使用 |
| page|int|否|1|页码 |
| page_size|int|否|500|每页结果数量 |
| fields|string|否|所有字段|返回值字段 |
| with_ancestors|bool|否|false|是否列出当前部门的所有上级部门 |

#
* input参数请求示例

/user/department?lookup_field=level&exact_lookups=3&page_size=1&page=2

#
* ouput 

```
{
  "result": true,
  "bk_error_code": 0,
  "bk_error_msg": "",
  "permission": null,
  "data": {
    "count": 335,
    "results": [
      {
        "id": 44548,
        "parent": 44547,
        "name": "工装科",
        "full_name": "保定市长/生产一工/生管物流/工装科",
        "level": 3,
        "has_children": false
      }
    ]
  }
}
```
*  output字段说明

| 名称  | 类型  | 说明 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| bk_error_code | int | 错误编码。 0表示success，>0表示失败错误 |
| bk_error_msg | string | 请求失败返回的错误信息 |
| bk_error_msg | string | 请求失败返回的错误信息 |
| permission | object | 权限信息 |
| data | object| 请求返回的数据 |


*  data字段说明：

| 名称  | 类型  | 说明 |
|---|---|---|
| count|int| 符合查询条件的部门数量，包含所有的分页|
| results|object| 部门详情，因为只是当展示当前分页，详情条目数不一定和count数一致|


*  results字段说明：

| 名称  | 类型  | 说明 |
|---|---|---|
| id|int| 部门id|
| parent|id| 直接上级部门id|
| name|string| 部门名，只展示当前层级|
| full_name|string| 部门全名，展示从顶层到当前层级的全部路径|
| level|int|部门层级，顶层为0级，依次类推|
| has_children|bool| 是否有下级部门|
