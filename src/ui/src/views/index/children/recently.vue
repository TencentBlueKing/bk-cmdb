<template>
    <div class="recently-layout">
        <ul class="model-list">
            <li class="model-item"
                v-for="(model, index) in displayModels"
                :key="index"
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
                {{model.bk_obj_name}}
            </li>
            <li class="model-item more"
                @click="toggleClassify">
                <i class="more-icon bk-icon icon-angle-down"
                    :class="{
                        'is-open': showClassify
                    }">
                </i>
            </li>
        </ul>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        computed: {
            ...mapGetters(['isAdminView']),
            ...mapGetters('index', ['showClassify']),
            ...mapGetters('userCustom', ['usercustom']),
            ...mapGetters('objectModelClassify', ['models']),
            recentlyModels () {
                const usercustomData = this.usercustom.recently_models || []
                const recentlyModels = []
                const basicModels = ['resource', 'business']
                usercustomData.forEach(id => {
                    if (basicModels.includes(id)) {
                        console.warn(`${id} model has been ignored in history.`)
                    } else {
                        const model = this.models.find(model => model.id === id)
                        if (model && !model.bk_ispaused) {
                            recentlyModels.push(model)
                        }
                    }
                })
                return recentlyModels
            },
            avaliableModels () {
                return this.models.filter(model => {
                    return !model.bk_ispaused
                        && !['bk_host_manage', 'bk_biz_topo', 'bk_organization'].includes(model.bk_classification_id)
                })
            },
            displayModels () {
                const displayModels = []
                const allModels = [...this.recentlyModels, ...this.avaliableModels]
                allModels.forEach(model => {
                    if (!displayModels.some(target => target.id === model.id)) {
                        displayModels.push(model)
                    }
                })
                return displayModels.slice(0, 5)
            }
        },
        methods: {
            isCustomModel (model) {
                if (this.isAdminView) {
                    return model.ispre
                }
                return !!this.$tools.getMetadataBiz(model)
            },
            handleClick (model) {
                const router = model.router || {
                    name: 'generalModel',
                    params: {
                        objId: model.bk_obj_id
                    },
                    query: {
                        from: this.$route.fullPath
                    }
                }
                this.$router.push(router)
            },
            toggleClassify () {
                this.$store.commit('index/toggleClassify', !this.showClassify)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .recently-layout {
        width: 50%;
        margin: 0 auto;
    }
    .model-list {
        display: flex;
        margin: 27px 0 0;
        color: #979ba5;
        overflow: hidden;
        .model-item {
            flex: 0 0 calc((100% - 34px) * 0.2 - 15px);
            height: 26px;
            padding: 0 13px;
            margin: 0 15px 0 0;
            line-height: 26px;
            border-radius: 15px;
            background-color: #fff;
            font-size: 12px;
            cursor: pointer;
            box-shadow:0px 2px 4px 0px rgba(51,60,72,0.06);
            @include ellipsis;
            &:hover {
                background-color: #e1ecff;
                color: #3a84ff;
                .model-icon {
                    color: #3a84ff;
                }
            }
            &.more {
                flex: 0 0 34px;
                margin: 0;
                padding: 0;
                font-size: 0;
                text-align: center;
                .more-icon {
                    font-size: 12px;
                    font-weight: bold;
                    line-height: inherit;
                    transition: transform .2s linear;
                    &.is-open {
                        transform: rotate(180deg);
                    }
                }
                &:hover .more-icon {
                    color: #0082ff;
                }
            }
            .model-icon {
                vertical-align: text-top;
                font-size: 14px;
                margin: 2px 4px 0 0;
                color: #798aad;
                &.is-custom {
                    color: #3a84ff;
                }
            }
        }
    }
</style>
