# apidoc
存放cmdb相关的接口文档。

### cc
存放cmdb在API Gateway的接口文档与映射关系。

### esb-cc 
早期由ESB维护的组件文档，ESB封装的CMDB组件，此文件夹里的内容需要保留，请勿往里面添加任何新的接口。

### cc
存放cmdb在esb的接口文档与映射关系。

### archived
此文件夹保存归档的接口，将没有上esb的接口或者从esb下线接口文档放到这里。

#### 如何找到接口对应的文档：
1. 先到cc文件夹中，根据cc.yaml文件中的dest_path和dest_http_method找到cc对应的接口
2. 根据1中找到的name或path，确定cc在esb接口文档

注：
1. 这些接口文档是cmdb在esb的接口文档，与cmdb的实际接口有略微差异，如：一些参数原本在cc请求url中，在esb处是放在了请求体里；
2. 关于cc.yaml文件中的参数含义，可看cc文件夹中的README.md。
3. esb为蓝鲸PaaS平台中的蓝鲸API网关: [github地址](https://github.com/TencentBlueKing/legacy-bk-paas)
4. API Gateway为蓝鲸PaaS平台中的蓝鲸API网关: [github地址](https://github.com/TencentBlueKing/blueking-apigateway)
