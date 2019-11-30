<template>
    <section>
        <p class="title">{{`${$t('服务分类')}：${serviceCategory}`}}</p>
        <bk-table ref="table"
            v-bkloading="{ isLoading: $loading() }"
            :data="processes"
            :show-header="!!processes.length"
            :height="276"
            :max-height="276">
            <bk-table-column v-for="head in header"
                :key="head.id"
                :prop="head.id"
                :label="head.name">
                <template slot-scope="{ row }">{{row[head.id] | formatter(head.type)}}</template>
            </bk-table-column>
        </bk-table>
    </section>
</template>

<script>
    export default {
        name: 'serviceTemplateInfo',
        props: {
            id: {
                type: Number,
                required: true
            }
        },
        data () {
            return {
                serviceCategory: '',
                header: [{
                    id: 'bk_func_name',
                    name: this.$t('进程名称'),
                    type: 'singlechar'
                }, {
                    id: 'bind_ip',
                    name: this.$t('监听IP'),
                    type: 'singlechar'
                }, {
                    id: 'port',
                    name: this.$t('端口'),
                    type: 'singlechar'
                }, {
                    id: 'work_path',
                    name: this.$t('启动路径'),
                    type: 'longchar'
                }, {
                    id: 'user',
                    name: this.$t('启动用户'),
                    type: 'singlechar'
                }],
                processes: []
            }
        },
        created () {
            setTimeout(() => {
                this.$refs.table.doLayout()
            }, 0)
            this.getTitle()
            this.getServiceProcesses()
        },
        methods: {
            close () {
                this.visible = false
            },
            async getServiceProcesses () {
                try {
                    const result = await this.$store.dispatch('processTemplate/getBatchProcessTemplate', {
                        params: this.$injectMetadata({
                            service_template_id: this.id
                        }),
                        config: {
                            requestId: 'getServiceProcesses'
                        }
                    })
                    this.processes = result.info.map(data => {
                        const process = {}
                        Object.keys(data.property).forEach(key => {
                            process[key] = data.property[key].value
                        })
                        return process
                    })
                } catch (e) {
                    console.error(e)
                    this.processes = []
                }
            },
            async getTitle () {
                try {
                    const [details, categoryData] = await Promise.all([
                        this.getServiceDetails(),
                        this.getServiceCategory()
                    ])
                    const categories = categoryData.info
                    const categoryId = details.template.service_category_id
                    const subCategory = categories.find(data => data.id === categoryId) || {}
                    const category = categories.find(data => data.id === subCategory.bk_parent_id) || {}
                    this.serviceCategory = `${category.name} / ${subCategory.name}`
                } catch (e) {
                    console.error(e)
                    this.serviceCategory = ''
                }
            },
            getServiceDetails () {
                return this.$store.dispatch('serviceTemplate/findServiceTemplate', {
                    id: this.id,
                    config: {
                        requestId: 'getServiceDetails'
                    }
                })
            },
            getServiceCategory () {
                return this.$store.dispatch('serviceClassification/searchServiceCategoryWithoutAmout', {
                    params: this.$injectMetadata({}),
                    config: {
                        requestId: 'getServiceCategoryWithoutAmount'
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .title {
        font-size: 14px;
        line-height: 20px;
        margin-top: -6px;
        padding: 0 0 16px 0;
    }
</style>
