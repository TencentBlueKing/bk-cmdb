<template>
    <div class="property">
        <div class="group"
            v-for="(group, index) in groupedProperties"
            :key="index">
            <h2 class="group-name">{{group.bk_group_name}}</h2>
            <ul class="property-list">
                <li class="property-item"
                    v-for="property in group.properties"
                    :key="property.id"
                    :id="`property-item-${property.id}`">
                    <span class="property-name" v-bk-overflow-tips>
                        {{property.bk_property_name}}
                    </span>
                    <span :class="['property-value', { 'is-loading': loadingState.includes(property) }]"
                        v-bk-overflow-tips
                        v-if="property !== editState.property">
                        <cmdb-property-value
                            :ref="`property-value-${property.bk_property_id}`"
                            :value="host[property.bk_property_id]"
                            :property="property">
                        </cmdb-property-value>
                    </span>
                    <template v-if="!loadingState.includes(property)">
                        <template v-if="hasRelatedRules(property) || !isPropertyEditable(property)">
                            <span :id="`rule-${property.id}`">
                                <i18n path="已配置属性自动应用提示" v-if="hasRelatedRules(property)">
                                    <bk-button text place="link" @click="handleViewRules(property)">{{$t('点击跳转查看配置详情')}}</bk-button>
                                </i18n>
                                <span v-else>{{$t('系统限定不可修改')}}</span>
                            </span>
                            <i class="is-related property-edit icon-cc-edit"
                                v-bk-tooltips="{
                                    allowHtml: true,
                                    content: `#rule-${property.id}`,
                                    placement: 'top',
                                    onShow: () => {
                                        setFocus(`#property-item-${property.id}`, true)
                                    },
                                    onHide: () => {
                                        setFocus(`#property-item-${property.id}`, false)
                                    }
                                }">
                            </i>
                        </template>
                        <template v-else>
                            <cmdb-auth style="margin: 8px 0 0 8px; font-size: 0;" :auth="HOST_AUTH.U_HOST" v-show="property !== editState.property">
                                <bk-button slot-scope="{ disabled }"
                                    text
                                    theme="primary"
                                    class="property-edit-btn"
                                    :disabled="disabled"
                                    @click="setEditState(property)">
                                    <i class="property-edit icon-cc-edit"></i>
                                </bk-button>
                            </cmdb-auth>
                            <div class="property-form" v-if="property === editState.property">
                                <div :class="['form-component', property.bk_property_type]">
                                    <component
                                        :is="`cmdb-form-${property.bk_property_type}`"
                                        :class="[property.bk_property_type, { error: errors.has(property.bk_property_id) }]"
                                        :unit="property.unit"
                                        :options="property.option || []"
                                        :data-vv-name="property.bk_property_id"
                                        :data-vv-as="property.bk_property_name"
                                        :placeholder="getPlaceholder(property)"
                                        :auto-check="false"
                                        v-validate="$tools.getValidateRules(property)"
                                        v-model.trim="editState.value"
                                        :ref="`component-${property.bk_property_id}`">
                                    </component>
                                </div>
                                <i class="form-confirm bk-icon icon-check-1" @click="confirm"></i>
                                <i class="form-cancel bk-icon icon-close" @click="exitForm"></i>
                                <span class="form-error"
                                    v-if="errors.has(property.bk_property_id)">
                                    {{errors.first(property.bk_property_id)}}
                                </span>
                            </div>
                        </template>
                        <template v-if="host[property.bk_property_id] && property !== editState.property">
                            <div class="copy-box">
                                <i class="property-copy icon-cc-details-copy" @click="handleCopy(property.bk_property_id)"></i>
                                <transition name="fade">
                                    <span class="copy-tips"
                                        :style="{ width: $i18n.locale === 'en' ? '100px' : '70px' }"
                                        v-if="showCopyTips === property.bk_property_id">
                                        {{$t('复制成功')}}
                                    </span>
                                </transition>
                            </div>
                        </template>
                    </template>
                </li>
            </ul>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapState } from 'vuex'
    import { MENU_BUSINESS_HOST_APPLY } from '@/dictionary/menu-symbol'
    import authMixin from '../mixin-auth'
    export default {
        name: 'cmdb-host-property',
        filters: {
            filterShowText (value, unit) {
                return value === '--' ? '--' : value + unit
            }
        },
        mixins: [authMixin],
        data () {
            return {
                editState: {
                    property: null,
                    value: null
                },
                loadingState: [],
                showCopyTips: false,
                hostRelatedRules: [],
                request: {
                    rules: Symbol('rules')
                }
            }
        },
        computed: {
            ...mapState('hostDetails', ['info']),
            ...mapGetters('hostDetails', ['groupedProperties']),
            host () {
                return this.info.host || {}
            }
        },
        watch: {
            host () {
                this.getHostRelatedRules()
            }
        },
        methods: {
            setFocus (id, focus) {
                const item = this.$el.querySelector(id)
                focus ? item.classList.add('focus') : item.classList.remove('focus')
            },
            handleViewRules (property) {
                const rule = this.hostRelatedRules.find(rule => rule.bk_attribute_id === property.id) || {}
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_HOST_APPLY,
                    query: {
                        module: rule.bk_module_id
                    },
                    history: true
                })
            },
            hasRelatedRules (property) {
                return this.hostRelatedRules.some(rule => rule.bk_attribute_id === property.id)
            },
            async getHostRelatedRules () {
                try {
                    const defaultType = this.$tools.getValue(this.info, 'biz.0.default')
                    if (defaultType) { // 为0时非业务模块，不存在属性自动应用
                        this.hostRelatedRules = []
                    } else {
                        const data = await this.$store.dispatch('hostApply/getHostRelatedRules', {
                            bizId: this.$tools.getValue(this.info, 'biz.0.bk_biz_id'),
                            params: {
                                bk_host_ids: [this.host.bk_host_id]
                            },
                            config: {
                                requestId: this.request.rules
                            }
                        })
                        this.hostRelatedRules = data[this.host.bk_host_id] || []
                    }
                } catch (e) {
                    this.hostRelatedRules = []
                    console.error(e)
                }
            },
            getPlaceholder (property) {
                const placeholderTxt = ['enum', 'list', 'organization'].includes(property.bk_property_type) ? '请选择xx' : '请输入xx'
                return this.$t(placeholderTxt, { name: property.bk_property_name })
            },
            isPropertyEditable (property) {
                const isSystemLimited = property.editable && !property.bk_isapi
                // bk_cloud_inst_id 有值得为云主机，云主机的外网IP不可编辑
                const isCloudHost = property.bk_property_id === 'bk_host_outerip' && this.host.bk_cloud_inst_id
                return isSystemLimited || isCloudHost
            },
            setEditState (property) {
                const value = this.host[property.bk_property_id]
                this.editState.value = (value === null || value === undefined) ? '' : value
                this.editState.property = property
                this.$nextTick(() => {
                    const component = this.$refs[`component-${property.bk_property_id}`]
                    component[0] && component[0].focus && component[0].focus()
                })
            },
            async confirm () {
                const { property, value } = this.editState
                try {
                    const isValid = await this.$validator.validateAll()
                    if (!isValid) {
                        return false
                    }
                    this.loadingState.push(property)
                    this.exitForm()
                    await this.$store.dispatch('hostUpdate/updateHost', {
                        params: {
                            [property.bk_property_id]: value,
                            bk_host_id: String(this.host.bk_host_id)
                        },
                        config: {
                            requestId: 'updateHostInfo'
                        }
                    })
                    this.$store.commit('hostDetails/updateInfo', {
                        [property.bk_property_id]: value
                    })
                    this.loadingState = this.loadingState.filter(exist => exist !== property)
                } catch (e) {
                    console.error(e)
                    this.loadingState = this.loadingState.filter(exist => exist !== property)
                }
            },
            exitForm () {
                this.editState.property = null
                this.editState.value = null
            },
            handleCopy (propertyId) {
                const component = this.$refs[`property-value-${propertyId}`]
                const copyText = component[0] ? component[0].$el.innerText : ''
                this.$copyText(copyText).then(() => {
                    this.showCopyTips = propertyId
                    const timer = setTimeout(() => {
                        this.showCopyTips = false
                        clearTimeout(timer)
                    }, 200)
                }, () => {
                    this.$error(this.$t('复制失败'))
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .property {
        height: 100%;
        overflow: auto;
        @include scrollbar-y;
    }
    .group {
        margin: 22px 0 0 0;
        .group-name {
            line-height: 21px;
            font-size: 16px;
            font-weight: normal;
            color: #333948;
            &:before {
                content: "";
                display: inline-block;
                vertical-align: -2px;
                width: 4px;
                height: 14px;
                margin-right: 9px;
                background-color: $cmdbBorderColor;
            }
        }
    }
    .property-list {
        width: 1000px;
        margin: 25px 0 0 0;
        color: #63656e;
        display: flex;
        flex-wrap: wrap;
        .property-item {
            flex: 0 0 50%;
            max-width: 50%;
            padding-bottom: 8px;
            display: flex;
            &:hover,
            &.focus {
                .property-edit {
                    opacity: 1;
                }
                .property-copy {
                    display: inline-block;
                }
            }
            .property-name {
                position: relative;
                width: 160px;
                line-height: 32px;
                padding: 0 16px 0 36px;
                font-size: 14px;
                color: #63656E;
                @include ellipsis;
                &:after {
                    position: absolute;
                    right: 2px;
                    content: "：";
                }
            }
            .property-value {
                margin: 6px 0 0 4px;
                max-width: 286px;
                font-size: 14px;
                color: #313237;
                overflow:hidden;
                text-overflow:ellipsis;
                word-break: break-all;
                display: -webkit-box;
                -webkit-line-clamp: 2;
                -webkit-box-orient: vertical;
                &.is-loading {
                    font-size: 0;
                    &:before {
                        content: "";
                        display: inline-block;
                        width: 16px;
                        height: 16px;
                        margin: 2px 0;
                        background-image: url("../../../assets/images/icon/loading.svg");
                    }
                }
            }
            .property-edit-btn {
                height: auto;
                font-size: 0;
            }
            .property-edit {
                font-size: 16px;
                opacity: 0;
                &.is-related {
                    display: inline-block;
                    vertical-align: middle;
                    width: 16px;
                    height: 16px;
                    margin: 8px 0 0 8px;
                    line-height: 1;
                }
                &:hover {
                    opacity: .8;
                }
            }
            .property-copy {
                margin: 8px 0 0 8px;
                color: #3c96ff;
                cursor: pointer;
                display: none;
                font-size: 16px;
            }
            .copy-box {
                position: relative;
                font-size: 0;
                .copy-tips {
                    position: absolute;
                    top: -22px;
                    left: -18px;
                    min-width: 70px;
                    height: 26px;
                    line-height: 26px;
                    font-size: 12px;
                    color: #ffffff;
                    text-align: center;
                    background-color: #9f9f9f;
                    border-radius: 2px;
                }
                .fade-enter-active, .fade-leave-active {
                    transition: all 0.5s;
                }
                .fade-enter {
                    top: -14px;
                    opacity: 0;
                }
                .fade-leave-to {
                    top: -28px;
                    opacity: 0;
                }
            }
        }
    }
    .property-form {
        font-size: 0;
        position: relative;
        .bk-icon {
            display: inline-block;
            vertical-align: middle;
            width: 32px;
            height: 32px;
            margin: 0 0 0 6px;
            border-radius: 2px;
            border: 1px solid #c4c6cc;
            line-height: 30px;
            font-size: 12px;
            text-align: center;
            cursor: pointer;
            &.form-confirm {
                color: #0082ff;
                font-size: 20px;
                &:before {
                    display: inline-block;
                }
            }
            &.form-cancel {
                color: #979ba5;
                font-size: 20px;
                &:before {
                    display: inline-block;
                }
            }
            &:hover {
                font-weight: bold;
            }
        }
        .form-error {
            position: absolute;
            top: 100%;
            left: 0;
            font-size: 12px;
            line-height: 1;
            color: $cmdbDangerColor;
        }
        .form-component {
            display: inline-flex;
            align-items: center;
            vertical-align: middle;
            height: 32px;
            width: 260px;
            margin: 0 4px 0 0;
            &.bool {
                width: 42px;
                height: 24px;
            }
            &.longchar {
                height: auto;
                vertical-align: top;
            }
        }
    }
</style>
