<template>
    <bk-dialog
        v-model="isShow"
        :draggable="false"
        :width="730">
        <span class="title" slot="tools">{{$t('筛选条件')}}</span>
        <section class="property-selector">
            <div class="group"
                v-for="group in groups"
                :key="group.id">
                <h2 class="group-title">
                    {{group.name}}
                    <span class="group-count">（{{group.children.length}}）</span>
                </h2>
                <ul class="property-list clearfix">
                    <li class="property-item fl"
                        v-for="property in group.children"
                        :key="property.bk_property_id">
                        <bk-checkbox class="property-checkbox"
                            v-model="selectionMap[group.id][property.bk_property_id].__selected__">
                            {{property.bk_property_name}}
                        </bk-checkbox>
                    </li>
                </ul>
            </div>
        </section>
        <footer class="footer clearfix" slot="footer">
            <i18n class="selected-count fl"
                v-if="selectedSelection.length"
                path="已选择条数"
                tag="div">
                <span class="count" place="count">{{selectedSelection.length}}</span>
            </i18n>
            <div class="selected-options">
                <bk-button theme="primary" @click="confirm">{{$t('确定')}}</bk-button>
                <bk-button theme="default" @click="cancel">{{$t('取消')}}</bk-button>
            </div>
        </footer>
    </bk-dialog>
</template>

<script>
    import { mapState } from 'vuex'
    export default {
        props: {
            properties: {
                type: Object,
                default () {
                    return {}
                }
            }
        },
        data () {
            return {
                isShow: false,
                selection: []
            }
        },
        computed: {
            ...mapState('hosts', ['filterList']),
            groups () {
                return Object.keys(this.properties).map(modelId => {
                    const model = this.$store.getters['objectModelClassify/getModelById'](modelId) || {}
                    return {
                        id: modelId,
                        name: model.bk_obj_name,
                        children: this.properties[modelId]
                    }
                })
            },
            selectionMap () {
                const map = {}
                this.selection.forEach(select => {
                    if (map.hasOwnProperty(select.bk_obj_id)) {
                        map[select.bk_obj_id][select.bk_property_id] = select
                    } else {
                        map[select.bk_obj_id] = {
                            [select.bk_property_id]: select
                        }
                    }
                })
                return map
            },
            selectedSelection () {
                return this.selection.filter(select => select.__selected__)
            }
        },
        watch: {
            async properties () {
                await this.setSelection()
                this.setSelectionState()
            },
            filterList () {
                this.setSelectionState()
            }
        },
        created () {
            this.setSelection()
        },
        methods: {
            setSelection () {
                const selection = []
                Object.keys(this.properties).forEach(modelId => {
                    this.properties[modelId].forEach(property => {
                        selection.push({ ...property, __selected__: false })
                    })
                })
                this.selection = selection
            },
            setSelectionState () {
                const list = this.filterList
                if (list.length) {
                    this.selection.forEach(select => {
                        const selected = list.some(condition => {
                            return condition.bk_property_id === select.bk_property_id
                                && condition.bk_obj_id === select.bk_obj_id
                        })
                        select.__selected__ = selected
                    })
                } else {
                    this.selection.forEach(select => {
                        select.__selected__ = false
                    })
                }
            },
            async confirm () {
                try {
                    const selectedList = this.selectedSelection.map(selection => {
                        return {
                            bk_obj_id: selection.bk_obj_id,
                            bk_property_id: selection.bk_property_id,
                            operator: '',
                            value: ''
                        }
                    })
                    const key = this.$route.meta.filterPropertyKey
                    await this.$store.dispatch('userCustom/saveUsercustom', {
                        [key]: selectedList
                    })
                    this.$store.commit('hosts/setShouldInjectAsset', false)
                    this.$store.commit('hosts/setFilterList', selectedList)
                    this.isShow = false
                } catch (e) {
                    console.error(e)
                }
            },
            cancel () {
                this.isShow = false
            }
        }
    }
</script>

<style lang="scss" scoped>
    .title {
        display: inline-block;
        vertical-align: middle;
        line-height: 31px;
        font-size: 24px;
        color: #444;
        padding: 15px 0 0 24px;
    }
    .property-selector {
        margin: 22px 0 0 0;
        max-height: calc((100vh * 0.3) + 0px);
        @include scrollbar-y;
    }
    .group {
        margin-top: 22px;
        &:first-child {
            margin-top: 0;
        }
        .group-title {
            position: relative;
            padding: 0 0 0 15px;
            line-height: 20px;
            font-size: 15px;
            font-weight: bold;
            color: #63656E;
            &:before {
                content: "";
                position: absolute;
                left: 0;
                top: 3px;
                width: 4px;
                height: 14px;
                background-color: #C4C6CC;
            }
            .group-count {
                color: #C4C6CC;
                font-weight: normal;
            }
        }
    }
    .property-list {
        padding: 10px 0 6px 0;
        .property-item {
            width: 33%;
        }
    }
    .property-checkbox {
        display: block;
        margin: 8px 20px 8px 0;
        @include ellipsis;
        /deep/ {
            .bk-checkbox-text {
                max-width: calc(100% - 25px);
                @include ellipsis;
            }
        }
    }
    .selected-count {
        font-size: 14px;
        line-height: 32px;
        .count {
            color: #2DCB56;
            padding: 0 4px;
        }
    }
</style>
