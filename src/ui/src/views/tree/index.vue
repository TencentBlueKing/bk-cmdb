<template>
    <div class="tree-layout">
        <tree ref="tree" :options="{
            idKey: idGenerator,
            nameKey: 'bk_inst_name',
            childrenKey: 'child'
        }"></tree>
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
                data: []
            }
        },
        created () {
            this.getData()
        },
        methods: {
            async getData () {
                const data = await this.$http.post('find/topoinst/biz/2?level=-1')
                this.data = data
                this.$refs.tree.setData(data)
            },
            idGenerator (node) {
                return `${node.bk_obj_id}_${node.bk_inst_id}`
            }
        }
    }
</script>
