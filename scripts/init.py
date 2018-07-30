#!/usr/bin/python 
# -*- coding: utf-8 -*-   

import sys,getopt,os,shutil
from string import Template


class FileTemplate(Template):
    delimiter='$'

def generate_config_file(rd_server_v,db_name_v,redis_ip_v,redis_port_v,redis_user_v,redis_pass_v,mongo_ip_v,mongo_port_v,mongo_user_v,mongo_pass_v,cc_url_v,paas_url_v):
    
    output = os.getcwd()+"/cmdb_adminserver/configures/"
 
    if not os.path.exists(output):
        os.mkdir(output)

    # apiserver.conf
    apiserver_file_template_str ='''
    '''

    template = FileTemplate(apiserver_file_template_str)
    result = template.substitute()
    with open( output + "apiserver.conf",'w') as tmp_file:
        tmp_file.write(result)

    # auditcontroller.conf
    auditcontroller_file_template_str='''
    [mongodb]
    host = $mongo_host
    usr = $mongo_user
    pwd = $mongo_pass
    database = $db
    port = $mongo_port
    maxOpenConns = 3000
    maxIdleConns = 1000
    '''
    template = FileTemplate(auditcontroller_file_template_str)
    result = template.substitute(dict(db=db_name_v,mongo_user=mongo_user_v,mongo_host=mongo_ip_v,mongo_pass=mongo_pass_v,mongo_port=mongo_port_v))
    with open( output + "auditcontroller.conf",'w') as tmp_file:
        tmp_file.write(result)

    # datacollection.conf
    datacollection_file_template_str='''
    [mongodb]
    host = $mongo_host
    usr = $mongo_user
    pwd = $mongo_pass
    database = $db
    port = $mongo_port
    maxOpenConns = 3000
    maxIdleConns = 1000

    [snap-redis]
    host = $redis_host
    usr = $redis_user
    pwd = $redis_pass
    database = 0

    [discover-redis]
    host = $redis_host
    usr = $redis_user
    pwd = $redis_pass
    database = 0

    [discover-redis]
    host = $redis_host
    usr = $redis_user
    pwd = $redis_pass
    database = 0
    chan = 3_snapshot

    [redis]
    host = $redis_host
    usr = $redis_user
    pwd = $redis_pass
    database = 0
    '''

    template = FileTemplate(datacollection_file_template_str)
    result = template.substitute(dict(db=db_name_v,redis_host=redis_ip_v+":"+str(redis_port_v),redis_user=redis_user_v,redis_pass=redis_pass_v, mongo_user=mongo_user_v,mongo_host=mongo_ip_v,mongo_pass=mongo_pass_v,mongo_port=mongo_port_v))
    with open( output + "datacollection.conf",'w') as tmp_file:
        tmp_file.write(result)

    # eventserver.conf
    eventserver_file_template_str='''
    [mongodb]
    host=$mongo_host
    usr=$mongo_user
    pwd=$mongo_pass
    database=$db
    port=$mongo_port
    maxOpenConns=3000
    maxIDleConns=1000
    [redis]
    host=$redis_host
    usr=$redis_user
    pwd=$redis_pass
    database=0
    port=$redis_port
    maxOpenConns=3000
    maxIDleConns=1000
    '''
    
    template = FileTemplate(eventserver_file_template_str)
    result = template.substitute(dict(db=db_name_v,redis_host=redis_ip_v,redis_port=redis_port_v,redis_user=redis_user_v,redis_pass=redis_pass_v, mongo_user=mongo_user_v,mongo_host=mongo_ip_v,mongo_pass=mongo_pass_v,mongo_port=mongo_port_v))
    with open( output + "eventserver.conf",'w') as tmp_file:
        tmp_file.write(result)

    # host.conf
    host_file_template_str='''
    [gse]
    addr=$rd_server
    user=bkzk
    pwd=L%blKas
    '''
    template = FileTemplate(host_file_template_str)
    result = template.substitute(dict(rd_server=rd_server_v))
    with open( output + "host.conf",'w') as tmp_file:
        tmp_file.write(result)

    # hostcontroller.conf
    hostcontroller_file_template_str='''
    [mongodb]
    host=$mongo_host
    usr=$mongo_user
    pwd=$mongo_pass
    database=$db
    port=$mongo_port
    maxOpenConns=3000
    maxIDleConns=1000
    [redis]
    host=$redis_host
    usr=$redis_user
    pwd=$redis_pass
    database=0
    port=$redis_port
    maxOpenConns=3000
    maxIDleConns=1000
    '''

    template = FileTemplate(hostcontroller_file_template_str)
    result = template.substitute(dict(db=db_name_v,redis_host=redis_ip_v,redis_port=redis_port_v,redis_user=redis_user_v,redis_pass=redis_pass_v, mongo_user=mongo_user_v,mongo_host=mongo_ip_v,mongo_pass=mongo_pass_v,mongo_port=mongo_port_v))
    with open( output + "hostcontroller.conf",'w') as tmp_file:
        tmp_file.write(result)

    # migrate.conf
    migrate_file_template_str='''
    [config-server]
    addrs=$rd_server
    usr=
    pwd=
    [register-server]
    addrs=$rd_server
    usr=
    pwd=

    [mongodb]
    host =$mongo_host
    usr = $mongo_user
    pwd = $mongo_pass
    database = $db
    port = $mongo_port
    maxOpenConns = 3000
    maxIDleConns = 1000

    [confs]
    dir = $configures_dir
    [errors]
    res=conf/errors
    [language]
    res=conf/language
    '''

    template = FileTemplate(migrate_file_template_str)
    result = template.substitute(dict(db=db_name_v,configures_dir=output,rd_server=rd_server_v,redis_host=redis_ip_v,redis_port=redis_port_v,redis_user=redis_user_v,redis_pass=redis_pass_v, mongo_user=mongo_user_v,mongo_host=mongo_ip_v,mongo_pass=mongo_pass_v,mongo_port=mongo_port_v))
    with open( output + "migrate.conf",'w') as tmp_file:
        tmp_file.write(result)

    # objectcontroller.conf
    objectcontroller_file_template_str='''
    [mongodb]
    host=$mongo_host
    usr=$mongo_user
    pwd=$mongo_pass
    database=$db
    port=$mongo_port
    maxOpenConns=3000
    maxIDleConns=1000
    [redis]
    host=$redis_host
    usr=$redis_user
    pwd=$redis_pass
    database=0
    port=$redis_port
    maxOpenConns=3000
    maxIDleConns=1000
    '''

    template = FileTemplate(objectcontroller_file_template_str)
    result = template.substitute(dict(db=db_name_v,redis_host=redis_ip_v,redis_port=redis_port_v,redis_user=redis_user_v,redis_pass=redis_pass_v, mongo_user=mongo_user_v,mongo_host=mongo_ip_v,mongo_pass=mongo_pass_v,mongo_port=mongo_port_v))
    with open( output + "objectcontroller.conf",'w') as tmp_file:
        tmp_file.write(result)

    # proc.conf
    proc_file_template_str='''
    '''
    template = FileTemplate(proc_file_template_str)
    result = template.substitute()
    with open( output + "proc.conf",'w') as tmp_file:
        tmp_file.write(result)

    # proccontroller.conf
    proccontroller_file_template_str='''
    [mongodb]
    host=$mongo_host
    usr=$mongo_user
    pwd=$mongo_pass
    database=$db
    port=$mongo_port
    maxOpenConns=3000
    maxIDleConns=1000
    [redis]
    host=$redis_host
    usr=$redis_user
    pwd=$redis_pass
    database=0
    port=$redis_port
    maxOpenConns=3000
    maxIDleConns=1000
    '''

    template = FileTemplate(proccontroller_file_template_str)
    result = template.substitute(dict(db=db_name_v,redis_host=redis_ip_v,redis_port=redis_port_v,redis_user=redis_user_v,redis_pass=redis_pass_v, mongo_user=mongo_user_v,mongo_host=mongo_ip_v,mongo_pass=mongo_pass_v,mongo_port=mongo_port_v))
    with open( output + "proccontroller.conf",'w') as tmp_file:
        tmp_file.write(result)

    # topo.conf
    topo_file_template_str='''
    [mongodb]
    host=$mongo_host
    usr=$mongo_user
    pwd=$mongo_pass
    database=$db
    port=$mongo_port
    maxOpenConns=3000
    maxIDleConns=1000
    '''

    template = FileTemplate(topo_file_template_str)
    result = template.substitute(dict(db=db_name_v,mongo_user=mongo_user_v,mongo_host=mongo_ip_v,mongo_pass=mongo_pass_v,mongo_port=mongo_port_v))
    with open( output + "topo.conf",'w') as tmp_file:
        tmp_file.write(result)

    # webserver.conf
    webserver_file_template_str='''
    [api]
    version=v3
    [session]
    name=cc3
    skip=1
    defaultlanguage=zh-cn
    host=$redis_host
    port=$redis_port
    secret=$redis_pass
    multiple_owner=0
    [site]
    domain_url=${cc_url}
    bk_login_url=${paas_url}/login/?app_id=%s&c_url=%s
    app_code=cc
    check_url=${paas_url}/login/accounts/get_user/?bk_token=
    bk_account_url=${paas_url}/login/accounts/get_all_user/?bk_token=%s
    resources_path=/tmp/
    html_root=$ui_root
    [app]
    agent_app_url=${agent_url}/console/?app=bk_agent_setup
    '''
    ui_root_v = os.getcwd()+"/web"
    template = FileTemplate(webserver_file_template_str)
    result = template.substitute(dict(redis_host=redis_ip_v,redis_port=redis_port_v,redis_pass=redis_pass_v,cc_url=cc_url_v,paas_url=paas_url_v,ui_root=ui_root_v,agent_url=paas_url_v))
    with open( output + "webserver.conf",'w') as tmp_file:
        tmp_file.write(result)

