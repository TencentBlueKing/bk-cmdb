<template>
    <div class="config-wrapper">
        <bk-tab class="config-tab" type="unborder-card" :active.sync="activeTabName">
            <bk-tab-panel name="collect" :label="$t('采集器管理')">
                <v-collect v-if="activeTabName === 'collect'"></v-collect>
            </bk-tab-panel>
            <bk-tab-panel name="device" :label="$t('设备管理')">
                <v-device v-if="activeTabName === 'device'"></v-device>
            </bk-tab-panel>
            <bk-tab-panel name="property" :label="$t('属性管理')">
                <v-property v-if="activeTabName === 'property'"></v-property>
            </bk-tab-panel>
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
        },
        created () {
            this.$route.meta.title = this.$t('网络发现配置')
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
