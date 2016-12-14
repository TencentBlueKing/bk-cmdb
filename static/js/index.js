/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

window.cookie = {
    get : function(name){
        if(document.cookie.length>0){
            name += '=';
            var allCookie = document.cookie;
            var startIndex = allCookie.indexOf(name);
            if(startIndex > -1) {
                if(startIndex > 1){
                    name = '; '+name;
                    startIndex = allCookie.indexOf(name)+name.length;
                }else{
                    startIndex = startIndex+name.length;
                }
            }else{
                return '';
            }
            var endIndex = allCookie.indexOf(';', startIndex);
            if(endIndex === -1){
                endIndex = document.cookie.length;
            }
            var value = decodeURIComponent(allCookie.substring(startIndex, endIndex));
            return value;
        }

        return '';
    },

    set : function(name, value, expires){
        var exdate = new Date();
        if(expires!=null){
            exdate.setTime(exdate.getTime() + expires*1000);
        }
        document.cookie = name +'='+ encodeURIComponent(value) + (expires==null ? ';path=/' : ';expires='+exdate.toUTCString()+';path=/');
    }
};

window.showWindows = function(msg, level) {      //代码提示框
    if(level=='success') {
        var d = dialog({
            width: 160,
            content: '<div class="c-dialogdiv2"><i class="c-dialogimg-success"></i>' + msg + '</div>'
        });
        d.show();
        setTimeout(function(){d.close();}, 1000);
    } else if(level =='error') {
        var d = dialog({
            title:'错误',
            width:300,
            height:45,
            okValue:"确定",
            ok:function(){},
            content: '<div class="c-dialogdiv2"><i class="c-dialogimg-failure"></i>' + msg + '</div>'
        });
        d.showModal();
    }
    else {
        var d = dialog({
            title:'警告',
            width:300,
            height:45,
            okValue:"确定",
            ok:function(){},
            content: '<div class="c-dialogdiv2"><i class="c-dialogimg-prompt"></i>' + msg + '</div>'
        });
        d.showModal();
    }
}


$(document).ready(function(){
    $.ajaxSetup({
        beforeSend : function(xhr){
            xhr.setRequestHeader('Token', cookie.get('token'));
        },
        complete:function(xhr,response)
        {
            if(xhr.getResponseHeader('Timeout') == 'true')
            {
                window.location.href="/welcome/logout";
                return;
            }
        }

    });

    /**鼠标点击存储cookie*/
    $('.c-sidebar-toggle').on('click',function(){
        var className = $('body').prop('class');
        var isCollapse = className.indexOf('sidebar-collapse') == -1 ? true:false;
        if(isCollapse){
            cookie.set('isCollapse',1);
        }else{
            cookie.set('isCollapse',0);
        }
    });

});
