<template>
    <div class="models-layout" v-show="isEditMode">
        <ul class="groups-list">
            <li class="groups-item"
                v-for="(group, index) in groups"
                :key="index">
                <div class="group-info"
                    :class="{
                        open: openGroup === group['bk_classification_id']
                    }"
                    @click="toggleCollapse(group)">
                    <span class="group-name"
                        :title="group['bk_classification_name']">
                        {{group['bk_classification_name']}}
                    </span>
                    <span class="group-hidden-count">{{group.hiddenModels.length}}</span>
                    <i class="group-collapse-icon bk-icon icon-angle-right"></i>
                </div>
                <template v-if="group.hiddenModels.length">
                    <cmdb-collapse-transition>
                        <ul class="models-list" v-show="openGroup === group['bk_classification_id']">
                            <li class="models-item" draggable
                                v-for="(model, index) in group.hiddenModels"
                                :key="index"
                                @dragstart="handleDragStart($event, model)">
                                <i class="model-icon icon fl" :class="model['bk_obj_icon']"></i>
                                <div class="model-info">
                                    <span class="model-name"
                                        :title="model['bk_obj_name']">
                                        {{model['bk_obj_name']}}
                                    </span>
                                    <span class="model-id"
                                        :title="model['bk_obj_id']">
                                        {{model['bk_obj_id']}}
                                    </span>
                                </div>
                            </li>
                        </ul>
                    </cmdb-collapse-transition>
                </template>
            </li>
        </ul>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        name: 'cmdb-graphics-models',
        data () {
            return {
                unrenderedModels: ['process', 'plat'],
                openGroup: null
            }
        },
        computed: {
            ...mapGetters('globalModels', ['topologyData', 'isEditMode']),
            ...mapGetters('objectAssociation', ['associationList']),
            ...mapGetters('objectModelClassify', ['classifications', 'models']),
            hiddenModels () {
                const hiddenModels = []
                this.topologyData.forEach(data => {
                    const modelId = data['bk_obj_id']
                    if (!this.unrenderedModels.includes(modelId)) {
                        const position = data.position || {}
                        if (typeof position.x !== 'number') {
                            const model = this.models.find(model => model['bk_obj_id'] === modelId)
                            hiddenModels.push(model)
                        }
                    }
                })
                return hiddenModels
            },
            groups () {
                const groups = []
                this.classifications.forEach(classify => {
                    const classifyModels = classify['bk_objects'] || []
                    const classifyHiddenModels = this.hiddenModels.filter(hiddenModel => {
                        return classifyModels.some(model => model['bk_obj_id'] === hiddenModel['bk_obj_id'])
                    })
                    groups.push({
                        ...classify,
                        hiddenModels: classifyHiddenModels
                    })
                })
                return groups
            }
        },
        methods: {
            toggleCollapse (group) {
                const groupId = group['bk_classification_id']
                this.openGroup = groupId === this.openGroup ? null : groupId
            },
            handleDragStart (event, model) {
                event.dataTransfer.setData('modelId', model['bk_obj_id'])
            }
        }
    }
</script>

<style lang="scss" scoped>
    .models-layout {
        width: 200px;
        border: 1px solid $cmdbTableBorderColor;
        border-left: none;
    }
    .group-info {
        position: relative;
        height: 42px;
        padding: 0 20px 0 15px;
        line-height: 42px;
        font-size: 0;
        cursor: pointer;
        &:hover,
        &.open {
            color: #fff;
            background-color: $cmdbBorderFocusColor;
        }
        &:hover {
            opacity: .65;
        }
        &.open {
            opacity: 1;
            .group-hidden-count {
                color: $cmdbBorderFocusColor;
                background-color: #fff;
            }
            .group-collapse-icon {
                transform: rotate(90deg);
            }
        }
        .group-name {
            display: inline-block;
            max-width: 120px;
            vertical-align: middle;
            font-size: 14px;
            @include ellipsis;
        }
        .group-hidden-count {
            display: inline-block;
            padding: 0 4px;
            margin: 0 4px;
            border-radius: 4px;
            vertical-align: middle;
            background-color: #ebf4ff;
            line-height: 16px;
            font-size: 12px;
        }
        .group-collapse-icon {
            position: absolute;
            top: 15px;
            right: 15px;
            font-size: 12px;
            transition: transform .2s linear;
        }
    }
    .models-list {
        .models-item {
            height: 56px;
            padding: 10px 12px;
            cursor: move;
            &:hover {
                background-color: #ebf4ff;
            }
            .model-icon {
                width: 36px;
                height: 36px;
                margin: 0 5px 0 0;
                border: 1px solid $cmdbTableBorderColor;
                border-radius: 50%;
                line-height: 34px;
                font-size: 20px;
                text-align: center;
                color: $cmdbBorderFocusColor;
            }
            .model-info {
                line-height: 18px;
                font-size: 12px;
                overflow: hidden;
            }
        }
    }
    .model-info {
        .model-name,
        .model-id {
            display: block;
            @include ellipsis;
        }
        .model-id {
            color: $cmdbBorderColor;
        }
    }
</style>
