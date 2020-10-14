<template>
    <div class="form-enum-layout">
        <div class="toolbar">
            <p class="title">{{$t('枚举值')}}</p>
            <i
                v-bk-tooltips.top-start="$t('按照0-9a-z排序')"
                :class="['sort-icon', `icon-cc-sort-${order > 0 ? 'up' : 'down'}`]"
                @click="handleSort">
            </i>
        </div>
        <vue-draggable
            class="form-enum-wrapper"
            tag="ul"
            v-model="enumList"
            :options="dragOptions"
            @end="handleDragEnd">
            <li class="form-item" v-for="(item, index) in enumList" :key="index">
                <div class="enum-id">
                    <div class="cmdb-form-item" :class="{ 'is-error': errors.has(`id${index}`) }">
                        <bk-input type="text"
                            class="cmdb-form-input"
                            :placeholder="$t('请输入ID')"
                            v-model.trim="item.id"
                            v-validate="`required|enumId|length:128|repeat:${getOtherId(index)}`"
                            @input="handleInput"
                            :disabled="isReadOnly"
                            :name="`id${index}`"
                            :ref="`id${index}`">
                        </bk-input>
                        <p class="form-error" :title="errors.first(`id${index}`)">{{errors.first(`id${index}`)}}</p>
                    </div>
                </div>
                <div class="enum-name">
                    <div class="cmdb-form-item" :class="{ 'is-error': errors.has(`name${index}`) }">
                        <bk-input type="text"
                            class="cmdb-form-input"
                            :placeholder="$t('请输入值')"
                            v-model.trim="item.name"
                            v-validate="`required|enumName|length:128|repeat:${getOtherName(index)}`"
                            @input="handleInput"
                            :disabled="isReadOnly"
                            :name="`name${index}`">
                        </bk-input>
                        <p class="form-error" :title="errors.first(`name${index}`)">{{errors.first(`name${index}`)}}</p>
                    </div>
                </div>
                <bk-button text class="enum-btn" @click="deleteEnum(index)" :disabled="enumList.length === 1 || isReadOnly">
                    <i class="bk-icon icon-minus-circle-shape"></i>
                </bk-button>
                <bk-button text class="enum-btn" @click="addEnum(index)" :disabled="isReadOnly" v-if="index === enumList.length - 1">
                    <i class="bk-icon icon-plus-circle-shape"></i>
                </bk-button>
            </li>
        </vue-draggable>
        <div class="default-setting">
            <p class="title mb10">{{$t('默认值设置')}}</p>
            <bk-select style="width: 100%;"
                :clearable="false"
                :disabled="isReadOnly"
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
    import vueDraggable from 'vuedraggable'
    export default {
        components: {
            vueDraggable
        },
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
                enumList: [this.generateEnum()],
                defaultIndex: 0,
                settingList: [],
                defaultValue: '',
                dragOptions: {
                    animation: 300,
                    disabled: false,
                    filter: '.enum-btn, .enum-id, .enum-name',
                    preventOnFilter: false,
                    ghostClass: 'ghost'
                },
                order: 1
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
                    if (this.settingList.length && this.defaultIndex > -1) {
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
                this.enumList.forEach((item, enumIndex) => {
                    if (index !== enumIndex) {
                        idList.push(item.id)
                    }
                })
                return idList.join(',')
            },
            getOtherName (index) {
                const nameList = []
                this.enumList.forEach((item, enumIndex) => {
                    if (index !== enumIndex) {
                        nameList.push(item.name)
                    }
                })
                return nameList.join(',')
            },
            initValue () {
                if (this.value === '') {
                    this.enumList = [this.generateEnum()]
                } else {
                    this.enumList = this.value.map(data => ({ ...data, type: 'text' }))
                    this.defaultIndex = this.enumList.findIndex(({ is_default: isDefault }) => isDefault)
                }
            },
            handleInput () {
                this.$emit('input', this.enumList)
            },
            addEnum (index) {
                this.enumList.push(this.generateEnum({ is_default: false }))
                this.$nextTick(() => {
                    this.$refs[`id${index + 1}`] && this.$refs[`id${index + 1}`][0].focus()
                })
            },
            deleteEnum (index) {
                this.enumList.splice(index, 1)
                if (this.defaultIndex === index) {
                    this.defaultIndex = 0
                    this.enumList[0]['is_default'] = true
                }
                this.handleInput()
            },
            generateEnum (settings = {}) {
                const defaults = {
                    id: '',
                    is_default: true,
                    name: '',
                    type: 'text'
                }
                return { ...defaults, ...settings }
            },
            validate () {
                return this.$validator.validateAll()
            },
            handleSettingDefault (id) {
                const itemIndex = this.enumList.findIndex(item => item.id === id)
                if (itemIndex > -1) {
                    this.defaultIndex = itemIndex
                    this.enumList.forEach(item => {
                        item.is_default = item.id === id
                    })

                    this.$emit('input', this.enumList)
                }
            },
            handleDragEnd () {
                this.$emit('input', this.enumList)
            },
            handleSort () {
                this.order = this.order * -1
                this.enumList.sort((A, B) => {
                    return A.name.localeCompare(B.name, 'zh-Hans-CN', { sensitivity: 'accent' }) * this.order
                })

                this.$emit('input', this.enumList)
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
            position: relative;
            font-size: 0;
            margin-bottom: 16px;
            padding: 2px 2px 2px 28px;
            cursor: move;

            &::before {
                content: '';
                position: absolute;
                top: 12px;
                left: 8px;
                width: 3px;
                height: 3px;
                border-radius: 50%;
                background-color: #979ba5;
                box-shadow: 0 5px 0 0 #979ba5,
                    0 10px 0 0 #979ba5,
                    5px 0 0 0 #979ba5,
                    5px 5px 0 0 #979ba5,
                    5px 10px 0 0 #979ba5;
            }

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

                &.is-disabled {
                    color: #dcdee5;
                }
                &:not(.is-disabled):hover {
                    color: #979ba5;
                }
            }
        }
    }

    .toolbar {
        display: flex;
        margin-bottom: 10px;
        align-items: center;
        line-height: 20px;

        .sort-icon {
            width: 20px;
            height: 20px;
            margin-left: 10px;
            border: 1px solid #c4c6cc;
            background: #fff;
            border-radius: 2px;
            font-size: 16px;
            line-height: 18px;
            text-align: center;
            color: #c4c6cc;
            cursor: pointer;

            &:hover {
                color: #979ba5;
            }
        }
    }

    .ghost {
        border: 1px dashed $cmdbBorderFocusColor;
    }
</style>
