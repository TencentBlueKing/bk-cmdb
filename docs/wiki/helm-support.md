# 使用helm快速体验cmdb
为了能让大家最快的体验蓝鲸bk-cmdb，本文推出了利用社区提供的`ccr.ccs.tencentyun.com/bk.io/cmdb-standalone:latest`镜像，使用helm部署的方式，快速部署体验蓝鲸配置平台服务。

## 环境要求

- Kubernetes 1.20+
- Helm 3+

## 部署方式
1. 拉取cmdb代码到本地(或下载源码到本地解压)：
    ```shell
    cd /data/
    git clone https://github.com/TencentBlueKing/bk-cmdb.git
    ```
   
2. 修改容器运行的node节点IP，执行：
    ```shell
    cd /data/bk-cmdb/docs/wiki/
    vi cmdb-standalone/values.yaml
    ```
    修改`values.yaml`文件中`hostIP: 127.0.0.1`的内容为`hostIP: ${容器运行时node节点IP}`

3. 运行服务，在`/data/bk-cmdb/docs/wiki/`目录下执行：
    ```shell
    helm install cmdb-standalone ./cmdb-standalone
    ```
   安装完成后执行`kubectl get pods`查看对应的pod状态，待pod状态为Running即为服务运行成功

4. 打开浏览器，访问`http://${容器运行时node节点IP}:8090`即可体验最新版的蓝鲸配置平台服务，用户名/密码：bk-cmdb/blueking
