<template>
    <div class="relation-layout">
        <div class="relation-options clearfix">
            <div class="fl">
                <a href="javascript:void(0)"
                    :class="['options-tab', {active: activeComponent === 'cmdbRelationTopology'}]"
                    @click.prevent="activeComponent = 'cmdbRelationTopology'">
                    <i class="icon-cc-resources"></i>
                    {{$t('Association["拓扑"]')}}
                </a>
                <a  href="javascript:void(0)"
                    :class="['options-tab', {active: activeComponent === 'cmdbRelationTree'}]"
                    @click.prevent="activeComponent = 'cmdbRelationTree'">
                    <i class="icon-cc-tree"></i>
                    {{$t('Association["树形"]')}}
                </a>
            </div>
            <div class="fr" v-if="activeComponent === 'cmdbRelationTopology'">
                <span class="options-full-screen"
                    v-tooltip="$t('Common[\'全屏\']')"
                    @click="handleFullScreen">
                    <i class="icon-cc-resize-full"></i>
                </span>
                <bk-button class="options-create" size="small" type="primary"
                    :disabled="!hasRelation"
                    @click="activeComponent = 'cmdbRelationUpdate'">
                    {{$t('Association["新增关联"]')}}
                </bk-button>
            </div>
        </div>
        <div class="relation-component">
            <component ref="dynamicComponent"
                :is="activeComponent"
                @on-relation-loaded="handleRelationLoaded"
                @on-new-relation-close="activeComponent = 'cmdbRelationTopology'">
            </component>
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
                hasRelation: false,
                fullScreen: false,
                activeComponent: 'cmdbRelationTopology'
            }
        },
        methods: {
            handleFullScreen () {
                this.$refs.dynamicComponent.toggleFullScreen(true)
            },
            handleRelationLoaded (relation) {
                this.hasRelation = !!relation.length
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
        .options-full-screen {
            display: inline-block;
            width: 24px;
            height: 24px;
            margin-right: 10px;
            line-height: 22px;
            font-size: 14px;
            vertical-align: bottom;
            text-align: center;
            border: 1px solid $cmdbBorderColor;
            cursor: pointer;
        }
        .options-create {
            height: 24px;
            line-height: 22px;
            font-size: 12px;
        }
    }
    .relation-component {
        height: calc(100% - 54px);
    }
</style>