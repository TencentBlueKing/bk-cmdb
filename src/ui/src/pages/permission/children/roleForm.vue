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
    <bk-dialog 
    :is-show.sync="isShow" 
    :has-header="false" 
    :has-footer="false" 
    :quick-close="false" 
    :width="565" 
    :padding="0"
    @confirm="submitRoleForm"
    @cancel="closeRoleForm"
    >
        <form class="role-form" slot="content">
            <h2 class="role-form-title">{{title}}</h2>
            <div class="role-form-content">
                <div class="content-group clearfix">
                    <label for="groupName" class="fl">{{$t('Permission["角色名"]')}}</label>
                    <input type="text" class="bk-form-input fl" id="groupName" v-model.trim="data['group_name']" :disabled="isAdmin">
                </div>
                <div class="content-group clearfix" v-if="isAdmin">
                    <label for="paasUserList" class="fl">{{$t('Permission["与PaaS同步成员"]')}}</label>
                    <input type="text" class="fl" id="paasUserList" v-model="data.PaasUserList" disabled>
                    <a class="content-jump-link fl" href="javascript:void(0)"><i class="icon-cc-derivation"></i>{{$t('Permission["跳转配置"]')}}</a>
                </div>
                <div class="content-group clearfix">
                    <label for="userList" class="fl">{{$t('Permission["角色成员"]')}}</label>
                    <v-member-selector class="fl member-selector" :selected.sync="data['user_list']" :multiple="true"></v-member-selector>
                </div>
            </div>
            <div class="role-form-btn">
                <div class="fr">
                    <bk-button :loading="$loading('saveRole')" type="primary" class="form-btn" :disabled="!data['group_name'] || !data['user_list']" @click.prevent="submitRoleForm">{{$t('Common["确定"]')}}</bk-button>
                    <bk-button type="default" class="form-btn vice-btn" @click.prevent="closeRoleForm">{{$t('Common["取消"]')}}</bk-button>
                </div>
            </div>
        </form>
    </bk-dialog>
</template>
<script>
    import { mapGetters } from 'vuex'
    import vMemberSelector from '@/components/common/selector/member'
    export default {
        components: {
            vMemberSelector
        },
        props: {
            data: {
                type: Object,
                required: false,
                default () {
                    return {
                        group_name: '',
                        group_id: '',
                        supplier_account: '',
                        user_list: '',
                        PaasUserList: ''
                    }
                }
            },
            type: {
                type: String,
                required: true
            },
            isShow: {
                type: Boolean,
                require: true
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount'
            ]),
            title () {
                return this.type === 'create' ? this.$t('Permission["新增角色"]') : this.$t('Permission["编辑角色"]')
            },
            params () {
                let params = {
                    group_name: this.data['group_name'],
                    user_list: this.data['user_list'].split(',').join(';')
                }
                return params
            },
            isAdmin () {
                return false // this.data['group_name'] === 'admin'
            }
        },
        methods: {
            submitRoleForm () {
                if (this.type === 'create') {
                    this.createRole()
                } else {
                    this.updateRole()
                }
            },
            createRole () {
                this.$axios.post(`topo/privilege/group/${this.bkSupplierAccount}`, this.params, {id: 'saveRole'}).then((res) => {
                    if (res.result) {
                        this.closeRoleForm()
                        this.$alertMsg('新建角色成功', 'success')
                        this.$emit('on-success', res)
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                        this.$emit('on-error', res)
                    }
                })
            },
            updateRole () {
                this.$axios.put(`topo/privilege/group/${this.data['bk_supplier_account']}/${this.data['group_id']}`, this.params, {id: 'saveRole'}).then(res => {
                    if (res.result) {
                        this.closeRoleForm()
                        this.$alertMsg('更新角色成功', 'success')
                        this.$emit('on-success', res)
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                        this.$emit('on-error', res)
                    }
                })
            },
            closeRoleForm () {
                this.$emit('update:isShow', false)
            }
        }
    }
</script>
<style lang="scss" scoped>
    .role-form{
        padding: 50px 50px 0;
        .role-form-title{
            line-height: 30px;
            text-align: center;
            font-size: 22px;
            font-weight: normal;
            color: #333948;
        }
        .role-form-content{
            padding: 15px 0 0 0;
        }
        .role-form-btn{
            // padding: 40px 0 0 0;
            border-top: 1px solid #e5e5e5;
            padding-right: 20px;
            margin: 50px -50px 0;
            text-align: center;
            font-size: 0;
            background: #fafafa;
            height: 60px;
            line-height: 60px;
            button.form-btn{
                margin-left: 10px;
            }
        }
    }
    .content-group{
        line-height: 36px;
        margin: 20px 0 0 0;
        label{
            width: 105px;
            // color: #6b7baa;
            padding: 0 20px 0 5px;
            text-align: right;
            &:after{
                content: '*';
                color: #ef4c4c;
            }
        }
        input{
            width: 350px;
            height: 36px;
        }
        .content-jump-link{
            color: #498fe0;
            .icon-cc-derivation{
                margin: 0 4px 0 10px;
                vertical-align: -1px;
            }
        }
        .member-selector{
            width: 350px;
            line-height: initial;
        }
    }
</style>
<style lang="scss">
    .role-form-content{
        .content-group{
            .bk-data-wrapper{
                float: left;
                width: 290px;
                min-height: 36px;
                border-radius: 2px;
                padding: 0;
            }
            .bk-data-editor{
                width: 100%;
                height: 100%;
            }
            .bk-data-item{
                height: 24px;
                line-height: 24px;
            }
            [name='bk-data-input']{
                vertical-align: top;
                width: 100%;
                height: 34px;
                padding: 0 8px;
                border: none;
            }
        }
    }
</style>