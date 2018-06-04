<template>
    <div class="relevance-content-wrapper">
        <div class="tab-wrapper clearfix">
            <ul class="relevance-tab">
                <li :class="{'active': currentComponent === 'v-topo'}" @click="currentComponent = 'v-topo'">
                    <i class="icon-cc-resources"></i>{{$t('Association["拓扑"]')}}
                </li>
                <li :class="{'active': currentComponent === 'v-tree'}" @click="currentComponent = 'v-tree'">
                    <i class="icon-cc-tree"></i>{{$t('Association["树形"]')}}
                </li>
            </ul>
            <div class="btn-group">
                <span class="resize-btn" @click="resizeFull()" v-if="currentComponent === 'v-topo'" :title="$t('Common[\'全屏\']')">
                    <i class="icon-cc-resize-full"></i>
                </span>
                <bk-button type="primary" class="btn btn-add"
                    :disabled="!hasAssociationProperty"
                    @click="currentComponent = 'v-new-association'">
                    {{$t('Association["新增关联"]')}}
                </bk-button>
            </div>
        </div>
        <component v-bind="componentProps"
            ref="component"
            :is="currentComponent"
            :class="{'new-association': currentComponent === 'v-new-association'}"
            @handleNewAssociationClose="handleNewAssociationClose">
        </component>
    </div>
</template>

<script>
    import vTopo from './topo'
    import vTree from './tree'
    import vNewAssociation from './new-association'
    import {mapGetters} from 'vuex'
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
            },
            instance: Object
        },
        data () {
            return {
                currentComponent: null,
                prevComponent: null
            }
        },
        computed: {
            ...mapGetters('object', ['attribute']),
            hasAssociationProperty () {
                if (this.objId) {
                    return (this.attribute[this.objId] || []).some(property => ['singleasst', 'multiasst'].includes(property['bk_property_type']))
                }
                return false
            },
            componentProps () {
                const component = this.currentComponent
                const props = {
                    'v-topo': {
                        isShow: component === 'v-topo',
                        objId: this.objId,
                        instId: this.ObjectID
                    },
                    'v-tree': {
                        objId: this.objId,
                        ObjectID: this.ObjectID
                    },
                    'v-new-association': {
                        objId: this.objId,
                        instance: this.instance
                    }
                }
                return component ? props[component] : {}
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    this.currentComponent = 'v-topo'
                } else {
                    this.currentComponent = null
                }
            },
            objId (objId) {
                if (this.objId && !this.attribute[this.objId]) {
                    this.$store.dispatch('object/getAttribute', this.objId)
                }
            },
            currentComponent (currentComponent, prevComponent) {
                this.prevComponent = prevComponent
            }
        },
        created () {
            if (this.objId && !this.attribute[this.objId]) {
                this.$store.dispatch('object/getAttribute', this.objId)
            }
        },
        methods: {
            handleNewAssociationClose () {
                this.currentComponent = this.prevComponent
            },
            resizeFull () {
                this.$refs.component.resizeCanvas(true)
            }
        },
        components: {
            vTopo,
            vTree,
            vNewAssociation
        }
    }
</script>

<style lang="scss" scoped>
    .relevance-content-wrapper {
        position: relative;
        height: 100%;
    }
    .tab-wrapper{
        padding: 20px 0 10px;
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
    .btn-group {
        float: right;
        font-size: 0;
    }
    .btn-add {
        height: 24px;
        line-height: 24px;
        font-size: 12px;
        border: none;
        &:disabled{
            cursor: not-allowed !important;
        }
    }
    .resize-btn {
        width: 24px;
        height: 24px;
        line-height: 22px;
        font-size: 14px;
        vertical-align: bottom;
        margin-right: 10px;
    }
    .new-association{
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
    }
</style>
