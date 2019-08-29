<template>
    <div class="business-topo-wrapper" :style="{ 'padding-top': showFeatureTips ? '10px' : '' }">
        <feature-tips
            style="text-align: left;"
            :feature-name="'modelBusiness'"
            :show-tips="showFeatureTips"
            :desc="$t('业务层级提示')"
            :more-href="'https://docs.bk.tencent.com/cmdb/Introduction.html#%EF%BC%882%EF%BC%89%E6%96%B0%E5%A2%9E%E8%87%AA%E5%AE%9A%E4%B9%89%E5%B1%82%E7%BA%A7'"
            @close-tips="showFeatureTips = false">
        </feature-tips>
        <div class="topo-level" v-bkloading="{ isLoading: $loading() }">
            <div class="topo-node"
                v-for="(model, index) in topo"
                :style="{
                    marginLeft: `${index * margin}px`
                }"
                :key="index">
                <a href="javascript:void(0);" class="node-circle"
                    :class="{
                        'is-first': index === 0,
                        'is-last': index === (topo.length - 1),
                        'is-inner': innerModel.includes(model['bk_obj_id'])
                    }"
                    @click="handleLinkClick(model)">
                    <i :class="['icon', model['bk_obj_icon']]"></i>
                </a>
                <div class="node-name" :title="model['bk_obj_name']">{{model['bk_obj_name']}}</div>
                <a href="javascript:void(0)" class="node-add"
                    :class="{
                        disabled: !createAuth
                    }"
                    v-cursor="{
                        active: !createAuth,
                        auth: [$OPERATION.SYSTEM_TOPOLOGY]
                    }"
                    v-if="canAddLevel(model)"
                    @click="handleAddLevel(model)">
                </a>
            </div>
        </div>
        <the-create-model
            :is-show.sync="addLevel.showDialog"
            :is-main-line="true"
            :title="$t('新建层级')"
            @confirm="handleCreateLevel"
        ></the-create-model>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import theCreateModel from '@/components/model-manage/_create-model'
    import featureTips from '@/components/feature-tips/index'
    const NODE_MARGIN = 62

    export default {
        components: {
            theCreateModel,
            featureTips
        },
        data () {
            return {
                showFeatureTips: false,
                margin: NODE_MARGIN * 1.5,
                topo: [],
                innerModel: ['biz', 'set', 'module', 'host'],
                addLevel: {
                    showDialog: false,
                    parent: null
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'userName', 'isAdminView', 'featureTipsParams']),
            ...mapGetters('objectModelClassify', ['models']),
            createAuth () {
                return this.$isAuthorized(this.$OPERATION.SYSTEM_TOPOLOGY)
            }
        },
        created () {
            this.showFeatureTips = this.featureTipsParams['modelBusiness']
            this.getMainLineModel()
        },
        methods: {
            ...mapActions('objectMainLineModule', [
                'searchMainlineObject',
                'createMainlineObject'
            ]),
            ...mapActions('objectModelClassify', [
                'searchClassificationsObjects'
            ]),
            async getMainLineModel () {
                try {
                    const topo = await this.searchMainlineObject({})
                    this.topo = topo.map(model => {
                        return {
                            ...model,
                            'bk_obj_icon': this.getModelIcon(model['bk_obj_id'])
                        }
                    })
                } catch (e) {
                    this.topo = []
                    console.log(e)
                }
            },
            getModelIcon (objId) {
                const model = this.models.find(model => model['bk_obj_id'] === objId)
                return (model || {})['bk_obj_icon']
            },
            canAddLevel (model) {
                return this.isAdminView && !['set', 'module', 'host'].includes(model['bk_obj_id'])
            },
            handleAddLevel (model) {
                if (this.createAuth) {
                    this.addLevel.parent = model
                    this.addLevel.showDialog = true
                }
            },
            async handleCreateLevel (data) {
                try {
                    await this.createMainlineObject({
                        params: this.$injectMetadata({
                            'bk_asst_obj_id': this.addLevel.parent['bk_obj_id'],
                            'bk_classification_id': 'bk_biz_topo',
                            'bk_obj_icon': data['bk_obj_icon'],
                            'bk_obj_id': data['bk_obj_id'],
                            'bk_obj_name': data['bk_obj_name'],
                            'bk_supplier_account': this.supplierAccount,
                            'creator': this.userName
                        })
                    })
                    await this.searchClassificationsObjects({
                        params: this.$injectMetadata({}),
                        config: {
                            clearCache: true,
                            requestId: 'post_searchClassificationsObjects'
                        }
                    })
                    this.getMainLineModel()
                    this.handleCancelCreateLevel()
                } catch (e) {
                    console.log(e)
                }
            },
            handleCancelCreateLevel () {
                this.addLevel.parent = null
                this.addLevel.showDialog = false
            },
            handleLinkClick (model) {
                this.$router.push({
                    name: 'modelDetails',
                    params: {
                        modelId: model.bk_obj_id
                    },
                    query: {
                        from: this.$route.fullPath
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .business-topo-wrapper {
        height: 100%;
        background-color: #f4f5f8;
        background-image: linear-gradient(#eef1f5 1px, transparent 0), linear-gradient(90deg, #eef1f5 1px, transparent 0);
        background-size: 10px 10px;
        text-align: center;
        &:after {
            content: "";
            display: inline-block;
            vertical-align: middle;
            width: 0;
            height: calc(100% - 20px);
        }
    }
    .topo-level {
        display: inline-block;
        vertical-align: middle;
    }
    .topo-node {
        position: relative;
        width: 62px;
        margin-top: 8px;
        .node-circle {
            position: relative;
            display: inline-block;
            width: 62px;
            height: 62px;
            line-height: 62px;
            background: #fff;
            box-shadow: 0px 2px 4px 0px rgba(147,147,147,0.5);
            border-radius: 50%;
            font-size: 0;
            color: #3c96ff;
            &.is-inner {
                color: #868b97;
            }
            &:before {
                content: "";
                position: absolute;
                right: 100%;
                top: 50%;
                width: 62px;
                height: 0;
                border-top: 2px dashed $cmdbBorderColor;
                pointer-events: none;
            }
            &:after {
                content: "";
                position: absolute;
                top: 100%;
                left: 50%;
                width: 0;
                height: 40px;
                margin: 0 0 0 -1px;
                border-right: 2px dashed $cmdbBorderColor;
                pointer-events: none;
            }
            &.is-first:before,
            &.is-last:after {
                display: none;
            }
            .icon {
                font-size: 24px;
            }
        }
        .node-name {
            position: absolute;
            width: 150px;
            top: 100%;
            left: 0;
            transform: translateX(-44px);
            font-size: 14px;
            @include ellipsis;
        }
        .node-add {
            position: absolute;
            top: 94px;
            left: 50%;
            width: 16px;
            height: 16px;
            margin: 0 0 0 -8px;
            border-radius: 2px;
            background-color: #3c96ff;
            z-index: 1;
            &:before,
            &:after {
                content: "";
                position: absolute;
                background-color: #fff;
            }
            &:before {
                left: 4px;
                top: 7px;
                width: 8px;
                height: 2px;
            }
            &:after {
                left: 7px;
                top: 4px;
                width: 2px;
                height: 8px;
            }
            &.disabled {
                background-color: #ccc;
            }
        }
    }
</style>
