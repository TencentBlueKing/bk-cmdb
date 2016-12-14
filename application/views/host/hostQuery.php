<?php
/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */
?>
<!-- datatables css -->
<link href="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/dataTables.bootstrap.css" rel="stylesheet"/>
<!-- datetimepicker -->
<link href="<?php echo STATIC_URL;?>/static/css/bootstrap-datetimepicker.min.css" rel="stylesheet">
<link href="<?php echo STATIC_URL;?>/static/assets/mCustomScrollbar-3.0.9/jquery.mCustomScrollbar.css" rel="stylesheet">
<div class="content-wrapper">
    <style type="text/css">
    .c-grid-moveIp{
        display: inline-block;
        border: 0;
        background: none;
        margin: 0!important;
        width: 200px;
        text-align: left;
        border: 1px solid #CCC;
         padding: 2px 0;
        background: #FFF;
        position: relative;
    }
    .c-grid-moveIp button{
        position: absolute;
        left: -1px;
        top: -1px;
    }
/*    .c-grid-moveIp.able button{
        display: inline-block;
    }*/
    </style>
    <!-- 主面板 Header  -->
    <!-- 主面板 Main-->
    <section class="content">
        <div class="host-sidebar-left">
            <!-- 分布拓扑 -->
            <div class="c-host-side" id="c_host_side">
                <h4>拓扑视图</h4>
                <div class="c-tree-box pl15 pr15">
                    <div id="treeContainer" class="mt15"></div>
                </div>
                <!-- 空闲机池 -->
                <div class="free-group pl15 pr15">
                    <h4>空闲机池</h4>
                    <div id="emptyContainer" class="mt15 k-widget k-treeview">
                    </div>
                </div>
                <!-- /空闲机池 -->
            </div>
            <!-- /分布拓扑  -->
            <div class="c-host-switch">
                <span class="glyphicon glyphicon-menu-left c-host-switch-img"></span>
            </div>
        </div>
        <div class="row host-main-right">
            <div class="col-md-12">
                <div class="c-search-box">
                    <div class="panel-heading" id="collapseOneBtn">
                        <h4 class="panel-title" style="display: inline-block;">查询条件 </h4>
                        <div class="c-collapse-btn-group">
                            <button class="btn btn-success btn-xs" id="filter_module_mine">我</button>
                            <button class="btn btn-success btn-xs" id="filter_module_all">ALL</button>
                            <button class="btn btn-success btn-xs" id="filter_module_empty">空闲机</button>
                            <span class="fa fa-angle-down" aria-hidden="true"></span>
                        </div>
                    </h4>
                    </div>
                    <div id="collapseOne" class="panel-collapse collapse">
                        <div class="search-content">
                            <div class="col-lg-6">
                                <div class="input-group pb10 input-group-height">
                                    <label for="InnerIP" class="input-group-addon">内网IP：</label>
                                    <div class="textarea-down"><div>
                                    <textarea placeholder="请输入内网IP" value="<?php echo isset($InnerIP) ? $InnerIP : '';?>" id="InnerIP" name="InnerIP" class="form-control column_filter" filter-name="InnerIP"><?php echo isset($InnerIP) ? $InnerIP : '';?></textarea>
                                    </div></div>
                                    <!-- <label class="checkbox-inline input-group-addon pl30 pr10"><input type="checkbox" checked="selected" name="IfInnerIPexact">精确</label> -->
                                </div>
                            </div>
                            <div class="col-lg-6">
                                <div class="input-group pb10 input-group-height">
                                    <label for="OuterIP" class="input-group-addon">外网IP：</label>
                                    <div class="textarea-down"><div>
                                    <textarea placeholder="请输入外网IP" value="<?php echo isset($OuterIP) ? $OuterIP : '';?>" id="OuterIP" name="OuterIP" class="form-control column_filter" filter-name="OuterIP"><?php echo isset($OuterIP) ? $OuterIP : '';?></textarea>
                                    </div></div>
                                    <!-- <label class="checkbox-inline input-group-addon pl30 pr10"><input type="checkbox" checked="selected" name="IfOuterexact">精确</label> -->
                                </div>
                            </div>
                            <?php if($level == 3){?>
                            <div class="input-group col-xs-12 pb10 pl15 pr15">
                                <div class="input-group-addon">集群名称</div>
                                <select placeholder="(不选为全部)" class="select2_src form-control" multiple="" name="SetID" id="set_select" style="display: none;">
                                    <?php
                                        foreach($set_select as $_sk=>$_sv){
                                            echo '<option value="'.$_sv['SetID'].'">'.$_sv['SetName'].'</option>';
                                        }
                                    ?>
                                </select>
                            </div>
                            <?php }?>
                            <div class="input-group col-xs-12 pb10 pl15 pr15">
                                <div class="input-group-addon">模块名称</div>
                                <select placeholder="(不选为全部)" class="select2_src form-control column_filter"  filter-name="ModuleName" multiple="" name="ModuleID" id="module_select" style="display: none;">
                                    <?php
                                        foreach($module_select as $_mk=>$_mv){
                                            foreach($set_select as $_sk=>$_sv){
                                                if($_sv['SetID']==$_mv['SetID']){
                                                    echo '<option value="'.$_mv['ModuleID'].'">'.($level==3?$_sv['SetName'] .'-':''). $_mv['ModuleName'].'</option>';
                                                }
                                            }
                                        }
                                    ?>
                                </select>
                            </div>
                            <?php foreach ($customerQueryFields as $key=>$item) : ?>
                                    <div class="col-lg-6" id="<?php echo $key;?>Label" style="<?php if(!in_array($key, $DefaultField)){echo 'display:none;';}?>" >
                                        <div class="input-group pb10 input-group-height">
                                            <label class="input-group-addon"><?php echo $item;?></label>

                                            <?php if(in_array($key, array('Operator', 'BakOperator', 'OSName'))) { ?>
                                                    <input type="hidden" style="width:100%;" class="operator-select column_filter" filter-name="<?php echo $key;?>" id="<?php echo $key;?>">
                                                    <div class="edit-area"></div>
                                            <?php }elseif(in_array($key, array('CreateTime', 'DeadLineTime'))){   ?>
                                                    <div class="textarea-down">
                                                        <div>
                                                            <!-- 日期选择器(Bootstrap)-2 Start -->
                                                            <div class="input-group date" id="<?php echo $key;?>" data-date="" data-date-format="dd MM yyyy" data-link-field="<?php echo $key;?>" data-link-format="yyyy-mm-dd">
                                                                <span class="input-group-addon"><span class="glyphicon glyphicon-calendar"></span></span>
                                                                <input class="form-control column_filter" filter-name="CreateTime" size="16" type="text" value="">
                                                                <input type="hidden" id="<?php echo $key;?>" value="">
                                                                <br>
                                                            </div>
                                                            <!-- 日期选择器(Bootstrap)-2 End -->
                                                        </div>
                                                    </div>
                                            <?php }else {   ?>
                                            <div class="textarea-down">
                                                <div>
                                                   <textarea placeholder="请输入<?php echo $item;?>" value="" id="<?php echo $key;?>" name="<?php echo $key;?>" class="form-control column_filter" filter-name="<?php echo $key;?>"></textarea>
                                                </div>
                                            </div>
                                            <?php } ?>
                                        </div>
                                    </div>
                                <?php endforeach;?>

                                    <div class="clearfix"></div>
                                    <div class="col-xs-12 text-center" style="float:none;">
                                        <div class="topology-conditionMore">
                                            <button class="btn btn-default mr10 conditionMore-button">更多条件<span class="glyphicon glyphicon-triangle-bottom"></span></button>
                                            <div class="conditionMore-div tl">

                                                <div class="close-div"><span class="button-cancel">×</span></div>

                                                <h5 class="mt0">必选条件</h5>
                                                <div class="select selected">
                                                    <div class="row">
                                                        <div class="col-4"><label class="iCheckbox_square"><input type="checkbox" disabled="" checked="checked" / ><span>内网IP</span></label></div>
                                                        <div class="col-4"><label class="iCheckbox_square"><input type="checkbox" checked="checked" disabled="" /><span>外网IP</span></label></div>
                                                    </div>
                                                    <div class="row">
                                                        <div class="col-4"><label class="iCheckbox_square"><input type="checkbox" checked="checked" disabled="" /><span>集群名称</span></label></div>
                                                        <div class="col-4"><label class="iCheckbox_square"><input type="checkbox" checked="checked" disabled="" /><span>模块名称</span></label></div>
                                                    </div>
                                                </div>
                                                <h5>可选条件</h5>
                                                <div class="select">
                                                    <?php $index = 1;?>
                                                <?php foreach ($customerQueryFields as $key=>$item) : ?>
                                                    <?php if($index == 1) { ?>
                                                    <div class="row pb10">
                                                    <?php }?>
                                                        <div class="col-4"><label class="iCheckbox"><input type="checkbox" <?php if(in_array($key, $DefaultField)){echo 'checked="checked"';}?> data-rel="#<?php echo $key;?>Label" /><span><?php echo $item;?></span></label></div>
                                                    <?php if($index % 3 == 0 && $index != 1) { ?>
                                                    </div>
                                                    <div class="row pb10">
                                                    <?php } $index++;?>
                                                 <?php endforeach;?>
                                                 </div>
                                                </div>

                                            </div>
                                        </div>
                                        <button class="btn btn-default" id="host_query_reset">重置</button>
                                        <button class="btn btn-primary" id="host_query_submit">查询</button>
                                    </div>
                        </div>
                    </div>
                </div>
                <!-- 主机列表 -->
                <div class="c-edit-set">
                    <h4 class="c-search-title">查询结果  <i class="fa fa-gear cc-column-edit"></i></h4>
                    <div class="column-edit-block" style="display:none">
                        <h4 class="pl15">主机显示字段
                         <button class="btn btn-info mr10 edit-btn-save">保存</button>
                         <button class="btn btn-default edit-btn-cancel">取消</button>
                        </h4>
                        <div class="row title">
                            <div class="col-9">可选字段</div>
                            <div class="col-3">当前选定的字段</div>
                        </div>
                        <div class="row detail">
                            <div class="col-9 column-before">
                                <h5>必选条件</h5>
                                <div class="select selected">
                                    <div class="row">
                                        <div class="col-6"><label class="iCheckbox_square"><input type="checkbox" disabled="" checked="checked" / ><span>内网IP</span></label></div>
                                        <div class="col-6"><label class="iCheckbox_square"><input type="checkbox" checked="checked" disabled="" /><span>外网IP</span></label></div>
                                    </div>
                                    <div class="row">
                                        <div class="col-6"><label class="iCheckbox_square"><input type="checkbox" checked="checked" disabled="" /><span>集群名称</span></label></div>
                                        <div class="col-6"><label class="iCheckbox_square"><input type="checkbox" checked="checked" disabled="" /><span>模块名称</span></label></div>
                                    </div>
                                </div>
                                <h5>可选条件</h5>
                                <div class="select">
                                    <?php $index=0; $FieldsCount = count($HostFields);
                                        $columnsArr = json_decode($columns);
                                       foreach($HostFields as $fie=>$des) {

                                           if($index%2 == 0) { ?>
                                    <div class="row pb10">
                                        <div class="col-6"><label class="iCheckbox"><input type="checkbox" target-rel="<?php echo $fie;?>" <?php if(in_array($fie,$columnsArr)){ ?>checked="checked" <?php }?>/><span><?php echo $des;?></span></label></div>
                                        <?php }elseif($index%2 == 1){   ?>
                                        <div class="col-6"><label class="iCheckbox"><input type="checkbox" target-rel="<?php echo $fie;?>" <?php if(in_array($fie,$columnsArr)){ ?>checked="checked" <?php }?>/><span><?php echo $des;?></span></label></div>
                                    </div>
                                <?php  } $index++; } if($FieldsCount%2 ==1){ ?>
                                </div>
                                    <?php }?>

                            </div>
                            </div>
                            <div class="col-3 column-after">
                                <ul id="sortable" style="overflow:auto;">
                                    <li name='InnerIP' class="ui-state-disabled"><i class="fa fa-ellipsis-v"></i><span>内网IP</span></li>
                                    <li name='OuterIP' class="ui-state-disabled"><i class="fa fa-ellipsis-v"></i><span>外网IP</span></li>
                                    <li name='SetName' class="ui-state-disabled"><i class="fa fa-ellipsis-v"></i><span>集群名称</span></li>
                                    <li name='ModuleName' class="ui-state-disabled"><i class="fa fa-ellipsis-v"></i><span>模块名称</span></li>
                                    <?php $columnsArr = json_decode($columns);foreach($HostFields as $fie=>$hf){ if(in_array($fie,$columnsArr)){?>
                                    <li name="<?php echo $fie;?>"><i class="fa fa-ellipsis-v"></i><span><?php echo $hf;?></span><i class="fa fa-close list-close"></i></li>
                                    <?php  }} ?>
                                </ul>
                            </div>
                        </div>
