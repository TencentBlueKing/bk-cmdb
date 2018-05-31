<template>
    <div class="association-table">
        <label class="table-label">{{property['bk_property_name']}}</label>
        <v-table
            :header="header"
            :list="list"
            :width="width"
            :height="height"
            :emptyHeight="emptyHeight"
            :colBorder="true"
            :sortable="false"
            :checked.sync="checked"
            :crossPageCheck="false">
        </v-table>
    </div>
</template>

<script>
    import vTable from '@/components/table/table'
    export default {
        components: {
            vTable
        },
        props: {
            property: {
                type: Object,
                required: true
            },
            list: {
                type: Array,
                default () {
                    return []
                }
            },
            width: {
                type: Number,
                default: 700
            },
            height: {
                type: Number,
                default: 200
            },
            emptyHeight: {
                type: Number,
                default: 40
            }
        },
        data () {
            return {
                checked: [],
                header: [{
                    id: 'bk_inst_id',
                    name: 'ID',
                    type: 'checkbox',
                    width: 50
                }, {
                    id: 'bk_inst_id',
                    name: 'ID'
                }, {
                    id: 'bk_inst_name',
                    name: this.$t('Association["实例名"]')
                }]
            }
        },
        watch: {
            checked (checked) {
                this.$emit('handleCheck', checked, this.property)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .association-table{
        margin: 10px 0 0 0;
        .table-label{
            display: block;
            line-height: 32px;
        }
    }
</style>