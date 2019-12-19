#!/usr/bin/python
# -*- coding: utf-8 -*-

import sys
import getopt
import os
import shutil
from string import Template


class FileTemplate(Template):
    delimiter = '$'


def generate_config_file(
        rd_server_v, db_name_v, redis_ip_v, redis_port_v,
        redis_pass_v, mongo_ip_v, mongo_port_v, mongo_user_v, mongo_pass_v, txn_enabled_v,
        cc_url_v, paas_url_v, full_text_search, es_url_v, auth_address, auth_app_code,
        auth_app_secret, auth_enabled, auth_scheme, auth_sync_workers, auth_sync_interval_minutes, log_level
):
    output = os.getcwd() + "/cmdb_adminserver/configures/"
    context = dict(
        db=db_name_v,
        mongo_user=mongo_user_v,
        mongo_host=mongo_ip_v,
        mongo_pass=mongo_pass_v,
        mongo_port=mongo_port_v,
        txn_enabled=txn_enabled_v,
        redis_host=redis_ip_v,
        redis_pass=redis_pass_v,
        redis_port=redis_port_v,
        cc_url=cc_url_v,
        paas_url=paas_url_v,
        es_url=es_url_v,
        ui_root="../web",
        agent_url=paas_url_v,
        configures_dir=output,
        rd_server=rd_server_v,
        auth_address=auth_address,
        auth_app_code=auth_app_code,
        auth_app_secret=auth_app_secret,
        auth_enabled=auth_enabled,
        auth_scheme=auth_scheme,
        auth_sync_workers=auth_sync_workers,
        auth_sync_interval_minutes=auth_sync_interval_minutes,
        full_text_search=full_text_search
    )
    if not os.path.exists(output):
        os.mkdir(output)

    # apiserver.conf
    apiserver_file_template_str = '''
[auth]
address = $auth_address
appCode = $auth_app_code
appSecret = $auth_app_secret
    '''

    template = FileTemplate(apiserver_file_template_str)
    result = template.substitute(**context)
    with open(output + "apiserver.conf", 'w') as tmp_file:
        tmp_file.write(result)

    # datacollection.conf
    datacollection_file_template_str = '''
[mongodb]
host = $mongo_host
port = $mongo_port
usr = $mongo_user
pwd = $mongo_pass
database = $db
maxOpenConns = 3000
maxIdleConns = 1000
mechanism = SCRAM-SHA-1
txnEnabled = $txn_enabled

[snap-redis]
host = $redis_host
port = $redis_port
pwd = $redis_pass
database = 0

[discover-redis]
host = $redis_host
port = $redis_port
pwd = $redis_pass
database = 0

[netcollect-redis]
host = $redis_host
port = $redis_port
pwd = $redis_pass
database = 0

[redis]
host = $redis_host
port = $redis_port
pwd = $redis_pass
database = 0

[auth]
address = $auth_address
appCode = $auth_app_code
appSecret = $auth_app_secret
'''

    template = FileTemplate(datacollection_file_template_str)
    result = template.substitute(**context)
    with open(output + "datacollection.conf", 'w') as tmp_file:
        tmp_file.write(result)

    # eventserver.conf
    eventserver_file_template_str = '''
[mongodb]
host = $mongo_host
port = $mongo_port
usr = $mongo_user
pwd = $mongo_pass
database = $db
maxOpenConns = 3000
maxIDleConns = 1000
mechanism = SCRAM-SHA-1
txnEnabled = $txn_enabled

[redis]
host = $redis_host
port = $redis_port
pwd = $redis_pass
database = 0
maxOpenConns = 3000
maxIDleConns = 1000
'''

    template = FileTemplate(eventserver_file_template_str)
    result = template.substitute(**context)
    with open(output + "eventserver.conf", 'w') as tmp_file:
        tmp_file.write(result)

    # host.conf
    host_file_template_str = '''
[gse]
addr = $rd_server
user = bkzk
pwd = L%blKas
[redis]
host = $redis_host
port = $redis_port
pwd = $redis_pass
database = 0
maxOpenConns = 3000
maxIDleConns = 1000
[auth]
address = $auth_address
appCode = $auth_app_code
appSecret = $auth_app_secret
enable = $auth_enabled
'''
    template = FileTemplate(host_file_template_str)
    result = template.substitute(**context)
    with open(output + "host.conf", 'w') as tmp_file:
        tmp_file.write(result)

    # migrate.conf
    migrate_file_template_str = '''
[config-server]
addrs = $rd_server
usr =
pwd =
[register-server]
addrs = $rd_server
usr =
pwd =
[mongodb]
host =$mongo_host
port = $mongo_port
usr = $mongo_user
pwd = $mongo_pass
database = $db
maxOpenConns = 3000
maxIDleConns = 1000
mechanism = SCRAM-SHA-1
txnEnabled = $txn_enabled

[confs]
dir = $configures_dir
[errors]
res = conf/errors
[language]
res = conf/language
[auth]
address = $auth_address
appCode = $auth_app_code
appSecret = $auth_app_secret
enableSync = false
syncWorkers = $auth_sync_workers
syncIntervalMinutes = $auth_sync_interval_minutes
    '''

    template = FileTemplate(migrate_file_template_str)
    result = template.substitute(**context)
    with open(output + "migrate.conf", 'w') as tmp_file:
        tmp_file.write(result)

    # coreservice.conf
    coreservice_file_template_str = '''
[mongodb]
host = $mongo_host
port = $mongo_port
usr = $mongo_user
pwd = $mongo_pass
database = $db
maxOpenConns = 3000
maxIDleConns = 1000
mechanism = SCRAM-SHA-1
txnEnabled = $txn_enabled

[redis]
host = $redis_host
port = $redis_port
pwd = $redis_pass
database = 0
maxOpenConns = 3000
maxIDleConns = 1000
'''

    template = FileTemplate(coreservice_file_template_str)
    result = template.substitute(**context)
    with open(output + "coreservice.conf", 'w') as tmp_file:
        tmp_file.write(result)

    # proc.conf
    proc_file_template_str = '''
[redis]
host = $redis_host
port = $redis_port
pwd = $redis_pass
database = 0

[auth]
address = $auth_address
appCode = $auth_app_code
appSecret = $auth_app_secret

[mongodb]
host = $mongo_host
port = $mongo_port
usr = $mongo_user
pwd = $mongo_pass
database = $db
maxOpenConns = 3000
maxIDleConns = 1000
txnEnabled = $txn_enabled
'''
    template = FileTemplate(proc_file_template_str)
    result = template.substitute(**context)
    with open(output + "proc.conf", 'w') as tmp_file:
        tmp_file.write(result)

    # operation.conf
    operation_file_template_str = '''
[mongodb]
host = $mongo_host
port = $mongo_port
usr = $mongo_user
pwd = $mongo_pass
database = $db
maxOpenConns = 3000
maxIDleConns = 1000
txnEnabled = $txn_enabled
'''
    template = FileTemplate(operation_file_template_str)
    result = template.substitute(**context)
    with open(output + "operation.conf", 'w') as tmp_file:
        tmp_file.write(result)

    # topo.conf
    topo_file_template_str = '''
[mongodb]
host = $mongo_host
port = $mongo_port
usr = $mongo_user
pwd = $mongo_pass
database = $db
maxOpenConns = 3000
maxIDleConns = 1000
mechanism = SCRAM-SHA-1
txnEnabled = $txn_enabled

[redis]
host = $redis_host
port = $redis_port
pwd = $redis_pass
database = 0
maxOpenConns = 3000
maxIDleConns = 1000

[level]
businessTopoMax = 7
[auth]
address = $auth_address
appCode = $auth_app_code
appSecret = $auth_app_secret

[es]
full_text_search = $full_text_search
url=$es_url
'''

    template = FileTemplate(topo_file_template_str)
    result = template.substitute(**context)
    with open(output + "topo.conf", 'w') as tmp_file:
        tmp_file.write(result)

    # webserver.conf
    webserver_file_template_str = '''
[api]
version = v3
[session]
name = cc3
skip = $skip
defaultlanguage = zh-cn
host = $redis_host
port = $redis_port
secret = $redis_pass
multiple_owner = 0
[site]
domain_url = ${cc_url}
bk_login_url = ${paas_url}/login/?app_id=%s&c_url=%s
app_code = cc
check_url = ${paas_url}/login/accounts/get_user/?bk_token=
bk_account_url = ${paas_url}/login/accounts/get_all_user/?bk_token=%s
resources_path = /tmp/
html_root = $ui_root
full_text_search = $full_text_search
[app]
agent_app_url = ${agent_url}/console/?app=bk_agent_setup
authscheme = $auth_scheme
'''
    template = FileTemplate(webserver_file_template_str)
    skip = '1'
    if auth_enabled == "true":
        skip = '0'
    result = template.substitute(skip=skip, **context)
    with open(output + "webserver.conf", 'w') as tmp_file:
        tmp_file.write(result)

    # task.conf
    taskserver_file_template_str = '''
[mongodb]
host = $mongo_host
port = $mongo_port
usr = $mongo_user
pwd = $mongo_pass
database = $db
maxOpenConns = 3000
maxIdleConns = 1000
mechanism = SCRAM-SHA-1
txnEnabled = $txn_enabled
maxIDleConns = 1000
mechanism = SCRAM-SHA-1
[redis]
host = $redis_host
port = $redis_port
pwd = $redis_pass
database = 0
maxOpenConns = 3000
maxIDleConns = 1000
'''
    template = FileTemplate(taskserver_file_template_str)
    result = template.substitute(**context)
    with open(output + "task.conf", 'w') as tmp_file:
        tmp_file.write(result)

