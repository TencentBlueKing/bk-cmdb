<template>
    <cmdb-collapse-transition>
        <div class="classify-layout" v-show="showClassify">
            <div class="classify"
                v-for="(classify, index) in displayClassifications"
                :key="index">
                <h2 class="classify-name">
                    {{`${classify.bk_classification_name}(${classify.bk_objects.length})`}}
                </h2>
                <ul class="model-list">
                    <li class="model-item"
                        v-for="(model, modelIndex) in classify.bk_objects"
                        :key="modelIndex">
                        <div class="model-info"
                            :title="model.bk_obj_name"
                            @click="handleClick(model)">
                            <i :class="[
                                'model-icon',
                                model.bk_obj_icon,
                                {
                                    'is-custom': isCustomModel(model)
                                }
                            ]">
                            </i>
                            <span class="model-name">{{model.bk_obj_name}}</span>
                            <i
                                :class="[
                                    'model-collect',
                                    'bk-icon',
                                    isCollected(model) ? 'icon-star-shape' : 'icon-star'
                                ]"
                                @click.stop="handleToggleCollection(model)">
                            </i>
                        </div>
                    </li>
                </ul>
            </div>
        </div>
    </cmdb-collapse-transition>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                maxCollectCount: 8
            }
        },
        computed: {
            ...mapGetters(['isAdminView']),
            ...mapGetters('index', ['showClassify']),
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('objectModelClassify', [
                'models',
                'activeClassifications'
            ]),
            collectedData () {
                return this.usercustom.collected_models || []
            },
            availableCollectedData () {
                return this.collectedData.filter(id => this.models.some(model => model.id === id))
            },
            displayClassifications () {
                const noDisplay = ['bk_host_manage', 'bk_biz_topo', 'bk_organization']
                return this.activeClassifications.filter(classify => {
                    return !noDisplay.includes(classify.bk_classification_id)
                })
            }
        },
        methods: {
            isCustomModel (model) {
                if (this.isAdminView) {
                    return model.ispre
                }
                return !!this.$tools.getMetadataBiz(model)
            },
            isCollected (model) {
                return this.availableCollectedData.includes(model.id)
            },
            handleClick (model) {
                this.$router.push({
                    name: 'generalModel',
                    params: {
                        objId: model.bk_obj_id
                    }
                })
            },
            async handleToggleCollection (model) {
                const isCollected = this.isCollected(model)
                let collectedData = [...this.collectedData]
                if (isCollected) {
                    collectedData = collectedData.filter(id => id !== model.id)
                } else if (this.availableCollectedData.length < this.maxCollectCount) {
                    collectedData.push(model.id)
                } else {
                    this.$warn(this.$t('限制添加导航提示', { max: this.maxCollectCount }))
                    return false
                }
                await this.$store.dispatch('userCustom/saveUsercustom', {
                    collected_models: collectedData
                })
                this.$success(
                    isCollected
                        ? this.$t('取消导航成功')
                        : this.$t('添加导航成功')
                )
            }
        }
    }
</script>

<style lang="scss" scoped>
    .classify-layout {
        width: 65%;
        margin: 0 auto;
        background:linear-gradient(180deg,rgba(245,246,250,0.48) 0%,rgba(245,246,250,0.93) 100%);
    }
    .classify {
        .classify-name {
            padding: 9px 0;
            font-size: 14px;
            font-weight: bold;
            line-height: 19px;
            border-bottom: 1px solid #dde6ef;
        }
    }
    .model-list {
        display: flex;
        padding: 6px 0 25px;
        flex-wrap: wrap;
        .model-item {
            flex: 0 0 25%;
            margin: 6px 0;
            overflow: hidden;
            .model-info {
                display: flex;
                height: 32px;
                margin: 0 24px;
                border-radius: 2px;
                line-height: 32px;
                cursor: pointer;
                &:hover {
                    background-color: #e1ecff;
                    box-shadow:0px 2px 4px 0px rgba(51, 60, 72, 0.06);
                    color: #3a84ff;
                    .model-collect {
                        font-size: 16px;
                    }
                }
                &:nth-child(5n+1) {
                    margin: 0 24px 0 14px;
                }
            }
            .model-icon {
                flex: 0 0 42px;
                padding: 0 14px 0 12px;
                font-size: 16px;
                line-height: inherit;
                &.is-custom {
                    color: #3a84ff;
                }
            }
            .model-name {
                flex: 1;
                font-size: 14px;
                @include ellipsis;
            }
            .model-collect {
                flex: 0 0 40px;
                padding: 0 17px 0 7px;
                font-size: 0px;
                line-height: inherit;
                color: #979ba5;
                &.icon-star-shape {
                    font-size: 16px;
                    color: #ffb400;
                }
            }
        }
    }
</style>
