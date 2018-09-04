<template>
    <ul class="form-enum-wrapper">
        <li class="form-item clearfix" v-for="(item, index) in enumList" :key="index">
            <div class="enum-default">
                <input type="radio" :value="index" name="enum-radio" v-model="defaultIndex" @change="handleChange(defaultIndex)" :disabled="isReadOnly">
            </div>
            <div class="enum-id">
                <input type="text"
                    class="cmdb-form-input"
                    :placeholder="$t('ModelManagement[\'请输入ID\']')"
                    v-model.trim="item.id"
                    v-validate="`required|enumId|repeat:${getOtherId(index)}`"
                    @input="handleInput"
                    :disabled="isReadOnly"
                    :name="`id${index}`">
                    <span v-show="errors.has(`id${index}`)" class="error-msg color-danger">{{ errors.first(`id${index}`) }}</span>
            </div>
            <div class="enum-name">
                <input type="text"
                    class="cmdb-form-input"
                    :placeholder="$t('ModelManagement[\'请输入名称英文数字\']')"
                    v-model.trim="item.name"
                    v-validate="`required|enumName|repeat:${getOtherName(index)}`"
                    @input="handleInput"
                    :disabled="isReadOnly"
                    :name="`name${index}`">
                    <span v-show="errors.has(`name${index}`)" class="error-msg color-danger">{{ errors.first(`name${index}`) }}</span>
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
            &:not(:first-child) {
                margin-top: 10px;
            }
            .enum-default {
                float: left;
                width: 80px;
                height: 30px;
                padding-right: 5px;
                font-size: 16px;
                text-align: right;
                line-height: 1;
                padding-top: 5px;
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
                width: 250px;
                input {
                    width: 100%;
                }
            }
            .enum-btn {
                display: inline-block;
                width: 30px;
                height: 30px;
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
