<template>
    <div class="search-item-layout">
        <div class="clearfix">
            <ul class="search-list fl"
                v-for="col in columns"
                :key="col">
                <li class="search-item"
                    v-for="row in itemPerCol"
                    v-if="getCellIndex(col, row) < list.length">
                    <v-search-item-host v-if="model === 'host'" :host="list[getCellIndex(col, row)]"></v-search-item-host>
                </li>
            </ul>
        </div>
    </div>
</template>

<script>
    import vSearchItemHost from './search-item-host'
    export default{
        components: {
            vSearchItemHost
        },
        props: {
            model: {
                type: String,
                required: true
            },
            list: {
                type: Array,
                required: true
            }
        },
        data () {
            return {
                itemPerCol: 5
            }
        },
        computed: {
            columns () {
                const totalColumns = Math.ceil(this.list.length / this.itemPerCol)
                return totalColumns > 3 ? 3 : totalColumns
            }
        },
        methods: {
            getCellIndex (col, row) {
                return (col - 1) * this.itemPerCol + row - 1
            }
        }
    }
</script>

<style lang="scss" scoped>
    .search-item-layout{
        padding: 20px 0;
        font-size: 12px;
    }
    .search-list{
        padding: 0 14px 0 18px;
        width: 33%;
        border-right: 1px solid #ebf0f5;
        &:first-child{
            padding-left: 0;
        }
        &:last-child{
            padding-right: 0;
        }
        &:nth-child(3){
            border-right: none;
        }
        .search-item{
            height: 24px;
            line-height: 24px;
            margin: 0 0 8px 0;
            cursor: pointer;
            &:last-child{
                margin-bottom: 0;
            }
            &:hover{
                background-color: #f4f7fa;
            }
        }
    }
</style>