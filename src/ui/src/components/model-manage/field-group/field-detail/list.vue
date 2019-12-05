<template>
    <div class="form-list-layout">
        <p class="title mb10">{{$t('列表值')}}</p>
        <ul class="form-list-wrapper">
            <li class="form-item clearfix" v-for="(item, index) in list" :key="index">
                <div class="list-name">
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
                <bk-button text class="list-btn" @click="deleteList(index)" :disabled="list.length === 1 || isReadOnly">
                    <i class="bk-icon icon-minus-circle-shape"></i>
                </bk-button>
                <bk-button text class="list-btn" @click="addList" :disabled="isReadOnly" v-if="index === list.length - 1">
                    <i class="bk-icon icon-plus-circle-shape"></i>
                </bk-button>
            </li>
        </ul>
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
    .title {
        font-size: 14px;
    }
    .form-list-wrapper {
        .form-item {
            display: flex;
            align-items: center;
            font-size: 0;
            &:not(:first-child) {
                margin-top: 20px;
            }
            .list-name {
                float: left;
                width: 200px;
                input {
                    width: 100%;
                }
            }
            .list-btn {
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
