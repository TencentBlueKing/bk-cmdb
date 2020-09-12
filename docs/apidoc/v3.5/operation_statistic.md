运营统计
-----------
**新增统计图表**

- API: POST /api/v3/create/operation/chart

- API名称：create_operation_chart

- 功能说明：
    中文：新增统计图表
    English： create_operation_chart

- input body:
  
  ```
  {
      "bk_obj_id": "host"
      "chart_type"": "pie"
      "config_id"": null
      "field": "bk_sla"
      "name": "123"
      "report_type": "custom"
      "width": "50"
      "x_axis_count": 10
  }
  ```
  
  
  
- input字段说明：

| 名称         | 类型   | 必填 | 默认值 | 说明                                         | Description                                                            |
| ------------ | ------ | ---- | ------ | -------------------------------------------- | ---------------------------------------------------------------------- |
| report_type  | string | 是   | 无     | 报表类型，自定义为custom，内置图表有特定名称 | Report type, customized to custom, built in chart with a specific name |
| name         | string | 是   | 无     | 统计图表的名称                               | The name of the chart                                                  |
| bk_obj_id    | string | 是   | 无     | 对象模型id                                   | the object identifier                                                  |
| chart_type   | string | 是   | 无     | 图表类型                                     | chart type                                                             |
| field        | string | 是   | 无     | 统计字段                                     | statistical field                                                      |
| width        | string | 是   | 无     | 表格宽度                                     | chart width                                                            |
| x_axis_count | int    | 是   | 10     | x轴显示数量                                  | X-axis display quantity                                                |




- output:
```
{
    bk_error_code: 0
    bk_error_msg: "success"
    data: {
        count: 1
        info: {
                bk_obj_id: "host"
                bk_supplier_account: "0"
                chart_type: "pie"
                config_id: 16
                create_time: "2019-08-09T16:04:20.145+08:00"
                field: "bk_sla"
                metadata: {label: null}
                name: "123"
                report_type: "custom"
                width: "50"
                x_axis_count: 10
            }
    }
    permission: null
    result: true
}
```

- output字段说明：

  | 字段 | 类型                | 说明                             | Description |
  | ---- | ------------------- | -------------------------------- | ----------- |
  | data | int  创建时记录的ID | the id of the target data record |


##### 删除统计报表
- API: DELETE /api/v3/delete/operation/chart/{id}

- API名称：delete_operation_chart

- 功能说明：

  中文：删除统计报表

  English：delete operation chart
  
- input body：

  无

- input字段说明：

  | 字段 | 类型 | 必填 | 默认值 | 说明                 | Description                      |
  | ---- | ---- | ---- | ------ | -------------------- | -------------------------------- |
  | id   | int  | 是   | 无     | 被删除的数据记录的ID | the id of the target data record |

- output：
```
  {
      "bk_error_code": 0
      "bk_error_msg": "success"
      "data": null
      "permission": null
      "result": true
  }
```

**修改统计报表**

- API：POST /api/v3/update/operation/chart

- API名称：update_operation_chart

- 功能说明：
    中文：修改统计报表
    English：update operation chart

- input body:
  
  ```
  {
      "bk_obj_id": "bk_switch"
      "bk_supplier_account": "0"
      "chart_type": "bar"
      "config_id": 9
      "field": "bk_biz_status"
      "metadata": {label: null}
      "name": "交换机"
      "report_type": "custom"
      "width": "50"
      "x_axis_count": 10
  }
  ```
  
  
  
- input字段说明：

| 名称         | 类型   | 必填 | 默认值 | 说明                                         | Description                                                            |
| ------------ | ------ | ---- | ------ | -------------------------------------------- | ---------------------------------------------------------------------- |
| report_type  | string | 是   | 无     | 报表类型，自定义为custom，内置图表有特定名称 | Report type, customized to custom, built in chart with a specific name |
| name         | string | 是   | 无     | 统计图表的名称                               | The name of the chart                                                  |
| bk_obj_id    | string | 是   | 无     | 对象模型id                                   | the object identifier                                                  |
| chart_type   | string | 是   | 无     | 图表类型                                     | chart type                                                             |
| field        | string | 是   | 无     | 统计字段                                     | statistical field                                                      |
| width        | string | 是   | 无     | 表格宽度                                     | chart width                                                            |
| x_axis_count | int    | 是   | 无     | x轴显示数量                                  | X-axis display quantity                                                |


