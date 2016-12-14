<?php
/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */
?>
<!-- 主面板Content -->
<link href="<?php echo STATIC_URL;?>/static/assets/bootstrap-switch-master/dist/css/bootstrap3/bootstrap-switch.css" rel="stylesheet">
<link href="<?php echo STATIC_URL;?>/static/assets/select2-3.5.2/select2.css" rel="stylesheet">
<style>
    #confTreeContainer>.k-group>li>.k-top .k-in{
        background: none !important;
        cursor: default !important;
    }
</style>
<div class="content-wrapper">
    <!-- 主面板 Main-->
    <section class="content">
        <div class="host-sidebar-left">
            <!-- 分布拓扑 -->
            <div class="c-host-side" id="c_host_side">
                <h4 class="mb0">拓扑模型<i class="fa fa-close pull-right mr15 conf-delete-link" style="display:none;"></i></h4>
                <form class="form-inline pl15 pr15 hide">
                    <div class="form-group">
                        <input id="key" class="form-control" />
                    </div>
                    <button type="submit" class="btn btn-default" id="searchBtn">查询</button>
                </form>
                <div class="c-conf-tree pl15 pr15">
                    <div id="confTreeContainer" class="mt15">
                    </div>
                </div>
            </div>
            <!-- /分布拓扑  -->
            <div class="c-host-switch">
                <span class="glyphicon glyphicon-menu-left c-host-switch-img"></span>
            </div>
        </div>
        <div class="row host-main-right">
            <div class="col-md-12">
                <div class="conf-right-box">
                    <div class="conf-right-empty">
            <?php if($Level==2){
                include_once('noModule.php');
            }else{
                include_once('noSet.php');
            } ?>
            <!-- 没有集群 -->
            <!-- /没有业务 -->
            <?php include_once('newSet.php')?>
            <!-- 新建模块 -->
            <?php include_once('newModule.php')?>
            <!-- /新建模块 -->
            <!-- 集群属性修改 -->
            <?php include_once('editSet.php')?>
            <!-- /集群属性修改 -->
            <!-- 模块属性修改 -->
            <?php include_once('editModule.php')?>
            <!-- /模块属性修改 -->
                    </div>
        </div>
                </div>
        </div>
    </section>
</div>
<div class="control-sidebar-bg"></div>

<!-- 项目js文件 -->
<script id="treeview-template" type="text/kendo-ui-template">
    <span cid='#: item.id #'>#: item.text #</span>
    # if (item.type=='application') { #
    # if (!item.noset){ #
    <span class='creat-group-btn btn btn-success btn-xs' cid='#: item.id #' style='position:absolute;right:0;top:0;'><i class='fa fa-plus'></i> 集群</span>
    # }else{ #
    <span class='creat-module-btn btn btn-success btn-xs' cid='#: item.id #' style='position:absolute;right:0;top:0;'><i class='fa fa-plus'></i> 模块</span>
    # } #
    # } #
    # if (item.type=='set') { #
    <span class='creat-module-btn btn btn-success btn-xs' cid='#: item.id #' style='position:absolute;right:0;top:0;'><i class='fa fa-plus'></i> 模块</span>
    # } #</script>
<script src="<?php echo STATIC_URL;?>/static/assets/mCustomScrollbar-3.0.9/jquery.mCustomScrollbar.concat.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/ZeroClipboard/ZeroClipboard.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/bootstrap-switch-master/dist/js/bootstrap-switch.js" rel="stylesheet"></script>

<link rel="stylesheet" href="<?php echo STATIC_URL;?>/static/assets/jstree-3.1.1/dist/themes/default/style.min.css" />
<script src="<?php echo STATIC_URL;?>/static/assets/jstree-3.1.1/dist/jstree.min.js"></script>

<script src="<?php echo STATIC_URL;?>/static/js/app.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/js/topo.index.js?version=<?php echo $version;?>"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/select2-3.5.2/select2.js"></script>
<script >
    var appId = '<?php echo $appId;?>';
    var appName = '<?php echo $appName;?>';
    var topo = <?php echo $topo;?>;
    if(topo.length){
        for(var i=0,j=topo.length;i<j;i++){
            var newId = topo[i].id+"plus";
            topo[i].id = newId;
        }
    }

    var level = <?php echo $Level;?>;
    var desetid = <?php echo $deSetID;?>;
    var defaultapp = <?php echo $Default;?>;
    var emptys = <?php echo $emptys;?>;
    $("body").css("overflow-y","hidden");
</script>