<!--                        <div class="row butn">-->
<!--                            <button class="btn btn-info mr10 edit-btn-save">保存</button>-->
<!--                            <button class="btn btn-default edit-btn-cancel">取消</button>-->
<!--                        </div>-->
                    </div>
                <div id="dialogs">  </div>
                </div>
                <div id="host-list" class="pb10" style="background:white;">
                <!-- 表格(DataTables)-4 Start -->
                <div class="table-box">
                <div class="c-header c-grid-toolbar pl10">
                    <a class="c-button c-button-icontext c-grid-copyIP" href="javascript:void(0);" disabled="disabled">
                        <span class=" "></span>
                        <div class="copy-menu-box">
                            <div>复制<span class="caret"></span></div>
                            <ul class="copy-menu">
                                <li class="copy-inner-ip" data-type="InnerIP">复制内网IP</li>
                                <li class="copy-outer-ip" data-type="HostID">复制外网IP</li>
                                <li class="copy-asset-id" data-type="copy-things-num">复制固资编号</li>
                            </ul>
                        </div>
                    </a>
                    <a class="c-button c-button-icontext c-grid-batEdit" id="batEdit" href="javascript:void(0);" disabled="disabled"><span class=" "></span>修改</a>
                    <div class="move-ip-panel" style="position: relative;display: inline-block;width: 200px;">
                        <a class="c-button c-button-icontext c-grid-moveIp" id="moveIp" href="javascript:void(0);" disabled="disabled"><span style="margin-left: 4px;font-size: 14px;" class="selectNote">转移主机至</span><span class="ui-icon ui-icon-triangle-1-s fr" style="opacity: .5;margin-right: 3px;"></span></a>
                        <div class="downList" id="downList" style="width: 370px;position: absolute;left: 0;top: 30px;width: 100%;z-index: 9999;background: #FFF;border: 1px solid #CCC;padding: 5px;">
                            <div class="searchPanel"><input type="text"></div>
                            <ul style="list-style: none;padding: 0; margin-top: 8px;" id="myScroll">
                                <?php
                                    foreach($module_select as $_mk=>$_mv){
                                        if($_mv['ModuleName']=='空闲机'){
                                            continue;
                                        }

                                        foreach($set_select as $_sk=>$_sv){
                                            if($_sv['SetID']==$_mv['SetID']){
                                                echo '<li value="'.$_mv['ModuleID'].'"><input type="checkbox" class="ui-checkbox">'.($level==3?$_sv['SetName'] .'-':''). $_mv['ModuleName'].'</li>';
                                            }
                                        }
                                    }
                                ?>
                            </ul>
                            <div class="text-center" id="modSelectMenu">
                                <button class="btn btn-xs btn-default mr10 operationSelect">取消</button>
                                <button class="btn btn-xs btn-primary operationSelect">转移</button>
                            </div>
                        </div>
                    </div>

                    <a class="c-button c-button-icontext c-grid-saveAsExcel" href="javascript:void(0);">
                        <span class=" "></span>导出Excel
                    </a>
                    <a class="c-button c-button-icontext c-grid-batRes" id="batRes" href="javascript:void(0);" style="display:none;float:right" disabled="disabled">
                        <span class=" "></span>上交
                    </a>
                    <a class="c-button c-button-icontext c-grid-batDel" id="batDel" href="javascript:void(0);" disabled="disabled">
                        <span class=" "></span>移至空闲机
                    </a>
                </div>
                <table id="table_topo" class="table table-bordered table-striped table-responsive">
                </table>
                </div>

                    <!-- 表格(DataTables)-4 End -->
                </div>
                <!-- /主机列表 -->
            </div>
        </div>
    </section>
