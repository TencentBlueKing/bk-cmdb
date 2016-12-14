<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class BaseParameterDataLogic extends Cc_Logic {

    public function __construct() {
        parent::__construct();
    }

    /**
     * 查询字典的code与name的对应关系数组
     * @return array()
     */
    public function getBaseParameterDataByDataType($dataType) {
        $data = array();
        $this->load->model('BaseParameterDataModel');
        try {
            $result = $this->BaseParameterDataModel->getBaseParameterDataByDataType($dataType);
            foreach ($result as $value) {
                $data[$value['ParameterCode']] = $value['ParameterName'];
            }
        } catch (Exception $e) {
            CCLog::LogErr("getBaseParameterDataByDataType" . $e->getMessage());
        }

        return $data;
    }

    /**
     * 根据类型获取字典的code与name的对应关系数组
     * @param $dataType   数据类型
     * @return array
     */
    public function getBaseParameterDataNameCodeByDataType($dataType) {
        $data = array();
        $this->load->model('BaseParameterDataModel');
        try {
            $result = $this->BaseParameterDataModel->getBaseParameterDataByDataType($dataType);
            foreach ($result as $value) {
                $data[$value['ParameterName']] = $value['ParameterCode'];
            }
        } catch (Exception $e) {
            CCLog::LogErr("getBaseParameterDataNameCodeByDataType" . $e->getMessage());
        }

        return $data;
    }

    /**
     * 查询字典的code与name的对应关系数组
     * @param $dataType   数据类型
     * @return array
     */
    public function getBaseParameterDataSelectByDataType($dataType) {
        $data = array();
        $this->load->model('BaseParameterDataModel');
        try {
            $result = $this->BaseParameterDataModel->getBaseParameterDataByDataTypeOrder($dataType);
            foreach ($result as $value) {
                $data[] = array('id' => $value['ParameterCode'], 'text' => $value['ParameterName']);
            }
        } catch (Exception $e) {
            CCLog::LogErr("getBaseParameterDataSelectByDataType" . $e->getMessage());
        }

        return $data;
    }

    /**
     * 查询支持的主机源
     * @return array
     */
    public function getSupportHostResource() {
        $data = array();
        $result = array();
        $this->load->model('HostSourceModel');
        try {
            $company = $this->session->userdata('company');
            $result = $this->HostSourceModel->getHostSourceOrder($company);
            foreach ($result as $value) {
                $data[] = array('id' => $value['SourceCode'], 'text' => $value['SourceName']);
            }
        } catch (Exception $e) {
            CCLog::LogErr("getHostSource" . $e->getMessage());
        }

        return $data;
    }

    /**
     * 根据类型获取字典的code与name的对应关系数组
     * @return array
     */
    public function getSupportHostSourceKv() {
        $data = array();
        $result = array();
        $this->load->model('HostSourceModel');
        try {
            $company = $this->session->userdata('company');
            $result = $this->HostSourceModel->getHostSourceOrder($company);
            foreach ($result as $value) {
                $data[$value['SourceCode']] = $value['SourceName'];
            }
        } catch (Exception $e) {
            CCLog::LogErr("getHostSourceKv" . $e->getMessage());
        }

        return $data;
    }
}