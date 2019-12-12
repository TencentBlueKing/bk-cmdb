<template>
    <div class="form-enum-layout">
        <p class="title mb10">{{$t('枚举值')}}</p>
        <ul class="form-enum-wrapper">
            <li class="form-item" v-for="(item, index) in enumList" :key="index">
                <div class="enum-id">
                    <div class="cmdb-form-item" :class="{ 'is-error': errors.has(`id${index}`) }">
                        <bk-input type="text"
                            class="cmdb-form-input"
                            :placeholder="$t('请输入ID')"
                            v-model.trim="item.id"
                            v-validate="`required|enumId|repeat:${getOtherId(index)}`"
                            @input="handleInput"
                            :disabled="isReadOnly"
                            :name="`id${index}`">
                        </bk-input>
                        <p class="form-error">{{errors.first(`id${index}`)}}</p>
                    </div>
                </div>
                <div class="enum-name">
                    <div class="cmdb-form-item" :class="{ 'is-error': errors.has(`name${index}`) }">
                        <bk-input type="text"
                            class="cmdb-form-input"
                            :placeholder="$t('请输入值')"
                            v-model.trim="item.name"
                            v-validate="`required|enumName|repeat:${getOtherName(index)}`"
                            @input="handleInput"
                            :disabled="isReadOnly"
                            :name="`name${index}`">
                        </bk-input>
                        <p class="form-error">{{errors.first(`name${index}`)}}</p>
                    </div>
                </div>
                <bk-button text class="enum-btn" @click="deleteEnum(index)" :disabled="enumList.length === 1 || isReadOnly">
                    <i class="bk-icon icon-minus-circle-shape"></i>
                </bk-button>
                <bk-button text class="enum-btn" @click="addEnum" :disabled="isReadOnly" v-if="index === enumList.length - 1">
                    <i class="bk-icon icon-plus-circle-shape"></i>
                </bk-button>
            </li>
        </ul>
        <div class="default-setting">
            <p class="title mb10">{{$t('默认值设置')}}</p>
            <bk-select style="width: 100%;"
                :clearable="false"
                v-model="defaultValue"
                @change="handleSettingDefault">
                <bk-option v-for="option in settingList"
                    :key="option.id"
                    :id="option.id"
                    :name="option.name">
                </bk-option>
            </bk-select>
        </div>
    </div>
</template>

<script>
    export default {
        props: {
            value: {
                type: [Array, String],
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
                defaultIndex: 0,
                settingList: [],
                defaultValue: ''
            }
        },
        watch: {
            value () {
                this.initValue()
            },
            enumList: {
                deep: true,
                handler (value) {
                    this.settingList = (value || []).filter(item => item.id && item.name)
                    if (this.settingList.length) {
                        this.defaultValue = this.settingList[this.defaultIndex].id
                    }
                }
            }
        },
        created () {
            this.initValue()
        },
        methods: {
            getOtherId (index) {
                const idList = []
                this.enumList.map((item, enumIndex) => {
                    if (index !== enumIndex) {
                        idList.push(item.id)
                    }
                })
                return idList.join(',')
            },
            getOtherName (index) {
                const nameList = []
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
                    this.defaultIndex = this.enumList.findIndex(({ is_default: isDefault }) => isDefault)
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
                this.handleInput()
            },
            validate () {
                return this.$validator.validateAll()
            },
            handleSettingDefault (id) {
                const itemIndex = this.enumList.findIndex(item => item.id === id)
                if (itemIndex > -1) {
                    this.defaultIndex = itemIndex
                    this.enumList = this.enumList.map(item => {
                        item.is_default = item.id === id
                        return item
                    })
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .title {
        font-size: 14px;
    }
    .form-enum-wrapper {
        .form-item {
            display: flex;
            align-items: center;
            font-size: 0;
            margin-bottom: 20px;
            .enum-id {
                width: 90px;
                margin-right: 10px;
                input {
                    width: 100%;
                }
            }
            .enum-name {
                width: 180px;
                input {
                    width: 100%;
                }
            }
            .enum-btn {
                font-size: 0;
                color: #c4c6cc;
                margin: -2px 0 0 6px;
                .bk-icon {
                    width: 18px;
                    height: 18px;
                    line-height: 18px;
                    font-size: 18px;
                    text-align: center;
                }
            }
        }
    }
</style>