</div>
<div class="control-sidebar-bg"></div>
</div>
<div id="window"></div>

<!-- Modal -->
<div class="modal fade" id="batEdit_window" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">
    <div class="modal-dialog" role="document">
        <div class="modal-content c-edit-all">
            <div class="modal-header">
                <button type="button" class="close"><span aria-hidden="true">&times;</span></button>
                <h4 class="modal-title">批量修改</h4>
            </div>
            <div class="modal-body">
                <table class="table table-bordered table-hovered">
                    <thead>
                        <tr>
                            <th style="width:40px;" class="select-area">
                                <input type="checkbox" data='selectAll' />
                            </th>
                            <th style="width:120px;">属性名</th>
                            <th class="host-attr-val">属性值</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td class="select-area">
                                <input type="checkbox" />
                            </td>
                            <td>主机名称</td>
                            <td class="edit-area">
                                <input type="text" value='' class="form-control" id="moduleHost_HostName" disabled style="cursor:default" maxlength="32">
                            </td>
                        </tr>
                        <tr>
                            <td class="select-area">
                                <input type="checkbox" data-field="Operator"/>
                            </td>
                            <td>维护人员</td>
                            <td class="edit-area">
                                <!-- <input type="text" value='' class="form-control" id="moduleHost_Operator"> -->
                                <div class="form-control" id="moduleHost_Operator"></div>
                                <div class="edit-area-mask"></div>
                            </td>
                        </tr>
                        <tr>
                            <td class="select-area">
                                <input type="checkbox" data-field="BakOperator"/>
                            </td>
                            <td>备份维护人</td>
                            <td class="edit-area">
                                <!-- <input type="text" value='' class="form-control" id="moduleHost_BakOperator" disabled> -->
                                <div class="form-control" id="moduleHost_BakOperator"></div>
                                <div class="edit-area-mask"></div>
                            </td>
                        </tr>
                        <tr>
                            <td class="select-area">
                                <input type="checkbox" />
                            </td>
                            <td>备注信息</td>
                            <td class="edit-area">
                                <input type="text" value='' class="form-control" id="moduleHost_Description" disabled style="cursor:default" maxlength="256">
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div class="modal-footer">
                <button type="button" class="btn btn-default" id="moduleHostHide">取消</button>
                <button type="button" class="btn btn-primary" id="moduleHostSubmit">修改</button>
            </div>
        </div>
    </div>
</div>
<script>
    var tablesFields = <?php echo $tablesFields;?>;
</script>
<!-- 项目需要引用的js文件 -->
<script src="<?php echo STATIC_URL;?>/static/assets/js/jquery-1.10.2.min.js?version=<?php echo $version;?>"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/bootstrap-3.3.4/js/bootstrap.min.js?version=<?php echo $version;?>"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/mCustomScrollbar-3.0.9/jquery.mCustomScrollbar.concat.min.js?version=<?php echo $version;?>"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/ZeroClipboard/ZeroClipboard.min.js?version=<?php echo $version;?>"></script>
<!-- jstree -->
<link rel="stylesheet" href="<?php echo STATIC_URL;?>/static/assets/jstree-3.1.1/dist/themes/default/style.min.css" />
<script src="<?php echo STATIC_URL;?>/static/assets/jstree-3.1.1/dist/jstree.min.js"></script>
<!-- datetimepicker -->
<script src="<?php echo STATIC_URL;?>/static/assets/datetimepicker/bootstrap-datetimepicker.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/datetimepicker/bootstrap-datetimepicker.zh-TW.js" charset="UTF-8"></script>

<script src="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/jquery.dataTables.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/dataTables.bootstrap.min.js"></script>

<!-- jquery-ui js -->
<script type="text/javascript" src="/static/assets/selectmenu/jquery-ui.min.js?version=<?php echo $version;?>"></script>

<!-- 项目js文件 -->
<script src="<?php echo STATIC_URL;?>/static/assets/icheck-1.x/icheck.min.js?version=<?php echo $version;?>"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/select2-3.5.2/select2.min.js?version=<?php echo $version;?>"></script>
<script src="<?php echo STATIC_URL;?>/static/js/hostQuery.js?version=<?php echo $version;?>"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/jquery-ui-1.11.0.custom/jquery-ui.js?version=<?php echo $version;?>"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/mCustomScrollbar-3.0.9/jquery.mCustomScrollbar.concat.min.js?version=<?php echo $version;?>"></script>

