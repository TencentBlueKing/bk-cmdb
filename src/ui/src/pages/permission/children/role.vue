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
    <div class="role">
        <div class="role-options clearfix">
            <div class="role-options-search fl clearfix">
                <div class="search-group fl">
                    <label for="searchGroupName">{{$t('Permission["角色名搜索"]')}}</label>
                    <input class="bk-form-input" type="text" id="searchGroupName" v-model.trim="filter['group_name']" @keyup.enter="getRoleList">
                </div>
                <div class="search-group fl">
                    <label for="SearchUserName">{{$t('Permission["成员搜索"]')}}</label>
                    <input class="bk-form-input" type="text" id="SearchUserName" v-model.trim="filter['user_list']" @keyup.enter="getRoleList">
                </div>
            </div>
            <div class="role-options-create fr">
                <bk-button type="primary" @click.prevent="createRole">
                    {{$t('Permission["新增角色"]')}}
                </bk-button>
            </div>
        </div>
        <v-role-table class="role-table"
            :sortable="false"
            :loading="isLoading"
            :header="table.header"
            :list="table.list"
            :wrapperMinusHeight="240">
            <template slot="operation" slot-scope="{ item }" style="text-align: center;font-size: 0">
                <span class="color-info" @click="skipToUser(item)">{{$t('Permission["跳转配置"]')}}</span>
                <span class="color-info" @click.stop="editRole(item)">{{$t('Common["编辑"]')}}</span>
                <span class="color-danger" @click.stop="confirmDeleteRole(item)">{{$t('Common["删除"]')}}</span>
            </template>        
        </v-role-table>
        <v-role-form 
            :data="form.data" 
            :isShow.sync="form.isShow" 
            :type="form.type"
            @on-success="handleCreateSuccess">
        </v-role-form>
        <bk-dialog 
            :has-header="false" 
            :is-show="deleteInfo.isShow" 
            :content="deleteInfo.content"
            @confirm="deleteRole(deleteInfo.role)"
            @cancel="cancelDeleteRole">
        </bk-dialog>
    </div>
</template>

<script>
    import vRoleForm from './roleForm'
    import vRoleTable from '@/components/table/table'
    import bus from '@/eventbus/bus'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            vRoleForm,
            vRoleTable
        },
        data () {
            return {
                isLoading: false,
                deleteInfo: {
                    isShow: false,
                    content: '',
                    role: null
                },
                filter: {
                    group_name: '',
                    user_list: ''
                },
                form: {
                    isShow: false,
                    type: 'create',
                    data: {
                        group_id: '',
                        group_name: '',
                        user_list: '',
                        PaasUserList: ''
                    }
                },
                table: {
                    header: [{
                        id: 'group_name',
                        name: this.$t('Permission["角色名"]')
                    }, {
                        id: 'user_list',
                        name: this.$t('Permission["角色成员"]')
                    }, {
                        id: 'operation',
                        name: this.$t('Permission["操作"]'),
                        attr: {
                            align: 'center'
                        }
                    }],
                    list: []
                }
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount',
                'language'
            ]),
            hasFilter () {
                return this.filter.group_name !== '' || this.filter.user_list !== ''
            }
        },
        watch: {
            'language' () {
                this.table.header = [{
                    id: 'group_name',
                    name: this.$t('Permission["角色名"]')
                }, {
                    id: 'user_list',
                    name: this.$t('Permission["角色成员"]')
                }, {
                    id: 'operation',
                    name: this.$t('Permission["操作"]'),
                    attr: {
                        align: 'center'
                    }
                }]
            }
        },
        methods: {
            skipToUser (item) {
                this.$emit('skipToUser', item)
            },
            handleCreateSuccess () {
                this.filter.group_name = ''
                this.filter.user_list = ''
                this.$nextTick(() => {
                    this.getRoleList()
                })
            },
            getRoleList () {
                this.isLoading = true
                this.$axios.post(`topo/privilege/group/${this.bkSupplierAccount}/search`, this.filter).then((res) => {
                    this.isLoading = false
                    if (res.result) {
                        this.table.list = res.data && res.data.length ? res.data : []
                        this.$emit('on-search-success', this.table.list, this.hasFilter)
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                }).catch(() => {
                    this.isLoading = false
                })
            },
            editRole (role) {
                this.form.data = Object.assign({}, role, {'user_list': role['user_list'].split(';').join(',')})
                this.form.type = 'edit'
                this.form.isShow = true
            },
            confirmDeleteRole (role) {
                this.deleteInfo.content = this.$tc('Permission["确认删除角色"]', role['group_name'], {name: role['group_name']})
                this.deleteInfo.role = role
                this.deleteInfo.isShow = true
            },
            cancelDeleteRole () {
                this.deleteInfo.isShow = false
            },
            deleteRole (role) {
                this.$axios.delete(`topo/privilege/group/${this.bkSupplierAccount}/${role['group_id']}`).then((res) => {
                    if (res.result) {
                        this.$alertMsg(this.$t('Permission["删除成功"]'), 'success')
                        this.getRoleList()
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
                this.deleteInfo.isShow = false
            },
            createRole () {
                this.form.data = {
                    group_id: '',
                    group_name: '',
                    user_list: '',
                    PaasUserList: ''
                }
                this.form.type = 'create'
                this.form.isShow = true
            }
        },
        mounted () {
            this.getRoleList()
        },
        created () {
            bus.$on('changePermissionTab', () => {
                this.createRole()
            })
        }
    }
</script>

<style lang="scss" scoped>
    $fontColor: #737987;
    .role-options{
        padding: 20px 0;
        font-size: 14px;
        .role-options-search{
            height: 36px;
            line-height: 36px;
        }
        .role-options-create{
            .btn-create{
                width: 124px;
                height: 36px;
                background-color: #6b7baa;
                border-radius: 2px;
                outline: 0;
                border: none;
                color: #fff;
            }
        }
    }
    .search-group{
        margin: 0 48px 0 0;
        font-size: 0;
        color: $fontColor;
        label{
            padding: 0 9px 0 0;
            font-size: 14px;
        }
        input{
            width: 210px;
            vertical-align: initial;
        }
    }
    .role-table{
        .color-info {
            margin-right: 10px;
            &:hover {
                color: #498fe0;
            }
        }
        .color-danger {
            &:hover {
                color: #ef4c4c;
            }
        }
    }
</style>