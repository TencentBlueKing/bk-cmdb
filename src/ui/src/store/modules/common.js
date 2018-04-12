/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

import Vue from 'vue'
import Cookies from 'js-cookie'
import {$Axios, $axios} from '@/api/axios'

let alertMsg = Vue.prototype.$alertMsg

const state = {
    bkSupplierAccount: '0',
    allClassify: [], // 所有模型分类及每个分类下的模型
    timezoneList: [],   // 时区列表
    biz: {     // 业务列表
        list: [],
        selected: parseInt(Cookies.get('bk_biz_id')) || -1
    },
    memberList: [],
    authority: null,
    adminAuthority: null,
    isAdmin: window.isAdmin === '1',
    navigation: {},
    usercustom: {},          // 用户字段配置
    globalLoading: false,
    memberLoading: false,
    language: Cookies.get('blueking_language') || 'zh_CN'
}

const getters = {
    bkSupplierAccount: state => state.bkSupplierAccount,
    bkBizId: state => state.biz.selected,
    bkBizList: state => state.biz.list,
    allClassify: state => state.allClassify,
    memberList: state => state.memberList,
    authority: state => state.authority,
    adminAuthority: state => state.adminAuthority,
    isAdmin: state => state.isAdmin,
    navigation: state => state.navigation,
    timezoneList: state => state.timezoneList,
    usercustom: state => state.usercustom,
    globalLoading: state => state.globalLoading,
    memberLoading: state => state.memberLoading,
    language: state => state.language
}

const actions = {
    getBkBizList ({commit, state}) {
        $axios.post(`biz/search/${state.bkSupplierAccount}`, {fields: ['bk_biz_id', 'bk_biz_name']}).then((res) => {
            if (res.result) {
                if (res.data.info && res.data.info.length) {
                    commit('setBkBizList', res.data.info)
                    if (state.biz.selected === -1) { // 如果未选择过，则选中第一个业务
                        commit('setBkBizId', state['biz']['list'][0]['bk_biz_id'])
                    } else { // 如果已经选择过，则需判断缓存的已选择业务是否被删除
                        let isExist = false
                        state.biz.list.map((biz) => {
                            if (state.biz.selected === biz['bk_biz_id']) {
                                isExist = true
                            }
                        })
                        if (!isExist) {
                            commit('setBkBizId', state.biz.list[0]['bk_biz_id'])
                        }
                    }
                } else {
                    commit('setBkBizList', [])
                    commit('setBkBizId', -1)
                }
            } else {
                alertMsg(res['bk_error_msg'])
            }
        })
    },
    getMemberList ({commit, state}, type) {
        state.memberLoading = true
        let baseURL = $axios.defaults.baseURL
        $axios.defaults.baseURL = window.siteUrl
        $axios.get('/user/list', { type }).then((res) => {
            if (res.result) {
                commit('setMemberList', res.data)
            } else {
                alertMsg(res['bk_error_msg'])
            }
            state.memberLoading = false
        }).catch(() => {
            state.memberLoading = false
        })
        $axios.defaults.baseURL = baseURL
    },
    getAuthority ({commit, state}) {
        return $axios.get(`topo/privilege/user/detail/${state.bkSupplierAccount}/${window.userName}`).then(res => {
            if (res.result) {
                commit('setAuthority', res.data)
            } else {
                alertMsg(res['bk_error_msg'])
            }
            return res
        })
    },
    getAllClassify ({commit, state}) {
        if (state.allClassify.length) {
            return Promise.resolve({data: state.allClassify})
        }
        return $axios.post('object/classifications', {}).then(res => {
            if (res.result) {
                $Axios.all(res.data.map(classify => {
                    return $axios.post(`object/classification/${state.bkSupplierAccount}/objects`, {'bk_classification_id': classify['bk_classification_id']}).then(modelRes => {
                        if (!modelRes.result) {
                            alertMsg(modelRes.message)
                        }
                        return modelRes.data || []
                    })
                })).then($Axios.spread(function () {
                    let results = [...arguments]
                    let allClassify = []
                    results.map(classify => {
                        allClassify.push({
                            'bk_classification_icon': classify[0]['bk_classification_icon'],
                            'bk_classification_id': classify[0]['bk_classification_id'],
                            'bk_classification_name': classify[0]['bk_classification_name'],
                            'bk_classification_type': classify[0]['bk_classification_type'],
                            'id': classify[0]['id'],
                            'bk_objects': classify[0]['bk_objects'],
                            'bk_asst_objects': classify[0]['bk_asst_objects']
                        })
                    })
                    commit('setAllClassify', allClassify)
                    commit('setAdminAuthority', allClassify)
                }))
            } else {
                alertMsg(res['bk_error_msg'])
            }
            return res
        })
    }
}

