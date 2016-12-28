# 蓝鲸智云配置平台安装步骤


文档中假设用户的服务器为linux整个安装流程也适合其它系统的安装。


### 服务器配置


  - 推荐采用nginx+php-fpm 的运行模式
  - php版本不低于5.6.9,nginx版本不低于1.8.0
  - nginx编译参数，需编译进pcre
  - php编译参数扩展 ./configure --prefix= -enable-fpm，另还需要（mysql、curl、pcntl、mbregex、mhash、zip、mbstring、openssl）等扩展

### 包文件

开发者在git clone代码之后，里面不仅包含代码包文件，还包含初始化所需要的数据库sql文件



### 安装步骤
* 数据库服务器上创建数据库cmdb,导入sql文件bk-cmdb.sql文件
```sh
server {
        listen       80;
        server_name  cmdb.bk.com;
        root   /data/htdocs/cc_openSource;
        #charset koi8-r;

        #access_log  logs/host.access.log  main;
         
        location / {
            index  index.php index.html index.htm;
            if (!-e $request_filename) {
               rewrite ^(.*)$ /index.php?s=$1 last;
               break;}
        }


        #error_page  404              /404.html;
    
        # redirect server error pages to the static page /50x.html
        #
        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }
    
        # proxy the PHP scripts to Apache listening on 127.0.0.1:80
        # 
        #  location ~ \.php$ {
        #      proxy_pass   http://127.0.0.1;
        #  }

        # pass the PHP scripts to FastCGI server listening on 127.0.0.1:9000
        #
        location ~ \.php$ {
            fastcgi_connect_timeout 300;
            fastcgi_read_timeout 300;
            fastcgi_send_timeout 300;
            fastcgi_buffer_size 128k;
            fastcgi_buffers 32 32k;
            fastcgi_pass   127.0.0.1:9000;
            fastcgi_index  index.php;
            fastcgi_param  SCRIPT_FILENAME  $document_root$fastcgi_script_name;
            include        fastcgi_params;
        }

        # deny access to .htaccess files, if Apache's document root
        # concurs with nginx's one
        #
        location ~ /\.ht {
            deny  all;
        }
    }
    
```

* 代码配置

  - 根目录中index.php中配置环境类型为 (development/testing/production)之一。
  - db.php中配置数据访问
  - config.php中配置session访问目录，$config['sess_save_path'] = '/data/session',保证配置的session目录可读写
  - 请保证此路径可读写 application/resource/upload/importPrivateHostByExcel，确保文件上传功能的正确性
  - 根据配置的环境类型找到对应的常量文件,例如前面环境类型配置的为 development 则在/config/development/constants.php中定义
 ```sh
    define('BASE_URL', 'http://cmdb.bk.com');   //访问主域名,务必带上http://
    define('COMPANY_NAME', '公司名称');        //当前公司名
    define('COOKIE_DOMAIN', '.bk.com');         //cookie访问域
  ```
  * 创建数据库，导入sql文件执行数据初始化，切换到工程根目录下执行php index.php /cli/Init/initUserData
  * 启动nginx与php-fpm
  * 配置hosts，使用 admin/blueking账号即可登录访问
  
  
### 现有功能简介

* 用户管理
* 业务管理
* 拓扑（集群、模块）管理
*  资源池管理
* 主机管理
* 日志查询


### FAQ
*  问：蓝鲸配置平台是否可以独立使用？答：是，蓝鲸配置平台依赖系统非常少，可以独立的作为cmdb使用。
* 问：对于可扩展性方面蓝鲸配置平台有何考虑？答：蓝鲸配置平台提供全套可扩展的rest api。
* 问：蓝鲸配置平台在蓝鲸体系中的位置？答：蓝鲸配置平台作为蓝鲸其它系统配置与数据的来源，处于核心的位置。


### 已知的故障或错误列表
* 无


