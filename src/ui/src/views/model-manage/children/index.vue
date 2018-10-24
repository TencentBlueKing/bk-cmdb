<template>
    <div class="model-detail-wrapper">
        <div class="model-info">
            <div class="icon-box">
                <i class="icon icon-cc-host"></i>
            </div>
            <div class="model-text">
                <span>{{$t('ModelManagement["唯一标识"]')}}：</span>
                <span class="text-content name">asdf</span>
            </div>
            <div class="model-text">
                <span>{{$t('Hosts["名称"]')}}：</span>
                <template v-if="true">
                    <span class="text-content">asdf<i class="icon icon-cc-edit text-primary"></i></span>
                </template>
                <template v-else>
                    <input type="text" class="cmdb-form-input">
                    <span class="text-primary">{{$t("Common['保存']")}}</span>
                    <span class="text-primary">{{$t("Common['取消']")}}</span>
                </template>
            </div>
            <div class="btn-group">
                <label class="label-btn">
                    <i class="icon-cc-derivation"></i>
                    <span>{{$t('ModelManagement["导出"]')}}</span>
                </label>
                <label class="label-btn">
                    <i class="icon-cc-copy"></i>
                    <span>{{$t('Common["复制"]')}}</span>
                </label>
                <label class="label-btn">
                    <i class="bk-icon icon-minus-circle-shape"></i>
                    <span>{{$t('ModelManagement["停用"]')}}</span>
                </label>
                <label class="label-btn">
                    <i class="icon-cc-del"></i>
                    <span>{{$t("Common['删除']")}}</span>
                </label>
            </div>
        </div>
        <bk-tab :active-name.sync="tab.active">
            <bk-tabpanel name="field" :title="$t('ModelManagement[\'模型字段\']')">
                <the-field></the-field>
            </bk-tabpanel>
            <bk-tabpanel name="relation" :title="$t('ModelManagement[\'模型关系\']')">
                <the-relation></the-relation>
            </bk-tabpanel>
            <bk-tabpanel name="verification" :title="$t('ModelManagement[\'唯一校验\']')">
            </bk-tabpanel>
            <bk-tabpanel name="layout" :title="$t('ModelManagement[\'字段分组\']')">
            </bk-tabpanel>
            <bk-tabpanel name="history" :title="$t('ModelManagement[\'操作历史\']')">
            </bk-tabpanel>
        </bk-tab>
    </div>
</template>

<script>
    import theField from './field'
    import theRelation from './relation'
    import { mapActions, mapGetters } from 'vuex'
    export default {
        components: {
            theField,
            theRelation
        },
        data () {
            return {
                tab: {
                    active: 'field'
                }
            }
        },
        computed: {
            ...mapGetters([
                'supplierAccount'
            ])
        },
        watch: {
            '$route.params.modelId' () {
                this.initObject()
            }
        },
        created () {
            this.initObject()
        },
        methods: {
            ...mapActions('objectModel', [
                'searchObjects'
            ]),
            async initObject () {
                const res = await this.searchObjects({
                    params: {
                        bk_obj_id: this.$route.params.modelId,
                        bk_supplier_account: this.supplierAccount
                    }
                })
                this.$store.commit('objectModel/setActiveModel', res[0])
            }
        }
    }
</script>

<style lang="scss" scoped>
    .model-detail-wrapper {
        padding: 0;
    }
    .model-info {
        padding: 0 24px 0 38px;
        height: 100px;
        background: rgba(235, 244, 255, .6);
        font-size: 14px;
        .icon-box {
            float: left;
            margin-top: 14px;
            margin: 14px 30px 0 0;
            padding-top: 20px;
            width: 72px;
            height: 72px;
            border: 1px solid #dde4eb;
            border-radius: 50%;
            background: #fff;
            text-align: center;
            font-size: 32px;
            color: $cmdbBorderFocusColor;
            .icon {
                vertical-align: top;
            }
        }
        .model-text {
            float: left;
            margin: 32px 10px 32px 0;
            line-height: 36px;
            >span {
                display: inline-block;
                vertical-align: top;
            }
            .text-content {
                max-width: 110px;
                @include ellipsis;
                &.name {
                    width: 110px;
                }
                .icon {
                    margin-top: -4px;
                    margin: -4px 0 0 4px;
                }
            }
            .cmdb-form-input {
                display: inline-block;
                width: 200px;
                vertical-align: top;
            }
            .text-primary {
                cursor: pointer;
                margin-left: 5px;
            }
        }
        .btn-group {
            float: right;
            height: 100px;
            line-height: 100px;
            i,
            span {
                vertical-align: middle;
            }
        }
    }
</style>

<style lang="scss">
    @import '@/assets/scss/model-manage.scss';
</style>
