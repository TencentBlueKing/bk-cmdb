<template>
    <div class="display-wrapper">
        <div class="display-setting">
            <label class="cmdb-form-checkbox">
                <input type="checkbox">
                <span class="cmdb-checkbox-text">{{$t('ModelManagement["显示模型名称"]')}}</span>
            </label>
            <label class="cmdb-form-checkbox">
                <input type="checkbox">
                <span class="cmdb-checkbox-text">{{$t('ModelManagement["显示关系名称"]')}}</span>
            </label>
        </div>
        <ul class="display-list">
            <li class="group-item" v-for="(group, groupIndex) in topoList" :key="groupIndex">
                <p class="group-name">{{group['bk_classification_name']}}</p>
                <ul class="clearfix">
                    <li class="model-item" v-for="(model, modelIndex) in group['bk_objects']" :key="modelIndex">
                        <label class="cmdb-form-checkbox">
                            <input type="checkbox">
                            <span class="cmdb-checkbox-text">{{model['bk_obj_name']}}</span>
                            <span class="count">({{group.asstObjects[model['bk_obj_id']].assts.length}})</span>
                            <i class="bk-icon icon-angle-down"></i>
                        </label>
                        <div class="relation-detail">
                            <div class="detail-title clearfix">
                                <div class="fl">
                                    <span class="title">{{$t('ModelManagement["模型关系"]')}}</span>
                                    <span class="info">({{$t('ModelManagement["即视图中的连线"]')}})</span>
                                </div>
                                <label class="fr cmdb-form-checkbox">
                                    <input type="checkbox">
                                    <span class="cmdb-checkbox-text">{{$t('ModelManagement["全选"]')}}</span>
                                </label>
                            </div>
                            <ul class="relation-list clearfix">
                                <li class="fl" v-for="(asst, asstIndex) in group.asstObjects[model['bk_obj_id']].assts" :key="asstIndex">
                                    <label class="cmdb-form-checkbox">
                                        <input type="checkbox" v-model="asst.checked">
                                        <span class="cmdb-checkbox-text">{{asstLabel(model, asst)}}</span>
                                    </label>
                                </li>
                            </ul>
                        </div>
                    </li>
                </ul>
            </li>
        </ul>
        <div class="btn-group">
            <bk-button type="primary" :loading="$loading(['updateAssociationType', 'createAssociationType'])" @click="saveDisplay">
                {{$t('ModelManagement["确定"]')}}
            </bk-button>
            <bk-button type="default" @click="cancel">
                {{$t('ModelManagement["重置"]')}}
            </bk-button>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        props: {
            properties: {
                type: Object
            }
        },
        data () {
            return {

            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'classifications'
            ])
        },
        created () {
            this.initTopoList()
        },
        methods: {
            initTopoList () {
                let {
                    topoModelList
                } = this.properties
                let topoList = []
                this.classifications.forEach(group => {
                    let objects = []
                    let asstObjects = {}
                    group['bk_objects'].forEach(model => {
                        let asstInfo = topoModelList.find(obj => obj['bk_obj_id'] === model['bk_obj_id'] && obj.hasOwnProperty('assts') && obj.assts.length)
                        if (asstInfo) {
                            objects.push(model)
                            asstObjects[model['bk_obj_id']] = {
                                ...asstInfo,
                                ...{
                                    checked: false
                                }
                            }
                        }
                    })
                    if (objects.length) {
                        topoList.push({
                            ...group,
                            ...{
                                bk_objects: objects,
                                asstObjects
                            }
                        })
                    }
                })
                this.topoList = topoList
            },
            asstLabel (model, asst) {
                let asstModel = this.$allModels.find(model => {
                    return model['bk_obj_id'] === asst['bk_obj_id']
                })
                if (asstModel) {
                    return `${model['bk_obj_name']}->${asstModel['bk_obj_name']}`
                }
            },
            saveDisplay () {

            },
            cancel () {
                
            }
        }
    }
</script>


<style lang="scss" scoped>
    .display-wrapper {
        padding: 20px 30px;
    }
    .display-setting {
        .cmdb-form-checkbox {
            width: 154px;
        }
    }
    .display-list {
        .group-item {
            position: relative;
            margin-top: 25px;
            .group-name {
                margin-bottom: 15px;
                padding-left: 8px;
                border-left: 4px solid $cmdbBorderFocusColor;
                line-height: 16px;
            }
        }
        .model-item {
            float: left;
            width: 140px;
            &:hover {
                .relation-detail {
                    display: block;
                }
                .count {
                    &:after {
                        position: absolute;
                        content: '';
                        border: 5px solid transparent;
                        border-bottom-color: #fff;
                        top: 17px;
                        left: calc(50% - 5px);
                        z-index: 2;
                    }
                    &:before {
                        position: absolute;
                        content: '';
                        border: 5px solid transparent;
                        border-bottom-color: $cmdbTableBorderColor;
                        top: 16px;
                        left: calc(50% - 5px);
                        z-index: 2;
                    }
                }
            }
            >.cmdb-form-checkbox {
                margin-right: 5px;
                width: 135px;
                font-size: 0;
                cursor: pointer;
                >.cmdb-checkbox-text {
                    max-width: 65px;
                    font-size: 14px;
                    @include ellipsis;
                }
            }
            .count {
                position: relative;
                font-size: 14px;
                vertical-align: middle;
                color: $cmdbBorderFocusColor;
                
            }
            .relation-detail {
                position: absolute;
                display: none;
                left: 0;
                width: 455px;
                padding: 5px 20px 10px;
                background: #fff;
                border: 1px solid $cmdbTableBorderColor;
                z-index: 1;
                .detail-title {
                    font-size: 0;
                    line-height: 36px;
                    color: $cmdbBorderColor;
                    .title {
                        color: #333948;
                    }
                    .title,
                    .info {
                        font-size: 14px;
                        vertical-align: middle;
                    }
                    .cmdb-checkbox-text {
                        color: $cmdbBorderColor;
                    }
                }
            }
            .relation-list {
                font-size: 0;
                >li:nth-child(3n) {
                    .cmdb-form-checkbox {
                        margin: 0;
                    }
                }
                .cmdb-form-checkbox {
                    width: 130px;
                    margin-right: 10px;
                    @include ellipsis;
                }
            }
            .icon-angle-down {
                margin-left: 5px;
                font-size: 12px;
                color: $cmdbBorderColor;
            }
        }
    }
    .btn-group {
        margin-top: 20px;
        font-size: 0;
        .bk-button {
            margin-right: 10px;
        }
    }
</style>
