<template>
    <div class="cross-import-wrapper">
        <div class="search-container">
            <bk-select class="search-area search-field" :selected.sync="search.selectedPlat">
                <bk-select-option v-for="(plat, index) in search.plat" :key="index"
                    :label="plat['bk_cloud_name']"
                    :value="plat['bk_cloud_id']">
                </bk-select-option>
            </bk-select>
            <span class="search-ip search-field">
                <input class="bk-form-input" type="text"
                    v-model.trim="search.ip"
                    :placeholder="$t('BusinessTopology[\'请输入IP地址\']')"
                    @keydown.enter="doSearch">
                <i class="bk-icon icon-close-circle" v-show="search.ip.length" @click="reset"></i>
            </span>
            <bk-button class="search-btn search-field" type="primary" :disabled="!isValidIp" @click="doSearch">{{$t("Common['查询']")}}</bk-button>
        </div>
        <div class="search-result" v-if="Object.keys(result).length">
            <div class="attribute-group" v-if="hostRelation.length">
                <h3 class="title">{{$t("Hosts['主机拓扑']")}}</h3>
                <ul class="attribute-list">
                    <li class="attribute-item" v-for="(relation, index) in hostRelation" :key="index">{{relation}}</li>
                </ul>
            </div>
            <div class="attribute-group" v-for="(groupId, groupIndex) in hostPropertyGroupOrder"
                v-show="groupId !== 'none' || !isNoneGroupHide"
                :key="groupIndex">
                <h3 class="title">{{getGroupName(groupId)}}</h3>
                <ul class="clearfix attribute-list">
                    <template v-for="(property, propertyIndex) in groupedHostProperty[groupId]">
                        <li class="attribute-item fl" v-if="!property['bk_isapi']" :key="propertyIndex">
                            <template v-if="property['bk_property_type'] !== 'bool'">
                                <span class="attribute-item-label">{{property['bk_property_name']}} :</span>
                                <span class="attribute-item-value" :title="getFieldValue(property)">{{getFieldValue(property)}}</span>
                            </template>
                            <template v-else>
                                <span class="attribute-item-label">{{property['bk_property_name']}}</span>
                                <span class="attribute-item-value bk-form-checkbox">
                                    <input type="checkbox" :checked="getFieldValue(property)" disabled>
                                </span>
                            </template>
                        </li>
                    </template>
                </ul>
            </div>
             <div class="attribute-group-more" v-if="groupedHostProperty['none'].length">
                <a href="javascript:void(0)" class="group-more-link" :class="{'open': !isNoneGroupHide}" @click="isNoneGroupHide = !isNoneGroupHide">{{$t("Common['更多属性']")}}</a>
            </div>
        </div>
        <div class="search-tips" v-show="showTips" v-bkloading="{isLoading: loading}">
            <p v-if="noResult">{{$t("BusinessTopology['未查询到该IP地址对应的主机']")}}</p>
            <p v-else>{{$t("BusinessTopology['请输入完整IP地址进行查询']")}}</p>
        </div>
        <div class="search-footer">
            <bk-button type="primary" @click="doCrossImport" :disabled="!Object.keys(result).length">{{$t("Common['确定']")}}</bk-button>
            <button class="bk-button vice-btn" @click="cancelCrossImport">{{$t("Common['取消']")}}</button>
        </div>
    </div>
