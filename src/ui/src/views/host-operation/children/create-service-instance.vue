<template>
    <section class="create-layout">
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
            {{$t('新增服务实例提示')}}
        </cmdb-tips>
        <bk-select class="sort-options"
            :style="{ width: $i18n.locale === 'en' ? '110px' : '90px' }"
            :clearable="false"
            prefix-icon="icon-cc-order"
            v-model="sort"
            @selected="sortInfo">
            <bk-option id="ip" name="IP"></bk-option>
            <bk-option id="module" :name="$t('模块')"></bk-option>
        </bk-select>
        <service-instance-table class="service-instance-table"
            v-for="(instance, index) in instances"
            ref="serviceInstance"
            :key="`${instance.bk_module_id}-${instance.bk_host_id}`"
            :index="index"
            :id="instance.bk_host_id"
            :name="getName(instance)"
            :deletable="false"
            :topology="$parent.getModulePath(instance.bk_module_id)"
            :templates="getServiceTemplates(instance)"
            :source-processes="getSourceProcesses(instance)"
            :class="{ 'is-first': index === 0 }">
        </service-instance-table>
    </section>
</template>

<script>
    import ServiceInstanceTable from '@/components/service/instance-table'
    export default {
        name: 'create-service-instance',
        components: {
            ServiceInstanceTable
        },
        props: {
            info: {
                type: Array,
                required: true
            }
        },
        data () {
            return {
                sort: 'module',
                instances: []
            }
        },
        watch: {
            info: {
                immediate: true,
                handler () {
                    this.sortInfo()
                }
            }
        },
        methods: {
            sortInfo () {
                if (this.sort === 'module') {
                    const order = this.$parent.targetModules
                    this.instances = [...this.info].sort((A, B) => {
                        return order.indexOf(A.bk_module_id) - order.indexOf(B.bk_module_id)
                    })
                } else {
                    this.instances = [...this.info].sort((A, B) => {
                        return this.getName(A).localeCompare(this.getName(B))
                    })
                }
            },
            getName (instance) {
                const data = this.$parent.hostInfo.find(data => data.host.bk_host_id === instance.bk_host_id)
                if (data) {
                    return data.host.bk_host_innerip
                }
                return '--'
            },
            getServiceTemplates (instance) {
                if (instance.service_template) {
                    return instance.service_template.process_templates
                }
                return []
            },
            getSourceProcesses (instance) {
                const templates = this.getServiceTemplates(instance)
                return templates.map(template => {
                    const value = {}
                    const ip = ['127.0.0.1', '0.0.0.0']
                    Object.keys(template.property).forEach(key => {
                        if (key === 'bind_ip') {
                            value[key] = ip[template.property[key].value - 1]
                        } else {
                            value[key] = template.property[key].value
                        }
                    })
                    return value
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-layout {
        position: relative;
        .sort-options {
            position: absolute;
            top: -4px;
            right: 0;
            /deep/ .icon-cc-order {
                font-size: 14px;
                color: #979BA5;
            }
        }
    }
    .service-instance-table {
        &.is-first {
            margin-top: 8px;
        }
        & + .service-instance-table {
            margin: 15px 0 0 0;
        }
    }
</style>
