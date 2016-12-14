<?php
if (!defined('BASEPATH')) exit('No direct script access allowed');

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

/*
 * layout 公用布局库
 */
class Layout{
    public $obj;

    public $layout;
 
    public function __construct($layout = "layout"){
        $this->obj = & get_instance();
        $this->layout = $layout;
    }
 
    public function setLayout($layout){
        $this->layout = $layout;
    }
 
    public function view($view, $data=array(), $return = false){
        $data['version'] = STATIC_VERSION;
        $data['isNewUser'] = $this->obj->session->userdata('newUser');
        $data['environment'] = ENVIRONMENT;
        $data['content'] = $this->obj->load->view($view, $data, true);
 
        if($return) {
            $output = $this->obj->load->view($this->layout, $data, true);
            return $output;
        }else {
            $this->obj->load->view($this->layout, $data, false);
        }
    }
}

?>