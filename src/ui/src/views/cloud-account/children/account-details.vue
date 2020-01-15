<template>
    <div class="details-layout"
        v-bkloading="{ isLoading: $loading(requestId) }">
        <bk-form :label-width="105">
            <bk-form-item class="details-item" :label="$t('账户名称')">
                {{account.name | formatter}}
            </bk-form-item>
            <bk-form-item class="details-item" :label="$t('账户类型')">
                {{account.type | formatter}}
            </bk-form-item>
            <bk-form-item class="details-item" label="ID">
                {{account.id | formatter}}
            </bk-form-item>
            <bk-form-item class="details-item" label="Key">
                {{account.key | formatter}}
            </bk-form-item>
            <bk-form-item class="details-item" :label="$t('备注')">
                {{account.remarks | formatter}}
            </bk-form-item>
        </bk-form>
        <bk-form class="extra-info-form" :label-width="105">
            <bk-form-item class="details-item" :label="$t('创建人')">
                {{account.creator | formatter}}
            </bk-form-item>
            <bk-form-item class="details-item" :label="$t('创建时间')">
                {{account.create_at | formatter}}
            </bk-form-item>
            <bk-form-item class="details-item" :label="$t('修改人')">
                {{account.updator | formatter}}
            </bk-form-item>
            <bk-form-item class="details-item" :label="$t('修改时间')">
                {{account.update_at | formatter}}
            </bk-form-item>
            <bk-form-item class="details-options">
                <bk-button class="mr10" theme="primary" @click="handleEdit">{{$t('编辑')}}</bk-button>
                <bk-button @click="handleCancel">{{$t('取消')}}</bk-button>
            </bk-form-item>
        </bk-form>
    </div>
</template>

<script>
    export default {
        name: 'cloud-account-details',
        filters: {
            formatter (value) {
                return value || '--'
            }
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
                requestId: Symbol('getAccountData')
            }
        },
        created () {
            this.getAccountData()
        },
        methods: {
            async getAccountData () {
                try {
                    this.account = await Promise.resolve({ id: '1', name: '测试' })
                } catch (e) {
                    this.account = {}
                    console.error(e)
                }
            },
            handleEdit () {
                this.container.show({
                    type: 'form',
                    title: `${this.$t('编辑账户')} 【${this.account.name}】`,
                    props: {
                        mode: 'edit',
                        account: this.account
                    }
                })
            },
            handleCancel () {
                this.container.hide()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .details-layout {
        padding: 18px 27px;
        .details-item {
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
                        display: block;
                        @include ellipsis;
                    }
                }
                .bk-form-content {
                    color: #313238;
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
