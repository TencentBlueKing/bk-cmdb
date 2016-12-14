<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

require_once APPPATH . 'core/Api_Controller.php';

class Set extends Api_Controller {

    public function __construct() {
        parent::__construct();
    }

    /*
     * 获取set属性
     * @return json
     */
    public function getSetProperty() {
        $this->load->logic('SetBaseLogic');
        $result = $this->SetBaseLogic->getSetProperty();
        return $this->outSuccess($result);
    }

    /*
     * 根据set属性获取set
     * @return json
     */
    public function getSetsByProperty() {
        $setEnviType = intval($this->input->post('SetEnviType'));
        $setServiceStatus = $this->input->post('SetServiceStatus');
        $appId = intval($this->input->post('ApplicationID'));
        $this->load->logic('SetBaseLogic');

        $result = $this->SetBaseLogic->getSetsBySetProperty($setEnviType, $setServiceStatus, $appId);
        return $this->outSuccess($result);
    }

    /*
     * 根据set属性获取模块
     * @return json
     */
    public function getModulesByProperty() {
        $setEnviType = intval($this->input->post('SetEnviType'));
        $setServiceStatus = intval($this->input->post('SetServiceStatus'));
        $setId = intval($this->input->post('SetID'));
        $appId = intval($this->input->post('ApplicationID'));

        $this->load->logic('SetBaseLogic');
        $result = $this->SetBaseLogic->getModulesBySetProperty($setEnviType, $setServiceStatus, $appId, $setId);
        return $this->outSuccess($result);
    }

    /*
     * 根据set属性获取主机
     * @return json
     */
    public function getHostsByProperty() {
        $setEnvType = $this->input->post('SetEnviType');
        $setServiceStatus = $this->input->post('SetServiceStatus');
        $setId = $this->input->post('SetID');
        $moduleName = $this->input->post('ModuleName');
        $appId = intval($this->input->post('ApplicationID'));

        $this->load->logic('SetBaseLogic');
        $result = $this->SetBaseLogic->getHostsBySetProperty($setEnvType, $setServiceStatus, $appId, $setId, $moduleName);
        return $this->outSuccess($result);
    }

}