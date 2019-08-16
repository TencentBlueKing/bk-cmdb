<template>
    <div class="role-wrapper">
        <div class="role-options clearfix">
            <div class="role-options-create fl">
                <bk-button theme="primary" @click="createRole">
                    {{$t('新建')}}
                </bk-button>
            </div>
            <div class="role-options-search fr clearfix">
                <bk-select class="search-selector"
                    :clearable="false"
                    v-model="filter.type">
                    <bk-option v-for="(option, index) in typeList"
                        :key="index"
                        :id="option.id"
                        :name="option.name">
                    </bk-option>
                </bk-select>
                <bk-input class="search-input"
                    :right-icon="'bk-icon icon-search'"
                    :placeholder="$t('请输入')"
                    type="text"
                    id="SearchUserName"
                    v-model.trim="filter.text"
                    @enter="getRoleList">
                </bk-input>
                <!-- <i class="filter-search bk-icon icon-search"
                    @click="getRoleList">
                </i> -->
            </div>
        </div>
        <bk-table
            class="role-table"
            v-bkloading="{ isLoading: $loading('searchUserGroup') }"
            :data="table.list"
            :max-height="$APP.height - 240">
            <bk-table-column prop="group_name" :label="$t('角色名')"></bk-table-column>
            <bk-table-column prop="user_list" :label="$t('角色成员')"></bk-table-column>
            <bk-table-column :label="$t('操作')" align="center">
                <template slot-scope="{ row }">
                    <span class="text-primary" @click="showDetails(row)">{{$t('权限详情')}}</span>
                    <span class="text-primary" @click.stop="editRole(row)">{{$t('编辑')}}</span>
                    <span class="text-danger" @click.stop="confirmDeleteRole(row)">{{$t('删除')}}</span>
                </template>
            </bk-table-column>
            <div class="empty-info" slot="empty">
                <p>{{$t('暂时没有数据')}}</p>
                <p>{{$t('当前并无角色，可点击下方按钮新增')}}</p>
                <bk-button class="process-btn" theme="primary" @click="createRole">{{$t('新建角色')}}</bk-button>
            </div>
        </bk-table>
        <v-role-form
            ref="roleForm"
            v-if="form.isShow"
            :data="form.data"
            :type="form.type"
            @on-success="handleCreateSuccess"
            @closeRoleForm="form.isShow = false">
        </v-role-form>
        <bk-sideslider
            :width="600"
            :title="slider.title"
            :is-show.sync="slider.isShow">
            <vAuthority
                slot="content"
                v-if="slider.isShow"
                :group-id="slider.groupId"
                @cancel="slider.isShow = false">
            </vAuthority>
        </bk-sideslider>
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
                    name: this.$t('角色名')
                }, {
                    id: 'user_list',
                    name: this.$t('角色成员')
                }],
                table: {
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
                this.slider.title = `${item['group_name']} ${this.$t('权限详情')}`
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
                    title: this.$tc('确认删除角色', role['group_name'], { name: role['group_name'] }),
                    confirmFn: () => {
                        this.deleteRole(role)
                    }
                })
            },
            async deleteRole (role) {
                await this.deleteUserGroup({ bkGroupId: role['group_id'] })
                this.$success(this.$t('删除成功'))
                this.getRoleList()
            },
            editRole (role) {
                this.form.data = Object.assign({}, role, { 'user_list': role['user_list'].split(';').join(',') })
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
                const params = {}
                if (this.filter.type) {
                    params[this.filter.type] = this.filter.text
                }
                const res = await this.searchUserGroup({
                    params,
                    config: {
                        requestId: 'searchUserGroup'
                    }
                })
                this.table.list = res && res.length ? res : []
                this.setRoles(this.table.list)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .role-options{
        padding: 0 0 14px 0;
        .role-options-search{
            position: relative;
            height: 32px;
            line-height: 32px;
        }
        .search-selector {
            position: relative;
            float: left;
            width: 120px;
            margin-right: -1px;
            z-index: 1;
        }
        .search-input {
            position: relative;
            width: 300px;
            border-radius: 0 2px 2px 0;
            float: left;
            &:focus {
                z-index: 2;
            }
            /deep/ .bk-form-input {
                float: left;
            }
        }
        .icon-search {
            position: absolute;
            right: 10px;
            top: 11px;
            cursor: pointer;
            font-size: 14px;
            z-index: 3;
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
