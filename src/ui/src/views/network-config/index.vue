<template>
    <div class="config-wrapper">
        <bk-tab class="config-tab" :active-name.sync="activeTabName">
            <bk-tabpanel name="collect" :title="$t('NetworkConfig[\'采集器管理\']')">
                <v-collect v-if="activeTabName === 'collect'"></v-collect>
            </bk-tabpanel>
            <bk-tabpanel name="device" :title="$t('NetworkConfig[\'设备管理\']')">
                <v-device v-if="activeTabName === 'device'"></v-device>
            </bk-tabpanel>
            <bk-tabpanel name="property" :title="$t('NetworkConfig[\'属性管理\']')">
                <v-property v-if="activeTabName === 'property'"></v-property>
            </bk-tabpanel>
        </bk-tab>
    </div>
</template>

<script>
    import vCollect from './collect'
    import vDevice from './device'
    import vProperty from './property'
    export default {
        components: {
            vCollect,
            vDevice,
            vProperty
        },
        data () {
            return {
                activeTabName: 'collect'
            }
        }
    }
</script>

<style lang="scss" scoped>
    .config-wrapper {
        padding: 0;
        background: $cmdbBackgroundColor;
        .config-tab {
            background: $cmdbBackgroundColor;
        }
    }
</style>
