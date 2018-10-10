<template>
    <div class="base-info-wrapper">
        <div class="form-box">
            <div class="form-item">
                <label class="form-label">{{$t('ModelManagement["图标选择"]')}}<span class="color-danger">*</span></label>
                <div class="select-wrapper">
                    <div class="select-box clearfix" @click.stop.prevent="toggleDrop" :class="{'active': iconInfo.isIconDrop, 'disabled': isReadOnly}">
                        <div class="select-content">
                            <i :class="baseInfo['bk_obj_icon']"></i>
                        </div>
                        <span class="arrow"><i class="bk-icon icon-angle-down"></i></span>
                    </div>
                    <div class="mask" v-if="iconInfo.isIconDrop" @click="closeDrop"></div>
                    <div class="select-list" v-if="iconInfo.isIconDrop">
                        <ul class="clearfix select-icon-list">
                            <li v-tooltip="{content: language === 'zh-CN' ? item.nameZh : item.nameEn}" v-for="(item,index) in curIconList" :class="{'active': false}" :key="index" @click.stop.prevent="chooseIcon(item)">
                                <i :class="item.value"></i>
                            </li>
                        </ul>
                        <div class="page-wrapper clearfix">
                            <div class="input-wrapper">
                                <input type="text" class="cmdb-form-input" v-model="iconInfo.searchText" :placeholder="$t('ModelManagement[\'请输入关键词\']')">
                                <i class="bk-icon icon-search"></i>
                            </div>
                            <ul class="clearfix page">
                                <li v-for="(page, index) in iconInfo.totalPage"
                                class="page-item" :class="{'cur-page': iconInfo.curPage === page}"
                                :key="index"
                                @click="iconInfo.curPage = page"
                                >
                                    {{page}}
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>
            </div>
            <div class="form-item">
                <label for="name" class="form-label">{{$t('ModelManagement["中文名称"]')}}<span class="color-danger"> * </span></label>
                <div class="input-box">
                    <input type="text" id="name" class="cmdb-form-input"
                        maxlength="20"
                        :disabled="isReadOnly || activeModel['ispre']"
                        v-model.trim="baseInfo['bk_obj_name']"
                        :data-vv-name="$t('ModelManagement[\'中文名称\']')"
                        v-validate="'required|singlechar'">
                    <div v-show="errors.has($t('ModelManagement[\'中文名称\']'))" class="error-msg color-danger">{{ errors.first($t('ModelManagement[\'中文名称\']')) }}</div>
                </div>
            </div>
            <div class="form-item">
                <label for="name" class="form-label">{{$t('ModelManagement["英文名称"]')}}<span class="color-danger"> * </span></label>
                <div class="input-box">
                    <input type="text" id="desc" class="cmdb-form-input"
                        maxlength="20"
                        :disabled="isEdit"
                        v-model.trim="baseInfo['bk_obj_id']"
                        :data-vv-name="$t('ModelManagement[\'英文名称\']')"
                        v-validate="'required|modelId'">
                    <div v-show="errors.has($t('ModelManagement[\'英文名称\']'))" class="error-msg color-danger">{{ errors.first($t('ModelManagement[\'英文名称\']')) }}</div>
                </div>
            </div>
        </div>
        <footer class="footer-btn" v-if="!isReadOnly">
            <bk-button type="primary" @click="saveBaseInfo" :loading="$loading('createModel')">{{$t('Common["确定"]')}}</bk-button>
            <bk-button class="default" type="default" :title="$t('Common[\'取消\']')" @click="cancel">{{$t('Common["取消"]')}}</bk-button>
        </footer>
    </div>
</template>

