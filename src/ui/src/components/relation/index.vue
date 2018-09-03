<template>
    <div class="relation-layout">
        <div class="relation-options clearfix">
            <div class="fl">
                <a :class="['options-tab', {active: tab.active === 'topology'}]" href="javascript:void(0)" @click.prevent="tab.active = 'topology'">
                    <i class="icon-cc-resources"></i>
                    {{$t('Association["拓扑"]')}}
                </a>
                <a :class="['options-tab', {active: tab.active === 'tree'}]" href="javascript:void(0)" @click.prevent="tab.active = 'tree'">
                    <i class="icon-cc-tree"></i>
                    {{$t('Association["树形"]')}}
                </a>
            </div>
        </div>
        <div class="relation-component">
            <cmdb-relation-topology></cmdb-relation-topology>
        </div>
    </div>
</template>

<script>
    import cmdbRelationTopology from './_topology.vue'
    import cmdbRelationTree from './_tree.vue'
    import cmdbRelationUpdate from './_update.vue'
    export default {
        components: {
            cmdbRelationTopology,
            cmdbRelationTree,
            cmdbRelationUpdate
        },
        props: {
            objId: {
                type: String,
                required: true
            },
            instId: {
                type: Number,
                required: true
            }
        },
        data () {
            return {
                tab: {
                    active: 'topology'
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .relation-layout {
        height: 100%;
        .relation-options {
            height: 54px;
            padding: 20px 0 10px;
        }
    }
    .relation-options {
        line-height: 24px;
        font-size: 0;
        .options-tab {
            display: inline-block;
            padding: 0 20px;
            margin: 0 2px 0 0;
            vertical-align: middle;
            font-size: 12px;
            text-align: center;
            background-color: #ebf0f5;
            &.active {
                background-color: #3c96ff;
                color: #fff;
            }
        }
    }
    .relation-component {
        height: calc(100% - 54px);
    }
</style>