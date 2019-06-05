<template>
    <div class="table-layout">
        <div class="table-title" @click="localExpanded = !localExpanded">
            <cmdb-form-bool class="title-checkbox"
                :size="16"
                @click.native.stop>
            </cmdb-form-bool>
            <i class="title-icon bk-icon icon-right-shape"></i>
            <span class="title-label">{{instance.name}}</span>
        </div>
        <cmdb-table
            v-show="localExpanded"
            :header="header"
            :list="flattenList"
            :height="166"
            :empty-height="42"
            :visible="localExpanded"
            :sortable="false">
            <template slot="__operation__" slot-scope="{ item }">
                <a href="javascript:void(0)" class="text-primary"
                    @click="handleEditProcess(item)">
                    {{$t('Common["编辑"]')}}
                </a>
                <a href="javascript:void(0)" class="text-primary"
                    @click="handleDeleteProcess(item)">
                    {{$t('Common["删除"]')}}
                </a>
            </template>
        </cmdb-table>
    </div>
</template>

<script>
    export default {
        props: {
            instance: {
                type: Object,
                required: true
            },
            expanded: Boolean
        },
        data () {
            return {
                localExpanded: this.expanded,
                properties: [],
                header: [],
                list: [],
                processForm: {
                    show: false,
                    instance: null,
                    processTemplate: null
                }
            }
        },
        computed: {
            processTemplateMap () {
                return this.$store.state.businessTopology.processTemplateMap
            },
            withTemplate () {
                return !!this.instance.service_template_id
            },
            flattenList () {
                return this.$tools.flattenList(this.properties, this.list)
            }
        },
        watch: {
            localExpanded (expanded) {
                if (expanded) {
                    this.getServiceProcessList()
                }
            }
        },
        async created () {
            await this.getProcessProperties()
            if (this.expanded) {
                this.getServiceProcessList()
            }
        },
        methods: {
            async getProcessProperties () {
                const action = 'objectModelProperty/searchObjectAttribute'
                const properties = await this.$store.dispatch(action, {
                    params: {
                        bk_obj_id: 'process',
                        bk_supplier_account: this.$store.getters.supplierAccount
                    },
                    config: {
                        requestId: 'get_service_process_properties',
                        fromCache: true
                    }
                })
                this.properties = properties
                this.setHeader()
            },
            getServiceProcessList () {
                this.list = [{}, {}]
            },
            setHeader () {
                const display = [
                    'bk_process_name',
                    'bind_ip',
                    'port',
                    'work_path',
                    'user'
                ]
                const header = display.map(id => {
                    const property = this.properties.find(property => property.bk_property_id === id) || {}
                    return {
                        id: property.bk_property_id,
                        name: property.bk_property_name
                    }
                })
                header.push({
                    id: '__operation__',
                    name: this.$t('Common["操作"]')
                })
                if (!this.withTemplate) {
                    header.unshift({
                        id: 'id',
                        type: 'checkbox'
                    })
                }
                this.header = header
            },
            async handleEditProcess (item) {
                const processTemplateId = item.process_template_id
                if (processTemplateId) {
                    this.processForm.template = await this.getProcessTemplate(processTemplateId)
                }
                this.processForm.instance = item
                this.processForm.show = true
            },
            async getProcessTemplate (processTemplateId) {
                if (this.processTemplateMap.hasOwnProperty(processTemplateId)) {
                    return Promise.resolve(this.processTemplateMap[processTemplateId])
                }
                const data = await this.$store.dispatch('processTemplate/getProcessTemplate', {
                    params: { processTemplateId }
                })
                this.$store.commit('businessTopology/setProcessTemplate', {
                    id: processTemplateId,
                    template: data.template
                })
                return Promise.resolve(data.template)
            },
            handleDeleteProcess () {}
        }
    }
</script>

<style lang="scss" scoped>
    .table-layout {
        padding: 0 0 12px 0;
    }
    .table-title {
        height: 40px;
        padding: 0 16px;
        line-height: 40px;
        border-radius: 2px 2px 0 0;
        background-color: #DCDEE5;
        cursor: pointer;
        .title-icon {
            font-size: 12px;
            color: #63656E;
            @include inlineBlock;
        }
        .title-label {
            font-size: 14px;
            color: #313238;
            @include inlineBlock;
        }
    }
</style>
