<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class OperationLog extends Cc_Controller {

    public function __construct() {
        parent::__construct();
    }

    /**
     * 日志查询页面
     */
    public function index() {
        $appInfo = array();
        $this->load->Logic('ApplicationBaseLogic');
        $app = $this->ApplicationBaseLogic->getUserApp();
        foreach ($app as $value) {
            $temArr = array();
            $temArr['id'] = $value['ApplicationID'];
            $temArr['text'] = $value['ApplicationName'];
            $appInfo[] = $temArr;
        }

        $data = $this->buildPageDataArr($this->lang->line('log_query'), '/operationlog/index');
        $defaultApp = $this->session->userdata('defaultApp');
        $data['ApplicationID'] = $defaultApp['ApplicationID'];
        $data['app'] = json_encode($appInfo);
        $this->load->library('layout');
        $this->layout->view('operationlog/index', $data);
    }

    /**
     * 日志查询(近1000条)
     * @return json
     */
    public function getOperationLog() {
        $data = array();
        $this->load->Logic('OperationLogLogic');
        $appId = (int)$this->input->post('ApplicationID', true);
        $operator = $this->input->post('Operator', true);
        $opTarget = $this->input->post('OpTarget', true);
        $opContent = $this->input->post('OpContent', true);
        $start = $this->input->post('start', true);
        $end = $this->input->post('end', true);

        if (!$appId) {
            $appId = $this->session->userdata('defaultApp');
            $appId = $appId['ApplicationID'];
        }

        $data = $this->OperationLogLogic->getUserOperationLog($operator, $appId, '', $opTarget, $opContent, '', $start, $end, 0, 1000);
        $result = array('data'=>$data);
        return $this->output->set_output(json_encode($result));
    }
}