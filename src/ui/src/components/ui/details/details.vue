<template>
    <div class="details-layout">
        <div ref="detailsWrapper">
            <slot name="details-header"></slot>
            <template v-for="(group, groupIndex) in $sortedGroups">
                <div class="property-group"
                    :key="groupIndex"
                    v-if="$groupedProperties[groupIndex].length">
                    <cmdb-collapse
                        :label="group['bk_group_name']"
                        :collapse.sync="groupState[group['bk_group_id']]">
                        <ul class="property-list clearfix">
                            <li class="property-item clearfix fl"
                                v-for="(property, propertyIndex) in $groupedProperties[groupIndex]"
                                :key="propertyIndex">
                                <span class="property-name fl" :title="property['bk_property_name']">{{property['bk_property_name']}}</span>
                                <slot :name="property['bk_property_id']">
                                    <span class="property-value clearfix fl"
                                        v-if="property.unit"
                                        :title="`${getValue(property)} ${property.unit}`">
                                        <span class="property-value-text fl">{{getValue(property)}}</span>
                                        <span class="property-value-unit fl" v-if="getValue(property) !== '--'">{{property.unit}}</span>
                                    </span>
                                    <span class="property-value fl" v-else-if="property.bk_property_type === 'objuser'" v-user.title="getValue(property)"></span>
                                    <span class="property-value fl" v-else :title="getValue(property)">{{getValue(property)}}</span>
                                </slot>
                            </li>
                        </ul>
                    </cmdb-collapse>
                </div>
            </template>
        </div>
        <div class="details-options"
            v-if="showOptions"
            :class="{ sticky: scrollbar }">
            <slot name="details-options">
                <cmdb-auth v-if="showEdit" class="inline-block-middle" :auth="authResources(editAuth)">
                    <bk-button slot-scope="{ disabled }"
                        class="button-edit"
                        theme="primary"
                        :disabled="disabled"
                        @click="handleEdit">
                        {{editText}}
                    </bk-button>
                </cmdb-auth>
                <cmdb-auth v-if="showDelete" class="inline-block-middle" :auth="authResources(deleteAuth)">
                    <bk-button slot-scope="{ disabled }"
                        class="button-delete"
                        theme="danger"
                        :disabled="disabled"
                        @click="handleDelete">
                        {{deleteText}}
                    </bk-button>
                </cmdb-auth>
            </slot>
        </div>
    </div>
</template>

<script>
    import formMixins from '@/mixins/form'
    import RESIZE_EVENTS from '@/utils/resize-events'
    export default {
        name: 'cmdb-details',
        mixins: [formMixins],
        props: {
            inst: {
                type: Object,
                required: true
            },
            showOptions: {
                type: Boolean,
                default: true
            },
            editButtonText: {
                type: String,
                default: ''
            },
            deleteButtonText: {
                type: String,
                default: ''
            },
            showEdit: {
                type: Boolean,
                default: true
            },
            showDelete: {
                type: Boolean,
                default: true
            },
            editAuth: {
                type: [String, Array],
                default: ''
            },
            deleteAuth: {
                type: [String, Array],
                default: ''
            }
        },
        data () {
            return {
                scrollbar: false
            }
        },
        computed: {
            editText () {
                return this.editButtonText || this.$t('编辑')
            },
            deleteText () {
                return this.deleteButtonText || this.$t('删除')
            }
        },
        mounted () {
            RESIZE_EVENTS.addResizeListener(this.$refs.detailsWrapper, this.checkScrollbar)
        },
        beforeDestroy () {
            RESIZE_EVENTS.removeResizeListener(this.$el.detailsWrapper, this.checkScrollbar)
        },
        methods: {
            authResources (auth) {
                if (!auth) return {}
                if (Array.isArray(auth) && !auth.length) return {}
                return this.$authResources({ type: auth })
            },
            checkScrollbar () {
                const $layout = this.$el
                this.scrollbar = $layout.scrollHeight !== $layout.offsetHeight
            },
            handleToggleGroup (group) {
                const groupId = group['bk_group_id']
                const collapse = !!this.groupState[groupId]
                this.$set(this.groupState, groupId, !collapse)
            },
            getTitle (inst, property) {
                return `${property['bk_property_name']}: ${inst[property['bk_property_id']] || '--'} ${property.unit}`
            },
            getValue (property) {
                const value = this.inst[property['bk_property_id']]
                return String(value).length ? value : '--'
            },
            handleEdit () {
                this.$emit('on-edit', this.inst)
            },
            handleDelete () {
                this.$emit('on-delete', this.inst)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .details-layout{
        height: 100%;
        padding: 0 0 0 32px;
        @include scrollbar-y;
    }
    .property-group{
        padding: 7px 0 10px 0;
        &:first-child{
            padding: 28px 0 10px 0;
        }
    }
    .group-name{
        font-size: 16px;
        line-height: 16px;
        color: #333948;
        overflow: visible;
        .group-toggle {
            cursor: pointer;
            &.collapse .bk-icon {
                transform: rotate(-90deg);
            }
            .bk-icon {
                vertical-align: baseline;
                font-size: 12px;
                font-weight: bold;
                transition: transform .2s ease-in-out;
            }
        }
    }
    .property-list{
        padding: 4px 0;
        .property-item{
            width: 50%;
            max-width: 400px;
            margin: 12px 0 0;
            font-size: 14px;
            line-height: 26px;
            .property-name{
                position: relative;
                width: 35%;
                padding: 0 16px 0 0;
                color: #63656e;
                @include ellipsis;
                &:after{
                    content: ":";
                    position: absolute;
                    right: 10px;
                }
            }
            .property-value{
                width: 65%;
                padding: 0 15px 0 0;
                color: #313238;
                @include ellipsis;
                &-text{
                    display: block;
                    max-width: calc(100% - 60px);
                    @include ellipsis;
                }
                &-unit{
                    display: block;
                    width: 60px;
                    padding: 0 0 0 5px;
                    @include ellipsis;
                }
            }
        }
    }
    .details-options{
        position: sticky;
        bottom: 0;
        left: 0;
        width: 100%;
        padding: 28px 18px 0;
        &.sticky {
            width: calc(100% + 32px);
            margin: 0 0 0 -40px;
            padding: 10px 50px;
            border-top: 1px solid $cmdbBorderColor;
            background-color: #fff;
        }
        .button-edit{
            min-width: 76px;
            margin-right: 4px;
        }
        .button-delete{
            min-width: 76px;
            background-color: #fff;
            color: #ff5656;
            &:disabled {
                color: #dcdee5 !important;
            }
        }
    }
</style>
