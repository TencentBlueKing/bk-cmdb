<template>
    <div class="details-layout"
        v-bkloading="{ isLoading: $loading(request.data) }">
        <bk-form :label-width="105">
            <bk-form-item class="details-item" :label="$t('账户名称')">
                {{account.bk_account_name}}
            </bk-form-item>
            <bk-form-item class="details-item" :label="$t('账户类型')">
                <cmdb-vendor :type="account.bk_cloud_vendor"></cmdb-vendor>
            </bk-form-item>
            <bk-form-item class="details-item" label="ID">
                {{account.bk_secret_id}}
            </bk-form-item>
            <bk-form-item class="details-item" label="Key">{{account.bk_secret_key}}</bk-form-item>
            <bk-form-item class="details-item" :label="$t('备注')">
                {{account.bk_description | formatter('longchar')}}
            </bk-form-item>
        </bk-form>
        <bk-form class="extra-info-form" :label-width="105">
            <bk-form-item class="details-item" :label="$t('创建人')">
                {{account.bk_creator | formatter('singlechar')}}
            </bk-form-item>
            <bk-form-item class="details-item" :label="$t('创建时间')">
                {{account.create_time | formatter('time')}}
            </bk-form-item>
            <bk-form-item class="details-item" :label="$t('修改人')">
                {{account.bk_last_editor | formatter('singlechar')}}
            </bk-form-item>
            <bk-form-item class="details-item" :label="$t('修改时间')">
                {{account.last_time | formatter('time')}}
            </bk-form-item>
            <bk-form-item class="details-options">
                <cmdb-auth class="mr10" :auth="{ type: $OPERATION.U_CLOUD_ACCOUNT, relation: [id] }">
                    <bk-button theme="primary" slot-scope="{ disabled }"
                        :disabled="disabled || $loading(request.delete)"
                        @click="handleEdit">
                        {{$t('编辑')}}
                    </bk-button>
                </cmdb-auth>
                <cmdb-auth class="inline-block-middle"
                    :auth="{ type: $OPERATION.D_CLOUD_ACCOUNT, relation: [id] }"
                    v-bk-tooltips="{
                        disabled: account.bk_can_delete_account,
                        content: $t('云账户禁止删除提示')
                    }">
                    <bk-button slot-scope="{ disabled }"
                        :disabled="disabled || !account.bk_can_delete_account"
                        :loading="$loading(request.delete)"
                        @click="handleDelete">
                        {{$t('删除')}}
                    </bk-button>
                </cmdb-auth>
            </bk-form-item>
        </bk-form>
    </div>
</template>

<script>
    import CmdbVendor from '@/components/ui/other/vendor'
    import RouterQuery from '@/router/query'
    export default {
        name: 'cloud-account-details',
        components: {
            CmdbVendor
        },
        props: {
            id: {
                type: Number,
                required: true
            },
            container: {
                type: Object,
                default: () => ({})
            }
        },
        data () {
            return {
                account: {},
                request: {
                    data: Symbol('getAccountData'),
                    delete: Symbol('delete')
                }
            }
        },
        created () {
            this.getAccountData()
        },
        methods: {
            async getAccountData () {
                try {
                    const account = await this.$store.dispatch('cloud/account/findOne', {
                        id: this.id,
                        config: {
                            requestId: this.request.data
                        }
                    })
                    this.account = {
                        ...account,
                        bk_secret_key: '******'
                    }
                } catch (e) {
                    this.account = {}
                    console.error(e)
                }
            },
            handleEdit () {
                this.container.show({
                    type: 'form',
                    title: `${this.$t('编辑账户')} 【${this.account.bk_account_name}】`,
                    props: {
                        mode: 'edit',
                        account: this.account
                    }
                })
            },
            handleDelete () {
                const infoInstance = this.$bkInfo({
                    title: this.$t('确认删除xx', { instance: this.account.bk_account_name }),
                    closeIcon: false,
                    confirmFn: () => {
                        return new Promise(async resolve => {
                            infoInstance.buttonLoading = true
                            try {
                                await this.$store.dispatch('cloud/account/delete', {
                                    id: this.account.bk_account_id,
                                    config: {
                                        requestId: this.request.delete
                                    }
                                })
                                this.$success('删除成功')
                                this.container.hide()
                                RouterQuery.set({
                                    _t: Date.now(),
                                    page: RouterQuery.get('page', 1)
                                })
                            } catch (error) {
                                console.error(error)
                            } finally {
                                resolve(true)
                            }
                        })
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .details-layout {
        padding: 18px 27px;
        .details-item {
            margin-top: 10px;
            /deep/ {
                .bk-label {
                    position: relative;
                    padding: 0 20px 0 0;
                    text-align: left;
                    &:after {
                        position: absolute;
                        right: 4px;
                        top: 0;
                        content: '：'
                    }
                    span {
                        display: inline-block;
                        vertical-align: middle;
                        @include ellipsis;
                    }
                }
                .bk-form-content {
                    color: #313238;
                    font-size: 14px;
                    line-height: 32px;
                }
            }
        }
        .extra-info-form {
            margin-top: 24px;
            padding-top: 18px;
            border-top: 1px solid #F0F1F5;
        }
        .details-options {
            font-size: 0;
        }
    }
</style>
