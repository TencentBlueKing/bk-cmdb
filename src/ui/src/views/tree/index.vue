<template>
    <div class="tree-layout">
        <bk-button @click="addNode">添加</bk-button>
        <bk-button @click="deleteNode">删除</bk-button>
        <tree ref="tree"
            :options="{
                idKey: idGenerator,
                nameKey: 'bk_inst_name',
                childrenKey: 'child'
            }"
            :node-icon="getNodeIcon"
            :default-expand-node="[8, 20]"
            show-checkbox
            expand-icon="icon-cc-rect-sub"
            collapse-icon="icon-cc-rect-add">
        </tree>
    </div>
</template>

<script>
    import tree from './tree.vue'
    export default {
        components: {
            tree
        },
        data () {
            return {
                data: [],
                id: 0
            }
        },
        created () {
            this.getData()
        },
        methods: {
            async getData () {
                const data = await this.$http.post('find/topoinst/biz/5?level=-1')
                // this.data = data
                // this.$refs.tree.setData(data)
                const repeat = []
                for (let i = 0; i < 300; i++) {
                    repeat.push(data[0])
                }
                this.$refs.tree.setData(repeat)
            },
            idGenerator (node) {
                return this.id++
            },
            getNodeIcon (data) {
                return 'icon-cc-host'
            },
            addNode () {
                this.$refs.tree.addNode({
                    bk_inst_id: Math.random(),
                    bk_inst_name: Math.random(),
                    bk_obj_id: 'fuck',
                    bk_obj_name: 'fuck',
                    child: [],
                    default: 0
                }, 8)
            },
            deleteNode () {}
        }
    }
</script>
