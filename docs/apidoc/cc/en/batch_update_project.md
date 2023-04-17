### Function description

batch update project (version: v3.10.23+, permission: update permission of the project)

### Request parameters

{{ common_args_desc }}


#### Interface parameters

| field | type | required | description |
| ----------------------------|------------|----------|-------------------------------------------|
| ids | array| yes      | an array of ids uniquely identified in cc, limited to 200 at a time |
| data | object| yes      |fields that need to be updated|

#### data

| field | type | required | description                                                                                                                           |
|--------------------|------------|----------|---------------------------------------------------------------------------------------------------------------------------------------|
| bk_project_name | string | no      | project_name                                                                                                                          |
| bk_project_code | string | no      | project_code                                                                                                                          |
| bk_project_desc | string | no       | project_description                                                                                                                   |
| bk_project_type | enum | no       | project type, optional values: "mobile_game", "pc_game", "web_game", "platform_prod", "support_prod", "other", default value: "other" |
| bk_project_sec_lvl | enum | no       | confidentiality level, optional values: "public", "private", "classified", default: "public"                                          |
| bk_project_owner | string | no      | project owner                                                                                                                         |
| bk_project_team | array | no       | project team                                                                                                                          |
| bk_project_icon | string | no       | project icon                                                                                                                          |

### Request parameter examples

```json
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "ids":[
        1, 2, 3
    ],   
    "data": {
        "bk_project_name": "test",
        "bk_project_desc": "test project",
        "bk_project_type": "mobile_game",
        "bk_project_sec_lvl": "public",
        "bk_project_owner": "admin",
        "bk_project_team": [1, 2],
        "bk_status": "enable",
        "bk_project_icon": "https://127.0.0.1/file/png/11111"
    }
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "data": null,
    "message": "success",
    "permission": null,
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### Return result parameter description
#### response

| name | type | description |
| ------- | ------ | ------------------------------------- |
| result | bool | The success or failure of the request. true: the request was successful; false: the request failed.|
| code | int | The error code. 0 means success, >0 means failure error.|
| message | string | The error message returned by the failed request.|
| permission | object | Permission information |
| request_id | string | request_chain_id |
| data | object | data returned by the request|

