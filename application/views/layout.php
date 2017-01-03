<?php
/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

$app = $this->session->userdata('app');
if(!$app && $active !== '/app/index'){
    $url = BASE_URL.'/app/index';
    echo "<script language='javascript' type='text/javascript'>";
    echo "window.location.href='$url'";
    echo "</script>";
}
?>
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title><?php echo $header; ?></title>
    <meta content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no" name="viewport">
    <link href="<?php echo STATIC_URL;?>/static/assets/jquery-ui-1.11.0.custom/jquery-ui.css?version=<?php echo $version;?>" rel="stylesheet">
    <link href="<?php echo STATIC_URL;?>/static/assets/bootstrap-3.3.4/css/bootstrap.min.css?version=<?php echo $version;?>" rel="stylesheet" type="text/css" />
    <link href="<?php echo STATIC_URL;?>/static/assets/font-awesome/css/font-awesome.min.css?version=<?php echo $version;?>" rel="stylesheet" type="text/css" />
    <link href="<?php echo STATIC_URL;?>/static/css/AdminLTE.min.css?>" rel="stylesheet" type="text/css" />
    <link href="<?php echo STATIC_URL;?>/static/css/skin-blue.min.css?version=<?php echo $version;?>" rel="stylesheet" type="text/css" />
    <link href="<?php echo STATIC_URL;?>/static/assets/mCustomScrollbar-3.0.9/jquery.mCustomScrollbar.css?version=<?php echo $version;?>" rel="stylesheet">
    <link href="<?php echo STATIC_URL;?>/static/bk/css/bk_base.css?version=<?php echo $version;?>" rel="stylesheet" type="text/css">
    <link href="<?php echo STATIC_URL;?>/static/assets/select2-3.5.2/select2.css?version=<?php echo $version;?>" rel="stylesheet" />
    <link rel="stylesheet" type="text/css" href="<?php echo STATIC_URL;?>/static/assets/selectmenu/jquery-ui.css?version=<?php echo $version;?>" />
    <link rel="stylesheet" type="text/css" href="<?php echo STATIC_URL;?>/static/assets/icheck-1.x/skins/all.css?version=<?php echo $version;?>" />
    <link href="<?php echo STATIC_URL;?>/static/css/index.css?version=<?php echo $version;?>" rel="stylesheet" type="text/css" />
    <link href="<?php echo STATIC_URL;?>/static/assets/plugin4/bkDialog-1.0/css/ui-dialog.css?version=<?php echo $version;?>" rel="stylesheet" type="text/css" />
    <script src="<?php echo STATIC_URL;?>/static/assets/js/jquery-1.10.2.min.js?version=<?php echo $version;?>"></script>
    <script src="<?php echo STATIC_URL;?>/static/assets/bootstrap-3.3.4/js/bootstrap.min.js?version=<?php echo $version;?>"></script>
    <script src="<?php echo STATIC_URL;?>/static/assets/plugin4/bkDialog-1.0/js/dialog.js?version=<?php echo $version;?>"></script>

    <script type='text/javascript'>
        var environment = '<?php echo $environment ?>';
    </script>
    <script src="<?php echo STATIC_URL;?>/static/js/index.js?version=<?php echo $version;?>"></script>
</head>

