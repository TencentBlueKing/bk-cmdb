<template>
    <div class="apply-layout">
        <cmdb-tips
            :tips-style="{
                background: 'none',
                border: 'none',
                fontSize: '12px'
            }"
            :icon-style="{
                color: '#63656E',
                fontSize: '14px'
            }">
            {{$t('转移属性变化确认提示')}}
        </cmdb-tips>
        <property-confirm-table
            ref="confirmTable"
            :list="info"
            :max-height="600"
            :render-icon="true">
        </property-confirm-table>
    </div>
</template>

<script>
    import propertyConfirmTable from '@/views/host-apply/children/property-confirm-table.vue'
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
        }
    }
</script>

<style lang="scss" scoped>

</style>
