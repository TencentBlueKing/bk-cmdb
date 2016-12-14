<?php
/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */
?>
<link href="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/dataTables.bootstrap.css" rel="stylesheet"/>
<style>
    #privatediv{border: 1px solid #ccc}
    #private_wrapper{padding: 10px 7px 0px 7px;}
    #private thead tr th:nth-child(1){padding-right:0px;}
    .c-import-fpBtn{margin: 5px 8px;}
    .c-import-block .c-grid-toolbar > a, #private .c-grid-toolbar > form > a {margin-top: 5px !important;}
    .pagination>.active>a{color: #333!important}
    .import-button{color: #555;}
</style>
<div class="content-wrapper">
    <!-- 主面板 Header  -->
    <section class="content-header">
        <h1>资源池管理</h1>
    </section>
    <!-- 主面板 Main-->
    <section class="content pt0">
        <div class="row">
            <div class="col-sm-12 col-lg-12 c-import-table">
                <div class="c-import-info f16 pt5" role="alert">1.导入主机数据到资源池;  2.从资源池分配主机到<span class="c-import-business-name"><?php $defaultApp = $this->session->userdata('defaultApp');echo $defaultApp['ApplicationName'];?></span>的<span class="c-import-business-name">空闲机池</span>中</div>
                <div class="c-import-block">
                    <ul class="nav nav-tabs hide">
                        <li><a href="#privatediv" data-toggle="tab">资源池</a></li>
                    </ul>
                    <div class="tab-content">
                        <div class="tab-pane active" id="privatediv">
                            <div class="c-header c-grid-toolbar">
                                <a class="c-button c-button-icontext c-grid-quickDistribute" href="javascript:void(0)" disabled="disabled">
                                    <span class=""></span>分配至
                                </a>
                                <a class="c-button import-button" id="importPrivateHostByExcel" href="javascript:void(0)">
                                    <span class=""></span>导入主机
                                </a>
                                <a class="c-button c-button-icontext c-grid-delete" href="javascript:void(0)" disabled="true">
                                    <span class=""></span>删除
                                </a>
                                <div class="c-import-fpBtn">
                                    <div class="cc_switch_btn cc_switch_btn_right" data-fp="1">
                                        <img class="switch" src="/static/img/cc_switch_btn.png">
                                        <div class="num">未分配()</div>
                                    </div>
                                </div>

                            </div>
                            <table id="private" class="table table-bordered table-striped"></table>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </section>
</div>
<div class="control-sidebar-bg"></div>
</div>
<iframe name="upload_proxy" id="upload_proxy" style="display:none"></iframe>
<!-- 项目需要引用的js文件 -->
<script src="<?php echo STATIC_URL;?>/static/assets/js/jquery-1.10.2.min.js?version=<?php echo $version;?>"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/bootstrap-3.3.4/js/bootstrap.min.js?version=<?php echo $version;?>"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/mCustomScrollbar-3.0.9/jquery.mCustomScrollbar.concat.min.js?version=<?php echo $version;?>"></script>
<link href="<?php echo STATIC_URL;?>/static/assets/bootstrap-switch-master/dist/css/bootstrap3/bootstrap-switch.css?version=<?php echo $version;?>" rel="stylesheet">
<style>
    .host-state-switcher-parent {
        position: relative;
        width: 134px;
        height: 34px;
        float: right;
        margin: 3px 6px;
    }
</style>
<script src="<?php echo STATIC_URL;?>/static/assets/bootstrap-switch-master/dist/js/bootstrap-switch.js?version=<?php echo $version;?>"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/jquery.dataTables.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/dataTables.bootstrap.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/js/quickImport.js?version=<?php echo $version;?>"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/select2-3.5.2/select2.min.js?version=<?php echo $version;?>"></script>
