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
import STATIC_NAVIGATION from '@/assets/json/static-navigation.json'

const state = {
    classifications: [],
    invisibleClassifications: ['bk_host_manage', 'bk_biz_topo'],
    interceptStaticModel: {
        'bk_host_manage': ['resource'],
        'bk_back_config': ['event', 'model', 'audit']
    }
}

const getters = {
    classifications: state => state.classifications,
    activeClassifications: state => {
        let classifications = state.classifications
        // 1.去掉停用模型
        let activeClassifications = classifications.map(classification => {
            let activeClassification = {...classification}
            activeClassification['bk_objects'] = activeClassification['bk_objects'].filter(model => !model['bk_ispaused'])
            return activeClassification
        })
        // 2.去掉无启用模型的分类和不显示的分类
        activeClassifications = activeClassifications.filter(classification => {
            let {
                'bk_classification_id': bkClassificationId,
                'bk_objects': bkObjects
            } = classification
            return !state.invisibleClassifications.includes(bkClassificationId) && Array.isArray(bkObjects) && bkObjects.length
        })
        return activeClassifications
    },
    // 可用分类中被授权的分类
    authorizedClassifications: (state, getters, rootState, rootGetters) => {
        let modelAuthority = rootGetters['userPrivilege/privilege']['model_config'] || {}
        let authorizedClassifications = JSON.parse(JSON.stringify(getters.activeClassifications))
        if (!rootGetters.admin) {
            // 1.去除无权限分类
            authorizedClassifications = authorizedClassifications.filter(classification => {
                return modelAuthority.hasOwnProperty(classification['bk_classification_id'])
            })
            // 2.去除分类下无权限的模型
            authorizedClassifications.forEach(classification => {
                classification['bk_objects'] = classification['bk_objects'].filter(model => {
                    return modelAuthority[classification['bk_classification_id']].hasOwnProperty(model['bk_obj_id'])
                })
            })
        }
        return authorizedClassifications.filter(({bk_objects: bkObjects}) => bkObjects.length)
    },
    // 被授权的导航(包含主机管理、后台配置、通用模型)
    authorizedNavigation: (state, getters, rootState, rootGetters) => {
        let authorizedClassifications = JSON.parse(JSON.stringify(getters.authorizedClassifications))
        // 构造模型导航数据
        let navigation = authorizedClassifications.map(classification => {
            return {
                'icon': classification['bk_classification_icon'],
                'id': classification['bk_classification_id'],
                'name': classification['bk_classification_name'],
                'children': classification['bk_objects'].map(model => {
                    return {
                        'path': model['bk_obj_id'] === 'biz' ? '/business' : `/general-model/${model['bk_obj_id']}`,
                        'id': model['bk_obj_id'],
                        'name': model['bk_obj_name'],
                        'icon': model['bk_obj_icon'],
                        'classificationId': model['bk_classification_id']
                    }
                })
            }
        })
        let staticNavigation = JSON.parse(JSON.stringify(STATIC_NAVIGATION))
        // 检查主机管理、后台配置权限
        if (!rootGetters.admin) {
            let sysConfig = {
                'bk_host_manage': rootGetters['userPrivilege/privilege']['sys_config']['global_busi'] || [],
                'bk_back_config': rootGetters['userPrivilege/privilege']['sys_config']['back_config'] || []
            }
            for (let classificationId in staticNavigation) {
                if (sysConfig.hasOwnProperty(classificationId)) {
                    staticNavigation[classificationId].children = STATIC_NAVIGATION[classificationId].children.filter(({id}) => {
                        if (state.interceptStaticModel[classificationId].includes(id)) {
                            return sysConfig[classificationId].includes(id)
                        }
                        return !['permission', 'model'].includes(id) // 权限管理、模型管理仅管理员拥有且后台接口不返回其配置
                    })
                }
            }
        }
        return [
            staticNavigation['bk_index'],
            staticNavigation['bk_host_manage'],
            ...navigation,
            staticNavigation['bk_back_config']
        ]
    }
}

const actions = {
    /**
     * 添加模型分类
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    createClassification ({ commit, state, dispatch }, { params }) {
        return $http.post(`object/classification`, params)
    },

    /**
     * 删除模型分类
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 分类数据记录id
     * @return {promises} promises 对象
     */
    deleteClassification ({ commit, state, dispatch }, { id }) {
        return $http.delete(`object/classification/${id}`)
    },

    /**
     * 更新模型分类数据
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Number} id 分类数据记录id
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    updateClassification ({ commit, state, dispatch }, { id, params }) {
        return $http.put(`object/classification/${id}`, params)
    },

    /**
     * 查询模型分类列表
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchClassifications ({ commit, state, dispatch }) {
        return $http.post(`object/classifications`)
    },

    /**
     * 查询模型分类及附属模型信息
     * @param {Function} commit store commit mutation hander
     * @param {Object} state store state
     * @param {String} dispatch store dispatch action hander
     * @param {Object} params 参数
     * @return {promises} promises 对象
     */
    searchClassificationsObjects ({ commit, state, dispatch, rootGetters }, { params = {}, config }) {
        return $http.post(`object/classification/${rootGetters.supplierAccount}/objects`, params, config).then(classifications => {
            commit('setClassificationsObjects', classifications)
        })
    }
}

const mutations = {
    setClassificationsObjects (state, classifications) {
        state.classifications = classifications
    },
    updateClassify (state, classification) {
        let activeClassification = state.classifications.find(({bk_classification_id: bkClassificationId}) => bkClassificationId === classification['bk_classification_id'])
        if (activeClassification) {
            activeClassification['bk_classification_icon'] = classification['bk_classification_icon']
            activeClassification['bk_classification_name'] = classification['bk_classification_name']
        } else {
            state.classifications.push({
                ...{
                    bk_asst_objects: {},
                    bk_classification_icon: 'icon-cc-default',
                    bk_classification_id: '',
                    bk_classification_name: '',
                    bk_classification_type: '',
                    bk_objects: [],
                    bk_supplier_account: '',
                    id: 0
                },
                ...classification
            })
        }
    },
    updateModel (state, updateModel) {
        let {
            bk_classification_id: bkClassificationId,
            bk_obj_id: bkObjId
        } = updateModel
        let currentClassify = state.classifications.find(classify => classify['bk_classification_id'] === bkClassificationId)
        let curModel = currentClassify['bk_objects'].find(model => model['bk_obj_id'] === bkObjId)
        if (updateModel.hasOwnProperty('position')) {
            curModel['position'] = updateModel['position']
        }
    },
    deleteClassify (state, classificationId) {
        let index = state.classifications.findIndex(({bk_classification_id: bkClassificationId}) => bkClassificationId === classificationId)
        state.classifications.splice(index, 1)
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
