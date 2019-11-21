<template>
    <ul class="form-list-wrapper">
        <li class="form-item clearfix" v-for="(item, index) in list" :key="index">
            <div class="list-label">
                {{$t('选项')}}{{index + 1}}
            </div>
            <div class="list-name">
                <div class="cmdb-form-item" :class="{ 'is-error': errors.has(`name${index}`) }">
                    <bk-input type="text"
                        class="cmdb-form-input"
                        :placeholder="$t('请输入名称英文数字')"
                        v-model.trim="item.name"
                        v-validate="`required|enumName|repeat:${getOtherName(index)}`"
                        @input="handleInput"
                        :disabled="isReadOnly"
                        :name="`name${index}`">
                    </bk-input>
                    <p class="form-error">{{errors.first(`name${index}`)}}</p>
                </div>
            </div>
            <button class="list-btn" @click="deleteList(index)" :disabled="list.length === 1 || isReadOnly">
                <i class="icon-cc-del"></i>
            </button>
            <button class="list-btn" @click="addList" :disabled="isReadOnly" v-if="index === list.length - 1">
                <i class="bk-icon icon-plus"></i>
            </button>
        </li>
    </ul>
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
                list: [{ name: '' }]
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
            getOtherName (index) {
                const nameList = []
                this.list.map((item, _index) => {
                    if (index !== _index) {
                        nameList.push(item.name)
                    }
                })
                return nameList.join(',')
            },
            initValue () {
                if (this.value === '') {
                    this.list = [{ name: '' }]
                } else {
                    this.list = this.value.map(name => ({ name }))
                }
            },
            handleInput () {
                this.$nextTick(async () => {
                    const res = await this.$validator.validateAll()
                    if (res) {
                        const list = this.list.map(item => item.name)
                        this.$emit('input', list)
                    }
                })
            },
            addList () {
                this.list.push({ name: '' })
            },
            deleteList (index) {
                this.list.splice(index, 1)
                this.handleInput()
            },
            validate () {
                return this.$validator.validateAll()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-list-wrapper {
        >.form-item {
            font-size: 0;
            &:not(:first-child) {
                margin-top: 15px;
            }
            .list-label {
                font-size: 14px;
                line-height: 36px;
                text-align: center;
                width: 55px;
            }
            .list-name {
                float: left;
                width: 200px;
                input {
                    width: 100%;
                }
            }
            .list-btn {
                display: inline-block;
                width: 32px;
                height: 32px;
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
