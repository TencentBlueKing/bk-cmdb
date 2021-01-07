<template>
    <div class="apply-layout">
        <cmdb-tips
            :tips-style="{
                background: 'none',
                border: 'none',
                fontSize: '12px',
                lineHeight: '30px',
                padding: 0
            }"
            :icon-style="{
                color: '#63656E',
                fontSize: '14px',
                lineHeight: '30px'
            }">
            {{$t('转移属性变化确认提示')}}
        </cmdb-tips>
        <property-confirm-table class="mt10"
            ref="confirmTable"
            max-height="auto"
            :list="list"
            :render-icon="true"
            :show-operation="!!conflictList.length">
        </property-confirm-table>
    </div>
</template>

<script>
    import propertyConfirmTable from '@/components/host-apply/property-confirm-table'
    export default {
        name: 'host-attrs-auto-apply',
        components: {
            propertyConfirmTable
        },
        props: {
            info: {
                type: Array,
                required: true
            }
        },
        computed: {
            conflictList () {
                return this.info.filter(item => item.unresolved_conflict_count > 0)
            },
            list () {
                return this.conflictList.length ? this.conflictList : this.info
            }
        },
        methods: {
            getHostApplyConflictResolvers () {
                const conflictResolveResult = this.$refs.confirmTable.conflictResolveResult
                const conflictResolvers = []
                Object.keys(conflictResolveResult).forEach(key => {
                    const propertyList = conflictResolveResult[key]
                    propertyList.forEach(property => {
                        conflictResolvers.push({
                            bk_host_id: Number(key),
                            bk_attribute_id: property.id,
                            bk_property_value: property.__extra__.value
                        })
                    })
                })
                return conflictResolvers
            }
        }
    }
</script>