</template>
<script>
    import { getHostRelation } from '@/utils/util'
    import {mapGetters} from 'vuex'
    export default {
        props: {
            isShow: {
                type: Boolean,
                default: false
            },
            bizId: {
                type: Number,
                default: -1
            },
            moduleId: {
                type: Number,
                default: -1
            }
        },
        data () {
            return {
                search: {
                    plat: [],
                    selectedPlat: -1,
                    ip: ''
                },
                result: {},
                resultPlat: 0,
                noResult: false,
                isNoneGroupHide: true,
                attribute: {},
                propertyGroups: {},
                loading: false,
                showTips: true,
                hostRelation: []
            }
        },
        computed: {
            ...mapGetters(['bkSupplierAccount']),
            isValidIp () {
                let reg = /^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/
                return reg.test(this.search.ip)
            },
            hostProperty () {
                return this.attribute['host'] || []
            },
            hostPropertyGroup () {
                let hostPropertyGroup = this.propertyGroups['host'] || []
                hostPropertyGroup.sort((groupA, groupB) => {
                    return groupA['bk_group_index'] - groupB['bk_group_index']
                })
                return hostPropertyGroup
            },
            hostPropertyGroupOrder () {
                let order = this.hostPropertyGroup.map(({bk_group_id: bkGroupId}) => bkGroupId)
                order.push('none')
                return order
            },
            groupedHostProperty () {
                let groupedHostProperty = {}
                this.hostPropertyGroupOrder.forEach(groupId => {
                    groupedHostProperty[groupId] = []
                })
                this.hostProperty.forEach(property => {
                    if (groupedHostProperty.hasOwnProperty(property['bk_property_group'])) {
                        groupedHostProperty[property['bk_property_group']].push(property)
                    }
                })
                return groupedHostProperty
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    this.reset()
                    this.getHostAttribute()
                    this.getPropertyGroups()
                    this.getPlat()
                }
            },
            'search.ip' () {
                this.noResult = false
            }
        },
        methods: {
            getHostAttribute () {
                let hostObjId = 'host'
                if (!this.attribute.hasOwnProperty(hostObjId)) {
                    this.$axios.post('object/attr/search', {
                        bk_obj_id: hostObjId,
                        bk_supplier_account: this.bkSupplierAccount
                    }).then(res => {
                        if (res.result) {
                            res.data.sort((objA, objB) => {
                                return objA['bk_property_index'] - objB['bk_property_index']
                            })
                            this.attribute[hostObjId] = res.data
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                }
            },
            getPropertyGroups () {
                let hostObjId = 'host'
                if (!this.propertyGroups.hasOwnProperty(hostObjId)) {
                    this.$axios.post(`objectatt/group/property/owner/${this.bkSupplierAccount}/object/${hostObjId}`, {}).then(res => {
                        if (res.result) {
                            this.propertyGroups[hostObjId] = res.data
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                }
            },
            getPlat () {
                if (!this.search.plat.length) {
                    this.$axios.post(`inst/search/owner/${this.bkSupplierAccount}/object/plat`, {
                        condition: {},
                        fields: [],
                        page: {}
                    }).then(res => {
                        if (res.result) {
                            this.search.plat = res.data.info
                            this.search.selectedPlat = 0
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    })
                } else {
                    this.search.selectedPlat = 0
                }
            },
            getGroupName (groupId) {
                if (groupId === 'none') {
                    return this.$t("Common['更多属性']")
                }
                return this.hostPropertyGroup.find(({bk_group_id: bkGroupId}) => bkGroupId === groupId)['bk_group_name']
            },
            getFieldValue (property) {
                let {
                    bk_property_id: bkPropertyId,
                    bk_property_type: bkPropertyType
                } = property
                let value = this.result[bkPropertyId]
                if (property['bk_asst_obj_id']) {
                    let associateName = []
                    if (Array.isArray(value)) {
                        value.map(({bk_inst_name: bkInstName}) => {
                            if (bkInstName) {
                                associateName.push(bkInstName)
                            }
                        })
                    }
                    return associateName.join(',')
                } else if (bkPropertyType === 'date') {
                    return this.$formatTime(value, 'YYYY-MM-DD')
                } else if (bkPropertyType === 'time') {
                    return this.$formatTime(value)
                } else if (bkPropertyType === 'enum') {
                    const option = (Array.isArray(property.option) ? property.option : []).find(({id}) => id === value)
                    return option ? option.name : ''
                } else {
                    return value
                }
            },
            reset () {
                this.search.ip = ''
                this.result = {}
                this.resultPlat = 0
                this.noResult = false
                this.showTips = true
            },
            doSearch () {
                if (this.isValidIp) {
                    this.loading = true
                    this.noResult = false
                    this.resultPlat = this.search.selectedPlat
                    this.$axios.post('hosts/search', {
                        bk_biz_id: -1,
                        condition: [{
                            bk_obj_id: 'host',
                            fields: [],
                            condition: [{
                                field: 'bk_cloud_id',
                                operator: '$eq',
                                value: this.search.selectedPlat
                            }]
                        }],
                        ip: {
                            data: [this.search.ip],
                            exact: 1,
                            flag: 'bk_host_innerip'
                        },
                        page: {
                            start: 0,
                            limit: 1
                        }
                    }).then(res => {
                        if (res.result) {
                            if (res.data.count) {
                                this.showTips = false
                                this.noResult = false
                                this.result = res.data.info[0]['host']
                                this.hostRelation = getHostRelation(res.data.info[0])
                            } else {
                                this.noResult = true
                                this.showTips = true
                                this.result = {}
                                this.hostRelation = []
                            }
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                        this.loading = false
                    }).catch(() => {
                        this.loading = false
                    })
                }
            },
            doCrossImport () {
                this.$axios.post('hosts/modules/biz/mutilple', {
                    bk_biz_id: this.bizId,
                    bk_module_id: this.moduleId,
                    host_info: [{
                        bk_host_innerip: this.result['bk_host_innerip'],
                        bk_cloud_id: this.resultPlat
                    }]
                }).then(res => {
                    if (res.result) {
                        this.$alertMsg(this.$t("Common['导入成功']"), 'success')
                        this.$emit('update:isShow', false)
                        this.$emit('handleCrossImportSuccess')
                    } else {
                        if (res.data && Array.isArray(res.data.error) && res.data.error.length) {
                            this.$alertMsg(res.data.error[0])
                        } else {
                            this.$alertMsg(res['bk_error_msg'])
                        }
                    }
                })
            },
            cancelCrossImport () {
                this.$emit('update:isShow', false)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .search-container{
        font-size: 0;
        padding: 20px 20px 0;
        .search-field{
            display: inline-block;
            vertical-align: top;
            font-size: 14px;
        }
        .search-area{
            width: 150px;
        }
        .search-ip{
            position: relative;
            width: 390px;
            margin: 0 10px;
            .bk-form-input{
                padding-right: 36px;
            }
            .icon-close-circle{
                position: absolute;
                font-size: 16px;
                top: 10px;
                right: 10px;
                cursor: pointer;
            }
        }
        .search-btn{
            width: 100px;
        }
    }
    .search-result{
        max-height: 400px;
        overflow-x: hidden;
        overflow-y: auto;
        margin-top: 20px;
        padding: 0 20px 20px;
        @include scrollbar;
        .title{
            margin: 0;
            font-size: 14px;
            line-height: 14px;
            overflow: visible;
            color: #333948;
        }
        .attribute-group{
            padding: 17px 0 0 0;
            &:first-child{
                padding: 0;
            }
        }
    }
    .attribute-list{
        padding: 4px 0;
        .attribute-item{
            width: 50%;
            font-size: 12px;
            line-height: 16px;
            margin: 12px 0 0 0;
            white-space: nowrap;
            .attribute-item-label{
                width: 100px;
                color: #6b7baa;
                text-align: right;
                display: inline-block;
                overflow: hidden;
                text-overflow: ellipsis;
                margin-right: 10px;
            }
            .attribute-item-value{
                max-width: 140px;
                display: inline-block;
                overflow: hidden;
                text-overflow: ellipsis;
                color: #4d597d;
            }
            .attribute-item-value.bk-form-checkbox{
                padding: 0;
                font-size: 0;
                transform: scale(0.889);
                vertical-align: -1px;
                vertical-align: top;
                input[type="checkbox"]{
                    &:checked{
                        background-position: -33px -62px;
                    }
                }
            }
        }
    }
    .attribute-group-more{
        text-align: center;
        margin-top: 26px;
        .group-more-link{
            color: #6b7baa;
            text-decoration: none;
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
                background-image: url('../../../common/images/icon/icon-result-slide-hover.png');
                transform: rotate(0deg);
            }
            &:after{
                content: '';
                display: inline-block;
                width: 11px;
                height: 10px;
                margin-left: 12px;
                background: url('../../../common/images/icon/icon-result-slide.png') no-repeat;
                transform: rotate(180deg);
            }
        }
    }
    .search-tips{
        padding: 20px;
        text-align: center;
    }
    .search-footer{
        text-align: right;
        height: 60px;
        line-height: 60px;
        border-top: 1px solid #e5e5e5;
        background-color: #fafafa;
        padding: 0 20px;
        font-size: 0;
        .bk-button{
            margin: 0 0 0 20px;
        }
    }
</style>