<template>
    <div class="topo-attribute-wrapper" id="wrapper" :class="{'full-screen': isFullScreen}" v-show="isShow">
        <div class="loading" v-bkloading="{isLoading: attr.isLoading}">
            <i class="bk-icon icon-close" @click="closePop"></i>
            <div class="attr-title">{{objName}} {{instName}}</div>
            <div class="attribute-padding">
                <div class="attribute-box" id="box">
                    <v-attribute 
                        ref="attribute"
                        :formFields="attr.formFields"
                        :formValues="attr.formValues"
                        :type="attr.type"
                        :showBtnGroup="false"
                        :active="isShow"
                        :objId="objId">
                    </v-attribute>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    import vAttribute from '@/components/object/attribute'
    import { mapGetters } from 'vuex'
    export default {
        props: {
            isShow: {
                type: Boolean,
                default: false
            },
            objId: {
                type: String,
                default: ''
            },
            instId: {
                default: ''
            },
            objName: {
                type: String,
                default: ''
            },
            instName: {
                type: String,
                default: ''
            },
            isFullScreen: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                attr: {
                    formFields: [],
                    formValues: {},
                    type: 'update',
                    isLoading: false
                }
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount'
            ]),
            ...mapGetters('object', [
                'attribute'
            ]),
            formValuesConfig () {
                let config = {
                    url: '',
                    type: 'post',
                    params: {
                        page: {},
                        fields: {},
                        condition: {}
                    }
                }
                if (this.objId === 'biz') {
                    config.url = `biz/search/${this.bkSupplierAccount}`
                    config.params.fields = []
                    config.params.condition = {
                        bk_biz_id: this.instId
                    }
                } else if (this.objId === 'host') {
                    config.type = 'get'
                    config.url = `/hosts/${this.bkSupplierAccount}/${this.instId}`
                } else {
                    config.url = `inst/association/search/owner/${this.bkSupplierAccount}/object/${this.objId}`
                    config.params.condition[this.objId] = [{
                        field: 'bk_inst_id',
                        operator: '$eq',
                        value: this.instId
                    }]
                }
                return config
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    this.initData()
                }
            }
        },
        methods: {
            closePop () {
                this.$emit('update:isShow', false)
            },
            resetAttributeBox () {
                let box = document.getElementById('box')
                let topo = document.getElementById('topo')
                box.style.maxHeight = `${topo.offsetHeight * 0.8}px`
            },
            async initData () {
                this.attr.isLoading = true
                await Promise.all([
                    this.$store.dispatch('object/getAttribute', {objId: this.objId}),
                    this.getFormValues()
                ])
                this.attr.formFields = this.attribute[this.objId]
                this.attr.isLoading = false
                this.$nextTick(() => {
                    setTimeout(() => {
                        this.resetAttributeBox()
                    }, 100)
                })
            },
            async getFormValues () {
                try {
                    let res = await this.$axios({
                        url: this.formValuesConfig.url,
                        data: this.formValuesConfig.params,
                        method: this.formValuesConfig.type
                    })
                    if (this.objId === 'host') {
                        let values = {}
                        res.data.map(({bk_property_id: bkPropertyId, bk_property_value: bkPropertyValue}) => {
                            values[bkPropertyId] = bkPropertyValue !== null ? bkPropertyValue : ''
                        })
                        this.attr.formValues = values
                    } else {
                        this.attr.formValues = res.data.info[0]
                    }
                } catch (e) {
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }
            }
        },
        created () {
            window.onresize = () => {
                this.resetAttributeBox()
            }
        },
        components: {
            vAttribute
        }
    }
</script>

<style lang="scss" scoped>
    .topo-attribute-wrapper {
        position: absolute;
        background: #fff;
        width: 710px;
        min-height: 200px;
        top: 30px;
        left: 50%;
        transform: translateX(-50%);
        box-shadow: 0px 2px 9.6px 0.4px rgba(0, 0, 0, 0.4);
        &.full-screen {
            top: 50%;
            transform: translate(-50%, -50%);
        }
        .loading {
            min-height: 200px;
            height: 100%;
        }
        .attr-title {
            font-size: 14px;
            padding-left: 16px;
            height: 49px;
            line-height: 48px;
            border-bottom: 1px solid #e5e5e5;
            background: #f7f7f7;
            color: #333948;
        }
        .icon-close {
            padding: 6px;
            font-size: 12px;
            position: absolute;
            right: 6px;
            top: 12px;
            cursor: pointer;
            color: #333948;
            border-radius: 50%;
            &:hover {
                background: #e5e5e5;
            }
        }
        .attribute-padding {
            height: 100%;
            padding-right: 2px;
        }
        .attribute-box {
            max-height: 300px;
            overflow: auto;
            @include scrollbar;
            padding-bottom: 40px;
            margin: 20px 0;
        }
    }
</style>

