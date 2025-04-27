#!/usr/bin/python
# -*- coding: utf-8 -*-

import sys
import getopt
import os, errno
import shutil
from string import Template


class FileTemplate(Template):
    delimiter = '$'

def mkdir_p(path):
    try:
        os.makedirs(path)
    except OSError as exc: # Python >2.5 (except OSError, exc: for Python <2.5)
        if exc.errno == errno.EEXIST and os.path.isdir(path):
            pass
        else: raise

def generate_config_file(
        rd_server_v, rd_cafile_v, rd_certfile_v, rd_keyfile_v, rd_skipverify_v, rd_certpassword_v,
        db_name_v, redis_ip_v, redis_port_v,
        redis_pass_v, sentinel_pass_v, redis_certfile_v, redis_keyfile_v, redis_cafile_v, redis_skipverify_v,
        mongo_ip_v, mongo_port_v, mongo_user_v, mongo_pass_v, mongo_certfile_v, mongo_keyfile_v, mongo_cafile_v, mongo_skipverify_v, rs_name, user_info,
        cc_url_v, paas_url_v, full_text_search, es_url_v, es_user_v, es_pass_v,es_shard_num_v,es_replica_num_v, auth_address, auth_app_code,
        auth_app_secret, auth_enabled, auth_scheme, auth_sync_workers, auth_sync_interval_minutes, log_level, register_ip,
        enable_cryptor_v, secret_key_url_v, secrets_addrs_v, secrets_token_v, secrets_project_v, secrets_env_v
):
    output = os.getcwd() + "/cmdb_adminserver/configures/"
    context = dict(
        db=db_name_v,
        mongo_user=mongo_user_v,
        mongo_host=mongo_ip_v,
        mongo_pass=mongo_pass_v,
        mongo_port=mongo_port_v,
        mongo_certfile=mongo_certfile_v,
        mongo_keyfile=mongo_keyfile_v,
        mongo_cafile=mongo_cafile_v,
        mongo_skipverify=mongo_skipverify_v,
        redis_host=redis_ip_v,
        redis_pass=redis_pass_v,
        sentinel_pass=sentinel_pass_v,
        redis_port=redis_port_v,
        redis_certfile=redis_certfile_v,
        redis_keyfile=redis_keyfile_v,
        redis_cafile=redis_cafile_v,
        redis_skipverify=redis_skipverify_v,
        cc_url=cc_url_v,
        paas_url=paas_url_v,
        es_url=es_url_v,
        es_user=es_user_v,
        es_pass=es_pass_v,
        es_shard_num=es_shard_num_v,
        es_replica_num=es_replica_num_v,
        ui_root="../web",
        agent_url=paas_url_v,
        configures_dir=output,
        rd_server=rd_server_v,
        rd_cafile=rd_cafile_v,
        rd_certfile=rd_certfile_v,
        rd_keyfile=rd_keyfile_v,
        rd_skipverify=rd_skipverify_v,
        rd_certpassword=rd_certpassword_v,
        auth_address=auth_address,
        auth_app_code=auth_app_code,
        auth_app_secret=auth_app_secret,
        auth_enabled=auth_enabled,
        auth_scheme=auth_scheme,
        auth_sync_workers=auth_sync_workers,
        auth_sync_interval_minutes=auth_sync_interval_minutes,
        full_text_search=full_text_search,
        rs_name=rs_name,
        user_info=user_info,
        enable_cryptor = enable_cryptor_v,
        secret_key_url = secret_key_url_v,
        secrets_addrs = secrets_addrs_v,
        secrets_token = secrets_token_v,
        secrets_project = secrets_project_v,
        secrets_env = secrets_env_v,
    )
    if not os.path.exists(output):
        os.mkdir(output)



    #redis.yaml
    redis_file_template_str = '''
#redis:
#  host: 127.0.0.1:6379
#  pwd: "123456"
#  database: "0"
#  maxOpenConns: 3000
#  maxIDleConns: 1000
#  snap:
#    host: 127.0.0.1:6379
#    pwd: 123456
#    database: "0"
#  discover:
#    host: 127.0.0.1:6379
#    pwd: 123456
#    database: "0"
#  netcollect:
#    host: 127.0.0.1:6379
#    pwd: 123456
#    database: "0"
redis:
  #公共redis配置信息,用于存取缓存，用户信息等数据
  host: $redis_host:$redis_port
  pwd: "$redis_pass"
  sentinelPwd: "$sentinel_pass"
  database: "0"
  maxOpenConns: 3000
  maxIDleConns: 1000
  #TLS配置信息
  tls:
    #证书文件路径
    certFile: "$redis_certfile"
    #密钥文件路径
    keyFile: "$redis_keyfile"
    #CA证书文件路径
    caFile: "$redis_cafile"
    #是否跳过证书验证
    insecureSkipVerify: $redis_skipverify
  #以下几个redis配置为datacollection模块所需的配置,用于接收第三方提供的数据
  #接收主机信息数据的redis
  snap:
    host: $redis_host:$redis_port
    pwd: "$redis_pass"
    sentinelPwd: "$sentinel_pass"
    database: "0"
    tls:
      certFile: "$redis_certfile"
      keyFile: "$redis_keyfile"
      caFile: "$redis_cafile"
      insecureSkipVerify: $redis_skipverify
  #接收模型实例数据的redis
  discover:
    host: $redis_host:$redis_port
    pwd: "$redis_pass"
    sentinelPwd: "$sentinel_pass"
    database: "0"
    tls:
      certFile: "$redis_certfile"
      keyFile: "$redis_keyfile"
      caFile: "$redis_cafile"
      insecureSkipVerify: $redis_skipverify
  #接受硬件数据的redis
  netcollect:
    host: $redis_host:$redis_port
    pwd: "$redis_pass"
    sentinelPwd: "$sentinel_pass"
    database: "0"
    tls:
      certFile: "$redis_certfile"
      keyFile: "$redis_keyfile"
      caFile: "$redis_cafile"
      insecureSkipVerify: $redis_skipverify
    '''

    template = FileTemplate(redis_file_template_str)
    result = template.substitute(**context)
    with open(output + "redis.yaml", 'w') as tmp_file:
        tmp_file.write(result)

    # mongodb.yaml
    mongodb_file_template_str = '''
#mongodb:
#  host: 127.0.0.1
#  port: 27017
#  usr: cc
#  pwd: cc
#  database: cmdb
#  maxOpenConns: 3000
#  maxIdleConns: 100
#  mechanism: SCRAM-SHA-1
#  rsName: rs0

# mongodb配置
mongodb:
  host: $mongo_host
  port: $mongo_port
  usr: $mongo_user
  pwd: "$mongo_pass"
  database: $db
  maxOpenConns: 3000
  maxIdleConns: 100
  mechanism: SCRAM-SHA-1
  rsName: $rs_name
  #mongo的socket连接的超时时间，以秒为单位，默认10s，最小5s，最大30s。
  socketTimeoutSeconds: 10
  #TLS配置信息
  tls:
    #证书文件路径
    certFile: "$mongo_certfile"
    #密钥文件路径
    keyFile: "$mongo_keyfile"
    #CA证书文件路径
    caFile: "$mongo_cafile"
    #是否跳过证书验证
    insecureSkipVerify: $mongo_skipverify
  # mongodb事件监听存储事件链的mongodb配置
watch:
  host: $mongo_host
  port: $mongo_port
  usr: $mongo_user
  pwd: "$mongo_pass"
  database: $db
  maxOpenConns: 10
  maxIdleConns: 5
  mechanism: SCRAM-SHA-1
  rsName: $rs_name
  socketTimeoutSeconds: 10
  tls:
    certFile: "$mongo_certfile"
    keyFile: "$mongo_keyfile"
    caFile: "$mongo_cafile"
    insecureSkipVerify: $mongo_skipverify
    '''
    template = FileTemplate(mongodb_file_template_str)
    result = template.substitute(**context)
    with open(output + "mongodb.yaml", 'w') as tmp_file:
        tmp_file.write(result)

    outputMonstache = os.getcwd() + "/monstache/etc/"
    if not os.path.exists(outputMonstache):
        mkdir_p(outputMonstache)
    # monstache.so config.toml
    mongo_tls_params = ""
    if mongo_certfile_v or mongo_keyfile_v or mongo_cafile_v:
        mongo_tls_params = "?tls=true"
        if mongo_cafile_v:
            mongo_tls_params += "&tlsCAFile=" + mongo_cafile_v
        if mongo_skipverify_v:
            mongo_tls_params += "&tlsInsecure=" + mongo_skipverify_v
    context['mongo_tls_params'] = mongo_tls_params
    
    monstachesoconfig_file_template_str = '''
# mongodb settings
mongo-url = "mongodb://$mongo_user:$mongo_pass@$mongo_host:$mongo_port/$db$mongo_tls_params"

# elasticsearch settings
elasticsearch-urls = ["$es_url"]
elasticsearch-user = "$es_user"
elasticsearch-password = "$es_pass"
gzip = true

# metadata collections.
change-stream-namespaces = [""]
direct-read-namespaces = [""]
direct-read-dynamic-include-regex = "cmdb.cc_ApplicationBase$$|cc_SetBase$$|cc_ModuleBase$$|cmdb.cc_HostBase$$|cmdb.cc_ObjDes$$|cc_ObjAttDes$$|cmdb.cc_ObjectBase_(.*)_pub_"
namespace-regex = "cmdb.cc_ApplicationBase$$|cc_SetBase$$|cc_ModuleBase$$|cmdb.cc_HostBase$$|cmdb.cc_ObjDes$$|cc_ObjAttDes$$|cmdb.cc_ObjectBase_(.*)_pub_"

# plugin
mapper-plugin-path = "etc/monstache-plugin.so"

# resume mode
resume = true
    '''
    template = FileTemplate(monstachesoconfig_file_template_str)
    result = template.substitute(**context)
    with open(outputMonstache + "config.toml", 'w') as tmp_file:
        tmp_file.write(result)

    # monstache.so extra.toml.toml
    monstachesoextra_file_template_str = '''
# elasticsearch settings

# the param must be assigned
elasticsearch-shard-num = "$es_shard_num"
elasticsearch-replica-num = "$es_replica_num"
    '''
    template = FileTemplate(monstachesoextra_file_template_str)
    result = template.substitute(**context)
    with open(outputMonstache + "extra.toml", 'w') as tmp_file:
        tmp_file.write(result)

    # common.yaml
    common_file_template_str = '''
#topoServer:
#  es:
#    fullTextSearch: "off"
#    url: http://127.0.0.1:9200
#    usr: cc
#    pwd: cc
#webServer:
#  api:
#    version: v3
#  session:
#    name: cc3
#    defaultlanguage: zh-cn
#    multipleOwner: 0
#    userInfo: cc:cc
#  site:
#    domainUrl: http://127.0.0.1:80/
#    bkLoginUrl: http://127.0.0.1/login/?app_id=%s&amp;c_url=%s
#    appCode: cc
#    checkUrl: http://127.0.0.1/login/accounts/get_user/?bk_token=
#    bkAccountUrl: http://127.0.0.1/login/accounts/get_all_user/?bk_token=%s
#    resourcesPath: /tmp/
#    htmlRoot: /data/cmdb/web
#    fullTextSearch: off
#  app:
#    agentAppUrl: http://127.0.0.1/console/?app=bk_agent_setup
#    authscheme: internal
#  login:
#    version: opensource
#operationServer:
#  timer:
#    spec: 00:30
#authServer:
#  address: 127.0.0.1
#  appCode: bk_cmdb
#  appSecret: 123456
#cloudServer:
#  cryptor:
#    enableCryptor: false
#    secretKeyUrl:
#    secretsAddrs:
#    secretsToken:
#    secretsProject:
#    secretsEnv:

#elasticsearch配置
es:
  #全文检索功能开关(取值：off/on)，默认是off，开启是on
  fullTextSearch: "$full_text_search"
  #elasticsearch服务监听url，默认是[http://127.0.0.1:9200](http://127.0.0.1:9200/)
  url: $es_url
  #用户
  usr: $es_user
  #密码
  pwd: $es_pass
# adminServer专属配置
adminServer:
  #同步IAM动态模型的周期,单位为分钟，最小为1分钟,默认为5分钟
  syncIAMPeriodMinutes: 5

# operation_server专属配置
operationServer:
  timer:
    # 00:00-23:59,运营统计定时收集数据时间点,默认是为00:30
    spec: 00:30  # 00:00 - 23:59
  # 禁用运营统计数据统计功能，默认false
  disableOperationStatistic: false

#auth_server专属配置
authServer:
  #蓝鲸权限中心地址,可配置多个,用,(逗号)分割
  address: $auth_address
  #cmdb项目在蓝鲸权限中心的应用编码
  appCode: $auth_app_code
  #cmdb项目在蓝鲸权限中心的应用密钥
  appSecret: $auth_app_secret
  authCenter:
     # 权限中心tls配置
     tls:
       # 客户端是否验证服务端证书，包含证书链和主机名，bool值, true为不校验, false为校验
       insecureSkipVerify:
       # 服务使用的证书的路径,如:/data/cmdb/cert/server.crt
       certFile:
       # 服务使用的证书对应的密钥的路径,如:/data/cmdb/cert/server.key
       keyFile:
       # CA证书路径，用于验证对方证书,如:/data/cmdb/cert/ca.crt
       caFile:
       # 用于解密根据RFC1423加密的证书密钥的PEM块
       password:
  # 空间级权限版本中，find_module_host_relation、find_host_by_service_template、find_host_by_set_template、
  # list_biz_hosts、list_biz_hosts_topo、find_host_by_topo、list_host_total_mainline_topo这几个上esb接口,
  # 可以配置不进行业务访问鉴权
  skipViewBizAuth: false

#cloudServer专属配置
cloudServer:
  # 加密服务使用
  cryptor:
    enableCryptor: ${enable_cryptor}
    secretKeyUrl: ${secret_key_url}
    secretsAddrs: ${secrets_addrs}
    secretsToken: ${secrets_token}
    secretsProject: ${secrets_project}
    secretsEnv: ${secrets_env}
  # 云同步任务
  syncTask:
    # 同步周期,最小为5分钟
    syncPeriodMinutes: 5

#datacollection专属配置
datacollection:
  hostsnap:
    # 主机静态数据采集模式，将数据导入kafka或者redis，可选值是 kafka、redis，默认值为redis（仅用于新插件bkmonitorbeat）
    reportMode: redis
    # 当主机快照数据属性,如cpu,bk_cpu_mhz,bk_disk,bk_mem这些数值型数据变动的范围大于该配置的值时，进行db数据的更新，默认值为10%，最小值为5%，以百分比为单位
    changeRangePercent: 10
    # 用于设置主机快照key在redis中的过期时间，该时间会有上下50%的波动，当key存在时，同一id的主机数据不会更新，默认值为10分钟，最小值为5分钟，以分钟为单位
    changeCountExpireMinute: 10
    # 用于设置单个服务主机快照处理请求能力，起到限流的作用，令牌桶配置，最多请求数通过burst设置。qps的默认值为40，burst的默认值为100
    rateLimiter:
      qps: 40
      burst: 100

# 监控配置， monitor配置项必须存在
monitor:
    # 监控插件名称，有noop，blueking， 不填时默认为noop
    pluginName: noop
    # 是否开启监控
    enableMonitor: false
    # 当使用blueking监控插件时，上报数据所需要的数据通道标识,如1500000
    dataID: 0
    # 采集数据后能够缓存的队列长度，设置范围为1～1000, 默认为100
    queueSize: 100
    # 当使用blueking监控插件时，上报数据的接收端URL
    bkMonitorReportUrl:
    # 当使用blueking监控插件时，上报数据所需要的数据通道标识验证码，防止数据误上报至其他管道
    accessToken:
    # 用于对数据上报进行频率控制和限流
    # qps的设置范围为1～50，默认值为10
    # burst的设置范围为1～100，默认值为20
    rateLimiter:
      qps: 10
      burst: 20

# openTelemetry跟踪链接入相关配置
openTelemetry:
  # 表示是否开启openTelemetry跟踪链接入相关功能，布尔值, 默认值为false不开启
  enable: false
  # openTelemetry跟踪链功能的自定义上报服务地址
  endpoint:
  # openTelemetry跟踪链功能的上报data_id, 用于旧版的数据上报，和bkDataToken同时配置时，会取bkDataToken
  bkDataID:
  # openTelemetry跟踪链功能的bk.data.token, 用于新版的数据上报，和bkDataID同时配置时，会取bkDataToken
  bkDataToken:
  # 服务tls配置
  tls:
    # 客户端是否验证服务端证书，包含证书链和主机名，bool值, true为不校验, false为校验
    insecureSkipVerify:
    # 服务使用的证书的路径,如:/data/cmdb/cert/server.crt
    certFile:
    # 服务使用的证书对应的密钥的路径,如:/data/cmdb/cert/server.key
    keyFile:
    # CA证书路径，用于验证对方证书,如:/data/cmdb/cert/ca.crt
    caFile:
    # 用于解密根据RFC1423加密的证书密钥的PEM块
    password:

# eventServer相关配置
eventServer:
  # 下发主机身份相关配置
  hostIdentifier:
    # 是否开始下发主机身份功能, 有两个值，true和false，当处于true时，开启下发主机身份功能，false时，关闭该功能
    startUp: false
    # 可选值为v1或v2, v1表示采用的是使用thrift方式连接gse，需要配置gse.apiServer和gse.taskServer
    # v2表示使用的是apigw接口与gse交互，此接口在gse2.0提供，需要配置gse.apiGW
    version: v1
    # 每隔多少个小时进行一次全量主机身份批量的同步操作，整数值，单位为小时，注：刚启动服务时，会等一个周期后再进行全量同步操作
    batchSyncIntervalHours: 6
    # 用于设置推送主机身份请求gse的taskServer能力，起到限流的作用。qps的默认值为200, 代表每秒最多推送的主机数量，burst的默认值为200
    rateLimiter:
      qps: 200
      burst: 200
    # 下发主机身份文件名
    fileName: "hostid"
    # 当下发主机为linux操作系统时，相关配置
    linux:
      # 下发主机身份文件路径
      filePath: "/var/lib/gse/host"
      # 下发主机身份文件所有者
      fileOwner: "root"
      # 下发主机身份文件权限值
      filePrivilege: 644
    # 当下发主机为windows操作系统时，相关配置
    windows:
      # 下发主机身份文件路径
      filePath: "c:/gse/data/host"
      # 下发主机身份文件所有者
      fileOwner: "root"
      # 下发主机身份文件权限值
      filePrivilege: 644

# apiServer相关配置
apiServer:
  # api-server使用的jwt配置
  jwt:
    # 是否开启jwt认证功能
    enabled: false
    # jwt公钥
    publicKey:

# 直接调用gse服务相关配置
gse:
  # 调用gse的apiServer服务时相关配置
  apiServer:
    # 此配置为数组类型，可配置连接gse的apiServer的多个host:port格式的值，去建立连接
    endpoints:
      # 证书相关信息
    insecureSkipVerify: true
    certFile:
    keyFile:
    caFile:
    password:
  # 调用gse的taskServer服务时相关配置
  taskServer:
    # 此配置为数组类型，可配置连接gse的taskServer的多个host:port格式的值，去建立连接
    endpoints:
    # 证书相关信息
    insecureSkipVerify: true
    certFile:
    keyFile:
    caFile:
    password:

# 当主机静态数据采集模式为kafka时，datacollection处理插件bkmonitorbeat采集上来的主机静态数据，选择kafka作为数据导入组件时的相关配置
# 此配置与migrate.yaml文件中的reportMode相关联
kafka:
  snap:
    brokers:
    # groupID为固定值，请勿随便修改，修改后会导致重复消费过去的数据
    groupID: bk_cmdb_snapshot_group
    # partition数量固定为1，保证消息的顺序性
    partition: 1
    # 安全协议SASL_PLAINTEXT，SASL机制SCRAM-SHA-512的账号、密码信息
    user:
    password:

# cmdb服务tls配置
tls:
  # 客户端是否验证服务端证书，包含证书链和主机名，bool值, true为不校验, false为校验
  insecureSkipVerify:
  # 服务使用的证书的路径,如:/data/cmdb/cert/server.crt
  certFile:
  # 服务使用的证书对应的密钥的路径,如:/data/cmdb/cert/server.key
  keyFile:
  # CA证书路径，用于验证对方证书,如:/data/cmdb/cert/ca.crt
  caFile:
  # 用于解密根据RFC1423加密的证书密钥的PEM块
  password:

# 调用apigw接口相关配置
apiGW:
  bkGseApiGatewayUrl:
  bkNoticeApiGatewayUrl:
  bkCmdbApiGatewayUrl:
  appCode:
  appSecret:
  username:
  tls:
    # 客户端是否验证服务端证书，包含证书链和主机名，bool值, true为不校验, false为校验
    insecureSkipVerify:
    # 服务使用的证书的路径,如:/data/cmdb/cert/server.crt
    certFile:
    # 服务使用的证书对应的密钥的路径,如:/data/cmdb/cert/server.key
    keyFile:
    # CA证书的路径，用于验证对方证书,如:/data/cmdb/cert/ca.crt
    caFile:
    # 用于解密根据RFC1423加密的证书密钥的PEM块
    password:

# CMDB同步服务配置
transferService:
  # 是否开启同步, 默认为false
  enableSync: false
  # 是否开启增量同步, 默认为false
  enableIncrSync: false
  # 同步服务的名称
  name:
  # 同步角色，src表示源环境(将数据同步到目标环境)，dest表示目标环境(接收源环境的同步数据)
  role:
  # 全量同步周期，单位：小时，仅源环境需要配置
  syncIntervalHours:
  # 传输介质地址
  transferMediumAddress:
    - transfer.example.com
    '''

    template = FileTemplate(common_file_template_str)
    result = template.substitute(**context)
    with open(output + "common.yaml", 'w') as tmp_file:
        tmp_file.write(result)

    # extra.yaml
    extra_file_template_str = ''

    template = FileTemplate(extra_file_template_str)
    result = template.substitute(**context)
    with open(output + "extra.yaml", 'w') as tmp_file:
        tmp_file.write(result)


    # web.yaml
    web_file_template_str = '''
#elasticsearch配置
es:
  #全文检索功能开关(取值：off/on)，默认是off，开启是on
  fullTextSearch: "$full_text_search"
  #elasticsearch服务监听url，默认是[http://127.0.0.1:9200](http://127.0.0.1:9200/)
  url: $es_url
  #用户
  usr: $es_user
  #密码
  pwd: $es_pass
# 指定errors的路径
errors:
  res: conf/errors
# 指定language的路径
language:
  res: conf/language
# web_server专属配置
webServer:
  api:
    #显示版本，比如v3为3.x
    version: v3
  #会话相关
  session:
    #会话名
    name: cc3
    #语言
    defaultlanguage: zh-cn
    #是否支持同时登录同一用户，0为不支持，1为支持
    multipleOwner: "0"
    #账号密码，以 : 分割
    userInfo: $user_info
  site:
    #该值表示部署完成后,输入到浏览器中访问的cmdb 网址
    domainUrl: ${cc_url}
    #登录地址
    bkLoginUrl: ${paas_url}/login/?app_id=%s&c_url=%s
    appCode: cc
    checkUrl: ${paas_url}/login/accounts/get_user/?bk_token=
    bkAccountUrl: ${paas_url}/login/accounts/get_all_user/?bk_token=%s
    resourcesPath: /tmp/
    #前端基础页面位置
    htmlRoot: $ui_root
    #帮助文档地址
    helpDocUrl: https://bk.tencent.com/docs/markdown/配置平台/产品白皮书/产品简介/Overview.md
    paas:
      # pass的tls相关配置
      tls:
        # 客户端是否验证服务端证书，包含证书链和主机名，bool值, true为不校验, false为校验
        insecureSkipVerify:
        # 服务使用的证书的路径,如:/data/cmdb/cert/server.crt
        certFile:
        # 服务使用的证书对应的密钥的路径,如:/data/cmdb/cert/server.key
        keyFile:
        # CA证书路径，用于验证对方证书,如:/data/cmdb/cert/ca.crt
        caFile:
        # 用于解密根据RFC1423加密的证书密钥的PEM块
        password:
  app:
    agentAppUrl: ${agent_url}/console/?app=bk_agent_setup
    #权限模式，web页面使用，可选值: internal, iam
    authscheme: $auth_scheme
  login:
    #登录模式
    version: $loginVersion
  #cmdb版本日志存放路径配置
  changelogPath:
    #中文版版本日志存放路径
    ch: ../changelog_user/ch
    #英文版版本日志存放路径
    en: ../changelog_user/en
  # web-server使用的jwt配置
  jwt:
    # 是否开启jwt认证功能
    enabled: false
    # jwt公钥
    publicKey:
    # jwt私钥
    privateKey:

# 监控配置， monitor配置项必须存在
monitor:
    # 监控插件名称，有noop，blueking， 不填时默认为noop
    pluginName: noop
    # 是否开启监控
    enableMonitor: false
    # 当使用blueking监控插件时，上报数据所需要的数据通道标识,如1500000
    dataID: 0
    # 采集数据后能够缓存的队列长度，设置范围为1～1000, 默认为100
    queueSize: 100
    # 采集数据用的gsecmdline命令绝对路径，默认版本需要大于等于2.0.2 ，默认路径：/usr/local/gse/plugins/bin/gsecmdline
    gsecmdlinePath:
    # 对应的domainsocket绝对路径，默认路径：/usr/local/gse/agent/data/ipc.state.report
    domainSocketPath:
    # 用于对数据上报进行频率控制和限流
    # qps的设置范围为1～50，默认值为10
    # burst的设置范围为1～100，默认值为20
    rateLimiter:
      qps: 10
      burst: 20

# openTelemetry跟踪链接入相关配置
openTelemetry:
  # 表示是否开启openTelemetry跟踪链接入相关功能，布尔值, 默认值为false不开启
  enable: false
  # openTelemetry跟踪链功能的自定义上报服务地址
  endpoint:
  # openTelemetry跟踪链功能的上报data_id, 用于旧版的数据上报，和bkDataToken同时配置时，会取bkDataToken
  bkDataID:
  # openTelemetry跟踪链功能的bk.data.token, 用于新版的数据上报，和bkDataID同时配置时，会取bkDataToken
  bkDataToken:
  # 服务tls配置
  tls:
    # 客户端是否验证服务端证书，包含证书链和主机名，bool值, true为不校验, false为校验
    insecureSkipVerify:
    # 服务使用的证书的路径,如:/data/cmdb/cert/server.crt
    certFile:
    # 服务使用的证书对应的密钥的路径,如:/data/cmdb/cert/server.key
    keyFile:
    # CA证书路径，用于验证对方证书,如:/data/cmdb/cert/ca.crt
    caFile:
    # 用于解密根据RFC1423加密的证书密钥的PEM块
    password:

redis:
  #公共redis配置信息,用于存取缓存，用户信息等数据
  host: $redis_host:$redis_port
  pwd: "$redis_pass"
  sentinelPwd: "$sentinel_pass"
  database: "0"
  maxOpenConns: 3000
  maxIDleConns: 1000
  #TLS配置信息
  tls:
    #证书文件路径
    certFile: "$redis_certfile"
    #密钥文件路径
    keyFile: "$redis_keyfile"
    #CA证书文件路径
    caFile: "$redis_cafile"
    #是否跳过证书验证
    insecureSkipVerify: $redis_skipverify

mongodb:
  host: $mongo_host
  port: $mongo_port
  usr: $mongo_user
  pwd: "$mongo_pass"
  database: $db
  maxOpenConns: 3000
  maxIdleConns: 100
  mechanism: SCRAM-SHA-1
  rsName: $rs_name
  #mongo的socket连接的超时时间，以秒为单位，默认10s，最小5s，最大30s。
  socketTimeoutSeconds: 10
  tls:
    #证书文件路径
    certFile: "$mongo_certfile"
    #密钥文件路径
    keyFile: "$mongo_keyfile"
    #CA证书文件路径
    caFile: "$mongo_cafile"
    #是否跳过证书验证
    insecureSkipVerify: $mongo_skipverify
    '''

    template = FileTemplate(web_file_template_str)
    loginVersion = 'opensource'
    if auth_enabled == "true":
        loginVersion = 'blueking'
    result = template.substitute(loginVersion=loginVersion, **context)
    web_dir = os.getcwd() + "/cmdb_webserver/"
    with open(web_dir + "web.yaml", 'w') as tmp_file:
        tmp_file.write(result)

    # migrate.yaml
    migrate_file_template_str = '''
#configServer:
#  addrs: 127.0.0.1:2181
#  usr: cc
#  pwd: cc
#registerServer:
#  addrs: 127.0.0.1:2181
#  usr: cc
#  pwd: cc
#mongodb:
#  host: 127.0.0.1
#  port: 27017
#  usr: cc
#  pwd: cc
#  database: cmdb
#  maxOpenConns: 5
#  maxIdleConns: 1
#  mechanism: SCRAM-SHA-1
#  rsName: rs0
#redis:
#  host: 127.0.0.1:6379
#  pwd: 123456
#  database: "0"
#  maxOpenConns: 5
#  maxIDleConns: 1
#confs:
#  dir: /data/cmdb/cmdb_adminserver/configures/
#errors:
#  res: /data/cmdb/cmdb_adminserver/conf/errors
#language:
#  res: /data/cmdb/cmdb_adminserver/conf/language
#monstache:
#  dir: /data/cmdb/monstache/etc
#auth:
#  address: 127.0.0.1
#  appCode: bk_cmdb
#  appSecret: 123456

# 配置中心
configServer:
  addrs: $rd_server
  usr:
  pwd:
  # ZooKeeper tls配置信息
  tls:
    # CA证书文件路径
    caFile: "$rd_cafile"
    # 证书文件路径
    certFile: "$rd_certfile"
    # 密钥文件路径
    keyFile: "$rd_keyfile"
    # 是否跳过证书验证
    insecureSkipVerify: $rd_skipverify
    # 证书密码
    password: $rd_certpassword
# 注册中心
registerServer:
  addrs: $rd_server
  usr:
  pwd:
  # ZooKeeper tls配置信息
  tls:
    # CA证书文件路径
    caFile: "$rd_cafile"
    # 证书文件路径
    certFile: "$rd_certfile"
    # 密钥文件路径
    keyFile: "$rd_keyfile"
    # 是否跳过证书验证
    insecureSkipVerify: $rd_skipverify
    # 证书密码
    password: $rd_certpassword
# 指定configures的路径，通过这个路径找到其他的配置文件
confs:
  dir: $configures_dir
# 指定errors的路径
errors:
  res: conf/errors
# 指定language的路径
language:
  res: conf/language
# 指定monstache相关配置文件
monstache:
  res: monstache/etc
# migrate dataid时的相关配置
dataid:
  # 通过何种方式调用gse接口注册dataid, 可选值esb和apigw, 如果填写esb则必须配置common.yaml的esb配置, 如果填写apigw则必须配置common.yaml的apiGW配置,不填默认为esb
  migrateWay: esb
    '''

    template = FileTemplate(migrate_file_template_str)
    result = template.substitute(**context)
    with open(output + "migrate.yaml", 'w') as tmp_file:
        tmp_file.write(result)

