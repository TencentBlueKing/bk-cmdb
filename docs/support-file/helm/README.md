# BK-CMDB

蓝鲸配置平台（蓝鲸CMDB）是一个面向资产及应用的企业级配置管理平台，本文档内容为如何在 Kubernetes 集群上部署 BK-CMDB 服务。

CMDB分为backend服务以及web服务，web服务提供页面访问功能，通过apiserver调用backend服务提供接口功能；

实际部署时请先部署backend服务，待backend服务完全启动后再部署web服务，服务部署文档：
- [backend服务](backend/README.md)
- [web服务](web/README.md)