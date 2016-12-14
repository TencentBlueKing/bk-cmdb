<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class HostSourceModel extends Cc_Model {

	public function _construct(){
		parent::_construct();
	}

    /**
     * @增加主机源
     * @return array
     */
    public function listHostSource() {
        $this->db->select('*');
        $this->db->from('cc_HostSource');
        $query = $this->db->get();
        return $query->num_rows() ? $query->result_array() : array();
    }

    /**
     * @更新主机源
     */
    public function updateHostSource($Id, $data) {
        $this->db->where('ID', $Id);
        $this->db->update('cc_HostSource', $data);
    }

    /**
     * @删除主机源
     * @return array
     */
    public function deleteHostSource($Id) {
        $this->db->where('ID', $Id);
        $this->db->delete('cc_HostSource');
    }

    /**
     * @根据类型获取字典集合
     */
    public function getHostSourceOrder($companyCode) {
        $this->db->select('*');
        $this->db->from('HostSource');
        $this->db->where('IsPublic', 1);
        if(!empty($companyCode)) {
            $this->db->or_where('CompanyCode', $companyCode);
        }
        $this->db->order_by('SourceCode', 'desc');
        $query = $this->db->get();
        return $query && $query->num_rows() >0 ? $query->result_array() : array();
    }
}