### Functional description

update container cluster type (v3.12.1+, permission: kube cluster editing permissions)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type   | Required | Description                                            |
|-----------|--------|----------|--------------------------------------------------------|
| bk_biz_id | int    | yes      | business ID                                            |
| id        | int    | yes      | unique ID of the cluster in cmdb                       |
| type      | string | yes      | cluster type. enum: INDEPENDENT_CLUSTER, SHARE_CLUSTER |

### Request Parameters Example

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 2,
  "id": 1,
  "type": "INDEPENDENT_CLUSTER"
}
```

### Return Result Example

```json
 {
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7",
  "data": null
}
```

### Return Result Parameters Description

#### response

| Name       | Type   | Description                                                                             |
|------------|--------|-----------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. True: request succeeded;false request failed |
| code       | int    | Wrong code. 0 indicates success,>0 indicates failure error                              |
| message    | string | Error message returned by request failure                                               |
| permission | object | Permission information                                                                  |
| request_id | string | Request chain id                                                                        |
| data       | object | Data returned by request                                                                |
