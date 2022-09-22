### 后端代理
* API:  POST /proxy/{method}/{target}/{target_url}
* 功能说明：

  后端代理前端跨域请求，前端类似于向用户管理等其他saas发起的的请求交由后端转发，避免跨域问题。POST或PUT请求等参数存放在body的仍然按照之前的请求方式请求，只需要更换请求对象为CMDB后端即可，CMDB后台转发路由不会修改body中的内容，只起到路由转发的作用

* Url字段说明：

| 名称         | 类型        | 必填 | 默认值 | 说明       | Description                                                   |
| ------------ | ----------- | ---- | ------ | ---------- | ------------------------------------------------------------- |
| method       | string     | 是   | GET     | 前端请求的方式 | 前端请求的方式                                             |
| target       | string      | 是   | 无     | 请求目标   | 请求目标，例如usermanage，用户管理                                  |
| target_url | string      | 是   | 无     | 请求的url      | 请求目标的url，例如：/api/c/compapi/v2/usermanage/fs_list_users/?app_code=bk-magicbox&page=1&page_size=20&fuzzy_lookups=c&callback=USER_LIST_CALLBACK_1 |

* 样例:

```
请求用户管理：

原请求由前端发起：
http://paas.xxx/api/c/compapi/v2/usermanage/fs_list_users/?app_code=bk-magicbox&page=1&page_size=20&fuzzy_lookups=a&callback=USER_LIST_CALLBACK_1

现在请求后端由后端转发：
http://cmdb.xxx/proxy/get/usermanage/api/c/compapi/v2/usermanage/fs_list_users/?app_code=bk-magicbox&page=1&page_size=20&callback=USER_LIST_CALLBACK_1&fuzzy_lookups=a
```
* 返回:

```
USER_LIST_CALLBACK_1({"message": "success", "code": 0, "data": {"count": 1, "results": [{"username": "admin", "domain": "default.local", "display_name": "", "logo": null, "category_id": 1, "id": 1, "category_name": "\u9ed8\u8ba4\u76ee\u5f55"}]}, "result": true, "request_id": "1d0f7dbd8f404efab4ab1bdfe3b79ab4"})
```

