<template>
    <div class="layout">
        <template v-if="instances.length">
            <div class="options">
                <bk-button class="options-button" type="primary">
                    {{$t('BusinessTopology["添加服务实例"]')}}
                </bk-button>
                <bk-dropdown-menu trigger="click">
                    <bk-button class="options-button clipboard-trigger" type="default" slot="dropdown-trigger">
                        {{$t('Common["更多"]')}}
                        <i class="bk-icon icon-angle-down"></i>
                    </bk-button>
                    <ul class="clipboard-list" slot="dropdown-content">
                        <li v-for="(item, index) in menuItem"
                            :class="['clipboard-item', { 'is-disabled': item.disabled }]"
                            :key="index"
                            @click="item.handler(item.disabled)">
                            {{item.name}}
                        </li>
                    </ul>
                </bk-dropdown-menu>
                <cmdb-form-bool class="options-checkbox"
                    :size="16">
                    <span class="checkbox-label">{{$t('Common["全选本页"]')}}</span>
                </cmdb-form-bool>
                <cmdb-form-bool class="options-checkbox"
                    :size="16">
                    <span class="checkbox-label">{{$t('Common["全部展开"]')}}</span>
                </cmdb-form-bool>
                <cmdb-form-singlechar class="options-search fr"></cmdb-form-singlechar>
            </div>
            <div class="tables">
                <service-instance-table v-for="n in 15" :key="n"></service-instance-table>
            </div>
        </template>
        <service-instance-empty v-else></service-instance-empty>
    </div>
</template>

<script>
    import serviceInstanceTable from './service-instance-table.vue'
    import serviceInstanceEmpty from './service-instance-empty.vue'
    export default {
        components: {
            serviceInstanceTable,
            serviceInstanceEmpty
        },
        data () {
            return {
                instances: []
            }
        },
        computed: {
            currentModule () {
                const node = this.$store.state.businessTopology.selectedNode
                if (node && node.data.bk_obj_id === 'module') {
                    return node.data
                }
                return null
            },
            menuItem () {
                return [{
                    name: this.$t('BusinessTopology["批量编辑"]'),
                    handler: this.batchEdit,
                    disabled: true
                }, {
                    name: this.$t('BusinessTopology["批量删除"]'),
                    handler: this.batchDelete,
                    disabled: true
                }, {
                    name: this.$t('BusinessTopology["复制IP"]'),
                    handler: this.copyIp,
                    disabled: false
                }]
            }
        },
        watch: {
            currentModule (module) {
                if (module) {
                    this.getServiceInstances()
                }
            }
        },
        methods: {
            async getServiceInstances () {
                try {
                    const data = await this.$store.dispatch('serviceInstance/getModuleServiceInstances', {
                        params: this.$injectMetadata({
                            module_id: this.currentModule.bk_inst_id
                        }),
                        config: {
                            requestId: 'getModuleServiceInstances',
                            cancelPrevious: true
                        }
                    })
                    this.instances = data.info
                } catch (e) {
                    console.error(e)
                    this.instances = []
                }
            },
            batchEdit (disabled) {
                if (disabled) {
                    return false
                }
            },
            batchDelete (disabled) {
                if (disabled) {
                    return false
                }
            },
            copyIp () {}
        }
    }
</script>

<style lang="scss" scoped>
    .options {
        padding: 15px 0;
    }
    .options-button {
        height: 32px;
        padding: 0 8px;
        margin: 0 6px 0 0;
        line-height: 30px;
    }
    .options-checkbox {
        margin: 0 19px 0 10px;
        .checkbox-label {
            padding: 0 0 0 9px;
            line-height: 1.5;
        }
    }
    .options-search {
        /deep/ {
            .cmdb-form-input {
                height: 32px;
                line-height: 30px;
            }
        }
    }
    .clipboard-trigger{
        padding: 0 16px;
        .icon-angle-down {
            font-size: 12px;
            top: 0;
        }
    }
    .clipboard-list{
        width: 100%;
        font-size: 14px;
        line-height: 40px;
        max-height: 160px;
        @include scrollbar-y;
        &::-webkit-scrollbar{
            width: 3px;
            height: 3px;
        }
        .clipboard-item{
            padding: 0 15px;
            cursor: pointer;
            @include ellipsis;
            &:not(.is-disabled):hover{
                background-color: #ebf4ff;
                color: #3c96ff;
            }
            &.is-disabled {
                color: #c4c6cc;
                cursor: not-allowed;
            }
        }
    }
    .tables {
        height: calc(100% - 42px);
        @include scrollbar-y;
    }
</style>