def update_start_script(rd_server, server_ports, enable_auth, log_level):
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
                    filedata = filedata.replace('rd_server_placeholder', "configures/migrate.conf")
                    filedata = filedata.replace('regdiscv', "config")
                else:
                    filedata = filedata.replace('rd_server_placeholder', rd_server)

                extend_flag = ''
                if d in ['cmdb_apiserver', 'cmdb_hostserver', 'cmdb_datacollection', 'cmdb_procserver', 'cmdb_toposerver', 'cmdb_eventserver']:
                    extend_flag += ' --enable-auth=%s ' % enable_auth
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
    mongo_ip = ''
    mongo_port = 27017
    mongo_user = ''
    mongo_pass = ''
    txn_enabled = 'false'
    cc_url = ''
    paas_url = 'http://127.0.0.1'
    auth = {
        "auth_scheme": "internal",
        # iam options
        "auth_address": "",
        "auth_enabled": "false",
        "auth_app_code": "bk_cmdb",
        "auth_app_secret": "",
        "auth_sync_workers": "1",
        "auth_sync_interval_minutes": "45",
    }
    full_text_search = 'off'
    es_url = 'http://127.0.0.1:9200'
    log_level = '3'

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
        "cmdb_taskserver": 60012
    }
    arr = [
        "help", "discovery=", "database=", "redis_ip=", "redis_port=",
        "redis_pass=", "mongo_ip=", "mongo_port=",
        "mongo_user=", "mongo_pass=", "txn_enabled=", "blueking_cmdb_url=",
        "blueking_paas_url=", "listen_port=", "es_url=", "auth_address=",
        "auth_app_code=", "auth_app_secret=", "auth_enabled=",
        "auth_scheme=", "auth_sync_workers=", "auth_sync_interval_minutes=", "full_text_search=", "log_level="
    ]
    usage = '''
    usage:
      --discovery          <discovery>            the ZooKeeper server address, eg:127.0.0.1:2181
      --database           <database>             the database name, default cmdb
      --redis_ip           <redis_ip>             the redis ip, eg:127.0.0.1
      --redis_port         <redis_port>           the redis port, default:6379
      --redis_pass         <redis_pass>           the redis user password
      --mongo_ip           <mongo_ip>             the mongo ip ,eg:127.0.0.1
      --mongo_port         <mongo_port>           the mongo port, eg:27017
      --mongo_user         <mongo_user>           the mongo user name, default:cc
      --mongo_pass         <mongo_pass>           the mongo password
      --txn_enabled        <txn_enabled>          txn_enabled, true or false
      --blueking_cmdb_url  <blueking_cmdb_url>    the cmdb site url, eg: http://127.0.0.1:8088 or http://bk.tencent.com
      --blueking_paas_url  <blueking_paas_url>    the blueking paas url, eg: http://127.0.0.1:8088 or http://bk.tencent.com
      --listen_port        <listen_port>          the cmdb_webserver listen port, should be the port as same as -c <blueking_cmdb_url> specified, default:8083
      --auth_scheme        <auth_scheme>          auth scheme, ex: internal, iam
      --auth_enabled       <auth_enabled>         iam auth enabled, true or false
      --auth_address       <auth_address>         iam address
      --auth_app_code      <auth_app_code>        app code for iam, default bk_cmdb
      --auth_app_secret    <auth_app_secret>      app code for iam
      --full_text_search   <full_text_search>     full text search on or off
      --es_url             <es_url>               the es listen url, see in es dir config/elasticsearch.yml, (network.host, http.port), default: http://127.0.0.1:9200
      --log_level          <log_level>            log level to start cmdb process, default: 3


    demo:
    python init.py  \\
      --discovery          127.0.0.1:2181 \\
      --database           cmdb \\
      --redis_ip           127.0.0.1 \\
      --redis_port         6379 \\
      --redis_pass         1111 \\
      --mongo_ip           127.0.0.1 \\
      --mongo_port         27017 \\
      --mongo_user         cc \\
      --mongo_pass         cc \\
      --txn_enabled        false \\
      --blueking_cmdb_url  http://127.0.0.1:8080/ \\
      --blueking_paas_url  http://paas.domain.com \\
      --listen_port        8080 \\
      --auth_scheme        internal \\
      --auth_enabled       false \\
      --auth_address       https://iam.domain.com/ \\
      --auth_app_code      bk_cmdb \\
      --auth_app_secret    xxxxxxx \\
      --auth_sync_workers  1 \\
      --auth_sync_interval_minutes  45 \\
      --full_text_search   off \\
      --es_url             http://127.0.0.1:9200 \\
      --log_level          3
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
        elif opt in ("--txn_enabled",):
            txn_enabled = arg
            print('txn_enabled:', txn_enabled)
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
        elif opt in("-v","--log_level",):
            log_level = arg
            print('log_level:', log_level)

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

    if txn_enabled not in ["true", "false"]:
        print('txn_enabled value invalid, can only be `true` or `false`')
        sys.exit()

    if full_text_search not in ["off", "on"]:
        print('full_text_search can only be off or on')
        sys.exit()
    if full_text_search == "on":
        if not es_url.startswith("http://"):
            print('es url not start with http://')
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
        mongo_ip_v=mongo_ip,
        mongo_port_v=mongo_port,
        mongo_user_v=mongo_user,
        mongo_pass_v=mongo_pass,
        txn_enabled_v=txn_enabled,
        cc_url_v=cc_url,
        paas_url_v=paas_url,
        full_text_search=full_text_search,
        es_url_v=es_url,
        log_level=log_level,
        **auth
    )
    update_start_script(rd_server, server_ports, auth['auth_enabled'], log_level)
    print('initial configurations success, configs could be found at cmdb_adminserver/configures')


if __name__ == "__main__":
    main(sys.argv[1:])