- output:
```
{
    "bk_error_code": 0
    "bk_error_msg": "success"
    "data": 9
    "permission": null
    "result": true
}
```

**获取所有图表**

- API：GET /api/v3/search/operation/chart

- API名称：get all the charts

- 功能说明：
    中文：获取所有正在统计的图表
    English：get all the charts

- input body：
```
{}
```

- input字段说明：

- output：
```
bk_error_code: 0
bk_error_msg: "success"
data: {
	count: 9
	info: {
		host: [
			0: {
                "bk_obj_id": ""
                "bk_supplier_account": "0"
                "chart_type": ""
                "config_id": 6
                "create_time": "2019-08-07T15:06:08.9+08:00"
                "field": ""
                "metadata": {label: null}
                "name": "主机数量变化趋势"
                "report_type": "host_change_biz_chart"
                "width": "100"
                "x_axis_count": 20
            }]
		inst: [
		0: {
            "bk_obj_id": ""
            "bk_supplier_account": "0"
            "chart_type": "bar"
            "config_id": 8
            "create_time": "2019-08-07T15:06:08.902+08:00"
            "field": ""
            "metadata": {label: null}
            "name": "实例变更统计"
            "report_type": "model_inst_change_chart"
            "width": "50"
            "x_axis_count": 10
		}]
		nav: [
		0: {
            "bk_obj_id": ""
            "bk_supplier_account": "0"
            "chart_type": ""
            "config_id": 1
            "create_time": "2019-08-07T15:06:08.893+08:00"
            "field": ""
            "metadata": {label: null}
            "name": ""
            "report_type": "biz_module_host_chart"
            "width": ""
            "x_axis_count": 0
		}]
	}
}
permission: null
result: true
```

- output字段说明：

| 名称 | 类型         | 说明                                              | Description                                                                                      |
| ---- | ------------ | ------------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| host | object array | 主机统计图表的列表，列表中的index代表前端展示顺序 | List of host statistics charts, the index in the list represents the front-end display order     |
| inst | object array | 实例统计图表的列表，列表中的index代表前端展示顺序 | List of instance statistics charts, the index in the list represents the front-end display order |
| nav  | object array | 统计报表上方的导航栏                              | Navigation bar above the statistics report                                                       |


**获取统计图表数据**

- API：POST /api/v3/search/operation/chart/data

- API名称：get_statistics_chart_data

- 功能说明：
    中文：获取统计图表数据
    English：get statistical chart data

- input body：
  ```
  {
    "config_id": 1
  }
  ```
- input字段说明：

| 名称 | 类型 | 必填 | 默认值 | 说明     | Descripti             |
| ---- | ---- | ---- | ------ | -------- | --------------------- |
| id   | int  | 是   | 无     | 图表的ID | the object identifier |

- output：
```
{
    "bk_error_code": 0
    "bk_error_msg": "success"
    "data": [
        0: {
        	count: 0
       		id: "蓝鲸"
       	},
        1: {
        	count: 1
       		id: "测试业务"
       }
    ]
    permission: null
    result: true
}
```

- output字段说明

| 名称  | 类型   | 说明                                             | Description                                                                                                           |
| ----- | ------ | ------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------- |
| id    | string | 统计维度的具体值，例如：按省份统计中，id为各省名 | The specific value of the statistical dimension, for example: by province statistics, id is the name of each province |
| count | int    | 统计的数值                                       | Statistical value                                                                                                     |

**更新图表位置信息**

- API：POST /api/v3/update/operation/chart/position

- API名称：update_chart_position

- 功能说明：
    中文：更新图表位置信息
    English：update chart position

- input body：
  ```
  {
      "position":
      {
          "host": [4, 6, 3, 15]
          "inst": [8, 7, 9]
     }
  }
  ```
- input字段说明：

| 名称 | 类型      | 必填 | 默认值 | 说明                         | Descripti                                                   |
| ---- | --------- | ---- | ------ | ---------------------------- | ----------------------------------------------------------- |
| host | int array | 是   | 无     | 主机统计图表展示顺序的ID数组 | Host statistic chart showing the ID array of the order      |
| int  | int array | 是   | 无     | 实例统计图表展示顺序的ID数组 | Instance statistics chart showing the sequence of ID arrays |

- output：
```
{
    "bk_error_code": 0
    "bk_error_msg": "success"
    "data": null
    "permission": null
    "result": true
}
```

- output字段说明
