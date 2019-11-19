<template>
    <div class="topology-details-layout" :class="{ 'full-screen': fullScreen }">
        <div class="details-container" ref="detailsContainer" v-bkloading="{ isLoading: $loading() }" v-click-outside="handleHideDetails">
            <div class="details-title" ref="detailsTitle">
                {{title}}
                <i class="bk-icon icon-close" @click="handleHideDetails"></i>
            </div>
            <cmdb-details ref="detailsPopup" class="details-popup"
                :show-options="false"
                :inst="inst"
                :properties="properties"
                :property-groups="propertyGroups">
            </cmdb-details>
        </div>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        props: {
            fullScreen: {
                type: Boolean,
                default: false
            },
            objId: {
                type: String,
                required: true
            },
            instId: {
                type: Number,
                required: true
            },
            title: {
                type: String,
                default: ''
            },
            show: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                properties: [],
                propertyGroups: [],
                inst: {}
            }
        },
        mounted () {
            this.init()
        },
        methods: {
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('objectModelFieldGroup', ['searchGroup']),
            ...mapActions('objectCommonInst', ['searchInst']),
            ...mapActions('objectBiz', ['searchBusiness']),
            ...mapActions('hostSearch', ['getHostBaseInfo']),
            async init () {
                const [properties, propertyGroups] = await Promise.all([
                    this.getObjectProperties(),
                    this.searchGroup({
                        objId: this.objId,
                        params: this.$injectMetadata(),
                        config: {
                            requestId: `get_${this.objId}_property_groups`
                        }
                    })
                ])
                this.properties = properties
                this.propertyGroups = propertyGroups
                const inst = await this.getDetails()
                this.inst = this.$tools.flattenItem(properties, inst)
                this.$nextTick(() => {
                    const detailsContainerHeight = this.$refs.detailsContainer.getBoundingClientRect().height
                    const detailsTitleHeight = this.$refs.detailsTitle.getBoundingClientRect().height
                    this.$refs.detailsPopup.$el.style.height = detailsContainerHeight - detailsTitleHeight + 'px'
                })
            },
            getObjectProperties () {
                return this.searchObjectAttribute({
                    params: this.$injectMetadata({
                        'bk_supplier_account': this.supplierAccount,
                        'bk_obj_id': this.objId
                    }),
                    config: {
                        requestId: `post_searchObjectAttribute_${this.objId}`
                    }
                })
            },
            getDetails () {
                let promise
                if (this.objId === 'host') {
                    promise = this.getHostDetails()
                } else if (this.objId === 'biz') {
                    promise = this.getBusinessDetails()
                } else {
                    promise = this.getInstDetails()
                }
                return promise
            },
            getHostDetails () {
                return this.getHostBaseInfo({ hostId: this.instId }).then(data => {
                    const inst = {}
                    data.forEach(field => {
                        inst[field['bk_property_id']] = field['bk_property_value']
                    })
                    return inst
                })
            },
            getBusinessDetails () {
                return this.searchBusiness({
                    params: {
                        condition: { 'bk_biz_id': this.instId },
                        fields: [],
                        page: { start: 0, limit: 1 }
                    }
                }).then(({ info }) => info[0])
            },
            getInstDetails () {
                return this.searchInst({
                    objId: this.objId,
                    params: this.$injectMetadata({
                        condition: {
                            [this.objId]: [{
                                field: 'bk_inst_id',
                                operator: '$eq',
                                value: this.instId
                            }]
                        },
                        fields: {},
                        page: { start: 0, limit: 1 }
                    })
                }).then(({ info }) => info[0])
            },
            handleHideDetails () {
                this.$emit('update:show', false)
                this.inst = {}
                this.properties = []
                this.propertyGroups = []
            }
        }
    }
</script>

<style lang="scss" scoped>
    .topology-details-layout {
        position: fixed;
        left: 0;
        right: 0;
        top: 0;
        bottom: 0;
        text-align: right;
        &.full-screen {
            text-align: center;
        }
        &:before {
            content: "";
            display: inline-block;
            vertical-align: middle;
            width: 0;
            height: 100%;
        }
        .details-container {
            position: relative;
            display: inline-block;
            width: 710px;
            min-height: 250px;
            max-height: 80%;
            margin: 0 45px 0 0;
            vertical-align: middle;
            text-align: left;
            background-color: #fff;
            box-shadow: 0px 2px 9.6px 0.4px rgba(0, 0, 0, 0.4);
            z-index: 100;
            .details-title {
                position: relative;
                height: 49px;
                padding: 0 0 0 16px;
                border-bottom: 1px solid $cmdbBorderColor;
                line-height: 48px;
                color: #333948;
                background-color: #f7f7f7;
                .icon-close {
                    position: absolute;
                    right: 6px;
                    top: 12px;
                    padding: 6px;
                    font-size: 12px;
                    cursor: pointer;
                    color: #333948;
                    border-radius: 50%;
                    &:hover {
                        background-color: #e5e5e5;
                    }
                }
            }
            .details-popup {
                padding-bottom: 20px;
                @include scrollbar-y;
            }
        }
    }
</style>
