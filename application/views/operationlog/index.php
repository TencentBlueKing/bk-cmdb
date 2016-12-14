<?php
/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */
?>
<style rel="stylesheet">
.select2-container .select2-choice{height: 34px;}
.operation_log_table tbody{}
</style>
<link href="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/dataTables.bootstrap.css" rel="stylesheet"/>
<link href="<?php echo STATIC_URL;?>/static/css/bootstrap-datetimepicker.min.css" rel="stylesheet">
<!-- 主面板Content -->
<div class="content-wrapper">
    <!-- 主面板 Header  -->
    <section class="content-header">
        <h1>操作日志</h1>
    </section>
    <!-- 主面板 main-sidebarn-->
    <section class="content p20">
        <div class="option-log">
            <div class="c-import-table">
                <div class="c-search-box">
                    <div class="panel-heading" id="collapseOneBtn">
                        <h4 class="panel-title" style="    display: inline-block;">查询条件 </h4>
                    </div>
                    <div id="collapseOne" class="panel-collapse collapse in">
                        <div class="search-content">
                        <div class="col-lg-4">
                            <div class="input-group pb10">
                                <label for="inputCount3" class="input-group-addon">业务：</label>
                                <input type="hidden" class="bigdrop form-control form-control" id="ApplicationID" style="width:100%;height:34px;">
                                </div>
                            </div>
                            <div class="col-lg-4">
                                <div class="input-group pb10">
                                    <label for="Operator" class="input-group-addon">操作人：</label>
                                    <input type="text" placeholder="" value="" id="Operator" name="Operator" class="form-control">
                                </div>
                            </div>
                            <div class="col-lg-4">
                                <div class="input-group pb10">
                                    <label for="OpTarget" class="input-group-addon">操作目标：</label>
                                    <input type="text" placeholder="" value="" id="OpTarget" name="OpTarget" class="form-control">
                                </div>
                            </div>
                            <div class="col-lg-4">
                                <div class="input-group pb10">
                                    <label for="OpContent" class="input-group-addon">操作内容：</label>
                                    <input type="text" placeholder="" value="" id="OpContent" name="OpContent" class="form-control">
                                </div>
                            </div>
                            <!-- 时间范围 -->
                            <div class="col-lg-4">
                                <div class="input-group pb10">
                                    <label for="start" class="input-group-addon">开始日期：</label>
                                    <!-- <div class="p0 form-control">
                                        <input type="text" id="start" placeholder="格式 2016-01-01" style="width:100%;">
                                    </div> -->
                                    <div class="form-control">
                                        <div class="input-group date" id="start_date" data-date="" data-date-format="dd MM yyyy" data-link-field="dtp_input2" data-link-format="yyyy-mm-dd">
                                            <span class="input-group-addon" style="position: absolute;left: 0;top: 0;"><span class="glyphicon-calendar"></span></span>
                                            <input class="form-control" id="start" size="16" type="text" placeholder="格式 2016-01-01">
                                            <input type="hidden" id="dtp_input2" value="">
                                            <span class="input-group-addon"><span class="glyphicon glyphicon-calendar"></span></span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div class="col-lg-4">
                                <div class="input-group pb10">
                                    <label for="end" class="input-group-addon">结束日期：</label>
                                    <div class="form-control">
                                        <div class="input-group date" id="end_date" data-date="" data-date-format="dd MM yyyy" data-link-field="dtp_input2" data-link-format="yyyy-mm-dd">
                                            <span class="input-group-addon" style="position: absolute;left: 0;top: 0;"><span class="glyphicon-calendar"></span></span>
                                            <input class="form-control" id="end" size="16" type="text" placeholder="格式 2016-01-01">
                                            <input type="hidden" id="dtp_input2" value="">
                                            <span class="input-group-addon"><span class="glyphicon glyphicon-calendar"></span></span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <!-- 时间范围 end -->
                            <div class="col-xs-12 text-center">
                                <button class="btn btn-default" id="log_query_reset">重置</button>
                                <button class="btn btn-primary" id="log_query_submit">查询</button>
                            </div>
                            <div class="clearfix"></div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <!-- <div class="c-import-search">
            <input id="filter" type="text" class="form-control pull-left w200" placeholder="搜索..." /><i class="glyphicon glyphicon-search"></i>
        </div> -->
        <table id="operation_log_table" class="table table-bordered table-striped"></table>
        <!-- <div id="operation_log_table" class="operation-log-table"></div> -->
        <div class="control-sidebar-bg"></div>
    </section>
</div>

