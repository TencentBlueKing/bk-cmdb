/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

import $http from '@/api'

const state = {

}

const getters = {

}

const actions = {
    getApplyRulePreview ({ commit, state, dispatch }, { params, config }) {
        console.log($http)
        return Promise.resolve({
            count: 123,
            info: [
                {
                    id: 1,
                    bk_host_innerip: '120.23.3.3',
                    bk_cloud_id: '直连区域',
                    bk_asset_id: 'No90308',
                    bk_host_name: 'nginx1_lol',
                    diff_value: [
                        {
                            bk_attribute_id: 30,
                            bk_property_id: 'bk_asset_id',
                            bk_property_value: 'a45413',
                            is_conflict: true
                        },
                        {
                            bk_attribute_id: 36,
                            bk_property_id: 'bk_state_name',
                            bk_property_value: 'CN',
                            is_conflict: true
                        },
                        {
                            bk_attribute_id: 45,
                            bk_property_id: 'bk_cpu_mhz',
                            bk_property_value: '425'
                        },
                        {
                            bk_attribute_id: 48,
                            bk_property_id: 'bk_disk',
                            bk_property_value: '512'
                        },
                        {
                            bk_attribute_id: 38,
                            bk_property_id: 'bk_isp_name',
                            bk_property_value: '2'
                        },
                        {
                            bk_attribute_id: 51,
                            bk_property_id: 'create_time',
                            bk_property_value: 1574149251628
                        }
                    ]
                },
                {
                    id: 1,
                    bk_host_innerip: '159.15.52.3',
                    bk_cloud_id: '直连区域',
                    bk_asset_id: 'No90378',
                    bk_host_name: 'nginx1_lol',
                    diff_value: [
                        {
                            bk_attribute_id: 39,
                            bk_property_id: 'bk_host_name',
                            bk_property_value: 'dellx',
                            is_conflict: false
                        },
                        {
                            bk_attribute_id: 42,
                            bk_property_id: 'bk_os_version',
                            bk_property_value: 'window10.27.2',
                            is_conflict: false
                        }
                    ]
                },
                {
                    id: 1,
                    bk_host_innerip: '186.223.47.25',
                    bk_cloud_id: '直连区域',
                    bk_asset_id: 'No90289',
                    bk_host_name: 'nginx2_lol',
                    diff_value: [
                        {
                            bk_attribute_id: 39,
                            bk_property_id: 'bk_host_name',
                            bk_property_value: 'dellx',
                            is_conflict: true
                        },
                        {
                            bk_attribute_id: 52,
                            bk_property_id: 'import_from',
                            bk_property_value: '3',
                            is_conflict: true
                        }
                    ]
                }
            ]
        })
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions
}
