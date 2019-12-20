<template>
    <div class="property">
        <div class="group"
            v-for="(group, index) in groupedProperties"
            :key="index">
            <h2 class="group-name">{{group.bk_group_name}}</h2>
            <ul class="property-list">
                <li class="property-item"
                    v-for="property in group.properties"
                    :key="property.id">
                    <span class="property-name" v-overflow-tips>
                        {{property.bk_property_name}}
                    </span>
                    <span :class="['property-value', { 'is-loading': loadingState.includes(property) }]"
                        v-overflow-tips
                        v-show="property !== editState.property">
                        {{$tools.getPropertyText(property, host) | filterShowText(property.unit)}}
                    </span>
                    <template v-if="!loadingState.includes(property)">
                        <template v-if="isPropertyEditable(property)">
                            <cmdb-auth style="margin: 8px 0 0 8px; font-size: 0;" :auth="updateAuthResources" v-show="property !== editState.property">
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
                                <component class="form-component"
                                    :is="`cmdb-form-${property.bk_property_type}`"
                                    :class="[property.bk_property_type, { error: errors.has(property.bk_property_id) }]"
                                    :options="property.option || []"
                                    :data-vv-name="property.bk_property_id"
                                    :data-vv-as="property.bk_property_name"
                                    :placeholder="getPlaceholder(property)"
                                    :auto-check="false"
                                    v-validate="$tools.getValidateRules(property)"
                                    v-model.trim="editState.value"
                                    :ref="`component-${property.bk_property_id}`"
                                    @enter="confirm">
                                </component>
                                <i class="form-confirm bk-icon icon-check-1" @click="confirm"></i>
                                <i class="form-cancel bk-icon icon-close" @click="exitForm"></i>
                                <span class="form-error"
                                    v-if="errors.has(property.bk_property_id)">
                                    {{errors.first(property.bk_property_id)}}
                                </span>
                            </div>
                        </template>
                        <template v-if="$tools.getPropertyText(property, host) !== '--' && property !== editState.property">
                            <div class="copy-box">
                                <i class="property-copy icon-cc-details-copy" @click="handleCopy($tools.getPropertyText(property, host), property.bk_property_id)"></i>
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
    import { MENU_RESOURCE_HOST_DETAILS } from '@/dictionary/menu-symbol'
    export default {
        name: 'cmdb-host-property',
        filters: {
            filterShowText (value, unit) {
                return value === '--' ? '--' : value + unit
            }
        },
        data () {
            return {
                editState: {
                    property: null,
                    value: null
                },
                loadingState: [],
                showCopyTips: false
            }
        },
        computed: {
            ...mapState('hostDetails', ['info']),
            ...mapGetters('hostDetails', ['groupedProperties']),
            host () {
                return this.info.host || {}
            },
            updateAuthResources () {
                const isResourceHost = this.$route.name === MENU_RESOURCE_HOST_DETAILS
                if (isResourceHost) {
                    return this.$authResources({ type: this.$OPERATION.U_RESOURCE_HOST })
                }
                return this.$authResources({ type: this.$OPERATION.U_HOST })
            }
        },
        methods: {
            getPlaceholder (property) {
                const placeholderTxt = ['enum', 'list'].includes(property.bk_property_type) ? '请选择xx' : '请输入xx'
                return this.$t(placeholderTxt, { name: property.bk_property_name })
            },
            isPropertyEditable (property) {
                return property.editable && !property.bk_isapi
            },
            setEditState (property) {
                const value = this.host[property.bk_property_id]
                this.editState.value = value === null ? '' : value
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
                        params: this.$injectMetadata({
                            [property.bk_property_id]: value,
                            bk_host_id: String(this.host.bk_host_id)
                        }),
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
            handleCopy (copyText, propertyId) {
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
            &:hover {
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
                &:before {
                    display: inline-block;
                    transform: scale(0.83);
                }
            }
            &.form-cancel {
                color: #979ba5;
                font-size: 14px;
                &:before {
                    display: inline-block;
                    transform: scale(0.66);
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
            display: inline-block;
            vertical-align: middle;
            width: 260px;
            margin: 0 4px 0 0;
            &.bool {
                width: 42px;
            }
        }
    }
</style>
