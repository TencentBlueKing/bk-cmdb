
### 获取某部门的用户信息
#
* API:  GET /user/departmentprofile
* API名称： department_profile
* 功能说明：
	* 中文： 获取某部门的用户信息
	* English ：get users info of one department

#
*  input body:
无

#
* input参数说明

| 名称  | 类型 |必填| 默认值 | 说明 | 
| ---  | --- |---| --- | --- | 
| id|int|是|无|部门id |
| recursive|bool|否|false|是否递归查询当前部门下的所有子部门用户信息 |
| page|int|否|1|页码 |
| page_size|int|否|500|每页结果数量 |

#
* input参数请求示例

/user/departmentprofile?id=84631&recursive=true&page_size=2&page=1

#
* ouput 

```
{
  "result": true,
  "bk_error_code": 0,
  "bk_error_msg": "",
  "permission": null,
  "data": {
    "count": 32558,
    "results": [
      {
        "id": 42353,
        "username": "fdsa"
      },
      {
        "id": 42350,
        "username": "GW32559"
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
| permission | object | 权限信息 |
| data | object| 请求返回的数据 |


*  data字段说明：

| 名称  | 类型  | 说明 |
|---|---|---|
| count|int| 符合查询条件的用户数量，包含所有的分页|
| results|object| 用户详情，因为只是当展示当前分页，详情条目数不一定和count数一致|


*  results字段说明：

| 名称  | 类型  | 说明 |
|---|---|---|
| id|int| 用户id|
| username|string| 用户名|