def update_start_script(rd_server, server_ports, enable_auth, log_level, register_ip, enable_cryptor, rd_cafile, rd_certfile, rd_keyfile, rd_skipverify, rd_certpassword):
    list_dirs = os.walk(os.getcwd()+"/")
    for root, dirs, _ in list_dirs:
        for d in dirs:
            if not d.startswith("cmdb_"):
                continue

            if d == "cmdb_adminserver":
                if os.path.exists(root+d+"/init_db.sh"):
                    shutil.copy(root + d + "/init_db.sh", os.getcwd() + "/init_db.sh")

            target_file = root + d + "/start.sh"
            if not os.path.exists(target_file) or not os.path.exists(root+d + "/template.sh.start"):
                continue

            # Read in the file
            with open(root+d + "/template.sh.start", 'r') as template_file:
                filedata = template_file.read()
                # Replace the target string
                filedata = filedata.replace('cmdb-name-placeholder', d)
                filedata = filedata.replace('cmdb-port-placeholder', str(server_ports.get(d, 9999)))
                if d == "cmdb_adminserver":
                    filedata = filedata.replace('rd_server_placeholder', "configures/migrate.yaml")
                    filedata = filedata.replace('regdiscv', "config")
                elif d == "cmdb_webserver":
                    filedata = filedata.replace('rd_server_placeholder', rd_server + " --config=web.yaml")
                else:
                    filedata = filedata.replace('rd_server_placeholder', rd_server)

                extend_flag = ''
                if d in ['cmdb_apiserver', 'cmdb_hostserver', 'cmdb_datacollection', 'cmdb_procserver',
                         'cmdb_toposerver', 'cmdb_eventserver', 'cmdb_operationserver', 'cmdb_cloudserver',
                         'cmdb_authserver','cmdb_adminserver','cmdb_cacheservice']:
                    extend_flag += ' --enable-auth=%s ' % enable_auth
                if d in ['cmdb_cloudserver']:
                     extend_flag += ' --enable_cryptor=%s ' % enable_cryptor
                if register_ip != '':
                    extend_flag += ' --register-ip=%s ' % register_ip

                if d != "cmdb_adminserver":
                    if rd_cafile != '':
                        extend_flag += ' --regdiscv-cafile=%s ' % rd_cafile
                    if rd_certfile != '':
                        extend_flag += ' --regdiscv-certfile=%s ' % rd_certfile
                    if rd_keyfile != '':
                        extend_flag += ' --regdiscv-keyfile=%s ' % rd_keyfile
                    if rd_skipverify != '':
                        extend_flag += ' --regdiscv-skipverify=%s ' % rd_skipverify
                    if rd_certpassword != '':
                        extend_flag += " --regdiscv-certpassword=%s " % rd_certpassword

                filedata = filedata.replace('extend_flag_placeholder', extend_flag)

                filedata = filedata.replace('log_level_placeholder', log_level)

                # Write the file out again
                with open(target_file, 'w') as new_file:
                    new_file.write(filedata)