<script>
    var currUser = <?php echo json_encode($this->session->userdata());?>;
    var userList = <?php echo $userList;?>;
    var static_url = "<?php echo STATIC_URL; ?>";
    var empty = {"id":"", text:"请选择"};
    userList.splice(0, 0, empty);var OSName= <?php echo $OSName;?>;
    OSName.splice(0, 0, empty);
    var customerQueryFields= <?php echo json_encode($customerQueryFields);?>;
    $('#set_select').selectDialog({});
    $('#module_select').selectDialog({});
    //获取编辑后的 列
    function setColumns(){
        var columns=[];
        $('#sortable li').each(function(){
            var text=$(this).find('span').text();
            var name=$(this).attr('name');
            columns.push({title:text,field:name,width:150});
        })

        var param = {};
        param['ApplicationID'] = cookie.get('defaultAppId');
        param['DefaultColumn'] = [];
        for(var i in columns){
            param['DefaultColumn'].push(columns[i].field);
        }

        $.ajax({
            url:'/UserCustom/setUserCustom/',
            dataType:'json',
            data:param,
            method:'post',
            success:function(response){
                CC.host.hostlist.defaultColumns = param['DefaultColumn'];
                console.log(param['DefaultColumn'])
                var table = $("#table_topo").DataTable();
                var ColumnData = param['DefaultColumn'];
                table.columns().visible(false);
                table.columns(0).visible(true);
                for (var i = 0; i < ColumnData.length; i++) {
                    table.column('.'+param.DefaultColumn[i]).visible(true);
                };
                // CC.host.hostlist.init();
                return true;
            }
        });
    }
    $(function(){
        $(".c-host-switch").click(function(){
            if($('.c-host-switch-img').hasClass("glyphicon-menu-left")){
                $('.c-host-switch-img').removeClass('glyphicon-menu-left').addClass('glyphicon-menu-right');
                $(".host-sidebar-left").animate({"width": "0px"}, "fast");
                $(".host-main-right").animate({"right": "+=320px","width": "+=320px"}, "fast");
            }else if($('.c-host-switch-img').hasClass("glyphicon-menu-right")){
                $('.c-host-switch-img').removeClass('glyphicon-menu-right').addClass('glyphicon-menu-left');
                $(".host-sidebar-left").animate({"width": "320px"}, "fast");
                $(".host-main-right").animate({"right": "-=320px","width": "-=320px"}, "fast");
            }
        })
        // 拓扑视图自定义滚动条
        $(".c-tree-box").mCustomScrollbar({
            //setHeight: 400, //设置高度
            theme: "minimal-dark" //设置风格
        }).css('maxHeight', 400);
        // CC.host.topology.init(<?php echo $topo;?>);
        // console.log(<?php echo $topo;?>);
        CC.host.hostlist.defaultColumns=<?php echo $columns;?>;
        if($('#InnerIP').val().length>0 || $('#OuterIP').val().length>0 || $('#AssetID').val().length>0){
            $('#collapseOneBtn').trigger('click');
            $('#host_query_submit').trigger('click');
        }else if(location.hash==='#em'){
            $('.module:first', '#emptyContainer').click();
        }else{
            $('.application', "#treeContainer").click();
        }

        //下拉选择框
        $("#e1 , #e2 ,#e3 , #e4 ,#e5").select2();

        $("#Operator").select2({
                placeholder:'请选择',
                data:window.userList,
                formatResult:function format(state) {
                return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
            },
            formatSelection:function format(state) {
                return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
            }
        }).select2('val','');

        $("#BakOperator").select2({
                placeholder:'请选择',
                data:window.userList,
                formatResult:function format(state) {
                return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
            },
            formatSelection:function format(state) {
                return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
            }
        }).select2('val','');

        $("#OSName").select2({
            placeholder:'请选择',
            data:window.OSName,
            formatResult:function format(state) {
                return state.text;
            },
            formatSelection:function format(state) {
                return state.text;
            }
        }).select2('val','');
        //checkbox和radio
        $('.iCheckbox').iCheck({
            checkboxClass: 'icheckbox_minimal-blue',
            radioClass: 'iradio_minimal-blue'
        });

        $('.iCheckbox_square').iCheck({
            checkboxClass: 'icheckbox_square-blue',
        });

        //取消和关闭按钮事件
        $('.button-cancel').click(function(){
            $('.conditionMore-div').hide();
        })

        //更多条件按钮事件
        $('.conditionMore-button').bind("click",function(event){
            $('.conditionMore-div').show();
            event.stopPropagation();
        })

        $('input[data-rel]').on('ifChanged', function(event){
            var id=$(this).attr('data-rel');
            $(''+id+'').toggle();
            var field = id.replace('Label','');
            var field = field.replace('#','');
            if($(this).prop("checked")){
                $.ajax({
                    url: "/host/setDefaultField",
                    type: "POST",
                    data: "key=" + field + "&type=a&field=DefaultField",
                    dataType: "json",
                });
            }else{
                $.ajax({
                    url: "/host/setDefaultField",
                    type: "POST",
                    data: "key=" + field + "&type=d&field=DefaultField",
                    dataType: "json"
                });
            }

            var id = id.replace('Label','');
            var selectarr = ["#Operator", "#BakOperator"];
            var result = $.inArray(id, selectarr);
            if(result >= 0){
                $(''+id+'').select2('val','');
            }

            var arr = ["#HostName", "#SN", "#CreateTime","#DeadLineTime"];
            var resultnew = $.inArray(id, arr);
            if(resultnew >= 0){
                $(''+id+'').val('');
            }
        });

        $(document).bind("click",function(event){
            var target = $(event.target);
            if (!target.closest('.conditionMore-div').length){
                $('.conditionMore-div').hide();
            }
        });

        //设置按钮
        $('.cc-column-edit').on('click',function(){
            $('.column-edit-block').toggle();
        })
        //关闭、取消 按钮
        $('.cc-column-close,.edit-btn-cancel').on('click',function(){
            $('.column-edit-block').hide();
        })
        //编辑显示字段
        $('input[target-rel]').on('ifChecked', function(event){
            var text=$(this).closest('label').find('span').text();
            var name=$(this).attr('target-rel');
            $('.column-after > ul').append('<li name='+name+'><i class="fa fa-ellipsis-v"></i><span>'+text+'</span><i class="fa fa-close list-close"></i>'+
            '</li>');
            $('.column-edit-block .list-close').on('click',function(){
                var name=$(this).closest('li').attr('name');
                $('input[target-rel='+name+']').iCheck('uncheck');
                $(this).closest('li').remove();
            })
        });
        //编辑显示字段
        $('input[target-rel]').on('ifUnchecked', function(event){
            var name=$(this).attr('target-rel');
            $('.column-after > ul > li').remove('li[name='+name+']');
        });
        //拖拽
        $('#sortable').sortable({
            items: "li:not(.ui-state-disabled)"
        });
        //保存按钮
        $(".edit-btn-save").on('click',function(){
            //隐藏弹出层
            $('.column-edit-block').hide();
            setColumns();
        });

        //编辑后行的删除按钮
        $('.column-edit-block .list-close').on('click',function(){
            var name=$(this).closest('li').attr('name');
            $('input[target-rel='+name+']').iCheck('uncheck');
            $(this).closest('li').remove();
        })

        $(".column-edit-block .detail").mCustomScrollbar({
            setHeight:575, //设置高度
            theme:"minimal-dark" //设置风格
        });
    });
