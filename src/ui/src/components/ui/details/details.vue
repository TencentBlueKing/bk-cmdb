<template>
    <div class="details-layout">
        <template v-for="(group, groupIndex) in $sortedGroups">
            <div class="property-group"
                :key="groupIndex"
                v-if="$groupedProperties[groupIndex].length">
                <h3 class="group-name">
                    <span class="group-toggle"
                        @click="handleToggleGroup(group)"
                        :class="{collapse: collapseStatus[group['bk_group_id']]}">
                        <i class="bk-icon icon-angle-down"></i>
                        {{group['bk_group_name']}}
                    </span>
                </h3>
                <bk-collapse-transition>
                    <ul class="property-list clearfix"
                        v-show="!collapseStatus[group['bk_group_id']]">
                        <li class="property-item clearfix fl"
                            v-for="(property, propertyIndex) in $groupedProperties[groupIndex]"
                            :key="propertyIndex"
                            :title="getTitle(inst, property)">
                            <span class="property-name fl">{{property['bk_property_name']}}</span>
                            <span class="property-value clearfix fl" v-if="property.unit">
                                <span class="property-value-text fl">{{inst[property['bk_property_id']] || '--'}}</span>
                                <span class="property-value-unit fl">{{property.unit}}</span>
                            </span>
                            <span class="property-value fl" v-else>{{inst[property['bk_property_id']] || '--'}}</span>
                        </li>
                    </ul>
                </bk-collapse-transition>
            </div>
        </template>
        <slot name="details-options" >
            <div class="details-options" v-if="showOptions" ref="test">
                <bk-button class="button-edit" type="primary"
                    v-if="showEdit"
                    :disabled="!$authorized.update"
                    @click="handleEdit">
                    {{editText}}
                </bk-button>
                <bk-button class="button-delete" type="danger"
                    v-if="showDelete"
                    :disabled="!$authorized.delete"
                    @click="handleDelete">
                    {{deleteText}}
                </bk-button>
            </div>
        </slot>
    </div>
</template>

<script>
    import formMixins from '@/mixins/form'
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
            }
        },
        data () {
            return {
                collapseStatus: {
                    none: true
                }
            }
        },
        computed: {
            editText () {
                return this.editButtonText || this.$t("Common['属性编辑']")
            },
            deleteText () {
                return this.deleteButtonText || this.$t("Common['删除']")
            }
        },
        mounted () {
            console.log(this.$refs)
        },
        methods: {
            handleToggleGroup (group) {
                const groupId = group['bk_group_id']
                const collapse = !!this.collapseStatus[groupId]
                this.$set(this.collapseStatus, groupId, !collapse)
            },
            getTitle (inst, property) {
                return `${property['bk_property_name']}: ${inst[property['bk_property_id']] || '--'} ${property.unit}`
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
        overflow: auto;
        @include scrollbar;
    }
    .property-group{
        padding: 17px 0 0 0;
        &:first-child{
            padding: 28px 0 0 0;
        }
    }
    .group-name{
        font-size: 14px;
        line-height: 14px;
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
            margin: 12px 0 0;
            font-size: 12px;
            line-height: 16px;
            .property-name{
                position: relative;
                width: 35%;
                padding: 0 16px 0 0;
                text-align: right;
                color: $cmdbTextColor;
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
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;
        height: 62px;
        padding: 14px 20px;
        background-color: #f9f9f9;
        .button-edit{
            width: 110px;
            margin-right: 4px;
        }
        .button-delete{
            width: 110px;
            background-color: #fff;
            color: #ff5656;
        }
    }
</style>