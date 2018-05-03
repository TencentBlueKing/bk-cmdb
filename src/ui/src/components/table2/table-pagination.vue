<template>
    <div class="table-pagination clearfix">
        <div class="table-pagination-info fl">
            <span>已选0行</span>
            <span>第1页/总1页</span>
            <span>每页显示x行</span>
        </div>
        <bk-paging class="table-pagination-list fr" v-bind="pagination"></bk-paging>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                pagination: {
                    curPage: 1,
                    totalPage: 20
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .table-pagination{
        height: 42px;
        line-height: 42px;
        padding: 0 20px;
        background-color: #fafbfd;
        border: 1px solid $borderColor;
        border-top: none;
        .pagination-info{
        }
        .table-pagination-list{
        }
    }
</style>
<style lang="scss">
    .table-pagination-list{
        .page-item{
            height: 32px;
            line-height: 32px;
            margin: 0 0 0 -1px;
        }
    }
</style>