### Functional description

find module host relation by module ids(v3.8.7)

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field         | Type      | Required | Description                                                  |
| ------------- | --------- | -------- | ------------------------------------------------------------ |
| bk_biz_id     | int       | Yes      | bussiness ID                                                 |
| bk_module_ids | int array | Yes      | module id array, length must be less than 200                |
| module_fields | array     | Yes      | module property list, the specified module property feilds will be returned |
| host_fields   | array     | Yes      | host property list, the specified host property feilds will be returned |
| page          | object    | Yes      | page condition                                               |

#### page

| Field | Type | Required | Description                      |
| ----- | ---- | -------- | -------------------------------- |
| start | int  | Yes      | start record, default is 0       |
| limit | int  | Yes      | page limit, maximum value is 500 |

### Request Parameters Example

```json
{
    "bk_module_ids": [
        1,
        2,
        3
    ],
    "module_fields": [
        "bk_module_id",
        "bk_module_name"
    ],
    "host_fields": [
        "bk_host_innerip",
        "bk_host_id"
    ],
    "page": {
        "start": 0,
        "limit": 500
    }
}
```

### Return Result Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "data": {
    "count": 2,
    "relation": [
      {
        "host": {
          "bk_host_id": 1,
          "bk_host_innerip": "127.0.0.1",
        },
        "modules": [
          {
            "bk_module_id": 1,
            "bk_module_name": "m1",
          },
          {
            "bk_module_id": 2,
            "bk_module_name": "m2",
          }
        ]
      },
      {
        "host": {
          "bk_host_id": 2,
          "bk_host_innerip": "127.0.0.2",
        },
        "modules": [
          {
            "bk_module_id": 3,
            "bk_module_name": "m3",
          }
        ]
      }
    ]
  }
}
```


### Return Result Parameters Description

| Field   | Type   | Description                                            |
| ------- | ------ | ------------------------------------------------------ |
| result  | bool   | request success or failed. true:success；false: failed |
| code    | int    | error code. 0: success, >0: something error            |
| message | string | error info description                                 |
| data    | object | response data                                          |

#### data

| Field    | Type  | Description          |
| -------- | ----- | -------------------- |
| count    | int   | the count of result  |
| relation | array | host and module info |

#### relation

| Field   | Type         | Description             |
| ------- | ------------ | ----------------------- |
| host    | object       | host information        |
| modules | object array | host module information |
