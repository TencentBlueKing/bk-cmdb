<template>
    <div class="business-wrapper">
        <div class="selector-wrapper clearfix">
            <h2 class="selector-title fl">{{$t('Permission["业务角色"]')}}</h2>
            <div class="selector-container fl clearfix">
                <bk-selector class="role-selector fl"
                    :list="businessRoles"
                    :selected.sync="selectedBusinessRole">
                </bk-selector>
            </div>
        </div>
        <div class="authority-wrapper clearfix">
            <h2 class="authority-title fl">{{$t('Permission["权限设置"]')}}</h2>
            <div class="checkbox-container clearfix">
                <span v-for="(authority, index) in authorities.list" class="checkbox-span fl" :key="index">
                    <label class="cmdb-form-checkbox cmdb-checkbox-small authority-checkbox"
                        :class="{'disabled': isMaintainers}"
                        :for="'business-authority-' + authority.id"
                        :title="$t(authority.name)">
                        <input type="checkbox" 
                            :value="authority.id" 
                            :id="'business-authority-' + authority.id" 
                            :disabled="isMaintainers"
                            v-model="authorities.selected">
                        {{$t(authority.name)}}
                    </label>
                </span>
            </div>
        </div>
        <footer class="footer" v-if="!isMaintainers">
            <bk-button type="primary" @click="updateAuthorities" :loading="$loading('updateAuthorities')">
                {{$t("Common['保存']")}}
            </bk-button>
        </footer>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        data () {
            return {
                selectedBusinessRole: '',
                businessRoles: [],
                authorities: {
                    list: [{
                        id: 'hostupdate',
                        name: 'Permission["主机编辑"]'
                    }, {
                        id: 'hosttrans',
                        name: 'Permission["主机转移"]'
                    }, {
                        id: 'topoupdate',
                        name: 'Permission["拓扑编辑"]'
                    }, {
                        id: 'customapi',
                        name: 'Permission["自定义查询"]'
                    }, {
                        id: 'proconfig',
                        name: 'Permission["进程管理"]'
                    }],
                    selected: []
                }
            }
        },
        computed: {
            isMaintainers () {
                return this.selectedBusinessRole === 'bk_biz_maintainer'
            }
        },
        watch: {
            isMaintainers (isMaintainers) {
                if (isMaintainers) {
                    this.authorities.selected = this.authorities.list.map(authority => {
                        return authority.id
                    })
                }
            },
            selectedBusinessRole () {
                this.getBusinessRoleAuthorities()
            }
        },
        methods: {
            ...mapActions('userPrivilege', [
                'getRolePrivilege',
                'bindRolePrivilege'
            ]),
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            async getBusinessRoles () {
                const res = await this.searchObjectAttribute({
                    params: {bk_obj_id: 'biz'},
                    config: {
                        requestId: 'post_searchObjectAttribute_biz',
                        fromCache: true
                    }
                })
                let roles = []
                res.map(role => {
                    if (role['bk_property_type'] === 'objuser') {
                        roles.push({
                            id: role['bk_property_id'],
                            name: role['bk_property_name'],
                            selectedAuthorities: []
                        })
                    }
                })
                this.businessRoles = roles
                this.selectedBusinessRole = roles[0].id
            },
            async getBusinessRoleAuthorities () {
                const res = await this.getRolePrivilege({bkObjId: 'biz', bkPropertyId: this.selectedBusinessRole})
                if (!this.isMaintainers) {
                    this.authorities.selected = res.length ? res : []
                }
            },
            async updateAuthorities () {
                await this.bindRolePrivilege({bkObjId: 'biz', bkPropertyId: this.selectedBusinessRole, params: this.authorities.selected, config: {requestId: 'updateAuthorities'}})
                this.$success(this.$t('Common[\'保存成功\']'))
            }
        },
        created () {
            this.getBusinessRoles()
        }
    }
</script>

<style lang="scss" scoped>
    .business-wrapper{
        padding: 50px 0 0 0;
        color: #737987;
        font-size: 14px;
    }
    .selector-wrapper{
        .selector-title{
            text-align: right;
            width: 137px;
            height: 36px;
            line-height: 36px;
            font-size: 14px;
            margin: 0;
            padding: 0 30px 0 0;
        }
        .selector-container{
            overflow: visible;
            .role-selector{
                width: 286px;
            }
            a{
                height: 14px;
                line-height: 1;
                display: inline-block;
                vertical-align: unset;
                margin: 11px 0 0 26px;
                color: #498fe0;
                text-decoration: none;
                .icon-cc-derivation{
                    vertical-align: -1px;
                    margin: 0 4px 0 0;
                }
            }
        }
    }
    .authority-wrapper{
        padding: 40px 0 0 0;
        .authority-title{
            width: 137px;
            height: 14px;
            text-align: right;
            line-height: 14px;
            font-size: 14px;
            padding: 0 30px 0 0;
            margin: 0;
        }
        .checkbox-container{
            width: 741px;
            padding-left: 137px;
            overflow: visible;
            line-height: 14px;
        }
    }
    .checkbox-span{
        width: 180px;
        height: 14px;
        margin: 0 16px 28px 0;
        .authority-checkbox{
            max-width: 180px;
            line-height: 14px;
            padding: 0;
            margin: 0;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            cursor: pointer;
            &.disabled{
                cursor: not-allowed;
            }
            input{
                width: 14px;
                height: 14px;
                vertical-align: top;
                &:checked,
                &:focus{
                    border-color: transparent !important;
                }
            }
        }
    }
    .footer{
        margin-top: 20px;
        padding-left: 137px;
    }
</style>
