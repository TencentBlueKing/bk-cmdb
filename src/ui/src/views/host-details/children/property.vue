<template>
    <div class="property" v-bkloading="{
        isLoading: $loading('updateHostInfo')
    }">
        <div class="group"
            v-for="(group, index) in groupedProperties"
            :key="index">
            <h2 class="group-name">{{group.bk_group_name}}</h2>
            <ul class="property-list">
                <li class="property-item"
                    v-for="property in group.properties"
                    :key="property.id">
                    <span class="property-name"
                        :title="property.bk_property_name">
                        {{property.bk_property_name}}
                    </span>
                    <bk-popover class="property-popover"
                        placement="bottom"
                        :disabled="!tooltipState[property.bk_property_id]"
                        :delay="300"
                        :offset="-5">
                        <span class="property-value"
                            v-show="property !== editState.property"
                            @mouseover="handleHover($event, property)">
                            {{$tools.getPropertyText(property, host) | filterShowText(property.unit)}}
                        </span>
                        <span class="popover-content" slot="content">
                            {{$tools.getPropertyText(property, host) | filterShowText(property.unit)}}
                        </span>
                    </bk-popover>
                    <template v-if="isPropertyEditable(property)">
                        <i class="property-edit icon-cc-edit"
                            v-if="$isAuthorized(updateAuth)"
                            v-show="property !== editState.property"
                            @click="setEditState(property)">
                        </i>
                        <i class="property-edit icon-cc-edit disabled"
                            v-else
                            v-cursor="{
                                active: true,
                                auth: [updateAuth]
                            }">
                        </i>
                        <div class="property-form" v-if="property === editState.property">
                            <component class="form-component"
                                :is="`cmdb-form-${property.bk_property_type}`"
                                :class="[property.bk_property_type, { error: errors.has(property.bk_property_id) }]"
                                :options="property.option || []"
                                :data-vv-name="property.bk_property_id"
                                :data-vv-as="property.bk_property_name"
                                :placeholder="$t('请输入xx', { name: property.bk_property_name })"
                                v-bind="$tools.getValidateEvents(property)"
                                v-validate="$tools.getValidateRules(property)"
                                v-model.trim="editState.value"
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
                tooltipState: {},
                showCopyTips: false
            }
        },
        computed: {
            ...mapState('hostDetails', ['info']),
            ...mapGetters('hostDetails', ['groupedProperties']),
            host () {
                return this.info.host || {}
            },
            updateAuth () {
                const isResourceHost = this.$route.name === MENU_RESOURCE_HOST_DETAILS
                if (isResourceHost) {
                    return this.$OPERATION.U_RESOURCE_HOST
                }
                return this.$OPERATION.U_HOST
            }
        },
        methods: {
            isPropertyEditable (property) {
                return property.editable && !property.bk_isapi
            },
            setEditState (property) {
                const value = this.host[property.bk_property_id]
                this.editState.value = value === null ? '' : value
                this.editState.property = property
            },
            async confirm () {
                try {
                    const isValid = await this.$validator.validateAll()
                    if (!isValid) {
                        return false
                    }
                    const { property, value } = this.editState
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
                    this.exitForm()
                } catch (e) {
                    console.error(e)
                }
            },
            exitForm () {
                this.editState.property = null
                this.editState.value = null
            },
            handleHover (event, property) {
                const target = event.target
                const range = document.createRange()
                range.selectNode(target)
                const rangeWidth = range.getBoundingClientRect().width
                const threshold = Math.max(rangeWidth - target.offsetWidth, target.scrollWidth - target.offsetWidth)
                this.$set(this.tooltipState, property.bk_property_id, threshold > 0.5)
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
                width: 150px;
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
                max-width: 296px;
                font-size: 14px;
                color: #313237;
                overflow:hidden;
                text-overflow:ellipsis;
                word-break: break-all;
                display: -webkit-box;
                -webkit-line-clamp: 2;
                -webkit-box-orient: vertical;
            }
            .property-edit {
                opacity: 0;
                margin: 8px 0 0 8px;
                font-size: 16px;
                color: #3c96ff;
                cursor: pointer;
                &:hover {
                    opacity: .8;
                }
                &.disabled {
                    color: #ccc;
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
    .property-popover {
        display: inline-block;
        /deep/ .bk-tooltip-ref {
            outline: none;
        }
    }
    .popover-content {
        display: inline-block;
        max-width: 300px;
        white-space: normal;
        word-break: break-all;
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
            width: 270px;
            margin: 0 4px 0 0;
            &.bool {
                width: 42px;
            }
        }
    }
</style>
