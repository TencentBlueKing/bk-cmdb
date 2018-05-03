<template>
    <div class="test-wrapper">
        <v-table v-bind="table"></v-table>
    </div>
</template>
<script>
    import vTable from '@/components/table2/table'
    export default {
        data () {
            return {
                table: {
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
                    }, {
                        z: 6,
                        a: 'A-6',
                        b: 'B-6'
                    }, {
                        z: 7,
                        a: 'A-7',
                        b: 'B-7'
                    }, {
                        z: 8,
                        a: 'A-8',
                        b: 'B-8'
                    }]
                    // header: [],
                    // list: []
                }
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