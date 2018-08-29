<template>
    <div class="details-layout">
        <template v-for="(group, groupIndex) in $sortedGroups">
            <div class="property-group"
                :key="groupIndex"
                v-if="$groupedProperties[groupIndex].length"
                v-show="group['bk_group_id'] !== 'none' || showNoneGroup">
                <h3 class="group-name">{{group['bk_group_name']}}</h3>
                <ul class="property-list">
                    <li class="property-item clearfix"
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
            </div>
            <div class="none-group" v-if="group['bk_group_id'] === 'none'">
                <a href="javascript:void(0)" class="none-group-link"
                    :class="{'open': showNoneGroup}"
                    @click="showNoneGroup = !showNoneGroup">
                    {{$t("Common['更多属性']")}}
                </a>
            </div>
        </template>
        <div class="details-options" v-if="showOptions">
            <bk-button class="button-edit" type="primary"
                v-if="showEdit"
                @click="handleEdit">
                {{editText}}
            </bk-button>
        </div>
    </div>
</template>

<script>
    import formMixins from '@/mixins/form'
    export default {
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
            }
        },
        data () {
            return {
                showNoneGroup: false
            }
        },
        computed: {
            editText () {
                return this.editButtonText || this.$t("Common['编辑']")
            },
            deleteText () {
                return this.deleteButtonText || this.$t("Common['删除']")
            }
        },
        methods: {
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
        position: relative;
        height: 100%;
        padding: 0 0 0 32px;
        @include scrollbar-y;
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
    }
    .property-list{
        padding: 4px 0;
        .property-item{
            margin: 12px 0 0;
            font-size: 14px;
            height: 36px;
            line-height: 36px;
            .property-name{
                position: relative;
                width: 120px;
                padding: 0 26px 0 0;
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
                width: calc(100% - 120px);
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
    .none-group{
        text-align: center;
        margin: 26px 0 0 0;
        .none-group-link{
            color: #6b7baa;
            font-size: 12px;
            &.open:after{
                transform: rotate(0deg);
            }
            &.open:hover:after{
                transform: rotate(180deg);
            }
            &:hover{
                color: #498fe0;
            }
            &:hover:after{
                background-image: url('../../../assets/images/icon/icon-result-slide-hover.png');
                transform: rotate(0deg);
            }
            &:after{
                content: '';
                display: inline-block;
                width: 11px;
                height: 10px;
                margin-left: 12px;
                background: url('../../../assets/images/icon/icon-result-slide.png') no-repeat;
                transform: rotate(180deg);
            }
        }
    }
    .details-options{
        position: sticky;
        bottom: 0;
        left: 0;
        background-color: #fff;
        padding: 20px 0 0 120px;
        .button-edit{
            margin-right: 4px;
        }
    }
</style>