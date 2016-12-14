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
<style>
    #cluster_table thead tr th:nth-child(1),#module_table thead tr th:nth-child(1){padding-right:0px;}
    #cluster_table thead tr th:nth-child(2),#module_table thead tr th:nth-child(2){display: none;}
    #cluster_table tbody tr td:nth-child(2),#module_table tbody tr td:nth-child(2){display: none;}
    .c-import-block .pagination>.active>a{color: #333;}
</style>
<div class="content-wrapper">
    <!-- 主面板 Header  -->
    <section class="content-header">
        <h1>拓扑列表</h1>
    </section>
    <!-- 主面板 main-sidebarn-->
    <section class="content p20">

        <div class="row">
            <div class="col-sm-12 col-lg-12 c-import-table">

                <div class="c-import-block">
                    <ul class="nav nav-tabs">
                        <?php if(3 == $Level){?>
                            <li class="active"><a href="#cluster" data-toggle="tab">集群</a></li>
                            <li><a href="#module" data-toggle="tab">模块</a></li>
                        <?php }else{?>
                            <li class="active" ><a id="moduletab" href="#module" data-toggle="tab">模块</a></li>
                        <?php }?>
                    </ul>

                    <div class="tab-content">
                        <?php if(3 == $Level){?>
                        <div class="tab-pane fade in active" id="cluster">
                            <div class="searchButton">
                                <a class="c-button b_edit" href="javascript:void(0)" disabled><span class=""></span>修改</a>
                                <a class="c-button clone" href="javascript:void(0)" disabled><span class=""></span>克隆此集群</a>
                                <!-- <div class="c-import-search">
                                    <input id="filter-cluster" type="text" class="form-control pull-left" placeholder="搜索..." /><i class="glyphicon glyphicon-search"></i>
                                </div> -->
                            </div>
                            <!-- <div id="cluster_table" class="mb30"></div> -->
                            <div style="border:1px #ccc solid;padding:10px;padding-top:20px;">
                            <table id="cluster_table" class="table table-bordered table-striped mt10"></table>
                            </div>
                        </div>

                        <div class="tab-pane fade" id="module">
                            <div class="searchButton">
                                <a class="c-button b_edit" href="javascript:void(0)" disabled><span class=""></span>修改</a>
                                <!-- <div class="c-import-search">
                                    <input id="filter-module" type="text" class="form-control pull-left" placeholder="搜索..." /><i class="glyphicon glyphicon-search"></i>
                                </div> -->
                            </div>
                            <!-- <div id="module_table" class="mb30"></div> -->
                            <div style="border:1px #ccc solid;padding:10px;padding-top:20px;">
                            <table id="module_table" class="table table-bordered table-striped mt10"></table>
                            </div>
                        </div>
                    </div>
                        <?php }else{?>
                    <div class="tab-pane fade in active" id="module">
                        <div class="searchButton">
                            <a class="c-button b_edit" href="javascript:void(0)" disabled><span class=""></span>修改</a>
                            <!-- <div class="c-import-search">
                                <input id="filter-module" type="text" class="form-control pull-left" placeholder="搜索..." /><i class="glyphicon glyphicon-search"></i>
                            </div> -->
                        </div>
                        <!-- <div id="module_table" class="mb30"></div> -->
                        <div style="border:1px #ccc solid;padding:10px;padding-top:20px;">
                        <table id="module_table" class="table table-bordered table-striped mt10"></table>
                        </div>
                    </div>
                </div>
                        <?php }?>

                </div>
            </div>
        </div>
    </section>
</div>
<div class="control-sidebar-bg"></div>
<!-- 克隆 Modal -->
<div class="modal fade" id="cloneModal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">
    <div class="modal-dialog" role="document">
        <div class="modal-content">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title" id="myModalLabel">克隆集群</h4>
            </div>
            <div class="modal-body">
                <!-- <div class="form-group row">
                    <label for="cloneList" class="col-sm-3 control-label">克隆已有集群：</label>
                    <div class="col-sm-8">
                        <input type="hidden" id="cloneList">
                    </div>
                </div> -->
                <div class="form-group">
                    <label for="cloneTextarea">克隆出的集群列表</label>
                    <textarea class="form-control" id="cloneTextarea" placeholder="请用换行分割" ></textarea>
                </div>
                1. 克隆多个集群，请换行分割
                <br>
                2. 克隆操作会连同原集群包含的模块一起克隆
            </div>
            <p class="text-danger tc mt10" id="cloneerrtips" ></p>
            <div class="modal-footer tc">
                <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                <button type="button" class="btn btn-primary btn-save">确定</button>
            </div>
        </div>
    </div>
</div>
<!-- Modal -->

<!-- 批量修改 Modal -->
<div class="modal fade" id="b_edit_modal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">
    <div class="modal-dialog" role="document">
        <div class="modal-content c-edit-all">
            <div class="modal-header">
                <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title" id="myModalLabel">批量修改</h4>
            </div>
            <div class="modal-body pb0">
                <div class="form-horizontal">
                    <table class="table table-bordered table-hovered cluster-edit-cluster">
                        <thead>
                        <tr>
                            <th style="width:40px;" class="select-area tc"> <input type="checkbox" class="selectAll"> </th>
                            <th style="width:120px;" class="tc">属性名</th>
                            <th class="host-attr-val tc">属性值</th>
                        </tr>
                        </thead>
                        <tbody>
                        <tr>
                            <td class="lh34"><input type="checkbox" data-type="select" data-id="e_EnviType" class="tr_select"  id="c_EnviType"></td>
                            <td class="lh34">环境类型</td>
                            <td>
                                <input type="hidden" id="e_EnviType" class="e_input" placeholder="请输入环境类型">
                                <div data-for="e_EnviType" class="edit-area-mask"></div>
                            </td>
                        </tr>
                        <tr>
                            <td class="lh34"><input type="checkbox" data-type="select" data-id="e_ServiceStatus" class="tr_select" id="c_ServiceStatus"></td>
                            <td class="lh34">服务状态</td>
                            <td>
                                <input type="hidden" id="e_ServiceStatus" class="e_input" placeholder="请输入服务状态" >
                                <div data-for="e_ServiceStatus" class="edit-area-mask"></div>
                            </td>
                        </tr>
                        <tr>
                            <td class="lh34"><input type="checkbox" data-type="text" data-id="e_Capacity" class="tr_select" id="c_Capacity"></td>
                            <td class="lh34">设计容量</td>
                            <td>
                                <input type="text" id="e_Capacity" class="form-control regNumber" class="e_input" placeholder="" disabled>
                                <div data-for="e_Capacity" class="edit-area-mask"></div>
                            </td>
                        </tr>
                        <tr>
                            <td class="lh34"><input type="checkbox" data-type="text" data-id="e_Openstatus" class="tr_select" id="c_Openstatus"></td>
                            <td class="lh34">Openstatus</td>
                            <td>
                                <input type="text" id="e_Openstatus" class="form-control" class="e_input" maxlength="16" placeholder="" disabled>
                                <div data-for="e_Openstatus" class="edit-area-mask"></div>
                            </td>
                        </tr>
                        </tbody>
                    </table>

                    <table class="table table-bordered table-hovered cluster-edit-module none">
                        <thead>
                        <tr>
                            <th style="width:40px;" class="select-area tc"> <input type="checkbox" class="selectAll"> </th>
                            <th style="width:120px;" class="tc">属性名</th>
                            <th class="host-attr-val tc">属性值</th>
                        </tr>
                        </thead>
                        <tbody>
                        <tr>
                            <td class="lh34"><input type="checkbox" data-id="e_Operator" class="tr_select" id="c_Operator"></td>
                            <td class="lh34">维护人</td>
                            <td>
                                <input type="hidden" id="e_Operator" class="e_input" placeholder="请输入维护人">
                                <div data-for="e_Operator" class="edit-area-mask"></div>
                            </td>
                        </tr>
                        <tr>
                            <td class="lh34"><input type="checkbox" data-id="e_BakOperator" class="tr_select" id="c_BakOperator"></td>
                            <td class="lh34">备份维护人</td>
                            <td>
                                <input type="hidden" id="e_BakOperator" class="e_input" placeholder="请输入备份维护人">
                                <div data-for="e_BakOperator" class="edit-area-mask"></div>
                            </td>
                        </tr>
                        </tbody>
                    </table>
                    <p class="text-danger tc mt10" id="errtips" ></p>
                </div>
            </div>
            <div class="modal-footer pt0">
                <button type="button" class="btn btn-default" data-dismiss="modal">取消</button>
                <button type="button" class="btn btn-primary btn-save">确定</button>
            </div>
        </div>
    </div>
</div>
<!-- Modal -->
<link href="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/dataTables.bootstrap.css" rel="stylesheet"/>
<script src="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/jquery.dataTables.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/dataTables.bootstrap.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/mCustomScrollbar-3.0.9/jquery.mCustomScrollbar.concat.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/ZeroClipboard/ZeroClipboard.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/bootstrap-switch-master/dist/js/bootstrap-switch.js" rel="stylesheet"></script>
<script src="<?php echo STATIC_URL;?>/static/js/app.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/select2-3.5.2/select2.js"></script>
<script src="<?php echo STATIC_URL;?>/static/js/top.list.js?version=<?php echo $version;?>"></script>
<script >
    var userlist = <?php echo $list;?>;
    var userkv   = <?php echo $kv;?>;
    var alevel = <?php echo $Level; ?>;
    level = '';

    $(function () {

    })
</script>