<script src="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/jquery.dataTables.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/datatables-1.10.7/dataTables.bootstrap.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/mCustomScrollbar-3.0.9/jquery.mCustomScrollbar.concat.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/ZeroClipboard/ZeroClipboard.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/bootstrap-switch-master/dist/js/bootstrap-switch.js" rel="stylesheet"></script>
<script src="<?php echo STATIC_URL;?>/static/js/app.min.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/select2-3.5.2/select2.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/datetimepicker/bootstrap-datetimepicker.js"></script>
<script src="<?php echo STATIC_URL;?>/static/assets/datetimepicker/bootstrap-datetimepicker.zh-TW.js" charset="UTF-8"></script>
<script type="text/javascript">
$(function () {
    var params = {};
    query(params);
    var app = <?php echo $app;?>;
    $("#ApplicationID").select2({ data: app });
    $("#ApplicationID").select2('val',<?php echo $ApplicationID;?>);
    /**
    * 查询条件“重置”按钮点击事件处理函数
    */
    $('#log_query_reset').click(function(){
        $("#ApplicationID").select2('val','');
        $('#Operator').val('');
        $('#OpTarget').val('');
        $('#OpContent').val('');
        $('#start').val('');
        $('#end').val('');
    });
    /**
    * 查询条件“查询”按钮点击事件处理函数
    */
    $('#log_query_submit').click(function(){
        var param = {};
        var ApplicationID = $('#ApplicationID').val();
        if(!ApplicationID){
            var promptDialog = dialog({
                title:'提示',
                width:300,
                okValue:"确定",
                zIndex:1051,
                ok:function(){},
                content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>业务不能为空</div>'
            });
            promptDialog.showModal();
            return true;
        }
        param['ApplicationID'] = ApplicationID;
        var Operator = $('#Operator').val();
        if(Operator){
            param['Operator'] = Operator;
        }
        var OpTarget = $('#OpTarget').val();
        if(OpTarget){
            param['OpTarget'] = OpTarget;
        }
        var OpContent = $('#OpContent').val();
        if(OpContent){
            param['OpContent'] = OpContent;
        }
        var start = $('#start').val();
        if(start){
            param['start'] = start;
        }
        var end = $('#end').val();
        if(end){
            param['end'] = end;
        }
        query(param);
    });

    function query(params){
         //console.log(params);
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
        $("#operation_log_table").dataTable({
            paging: true, //隐藏分页
            ordering: false, //关闭排序
            info: false, //隐藏左下角分页信息
            //searching: false, //关闭搜索
            lengthChange: false, //不允许用户改变表格每页显示的记录数
            language: language, //汉化
            autoWidth: false,
            destroy: true,
            "ajax": {
                url: "/operationLog/getOperationLog",
                type:"post",
                data: function ( d ) {
                   return d = $.extend({},d, params);
                }
            },
            columns : [
              {title : '姓名', data:"ID",visible:false,width:"120px"},
              {title : '操作者', data:"Operator",width: "120px"},
              {title : '内容', data:"OpContent"},
              {title : '对象', data:"OpTarget",width:"50px"},
              {title : '类型', data:"OpName",width: "100px"},
              {title : '耗时(ms)', data:"ExecTime",width: "100px"},
              {title : '时间', data:"OpTime",width: "200px"}
            ]
        });
        // $("#operation_log_table").kendoGrid({
        //     dataSource: {//数据源配置项
        //         transport: {
        //             read: {
        //                 url: "/operationLog/getOperationLog",
        //                 data:params,
        //                 type:"post",
        //                 dataType: "json"
        //             }
        //         },
        //         pageSize:20,
        //         schema: {
        //             model: { id: "ID",
        //                     fields: {
        //                             "ID": {},
        //                             "Operator": { type: "string"},
        //                             "OpContent": { type: "string"},
        //                             "OpTarget": { type: "string"},
        //                             "OpName": { type: "string"},
        //                             "ExecTime": { type: "string"},
        //                             "OpTime": { type: "datetime"}
        //             } }
        //         }
        //     },
        //     pageable: true,
        //     resizable:true,
        //     scrollable: true,
        //     selectable:"multiple cell",
        //     allowCopy:{delimiter : ';'},
        //     height: 600,
        //     columns: [
        //         {field:'ID',title:"#",inputType:'text',width:"50px", hidden:true },
        //         {field: "Operator", title: "操作者",inputType:'text', width: "120px", filterable: true },
        //         //依赖强大的Template自定义显示内容,最牛逼的是可以模板可以是函数,最后事件绑定自己想办法,与extjs不一样，field不能重复，重复了结果很严重，可以自己试试。
        //         {field:"OpContent",title:'内容',inputType:'text',filterable: true},
        //         {field:"OpTarget",title:'对象',inputType:'text',width:"50px", filterable: true},
        //         {field:"OpName",title:'类型', inputType:'text',width: "100px",filterable: true},
        //         {field:"ExecTime",title:'耗时(ms)', inputType:'text',width: "100px",filterable: false},
        //         {field:"OpTime",title:'时间', inputType:'text',width: "200px",filterable: true}
        //     ],
        //     messages: {
        //         noRows: "没有记录",
        //         loading: "正在加载...",
        //         requestFailed: "加载失败...",
        //         retry: "重新加载"
        //     },
        //     editable: "inline"
        // });
    }

    //选择器的值改变时，更改两个选择器的最大值和最小值
    function startChange() {
        var startDate = start.value(),
        endDate = end.value();

        if (startDate) {
            startDate = new Date(startDate);
            startDate.setDate(startDate.getDate());
            end.min(startDate);
        } else if (endDate) {
            start.max(new Date(endDate));
        } else {
            endDate = new Date();
            start.max(endDate);
            end.min(endDate);
        }
    }

    function endChange() {
        var endDate = end.value(),
        startDate = start.value();

        if (endDate) {
            endDate = new Date(endDate);
            endDate.setDate(endDate.getDate());
            start.max(endDate);
        } else if (startDate) {
            end.min(new Date(startDate));
        } else {
            endDate = new Date();
            start.max(endDate);
            end.min(endDate);
        }
    }

    // var start = $("#start").kendoDatePicker({
    //     change: startChange,//事件绑定
    //     format : "yyyy-MM-dd"
    // }).data("kendoDatePicker");

    // var end = $("#end").kendoDatePicker({
    //     change: endChange,//事件绑定
    //     format : "yyyy-MM-dd"
    // }).data("kendoDatePicker");

    // start.max(end.value());
    // end.min(start.value());

});
 $('#start_date,#end_date').datetimepicker({
        language:  'zh-TW',
        weekStart: 1,
        todayBtn:  true,
        autoclose: true,
        todayHighlight: true,
        startView: 2,
        minView: 2,
        forceParse: false,
        format:"yyyy-mm-dd",
        pickerPosition: "bottom-right"
    });

</script>
