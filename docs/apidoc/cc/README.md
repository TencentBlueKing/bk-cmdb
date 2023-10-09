# cc
存放cmdb在esb的接口文档与映射关系。

### en
英文文档目录。

### zh_hans
中文文档目录。

### cc.yaml
存放cmdb与esb的接口映射关系。
 ```
举例：

  path: /v2/cc/list_host_total_mainline_topo/
  name: list_host_total_mainline_topo
  label: 查询主机及其对应topo
  label_en: query host and its corresponding Topo.
  suggest_method: POST
  api_type: operate
  comp_codename: generic.v2.cc.cc_component
  dest_path: /api/v3/findmany/hosts/total_mainline_topo/biz/{bk_biz_id}
  dest_http_method: POST

 ```
  
| 字段                  | 必选	 |	   描述                |
|----------------------|---------------------|---------------------|
|        path     |是	 |	      在esb处的映射地址     |
|       name      |  是	 |	   在esb处的接口名称，此名称需与接口文档名称一致     |
|       label      | 是	 |	    在esb处的中文接口描述     |
|       label_en      | 是	 |	    在esb处的英文接口描述     |
|    suggest_method         |  是	 |	   在esb处的请求调用方法     |
|      api_type       | 是	 |	     操作类型    |
|      comp_codename       |是	 |	  组建名称        |
|      dest_path       | 是	 |	   对应cmdb的映射地址，此地址是cmdb的接口地址      |
|      dest_http_method       |  是	 |	   对应cmdb的接口请求方法     |
|      is_hidden       |  否	 |	   接口在esb处是否隐藏，若为true，表示可以通过esb调用cmdb的接口，但是在esb的接口文档页面看不到此接口。默认值为false不隐藏     |

- 可根据esb与cmdb的映射关系，找到对应接口的文档。
- 当新增接口时，若需要通过esb进行调用，那么添加完中英文文档之后，还需要在cc.yaml文件中增加esb到cmdb接口的映射关系。
