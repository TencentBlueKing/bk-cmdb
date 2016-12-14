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
<link href="<?php echo STATIC_URL;?>/static/bk/css/bk_app_theme.css" rel="stylesheet">
<link href="<?php echo STATIC_URL; ?>/static/assets/select2-3.5.2/select2.css" rel="stylesheet">
<div class="content-wrapper">
    <!-- 主面板 Main-->
    <section class="content">

        <div class="row">
            <div class="col-sm-12 col-lg-12">
                <div class="c-panel panelMark c-panel-pageOne mb30">
                    <div class="c-panel-header">
                        <span>step 1：请选择您的业务类型</span><i></i>
                    </div>
                    <div class="c-panel-content">
                        <div name="c-button-game">游戏</div>
                        <div name="c-button-notgame" class="c-button-notgame">非游戏<img src="<?php echo STATIC_URL; ?>/static/img/second.jpg"/></div>
                        <i class="clearfix"></i>
                    </div>
                </div>
            </div>
        </div>

        <div class="row">
            <div class="col-sm-12 col-lg-12">
                <div class="c-panel panelMark c-panel-pageTwo mb30">
                    <div class="c-panel-header">
                        <span>step2：您的游戏是分区分服还是全区全服？</span><i></i>
                    </div>
                    <div class="c-panel-content">
                        <div name="c-button-fqff" class="c-button-fqff">分区分服<img src="<?php echo STATIC_URL; ?>/static/img/third.jpg"/></div>
                        <div name="c-button-qqqf" class="c-button-qqqf">全区全服<img src="<?php echo STATIC_URL; ?>/static/img/second.jpg"/></div>
                        <i class="clearfix"></i>
                    </div>
                </div>
            </div>
        </div>




        <div class="row">
            <div class="col-sm-12 col-lg-12">
                <div class="c-panel c-panel-pageThree">

                    <div class="c-panel-header">
                        <span id="stepMark">step 3：请填写详细的业务信息</span><i></i>
                    </div>

                    <div class="c-panel-stepFour">
                        <table class="table table-bordered">
                            <thead >
                            <tr>
                                <th style="width:20%;">属性分组</th>
                                <th>业务属性</th>
                                <th style="width:40%;">属性值</th>
                            </tr>
                            </thead>
                            <tbody>
                            <tr>
                                <td rowspan="2">基本属性</td>
                                <td >名称</td>
                                <td class="tl"><input type="text" value="" class="c-newyw-input c-gridinput" id="AppName"  maxlength="10" placeholder="输入创建的业务名称,必填">
                                    <span class="c-gridinputmust">*</span>
                                </td>
                            </tr>
                            <tr>
                                <td>运维人员</td>
                                <td class="tl">
                                    <div class="operaer c-gridinput" id="Maintainers" ></div>
                                <!--    <select id="Maintainers" style="width:80%" placeholder="请选择运维人员,选填">
                                        <?php /* foreach($UserList as $user):*/?>
                                            <option value="<?php /*echo $user;*/?>"><?php /*echo $user;*/?></option>
                                        <?php /*endforeach */?>
                                    </select>-->
                                    <span class="c-gridinputmust">*</span>
                                </td>
                            </tr>
                            <tr>
                                <td rowspan="2">扩展属性</td>
                                <td>产品人员</td>
                                <td class="tl">
                                    <div id="ProducterList" class="operaer c-gridinput"  ></div>
                                    <!--<select id="ProducterList" style="width:80%" placeholder="请选择产品人员,选填">
                                        <?php /*$length=count($UserList); for($i=0;$i!=$length;$i++){  */?>
                                            <option value="<?php /*echo $UserList[$i];*/?>"><?php /*echo $UserList[$i];*/?></option>
                                        <?php /*} */?>
                                    </select>-->
                                    <span class="c-gridinputmust">*</span>
                                </td>
                            </tr>
                            <tr id="lifecycle">
                                <td>生命周期</td>
                                <td class="tl"  id="LifeCycle">
                                    <div class="btn-group" data-toggle="buttons">
                                        <label class="btn btn-radio active">
                                            <input type="radio" name="options" id="option1" autocomplete="off" value="公测">公测
                                        </label>
                                        <label class="btn btn-radio">
                                            <input type="radio" name="options" id="option2" autocomplete="off" value="内测">内测
                                        </label>
                                        <label class="btn btn-radio">
                                            <input type="radio" name="options" id="option3" autocomplete="off" value="不删档">不删档
                                        </label>
                                    </div>
                                </td>
                            </tr>
                            </tbody>
                        </table>
                        <div class="text-center">
                            <button class="btn btn-default cancelb">取消</button>
                            <button class="btn btn-primary">保存</button>
                        </div>

                    </div>
                    <div class="clearfix"></div>
                </div>
            </div>
        </div>
        <div class="clearfix"></div>
    </section>
</div>
<div class="control-sidebar-bg"></div>
</div>
<script src="<?php echo STATIC_URL;?>/static/assets/select2-3.5.2/select2.js"></script>
<script>
    var userList = $.parseJSON('<?php echo $userListJ;?>');
    var uin = '<?php echo $uin;?>';
</script>
<script src="<?php echo STATIC_URL; ?>/static/js/app.new.js?version=<?php echo $version;?>"></script>

