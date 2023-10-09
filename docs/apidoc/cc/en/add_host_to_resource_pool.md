### Functional description

Adds a host to the resource pool with the specified id based on the host list information

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                  | Type        | Required	 |Description                |
|----------------------|--------------|--------|---------------------|
| host_info            |  object array |yes     | Host information              |
| directory            |  int          | no    | Resource directory ID |

#### host_info
| Field             | Type| Required| Description                    |
|-----------------|--------|-----|-------------------------|
| bk_host_innerip | string |yes| Host intranet ip                |
| bk_cloud_id | int |yes| Cloud area id     |
| bk_host_name    |  string |no| Host name, or any other property    |
| operator        |  string | no       | Main maintainer, or other attributes|
| bk_comment      |  string |no| Comments, or other attributes      |

### Request Parameters Example

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "host_info": [
        {
            "bk_host_innerip": "127.0.0.1",
            "bk_host_name": "host1",
            "bk_cloud_id": 0,
            "operator": "admin",
            "bk_comment": "comment"
        },
        {
            "bk_host_innerip": "127.0.0.2",
            "bk_host_name": "host2",
            "bk_cloud_id": 0,
            "operator": "admin",
            "bk_comment": "comment"
        }
    ],
    "directory": 1
}
```

### Return Result Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "data": {
      "success": [
          {
              "index": 0,
              "bk_host_id": 6
          },
          {
              "index": 1,
              "bk_host_id": 7
          }
      ]
  },
  "permission": null,
  "request_id": "e43da4ef221746868dc4c837d36f3807"
}

```

### Return Result Parameters Description

#### response

| Name    | Type   | Description                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                    |
| data    |  object |Data returned by request                           |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |

#### Data field Description

| Field     | Type| Description                |
| ------- | ----- | ------------------ |
| success | array |Host information array added successfully|
| error   |  array |Add failed host info array|

#### Success Field Description

| Field        | Type| Description             |
| ---------- | ---- | --------------- |
| index      |  int  |Add successful host subscripts|
| bk_host_id | int  |Successfully added host ID   |

#### Error Field Description

| Field           | Type   | Description             |
| ------------- | ------ | --------------- |
| index         |  int    | Add failed host subscript|
| error_message | string |Failure reason         |
