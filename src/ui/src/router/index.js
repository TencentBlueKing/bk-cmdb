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
import {mapMutations, mapGetters} from 'vuex'
import store from '@/store'
import Router from 'vue-router'

Vue.use(Router)

const pageIndex = () => import(/* webpackChunkName: "page-index" */ '@/pages/index/index-v3')
const pageHosts = () => import(/* webpackChunkName: "page-hosts" */ '@/pages/hosts/hosts')
const pageModel = () => import(/* webpackChunkName: "page-model" */ '@/pages/model/model')
const pageResource = () => import(/* webpackChunkName: "page-resource" */ '@/pages/resource/resource')
const pageProcess = () => import(/* webpackChunkName: "page-process" */ '@/pages/process/process')
const pagePermission = () => import(/* webpackChunkName: "page-permission" */ '@/pages/permission/permission')
const pageEventpush = () => import(/* webpackChunkName: "page-eventpush" */ '@/pages/eventpush/eventpush')
const pageAuditing = () => import(/* webpackChunkName: "page-auditing" */ '@/pages/auditing/auditing')
const pageOrganization = () => import(/* webpackChunkName: "page-organization" */ '@/pages/organization/object')
const pageTopology = () => import(/* webpackChunkName: "page-topology" */ '@/pages/topology/topology')
const pageCustomQuery = () => import(/* webpackChunkName: "page-customQuery" */ '@/pages/customQuery/customQuery')

var routerVue = new Vue({
    store: store,
    computed: {
        ...mapGetters('navigation', ['authorizedNavigation'])
    },
    methods: {
        ...mapMutations(['setGlobalLoading']),
        ...mapMutations('navigation', ['updateHistoryCount']),
        async isAuthorized (to) {
            await this.$store.dispatch('navigation/getAuthority')
            await Promise.all([this.$store.dispatch('navigation/getClassifications'), this.$store.dispatch('usercustom/getUserCustom')])
            let isAuthorized = false
            let authorizedPath = ['/index', '/403', '/404']
            if (authorizedPath.includes(to.path)) {
                isAuthorized = true
            } else {
                isAuthorized = this.authorizedNavigation.some(({id, children}) => {
                    return children.some(({path}) => path === to.path)
                })
            }
            return Promise.resolve(isAuthorized)
        }
    }
})

var router = new Router({
    linkActiveClass: 'active',
    routes: [{
        path: '/404',
        components: require('@/pages/404')
    }, {
        path: '/403',
        components: require('@/pages/403')
    }, {
        path: '/',
        redirect: '/index'
    }, {
        path: '/index',
        component: pageIndex
    }, {
        path: '/hosts',
        component: pageHosts,
        meta: {
            setBkBizId: true
        }
    }, {
        path: '/model',
        component: pageModel
    }, {
        path: '/resource',
        component: pageResource
    }, {
        path: '/process',
        component: pageProcess,
        meta: {
            setBkBizId: true
        }
    }, {
        path: '/permission',
        component: pagePermission
    }, {
        path: '/eventpush',
        component: pageEventpush
    }, {
        path: '/auditing',
        component: pageAuditing
    }, {
        path: '/organization/:objId',
        component: pageOrganization
    }, {
        path: '/topology',
        component: pageTopology,
        meta: {
            setBkBizId: true
        }
    }, {
        path: '/customQuery',
        component: pageCustomQuery,
        meta: {
            setBkBizId: true
        }
    }, {
        path: '*',
        redirect: '/404'
    }]
})

router.beforeEach(async (to, from, next) => {
    routerVue.setGlobalLoading(true)
    routerVue.updateHistoryCount(-1)
    let isAuthorized = await routerVue.isAuthorized(to)
    if (isAuthorized) {
        if (!to.matched.some(({meta}) => meta.setBkBizId)) {
            delete routerVue.$axios.defaults.headers.bk_biz_id
        }
        next()
    } else {
        next('/403')
    }
})
router.afterEach((to, from) => {
    routerVue.setGlobalLoading(false)
})
export default router
