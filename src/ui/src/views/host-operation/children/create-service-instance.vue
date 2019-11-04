<template>
    <section class="layout">
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
        <service-instance-table class="service-instance-table"
            v-for="(instance, index) in info"
            :key="index"
            :index="index"
            :id="instance.bk_host_id"
            :name="getName(instance)"
            :deleteable="false"
            :expanded="index === 0"
            :templates="getServiceTemplates(instance)"
            :source-processes="getSourceProcesses(instance)"
            :class="{ 'is-first': index === 0 }">
        </service-instance-table>
    </section>
</template>

<script>
    import ServiceInstanceTable from '@/components/service/instance-table'
    export default {
        components: {
            ServiceInstanceTable
        },
        props: {
            info: {
                type: Array,
                required: true
            }
        },
        methods: {
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
    .service-instance-table {
        &.is-first {
            margin-top: 8px;
        }
        & + .service-instance-table {
            margin: 15px 0 0 0;
        }
    }
</style>
