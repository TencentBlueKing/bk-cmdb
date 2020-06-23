<template>
    <span class="cmdb-vendor">
        <template v-if="type">
            <svg class="icon" aria-hidden="true">
                <use :xlink:href="icon"></use>
            </svg>
            <slot>{{vendor ? vendor.name : emptyText }}</slot>
        </template>
        <template v-else>{{emptyText}}</template>
    </span>
</template>

<script>
    import { mapGetters } from 'vuex'
    import { CLOUD_AREA_PROPERTIES } from '@/dictionary/request-symbol'
    export default {
        props: {
            type: String,
            emptyText: {
                type: String,
                default: '--'
            }
        },
        data () {
            return {
                vendors: []
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            vendor () {
                return this.vendors.find(vendor => vendor.id === this.type)
            },
            icon () {
                if (!this.vendor) {
                    return null
                }
                const iconMap = {
                    '1': '#icon-cc-cloud-aws',
                    '2': '#icon-cc-cloud-tencent',
                    '3': '#icon-cc-cloud-ali'
                }
                return iconMap[this.vendor.id] || null
            }
        },
        created () {
            this.getVendors()
        },
        methods: {
            async getVendors () {
                try {
                    const properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
                        params: {
                            bk_obj_id: 'plat',
                            bk_supplier_account: this.supplierAccount
                        },
                        config: {
                            requestId: CLOUD_AREA_PROPERTIES,
                            fromCache: true
                        }
                    })
                    const vendorProperty = properties.find(property => property.bk_property_id === 'bk_cloud_vendor')
                    this.vendors = vendorProperty ? vendorProperty.option || [] : []
                } catch (error) {
                    console.error(error)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cmdb-vendor {
        display: inline-flex;
        align-items: center;
        vertical-align: middle;
        @include ellipsis;
        .icon {
            width: 14px;
            height: 14px;
            margin-right: 4px;
        }
    }
</style>
