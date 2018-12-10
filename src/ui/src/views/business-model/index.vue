<template>
    <div class="business-topo-wrapper">
        <div class="topo-level" v-bkloading="{isLoading: $loading()}">
            <div class="topo-node" 
                v-for="(model, index) in topo"
                :style="{
                    marginLeft: `${index * margin}px`
                }"
                :key="index">
                <router-link :to="`/model/details/${model['bk_obj_id']}`" class="node-circle" 
                    :class="{
                        'is-first': index === 0,
                        'is-last': index === (topo.length - 1),
                        'is-inner': innerModel.includes(model['bk_obj_id'])
                    }">
                    <i :class="['icon', model['bk_obj_icon']]"></i>
                </router-link>
                <div class="node-name" :title="model['bk_obj_name']">{{model['bk_obj_name']}}</div>
                <a href="javascript:void(0)" class="node-add"
                    v-if="canAddLevel(model)"
                    @click="handleAddLevel(model)">
                </a>
            </div>
        </div>
        <bk-dialog
            :is-show.sync="addLevel.showDialog" 
            :close-icon="false" 
            :has-header="false"
            :width="600"
            :padding="0"
            @confirm="handleCreateLevel"
            @cancel="handleCancelCreateLevel">
            <div class="add-level-wrapper" slot="content">
                <h2 class="add-level-title">{{$t('ModelManagement["新建层级"]')}}</h2>
                <div class="add-level-form clearfix">
                    <a href="javascript:void(0)" class="add-level-icon fl" @click="addLevel.showIconSelector = true">
                        <i :class="['icon', addLevel.icon]"></i>
                        <span class="text">{{$t('ModelManagement["点击切换"]')}}</span>
                    </a>
                    <div class="add-level-info">
                        <label class="label">{{$t('ModelManagement["唯一标识"]')}}</label>
                        <input type="text" class="input cmdb-form-input" :placeholder="$t('ModelManagement[\'请输入英文标识\']')"
                            name="enName"
                            v-model.trim="addLevel.enName"
                            v-validate="'required|modelId'" />
                        <span class="error">{{errors.first('enName')}}</span>
                    </div>
                    <div class="add-level-info" style="margin-top: 10px">
                        <label class="label">{{$t('ModelManagement["名称"]')}}</label>
                        <input type="text" class="input cmdb-form-input" :placeholder="$t('ModelManagement[\'请输入名称\']')"
                            name="name"
                            v-model.trim="addLevel.name"
                            v-validate="'required|singlechar'" />
                        <span class="error">{{errors.first('name')}}</span>
                    </div>
                </div>
                <the-choose-icon class="icon-selector"
                    v-if="addLevel.showIconSelector"
                    v-model="addLevel.icon"
                    @chooseIcon="addLevel.showIconSelector = false">
                </the-choose-icon>
                <span class="back"
                    v-show="addLevel.showIconSelector"
                    @click="addLevel.showIconSelector = false">
                    <i class="bk-icon icon-back2"></i>
                </span>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import theChooseIcon from '@/components/model-manage/_choose-icon'

    const NODE_MARGIN = 62

    export default {
        components: {
            theChooseIcon
        },
        data () {
            return {
                margin: NODE_MARGIN * 1.5,
                topo: [],
                innerModel: ['biz', 'set', 'module', 'host'],
                addLevel: {
                    showDialog: false,
                    showIconSelector: false,
                    icon: 'icon-cc-default',
                    name: '',
                    enName: '',
                    parent: null
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'userName', 'admin']),
            authority () {
                return this.admin ? ['search', 'update', 'delete'] : []
            }
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t('Nav["业务模型"]'))
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
                    const topo = await this.searchMainlineObject()
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
                const model = this.$allModels.find(model => model['bk_obj_id'] === objId)
                return (model || {})['bk_obj_icon']
            },
            canAddLevel (model) {
                return this.authority.includes('update') && !['set', 'module', 'host'].includes(model['bk_obj_id'])
            },
            handleAddLevel (model) {
                this.addLevel.parent = model
                this.addLevel.showDialog = true
            },
            async handleCreateLevel () {
                const valid = await this.$validator.validateAll()
                if (!valid) {
                    return false
                }
                try {
                    await this.createMainlineObject({
                        params: {
                            'bk_asst_obj_id': this.addLevel.parent['bk_obj_id'],
                            'bk_classification_id': 'bk_biz_topo',
                            'bk_obj_icon': this.addLevel.icon,
                            'bk_obj_id': this.addLevel.enName,
                            'bk_obj_name': this.addLevel.name,
                            'bk_supplier_account': this.supplierAccount,
                            'creator': this.userName
                        }
                    })
                    await this.searchClassificationsObjects({
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
                this.addLevel.name = ''
                this.addLevel.enName = ''
                this.addLevel.parent = null
                this.addLevel.icon = 'icon-cc-default'
                this.addLevel.showDialog = false
                this.addLevel.showIconSelector = false
                this.$validator.reset()
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
            height: 100%;
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
            font-size: 24px;
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
            }
            &.is-first:before,
            &.is-last:after {
                display: none;
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
        }
    }
    .add-level-wrapper {
        text-align: left;
        padding: 23px 13px;
        position: relative;
        .add-level-title {
            font-size: 20px;
            line-height: 26px;
        }
        .add-level-form {
            margin: 30px 12px;
        }
        .add-level-icon {
            display: block;
            width: 93px;
            border-radius: 4px;
            border: 1px solid $cmdbBorderColor;
            text-align: center;
            .icon {
                display: block;
                height: 70px;
                line-height: 70px;
                font-size: 38px;
                color: #3c96ff;
            }
            .text {
                display: block;
                height: 30px;
                line-height: 30px;
                font-size: 12px;
                border-top: 1px solid $cmdbBorderColor;
                background-color: #ebf4ff;
                border-radius: 0 0 4px 4px;
            }
        }
        .add-level-info {
            position: relative;
            padding: 5px 0;
            margin: 0 0 0 93px;
            font-size: 0;
            .label {
                display: inline-block;
                width: 100px;
                padding: 0 4px;
                font-size: 16px;
                line-height: 36px;
                vertical-align: middle;
                text-align: right;
                &:after {
                    display: inline-block;
                    content: "*";
                    color: $cmdbDangerColor;
                }
            }
            .input {
                width: 330px;
                font-size: 16px;
                vertical-align: middle;
            }
            .error {
                position: absolute;
                top: 40px;
                left: 100px;
                color: $cmdbDangerColor;
                font-size: 12px;
            }
        }
        .icon-selector {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: calc(100% + 60px);
            background-color: #fff;
        }
        .back {
            position: absolute;
            left: 100%;
            top: 0;
            width: 44px;
            height: 44px;
            padding: 7px;
            margin: 0 0 0 3px;
            cursor: pointer;
            font-size: 18px;
            text-align: center;
            background: #2f2f2f;
            color: #fff;
        }
    }
</style>