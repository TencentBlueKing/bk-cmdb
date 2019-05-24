<template>
    <div class="property">
        <div class="group"
            v-for="(group, index) in groupedProperties"
            :key="index">
            <h2 class="group-name">{{group.bk_group_name}}</h2>
            <ul class="property-list clearfix">
                <li class="property-item fl"
                    v-for="property in group.properties"
                    :key="property.id">
                    <span class="property-name"
                        :title="property.bk_property_name">
                        {{property.bk_property_name}}
                    </span>
                    <span class="property-value"
                        v-show="property !== editState.property"
                        :title="$tools.getPropertyText(property, host)">
                        {{$tools.getPropertyText(property, host)}}
                    </span>
                    <template v-if="updateAuth && isPropertyEditable(property)">
                        <i class="property-edit icon-cc-edit"
                            v-show="property !== editState.property"
                            @click="setEditState(property)">
                        </i>
                        <div class="property-form" v-if="property === editState.property">
                            <component class="form-component"
                                :is="`cmdb-form-${property.bk_property_type}`"
                                :class="{ error: errors.has(property.bk_property_id) }"
                                :options="property.option || []"
                                :data-vv-name="property.bk_property_id"
                                :data-vv-as="property.bk_property_name"
                                :exclude="false"
                                v-validate="$tools.getValidateRules(property)"
                                v-model.trim="editState.value">
                            </component>
                            <i class="form-confirm bk-icon icon-check-1" @click="confirm"></i>
                            <i class="form-cancel bk-icon icon-close" @click="exitForm"></i>
                        </div>
                    </template>
                </li>
            </ul>
        </div>
    </div>
</template>

<script>
    import { mapGetters, mapState } from 'vuex'
    import { OPERATION, RESOURCE_HOST } from '../router.config.js'
    export default {
        name: 'cmdb-host-property',
        data () {
            return {
                OPERATION,
                editState: {
                    property: null,
                    value: null
                }
            }
        },
        computed: {
            ...mapState('hostDetails', ['info']),
            ...mapGetters('hostDetails', ['groupedProperties']),
            host () {
                return this.info.host || {}
            },
            updateAuth () {
                const isResourceHost = this.$route.name === RESOURCE_HOST
                if (isResourceHost) {
                    return this.$isAuthorized(OPERATION.U_RESOURCE_HOST)
                }
                return this.$isAuthorized(OPERATION.U_HOST)
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
            confirm () {
                const { property, value } = this.editState
                this.$store.dispatch('hostUpdate/updateHost', this.$injectMetadata({
                    [property.bk_property_id]: value,
                    bk_host_id: this.host.bk_host_id
                }))
                this.$store.commit('hostDetails/updateInfo', {
                    [property.bk_property_id]: value
                })
                this.exitForm()
            },
            exitForm () {
                this.editState.property = null
                this.editState.value = null
            }
        }
    }
</script>

<style lang="scss" scoped>
    .property {
        overflow: hidden;
    }
    .group {
        margin: 22px 0 0 0;
        .group-name {
            line-height: 19px;
            font-size: 14px;
            font-weight: normal;
            color: #313238;
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
        &:hover {
            .property-item .property-edit {
                display: inline-block;
            }
        }
    }
    .property-list {
        width: 740px;
        margin: 15px 0 0 0;
        line-height: 26px;
        color: #63656e;
        .property-item {
            width: 50%;
            font-size: 0;
            .property-name {
                position: relative;
                display: inline-block;
                width: 105px;
                padding: 0 16px 0 0;
                vertical-align: middle;
                text-align: right;
                font-size: 12px;
                @include ellipsis;
                &:after {
                    position: absolute;
                    right: 2px;
                    content: "：";
                }
            }
            .property-value {
                font-size: 12px;
                display: inline-block;
                max-width: 240px;
                vertical-align: middle;
                @include ellipsis;
            }
            .property-edit {
                display: none;
                margin: 0 0 0 8px;
                vertical-align: middle;
                font-size: 16px;
                color: #3c96ff;
                cursor: pointer;
            }
            .property-form {
                display: inline-block;
                vertical-align: middle;
            }
        }
    }
    .property-form {
        font-size: 0;
        .bk-icon {
            display: inline-block;
            vertical-align: middle;
            width: 26px;
            height: 26px;
            margin: 0 0 0 6px;
            border-radius: 2px;
            border: 1px solid #c4c6cc;
            line-height: 24px;
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
                &:before {
                    display: inline-block;
                    transform: scale(0.66);
                }
            }
            &:hover {
                font-weight: bold;
            }
        }
        .form-component {
            display: inline-block;
            vertical-align: middle;
            width: 195px;
            height: 26px;
            margin: 0 4px 0 0;
            /deep/ {
                .bk-date-picker,
                .bk-selector-input,
                .form-float-input,
                .form-singlechar-input,
                .form-longchar-input,
                .form-int-input,
                [name="date-select"] {
                    height: 26px ;
                    font-size: 12px !important;
                }
                .bk-date-picker:after {
                    width: 26px;
                    height: 26px;
                }
                .date-dropdown-panel,
                .bk-selector-list {
                    margin-top: -10px;
                }
                .bk-selector-icon {
                    top: 7px;
                }
                .bk-selector-node .text {
                    line-height: 26px;
                    font-size: 12px;
                }
                .objuser-layout {
                    font-size: 12px;
                    .objuser-container {
                        min-height: 26px;
                    }
                    .objuser-container.placeholder:after {
                        line-height: 24px;
                    }
                    .objuser-selected {
                        height: 18px;
                        margin: 1px 3px;
                        line-height: 16px;
                    }
                    .objuser-input {
                        height: 18px;
                        line-height: 18px;
                        margin: 1px 0 0;
                    }
                }
            }
        }
    }
</style>
