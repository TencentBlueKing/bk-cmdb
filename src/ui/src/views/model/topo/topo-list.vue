<template>
    <ul class="topolist-wrapper clearfix" v-bkloading="{isLoading: $loading(['searchObjects', 'searchMainlineObject'])}">
        <li class="line"
            :class="{'default': model['bk_obj_id'] === 'biz', 'custom-item': model['ispre']}"
            v-for="(model, index) in topoList"
            :key="index"
            @click="editModel(model)"
        >
            <div class="content">
                <i @click.stop="createModel(model)" class="icon-add icon-cc-round-plus" v-if="model['bk_obj_id'] !== 'biz' && model['bk_obj_id'] !== 'module'"></i>
                <div>
                    <i class="topo-icon" :class="model['bk_obj_icon']"></i>
                </div>
                <div class="content-name">
                    {{model['bk_obj_name']}}
                </div>
            </div>
        </li>
    </ul>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                topoList: [],
                topoStructure: []
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'classifications'
            ])
        },
        methods: {
            createModel (model) {
                this.$emit('createModel', this.findPrevModelId(model))
            },
            editModel (model) {
                if (model['bk_obj_id'] !== 'biz') {
                    this.$emit('editModel', model)
                }
            },
            getBiz () {
                return this.$store.dispatch('objectModel/searchObjects', {params: {bk_obj_id: 'biz'}, config: {requestId: 'searchObjects'}})
            },
            getTopoStructure () {
                return this.$store.dispatch('objectMainLineModule/searchMainlineObject', {requestId: 'searchMainlineObject'})
            },
            getTopoModel () {
                return this.classifications.find(({bk_classification_id: bkClassificationId}) => bkClassificationId === 'bk_biz_topo')['bk_objects']
            },
            findPrevModelId (model) {
                return this.topoStructure.find(({bk_obj_id: bkObjId}) => bkObjId === model['bk_obj_id'])['bk_pre_obj_id']
            },
            sortByStructure (list) {
                let topoList = []
                this.topoStructure.map(({bk_obj_id: bkObjId}) => {
                    let obj = list.find(item => item['bk_obj_id'] === bkObjId)
                    if (obj) {
                        topoList.push(obj)
                    }
                })
                this.topoList = topoList
            },
            async initTopo () {
                const res = await Promise.all([
                    this.getTopoStructure(),
                    this.getBiz()
                ])
                this.topoStructure = res[0]
                let list = this.getTopoModel()
                list.push(res[1][0])
                this.sortByStructure(list)
            }
        },
        created () {
            this.initTopo()
        }
    }
</script>

<style lang="scss" scoped>
    .topolist-wrapper{
        padding: 80px 0;
        margin: 0 auto;
        min-height: 100%;
        li{
            position: relative;
            float: left;
            width: 91px;
            height: 91px;
            border: 1px solid #d6d8df;
            box-shadow: 0 0 10px transparent;
            text-align: center;
            margin-right: 20px;
            margin-top: 30px;
            border-radius: 50%;
            padding: 0 5px;
            cursor: pointer;
            background: #fff;
            color: $cmdbMainBtnColor;
            font-size: 12px;
            font-weight: bold;
            &.line{
                float: none;
                margin: 60px auto 0;
                &::after{
                    content: "";
                    height: 60px;
                    position: absolute;
                    left: 50%;
                    top: 89px;
                    border: 1px dashed #d6d8df;
                }
                &:last-child{
                    &::after{
                        content: "";
                        height: 100%;
                        position: absolute;
                        left: 56px;
                        top: 50px;
                        border:none;
                    }
                }
            }
            &.default{
                margin-top: 0;
                border-style: solid;
                border-width: 1px;
                background: transparent;
                background: $cmdbDefaultColor !important;
                color: #d6d8df !important;
                cursor: default;
            }
            &:not(.default):hover{
                border: 1px solid #d6d8df;
                box-shadow: 0 2.8px 0 rgba(12, 34, 59, 0.05)
            }
            &.custom-item{
                background: $cmdbPrimaryHoverColor;
                color: $cmdbDefaultColor;
            }
            .content{
                white-space: nowrap;
                text-overflow: ellipsis;
                overflow: hidden;
                padding-top: 20px;
                .topo-icon{
                    font-size: 25px;
                }
                .content-name{
                    margin-top: 10px;
                    line-height: 1;
                }
            }
            .icon-add{
                position: absolute;
                display: inline-block;
                border-radius: 50%;
                top: -41px;
                left: 36px;
                padding: 1px 0;
                font-size: 18px;
                color: $cmdbMainBtnColor;
                background: #eee;
                &:hover{
                    color: #50abff;
                }
            }
        }
    }
</style>