def main(argv):
    db_name = 'cmdb'
    rd_server = ''
    redis_ip = ''
    redis_port = 6379
    redis_pass = ''
    sentinel_pass = ''
    redis_certfile = ''
    redis_keyfile = ''
    redis_cafile = ''
    redis_skipverify = 'true'
    mongo_ip = ''
    mongo_port = 27017
    mongo_user = ''
    mongo_pass = ''
    mongo_certfile = ''
    mongo_keyfile = ''
    mongo_cafile = ''
    mongo_skipverify = 'true'
    cc_url = ''
    paas_url = 'http://127.0.0.1'
    auth = {
        "auth_scheme": "internal",
        # iam options
        "auth_address": "",
        "auth_enabled": "false",
        "auth_app_code": "bk_cmdb",
        "auth_app_secret": "",
        "auth_sync_workers": "100",
        "auth_sync_interval_minutes": "45",
    }
    full_text_search = 'off'
    es_url = 'http://127.0.0.1:9200'
    es_user = ''
    es_pass = ''
    es_shard_num = 1
    es_replica_num = 1
    log_level = '3'
    register_ip = ''
    rs_name = 'rs0'
    user_info = ''
    enable_cryptor = 'false'
    secret_key_url = ''
    secrets_addrs = ''
    secrets_token = ''
    secrets_project = ''
    secrets_env = ''
    rd_cafile = ''
    rd_certfile = ''
    rd_keyfile = ''
    rd_skipverify = 'true'
    rd_certpassword = ''

    server_ports = {
        "cmdb_adminserver": 60004,
        "cmdb_apiserver": 8080,
        "cmdb_datacollection": 60005,
        "cmdb_eventserver": 60009,
        "cmdb_hostserver": 60001,
        "cmdb_coreservice": 50009,
        "cmdb_procserver": 60003,
        "cmdb_toposerver": 60002,
        "cmdb_webserver": 8083,
        "cmdb_synchronizeserver": 60010,
        "cmdb_operationserver": 60011,
        "cmdb_taskserver": 60012,
        "cmdb_cloudserver": 60013,
        "cmdb_authserver": 60014,
        "cmdb_cacheservice": 50010,
        "cmdb_transferservice": 50011
    }
    arr = [
        "help", "discovery=", "database=", "redis_ip=", "redis_port=",
        "redis_pass=", "sentinel_pass=", "redis_certfile=", "redis_keyfile=", "redis_cafile=", "redis_skipverify=",
        "mongo_ip=", "mongo_port=", "rs_name=",
        "mongo_user=", "mongo_pass=", "mongo_certfile=", "mongo_keyfile=", "mongo_cafile=", "mongo_skipverify=",
        "blueking_cmdb_url=", "user_info=",
        "blueking_paas_url=", "listen_port=", "es_url=", "es_user=", "es_pass=", "es_shard_num=","es_replica_num=","auth_address=",
        "auth_app_code=", "auth_app_secret=", "auth_enabled=",
        "auth_scheme=", "auth_sync_workers=", "auth_sync_interval_minutes=", "full_text_search=", "log_level=", "register_ip=",
        "enable_cryptor=", "secret_key_url=", "secrets_addrs=", "secrets_token=", "secrets_project=", "secrets_env=",
        "discovery_cafile=", "discovery_certfile=", "discovery_keyfile=", "discovery_skipverify=", "discovery_certpassword="
    ]
    usage = '''
    usage:
      --discovery               <discovery>               the ZooKeeper server address, eg:127.0.0.1:2181
      --database                <database>                the database name, default cmdb
      --redis_ip                <redis_ip>                the redis ip, eg:127.0.0.1
      --redis_port              <redis_port>              the redis port, default:6379
      --redis_pass              <redis_pass>              the redis user password
      --sentinel_pass           <sentinel_pass>           the redis sentinel password
      --redis_certfile          <redis_certfile>          the redis cert file path
      --redis_keyfile           <redis_keyfile>           the redis key file path
      --redis_cafile            <redis_cafile>            the redis ca cert file path
      --redis_skipverify        <redis_skipverify>        the redis skip verify
      --mongo_ip                <mongo_ip>                the mongo ip ,eg:127.0.0.1
      --mongo_port              <mongo_port>              the mongo port, eg:27017
      --mongo_user              <mongo_user>              the mongo user name, default:cc
      --mongo_pass              <mongo_pass>              the mongo password
      --mongo_certfile          <mongo_certfile>          the mongo cert file path
      --mongo_keyfile           <mongo_keyfile>           the mongo key file path
      --mongo_cafile            <mongo_cafile>            the mongo ca cert file path
      --mongo_skipverify        <mongo_skipverify>        the mongo skip verify
      --rs_name                 <rs_name>                 the mongo replica set name, default: rs0
      --blueking_cmdb_url       <blueking_cmdb_url>       the cmdb site url, eg: http://127.0.0.1:8088 or http://bk.tencent.com
      --blueking_paas_url       <blueking_paas_url>       the blueking paas url, eg: http://127.0.0.1:8088 or http://bk.tencent.com
      --listen_port             <listen_port>             the cmdb_webserver listen port, should be the port as same as -c <blueking_cmdb_url> specified, default:8083
      --auth_scheme             <auth_scheme>             auth scheme, ex: internal, iam
      --auth_enabled            <auth_enabled>            iam auth enabled, true or false
      --auth_address            <auth_address>            iam address
      --auth_app_code           <auth_app_code>           app code for iam, default bk_cmdb
      --auth_app_secret         <auth_app_secret>         app code for iam
      --full_text_search        <full_text_search>        full text search on or off
      --es_url                  <es_url>                  the es listen url, see in es dir config/elasticsearch.yml, (network.host, http.port), default: http://127.0.0.1:9200
      --es_user                 <es_user>                 the es user name
      --es_pass                 <es_pass>                 the es password
      --es_shard_num            <es_shard_num>            the es sharding num
      --es_replica_num          <es_replica_num>          the es es_replica_num
      --log_level               <log_level>               log level to start cmdb process, default: 3
      --register_ip             <register_ip>             the ip address registered on zookeeper, it can be domain
      --user_info               <user_info>               the system user info, user and password are combined by semicolon, multiple users are separated by comma. eg: user1:password1,user2:password2
      --enable_cryptor          <enable_cryptor>          enable cryptor,true or false, default is false
      --secret_key_url          <secret_key_url>          the url to get secret_key which used to encrypt and decrypt cloud account
      --secrets_addrs           <secrets_addrs>           secrets_addrs, the addrs of bk-secrets service, start with http:// or https://
      --secrets_token           <secrets_token>           secrets_token , as a header param for sending the api request to bk-secrets service
      --secrets_project         <secrets_project>         secrets_project, as a header param for sending the api request to bk-secrets service
      --secrets_env             <secrets_env>             secrets_env, as a header param for sending the api request to bk-secrets service
      --discovery_cafile        <discovery_cafile>        CA file for ZooKeeper TLS connection
      --discovery_certfile      <discovery_certfile>      cert file for ZooKeeper TLS connection
      --discovery_keyfile       <discovery_keyfile>       key file for ZooKeeper TLS connection
      --discovery_skipverify    <discovery_skipverify>    whether to skip verify ZooKeeper TLS connection, true or false, default true
      --discovery_certpassword  <discovery_certpassword>  password for ZooKeeper TLS connection

    demo:
    python init.py  \\
      --discovery               127.0.0.1:2181 \\
      --discovery_cafile        ./zk-ca.crt \\
      --discovery_certfile      ./zk.cert \\
      --discovery_keyfile       ./zk.key \\
      --discovery_skipverify    true \\
      --discovery_certpassword  password \\
      --database                cmdb \\
      --redis_ip                127.0.0.1 \\
      --redis_port              6379 \\
      --redis_pass              1111 \\
      --sentinel_pass           2222 \\
      --redis_certfile          ./redis.cert \\
      --redis_keyfile           ./redis.key \\
      --redis_cafile            ./redis-ca.crt \\
      --redis_skipverify        true \\
      --mongo_ip                127.0.0.1 \\
      --mongo_port              27017 \\
      --mongo_user              cc \\
      --mongo_pass              cc \\
      --mongo_certfile          ./mongo.cert \\
      --mongo_keyfile           ./mongo.key \\
      --mongo_cafile            ./mongo-ca.crt \\
      --mongo_skipverify        true \\
      --rs_name                 rs0 \\
      --blueking_cmdb_url       http://127.0.0.1:8080/ \\
      --blueking_paas_url       http://paas.domain.com \\
      --listen_port             8080 \\
      --auth_scheme             internal \\
      --auth_enabled            false \\
      --auth_address            https://iam.domain.com/ \\
      --auth_app_code           bk_cmdb \\
      --auth_app_secret         xxxxxxx \\
      --auth_sync_workers       1 \\
      --auth_sync_interval_     minutes  45 \\
      --full_text_search        off \\
      --es_url                  http://127.0.0.1:9200 \\
      --es_user                 cc \\
      --es_pass                 cc \\
      --es_shard_num            1 \\
      --es_replica_num          1 \\
      --log_level               3 \\
      --register_ip             cmdb.domain.com \\
      --user_info               user1:password1,user2:password2
    '''
    try:
        opts, _ = getopt.getopt(argv, "hd:D:r:p:x:s:m:P:X:S:u:U:a:l:es:v", arr)

    except getopt.GetoptError as e:
        print("\n \t", e.msg)
        print(usage)

        sys.exit(2)
    if len(opts) == 0:
        print(usage)
        sys.exit(2)

    for opt, arg in opts:
        if opt in ('-h', '--help'):
            print(usage)
            sys.exit()
        elif opt in ("-d", "--discovery"):
            rd_server = arg
            print('rd_server:', rd_server)
        elif opt in("--discovery_cafile",):
            rd_cafile = arg
            print('rd_cafile:', rd_cafile)
        elif opt in("--discovery_certfile",):
            rd_certfile = arg
            print('rd_certfile:', rd_certfile)
        elif opt in("--discovery_keyfile",):
            rd_keyfile = arg
            print('rd_keyfile:', rd_keyfile)
        elif opt in("--discovery_skipverify",):
            rd_skipverify = arg
            print('rd_skipverify:', rd_skipverify)
        elif opt in("--discovery_certpassword"):
            rd_certpassword = arg
            print('rd_certpassword:', rd_certpassword)
        elif opt in ("-D", "--database"):
            db_name = arg
            print('database:', db_name)
        elif opt in ("-r", "--redis_ip"):
            redis_ip = arg
            print('redis_ip:', redis_ip)
        elif opt in ("-p", "--redis_port"):
            redis_port = arg
            print('redis_port:', redis_port)
        elif opt in ("-s", "--redis_pass"):
            redis_pass = arg
            print('redis_pass:', redis_pass)
        elif opt in ("-s", "--sentinel_pass"):
            sentinel_pass = arg
            print('sentinel_pass:', sentinel_pass)
        elif opt in ("--redis_certfile",):
            redis_certfile = arg
            print('redis_certfile:', redis_certfile)
        elif opt in ("--redis_keyfile",):
            redis_keyfile = arg
            print('redis_keyfile:', redis_keyfile)
        elif opt in ("--redis_cafile",):
            redis_cafile = arg
            print('redis_cafile:', redis_cafile)
        elif opt in ("--redis_skipverify",):
            redis_skipverify = arg
            print('redis_skipverify:', redis_skipverify)
        elif opt in ("-m", "--mongo_ip"):
            mongo_ip = arg
            print('mongo_ip:', mongo_ip)
        elif opt in ("-P", "--mongo_port"):
            mongo_port = arg
            print('mongo_port:', mongo_port)
        elif opt in ("-X", "--mongo_user"):
            mongo_user = arg
            print('mongo_user:', mongo_user)
        elif opt in ("-S", "--mongo_pass"):
            mongo_pass = arg
            print('mongo_pass:', mongo_pass)
        elif opt in ("--mongo_certfile",):
            mongo_certfile = arg
            print('mongo_certfile:', mongo_certfile)
        elif opt in ("--mongo_keyfile",):
            mongo_keyfile = arg
            print('mongo_keyfile:', mongo_keyfile)
        elif opt in ("--mongo_cafile",):
            mongo_cafile = arg
            print('mongo_cafile:', mongo_cafile)
        elif opt in ("--mongo_skipverify",):
            mongo_skipverify = arg
            print('mongo_skipverify:', mongo_skipverify)
        elif opt in ("--rs_name",):
            rs_name = arg
            print('rs_name:', rs_name)
        elif opt in ("-u", "--blueking_cmdb_url"):
            cc_url = arg
            print('blueking_cmdb_url:', cc_url)
        elif opt in ("-U", "--blueking_paas_url"):
            paas_url = arg
            print('blueking_pass_url:', paas_url)
        elif opt in ("-l", "--listen_port"):
            server_ports["cmdb_webserver"] = arg
            print("listen_port:", server_ports["cmdb_webserver"])
        elif opt in ("--auth_address",):
            auth["auth_address"] = arg
            print("auth_address:", auth["auth_address"])
        elif opt in ("--auth_enabled",):
            auth["auth_enabled"] = arg
            print("auth_enabled:", auth["auth_enabled"])
        elif opt in ("--auth_scheme",):
            auth["auth_scheme"] = arg
            print("auth_scheme:", auth["auth_scheme"])
        elif opt in ("--auth_app_code",):
            auth["auth_app_code"] = arg
            print("auth_app_code:", auth["auth_app_code"])
        elif opt in ("--auth_app_secret",):
            auth["auth_app_secret"] = arg
            print("auth_app_secret:", auth["auth_app_secret"])
        elif opt in ("--auth_sync_workers",):
            auth["auth_sync_workers"] = arg
            print("auth_sync_workers:", auth["auth_sync_workers"])
        elif opt in ("--auth_sync_interval_minutes",):
            auth["auth_sync_interval_minutes"] = arg
            print("auth_sync_interval_minutes:", auth["auth_sync_interval_minutes"])
        elif opt in ("--full_text_search",):
            full_text_search = arg
            print('full_text_search:', full_text_search)
        elif opt in("-es","--es_url",):
            es_url = arg
            print('es_url:', es_url)
        elif opt in ("--es_user",):
            es_user = arg
            print('es_user:', es_user)
        elif opt in ("--es_pass",):
            es_pass = arg
            print('es_pass:', es_pass)
        elif opt in ("--es_shard_num",):
            es_shard_num = arg
            print('es_shard_num:', es_shard_num)
        elif opt in ("--es_replica_num",):
            es_replica_num = arg
            print('es_replica_num:', es_replica_num)
        elif opt in("-v","--log_level",):
            log_level = arg
            print('log_level:', log_level)
        elif opt in("--register_ip",):
            register_ip = arg
            print('register_ip:', register_ip)
        elif opt in("--user_info",):
            user_info = arg
            print('user_info:', user_info)
        elif opt in("--enable_cryptor",):
            enable_cryptor = arg
            print('enable_cryptor:', enable_cryptor)
        elif opt in("--secret_key_url",):
            secret_key_url = arg
            print('secret_key_url:', secret_key_url)
        elif opt in("--secrets_addrs",):
            secrets_addrs = arg
            print('secrets_addrs:', secrets_addrs)
        elif opt in("--secrets_token",):
            secrets_token = arg
            print('secrets_token:', secrets_token)
        elif opt in("--secrets_project",):
            secrets_project = arg
            print('secrets_project:', secrets_project)
        elif opt in("--secrets_env",):
            secrets_env = arg
            print('secrets_env:', secrets_env)

    if 0 == len(rd_server):
        print('please input the ZooKeeper address, eg:127.0.0.1:2181')
        sys.exit()
    if 0 == len(db_name):
        print('please input the database name, eg:cmdb')
        sys.exit()
    if 0 == len(redis_ip):
        print('please input the redis ip, eg: 127.0.0.1')
        sys.exit()
    if redis_port < 0:
        print('please input the redis port, eg:6379')
        sys.exit()
    if 0 == len(redis_pass):
        print('please input the redis password')
        sys.exit()
    if 0 == len(mongo_ip):
        print('please input the mongo ip, eg:127.0.0.1')
        sys.exit()
    if mongo_port < 0:
        print('please input the mongo port, eg:27017')
        sys.exit()
    if 0 == len(mongo_user):
        print('please input the mongo user, eg:cc')
        sys.exit()
    if 0 == len(mongo_pass):
        print('please input the mongo password')
        sys.exit()
    if 0 == len(cc_url):
        print('please input the blueking cmdb url')
        sys.exit()
    if 0 == len(paas_url):
        print('please input the blueking paas url')
        sys.exit()
    if not cc_url.startswith("http://"):
        print('blueking cmdb url not start with http://')
        sys.exit()

    if full_text_search not in ["off", "on"]:
        print('full_text_search can only be off or on')
        sys.exit()
    if full_text_search == "on":
        if not(es_url.startswith("http://") or es_url.startswith("https://")) :
            print('es url not start with http:// or https://')
            sys.exit()

    if enable_cryptor == "true":
        if len(secret_key_url) == 0 or len(secrets_addrs) == 0 or len(secrets_token) == 0 or len(secrets_project) == 0 or len(secrets_env) == 0:
            print('secret_key_url, secrets_addrs, secrets_token, secrets_project, secrets_env must be set when enable_cryptor is true')
            sys.exit()

    if auth["auth_scheme"] not in ["internal", "iam"]:
        print('auth_scheme can only be internal or iam')
        sys.exit()

    if auth["auth_enabled"] not in ["true", "false"]:
        print('auth_enabled value invalid, can only be `true` or `false`')
        sys.exit()

    if auth["auth_scheme"] == "iam" and auth["auth_enabled"] == 'true':
        if not auth["auth_address"]:
            print("auth_address can't be empty when iam auth enabled")
            sys.exit()

        if not auth["auth_app_code"]:
            print("auth_app_code can't be empty when iam auth enabled")
            sys.exit()

        if not auth["auth_app_secret"]:
            print("auth_app_secret can't be empty when iam auth enabled")
            sys.exit()

    availableLogLevel = [str(i) for i in range(0, 10)]
    if log_level not in availableLogLevel:
        print("available log_level value are: %s" %  availableLogLevel)
        sys.exit()

    generate_config_file(
        rd_server_v=rd_server,
        db_name_v=db_name,
        redis_ip_v=redis_ip,
        redis_port_v=redis_port,
        redis_pass_v=redis_pass,
        sentinel_pass_v=sentinel_pass,
        redis_certfile_v=redis_certfile,
        redis_keyfile_v=redis_keyfile,
        redis_cafile_v=redis_cafile,
        redis_skipverify_v=redis_skipverify,
        mongo_ip_v=mongo_ip,
        mongo_port_v=mongo_port,
        mongo_user_v=mongo_user,
        mongo_pass_v=mongo_pass,
        mongo_certfile_v=mongo_certfile,
        mongo_keyfile_v=mongo_keyfile,
        mongo_cafile_v=mongo_cafile,
        mongo_skipverify_v=mongo_skipverify,
        rs_name=rs_name,
        cc_url_v=cc_url,
        paas_url_v=paas_url,
        full_text_search=full_text_search,
        es_url_v=es_url,
        es_user_v=es_user,
        es_pass_v=es_pass,
        es_shard_num_v=es_shard_num,
        es_replica_num_v=es_replica_num,
        log_level=log_level,
        register_ip=register_ip,
        user_info=user_info,
        enable_cryptor_v=enable_cryptor,
        secret_key_url_v=secret_key_url,
        secrets_addrs_v=secrets_addrs,
        secrets_token_v = secrets_token,
        secrets_project_v = secrets_project,
        secrets_env_v = secrets_env,
        rd_cafile_v=rd_cafile,
        rd_certfile_v=rd_certfile,
        rd_keyfile_v=rd_keyfile,
        rd_skipverify_v=rd_skipverify,
        rd_certpassword_v=rd_certpassword,
        **auth
    )
    update_start_script(rd_server, server_ports, auth['auth_enabled'], log_level, register_ip, enable_cryptor, 
                       rd_cafile, rd_certfile, rd_keyfile, rd_skipverify, rd_certpassword)
    print('initial configurations success, configs could be found at cmdb_adminserver/configures')
    print('initial monstache config success, configs could be found at monstache/etc')


if __name__ == "__main__":
    main(sys.argv[1:])
