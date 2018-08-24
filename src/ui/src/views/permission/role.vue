<template>
    <div class="role-wrapper">
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
                <bk-button type="primary" @click="createRole">
                    {{$t('Permission["新增角色"]')}}
                </bk-button>
            </div>
        </div>
        <cmdb-table
            class="role-table"
            :sortable="false"
            :loading="$loading('searchUserGroup')"
            :header="table.header"
            :list="table.list"
            :wrapperMinusHeight="240">
            <template slot="operation" slot-scope="{ item }">
                <span class="text-primary" @click="skipToUser(item)">{{$t('Permission["跳转配置"]')}}</span>
                <span class="text-primary" @click.stop="editRole(item)">{{$t('Common["编辑"]')}}</span>
                <span class="text-danger" @click.stop="confirmDeleteRole(item)">{{$t('Common["删除"]')}}</span>
            </template>  
        </cmdb-table>
        <v-role-form 
            v-if="form.isShow"
            :data="form.data" 
            :type="form.type"
            @on-success="handleCreateSuccess"
            @closeRoleForm="form.isShow = false">
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
    import vRoleForm from './role-form'
    import { mapActions, mapMutations } from 'vuex'
    export default {
        data () {
            return {
                filter: {
                    group_name: '',
                    user_list: ''
                },
                deleteInfo: {
                    isShow: false,
                    content: '',
                    role: null
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
                }
            }
        },
        methods: {
            ...mapActions('userPrivilege', [
                'searchUserGroup',
                'deleteUserGroup'
            ]),
            ...mapMutations('userPrivilege', [
                'setRoles'
            ]),
            skipToUser (item) {
                this.$emit('skipToUser', item['group_id'])
            },
            handleCreateSuccess () {
                this.filter.group_name = ''
                this.filter.user_list = ''
                this.$nextTick(() => {
                    this.getRoleList()
                })
            },
            confirmDeleteRole (role) {
                this.deleteInfo.content = this.$tc('Permission["确认删除角色"]', role['group_name'], {name: role['group_name']})
                this.deleteInfo.role = role
                this.deleteInfo.isShow = true
            },
            async deleteRole (role) {
                await this.deleteUserGroup({bkGroupId: role['group_id']})
                this.$success(this.$t('Permission["删除成功"]'))
                this.getRoleList()
            },
            cancelDeleteRole () {
                this.deleteInfo.isShow = false
            },
            editRole (role) {
                this.form.data = Object.assign({}, role, {'user_list': role['user_list'].split(';').join(',')})
                this.form.type = 'edit'
                this.form.isShow = true
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
            },
            async getRoleList () {
                const res = await this.searchUserGroup({params: this.filter, config: {requestId: 'searchUserGroup'}})
                this.table.list = res && res.length ? res : []
                this.setRoles(this.table.list)
            }
        },
        created () {
            this.getRoleList()
        },
        components: {
            vRoleForm
        }
    }
</script>

<style lang="scss" scoped>
    .role-options{
        padding: 20px 0;
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
        .search-group{
            margin: 0 48px 0 0;
            font-size: 0;
            label{
                padding: 0 9px 0 0;
                font-size: 14px;
            }
            input{
                width: 210px;
                vertical-align: initial;
            }
        }
    }
    .role-table{
        .text-primary {
            margin-right: 10px;
        }
    }
</style>
