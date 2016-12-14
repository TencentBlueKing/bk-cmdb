/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

/*  拓扑配置 */
!(function(window){
    window.CC.hostConf=CC.hostConf||{};
    $('#editSetCapacity,#newSetCapacity').on('keyup',function(){
        var val = $(this).val().trim();
        if(val.substr(-1) == '.' && val.split('.').length >2){
            $(this).val(val.substring(0,val.length-1));
        }else{
            $(this).val(val.replace(/[^0-9|.]/gi,''));
        }
    });

    //首次加载的数据
    window.CC.hostConf.init = function() {
        if (level == 2) {
            var noset = true;
            var spanTxt="<span class='creat-module-btn btn btn-success btn-xs' cid='3'><i class='fa fa-plus'></i> 模块</span>"
        }
        else {
            var noset = false;
            var spanTxt="<span class='creat-group-btn btn btn-success btn-xs' cid='2'><i class='fa fa-plus'></i> 集群</span>"
        }

        $('#confTreeContainer').jstree({
            'core' : {
                 "themes" : {
                    "variant" : "large"//设置节点间距离
                 },
                'animation':false,//取消树列表的显示隐藏运动效果
                'data' : [{
                    id: appId,
                    text: "<span class='node-text'>"+appName+"</span>"+spanTxt,
                    noset: noset,
                    spriteCssClass: 'c-icon icon-app application',
                    icon: 'c-icon icon-app application',
                    type: 'application',
                    number: 200,
                    "state" : { "opened" : true },
                    children: topo
                }]
            }
        }).on("ready.jstree" , function (e,data){
            // 设置选中状态
            function selectStyle(selector){
                $('[jstree-active]').removeAttr('jstree-active');
                selector.attr('jstree-active',true);
            }
            // 新增集群事件
            $(".creat-group-btn").on("click", function(e) {
                ApplicationID = cookie.get('defaultAppId');
                $('.creat-container').css('display','none');
                $('.creat-group-container').fadeIn();
                e.preventDefault();
                e.stopPropagation();
                $("#newsetname").focus();

                var thisParent = $(this).parent('a')
                selectStyle(thisParent);

                $.post("/Set/getAllSetInfo",
                    {ApplicationID:ApplicationID}
                    ,function(result) {
                        rere = $.parseJSON(result);
                        if(rere.success == false) {
                            CC.hostConf.showWindow(rere.errInfo,'notice');
                            return;
                        }
                    });
                return false;
            });
            // 新建集群事件 end

            // 三级操作
            if(level==3) {
                $("#confTreeContainer>ul>li>a").attr('no-hovered', true);
                // 修改属性事件
                $('#confTreeContainer').on('click', '.node-text' , function() {
                    var $that = $(this);
                    var node = data.instance.get_node($that.closest('li'));
                    var nodeId = node.id+"_anchor";

                    if ($that.next('.creat-group-btn').length>0) {
                        return false;
                    };

                    if ($that.siblings('.icon-modal').length>0) {
                        selectStyle($that);
                        var panodeId = node.parent;// get parent id
                        var panode = data.instance.get_node( panodeId );// obj
                        console.log(panode)

                            if(typeof panode == 'undefined') {
                                return;
                            }
                            //二级结构，特殊处理
                            var SetID = panode.id;
                            var SetName = $(panode.text)[0].innerHTML;

                            $("#editmodulebe").text("所属集群");

                            $.post("/App/getMaintainers",
                                {}
                                ,function(result) {
                                rere = $.parseJSON(result);
                                var opstr='';
                                $.each( rere, function( i, x ) {
                                    var uin=i;
                                    opstr += ' <option value="' + uin + '">' + uin + '(' + x + ')' + '</option>';
                                });
                                $("#editOperator").html(opstr);
                                $("#editBakOperator").html(opstr);
                                $("#editOperator").val(node.original.operator);
                                $("#editBakOperator").val(node.original.bakoperator);
                                $("#editOperator").select2();
                                $("#editBakOperator").select2();
                            });

                            var ModuleID = node.id;
                            console.log($(node.text))
                            var ModuleName = $(node.text)[0].innerHTML;
                            $("#editmodulegroup").text(SetName);
                            $("#editmoduleModuleName").val(ModuleName);
                            $("#edit_module_property").html(ModuleName);
                            $("#editmoduleModuleName").attr('SetID',SetID);
                            $("#editmoduleModuleName").attr('ModuleID',ModuleID);
                            $('.creat-container').css('display','none');
                            $('.edit-module-container').fadeIn();
                            $("#editmoduleModuleName").focus();
                            return false;
                            // 修改模块属性 end
                    };

                    if ($that.siblings('.icon-group').length>0) {

                    selectStyle($that);
                    if(defaultapp==1) {
                        showWindows('资源池业务不能操作', 'notice');
                        return;
                    }

                    var ApplicationID = cookie.get('defaultAppId');
                    var SetID = node.id;
                    $("#editSetSetName").attr('setid',SetID);
                    $.post("/Set/getSetInfoById",
                        {ApplicationID:ApplicationID,SetID:SetID}
                        ,function(result) {
                            rere = $.parseJSON(result);
                            if (rere.success == false) {
                                showWindows(rere.errInfo, 'notice');
                            }
                            else {
                                var set = rere.set;
                                $("#editset_property").html(set.SetName);
                                $("#editSetSetName").val(set.SetName);
                                $("#editSetEnviType").val(set.EnviType);
                                $("#edit_set_setenctype label").removeClass('active');
                                if(set.EnviType == 1) {
                                    $("#edit_option1").parent().addClass('active');
                                }else if(set.EnviType == 2){
                                    $("#edit_option2").parent().addClass('active');
                                }
                                else{
                                    $("#edit_option3").parent().addClass('active');
                                }
                                $("#edit_set_sersta input").bootstrapSwitch({onText: '开放',
                                    offText: '关闭'});
                                $('#edit_set_sersta input').bootstrapSwitch('state', set.ServiceStatus == 1);
                                $("#editSetChnName").val(set.ChnName);
                                if(set.Capacity != 0) {
                                    $("#editSetCapacity").val(set.Capacity);
                                }
                                $("#editSetDes").val(set.Description);
                                $("#editOpenstatus").val(set.Openstatus);
                            }
                        });
                    $('.creat-container').css('display','none');
                    $('.edit-group-container').fadeIn();
                    $("#editSetSetName").focus();
                    e.preventDefault();
                    e.stopPropagation();
                    return false;
                    }

                });
                //修改属性事件 end

                // 创建模块 start
                $('#confTreeContainer').on('click', '.creat-module-btn' , function() {
                    function creatModule(node){
                        if(defaultapp==1) {
                            showWindows('资源池业务不能操作', 'notice');
                            return;
                        }
                        $.post("/App/getMaintainers",
                            {}
                            ,function(result) {
                                rere = $.parseJSON(result);
                                var opstr='';
                                $.each( rere, function( i, x ) {
                                    var uin=i;
                                    opstr += ' <option value="' + uin + '">' + uin + '(' + x + ')' + '</option>';
                                });
                                $("#Operator").html(opstr);
                                $("#BakOperator").html(opstr);
                                $("#Operator").select2();
                                $("#BakOperator").select2();

                            });

                        //对于二级业务来说重新取其setid
                        if(level==2) {
                            var SetID = desetid;
                            var SetName = cookie.get('defaultAppName');
                            $("#newmodulebe").text("所属业务");
                        } else {
                            var SetID = node.id;
                            var SetName = $(node.text)[0].innerHTML;
                            console.log(node)
                            $("#newmodulebe").text("所属集群");
                        }

                        $("#newmodulegroupname").text(SetName);
                        $("#newmoduleModuleName").attr('setid',SetID);
                        $('.creat-container').css('display','none');
                        $('.creat-module-container').fadeIn();
                        $("#newmoduleModuleName").focus();
                        e.preventDefault();
                        e.stopPropagation();
                        return false;
                    }

                    var $that = $(this);
                    var thisParent = $(this).parent('a');
                    var node = data.instance.get_node($that.closest('li'));
                    selectStyle(thisParent);
                    creatModule(node);
                    e.preventDefault();
                    e.stopPropagation();
                    return false;

                })
                // 创建模块 end
            }
            // 三级操作 end
        }).on("changed.jstree" , function(e, data) {
                function creatModule(node){
                    if(defaultapp==1) {
                        showWindows('资源池业务不能操作', 'notice');
                        return;
                    }
                    $.post("/App/getMaintainers",
                        {}
                        ,function(result) {
                            rere = $.parseJSON(result);
                            var opstr='';
                            $.each( rere, function( i, x ) {
                                var uin=i;
                                opstr += ' <option value="' + uin + '">' + uin + '(' + x + ')' + '</option>';
                            });
                            $("#Operator").html(opstr);
                            $("#BakOperator").html(opstr);
                            $("#Operator").select2();
                            $("#BakOperator").select2();

                        });

                    //对于二级业务来说重新取其setid
                    if(level==2) {
                        var SetID = desetid;
                        var SetName = cookie.get('defaultAppName');
                        $("#newmodulebe").text("所属业务");
                    }

                    $("#newmodulegroupname").text(SetName);
                    $("#newmoduleModuleName").attr('setid',SetID);
                    $('.creat-container').css('display','none');
                    $('.creat-module-container').fadeIn();
                    $("#newmoduleModuleName").focus();
                    e.preventDefault();
                    e.stopPropagation();
                    return false;
                }
                // 点击二级节点的操作
                if(level==2) {
                    // 判断节点有添加模块的按钮
                    if($(this).find('.jstree-clicked').find('.creat-module-btn').length>0){
                        $('[jstree-active]').removeAttr('jstree-active');
                        $('.jstree-clicked').attr('jstree-active',true);
                        //获取当前节点
                        var node = data.instance.get_node(data.selected[0]);
                        // 创建节点
                        creatModule(node);
                    }else{
                    // 修改模块属性 start
                    if (data.selected.length) {//获取当前节点
                        var node = data.instance.get_node(data.selected[0]);
                        var panodeId = data.instance.get_parent(data.selected[0]);// get parent id
                        var panode = data.instance.get_node( panodeId );// obj
                    }
                    if(typeof panode == 'undefined') {
                        return;
                    }
                    //二级结构，特殊处理
                    var SetID = panode.id;
                    if(level ==2) {
                        var SetID = desetid;
                        var SetName = cookie.get('defaultAppName');
                        $("#editmodulebe").text("所属业务");
                    } else {
                        var SetID = panode.id;
                        var SetName = $(panode).html();
                        $("#editmodulebe").text("所属集群");
                    }
                    $.post("/App/getMaintainers",
                        {}
                        ,function(result) {
                            rere = $.parseJSON(result);
                            {
                                var opstr='';
                                $.each( rere, function( i, x ) {
                                    var uin=i;
                                    opstr += ' <option value="' + uin + '">' + uin + '(' + x + ')' + '</option>';
                                });
                                $("#editOperator").html(opstr);
                                $("#editBakOperator").html(opstr);
                                $("#editOperator").val(node.original.operator);
                                $("#editBakOperator").val(node.original.bakoperator);
                                $("#editOperator").select2();
                                $("#editBakOperator").select2();
                            }
                        });

                        $('[jstree-active]').removeAttr('jstree-active');
                        $('.jstree-clicked').attr('jstree-active',true);

                        var ModuleID = node.id;
                        var ModuleName = $(node.text)[0].innerHTML;
                        $("#editmodulegroup").text(SetName);
                        $("#editmoduleModuleName").val(ModuleName);
                        $("#edit_module_property").html(ModuleName);
                        $("#editmoduleModuleName").attr('SetID',SetID);
                        $("#editmoduleModuleName").attr('ModuleID',ModuleID);
                        $('.creat-container').css('display','none');
                        $('.edit-module-container').fadeIn();
                        $("#editmoduleModuleName").focus();
                        return false;
                    // 修改模块属性 end
                }
            }
        });
        //plugin11_demo2_js_end
        CC.hostConf.optionFn();
        /*拓扑树根据浏览器调整自身高度*/
            function treeHeightChange () {
                $('.host-sidebar-left').css('position','fixed');
                setTimeout(function (){
                    $('.c-host-side').css('height',$(window).outerHeight()-70-20);
                    $('.c-conf-inner').css('height',$(window).outerHeight()-70-20-40);
                    if ($('.c-host-side').height()>820)$('.c-host-side').css('height',820);
                    if ($('.c-host-side').height()>780)$('.c-conf-inner').css('height',780);
                    $('.conf-right-empty').css('height',$(window).outerHeight()-70-20);
                    if ($('.conf-right-empty').height()>820)$('.conf-right-empty').css('height',820);
                    $(".c-conf-tree").css('height',$('.c-host-side').outerHeight()-$('.conf-free-group').outerHeight()-$('.c-host-side>h4').outerHeight()-37);
                },200)

            }
            treeHeightChange();
            $(window).resize(function (){
                treeHeightChange();
            })
        /*拓扑树根据浏览器调整自身高度*/
    }
    /*重新请求树的数据*/
    CC.hostConf.addItem = function() {
        var ApplicationID = cookie.get('defaultAppId');
        $.post("/topology/gettopdata",
            {ApplicationID:ApplicationID}
            ,function(result) {
                re = $.parseJSON(result);
                appId = re.appId;
                appName = re.appName;
                topo = $.parseJSON(re.topo);
                level = re.Level;
                desetid = re.deSetID;
                defaultapp = re.Default;
                emptys = re.emptys;
                if (level == 2) {
                    var noset = true;
                }
                else {
                    var noset = false;
                }
                var hostTreeview = $("#confTreeContainer").data("kendoTreeView");
                hostTreeview.setDataSource(new kendo.data.HierarchicalDataSource({
                    data: [{
                        id: appId,
                        text: appName,
                        noset: noset,
                        spriteCssClass: 'c-icon icon-app application',
                        icon: 'c-icon icon-app application',
                        type: 'application',
                        number: 200,
                        expanded: true,
                        items: topo
                    }]
                }));
                CC.hostConf.optionFn();
            });
    }
    //数据加载选项应该注册的事件
    CC.hostConf.optionFn=function(){
            //隐藏set菜单
            $('#confTreeContainer').find('.c-icon.hide').closest('.k-top').hide();
            $('#confTreeContainer').find('.c-icon.hide').closest('.k-mid').hide();
            $('#confTreeContainer').find('.c-icon.hide').closest('.k-bot').hide();

            //改变默认样式
            if(emptys == 1) {

                $('.conf-right-empty').css('height',$(window).outerHeight()-70-20);
                if ($('.conf-right-empty').height()>820)$('.conf-right-empty').css('height',820);

                //if ($('.c-host-side').height()>780)$('.c-conf-inner').css('height',780);
                if(level==3) {
                    $('.creat-container').css('display','none');
                    $('.creat-group-container').css('display','none');
                } else {
                    $('.creat-container').css('display','none');
                    $('.creat-module-container').css('display','none');
                }
            } else {
                    $('.creat-nothing-container').show();
                    $('.creat-nothing-container').find('.creat-module-btn').click(function(event) {
                         $('.jstree-anchor .creat-module-btn').trigger('click');
                    });
            }

            //新增模块事件
            // $(".creat-module-btn").on("click", function(e) {
            //     if(defaultapp==1) {
            //         showWindows('资源池业务不能操作', 'notice');
            //         return;
            //     }
            //     $.post("/App/getMaintainers",
            //         {}
            //         ,function(result) {
            //             rere = $.parseJSON(result);
            //             {
            //                 var opstr='';
            //                 $.each( rere, function( i, x ) {
            //                     var uin=i;
            //                     opstr += ' <option value="' + uin + '">' + uin + '(' + x + ')' + '</option>';
            //                 });
            //                 $("#Operator").html(opstr);
            //                 $("#BakOperator").html(opstr);
            //                 $("#Operator").select2();
            //                 $("#BakOperator").select2();
            //             }
            //         });
            //     var domNode = $(e.target).closest('.k-item'),
            //         node = $("#confTreeContainer").data('kendoTreeView').dataItem(domNode);
            //     //对于二级业务来说重新取其setid
            //     if(level==2) {
            //         var SetID = desetid;
            //         var SetName = cookie.get('defaultAppName');
            //         $("#newmodulebe").text("所属业务");
            //     } else {
            //         var SetID = node.id;
            //         var SetName = node.text;
            //         $("#newmodulebe").text("所属集群");
            //     }
            //     $("#newmodulegroupname").text(SetName);
            //     $("#newmoduleModuleName").attr('setid',SetID);
            //     $('.creat-container').css('display','none');
            //     $('.creat-module-container').fadeIn();
            //     $("#newmoduleModuleName").focus();
            //     e.preventDefault();
            //     e.stopPropagation();
            //     return false;
            // });

            //修改模块属性事件
            // $('.c-conf-tree .k-sprite.c-icon.icon-modal,.conf-free-group .k-in').parent().css('cursor','pointer').on("click dblclick", function(e) {
            //     var padomNode=$(e.target).closest('.k-item').closest('.k-group').closest('.k-item'),
            //         panode = $("#confTreeContainer").data('kendoTreeView').dataItem(padomNode);
            //     var domNode = $(e.target).closest('.k-item'),
            //         node = $("#confTreeContainer").data('kendoTreeView').dataItem(domNode);
            //     if(typeof panode == 'undefined') {
            //         return;
            //     }
            //     //二级结构，特殊处理
            //     var SetID = panode.id;
            //     if(level ==2) {
            //         var SetID = desetid;
            //         var SetName = cookie.get('defaultAppName');
            //         $("#editmodulebe").text("所属业务");
            //     } else {
            //         var SetID = panode.id;
            //         var SetName = panode.text;
            //         $("#editmodulebe").text("所属集群");
            //     }
            //     $.post("/App/getMaintainers",
            //         {}
            //         ,function(result) {
            //             rere = $.parseJSON(result);
            //             {
            //                 var opstr='';
            //                 $.each( rere, function( i, x ) {
            //                     var uin=i;
            //                     opstr += ' <option value="' + uin + '">' + uin + '(' + x + ')' + '</option>';
            //                 });
            //                 $("#editOperator").html(opstr);
            //                 $("#editBakOperator").html(opstr);
            //                 $("#editOperator").val(node.operator);
            //                 $("#editBakOperator").val(node.bakoperator);
            //                 $("#editOperator").select2();
            //                 $("#editBakOperator").select2();
            //             }
            //         });

            //     var ModuleID = node.id;
            //     var ModuleName = node.text;
            //     $("#editmodulegroup").text(SetName);
            //     $("#editmoduleModuleName").val(ModuleName);
            //     $("#edit_module_property").html(ModuleName);
            //     $("#editmoduleModuleName").attr('SetID',SetID);
            //     $("#editmoduleModuleName").attr('ModuleID',ModuleID);
            //     $('.creat-container').css('display','none');
            //     $('.edit-module-container').fadeIn();
            //     $("#editmoduleModuleName").focus();
            //     return false;
            // });

            //修改集群属性事件
            $('.creat-module-btn').parent().on('click',function (e){
                //二级业务不能操作
                if(level == 2) {
                    return;
                }
                if(defaultapp==1) {
                    showWindows('资源池业务不能操作', 'notice');
                    return;
                }
                var domNode = $(e.target).closest('.k-item'),
                    node = $("#confTreeContainer").data('kendoTreeView').dataItem(domNode);
                var ApplicationID = cookie.get('defaultAppId');
                var SetID = node.id;
                $("#editSetSetName").attr('setid',SetID);
                $.post("/Set/getSetInfoById",
                    {ApplicationID:ApplicationID,SetID:SetID}
                    ,function(result) {
                        rere = $.parseJSON(result);
                        if (rere.success == false) {
                            showWindows(rere.errInfo, 'notice');
                        }
                        else {
                            var set = rere.set;
                            $("#editset_property").html(set.SetName);
                            $("#editSetSetName").val(set.SetName);
                            $("#editSetEnviType").val(set.EnviType);
                            $("#edit_set_setenctype label").removeClass('active');
                            if(set.EnviType == 1) {
                                $("#edit_option1").parent().addClass('active');
                            }else if(set.EnviType == 2){
                                $("#edit_option2").parent().addClass('active');
                            }
                            else{
                                $("#edit_option3").parent().addClass('active');
                            }
                            $("#edit_set_sersta input").bootstrapSwitch({onText: '开放',
                                offText: '关闭'});
                            $('#edit_set_sersta input').bootstrapSwitch('state', set.ServiceStatus == 1);
                            $("#editSetChnName").val(set.ChnName);
                            if(set.Capacity != 0) {
                                $("#editSetCapacity").val(set.Capacity);
                            }
                            $("#editSetDes").val(set.Description);
                            $("#editOpenstatus").val(set.Openstatus);
                        }
                    });
                $('.creat-container').css('display','none');
                $('.edit-group-container').fadeIn();
                $("#editSetSetName").focus();
                e.preventDefault();
                e.stopPropagation();
                return false;
            })

            //select事件
            $("#exsitGroup").find("select").change(function() {
                    var IsClone = $("#newcloneset").attr("checked");
                    console.log(IsClone);
                }
            )

            //保存新增set
            $("#save_new_set").on('click',function(e){
                e.preventDefault();
                e.stopPropagation();
                var ApplicationID = cookie.get('defaultAppId');
                var SetName = $.trim($("#newsetname").val());
                var EnviType = $("#new_set_envtype .btn-group .active input").val();
                var SerSw =  $("#new_set_sersta input").bootstrapSwitch('state');
                var ServiceStatus = SerSw? 1:0;
                var ChnName = $.trim($("#newSetChnName").val());
                var Des = $.trim($("#newSetdes").val());
                var Capacity = $.trim($("#newSetCapacity").val());
                var Openstatus = $.trim($("#newOpenstatus").val());
                if(SetName.length==0 ){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">集群名不能为空</span>'
                    });
                    diaCopyMsg.show($("#newsetname").get(0));
                    return ;
                }
                if( SetName.length > 10){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">集群名过长</span>'
                    });
                    diaCopyMsg.show($("#newsetname").get(0));
                    return ;
                }
                if(ChnName.length > 32){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">集群中文名过长</span>'
                    });
                    diaCopyMsg.show($("#newSetChnName").get(0));
                    return ;
                }
                if(Des.length > 250){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">描述过长</span>'
                    });
                    diaCopyMsg.show($("#newSetdes").get(0));
                    return ;
                }
                if(Capacity >9990000) {
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">请输入合理容量</span>'
                    });
                    diaCopyMsg.show($("#newSetCapacity").get(0));
                    return ;
                }
                $.post("/Set/newSet",
                    {ApplicationID:ApplicationID,SetName:SetName,EnviType:EnviType,ServiceStatus:ServiceStatus,
                        ChnName:ChnName,Capacity:Capacity,Des:Des,Openstatus:Openstatus}
                    ,function(result) {
                        rere = $.parseJSON(result);
                        if (rere.success == false) {
                            showWindows(rere.errInfo, 'notice');
                            return;
                        }
                        else {
                            showWindows('新增集群成功！', 'success');
                            window.location.reload();
                        }
                    });
            });

            //删除set
            $("#editsetdelete").on('click',function(e){
                e.preventDefault();
                e.stopPropagation();
                var gridBatDel = dialog({
                    title:'确认',
                    width:250,
                    content: '是否确认删除集群？',
                    okValue:"确定",
                    cancelValue:"取消",
                    ok:function (){
                        var ApplicationID = cookie.get('defaultAppId');
                        var SetID = $("#editSetSetName").attr('setid');
                        $.post("/Set/delSet",
                            {ApplicationID:ApplicationID,SetID:SetID}
                            ,function(result) {
                                rere = $.parseJSON(result);
                                if (rere.success == false) {
                                    showWindows(rere.errInfo, 'notice');
                                    return;
                                } else {
                                    showWindows('删除集群成功！', 'success');
                                    window.location.reload();
                                }
                            });
                    },
                    cancel: function () {
                    }
                });
                gridBatDel.showModal();
            });

            /*保存set修改*/
            $("#editsetsave").on('click',function(e){
                e.preventDefault();
                e.stopPropagation();
                var ApplicationID = cookie.get('defaultAppId');
                var SetID = $("#editSetSetName").attr('setid');
                var SetName = $.trim($("#editSetSetName").val());
                var EnviType = $("#edit_set_setenctype .btn-group .active input").val();
                var ServiceStatus = $("#edit_set_sersta input").bootstrapSwitch('state') ? 1:0;
                var ChnName = $.trim($("#editSetChnName").val());
                var Des = $.trim($("#editSetDes").val());
                var Capacity = $.trim($("#editSetCapacity").val());
                var Openstatus = $.trim($("#editOpenstatus").val());
                if(SetName.length==0 ){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">集群名不能为空</span>'
                    });
                    diaCopyMsg.show($("#editSetSetName").get(0));
                    return ;
                }
                if(SetName.length > 64){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">集群名过长</span>'
                    });
                    diaCopyMsg.show($("#editSetSetName").get(0));
                    return ;
                }
                if(ChnName.length > 32){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">集群中文名过长</span>'
                    });
                    diaCopyMsg.show($("#editSetChnName").get(0));
                    return ;
                }
                if(Des.length > 250){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">描述过长</span>'
                    });
                    diaCopyMsg.show($("#editSetDes").get(0));
                    return ;
                }
                if(Capacity >9990000) {
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">请输入合理容量</span>'
                    });
                    diaCopyMsg.show($("#editSetCapacity").get(0));
                    return ;
                }
                $.post("/Set/editSet",
                    {ApplicationID:ApplicationID,SetName:SetName,SetEnviType:EnviType,ServiceStatus:ServiceStatus,
                        ChnName:ChnName,Capacity:Capacity,Des:Des,SetID:SetID,Openstatus:Openstatus}
                    ,function(result) {
                        rere = $.parseJSON(result);
                        if (rere.success == false) {
                            showWindows(rere.errInfo, 'notice');
                            return;
                        }
                        else {
                            showWindows('修改集群成功！', 'success');
                            window.location.reload();
                            return;
                        }
                    });
            });

        //新建集群取消事件
        $("#new_set_cancel").on('click',function(e) {
            e.preventDefault();
            e.stopPropagation();
            $(".col-md-12.creat-container.creat-group-container").hide();
        });

        //新建模块取消事件
        $("#new_module_cancel").on('click',function(e) {
            e.preventDefault();
            e.stopPropagation();
            $(".col-md-12.creat-container.creat-module-container").hide();
        });

            //保存新增模块
            $("#newsavemodule").on('click',function(e){
                e.preventDefault();
                e.stopPropagation();
                var ApplicationID = cookie.get('defaultAppId');
                var SetID = $("#newmoduleModuleName").attr('setid');
                var ModuleName = $.trim($("#newmoduleModuleName").val());
                var Operator = $("#Operator").val();
                var BakOperator = $("#BakOperator").val();
                if(ModuleName.length==0 ){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">模块名不能为空</span>'
                    });
                    diaCopyMsg.show($("#newmoduleModuleName").get(0));
                    return ;
                }
                if(ModuleName.length>10){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">模块名过长</span>'
                    });
                    diaCopyMsg.show($("#newmoduleModuleName").get(0));
                    return ;
                }
                $.post("/Module/newModule",
                    {ApplicationID:ApplicationID,SetID:SetID,ModuleName:ModuleName,Operator:Operator,BakOperator:BakOperator}
                    ,function(result) {
                        rere = $.parseJSON(result);
                        if (rere.success == false) {
                            showWindows(rere.errInfo, 'notice');
                            return;
                        }
                        else {
                            showWindows('新增模块成功！', 'success');
                            window.location.reload();
                            return;
                        }
                    });
            });

            //保存模块修改
            $("#editmodulesave").on('click',function(e){
                var ApplicationID = cookie.get('defaultAppId');
                var SetID = $("#editmoduleModuleName").attr("setid");
                var ModuleID = $("#editmoduleModuleName").attr("moduleid");
                var ModuleName = $.trim($("#editmoduleModuleName").val());
                var Operator = $("#editOperator").val();
                var BakOperator = $("#editBakOperator").val();
                if(ModuleName.length==0 ){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">模块名不能为空</span>'
                    });
                    diaCopyMsg.show($("#editmoduleModuleName").get(0));
                    return ;
                }
                if(ModuleName.length>10){
                    var diaCopyMsg = dialog({
                        quickClose: true,
                        align: 'left',
                        padding:'5px 5px 5px 10px',
                        skin: 'c-Popuplayer-remind-left',
                        content: '<span style="color:#fff">模块名过长</span>'
                    });
                    diaCopyMsg.show($("#editmoduleModuleName").get(0));
                    return ;
                }
                $.post("/Module/editModule",
                    {ApplicationID:ApplicationID,ModuleID:ModuleID,ModuleName:ModuleName,SetID:SetID,Operator:Operator,BakOperator:BakOperator}
                    ,function(result) {
                        rere = $.parseJSON(result);
                        if (rere.success == false) {
                            showWindows(rere.errInfo, 'notice');
                            return;
                        }
                        else {
                            cookie.set('level',2,5);
                            showWindows('更新模块成功！', 'success');
                            return window.location.reload();
                        }
                    });
            });

            //删除模块
            $("#editmoduledelete").on('click',function(e){
                e.preventDefault();
                e.stopPropagation();
                var gridBatDel = dialog({
                    title:'确认',
                    width:250,
                    content: '是否删除选中模块',
                    okValue:"确定",
                    cancelValue:"取消",
                    ok:function (){
                        var ApplicationID = cookie.get('defaultAppId');
                        var SetID = $("#editmoduleModuleName").attr("setid");
                        var ModuleID = $("#editmoduleModuleName").attr("moduleid");
                        $.post("/Module/delModule",
                            {ApplicationID:ApplicationID,ModuleID:ModuleID,SetID:SetID}
                            ,function(result) {
                                rere = $.parseJSON(result);
                                if (rere.success == false) {
                                    showWindows(rere.errInfo, 'notice');
                                    return;
                                }
                                else {
                                    showWindows('删除模块成功！', 'success');
                                    window.location.reload();
                                }
                            });
                    },
                    cancel: function () {
                    }
                });
                gridBatDel.showModal();

            });

        $(document).on("click", ".conf-delete-link", function(e) {
            e.preventDefault();
            var hostTreeview = $("#confTreeContainer").data("kendoTreeView");
            if ($(".k-state-selected").find('.application').length>0)return false;
            hostTreeview.remove($(".k-state-selected"));
        });

    }
})(window);
/* 拓扑配置 end */
$(function() {
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

    $("[name='my-checkbox']").bootstrapSwitch();
    $("#new_set_sersta input").bootstrapSwitch(
        {onText: '开放',
            offText: '关闭'}
    );
    $(".c-conf-tree").mCustomScrollbar({
        //setHeight: 400, //设置高度
        theme: "minimal-dark" //设置风格
    });
    // $('#date1').kendoDatePicker({
    //     value : new Date(),
    //     format : "yyyy-MM-dd"
    // });
    // $('#date2').kendoDatePicker({
    //     value : new Date(),
    //     format : "yyyy-MM-dd"
    // });
    CC.hostConf.init();

})