<?php if(isset($_COOKIE['isCollapse']) && $_COOKIE['isCollapse']==true){ ?>
<body class="skin-blue sidebar-mini sidebar-collapse" id="index">
<?php }else{ ?>
<body class="skin-blue sidebar-mini" id="index">
<?php }?>
<div class="wrapper">
    <!-- 头部 -->
    <header class="main-header">
        <!-- 头部左侧Logo -->
        <a href="/" class="logo">
            <img src="<?php echo STATIC_URL;?>/static/img/logo.png" alt="" height="35">
        </a>
        <!-- 头部右侧Navbar -->
        <nav class="navbar navbar-static-top" role="navigation">
            <a href="#" title="收起" class="c-sidebar-toggle" data-toggle="offcanvas" role="button">
                <span class="sr-only">Toggle navigation</span>
            </a>
            <div class="navbar-custom-menu pull-left">
                 <form target='_self' action="/host/hostQuery" method='post' id="quick_search">
                    <div class="form-group has-feedback" id="speed_search_input">
                        <!-- <input class="form-control mt10" type="text" value="" placeholder="快速查询"> -->
                        <textarea class="search-textarea" placeholder="快速查询" id="search-textarea"></textarea>
                        <input type="hidden" id="IfInnerIPexact" value="">
                        <input type="hidden" id="IfOuterexact" value="">
                        <span class="glyphicon glyphicon-search form-control-feedback search-btn" id ="speed_search" aria-hidden="true"></span>
                    </div>
                </form>
            </div>

            <!-- 头部右侧 Menu -->
            <div class="navbar-custom-menu">
                <ul class="nav navbar-nav">
                    <!-- 帮助 -->
                    <li>
                        <a href="http://www.bkclouds.cc/wiki/course/" target="_blank" class="c-nav-help">
                            <i class="fa fa-question-circle mr5 f20"></i>
                            <span>帮助</span>
                        </a>
                    </li>
                    <!-- 业务 -->
                    <?php $defaultApp = $this->session->userdata('defaultApp');?>
                    <?php $app = $this->session->userdata('app');?>
                    <?php $company_list = $this->session->userdata('company_list');?>
                    <li id="chooseBusiness" class="dropdown">
                            <a href="javascript:void(0);" class="dropdown-toggle text-center" data-toggle="dropdown">
                                <span class="first-letter mr5" title="<?php echo $defaultApp['ApplicationName']; ?>"><?php echo mb_substr($defaultApp['ApplicationName'],0,1,'UTF-8');?></span>
                                <span class="nav-group-name"><?php echo $defaultApp['ApplicationName']; ?><i class='caret'></i></span>
                            </a>
                            <ul class="dropdown-menu business-menu">
                                <?php foreach ($app as $item):?>
                                    <li <?php if($defaultApp['ApplicationID']==$item['ApplicationID']){echo 'class="active"';}?>><a class="defaultApp" value="<?php echo $item['ApplicationID'];?>"><?php echo $item['ApplicationName'];?><span class="text-warning">(<?php echo $item['ApplicationHostCount'];?>)</span></a></li>
                                <?php endforeach;?>
                            </ul>
                    </li>
                    <li>
                        <a href="/account/logout" class="c-nav-out">
                            <i class="fa fa-sign-out mr5 f20"></i>
                            <span>注销</span>
                        </a>
                    </li>
                </ul>
            </div>
        </nav>
    </header>
    </div>
    <!-- 左侧菜单 -->
    <aside class="main-sidebar">
        <section class="sidebar">
            <!-- 菜单 -->
            <ul class="sidebar-menu">
                <!-- Optionally, you can add icons to the links -->
                <li class="<?php echo isset($active) && $active=='/welcome/index' ? 'active' : '';?>"><a href="/welcome/index"><i class="fa fa-dashboard"></i> <span>总览</span></a></li>
                <li id="hostmng" class="<?php echo isset($active) && $active=='/host/hostQuery' ? 'active' : '';?>"><a href="/host/hostQuery"><i class="fa fa-desktop"></i> <span>主机管理</span></a></li>
                <li id="apptopo" class="<?php echo isset($active) && $active=='/topology' ? 'active' : '';?>">
                    <a href="javascript:void(0)"><i class="fa fa-database"></i> <span>拓扑配置</span> <i class="fa fa-angle-left pull-right"></i></a>
                    <ul class="treeview-menu">
                        <li class="<?php echo isset($subactive) && $subactive=='/index' ? 'active' : '';?>" ><a href="/topology/index"><i class="fa fa-sitemap"></i>树状视图</a></li>
                        <li class="<?php echo isset($subactive) && $subactive=='/topolist' ? 'active' : '';?>" ><a href="/topology/topolist"><i class="fa fa-list-ul"></i>列表视图</a></li>
                    </ul>
                </li>
                <li class="<?php echo isset($active) && $active=='/app/index' ? 'active' : '';?>"><a href="/app/index"><i class="fa fa-puzzle-piece"></i> <span>业务管理</span></a></li>
                <li id="hostdis" class="<?php echo isset($active) && ($active=='/host/quickImport' || $active=='/host/quickBuy')? 'active' : '';?>"><a href="/host/quickImport"><i class="fa fa-graduation-cap"></i> <span>资源池管理</span></a></li>
                <li class="<?php echo isset($active) && $active=='/operationlog/index' ? 'active' : '';?>"><a href="/operationLog/index"><i class="fa fa-calendar"></i> <span>操作日志</span></a></li>
                <li id="userlist" class="<?php echo isset($active) && ($active=='/account/index' || $active=='/host/quickBuy')? 'active' : '';?>"><a href="/account/index"><i class="fa fa-user"></i> <span>用户管理</span></a></li>
            </ul>
        </section>
    </aside>
    <!-- 左侧菜单 结束 -->
    <?php echo $content;?>
    <script src="<?php echo STATIC_URL;?>/static/js/common.js?version=<?php echo $version;?>"></script>
    <!-- 项目需要引用的js文件 -->
    <script type='text/javascript'>
        $(document).ready(function(){
            <?php if(!$app):?>
            $('.sidebar').find('li').divLoad('show');
            <?php endif;?>

            <?php if(count($app) > 1):?>
            $('.defaultApp').click(function(){
                var appId=$(this).attr('value');
                var defaultappId = cookie.get('defaultAppId');
                if(appId == defaultappId){
                    return;
                }
                $.ajax({
                    url: "/welcome/setDefaultApp?ApplicationID="+appId,
                    type: "POST",
                    dataType: "json",
                    success: function (response) {
                        if(response.success){
                            setTimeout(function () {
                                window.location.reload();
                            }, 500);
                        }

                        return true;
                    }
                });
            });
            <?php endif;?>
        });
    </script>
</body>
</html>
