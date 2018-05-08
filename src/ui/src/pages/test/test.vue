<template>
    <div class="test-wrapper">
        <v-table
            :header="table.header"
            :list="table.list"
            :pagination.sync="table.pagination"
            :loading="table.loading"
            :rowBorder="true"
            :colBorder="true"
            :maxHeight="300"
            :height="500">
            <div slot="b" slot-scope="{ item }">{{item.b}}</div>
        </v-table>
        <button @click="addHeader">添加表头</button>
        <button @click="removeHeader">删除表头</button>
        <button @click="addList">添加数据</button>
        <button @click="removeList">删除数据</button>
    </div>
</template>
<script>
    import vTable from '@/components/table2/table'
    export default {
        data () {
            return {
                table: {
                    checked: [1, 2, 3],
                    loading: false,
                    header: [{
                        id: 'z',
                        type: 'checkbox',
                        width: 50
                    }, {
                        id: 'a',
                        name: 'A',
                        width: 100
                    }, {
                        id: 'b',
                        name: 'B',
                        width: 150
                    }, {
                        id: 'c',
                        name: 'C'
                    }, {
                        id: 'd',
                        name: 'D'
                    }],
                    list: [{
                        z: 1,
                        a: 'A-1',
                        b: 'B-1'
                    }, {
                        z: 2,
                        a: 'A-2',
                        b: 'B-2'
                    }, {
                        z: 3,
                        a: 'A-3',
                        b: 'B-3'
                    }, {
                        z: 4,
                        a: 'A-4',
                        b: 'B-4'
                    }, {
                        z: 5,
                        a: 'A-5',
                        b: 'B-5'
                    }],
                    pagination: {
                        current: 1,
                        count: 123,
                        size: 10
                    }
                },
                seed: 1
            }
        },
        methods: {
            addHeader () {
                this.table.header.push({
                    id: this.seed,
                    name: `动态新增${this.seed}`
                })
                this.seed++
            },
            removeHeader () {
                this.table.header.pop()
            },
            addList () {
                this.table.list.push({
                    z: 8,
                    a: 'A-8',
                    b: 'B-8'
                })
            },
            removeList () {
                this.table.list.pop()
            }
        },
        components: {
            vTable
        }
    }
</script>
<style lang="scss" scoped>
    .test-wrapper{
        height: 100%;
        padding: 20px;
    }
</style>