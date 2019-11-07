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
        <div class="sort-options">
            <i class="sort-item icon-cc-instance-path"
                :class="{ active: sort === 'module' }"
                v-bk-tooltips="$t('按模块路径排序')"
                @click="setSort('module')">
            </i>
            <i class="sort-grep"></i>
            <i class="sort-item"
                :class="{ active: sort === 'ip' }"
                v-bk-tooltips="$t('按IP排序')"
                @click="setSort('ip')">
                IP
            </i>
        </div>
        <service-instance-table class="service-instance-table"
            v-for="(instance, index) in instances"
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
            setSort (type) {
                if (this.sort === type) {
                    return false
                }
                this.sort = type
                this.sortInfo()
            },
            sortInfo () {
                if (this.sort === 'module') {
                    const order = this.$parent.targetModules
                    this.instances = [...this.info].sort((A, B) => {
                        return order.indexOf(A.bk_module_id) - order.indexOf(B.bk_module_id)
                    })
                } else {
                    this.instances = [...this.info].sort((A, B) => {
                        return A.bk_host_id - B.bk_host_id
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
            border: 1px solid #C4C6CC;
            border-radius: 2px;
            font-size: 0;
            .sort-item {
                width: 30px;
                height: 30px;
                line-height: 30px;
                text-align: center;
                display: inline-block;
                vertical-align: middle;
                font-size: 12px;
                font-style: normal;
                cursor: pointer;
                outline: 0;
                &.active {
                    color: $primaryColor;
                }
            }
            .sort-grep {
                display: inline-block;
                vertical-align: middle;
                height: 14px;
                width: 1px;
                background-color: $borderColor;
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
