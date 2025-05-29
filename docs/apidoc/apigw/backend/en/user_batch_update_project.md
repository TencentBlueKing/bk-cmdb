### Description

Update Project (Version: v3.10.23+, Permission: Update Project Permission)

### Parameters

| Name | Type   | Required | Description                                                     |
|------|--------|----------|-----------------------------------------------------------------|
| ids  | array  | Yes      | Unique ID array in cc, a maximum of 200 can be passed at a time |
| data | object | Yes      | Fields to be updated                                            |

#### data

| Name               | Type   | Required | Description                                                                                                                                                                                      |
|--------------------|--------|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_project_name    | string | No       | Project name                                                                                                                                                                                     |
| bk_project_desc    | string | No       | Project description                                                                                                                                                                              |
| bk_project_type    | enum   | No       | Project type, optional values: "mobile_game" (mobile game), "pc_game" (PC game), "web_game" (web game), "platform_prod" (platform product), "support_prod" (supporting product), "other" (other) |
| bk_project_sec_lvl | enum   | No       | Confidentiality level, optional values: "public" (public), "private" (private), "classified" (classified)                                                                                        |
| bk_project_owner   | string | No       | Project owner                                                                                                                                                                                    |
| bk_project_team    | array  | No       | Belonging team                                                                                                                                                                                   |
| bk_project_icon    | string | No       | Project icon                                                                                                                                                                                     |
| bk_status          | string | No       | Project status, optional values: "enable" (enabled), "disabled" (disabled)                                                                                                                       |

### Request Example

```json
{
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

### Response Example

```json
{
    "result": true,
    "code": 0,
    "data": null,
    "message": "success",
    "permission": null,
}
```

### Response Parameters

| Name       | Type   | Description                                                                 |
|------------|--------|-----------------------------------------------------------------------------|
| result     | bool   | Indicates whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates failure error                 |
| message    | string | Error message returned in case of request failure                           |
| permission | object | Permission information                                                      |
| data       | object | Data returned in the request                                                |
