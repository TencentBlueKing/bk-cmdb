<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class BaseParameterDataModel extends Cc_Model {
    public function __construct() {
        parent::__construct();
    }

    /**
     * @获取所有字典
     * @return array
     */
    public function getBaseParameterDataList() {
        $this->db->order_by('ParentCode', 'desc');
        $this->db->order_by('DataType', 'desc');
        $this->db->order_by('ParameterName', 'desc');
        $this->db->order_by('ParameterCode', 'desc');
        $query = $this->db->get('BaseParameterData');
        return $query && $query->num_rows() > 0 ? $query->result_array() : array();
    }

    /**
     * @根据类型获取字典集合
     * @return array
     */
    public function getBaseParameterDataByDataType($dataType) {
        if (!$dataType) {
            return array();
        }
        $this->db->where('DataType', $dataType);
        $query = $this->db->get('BaseParameterData');
        return $query && $query->num_rows() > 0 ? $query->result_array() : array();
    }

    /**
     * @根据类型获取字典集合
     * @return array
     */
    public function getBaseParameterDataByDataTypeOrder($dataType, $ParentCode) {
        if (!$dataType) {
            return array();
        }
        $this->db->select('*');
        $this->db->from('BaseParameterData');
        $this->db->where('DataType', $dataType);
        if (isset($ParentCode)) {
            $this->db->group_start();
            $this->db->where('ParentCode', $ParentCode);
            $this->db->or_where('ParentCode', '');
            $this->db->group_end();
        }
        $this->db->order_by('ParameterName', 'desc');
        $query = $this->db->get();
        return $query && $query->num_rows() > 0 ? $query->result_array() : array();
    }

    /**
     * @更新字典信息
     */
    public function updateBaseParameterData($data) {
        $this->db->where('ParameterID', $data['ParameterID']);
        $this->db->update('BaseParameterData', $data);
        $log = array();
        $log['OpType'] = '更新';
        $log['OpTarget'] = '更新字典';
        $log['OpContent'] = '字典变更细节稍后补上';
        CCLog::addOpLogArr($log);
    }

    /**
     * @添加字典
     */
    public function addBaseParameterData($data) {
        $this->db->select('ParameterID');
        $this->db->where('DataType', $data['DataType']);
        $this->db->where('ParameterCode', $data['ParameterCode']);
        $query = $this->db->get('BaseParameterData');

        if ($query && $query->num_rows() == 0) {
            $query->free_result();

            $result = $this->db->insert('BaseParameterData', $data);

            if (!$result) {
                $this->_errInfo = '添加字典失败!';
                $err = $this->db->error();
                CCLog::LogErr('添加字典失败! mysql_errno: ' . $err['code'] . ', mysql_error: ' . $err['message']);
            }

            return $this->db->insert_id();
        }
        $this->_errInfo = '同名字典已存在[DataType=' . $data['DataType'] . ',ParameterCode=' . $data['ParameterCode'] . ']!';
        return false;
    }

    /**
     * @根据id删除字典
     * @return boolean 成功 or 失败
     */
    public function delBaseParameterData($data) {

        $this->db->where('ParameterID', $data['ParameterID']);
        $query = $this->db->get('BaseParameterData');
        if (!$query || $query->num_rows() == 0) {
            $this->_errInfo = '这条记录[ParameterID=' . $data['ParameterID'] . ']不存在';
            return false;
        }

        $this->db->where('ParameterID', $data['ParameterID']);
        $result = $this->db->delete('BaseParameterData');

        if (!$result) {
            $this->_errInfo = '删除字典失败！';
            $err = $this->db->error();
            CCLog::LogErr('删除字典失败！! mysql_errno: ' . $err['code'] . ', mysql_error: ' . $err['message']);
        }

        $log = array();
        $log['OpType'] = '删除';
        $log['OpTarget'] = '删除字典';
        $log['OpContent'] = '删除字典变更细节稍后补上';
        CCLog::addOpLogArr($log);
        return $result;
    }
}