<template>
    <bk-select v-if="display === 'selector'"
        searchable
        :readonly="readonly"
        :disabled="disabled"
        :placeholder="$t('请选择xx', { name: $t('账户名称') })"
        :loading="$loading(Object.values(request))"
        v-model="selected">
        <bk-option v-for="account in accounts"
            :key="account.bk_account_id"
            :name="account.bk_account_name"
            :id="account.bk_account_id">
        </bk-option>
    </bk-select>
    <span v-else>{{selectedAccount ? selectedAccount.bk_account_name : '--'}}</span>
</template>

<script>
    import symbols from '../common/symbol'
    import { CLOUD_AREA_PROPERTIES } from '@/dictionary/request-symbol'
    import { mapGetters } from 'vuex'
    export default {
        name: 'task-account-selector',
        props: {
            display: {
                type: String,
                default: 'selector'
            },
            readonly: Boolean,
            disabled: Boolean,
            value: {
                type: [String, Number]
            }
        },
        data () {
            return {
                accounts: [],
                vendors: [],
                request: {
                    account: symbols.get('account'),
                    properties: Symbol('properties')
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            selected: {
                get () {
                    return this.value
                },
                set (value, oldValue) {
                    this.$emit('input', value)
                    this.$emit('change', value, oldValue)
                }
            },
            selectedAccount () {
                return this.accounts.find(account => account.bk_account_id === this.selected)
            },
            accountVendor () {
                if (!this.selectedAccount) {
                    return null
                }
                // return this.vendors.find(vendor => vendor.name === this.selectedAccount.bk_cloud_vendor)
                const vendors = [{
                    id: '1',
                    name: 'aws'
                }, {
                    id: '2',
                    name: 'tencent_cloud'
                }]
                return vendors.find(vendor => vendor.name === this.selectedAccount.bk_cloud_vendor)
            }
        },
        async created () {
            try {
                const [{ info: accounts }, properties] = await Promise.all([
                    this.getAccounts(),
                    this.getCloudAreaProperties()
                ])
                this.accounts = accounts
                const venderProperty = properties.find(property => property.bk_property_id === 'bk_cloud_vendor')
                this.vendors = venderProperty ? venderProperty.option : []
            } catch (error) {
                this.accounts = []
                this.vendors = []
            }
        },
        methods: {
            getAccounts () {
                return this.$store.dispatch('cloud/account/findMany', {
                    params: {},
                    config: {
                        requestId: this.request.account,
                        fromCache: true
                    }
                })
            },
            getCloudAreaProperties () {
                return this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
                    params: {
                        bk_obj_id: 'plat',
                        bk_supplier_account: this.supplierAccount
                    },
                    config: {
                        requestId: CLOUD_AREA_PROPERTIES,
                        fromCache: true
                    }
                })
            }
        }
    }
</script>
