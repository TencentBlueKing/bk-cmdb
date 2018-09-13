<template>
    <bk-dialog 
    :is-show.sync="isShow" 
    :has-header="false" 
    :has-footer="false" 
    :quick-close="false" 
    :width="565" 
    :padding="0"
    @cancel="closeRoleForm">
    >
        <form class="role-form" slot="content">
            <h2 class="role-form-title">{{title}}</h2>
            <div class="role-form-content">
                <div class="content-group clearfix">
                    <label for="groupName" class="fl">{{$t('Permission["角色名"]')}}</label>
                    <input type="text" class="cmdb-form-input fl" id="groupName" v-model.trim="data['group_name']">
                </div>
                <div class="content-group clearfix">
                    <label for="userList" class="fl">{{$t('Permission["角色成员"]')}}</label>
                    <cmdb-form-objuser class="fl member-selector" v-model="data['user_list']" :multiple="true"></cmdb-form-objuser>
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
    import { mapActions } from 'vuex'
    export default {
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
            }
        },
        data () {
            return {
                isShow: true
            }
        },
        computed: {
            title () {
                return this.type === 'create' ? this.$t('Permission["新增角色"]') : this.$t('Permission["编辑角色"]')
            },
            params () {
                let params = {
                    group_name: this.data['group_name'],
                    user_list: this.data['user_list'].split(',').join(';')
                }
                return params
            }
        },
        methods: {
            ...mapActions('userPrivilege', [
                'createUserGroup',
                'updateUserGroup'
            ]),
            submitRoleForm () {
                if (this.type === 'create') {
                    this.createRole()
                } else {
                    this.updateRole()
                }
            },
            async createRole () {
                const res = await this.createUserGroup({params: this.params, config: {requestId: 'saveRole'}})
                this.$success(this.$t('Permission["新建角色成功"]'))
                this.$emit('on-success', res)
            },
            async updateRole () {
                const res = await this.updateUserGroup({bkGroupId: this.data['group_id'], params: this.params, config: {requestId: 'saveRole'}})
                this.$success(this.$t('Permission["更新角色成功"]'))
                this.$emit('on-success', res)
            },
            closeRoleForm () {
                this.isShow = false
                setTimeout(() => {
                    this.$emit('closeRoleForm')
                }, 300)
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
            padding: 0 20px 0 5px;
            text-align: right;
            &:after{
                content: '*';
                color: $cmdbDangerColor;
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
