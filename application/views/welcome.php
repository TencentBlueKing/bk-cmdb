<?php
/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

if($appCount):?>
     <!-- 主面板Content -->
        <div class="content-wrapper">
            <!-- 主面板 Header  -->
            <section class="content-header">
                <h1>我的服务资源</h1>
            </section>
            <!-- 主面板 Main-->
            <section class="content">
                <div class="row">
                    <!-- 业务数量 -->
                    <div class="col-md-4 col-xs-12">
                        <div class="small-box c-bg1">
                            <div class="div-img fl"><span class="today-icon icon1"></span></div>
                            <div class="div-text fl">
                                <?php if($appCount == 1):?>
                                <p>集群数</p>
                                <h3><?php echo $GroupCount;?></h3>
                                <?php else :?>
                                <p>业务数量</p>
                                <h3><?php echo $appCount; ?></h3>
                                <?php endif;?>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-4 col-xs-12">
                        <div class="small-box c-bg2">
                            <div class="div-img fl"><span class="today-icon icon2"></span></div>
                            <div class="div-text fl">
                                <p>设备数</p>
                                <h3><?php echo $hostCount; ?></h3>
                            </div>
                        </div>
                    </div>
                    <div class="col-md-4 col-xs-12">
                        <div class="small-box c-bg3">
                            <div class="div-img fl"><span class="today-icon icon4"></span></div>
                            <div class="div-text fl">
                                <p>空闲设备</p>
                                <h3><?php echo $emptyPoolHostCount; ?><span></span></h3>
                            </div>
                        </div>
                    </div>
                    <!-- /当前余额 -->
                </div>
                <div class="row">
                    <div class="col-lg-8 col-md-12 col-sm-12">
                    </div>
                    <div class="col-lg-4 col-md-12 col-sm-12">
                    </div>
                    <div class="col-md-6 col-sm-12">
                        <!-- 主机分布图 -->
                        <div class="c-host-wapper">
                            <div class="page-header">
                                <h4>主机分布图</h4>
                            </div>
                            <div class="c-host-area" id="host-area">
                            </div>
                        </div>
                        <!-- /主机分布图 -->
                    </div>
                    <div class="col-md-6 col-sm-12">
                        <!-- 操作日志 -->
                        <ul class="c-host-opt">
                            <div class="page-header">
                                <h4>操作日志 <small><a href="/operationLog/">查看全部</a></small></h4>
                            </div>
                            <div id="opt_box" class="opt-box">
                            <ul class="king-timeline-simple">
                                <?php foreach ($operationLog as $item):?>
                                <li class="info">
                                <div class="timeline-simple-wrap">
                                    <p class="info-title"><?php echo $item['Operator']. ' ' . $item['OpName'];?></p>
                                    <span class="timeline-simple-date"><?php echo Utility::tranTime($item['OpTime']);?></span>
                                            <p><?php echo $item['OpContent'];?></p>
                                </div>
                                </li>
                                <?php endforeach;?>
                            </ul>
                        </div>
                        <!-- /操作记录 -->
                    </div>
                </div>
            </section>
        </div>
        <div class="control-sidebar-bg"></div>
    </div>
        <!-- 项目js文件 -->
        <script src="<?php echo STATIC_URL;?>/static/assets/map/raphael-min.js"></script>
        <script src="<?php echo STATIC_URL;?>/static/assets/map/map.js"></script>
        <script src="<?php echo STATIC_URL;?>/static/assets/mCustomScrollbar-3.0.9/jquery.mCustomScrollbar.concat.min.js"></script>
        <script src="<?php echo STATIC_URL;?>/static/js/app.min.js"></script>
        <script>
            CC.index.deadLineTime=<?php echo json_encode(array());?>;
            CC.index.hostTrend=<?php echo json_encode(array());?>;
            CC.index.hostMap=<?php echo $hostIDC;?>;
            $(function(){
                 CC.index.init();
            })
        </script>
        <?php else :?>
    <div class="content-wrapper">
            <div class="no-host-content">
                <img src="<?php echo STATIC_URL;?>/static/img/expre_403.png">
                <h4 class="pt15">对不起，您当前没有可操作的业务，您可尝试如下操作</h4>
                <ul class="pt15" style="width:300px;">
                    <li class="text-left">点此 <a href="/App/newapp" id="home_creat_one">新建业务</a></li>
                    <li>联系您公司已有权限的同事为您开通权限</li>
                </ul>
            </div>
        </div>
        <div class="control-sidebar-bg"></div>
    </div>
    <?php endif;?>
