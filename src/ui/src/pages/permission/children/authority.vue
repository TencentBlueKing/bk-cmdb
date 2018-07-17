/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template>
    <div class="authority-wrapper">
        <template v-if="roles.length">
            <div class="authority-group clearfix">
                <h2 class="authority-group-title fl">{{$t('Permission["角色选择"]')}}</h2>
                <bk-select class="role-selector fl" :selected.sync="localRoles.selected">
                    <bk-select-option v-for="(role, index) in localRoles.list"
                        :key="index"
                        :value="role['group_id']"
                        :label="role['group_name']"
                    ></bk-select-option>
                </bk-select>
            </div>
            <div class="authority-group clearfix">
                <h2 class="authority-group-title fl">{{$t('Permission["系统相关"]')}}</h2>
                <div class="authority-group-content">
                    <div class="authority-type system clearfix" 
                        v-for="(config, configId) in sysConfig" 
                        v-if="config.authorities.length">
                        <h3 class="system-title fl">{{$t(config.name)}}</h3>
                        <ul class="system-list clearfix">
                            <li class="system-item fl"  v-for="authority in config.authorities">
                                <label class="bk-form-checkbox bk-checkbox-small"
                                    :for="'systemAuth-' + authority.id" 
                                    :title="$t(authority.name)"
                                    @click="updateGroupAuthorities">
                                    <input type="checkbox"
                                        :id="'systemAuth-' + authority.id" 
                                        :value="authority.id"
                                        v-model="config.selectedAuthorities">{{$t(authority.name)}}
                                </label>
                            </li>
                        </ul>
                    </div>
                </div>
            </div>
            <div class="authority-group model clearfix" style="margin-top:14px;">
                <h2 class="authority-group-title"><span>{{$t('Permission["模型相关"]')}}</span></h2>
                <div class="authority-group-content">
                    <div class="authority-type model" v-for="(classify,classifyIndex) in classifications" v-if="classify.models.length"> 
                        <h3 class="classify-name clearfix" :title="classify.name" @click="classify.open = !classify.open">
                            <span class="fl">{{classify.name}}</span>
                            <i class="bk-icon icon-angle-down angle fr" :class="{'open': classify.open}"></i>
                        </h3>
                        <transition name="slide">
                            <ul class="model-list" v-show="classify.open" :style="calcModelListStyle(classify.models.length)">
                                <li class="model-item clearfix" v-for="(model,modelIndex) in classify.models">
                                    <h4 class="model-authority fl" :title="model['bk_obj_name']">{{model['bk_obj_name']}}</h4>
                                    <span class="model-authority-checkbox fl">
                                        <label class="bk-form-checkbox bk-checkbox-small"
                                            :for="'model-all-'+model['bk_obj_id']" 
                                            @click="updateGroupAuthorities">
                                            <input type="checkbox"
                                                :id="'model-all-'+model['bk_obj_id']" 
                                                :checked="model.selectedAuthorities.length === 3"
                                                @change="checkAllModelAuthorities(classifyIndex,modelIndex,$event)">{{$t('Common["全选"]')}}
                                        </label>
                                    </span>
                                    <span class="model-authority-checkbox fl">
                                        <label class="bk-form-checkbox bk-checkbox-small"
                                            :for="'model-search-'+model['bk_obj_id']"
                                            @click="updateGroupAuthorities">
                                            <input type="checkbox" value='search' 
                                                :id="'model-search-'+model['bk_obj_id']" 
                                                v-model="model.selectedAuthorities"
                                                @change="checkOtherAuthorities(classifyIndex,modelIndex,$event)">{{$t('Common["查询"]')}}
                                        </label>
                                    </span>
                                    <span class="model-authority-checkbox fl">
                                        <label class="bk-form-checkbox bk-checkbox-small" 
                                            :for="'model-update-'+model['bk_obj_id']" 
                                            :class="{'disabled': model.selectedAuthorities.indexOf('search') === -1}"
                                            @click="updateGroupAuthorities">
                                            <input type="checkbox" value='update' 
                                                :id="'model-update-'+model['bk_obj_id']"
                                                :disabled="model.selectedAuthorities.indexOf('search') === -1"  
                                                v-model="model.selectedAuthorities">{{$t('Common["编辑"]')}}
                                        </label>
                                    </span>
                                    <span class="model-authority-checkbox fl">
                                        <label class="bk-form-checkbox bk-checkbox-small" 
                                            :for="'model-delete-'+model['bk_obj_id']" 
                                            :class="{'disabled': model.selectedAuthorities.indexOf('search') === -1}"
                                            @click="updateGroupAuthorities">
                                            <input type="checkbox" value='delete' 
                                                :id="'model-delete-'+model['bk_obj_id']"
                                                :disabled="model.selectedAuthorities.indexOf('search') === -1" 
                                                v-model="model.selectedAuthorities">{{$t('Common["删除"]')}}
                                        </label>
                                    </span>
                                </li>
                            </ul>
                        </transition>
                    </div>
                </div>
            </div>
        </template>
        <template v-else>
            <div class="user-none">
                <img src="../../../common/images/user-none.png" :alt="$t('Permission[\'没有创建角色\']')">
                <p>
                    {{$t('Permission["没有创建角色"]')}}
                    <span class="btn" @click="changeTab">{{$t('Permission["点击新增"]')}}</span>
                </p>
            </div>
        </template>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import Throttle from 'lodash.throttle'
    import bus from '@/eventbus/bus'
    export default {
        props: {
            roles: {
                type: Array,
                required: true
            },
            activeGroup: {
                type: Object
            }
        },
        data () {
            return {
                localRoles: {
                    list: [],
                    selected: ''
                },
                sysConfig: {
                    global_busi: {
                        id: 'global_busi',
                        name: 'Permission["全局业务"]',
                        authorities: [{
                            id: 'resource',
                            name: 'Permission["资源池管理"]'
                        }],
                        selectedAuthorities: []
                    },
                    back_config: {
                        id: 'back_config',
                        name: 'Permission["后台配置"]',
                        authorities: [{
                            id: 'event',
                            name: 'Permission["事件推送配置"]'
                        }, {
                            id: 'audit',
                            name: 'OperationAudit["操作审计"]'
                        }],
                        selectedAuthorities: []
                    }
                },
                classifications: [],
                groupAuthorities: null,
                hideClassify: ['bk_host_manage', 'bk_biz_topo']
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount'
            ]),
            ...mapGetters('navigation', ['activeClassifications']),
            updateParams () {
                let updateParams = {}
                for (let config in this.sysConfig) {
                    if (this.sysConfig[config].selectedAuthorities.length) {
                        updateParams.sys_config = updateParams.sys_config || {}
                        updateParams.sys_config[config] = this.sysConfig[config].selectedAuthorities
                    }
                }
                this.classifications.map((classify) => {
                    classify.models.map((model) => {
                        if (model.selectedAuthorities.length) {
                            updateParams['model_config'] = updateParams['model_config'] || {}
                            updateParams['model_config'][classify.id] = updateParams['model_config'][classify.id] || {}
                            updateParams['model_config'][classify.id][model['bk_obj_id']] = model.selectedAuthorities
                        }
                    })
                })
                return updateParams
            }
        },
        watch: {
            activeGroup (group) {
                if (JSON.stringify(group) !== '{}') {
                    this.localRoles.selected = group['group_id']
                }
            },
            roles (roles) {
                this.localRoles.list = roles.slice()
                if (this.localRoles.list.length) {
                    // 如果之前已经选中了角色，需判断该选中的角色是否被删除
                    if (this.localRoles.selected) {
                        let isRoleDeleted = true
                        this.localRoles.list.map((role) => {
                            if (this.localRoles.selected === role['group_id']) {
                                isRoleDeleted = false
                            }
                        })
                        if (isRoleDeleted) {
                            this.localRoles.selected = this.localRoles.list[0]['group_id']
                        }
                    } else {
                        this.localRoles.selected = this.localRoles.list[0]['group_id']
                    }
                } else {
                    this.localRoles.selected = ''
                }
            },
            activeClassifications () {
                // 查询分组权限接口先返回数据
                // 获取到模型后要做一次初始化
                if (this.groupAuthorities) {
                    this.initClassifications()
                }
            },
            'localRoles.selected' (groupID) {
                if (groupID) {
                    this.getGroupAuthorities(groupID)
                }
            }
        },
        methods: {
            changeTab () {
                this.$emit('update:activeTabName', 'role')
                bus.$emit('changePermissionTab')
            },
            getGroupAuthorities (groupID) {
                this.$axios.get(`topo/privilege/group/detail/${this.bkSupplierAccount}/${groupID}`).then((res) => {
                    if (res.result) {
                        this.groupAuthorities = res.data.privilege
                        this.initSystemAuthorities()
                        this.initClassifications()
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            // 勾选后直接发起请求，并做函数节流
            updateGroupAuthorities: Throttle(function () {
                this.$nextTick(() => {
                    this.$axios.post(`topo/privilege/group/detail/${this.bkSupplierAccount}/${this.localRoles.selected}`, this.updateParams)
                    .then((res) => {
                        if (!res.result) {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                })
            }, 500, {leading: false, trailing: true}),
            // 获取到分组权限后设置系统权限
            initSystemAuthorities () {
                if (this.groupAuthorities.hasOwnProperty('sys_config')) {
                    for (let configId in this.sysConfig) {
                        if (this.groupAuthorities['sys_config'].hasOwnProperty(configId)) {
                            this.sysConfig[configId].selectedAuthorities = this.groupAuthorities['sys_config'][configId] || []
                        } else {
                            this.sysConfig[configId].selectedAuthorities = []
                        }
                    }
                } else {
                    for (let configId in this.sysConfig) {
                        this.sysConfig[configId].selectedAuthorities = []
                    }
                }
            },
            // 获取到分组权限后设置模型权限
            initClassifications () {
                let classifications = []
                let authorities = this.groupAuthorities
                this.activeClassifications.forEach((classify) => {
                    let models = []
                    let classifyId = classify['bk_classification_id']
                    if (this.hideClassify.indexOf(classifyId) === -1) {
                        classify['bk_objects'].forEach((model) => {
                            let selectedAuthorities = []
                            if (authorities.hasOwnProperty('model_config') &&
                                authorities['model_config'].hasOwnProperty(classifyId) &&
                                authorities['model_config'][classifyId].hasOwnProperty(model['bk_obj_id'])
                            ) {
                                selectedAuthorities = authorities['model_config'][classifyId][model['bk_obj_id']]
                            }
                            models.push(Object.assign({}, model, {selectedAuthorities}))
                        })
                        classifications.push({
                            id: classify['bk_classification_id'],
                            name: classify['bk_classification_name'],
                            open: true,
                            models: models
                        })
                    }
                })
                this.classifications = classifications
            },
            // 模型全选
            checkAllModelAuthorities (classifyIndex, modelIndex, event) {
                let model = this.classifications[classifyIndex]['models'][modelIndex]
                if (event.target.checked) {
                    model.selectedAuthorities = ['search', 'update', 'delete']
                } else {
                    model.selectedAuthorities = []
                }
            },
            /* 模型权限没有选择'查询'，则无'新增'、编辑'、删除'权限 */
            checkOtherAuthorities (classifyIndex, modelIndex, event) {
                let model = this.classifications[classifyIndex]['models'][modelIndex]
                if (!event.target.checked) {
                    model.selectedAuthorities = []
                }
            },
            calcModelListStyle (total) {
                return {
                    height: `${total * 42 - 10}px`
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    $primaryColor: #737987;
    .authority-wrapper{
        position: relative;
        padding: 20px 0;
        height: 100%;
    }
    .authority-group{
        font-size: 14px;
        margin: 30px 0 0 0;
        &.model{
            .authority-group-title{
                width: auto;
                height: 44px;
                line-height: 44px;
                background-color: #f9f9f9;
                span{
                    display: block;
                    width: 107px;
                    text-align: right; 
                }
            }
            .authority-group-content{
                padding: 0;
            }
        } 
        .authority-group-title{
            width: 137px;
            height: 36px;
            margin: 0;
            line-height: 36px;
            font-weight: bold;
            font-size: 14px;
            color: $primaryColor;
            padding: 0 30px 0 0;
            text-align: right;
        }
        .authority-group-content{
            overflow: visible;
            padding: 44px 0 0 137px;
        }
        .role-selector{
            width: 286px;
        }
    }
    .authority-type.system{
        line-height: 32px;
        &:first-child{
            margin-top: 0;
        }
        .system-title{
            width: 100px;
            margin: 0;
            font-size: 14px;
            font-weight: normal;
            color: #498fe0;
        }
        .system-list{
            padding: 0 0 0 100px;
            white-space: nowrap;
            .system-item{
                width: 150px;
                height: 32px;
                margin: 0 0 10px 16px;
            }
        }
    }
    .authority-type.model{
        line-height: 32px;
        padding: 10px 0;
        border-bottom: 1px solid #eceef5;
        .classify-name{
            font-size: 12px;
            color: $primaryColor;
            cursor: pointer;
            margin: 0;
            span{
                display: block;
                width: 107px;
                text-align: right;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: nowrap;
            }
            .icon-angle-down{
                font-size: 14px;
                margin: 9px 30px 0 0;
                color: $primaryColor;
                transform: rotate(180deg);
                transition: transform .5s cubic-bezier(.23, 1, .23, 1);
                &.open{
                    transform: rotate(0);
                }
            }
        }
        .model-list{
            .model-item{
                padding: 10px 0 0 138px;
                &:first-child{
                    padding-top: 0;
                }
            }
        }
        .model-authority{
            width: 100px;
            margin: 0;
            font-size: 14px;
            font-weight: normal;
            color: #498fe0;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
        }
        .model-authority-checkbox{
            width: 150px;
            height: 32px;
            margin: 0 0 0 16px;
            &:last-child{
                width: auto;
            }
        }
    }
    label.bk-form-checkbox{
        max-width: 130px;
        line-height: 32px;
        padding: 0;
        margin: 0;
        overflow: hidden;
        text-overflow: ellipsis;
        cursor: pointer;
        border: none;
        &.disabled{
            cursor: not-allowed;
            color: #c3cdd7;
        }
        input[type='checkbox']{
            width: 14px;
            height: 14px;
            margin-right: 10px;
            vertical-align: -2px;
        }
    }
    .slide-enter-active, .slide-leave-active{
        transition: height .5s cubic-bezier(.23, 1, .23, 1);
        overflow: hidden;
    }
    .slide-enter, .slide-leave-to{
        height: 0 !important;
    }
    .user-none{
        position: absolute;
        top: 35%;
        left: 50%;
        transform: translate(-50%, -50%);
        color: $primaryColor;
        margin: 0 auto;
        text-align: center;
        img{
            width: 180px;
        }
        span{
            color: #3c96ff;
            cursor: pointer;
        }
    }
</style>