<template>
    <ul class="form-enum-wrapper">
        <li class="form-item clearfix" v-for="(item, index) in enumList" :key="index">
            <div class="clearfix">
                <div class="enum-default cmdb-form-radio">
                    <input type="radio" 
                    :value="index" 
                    name="enum-radio" 
                    v-model="defaultIndex" 
                    v-tooltip="$t('ModelManagement[\'将设置为下拉选项默认选项\']')"
                    @change="handleChange(defaultIndex)" :disabled="isReadOnly">
                </div>
                <div class="enum-label">
                    {{$t('ModelManagement["枚举"]')}}{{index + 1}}
                </div>
            </div>
            <div class="enum-id">
                <div class="cmdb-form-item" :class="{'is-error': errors.has(`id${index}`)}">
                    <input type="text"
                        class="cmdb-form-input"
                        :placeholder="$t('ModelManagement[\'请输入ID\']')"
                        v-model.trim="item.id"
                        v-validate="`required|enumId|repeat:${getOtherId(index)}`"
                        @input="handleInput"
                        :disabled="isReadOnly"
                        :name="`id${index}`">
                    <p class="form-error">{{errors.first(`id${index}`)}}</p>
                </div>
            </div>
            <div class="enum-name">
                <div class="cmdb-form-item" :class="{'is-error': errors.has(`name${index}`)}">
                    <input type="text"
                        class="cmdb-form-input"
                        :placeholder="$t('ModelManagement[\'请输入名称英文数字\']')"
                        v-model.trim="item.name"
                        v-validate="`required|enumName|repeat:${getOtherName(index)}`"
                        @input="handleInput"
                        :disabled="isReadOnly"
                        :name="`name${index}`">
                    <p class="form-error">{{errors.first(`name${index}`)}}</p>
                </div>
            </div>
            <button class="enum-btn" @click="deleteEnum(index)" :disabled="enumList.length === 1 || isReadOnly">
                <i class="icon-cc-del"></i>
            </button>
            <button class="enum-btn" @click="addEnum" :disabled="isReadOnly" v-if="index === enumList.length - 1">
                <i class="bk-icon icon-plus"></i>
            </button>
        </li>
    </ul>
</template>

<script>
    export default {
        props: {
            value: {
                default: ''
            },
            isReadOnly: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                enumList: [{
                    id: '',
                    is_default: true,
                    name: ''
                }],
                defaultIndex: 0
            }
        },
        watch: {
            value () {
                this.initValue()
            }
        },
        created () {
            this.initValue()
        },
        methods: {
            getOtherId (index) {
                let idList = []
                this.enumList.map((item, enumIndex) => {
                    if (index !== enumIndex) {
                        idList.push(item.id)
                    }
                })
                return idList.join(',')
            },
            getOtherName (index) {
                let nameList = []
                this.enumList.map((item, enumIndex) => {
                    if (index !== enumIndex) {
                        nameList.push(item.name)
                    }
                })
                return nameList.join(',')
            },
            initValue () {
                if (this.value === '') {
                    this.enumList = [{
                        id: '',
                        is_default: true,
                        name: ''
                    }]
                } else {
                    this.enumList = this.value
                    this.defaultIndex = this.enumList.findIndex(({is_default: isDefault}) => isDefault)
                }
            },
            handleInput () {
                this.$nextTick(async () => {
                    const res = await this.$validator.validateAll()
                    if (res) {
                        this.$emit('input', this.enumList)
                    }
                })
            },
            handleChange (index) {
                let defaultItem = this.enumList.find(({is_default: isDefault}) => isDefault)
                if (defaultItem) {
                    defaultItem['is_default'] = false
                }
                this.enumList[index]['is_default'] = true
                this.handleInput()
            },
            addEnum () {
                this.enumList.push({
                    id: '',
                    is_default: false,
                    name: ''
                })
                this.handleInput()
            },
            deleteEnum (index) {
                this.enumList.splice(index, 1)
                if (this.defaultIndex === index) {
                    this.defaultIndex = 0
                    this.enumList[0]['is_default'] = true
                }
            },
            validate () {
                return this.$validator.validateAll()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-enum-wrapper {
        >.form-item {
            font-size: 0;
            &:not(:first-child) {
                margin-top: 15px;
            }
            .enum-default {
                float: left;
                margin: 0;
                padding: 8px 5px 0 0;
                height: 36px;
                font-size: 16px;
                line-height: 1;
            }
            .enum-label {
                float: left;
                padding-right: 10px;
                font-size: 14px;
                line-height: 36px;
                text-align: center;
                width: 55px;
            }
            .enum-id {
                float: left;
                width: 90px;
                margin-right: 10px;
                input {
                    width: 100%;
                }
            }
            .enum-name {
                float: left;
                width: 180px;
                input {
                    width: 100%;
                }
            }
            .enum-btn {
                display: inline-block;
                width: 36px;
                height: 36px;
                margin-left: 5px;
                vertical-align: middle;
                text-align: center;
                font-size: 14px;
                line-height: 1;
                border: 1px solid $cmdbFnMainColor;
                background-color: $cmdbDefaultColor;
                outline: 0;
                &:disabled {
                    cursor: not-allowed;
                    background-color: #eee;
                    border-color: #eee;
                    color: $cmdbFnMainColor;
                }
            }
        }
    }
</style>