<script>
    import iconList from '@/assets/json/model-icon'
    import { mapGetters, mapActions } from 'vuex'
    export default {
        props: {
            isReadOnly: {
                type: Boolean,
                default: false
            },
            isEdit: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                iconInfo: {
                    isIconDrop: false,
                    searchText: '',
                    list: iconList,
                    count: 0,
                    curPage: 1,
                    totalPage: 0,
                    size: 24
                },
                baseInfo: {
                    bk_obj_name: '',
                    bk_obj_id: '',
                    bk_classification_id: '',
                    bk_supplier_account: '',
                    bk_obj_icon: 'icon-cc-default'
                },
                baseInfoCopy: {}
            }
        },
        computed: {
            ...mapGetters([
                'language',
                'userName',
                'supplierAccount'
            ]),
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            isMainLine () {
                return this.activeModel.hasOwnProperty('bk_asst_obj_id') && this.activeModel['bk_asst_obj_id']
            },
            curIconList () {
                let {
                    searchText,
                    list,
                    curPage,
                    size
                } = this.iconInfo
                let curIconList = list
                if (searchText.length) {
                    curIconList = list.filter(icon => {
                        return icon.nameZh.toLowerCase().indexOf(searchText.toLowerCase()) > -1 || icon.nameEn.toLowerCase().indexOf(searchText.toLowerCase()) > -1
                    })
                }
                this.iconInfo.count = list.length
                this.iconInfo.totalPage = Math.ceil(list.length / size)
                return curIconList.slice((curPage - 1) * size, curPage * size)
            }
        },
        created () {
            if (this.isEdit) {
                this.getObjInfo()
            } else {
                this.baseInfo['bk_classification_id'] = this.activeModel['bk_classification_id']
                this.baseInfoCopy = this.$tools.clone(this.baseInfo)
            }
        },
        methods: {
            ...mapActions('objectModel', [
                'createObject',
                'updateObject',
                'searchObjects'
            ]),
            ...mapActions('objectMainLineModule', [
                'createMainlineObject'
            ]),
            isCloseConfirmShow () {
                if (JSON.stringify(this.baseInfo) !== JSON.stringify(this.baseInfoCopy)) {
                    return true
                }
                return false
            },
            chooseIcon (item) {
                this.baseInfo['bk_obj_icon'] = item.value
                this.closeDrop()
            },
            toggleDrop () {
                if (this.isReadOnly) {
                    return
                }
                this.iconInfo.isIconDrop = !this.iconInfo.isIconDrop
            },
            closeDrop () {
                this.iconInfo.isIconDrop = false
            },
            async saveBaseInfo () {
                if (!await this.$validator.validateAll()) {
                    return
                }
                let params = {
                    bk_supplier_account: this.supplierAccount,
                    bk_obj_name: this.baseInfo['bk_obj_name'],
                    bk_obj_icon: this.baseInfo['bk_obj_icon'],
                    bk_classification_id: this.activeModel['bk_classification_id'],
                    bk_obj_id: this.baseInfo['bk_obj_id']
                }
                if (this.isEdit) {
                    if (this.baseInfo['bk_obj_name'] === this.baseInfoCopy['bk_obj_name'] && this.baseInfo['bk_obj_icon'] === this.baseInfoCopy['bk_obj_icon']) {
                        this.cancel()
                        return
                    }
                    await this.updateObject({
                        id: this.activeModel['id'],
                        params: {...params, ...{modifier: this.userName}}
                    }).then(() => {
                        this.$http.cancel('post_searchClassificationsObjects')
                    })
                    this.$emit('confirm')
                } else {
                    let res = null
                    if (this.isMainLine) {
                        res = await this.createMainlineObject({
                            params: {...params, ...{creator: this.userName, bk_asst_obj_id: this.activeModel['bk_asst_obj_id']}},
                            config: {
                                requestId: 'createModel'
                            }
                        }).then(res => {
                            this.$http.cancel('post_searchClassificationsObjects')
                            return res
                        })
                    } else {
                        res = await this.createObject({
                            params: {...params, ...{creator: this.userName}},
                            config: {
                                requestId: 'createModel'
                            }
                        }).then(res => {
                            this.$http.cancel('post_searchClassificationsObjects')
                            return res
                        })
                    }
                    this.$store.commit('objectModel/setActiveModel', res)
                    this.baseInfoCopy = this.$tools.clone(this.baseInfo)
                    this.$emit('createObject')
                }
            },
            async getObjInfo () {
                const res = await this.searchObjects({
                    params: {
                        bk_obj_id: this.activeModel['bk_obj_id']
                    }
                })
                this.baseInfo = res[0]
                this.baseInfoCopy = this.$tools.clone(res[0])
            },
            cancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .base-info-wrapper {
        padding: 20px 0;
        width: 100%;
        .form-box {
            display: flex;
        }
        .form-item {
            margin-right: 30px;
            font-size: 0;
            &:first-child {
                width: 150px;
            }
            &:last-child {
                margin-right: 0;
                text-align: right;
            }
            .form-label {
                display: inline-block;
                width: 65px;
                vertical-align: top;
                text-align: right;
                font-size: 14px;
                line-height: 36px;
                margin-right: 10px;
                span{
                    display: inline-block;
                }
            }
            .input-box {
                display: inline-block;
                text-align: left;
                width: 198px;
                .error-msg {
                    font-size: 12px;
                    line-height: 1;
                }
            }
            .cmdb-form-input {
                width: 198px;
            }
        }
        .select-wrapper {
            position: relative;
            display: inline-block;
            vertical-align: top;
            .mask {
                z-index: 10;
            }
            .select-box {
                border: 1px solid $cmdbBorderColor;
                text-align: center;
                height: 36px;
                line-height: 34px;
                cursor: pointer;
                &.disabled {
                    background: #fafafa;
                    cursor: not-allowed;
                }
                .select-content {
                    position: relative;
                    float: left;
                    height: 34px;
                    padding: 0 3px 0 10px;
                    color: $cmdbBorderFocusColor;
                    font-size: 24px;
                    i {
                        position: relative;
                        top: -1px;
                    }
                }
                .arrow {
                    float: left;
                    width: 26px;
                    font-size: 12px;
                }
            }
            .select-list {
                position: absolute;
                padding: 10px;
                top: 44px;
                left: 0;
                width: 382px;
                height: 248px;
                border: 1px solid $cmdbFnMainColor;
                z-index: 500;
                background: #fff;
                box-shadow: 0 2px 2px rgba(0, 0, 0, .1);
                overflow: auto;
                @include scrollbar;
                .select-icon-list {
                    padding: 0;
                    margin: 0;
                    width: 360px;
                    height: 184px;
                    li {
                        width: 60px;
                        height: 46px;
                        text-align: center;
                        line-height: 46px;
                        float: left;
                        cursor: pointer;
                        &.active {
                            color: $cmdbBorderFocusColor;
                            background: #e2efff;
                        }
                        i {
                            font-size: 24px;
                        }
                        &:hover {
                            background: #e2efff;
                        }
                        &:nth-child(6n) {
                            margin-right: 0;
                        }
                    }
                }
                .page-wrapper {
                    padding: 15px 18px 5px;
                    .input-wrapper {
                        float: left;
                        position: relative;
                        vertical-align: bottom;
                        color: $cmdbBorderColor;
                        input {
                            width: 116px;
                            height: 22px;
                            font-size: 12px;
                            padding: 0 25px 0 5px;
                        }
                        .bk-icon {
                            position: absolute;
                            font-size: 12px;
                            top: 5px;
                            right: 8px;
                        }
                    }
                }
                .page{
                    float: right;
                    li{
                        text-align: center;
                        float: left;
                        margin-right: 5px;
                        width: 22px;
                        height: 22px;
                        line-height: 20px;
                        border-radius: 2px;
                        font-size: 12px;
                        cursor: pointer;
                        border: 1px solid $cmdbBorderColor;
                        &.cur-page{
                            color: #fff;
                            background: $cmdbBorderFocusColor;
                            border-color: $cmdbBorderFocusColor;
                        }
                        &:last-child{
                            margin: 0;
                        }
                    }
                }
            }
        }
    }
</style>
