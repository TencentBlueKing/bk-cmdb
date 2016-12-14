/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

$(document).ready(function() {
    var level = 2;
    var appType = 1;

    function format(state) {
        return state.id + "<span class='select-info'>(" + state.text2 + ")</span>";
    }
    function format2(state) {
        return state.id ;
    }
    $('.btn-primary').click(function () {
        if ($(this).hasClass('active')) {
            $(this).removeClass('active');
            return false;
        }
    })

    var uinArr = new Array();
    uinArr.push(uin);
    $('#Maintainers').select2({
        placeholder: '选择运维人员',
        allowClear: true,
        data: userList,
        multiple: true,
        allowClear: true,
        formatResult: format,
        formatSelection: format
    }).select2('val', uinArr);


    $('#ProducterList').select2({
        placeholder: '选择产品人员',
        allowClear: true,
        data: userList,
        multiple: true,
        allowClear: true,
        formatResult: format,
        formatSelection: format
    }).select2('val', uinArr);

    $(".c-panel-content div").click(function(){
        $(this).addClass('buttonMark');
        $(this).siblings().removeClass('buttonMark');
        $(this).parent().parent().removeClass('panelMark');
        $(this).parent().siblings().addClass('headerMark');
        $(this).parent().siblings().find('i').addClass('step-imgMark');

        if($(this).attr('name') == 'c-button-game'){
            $("#stepMark").text('step 3：请填写详细的业务信息');
            $(".c-panel-pageTwo").fadeIn();
            $(".c-panel-pageThree").css("display","none");
            $(".c-panel-pageTwo").addClass('panelMark');
            $(".c-panel-pageTwo .c-panel-header").removeClass('headerMark');
            $(".c-panel-pageTwo .c-panel-header i").removeClass('step-imgMark');
            $(".c-panel-pageTwo .c-panel-content div").removeClass('buttonMark');
            appType = 1;
        }else if($(this).attr('name') == "c-button-notgame"){
            $("#stepMark").text('step 2：请填写详细的业务信息');
            $(".c-panel-pageTwo").css("display","none");
            $(".c-panel-pageThree").fadeIn();
            $("#AppName").val('');
            $("#AppName").focus();
            level = 2;
            appType = 0;
        }else if($(this).attr('name')=="c-button-fqff"){
            $(".c-panel-pageThree").fadeIn();
            $("#AppName").val('');
            $("#AppName").focus();
            level = 3;
        }else if($(this).attr('name')=="c-button-qqqf"){
            $(".c-panel-pageThree").fadeIn();
            $("#AppName").val('');
            $("#AppName").focus();
            level = 2;
        }
    })

    $(".cancelb").click(function(){
        $(".c-panel-pageTwo,.c-panel-pageThree").css("display","none");
        $(".c-panel-pageOne").addClass('panelMark');
        $(".c-panel-pageOne .c-panel-header").removeClass('headerMark');
        $(".c-panel-pageOne .c-panel-content div").removeClass('buttonMark');

    })


    //保存按钮
    $(".btn-primary").click(function () {
        var appName = $.trim($("#AppName").val());
        if ((appName.length > 32) || (appName.length == 0)) {
            var diaCopyMsg = dialog({
                quickClose: true,
                align: 'left',
                padding:'5px 5px 5px 10px',
                skin: 'c-Popuplayer-remind-left',
                content: '<span style="color:#fff">业务名称不合法</span>'
            });
            diaCopyMsg.show($("#AppName").get(0));
            return;
        }
        var maintainers = $("#Maintainers").select2('val');
        maintainers.sort();
        $.unique(maintainers);
        if (maintainers.length == 0) {
            var diaCopyMsg = dialog({
                quickClose: true,
                align: 'left',
                padding:'5px 5px 5px 10px',
                skin: 'c-Popuplayer-remind-left',
                content: '<span style="color:#fff">请配置运维人员</span>'
            });
            diaCopyMsg.show($("#s2id_Maintainers").get(0));
            return;
        }
        if (maintainers.length > 24) {
            var diaCopyMsg = dialog({
                quickClose: true,
                align: 'left',
                padding:'5px 5px 5px 10px',
                skin: 'c-Popuplayer-remind-left',
                content: '<span style="color:#fff">请配置24个以下运维人员</span>'
            });
            diaCopyMsg.show($("#s2id_Maintainers").get(0));
            return;
        }
        var productList = $("#s2id_ProducterList").select2('val');
        if (productList.length == 0) {
            var diaCopyMsg = dialog({
                quickClose: true,
                align: 'left',
                padding:'5px 5px 5px 10px',
                skin: 'c-Popuplayer-remind-left',
                content: '<span style="color:#fff">请配置产品人员</span>'
            });
            diaCopyMsg.show($("#s2id_ProducterList").get(0));
            return;
        }
        if (productList.length > 8) {
            var diaCopyMsg = dialog({
                quickClose: true,
                align: 'left',
                padding:'5px 5px 5px 10px',
                skin: 'c-Popuplayer-remind-left',
                content: '<span style="color:#fff">请配置8个以下产品人员</span>'
            });
            diaCopyMsg.show($("#ProducterList").get(0));
            return;
        }
        var lifeCycle = $("#LifeCycle .btn-group .active input").val();
        $.post("/app/add",
            {
                Level: level,
                Type: appType,
                ApplicationName: appName,
                Maintainers: maintainers,
                ProducterList: productList,
                LifeCycle: lifeCycle
            }
            , function (result) {
                re = $.parseJSON(result);
                if (re.success == false) {
                    showWindows(re.errInfo, 'notice');
                }
                else {
                    showWindows('新增成功！', 'success');
                    window.location.href = '/app/index';
                }
            });
    })

    $('#AppName').blur(function (){
        showstep();
    });
    $('#Maintainers').blur(function (){
        showstep();
    });
    $('#s2id_ProducterList').blur(function (){
        showstep();
    });

    function showstep() {
        var val = $('#AppName').val();
        var maintainers = $('#Maintainers').select2('val');
        var productList = $('#s2id_ProducterList').select2('val');
        var header=$('.c-panel-pageThree').find('.c-panel-header');
        if (val && $.trim(val)!='' && maintainers.length!=0 &&  productList.length != 0) {
            header.addClass('headerMark');
            header.find('i').addClass('step-imgMark');
        }else{
            header.removeClass('headerMark');
            header.removeClass('step-imgMark');
        }
    }
})