</script>
<script>
$(function (){
// 日期选择器
$('#CreateTime').datetimepicker({
    language:  'zh-TW',
    weekStart: 1,
    todayBtn:  true,
    autoclose: true,
    todayHighlight: true,
    startView: 2,
    minView: 2,
    forceParse: false,
    format:"yyyy-mm-dd",
});
// 日期选择器 end
// 拓扑视图 表格
var language = {
  search: '搜索：',
  lengthMenu: "每页显示 _MENU_ 记录",
  zeroRecords: "没找到相应的数据！",
  info: "分页 _PAGE_ / _PAGES_",
  infoEmpty: "暂无数据！",
  infoFiltered: "(从 _MAX_ 条数据中搜索)",
  paginate: {
    first: '首页',
    last: '尾页',
    previous: '上一页',
    next: '下一页',
  }
}
var tableData={
    lengthChange: false, //不允许用户改变表格每页显示的记录数
    pageLength : 20, //每页显示几条数据
    lengthMenu: [5, 10, 20], //每页显示选项
    pagingType: 'full_numbers',
    // ajax : '/static/js/j1.json',
    ajax : {
        url:'/host/getHostById',
        data:{}
    },
    destroy: true,
    columnReorder: true,
    responsive: true,
    columns : [
      {title:'<input type="checkbox" class="host-td-checkbox host-check-all">',data:null,"orderable":false,class:"hots-th-checkbox",render:function ( data, type, row){
            return '<input type="checkbox" class="host-td-checkbox" value="'+data.HostID+'"/>';
      }},
      {title:'内网IP',data:null,class:"InnerIP",render:function ( data,type,row ){
        // InnerIP
            return '<a href="javascript:void(0)" class="a-innerip" title="'+data.InnerIP+'">'+data.InnerIP+'</a>'
      }},
      {title:'外网IP',data:"OuterIP",class:"OuterIP"},
      {title:'集群名称',data:"SetName",class:"SetName"},
      {title:'模块名称',data:"ModuleName",class:"ModuleName"},
      {title:'内存',data:"Mem",class:"Mem"},
      {title:'操作系统',data:"OSName",class:"OSName"},
      {title:'维护人',data:"Operator",class:"Operator"},
      {title:'固资编号',data:"AssetID",class:"AssetID",visible:false},
      {title:'主机ID',data:"HostID",class:"HostID",visible:false},
      {title:'Cpu',data:"Cpu",class:"Cpu",visible:false},
      {title:'机房城市',data:"Region",class:"Region",visible:false},
      {title:'购买时间',data:"CreateTime",class:"CreateTime",visible:false},
      {title:'可用区ID',data:"ZoneID",class:"ZoneID",visible:false},
      {title:'可用区',data:"ZoneName",class:"ZoneName",visible:false},
      {title:'备注',data:"Description",class:"Description",visible:false},
      {title:'备份维护人',data:"BakOperator",class:"BakOperator",visible:false},
      {title:'运行状态',data:"Status",class:"Status",visible:false},
      {title:'主机名称',data:"HostName",class:"HostName",visible:false},
      {title:'设备类型',data:"DeviceClass",class:"DeviceClass",visible:false},
      {title:'负责人',data:"Customer001",class:"Customer001",visible:false}
    ],
    language:language,
    "initComplete":function (){
        $(this).css('width', '100%');
        setColumns();
        // 表格渲染时判断checkbox是否选中
        $('#table_topo').on('draw.dt', function() {
            var data = $('#table_topo').DataTable().data();
            for (var i=0,len=data.length; i<len; i++){
                var selector = 'input[value='+data[i].HostID+']';
                if (data[i].Checked){
                    $(selector).prop('checked','checked');
                }else{
                    $(selector).prop('checked',false);
                }
            }
        });
    }
}
// 拓扑视图 表格 end

// 自定义 search 事件
$('.column_filter').on('search', function() {
    var table = $('#table_topo').DataTable();
    var thisName="."+$(this).attr('filter-name');
    // 模块名称匹配条件
    var moduleReg='';
    // 查询过滤
    if (table.column(thisName).search() !== this.value ) {
        table = table.column(thisName).search(this.value).draw();
    }
    if(thisName=='.ModuleName'){//模块名称
        $('[filter-name="ModuleName"]').find("option:selected").each(function(index, el) {
            moduleReg+= $(this).text()+"|";
        });
        // 截取最后一个字符
        moduleReg=moduleReg.substr(0,moduleReg.length-1);
        table.column(thisName).search(moduleReg, true, false, true).draw();
    }
});

// 搜索过滤 start
$('#host_query_submit').on('click', function(event) {
    $('.column_filter').trigger('search');
});
// 搜索过滤 end

// 点击 checkbox
$('.table-box').on('change', '.host-td-checkbox', function(e){
    var target = e.target;
    var _tr=$(target).closest('tr').get(0);
    var checkedNum;
    var table= $('#table_topo').DataTable();
    var data = $('#table_topo').DataTable().data();
    var d = data;
    if($(target).parent('th').length>0){//表头 全选
        var checked = $(target).prop('checked');
        for(var i=0,len=d.length; i<len; i++){
            d[i].Checked = checked ? 'checked' : '';
            $('.host-td-checkbox').eq(i+1).prop('checked', checked )
        }
        checkedNum = checked ? d.length : 0;
        $('.host-td-checkbox').not('.host-check-all').prop('checked', checked )
        console.log(checkedNum);

    }else{
        var checked = $(target).prop('checked');
        if (!checked) {
            $('.host-check-all').prop('checked', false);
        };
        var checkItem=$('.host-td-checkbox').not('.host-check-all');
        var len = d.length;
        var _index = table.row(_tr).index();
        var _this=$(this).closest('tr').get(0);
        var dataItem = table.row( _this ).data();

        dataItem.Checked = checked ? 'checked' : '';

        var checkedNum=0;
        for (var i = 0; i < len; i++) {
            d[i]["Checked"]=='checked'?checkedNum++:'';
        };
    }

    if(checkedNum>0){
        $(".c-grid-copyIP,.c-grid-batEdit,.c-grid-saveAsExcel,.c-grid-batDel,.c-grid-batRes").css("color","#555");
        $('#moveIp,#batDel,#batEdit,#batRes').attr('disabled', false);
        $('.table-box').find('.c-grid-copyIP').attr('disabled', false);
        $('.c-grid-moveIp').css('color', 'rgb(85, 85, 85)');
    }else{
        $(".c-grid-copyIP,.c-grid-batEdit,.c-grid-saveAsExcel,.c-grid-batDel,.c-grid-batRes").css("color","#a1a1a1");
        $('#moveIp,#batDel,#batEdit,#batRes').attr('disabled', true);
        $('.table-box').find('.c-grid-copyIP').attr('disabled', true);
        $('.c-grid-moveIp').css('color', '#a1a1a1');
    }

    if(checkedNum==d.length){
        $('.host-check-all').prop('checked', true);
        var selectAllDialog = dialog({
            align: 'bottom',
            content: '<div class="c-dialogdiv2"><i class="c-dialogimg-success"></i>全选<i class="redFont">'+checkedNum+'</i>台主机</div>'
        });
        selectAllDialog.show(document.getElementById('dialogs'));
        setTimeout(function(){selectAllDialog.close().remove();}, 2000);
    }

});
// 点击checkbox end
$(".c-grid-moveIp").click(function(e){
    var note = $(this).attr("disabled");
    e.stopPropagation();
    e.preventDefault();
    if(note !== "disabled"){
        $(".downList").toggle();
    }
});
$(".downList").click(function(e){
    e.stopPropagation();
    $(".downList").show();
})
$(document).click(function(){
    $(".downList").hide();
});

$('.searchPanel input').bind('input propertychange', function() {
    var val = $(this).val();
    $(".downList li").hide();
    $(".downList li").each(function(){
        if($(this).text().indexOf(val)>=0){
            $(".downList li").eq($(this).index()).show();
        }
    });
});

var allSelect = 0;
$(".downList .ui-checkbox").click(function(e){
    e.stopPropagation();
    $(this).is(":checked")?allSelect++:allSelect--;
    if(allSelect>0){
        $(".selectNote").text(allSelect+" selected");
    }else{
        $(".selectNote").text("转移主机至");
    }
});
$(".downList li").click(function(){
    $(this).find(".ui-checkbox").trigger("click");
});

$("#myScroll").mCustomScrollbar({
    setHeight:175, //设置高度
    theme:"minimal-dark" //设置风格
});

// 显示详情
$('.table-box').on('click', '.a-innerip', function(e){
    var me = e.target;
    var _tr = $(me).closest('tr').get(0);
    var grid = $('#table_topo').DataTable();
    var data = grid.row( _tr ).data();

    var param = {};
    param['HostID'] = data.HostID;
    param['ApplicationID'] = cookie.get('defaultAppId');
    $.ajax({
        url:'/host/details',
        data:param,
        dataType:'html',
        method:'post',
        success:function(data){
            CC.rightPanel.show();
            CC.rightPanel.render(data);

            $(".sidebar-panel").mCustomScrollbar({
                theme: "minimal-dark" //设置风格
            });

            $('.show-all-details-info').on('click',function(e){
                $('.show-all-details-info').popover('destroy');
                $(e.target).popover('show');
            });

            $('.sidebar-panel-container-new').on('click',function(e){
                var className = $(e.target).prop('class');
                if(className != 'show-all-details-info' && className != 'popover-content'){
                  $('.show-all-details-info').popover('destroy');
                }
            });
        }
    });
});
// 显示详情 end

/* 点击修改选中按钮 */
$('#batEdit').on('click', function (e){
    var hostId = [];
    var appId = cookie.get('defaultAppId');
    var table = $('#table_topo').DataTable();
    var data = table.data();
    var newData = data;
    var hostInfo = [];
    for(var i=0,len=newData.length; i<len; i++){
        if(newData[i].Checked==='checked'){
            hostId.push(newData[i].HostID);
            hostInfo = newData[i];
        }
    }

    if(hostId.length==0){
        var noHostSelectDialog = dialog({
                content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请选择主机</div>'
            });
        noHostSelectDialog.show();
        setTimeout(function () {
            noHostSelectDialog.close().remove();
        }, 2000);
        return false;
    }

    if(hostId.length==1){
        var _this=$(this).closest('tr').get(0);
        var hostInfo = table.row( _this ).data();
        $('#moduleHost_HostName').val(hostInfo['HostName']);
        $('#moduleHost_Description').val(hostInfo['Description']);
        $('#moduleHost_Source').val(hostInfo['Source']);

        $("#moduleHost_Operator").select2({
            placeholder:'请选择',
            data:window.userList,
            formatResult:function format(state) {
            return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
        },
        formatSelection:function format(state) {
            return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
        }
    }).select2('val',hostInfo['Operator']).select2("enable", false);

        $("#moduleHost_BakOperator").select2({
            placeholder:'请选择',
            data:window.userList,
            formatResult:function format(state) {
                        return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
                    },
            formatSelection:function format(state) {
                        return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
                    }
        }).select2('val',hostInfo['BakOperator']).select2("enable", false);

        $("#moduleHost_Source").select2({
            placeholder:'请选择',
            data:window.hostSource,
            formatResult:function format(state) {
                return state.text;
            },
            formatSelection:function format(state) {
                return state.text;
            }
        }).select2('val',hostInfo['Source']).select2("enable", false);
    }else{
        $("#moduleHost_Operator,#moduleHost_BakOperator").select2({
            placeholder:'请选择',
            data:window.userList,
            formatResult:function format(state) {
                        return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
                    },
            formatSelection:function format(state) {
                        return state.id ? state.id + "<span class='select-info'>(" + state.text + ")</span>" : state.text;
                    }
        }).select2('val','').select2("enable", false);

        $("#moduleHost_Source").select2({
            placeholder:'请选择',
            data:window.hostSource,
            formatResult:function format(state) {
                return state.text;
            },
            formatSelection:function format(state) {
                return state.text;
            }
        }).select2('val',hostInfo['Source']).select2("enable", false);
    }

    $('#batEdit_window').modal('show');
});

/* 点击修改选中弹窗的保存按钮*/
$('#batEdit_window').on('click', '#moduleHostSubmit', function(){
    var hostInfo = {};
    var stdProperty = {};
    var cusProperty = {};

    if(!$('#moduleHost_HostName').prop('disabled')){
        stdProperty['HostName'] = $('#moduleHost_HostName').val();
        if(stdProperty['HostName']==null || stdProperty['HostName']==''){
            var d = dialog({
                    content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>主机名不能为空</div>',
                    zIndex:1051
                });
            d.show();
            setTimeout(function () {
                d.close().remove();
                $('#moduleHost_HostName').focus();
            }, 1500);
            $('#moduleHost_HostName').focus();
            return false;
        }
    }

    var empty_field = [];
    var empty_field_disable = [];
    if(!$('#moduleHost_Operator').prop('disabled')){
        stdProperty['Operator'] = $('#moduleHost_Operator').select2('val');
        stdProperty['Operator']=='' && empty_field.push('负责人');
    }

    if(!$('#moduleHost_BakOperator').prop('disabled')){
        stdProperty['BakOperator'] = $('#moduleHost_BakOperator').select2('val');
        stdProperty['BakOperator']=='' && empty_field.push('备份负责人');
    }

    if(!$('#moduleHost_Source').prop('disabled')){
        stdProperty['Source'] = $('#moduleHost_Source').val();
        stdProperty['Source']=='' && empty_field_disable.push('云供应商');
    }

    if(!$('#moduleHost_Description').prop('disabled')){
        stdProperty['Description'] = $('#moduleHost_Description').val();
        stdProperty['Description']=='' && empty_field.push('备注信息');
    }

    var table= $('#table_topo').DataTable();
    var data = $('#table_topo').DataTable().data();

    var newData = data;
    var hostId = [];
    for(var i in newData){
        if(newData[i].Checked==='checked'){
            hostId.push(newData[i].HostID);
        }
    };
    hostInfo['HostID'] = hostId.join(',');
    hostInfo['ApplicationID'] = cookie.get('defaultAppId');
    hostInfo['stdProperty'] = stdProperty;
    hostInfo['cusProperty'] = cusProperty;

    if(empty_field_disable.length > 0){
        var promptDialog = dialog({
            title:'提示',
            width:300,
            okValue:"确定",
            zIndex:1051,
            ok:function(){},
            content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>' + empty_field_disable.join('、') + '不能设置为空</div>'
        });

        promptDialog.showModal();
        return true;
    }

    if(empty_field.length > 0){
        var confirmDialog = dialog({
            title:'确认',
            width:300,
            zIndex:1051,
            content: empty_field.join('、')+'将设置为空，确认继续？',
            okValue:"继续",
            cancelValue:"取消",
            ok:function (){
                if(!$.isEmptyObject(hostInfo['stdProperty']) && hostInfo['ApplicationID']!='' && hostInfo['HostID']!=''){
                    $.ajax({
                        url:'/host/updateHostInfo/',
                        data:hostInfo,
                        dataType:'json',
                        method:'post',
                        success:function(response){
                            var content = '<i class="c-dialogimg-'+ (response.success==true?'success':'prompt') +'"></i>'+response.message;
                            var d = dialog({
                                    content: '<div class="c-dialogdiv2">'+content+'</div>'
                                });
                            d.show();
                            setTimeout(function () {
                                d.close().remove();
                            }, 2500);
                            $('.host-check-all').prop('checked', false);
                            table.ajax.reload();
                            // CC.host.hostlist.init();
                            return true;
                        }
                    });

                    $("#s2id_moduleHost_Operator,#s2id_moduleHost_BakOperator").select2('val', '');
                    $('#batEdit_window').find('input[type=checkbox]').prop('checked', false).end().find('input[type=text]').val('').attr('disabled', true).end().modal('hide');
                }
            },
            cancel: function () {
            }
        });

        confirmDialog.showModal();
        return true;
    }

    if(!$.isEmptyObject(hostInfo['stdProperty']) && hostInfo['ApplicationID']!='' && hostInfo['HostID']!=''){
        $.ajax({
            url:'/host/updateHostInfo/',
            data:hostInfo,
            dataType:'json',
            method:'post',
            success:function(response){
                var content = '<i class="c-dialogimg-'+ (response.success==true?'success':'prompt') +'"></i>'+response.message;
                var d = dialog({
                        content: '<div class="c-dialogdiv2">'+content+'</div>'
                    });
                d.show();
                setTimeout(function () {
                    d.close().remove();
                }, 2500);
                $('.host-check-all').prop('checked', false);
                table.ajax.reload();
                // CC.host.hostlist.init();
                return true;
            }
        });
    }

    $("#s2id_moduleHost_Operator,#s2id_moduleHost_BakOperator").select2('val', '');
    $('#batEdit_window').find('input[type=checkbox]').prop('checked', false).end().find('input[type=text]').val('').attr('disabled', true).end().modal('hide');
});

// 没有选择主机时不能复制
$('.c-grid-copyIP').on('mouseenter', function() {
    if($(this).attr('disabled')){
        $('.copy-menu').addClass('hide');
    }else{
        $('.copy-menu').removeClass('hide');
    };
});

// 复制函数
function _init_copy(btnclass,filedname){
        /**
         * copy
         */
        var copyIpBtn = $('.c-grid-toolbar').find(btnclass);
        ZeroClipboard.config({moviePath:'/assets/ZeroClipboard/ZeroClipboard.swf'});
        var clip = new ZeroClipboard(copyIpBtn.get(0));
        clip.on('copy',function(e){
                 var clipboarde=e.clipboardData;
                 _copyIp(copyIpBtn,clipboarde,filedname);
        });
        clip.on('aftercopy',function(e){
            if(e.data['text/plain']){
                var d = dialog({
                        content: '<div class="c-dialogdiv2"><i class="c-dialogimg-success"></i>复制成功</div>'
                    });
                d.show();
                setTimeout(function() {
                    d.close().remove();
                }, 2500);
            }
        });
    };

    function _copyIp(btn,clipboard,filedname){
        var me=this,
            list=_getCopyList(filedname);
        if(list.length){
            clipboard.setData('text/plain',list.join("\n"));
        }else{
            if($('.c-grid-copyIP').attr('disabled')=='disabled'){
                return false;
            }
            var noHostSelectDialog = dialog({
                    content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请选择主机</div>'
                });
            noHostSelectDialog.show();
            setTimeout(function () {
                noHostSelectDialog.close().remove();
            }, 2000);
            return false;
        }
    };
    function _getCopyList(key){
        var list=[];
        var newData = $('#table_topo').DataTable().data();
        for(var i=0,len=newData.length; i<len; i++){
            if(newData[i].Checked==='checked'){
                list.push(newData[i][key]);
            }
        }
        return list;
    }
    _init_copy('.copy-inner-ip','InnerIP');
    _init_copy('.copy-outer-ip','OuterIP');
    _init_copy('.copy-asset-id','AssetID');

/*导出到excel*/
$('.table-box').on('click','.c-grid-saveAsExcel',function(e){
    if($(e.target).attr('disabled')=='disabled'){
        return false;
    }

    if($('#hostExport').length>0){
        $('#hostExport').remove();
    }

    var grid = $('#table_topo').DataTable();
    var d = $('#table_topo').DataTable().data();
    var hostId = [];
    for(var i=0,len=d.length; i<len; i++){
        hostId.push(d[i].HostID);
    }

    if(hostId.length==0){
        var d = dialog({
                content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>没有主机可导出</div>'
            });
        d.show();
        setTimeout(function() {
            d.close().remove();
        }, 2500);
        return false;
    }
    $('body').append('<form id="hostExport" action="/host/hostExport" method="post" style="display:none;" target="_self"><input type="text" name="HostID" value="'+hostId.join(',')+'"><input type="hidden" name="ApplicationID" value="'+cookie.get('defaultAppId')+'"></form>');
    var dialog1 = dialog({
        content: '<div class="c-dialogdiv2"><img class="c-dialogimg-loading" src="/static/img/loading_2_24x24.gif"></img>导出中...</div>'
    });
    setTimeout(function(){
        dialog1.showModal();
        window.dintval = setInterval(function(){
            var dowcomplete = cookie.get('comdownload');
            if(dowcomplete == 1){
                dialog1.close().remove();
                cookie.set('comdownload','');
                clearInterval(dintval);
            }
        }, 500);
        $('#hostExport').submit();
    },500);

});

/* 点击移至空闲机按钮 */
$('.table-box').on('click', '#batDel', function(e){
    if($(e.target).attr('disabled')=='disabled'){
        return false;
    }

    var grid = $('#table_topo').DataTable();
    var newData = $('#table_topo').DataTable().data();
    var hostId = [];
    for(var i=0,len=newData.length; i<len; i++){
        if(newData[i].Checked=='checked'){
            hostId.push(newData[i].HostID);
        }
    }

    if(hostId.length==0){
        var noHostSelectDialog = dialog({
                content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请选择主机</div>'
            });
        noHostSelectDialog.show();
        setTimeout(function () {
            noHostSelectDialog.close().remove();
        }, 2000);
        return false;
    }

    var param = {};
    param['ApplicationID'] = cookie.get('defaultAppId');
    param['HostID'] = hostId.join(',');

    var gridBatDel = dialog({
        title:'确认',
        width:300,
        content: '确认是否将已勾选的<i class="redFont">'+hostId.length+'</i>台主机移动至空闲机?',
        okValue:"确定",
        cancelValue:"取消",
        ok:function (){
            var d = dialog({
                content: '<div class="c-dialogdiv2"><img class="c-dialogimg-loading" src="/static/img/loading_2_24x24.gif"></img>转移中...</div>'
            });
            d.showModal();
            $.ajax({
                url:'/host/delHostModule/',
                dialog:d,
                data:param,
                method:'post',
                dataType:'json',
                //async:false,
                success:function(response){

                    $.ajax({
                        url:'/host/getTopoTree4view',
                        // url:'/static/js/j3.json',
                        data:{ApplicationID:cookie.get('defaultAppId')},
                        dataType:'json',
                        method:'post',
                        success:function(response){
                            $('#treeContainer').jstree('destroy');
                            $('#emptyContainer').jstree('destroy');
                            createTopo();
                            grid.ajax.reload();
                        }
                    });

                    this.dialog.close().remove();
                    var content = '<i class="c-dialogimg-'+ (response.success==true ? 'success' : 'prompt') +'"></i>'+ response.message +'</div>';
                    var d = dialog({
                        content: '<div class="c-dialogdiv2">'+content+'</div>'
                    });
                    d.show();
                    setTimeout(function() {
                        d.close().remove();
                    }, 2500);
                    // CC.host.hostlist.init();
                    return true;
                }
            });
        },
        cancel: function () {
        }
    });

    gridBatDel.showModal();
});

//转移主机操作
$(document).on('change', '.ui-multiselect-checkboxes', function(e){
    var disabled = true;
    $(e.target).closest('ul').find('input[type=checkbox]').each(function(index, el){
        if($(el).prop('checked')){
            disabled = false;
            return false;
        }
    });

    $('.btn-primary', '#modSelectMenu').prop('disabled', disabled);
});

$(document).on('mouseleave', '.ui-multiselect-checkboxes', function(e){
    $(this).find('.ui-state-hover').removeClass('ui-state-hover');
});
// 转移主机按钮
$('.operationSelect').click(function(e){
    e.stopPropagation();
    if($(this).hasClass('btn-default')){
        $(".downList").hide();
    }else{
        if($(this).prop('disabled')!==false){
            return false;
        }
        var appId = cookie.get('defaultAppId');
        var moduleId = [];
        var hostId = [];
        $('#downList li').each(function(index){
            var mID = $(this).attr("value");
            if($(this).find("input[type=checkbox]").is(':checked')){
                moduleId.push(mID);
            }
        });
        // $(e.target).parents('.ui-multiselect-menu').find('.ui-icon-circle-close').click();
        $(".downList").hide();

        if(moduleId.length==0){
            var content = '';
            var d = dialog({
                    content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请至少选择一个模块</div>'
                });
            d.show();
            setTimeout(function () {
                d.close().remove();
            }, 2500);

            return false;
        }

        var grid = $('#table_topo').DataTable();
        var newData = $('#table_topo').DataTable().data();

        var d = newData;
        for(var i in d){
            if(d[i].Checked==='checked'){
                hostId.push(d[i]['HostID']);
            }
        }

        if(hostId.length==0){
            var content = '';
            var d = dialog({
                    content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请至少选择一台主机</div>'
                });
            d.show();
            setTimeout(function () {
                d.close().remove();
            }, 2500);

            return false;
        }

        var param = {};
        param['ApplicationID'] = appId;
        param['ModuleID'] = moduleId.join(',');
        param['HostID'] = hostId.join(',');

        var d = dialog({
            content: '<div class="c-dialogdiv2"><img class="c-dialogimg-loading" src="/static/img/loading_2_24x24.gif"></img>转移中...</div>'
        });
        d.showModal();

        setTimeout(function(){
                $.ajax({
                url:'/host/modHostModule/',
                dialog:d,
                data:param,
                dataType:'json',
                method:'post',
                success:function(response){
                    this.dialog.close().remove();
                    var content = '<i class="c-dialogimg-'+ (response.success==true?'success':'prompt') +'"></i>'+response.message;
                    var d = dialog({
                            content: '<div class="c-dialogdiv2">'+content+'</div>'
                        });
                    d.show();
                    setTimeout(function () {
                        d.close().remove();
                    }, 2500);
                    // CC.host.hostlist.init();
                    // CC.host.topology.refresh();
                    $('#treeContainer').jstree('destroy');
                    $('#emptyContainer').jstree('destroy');
                    createTopo();
                    grid.ajax.reload();
                    return true;
                }
            });
        },100);
        $(".c-grid-copyIP,.c-grid-batEdit,.c-grid-saveAsExcel,.c-grid-batDel,.c-grid-batRes").css("color","#a1a1a1");
        $('#moveIp,#batDel,#batEdit,#batRes').attr('disabled', true);
        $('.table-box').find('.c-grid-copyIP').attr('disabled', true);
        $('.c-grid-moveIp').css('color', '#a1a1a1');
    }
});

/* 点击上交按钮 */
$('.table-box').on('click', '#batRes', function (e){
    if($(e.target).attr('disabled')=='disabled'){
        return false;
    }

    var grid = $('#table_topo').DataTable();
    var newData = $('#table_topo').DataTable().data();
    var hostId = [];
    for(var i=0,len=newData.length; i<len; i++){
        if(newData[i].Checked==='checked'){
            hostId.push(newData[i].HostID);
        }
    }

    if(hostId.length==0){
        var noHostSelectDialog = dialog({
                content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>请选择主机</div>'
            });
        noHostSelectDialog.show();
        setTimeout(function () {
            noHostSelectDialog.close().remove();
        }, 2000);
        return false;
    }


    var param = {};
    param['ApplicationID'] = cookie.get('defaultAppId');
    param['HostID'] = hostId.join(',');

    var gridBatRes = dialog({
        title:'确认',
        width:300,
        content: '确认是否将已勾选的<i class="redFont">'+hostId.length+'</i>台主机上交至资源池',
        okValue:"确定",
        cancelValue:"取消",
        ok:function (){
            var d = dialog({
                content: '<div class="c-dialogdiv2"><img class="c-dialogimg-loading" src="/static/img/loading_2_24x24.gif"></img>上交中...</div>'
            });
            d.showModal();
            $.ajax({
                url:'/host/resHostModule/',
                dialog:d,
                data:param,
                method:'post',
                dataType:'json',
                //async:false,
                success:function(response){
                    this.dialog.close().remove();
                    var content = '<i class="c-dialogimg-'+ (response.success==true ? 'success' : 'prompt') +'"></i>'+ response.message +'</div>';
                    var d = dialog({
                        content: '<div class="c-dialogdiv2">'+content+'</div>'
                    });
                    d.show();
                    setTimeout(function() {
                        d.close().remove();
                    }, 2500);

                    $('#treeContainer').jstree('destroy');
                    $('#emptyContainer').jstree('destroy');
                    createTopo();
                    reDrawTable();

                    return true;
                }
            });
        },
        cancel: function () {
        }
    });

    gridBatRes.showModal();
});


function clearTable(){
    // 把旧的表格删除，重新生成一个空表格
    var table = $('#table_topo').DataTable();
    table.destroy();
    $('#table_topo').remove();
    $('.table-box').append('<table id="table_topo" class="table table-bordered table-striped table-responsive">');
}

var curParam = null;
function drawTable(data){
    // 节点的数据
    var newTableData=tableData;
    // 判断类型，传参
    switch(data.node.original.type)
    {
    case "application":
        //2级拓扑
        if (data.node.original.appId==2) {
            var curParam={ ApplicationID:data.node.original.appId,ModuleID:data.node.children.join(',') };
        }else{//3级拓扑
            var curParam={ ApplicationID:data.node.original.appId,ModuleID:data.node.children_d.join(',') };
        }
      break;
    case "set":
      curParam={ ApplicationID:data.node.original.appId,SetID:data.node.id };
      break;
    case "module":
      curParam = { ApplicationID:data.node.original.appId,ModuleID:data.node.id }
      break;
    case "all":
      curParam = { ApplicationID:data.node.original.appId }
    }
    // 重新绘制表格
    newTableData['ajax'].data= curParam ;
    $('#table_topo').DataTable(newTableData);
}

function getCurParam(data){
    // 节点的数据
    var newTableData=tableData;
    // 判断类型，传参
    switch(data.node.original.type)
    {
    case "application":
        //2级拓扑
        if (data.node.original.appId==2) {
            var curParam={ ApplicationID:data.node.original.appId,ModuleID:data.node.children.join(',') };
        }else{//3级拓扑
            var curParam={ ApplicationID:data.node.original.appId,ModuleID:data.node.children_d.join(',') };
        }
      break;
    case "set":
      var curParam={ ApplicationID:data.node.original.appId,SetID:data.node.id };
      break;
    case "module":
      var curParam = { ApplicationID:data.node.original.appId,ModuleID:data.node.id }
      break;
    case "all":
      var curParam = { ApplicationID:data.node.original.appId }
    }
    return curParam;
}

var globalParam = {};
function reDrawTable(params){
    if (params){
        globalParam = params;
    }

    var newTableData=tableData;
    newTableData['ajax'].data= globalParam;
    $('#table_topo').DataTable(newTableData);
}
function createTopo(callback){
    $.ajax({
        url: '/host/getTopoTree4view',
        // url:'/static/js/j3.json',
        type: 'POST',
        dataType: 'json',
        success:function (result){
            var topo=result.topo;
            var empty=result.empty;
            if(topo[0].children.length){
                for(var i=0,j=topo[0].children.length;i<j;i++){
                    var newId = topo[0].children[i].id+"plus";
                    topo[0].children[i].id = newId;
                }
            }


            //2级拓扑
            if (topo[0]['lvl']==2) {
                var modules = topo[0]['children'];
                var module_id=[];
                for (var i = 0; i < modules.length; i++) {
                    module_id.push(modules[i]['id']);
                };

            }else{//3级拓扑
                var modules = topo[0]['children'];
                var module_id=[];
                for (var i = 0; i < modules.length; i++) {
                    for (var j = 0; j < modules[i]['children'].length; j++) {
                        module_id.push(modules[i]['children'][j]['id']);
                    };
                };

            }

            var application_id = topo[0]['id']
            module_id = module_id.join(',');
            curParam={ ApplicationID:application_id,ModuleID:module_id };
            callback && callback(curParam);

            // 拓扑树
            $('#treeContainer').jstree({
                'core' : {
                    'data' : topo
                }
            }).on('loaded.jstree',function (e,data){
                // console.log(data);
            })
              .on('changed.jstree',function (e,data){
                // console.log(data)
                 // 显示上交按钮
                $('#batDel').show();
                $('#batRes').hide();
                // 重新渲染表格
                clearTable();
                drawTable(data);
            });
            // 空闲机
            $('#emptyContainer').jstree({
                'core' : {
                    'data' : empty
                }
            }).on('changed.jstree',function (e,data){
                 // 显示上交按钮
                $('#batDel').hide();
                $('#batRes').show();
                // 重新渲染表格
                clearTable();
                var params = getCurParam(data);
                reDrawTable(params);
            });
        }
    });
}
createTopo(function(curParam){
    var newTableData=tableData;
    newTableData['ajax'].data= curParam ;
    $('#table_topo').DataTable(newTableData);
});
/**
* 查询条件按钮"空闲机"点击事件处理函数
*/
$('#filter_module_empty').click(function(e){
    e.stopPropagation();
    $('#emptyContainer .jstree-anchor').trigger('click');
});
/**
* 查询条件按钮"ALL"点击事件处理函数
*/
$('#filter_module_all').click(function(e){
    e.stopPropagation();
    clearTable();
    // 节点的数据
    var newTableData=tableData;
    // 重新绘制表格
    newTableData['ajax'].url = '/host/getHostByCondition' ;
    newTableData['ajax'].data = {'ApplicationID' : cookie.get('defaultAppId')};
    $('#table_topo').DataTable(newTableData);
    // CC.host.hostlist.init();
    $('#batRes').hide();
    $('#batDel').show();
});
/**
* 查询条件按钮"我"点击事件处理函数
*/
$('#filter_module_mine').click(function(e){
    e.stopPropagation();
    clearTable();
    // 节点的数据
    var newTableData=tableData;
    // 重新绘制表格
    newTableData['ajax'].url = '/host/getHostByCondition' ;
    newTableData['ajax'].data = {'ApplicationID' : cookie.get('defaultAppId'), 'Operator' : window.currUser.username};
    $('#table_topo').DataTable(newTableData);
    // CC.host.hostlist.init();
    $('#batRes').hide();
    $('#batDel').show();

});


// 拓扑树根据浏览器调整自身高度
function treeHeightChange () {
    $('.host-sidebar-left').css('position','fixed');
    setTimeout(function (){
        $('.c-host-side').css('height',$(window).outerHeight()-70-20);
        if ($('.c-host-side').height()>820)$('.c-host-side').css('height',820);
        $(".c-tree-box").css('height',$('.c-host-side').outerHeight()-$('.free-group').outerHeight()-$('.c-host-side>h4').outerHeight()-40);
    },200)

}
treeHeightChange();
$(window).resize(function (){
    treeHeightChange();
})
// 拓扑树根据浏览器调整自身高度 end

})
</script>
