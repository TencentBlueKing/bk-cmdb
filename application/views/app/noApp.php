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
<div class="content-wrapper">
    <div class="no-host-content">
        <img style="height:201px;" src="<?php echo STATIC_URL;?>/static/img/expre_403.png" style="height:201px;">
        <h4 class="pt15">对不起，您当前没有可操作的业务，您可尝试如下操作</h4>
        <ul class="pt15" style="width:300px;">
                <li class="text-left">点此
                    <div id="home_creat_one" style="display:inline-block;">
                        <a href="/app/newapp" id="home_creat_one">新建业务</a>
                    </div>
                </li>
            <li class="text-left">联系您公司已有权限的同事为您开通权限</li>
        </ul>
    </div>
</div>
<div class="control-sidebar-bg"></div>
</div>


<script src="<?php echo STATIC_URL;?>/static/assets/mCustomScrollbar-3.0.9/jquery.mCustomScrollbar.concat.min.js?version=<?php echo $version;?>"></script>