const mutations = {
    setLang (state, language) {
        state.language = language
    },
    setAllClassify (state, classify) {
        state.allClassify = classify
    },
    // 新增分类
    createClassify (state, classify) {
        state.allClassify.push(classify)
    },
    // 删除分类
    deleteClassify (state, classify) {
        let index = state.allClassify.findIndex(({ bk_classification_id: bkClassificationId }) => {
            return classify['bk_classification_id'] === bkClassificationId
        })
        state.allClassify.splice(index, 1)
    },
    // 更新分类位置信息
    updateClassifyPosition (state, classify) {
        let allClassify = state.allClassify
        for (let i = 0; i < allClassify.length; i++) {
            if (allClassify[i]['bk_classification_id'] === classify['bk_classification_id']) {
                for (let j = 0, model = allClassify[i]['bk_objects']; j < model.length; j++) {
                    if (classify['bk_obj_id'] === model[j]['bk_obj_id']) {
                        model[j]['position'] = classify['position']
                        break
                    }
                }
                break
            }
        }
    },
    // 新增模型
    createModel (state, model) {
        let activeClassify = state.allClassify.find(({ bk_classification_id: bkClassificationId }) => {
            return bkClassificationId === model['bk_classification_id']
        })
        activeClassify.push(model)
    },
    // 更新模型
    updateModel (state, model) {
        let activeClassify = state.allClassify.find(({ bk_classification_id: bkClassificationId }) => {
            return bkClassificationId === model['bk_classification_id']
        })
        let activeModel = activeClassify['bk_objects'].find(({ bk_obj_id: bkObjId }) => {
            return bkObjId === model['bk_obj_id']
        })
        activeModel = model
    },
    // 删除模型
    deleteModel (state, model) {
        let allClassify = state.allClassify
        for (let i = 0; i < allClassify.length; i++) {
            if (allClassify[i]['bk_classification_id'] === model['bk_classification_id']) {
                let index = allClassify[i]['bk_objects'].findIndex(({ bk_obj_id: bkObjId }) => {
                    return bkObjId === model['bk_obj_id']
                })
                allClassify[i]['bk_objects'].splice(index, 1)
            }
        }
        state.allClassify = Vue.prototype.$deepClone(allClassify)
    },
    // 修改分类名称后需同步
    updateClassify (state, payload) {
        let allClassify = state.allClassify
        for (let i = 0; i < allClassify.length; i++) {
            if (allClassify[i]['bk_classification_id'] === payload['bk_classification_id']) {
                // 修改分类名
                if (payload.hasOwnProperty('bk_classification_name')) {
                    allClassify[i]['bk_classification_name'] = payload['bk_classification_name']
                }
                // 修改icon
                if (payload.hasOwnProperty('bk_classification_icon')) {
                    allClassify[i]['bk_classification_icon'] = payload['bk_classification_icon']
                }
                // 修改分类下的模型名
                if (allClassify[i].hasOwnProperty('bk_objects')) {
                    let obj = allClassify[i]['bk_objects']
                    let isModelExist = false
                    for (var j = 0; j < obj.length; j++) {
                        if (obj[j]['bk_obj_id'] === payload['bk_obj_id']) {
                            obj[j]['bk_obj_name'] = payload['bk_obj_name']
                            break
                        }
                    }
                }
                break
            }
        }
    },
    setBkBizList (state, list) {
        state.biz.list = list
    },
    setBkBizId (state, selected) {
        Cookies.set('bk_biz_id', selected, { expires: 30, path: '' })
        state.biz.selected = selected
    },
    deleteApplication (state, appId) {
        let applicationList = state.application.list
        for (let i = 0; i < applicationList.length; i++) {
            if (applicationList[i]['ApplicationID'] === appId) {
                applicationList.splice(i, 1)
                if (state.application.selected === appId) {
                    // 如果删除的业务是已选中的，则重新设置当前选中的业务
                    if (applicationList.length) {
                        Cookies.set('selectedApplicationId', applicationList[0]['ApplicationID'], { expires: 30, path: '' })
                        state.application.selected = applicationList[0]['ApplicationID']
                    } else {
                        Cookies.set('selectedApplicationId', '', { expires: 30, path: '' })
                        state.application.selected = ''
                    }
                }
                break
            }
        }
    },
    setMemberList (state, memberList) {
        state.memberList = memberList
    },
    setAuthority (state, authority) {
        state.authority = authority
    },
    setAdminAuthority (state, allClassify) {
        let authority = {
            'sys_config': {
                'back_config': ['audit', 'event', 'model', 'permission'],
                'global_busi': ['resource']
            },
            'model_config': {}
        }
        allClassify.map((classify) => {
            authority['model_config'][classify['bk_classification_id']] = {}
            classify['bk_objects'].map(model => {
                authority['model_config'][classify['bk_classification_id']][model['bk_obj_id']] = ['search', 'create', 'update', 'delete']
            })
        })
        state.adminAuthority = authority
    },
    setNavigation (state, navigation) {
        state.navigation = navigation
    },
    setTimezoneList (state, timezoneList) {
        state.timezoneList = timezoneList
    },
    setUsercustom (state, usercustom) {
        state.usercustom = usercustom
    },
    setGlobalLoading (state, isLoading) {
        state.globalLoading = isLoading
    }
}

export default {
    state,
    getters,
    actions,
    mutations
}
