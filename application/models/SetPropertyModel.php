<?php

/**
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2016 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 */

class SetPropertyModel extends Cc_Model{

	public function __construct(){
		parent::__construct();
	}

	/**
	* @查询所有set属性编码
	*/
	public function getSetPropertyCode() {
        $this->db->select('*');
        $this->db->from('cc_SetProperty');
        $query = $this->db->get();
        return $query->num_rows() ? $query->result_array() : array();
	}

    /**
     * @新增set属性编码
     */
    public function addSetPropertyCode($data) {
        $this->db->insert('cc_SetProperty',$data);
    }

    /**
     * @更新set属性编码
     */
    public function updateSetPropertyCode($data, $Id) {
        $this->db->where('ID',$Id);
        $this->db->update('cc_SetProperty',$data);
    }

    /**
     * @删除set属性编码
     */
    public function deleteSetPropertyCode($Id) {
        $this->db->where('ID',$Id);
        $this->db->delete('cc_SetProperty');
    }

    /**
     * @获取Set属性
     */
    public function getSetProperty() {
        $this->db->select('*');
        $this->db->from('cc_SetProperty');
        $query = $this->db->get();
        return $query->num_rows() ? $query->result_array() :array();
    }

    /**
     * @根据属性获取set
     */
    public function getSetsByProperty($appId, $serviceStatus, $enviType) {
        $this->db->select('SetID,SetName');
        $this->db->from('SetBase');
        $this->db->where('ApplicationID', $appId);
        if(!empty($serviceStatus)) {
            $this->db->where('ServiceStatus', $serviceStatus);
        }
        if(!empty($enviType)) {
            $this->db->where('EnviType', $enviType);
        }

        $query = $this->db->get();
        return $query->num_rows() ? $query->result_array() :array();
    }

    /**
     * @根据set属性和setId获取集群
     */
    public function getSetsByPropertyAndSetId($appId, $setIdArr, $setServiceStatusArr, $setEnviTypeArr) {
        $this->db->select('SetID,SetName');
        $this->db->from('SetBase');
        $this->db->where('ApplicationID', $appId);
        if( !empty($setIDArr)) {
            $this->db->where_in('SetID', $setIdArr);
        }
        if(!empty($setEnviTypeArr)) {
            $this->db->where_in('EnviType', $setEnviTypeArr);
        }
        if(!empty($setServiceStatusArr)) {
            $this->db->where_in('ServiceStatus', $setServiceStatusArr);
        }
        $query = $this->db->get();
        return $query->num_rows() ? $query->result_array() : array();
    }

}