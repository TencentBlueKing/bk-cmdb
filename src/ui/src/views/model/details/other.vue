<template>
    <div class="other-wrapper">
        <div class="other-content">
            <h3>{{$t('ModelManagement["模型停用"]')}}</h3>
            <p>{{$t('ModelManagement["保留模型和相应实例，隐藏关联关系"]')}}</p>
            <div class="btn-wrapper">
                <template v-if="ispaused">
                    <bk-button type="primary" @click="dialogConfirm('restart')" :loading="$loading('updateModel')">
                        {{$t('ModelManagement["启用模型"]')}}
                    </bk-button>
                </template>
                <template v-else>
                    <bk-button type="primary" 
                    :loading="$loading('updateModel')"
                    :disabled="ispre || isMainLine"
                    @click="dialogConfirm('stop')">
                        {{$t('ModelManagement["停用模型"]')}}
                    </bk-button>
                </template>
                <span class="btn-tip-content" v-if="ispre">
                    <i class="icon-cc-attribute"></i>
                    <span class="btn-tip" :class="{'en': language !== 'zh_CN'}">
                        {{$t('ModelManagement["系统内建模型不可停用"]')}}
                    </span>
                </span>
            </div>
        </div>
        <div class="other-content">
            <h3>{{$t('ModelManagement["模型删除"]')}}</h3>
            <p>{{$t('ModelManagement["删除模型和其下所有实例，此动作不可逆，请谨慎操作"]')}}</p>
            <div class="btn-wrapper">
                <bk-button type="primary" :disabled="ispre" @click="dialogConfirm('delete')" :loading="$loading('deleteModel')">
                    {{$t('ModelManagement["删除模型"]')}}
                </bk-button>
                <span class="btn-tip-content" v-if="ispre">
                    <i class="icon-cc-attribute"></i>
                    <span class="btn-tip" :class="{'en': language !== 'zh_CN'}">
                        {{$t('ModelManagement["系统内建模型不可删除"]')}}
                    </span>
                </span>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    export default {
        computed: {
            ...mapGetters([
                'language'
            ]),
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            ispre () {
                return this.activeModel['ispre']
            },
            ispaused () {
                return this.activeModel['bk_ispaused']
            },
            isMainLine () {
                return this.activeModel['bk_classification_id'] === 'bk_biz_topo'
            }
        },
        methods: {
            ...mapActions('objectModel', [
                'updateObject',
                'deleteObject'
            ]),
            ...mapActions('objectMainLineModule', [
                'deleteMainlineObject'
            ]),
            dialogConfirm (type) {
                switch (type) {
                    case 'restart':
                        this.$bkInfo({
                            title: this.$t('ModelManagement["确认要启用该模型？"]'),
                            confirmFn: () => {
                                this.updateModel(false)
                            }
                        })
                        break
                    case 'stop':
                        this.$bkInfo({
                            title: this.$t('ModelManagement["确认要停用该模型？"]'),
                            confirmFn: () => {
                                this.updateModel(true)
                            }
                        })
                        break
                    case 'delete':
                        this.$bkInfo({
                            title: this.$t('ModelManagement["确认要删除该模型？"]'),
                            confirmFn: () => {
                                this.deleteModel()
                            }
                        })
                        break
                    default:
                }
            },
            async updateModel (ispaused) {
                await this.updateObject({
                    id: this.activeModel['id'],
                    params: {
                        bk_ispaused: ispaused
                    },
                    config: {
                        requestId: 'updateModel'
                    }
                })
                this.$emit('closeSlider')
            },
            async deleteModel () {
                if (this.isMainLine) {
                    await this.deleteMainlineObject({
                        bkObjId: this.activeModel['bk_obj_id'],
                        config: {
                            requestId: 'deleteModel'
                        }
                    })
                } else {
                    await this.deleteObject({
                        id: this.activeModel['id'],
                        config: {
                            requestId: 'deleteModel'
                        }
                    })
                }
                this.$emit('closeSlider')
            }
        }
    }
</script>


<style lang="scss" scoped>
    .other-wrapper {
        padding: 55px 35px;
        .other-content {
            &:first-child {
                margin-bottom: 50px;
            }
            >h3 {
                margin-bottom: 10px;
                padding-left: 4px;
                font-size: 14px;
                line-height: 1;
                border-left: 4px solid $cmdbTextColor;
            }
        }
        .btn-wrapper {
            margin-top: 20px;
            font-size: 0;
            .btn-tip-content {
                .icon-cc-attribute {
                    margin-left: 10px;
                    cursor: pointer;
                    color: #ffb400;
                    font-size: 16px;
                    &:hover {
                        +span {
                            display: inline-block;
                        }
                    }
                }
                .btn-tip {
                    display: none;
                    min-width: 170px;
                    height: 36px;
                    line-height: 36px;
                    text-align: center;
                    box-shadow: 0 0 5px #ebedef;
                    margin-left: 8px;
                    position:relative;
                    background: #333333;
                    color: #fff;
                    border-radius: 2px;
                    font-size: 12px;
                    vertical-align: middle;
                    &.en {
                        min-width: 300px;
                    }
                    &:before {
                        content: '';
                        position: absolute;
                        left: -7px;
                        top: 9px;
                        width: 0;
                        height: 0;
                        border-top: 7px solid transparent;
                        border-right: 7px solid #333;
                        border-bottom: 7px solid transparent;
                    }
                }
            }
        }
    }
</style>
