import {$axios, $Axios, $alertMsg} from '@/api/axios'
import STATIC_NAVIGATION from '@/common/json/static_navigation.json'
const state = {
    fold: window.localStorage.getItem('navStick') === 'false',
    navStick: window.localStorage.getItem('navStick') !== 'false',
    historyCount: 0,
    classifications: [],
    invisibleClassifications: ['bk_host_manage', 'bk_biz_topo'],
    notCustomClassifications: ['bk_index', 'bk_host_manage', 'bk_back_config'],
    authority: {
        model_config: {},
        sys_config: {
            global_busi: null,
            back_config: null
        }
    },
    interceptStaticModel: {
        'bk_host_manage': ['resource'],
        'bk_back_config': ['event', 'model', 'audit']
    },
    result: {
        classification: false,
        authority: false
    }
}

const getters = {
    fold: state => state.fold,
    navStick: state => state.navStick,
    historyCount: state => state.historyCount,
    classifications: state => state.classifications,
    result: state => state.result,
    authority: state => state.authority,
    // 可用分类
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
        let modelAuthority = state.authority['model_config'] || {}
        let authorizedClassifications = JSON.parse(JSON.stringify(getters.activeClassifications))
        if (!rootGetters.isAdmin) {
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
        let usercustomNavigation = rootGetters['usercustom/usercustom']['navigation'] || {}
        let authorizedClassifications = JSON.parse(JSON.stringify(getters.authorizedClassifications))
        // 构造模型导航数据
        let navigation = authorizedClassifications.map(classification => {
            return {
                'icon': classification['bk_classification_icon'],
                'id': classification['bk_classification_id'],
                'name': classification['bk_classification_name'],
                'children': classification['bk_objects'].map(model => {
                    return {
                        'path': `/organization/${model['bk_obj_id']}`,
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
        if (!rootGetters.isAdmin) {
            let sysConfig = {
                'bk_host_manage': state.authority['sys_config']['global_busi'] || [],
                'bk_back_config': state.authority['sys_config']['back_config'] || []
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
    },
    // 用户自定义导航
    customNavigation: (state, getters, rootState, rootGetters) => {
        let navigation = JSON.parse(JSON.stringify(getters.authorizedNavigation))
        let usercustomNavigation = rootGetters['usercustom/usercustom']['navigation'] || {}
        navigation = navigation.filter(({id}) => {
            if (state.notCustomClassifications.includes(id)) {
                return true
            }
            return usercustomNavigation.hasOwnProperty(id) && usercustomNavigation[id].length
        })
        navigation.forEach(({id}, index) => {
            if (!state.notCustomClassifications.includes(id)) {
                navigation[index]['children'] = navigation[index]['children'].filter(model => usercustomNavigation[id].includes(model.id))
            }
        })
        navigation = navigation.filter(({id, children}) => {
            if (['bk_index', 'bk_host_manage'].includes(id)) {
                return true
            }
            return children.length
        })
        return navigation
    }
}
const actions = {
    async getClassifications ({commit, state, rootState}, forceUpdate = false) {
        if (state.result.classification && !forceUpdate) {
            return Promise.resolve({result: true, data: state.classifications})
        }
        let classifications = []
        await $axios.post(`object/classification/${rootState.common.bkSupplierAccount}/objects`, {
        }).then(res => {
            if (res.result) {
                classifications = res.data
                commit('setClassifications', classifications)
            } else {
                $alertMsg(res['bk_error_msg'])
            }
        })
        return Promise.resolve({result: state.result.classification, data: state.classifications})
    },
    getAuthority ({commit, state, rootState}) {
        if (state.result.authority) {
            return Promise.resolve({result: true, data: state.authority})
        }
        return $axios.get(`topo/privilege/user/detail/${rootState.common.bkSupplierAccount}/${window.userName}`).then(res => {
            if (res.result) {
                state.authority = res.data
            } else {
                $alertMsg(res['bk_error_msg'])
            }
            state.result.authority = res.result
            return res
        })
    }
}

const mutations = {
    updateNavFold (state, fold) {
        state.fold = fold
    },
    updateNavStick (state, stick) {
        state.navStick = !!stick
        window.localStorage.setItem('navStick', !!stick)
    },
    updateHistoryCount (state, step = 1) {
        const count = state.historyCount + step
        state.historyCount = count > 0 ? count : 0
    },
    setClassifications (state, classifications) {
        state.result.classification = true
        state.classifications = classifications
    },
    updateClassification (state, classification) {
        if (classification.hasOwnProperty('bk_classification_id')) {
            let targetClassification = state.classifications.find(({bk_classification_id: bkClassificationId}) => {
                return bkClassificationId === classification['bk_classification_id']
            })
            if (targetClassification) {
                Object.assign(targetClassification, classification)
            }
        }
    },
    updateModel (state, model) {
        if (model.hasOwnProperty('bk_classification_id') && model.hasOwnProperty('bk_obj_id')) {
            let targetClassification = state.classifications.find(({bk_classification_id: bkClassificationId}) => {
                return bkClassificationId === model['bk_classification_id']
            })
            if (targetClassification) {
                let targetModel = targetClassification['bk_objects'].find(({bk_obj_id: bkObjId}) => {
                    return bkObjId === model['bk_obj_id']
                })
                if (targetModel) {
                    Object.assign(targetModel, model)
                }
            }
        }
    }
}

export default {
    namespaced: true,
    state,
    getters,
    actions,
    mutations
}
