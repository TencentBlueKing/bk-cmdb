<template>
    <div class="authority-wrapper">
        <template v-if="roles.length">
            <div class="authority-group clearfix">
                <h2 class="authority-group-title fl">{{$t('Permission["角色选择"]')}}</h2>
                <bk-selector class="role-selector fl"
                    :list="localRoles.list"
                    :selected.sync="localRoles.selected"
                    :searchable="true">
                </bk-selector>
            </div>
            <div class="authority-group clearfix">
                <h2 class="authority-group-title fl">{{$t('Permission["系统相关"]')}}</h2>
                <div class="authority-group-content">
                    <div class="authority-type system clearfix" 
                        v-for="(config, configId) in sysConfig" 
                        v-if="config.authorities.length">
                        <h3 class="system-title fl">{{$t(config.name)}}</h3>
                        <ul class="system-list fl">
                            <li class="system-item fl"  v-for="authority in config.authorities">
                                <label class="cmdb-form-checkbox cmdb-checkbox-small"
                                    :for="'systemAuth-' + authority.id" 
                                    :title="$t(authority.name)">
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
            <div class="authority-group model clearfix">
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
                                        <label class="cmdb-form-checkbox cmdb-checkbox-small"
                                            :for="'model-all-'+model['bk_obj_id']">
                                            <input type="checkbox"
                                                :id="'model-all-'+model['bk_obj_id']" 
                                                :checked="model.selectedAuthorities.length === 3"
                                                @change="checkAllModelAuthorities(classifyIndex,modelIndex,$event)">{{$t('Common["全选"]')}}
                                        </label>
                                    </span>
                                    <span class="model-authority-checkbox fl">
                                        <label class="cmdb-form-checkbox cmdb-checkbox-small"
                                            :for="'model-search-'+model['bk_obj_id']">
                                            <input type="checkbox" value='search' 
                                                :id="'model-search-'+model['bk_obj_id']" 
                                                v-model="model.selectedAuthorities"
                                                @change="checkOtherAuthorities(classifyIndex,modelIndex,$event)">{{$t('Common["查询"]')}}
                                        </label>
                                    </span>
                                    <span class="model-authority-checkbox fl">
                                        <label class="cmdb-form-checkbox cmdb-checkbox-small" 
                                            :for="'model-update-'+model['bk_obj_id']" 
                                            :class="{'disabled': model.selectedAuthorities.indexOf('search') === -1}">
                                            <input type="checkbox" value='update' 
                                                :id="'model-update-'+model['bk_obj_id']"
                                                :disabled="model.selectedAuthorities.indexOf('search') === -1"  
                                                v-model="model.selectedAuthorities">{{$t('Common["编辑"]')}}
                                        </label>
                                    </span>
                                    <span class="model-authority-checkbox fl">
                                        <label class="cmdb-form-checkbox cmdb-checkbox-small" 
                                            :for="'model-delete-'+model['bk_obj_id']" 
                                            :class="{'disabled': model.selectedAuthorities.indexOf('search') === -1}">
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
            <footer class="footer">
                <bk-button type="primary" :loading="$loading('updateGroupAuthorities')" @click="updateGroupAuthorities">
                    {{$t('Common["保存"]')}}
                </bk-button>
            </footer>
        </template>
        <template v-else>
            <div class="user-none">
                <img src="../../assets/images/user-none.png" :alt="$t('Permission[\'没有创建角色\']')">
                <p class="mt20">
                    {{$t('Permission["没有创建角色"]')}}
                    <span class="btn" @click="changeTab">{{$t('Permission["点击新增"]')}}</span>
                </p>
            </div>
        </template>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    export default {
        props: {
            groupId: {
                type: String
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
                groupAuthorities: {},
                hideClassify: ['bk_host_manage', 'bk_biz_topo']
            }
        },
        computed: {
            ...mapGetters('userPrivilege', [
                'roles'
            ]),
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
            'localRoles.selected' (groupId) {
                if (groupId) {
                    this.getGroupAuthorities(groupId)
                }
            }
        },
        methods: {
            ...mapActions('userPrivilege', [
                'searchUserPrivilege',
                'updateGroupPrivilege'
            ]),
            changeTab () {
                this.$emit('createRole')
            },
            updateGroupAuthorities () {
                this.updateGroupPrivilege({bkGroupId: this.localRoles.selected, params: this.updateParams, config: {requestId: 'updateGroupAuthorities'}})
            },
            checkAllModelAuthorities (classifyIndex, modelIndex, event) {
                let model = this.classifications[classifyIndex]['models'][modelIndex]
                if (event.target.checked) {
                    model.selectedAuthorities = ['search', 'update', 'delete']
                } else {
                    model.selectedAuthorities = []
                }
            },
            async getGroupAuthorities (groupId) {
                const res = await this.searchUserPrivilege({bkGroupId: groupId})
                this.groupAuthorities = res.privilege
                this.initSystemAuthorities()
                this.initClassifications()
            },
            calcModelListStyle (total) {
                return {
                    height: `${total * 42 - 10}px`
                }
            },
            /* 模型权限没有选择'查询'，则无'新增'、编辑'、删除'权限 */
            checkOtherAuthorities (classifyIndex, modelIndex, event) {
                let model = this.classifications[classifyIndex]['models'][modelIndex]
                if (!event.target.checked) {
                    model.selectedAuthorities = []
                }
            },
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
            initClassifications () {
                let classifications = []
                // 1.去掉停用模型
                let activeClassifications = this.$classifications.map(classification => {
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
                    return !this.hideClassify.includes(bkClassificationId) && Array.isArray(bkObjects) && bkObjects.length
                })
                let authorities = this.groupAuthorities
                activeClassifications.map(classify => {
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
            initRoles () {
                this.localRoles.list = this.roles.map(role => {
                    return Object.assign(role, {
                        id: role['group_id'],
                        name: role['group_name']
                    })
                })
                if (this.localRoles.list.length) {
                    this.localRoles.selected = this.groupId === '' ? this.localRoles.list[0].id : this.groupId
                } else {
                    this.localRoles.selected = ''
                }
            }
        },
        created () {
            this.initRoles()
        }
    }
</script>

<style lang="scss" scoped>
    .authority-wrapper{
        padding: 20px 0;
        height: 100%;
    }
    .authority-group{
        font-size: 14px;
        margin: 30px 0 0 0;
        &.model{
            margin-top: 14px;
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
    label.cmdb-form-checkbox{
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
    .footer{
        width: 100%;
        position: absolute;
        left: 0;
        bottom: 0;
        padding: 14px 20px;
        background: #f9f9f9;
        .bk-button{
            height: 36px;
            line-height: 34px;
            border-radius: 2px;
            display: inline-block;
            min-width: 110px;
            vertical-align: bottom;
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
        top: 40%;
        left: 50%;
        transform: translate(-50%, -50%);
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
