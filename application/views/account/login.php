<?php
/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */
?>
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<link href="<?php echo STATIC_URL; ?>/static/css/index.css" rel="stylesheet" type="text/css" />
	<title>logo</title>
</head>

<body style="padding:0px;margin:0px;min-width: 1200px;">

	<div class="head_level1">
		<div><a class="logoimg"></a></div>
	</div>
	<div class="head_level2"></div>

	<div class="bg_body">

		<div class="login_wrapper" style="min-width: 1200px; height: 739px;">

			<div class="login_text login_text1"></div>
		    <div class="login_text login_text2"></div>

			<div id="logo">
				<h1>蓝鲸平台</h1>
				<div class="center">
                    <form action="/account/doLogin" method="post">
                        <input type="text" name="UserName" class="ipt" placeholder="用户名"/>
                        <input type="password" name="Password" class="ipt" placeholder="密码"/>
                        <input type="submit" class="btn_login" value="登入">
                        <?php
                        if(isset($message) && ($message != '')){
                            echo '<p class="login_error_tips"><i class="icon icon_error"></i>' . $message . '</p>';
                        }
						if(isset($cburl) && ($cburl != '')){
                            echo '<input type="hidden" name="cburl" value="' . $cburl . '" />';
                        }
                        ?>
                    </form>

				</div>

			</div>

		</div>
	</div>

	<div class="footer-menu">Copyright © 2013-2016 Tencent. All Rights Reserved.</div>

</body>
</html>