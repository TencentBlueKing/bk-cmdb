<template>
    <div class="role-wrapper">
        <div class="role-options clearfix">
            <div class="role-options-create fl">
                <bk-button type="primary" @click="createRole">
                    {{$t('Permission["新增角色"]')}}
                </bk-button>
            </div>
            <div class="role-options-search fr clearfix">
                <bk-selector
                    class="search-selector"
                    :list="typeList"
                    :selected.sync="filter.type"
                ></bk-selector>
                <input class="cmdb-form-input" :placeholder="$t('Common[\'请输入\']')" type="text" id="SearchUserName" v-model.trim="filter.text" @keyup.enter="getRoleList">
                <i class="filter-search bk-icon icon-search"
                @click="getRoleList"></i>
            </div>
        </div>
        <cmdb-table
            class="role-table"
            rowCursor="default"
            :sortable="false"
            :loading="$loading('searchUserGroup')"
            :header="table.header"
            :list="table.list"
            :wrapperMinusHeight="240">
            <template slot="operation" slot-scope="{ item }">
                <span class="text-primary" @click="showDetails(item)">{{$t('Permission["权限详情"]')}}</span>
                <span class="text-primary" @click.stop="editRole(item)">{{$t('Common["编辑"]')}}</span>
                <span class="text-danger" @click.stop="confirmDeleteRole(item)">{{$t('Common["删除"]')}}</span>
            </template>
            <div class="empty-info" slot="data-empty">
                <p>{{$t("Common['暂时没有数据']")}}</p>
                <p>{{$t("Permission['当前并无角色，可点击下方按钮新增']")}}</p>
                <bk-button class="process-btn" type="primary" @click="createRole">{{$t("Permission['新增角色']")}}</bk-button>
            </div>
        </cmdb-table>
        <v-role-form 
            ref="roleForm"
            v-if="form.isShow"
            :data="form.data" 
            :type="form.type"
            @on-success="handleCreateSuccess"
            @closeRoleForm="form.isShow = false">
        </v-role-form>
        <cmdb-slider
        :width="600"
        :title="slider.title"
        :isShow.sync="slider.isShow">
            <vAuthority
                slot="content"
                v-if="slider.isShow"
                :groupId="slider.groupId"
                @cancel="slider.isShow = false"
            ></vAuthority>
        </cmdb-slider>
    </div>
</template>

<script>
    import vRoleForm from './role-form'
    import vAuthority from './authority'
    import { mapActions, mapMutations } from 'vuex'
    export default {
        components: {
            vRoleForm,
            vAuthority
        },
        data () {
            return {
                filter: {
                    type: 'group_name',
                    text: ''
                },
                typeList: [{
                    id: 'group_name',
                    name: this.$t('Permission["角色名"]')
                }, {
                    id: 'user_list',
                    name: this.$t('Permission["角色成员"]')
                }],
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
                },
                slider: {
                    isShow: false,
                    groupId: '',
                    title: ''
                }
            }
        },
        created () {
            this.getRoleList()
        },
        methods: {
            ...mapActions('userPrivilege', [
                'searchUserGroup',
                'deleteUserGroup'
            ]),
            ...mapMutations('userPrivilege', [
                'setRoles'
            ]),
            showDetails (item) {
                this.slider.groupId = item['group_id']
                this.slider.title = `${item['group_name']} ${this.$t('Permission["权限详情"]')}`
                this.slider.isShow = true
            },
            handleCreateSuccess () {
                this.filter.group_name = ''
                this.filter.user_list = ''
                this.$refs.roleForm.closeRoleForm()
                this.$nextTick(() => {
                    this.getRoleList()
                })
            },
            confirmDeleteRole (role) {
                this.$bkInfo({
                    title: this.$tc('Permission["确认删除角色"]', role['group_name'], {name: role['group_name']}),
                    confirmFn: () => {
                        this.deleteRole(role)
                    }
                })
            },
            async deleteRole (role) {
                await this.deleteUserGroup({bkGroupId: role['group_id']})
                this.$success(this.$t('Permission["删除成功"]'))
                this.getRoleList()
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
                let params = {}
                if (this.filter.type) {
                    params[this.filter.type] = this.filter.text
                }
                const res = await this.searchUserGroup({params, config: {requestId: 'searchUserGroup'}})
                this.table.list = res && res.length ? res : []
                this.setRoles(this.table.list)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .role-options{
        padding: 20px 0;
        .role-options-search{
            position: relative;
            height: 36px;
            line-height: 36px;
        }
        .search-selector {
            position: relative;
            float: left;
            width: 120px;
            margin-right: -1px;
            z-index: 1;
        }
        .cmdb-form-input {
            position: relative;
            width: 300px;
            border-radius: 0 2px 2px 0;
            &:focus {
                z-index: 2;
            }
        }
        .icon-search {
            position: absolute;
            right: 10px;
            top: 11px;
            cursor: pointer;
            font-size: 14px;
        }
        label{
            padding: 0 9px 0 0;
            font-size: 14px;
        }
    }
    .role-table{
        .text-primary {
            margin-right: 10px;
        }
    }
</style>

<style lang="scss">
    .role-wrapper {
        .role-options-search {
            .search-selector {
                input {
                    border-radius: 2px 0 0 2px;
                }
            }
        }
    }
</style>
