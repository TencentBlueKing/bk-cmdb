<template>
    <div class="service-table-layout">
        <div class="title">
            <div class="fl" @click="localExpanded = !localExpanded">
                <i class="bk-icon icon-down-shape" v-if="localExpanded"></i>
                <i class="bk-icon icon-right-shape" v-else></i>
                {{name}}
            </div>
            <div class="fr">
                <i class="bk-icon icon-close" v-if="deletable" @click="handleDelete"></i>
            </div>
        </div>
        <cmdb-table
            :header="header"
            :list="processList"
            :empty-height="58">
            <template slot="data-empty">
                <a href="javascript:void(0)" class="text-primary">
                    <i class="bk-icon icon-plus"></i>
                    <span>{{$t('BusinessTopology["添加进程"]')}}</span>
                </a>
            </template>
        </cmdb-table>
    </div>
</template>

<script>
    export default {
        props: {
            deletable: Boolean,
            expanded: Boolean,
            id: {
                type: Number,
                required: true
            },
            index: {
                type: Number,
                required: true
            },
            name: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                localExpanded: this.expanded,
                processList: [],
                processProperties: []
            }
        },
        computed: {
            header () {
                const display = [
                    'bk_process_name',
                    'bind_ip',
                    'port',
                    'work_path',
                    'user'
                ]
                const header = []
                display.map(id => {
                    const property = this.processProperties.find(property => property.bk_property_id === id)
                    if (property) {
                        header.push({
                            id: property.bk_property_id,
                            name: property.bk_property_name
                        })
                    }
                })
                header.push({
                    id: '__operation__',
                    name: this.$t('Common["操作"]')
                })
                return header
            }
        },
        created () {
            this.getProcessProperties()
        },
        methods: {
            async getProcessProperties () {
                try {
                    const action = 'objectModelProperty/searchObjectAttribute'
                    this.processProperties = await this.$store.dispatch(action, {
                        params: {
                            bk_obj_id: 'process',
                            bk_supplier_account: this.$store.getters.supplierAccount
                        },
                        config: {
                            requestId: 'get_service_process_properties',
                            fromCache: true
                        }
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            handleDelete () {
                this.$emit('delete-instance', this.index)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .title {
        height: 40px;
        padding: 0 10px;
        line-height: 40px;
        border-radius: 2px 2px 0 0;
        background-color: #DCDEE5;
        .bk-icon {
            font-size: 12px;
            font-weight: bold;
            width: 24px;
            height: 24px;
            line-height: 24px;
            text-align: center;
            cursor: pointer;
            @include inlineBlock;
        }
    }
</style>
