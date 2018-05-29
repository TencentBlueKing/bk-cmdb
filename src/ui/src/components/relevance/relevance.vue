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
            <bk-button type="primary" class="btn btn-add" @click="currentComponent = 'v-new-association'">{{$t('Association["新增关联"]')}}</bk-button>
        </div>
        <component v-bind="componentProps"
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
            currentComponent (currentComponent, prevComponent) {
                this.prevComponent = prevComponent
            }
        },
        methods: {
            handleNewAssociationClose () {
                this.currentComponent = this.prevComponent
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
    .new-association{
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
    }
</style>
