<template>
    <div class="relevance-content-wrapper">
        <div class="tab-wrapper clearfix">
            <ul class="relevance-tab">
                <li :class="{'active': relevanceTabName === 'topo'}" @click="relevanceTabName = 'topo'">
                    <i class="icon-cc-resources"></i>拓扑
                </li>
                <li :class="{'active': relevanceTabName === 'tree'}" @click="relevanceTabName = 'tree'">
                    <i class="icon-cc-tree"></i>树形
                </li>
            </ul>
            <bk-button type="primary" class="btn btn-add">新增关联</bk-button>
        </div>
        <v-topo v-show="relevanceTabName === 'topo'"
            :isShow="relevanceTabName === 'topo'"
            :objId="objId"
            :instId="ObjectID"
        ></v-topo>
        <v-tree v-show="relevanceTabName === 'tree'" :isShow="relevanceTabName === 'tree'"
            :objId="objId"
            :ObjectID="ObjectID"
        ></v-tree>
    </div>
</template>

<script>
    import vTopo from './topo'
    import vTree from './tree'
    export default {
        props: {
            isShow: {
                type: Boolean,
                default: false
            },
            objId: {
                required: true
            },
            ObjectID: {
                required: true
            }
        },
        data () {
            return {
                relevanceTabName: ''
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    this.relevanceTabName = 'topo'
                } else {
                    this.relevanceTabName = ''
                }
            }
        },
        components: {
            vTopo,
            vTree
        }
    }
</script>

<style lang="scss" scoped>
    .relevance-content-wrapper {
        height: 100%;
    }
    .tab-wrapper{
        padding: 20px 30px;
    }
    .relevance-tab{
        >li{
            float: left;
            margin-right: 2px;
            width: 80px;
            height: 24px;
            line-height: 24px;
            font-size: 12px;
            text-align: center;
            background: #ebf0f5;
            color: #737987;
            cursor: pointer;
            &.active{
                background: #3c96ff;
                color: #fff;
            }
            i{
                position: relative;
                top: -1px;
                margin-right: 5px;
            }
        }
    }
    .btn{
        padding: 0 10px;
    }
    .btn-add {
        float: right;
        height: 24px;
        line-height: 24px;
    }
</style>
