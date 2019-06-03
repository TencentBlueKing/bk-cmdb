<template>
    <div class="table-layout">
        <div class="table-title">
            <cmdb-form-bool class="title-checkbox"
                :size="16">
            </cmdb-form-bool>
            <i class="title-icon bk-icon icon-right-shape"></i>
            <span class="title-label">192.168.1.1</span>
        </div>
        <cmdb-table
            :header="header"
            :list="list"
            :height="166"
            :empty-height="42">
            <template slot="__operation__" slot-scope="{ item }">
                <a href="javascript:void(0)" class="text-primary"
                    @click="handleEditProcess(item)">
                    {{$t('Common["编辑"]')}}
                </a>
            </template>
        </cmdb-table>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                properties: [],
                header: [],
                list: [{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}]
            }
        },
        computed: {},
        created () {
            this.getProcessProperties()
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
                this.header = header
            },
            handleEditProcess (item) {}
        }
    }
</script>

<style lang="scss" scoped>
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