def update_start_script(rd_server, server_ports):
    list_dirs = os.walk(os.getcwd()+"/")
    for root, dirs, _ in list_dirs:
        for d in dirs:
            if d.startswith("cmdb_"):
                if d == "cmdb_adminserver":
                    if os.path.exists(root+d+"/init_db.sh"):
                        shutil.copy(root + d + "/init_db.sh", os.getcwd()+"/init_db.sh")
                target_file = root+d +"/start.sh"
                if os.path.exists(target_file) and os.path.exists(root+d +"/template.sh.start"):
                    # Read in the file
                    with open(root+d +"/template.sh.start", 'r') as template_file :
                        filedata = template_file.read()
                        # Replace the target string
                        filedata = filedata.replace('cmdb-name-placeholder', d)
                        filedata = filedata.replace('cmdb-port-placeholder', str(server_ports[d]))
                        if d != "cmdb_adminserver":
                            filedata = filedata.replace('rd_server_placeholer', rd_server)
                        else:
                            filedata = filedata.replace('rd_server_placeholer', "configures/migrate.conf")
                            filedata = filedata.replace('regdiscv', "config")
                        # Write the file out again
                        with open(target_file, 'w') as new_file:
                            new_file.write(filedata)
    

def main(argv):
    db_name='cmdb'
    rd_server=''
    redis_ip=''
    redis_port=6379
    redis_user='cc'
    redis_pass=''
    mongo_ip=''
    mongo_port=27017
    mongo_user=''
    mongo_pass=''
    cc_url=''
    paas_url='http://127.0.0.1'

    server_ports={"cmdb_adminserver":60004,"cmdb_apiserver":8080,\
    "cmdb_auditcontroller":50005,"cmdb_datacollection":60005,\
    "cmdb_eventserver":60009,"cmdb_hostcontroller":50002,\
    "cmdb_hostserver":60001,"cmdb_objectcontroller":50001,\
    "cmdb_proccontroller":50003,"cmdb_procserver":60003,\
    "cmdb_toposerver":60002,"cmdb_webserver":8083}
    try:
        opts, _ = getopt.getopt(argv,"hd:D:r:p:x:s:m:P:X:S:u:U:a:l:"\
        ,["help","discovery=","database=","redis_ip=","redis_port="\
        ,"redis_user=","redis_pass=","mongo_ip=","mongo_port="\
        ,"mongo_user=","mongo_pass=","blueking_cmdb_url=","blueking_paas_url=","listen_port="])

    except getopt.GetoptError as e:
        print "\n \t",e.msg
        print "\n\tusage:" \
        ,"\n\t-discovery           <discovery>            the ZooKeeper server address, eg:127.0.0.1:2181"\
        ,"\n\t--database           <database>             the database name, default cmdb"\
        ,"\n\t--redis_ip           <redis_ip>             the redis ip, eg:127.0.0.1"\
        ,"\n\t--redis_port         <redis_port>           the redis port, default:6379"\
        ,"\n\t--redis_pass         <redis_pass>           the redis user password"\
        ,"\n\t--mongo_ip           <mongo_ip>             the mongo ip ,eg:127.0.0.1"\
        ,"\n\t--mongo_port         <mongo_port>           the mongo port, eg:27017"\
        ,"\n\t--mongo_user         <mongo_user>           the mongo user name, default:cc"\
        ,"\n\t--mongo_pass         <mongo_pass>           the mongo password"\
        ,"\n\t--blueking_cmdb_url  <blueking_cmdb_url>    the cmdb site url, eg: http://127.0.0.1:8088 or http://bk.tencent.com"\
        ,"\n\t--blueking_paas_url  <blueking_paas_url>    the blueking pass url, eg: http://127.0.0.1:8088 or http://bk.tencent.com"\
        ,"\n\t--listen_port        <listen_port>          the cmdb_webserver listen port, should be the port as same as -c <blueking_cmdb_url> specified, default:8083"
       
        sys.exit(2)
    if len(opts) == 0:
        print "\n\tusage:" \
        ,"\n\t--discovery          <discovery>            the ZooKeeper server address, eg:127.0.0.1:2181"\
        ,"\n\t--database           <database>             the database name, default cmdb"\
        ,"\n\t--redis_ip           <redis_ip>             the redis ip, eg:127.0.0.1"\
        ,"\n\t--redis_port         <redis_port>           the redis port, default:6379"\
        ,"\n\t--redis_pass         <redis_pass>           the redis user password"\
        ,"\n\t--mongo_ip           <mongo_ip>             the mongo ip ,eg:127.0.0.1"\
        ,"\n\t--mongo_port         <mongo_port>           the mongo port, eg:27017"\
        ,"\n\t--mongo_user         <mongo_user>           the mongo user name, default:cc"\
        ,"\n\t--mongo_pass         <mongo_pass>           the mongo password"\
        ,"\n\t--blueking_cmdb_url  <blueking_cmdb_url>    the cmdb site url, eg: http://127.0.0.1:8088 or http://bk.tencent.com"\
        ,"\n\t--blueking_paas_url  <blueking_paas_url>    the blueking paas url, eg: http://127.0.0.1:8088 or http://bk.tencent.com"\
        ,"\n\t--listen_port        <listen_port>          the cmdb_webserver listen port, should be the port as same as -c <blueking_cmdb_url> specified, default:8083"
       
        sys.exit(2)
    #print opts
    for opt, arg in opts:
        if opt in ('-h','--help'):
            print 'init.py --discovery <discovery> --database <database>  --redis_ip <redis_ip> --redis_port <redis_port> --redis_pass <redis_pass> --mongo_ip <mongo_ip> --mongo_port <mongo_port> --mongo_user <mongo_user> --mongo_pass <mongo_pass> --blueking_cmdb_url <blueking_cmdb_url> --blueking_paas_url <blueking_paas_url> --listen_port <listen_port>'
            sys.exit()
        elif opt in ("-d", "--discovery"):
            rd_server=arg
            print 'rd_server:'+rd_server
        elif opt in ("-D", "--database"):
            db_name=arg
            print 'database:',db_name
        elif opt in ("-r", "--redis_ip"):
            redis_ip= arg
            print 'redis_ip:'+redis_ip
        elif opt in ("-p","--redis_port"):
            redis_port = arg
            print 'redis_port:',redis_port
        elif opt in ("-s","--redis_pass"):
            redis_pass = arg
            print 'redis_pass:',redis_pass
        elif opt in ("-m","--mongo_ip"):
            mongo_ip=arg
            print 'mongo_ip:',mongo_ip
        elif opt in("-P","--mongo_port"):
            mongo_port=arg
            print 'mongo_port:',mongo_port
        elif opt in ("-X","--mongo_user"):
            mongo_user = arg
            print 'mongo_user:',mongo_user
        elif opt in ("-S","--mongo_pass"):
            mongo_pass = arg
            print 'mongo_pass:',mongo_pass
        elif opt in("-u","--blueking_cmdb_url"):
            cc_url=arg
            print 'blueking_cmdb_url:',cc_url
        elif opt in("-U","--blueking_paas_url"):
            paas_url = arg
            print 'blueking_pass_url:',paas_url
        elif opt in("-l","--listen_port"):
            server_ports["cmdb_webserver"]=arg
            print "listen_port:",server_ports["cmdb_webserver"]
    
    if 0 == len(rd_server):
        print 'please input the ZooKeeper address, eg:127.0.0.1:2181'
        sys.exit()
    if 0 == len(db_name):
        print 'please input the database name, eg:cmdb'
        sys.exit()
    if 0 == len(redis_ip):
        print 'please input the redis ip, eg: 127.0.0.1'
        sys.exit()
    if redis_port < 0:
        print 'please input the redis port, eg:6379'
        sys.exit()
    if 0 == len(redis_pass):
        print 'please input the redis password'
        sys.exit()
    if 0 == len(mongo_ip):
        print 'please input the mongo ip, eg:127.0.0.1'
        sys.exit()
    if mongo_port < 0:
        print 'please input the mongo port, eg:27017'
        sys.exit()
    if 0 == len(mongo_user):
        print 'please input the mongo user, eg:cc'
        sys.exit()
    if 0 == len(mongo_pass):
        print  'please input the mongo password'
        sys.exit()
    if 0 == len(cc_url):
        print 'please input the blueking cmdb url'
        sys.exit()
    if 0 == len(paas_url):
        print 'please input the blueking paas url'
        sys.exit()
    if not cc_url.startswith("http://"):
        print 'blueking cmdb url not start with http://'
        sys.exit()

    generate_config_file(rd_server,db_name,redis_ip,redis_port,redis_user,redis_pass,mongo_ip,mongo_port,mongo_user,mongo_pass,cc_url,paas_url)
    update_start_script(rd_server, server_ports)
    print 'initial configurations success, configs could be found at cmdb_adminserver/configures'
if __name__=="__main__":
    main(sys.argv[1:])
