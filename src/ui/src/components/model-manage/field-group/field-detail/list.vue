<template>
    <div class="form-list-layout">
        <div class="toolbar">
            <p class="title">{{$t('列表值')}}</p>
            <i
                v-bk-tooltips.top-start="$t('按照0-9a-z排序')"
                :class="['sort-icon', `icon-cc-sort-${order > 0 ? 'up' : 'down'}`]"
                @click="handleSort">
            </i>
        </div>
        <vue-draggable
            class="form-list-wrapper"
            tag="ul"
            v-model="list"
            :options="dragOptions"
            @end="handleDragEnd">
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
                            :name="`name${index}`"
                            :ref="`name${index}`">
                        </bk-input>
                        <p class="form-error">{{errors.first(`name${index}`)}}</p>
                    </div>
                </div>
                <bk-button text class="list-btn" @click="deleteList(index)" :disabled="list.length === 1 || isReadOnly">
                    <i class="bk-icon icon-minus-circle-shape"></i>
                </bk-button>
                <bk-button text class="list-btn" @click="addList(index)" :disabled="isReadOnly" v-if="index === list.length - 1">
                    <i class="bk-icon icon-plus-circle-shape"></i>
                </bk-button>
            </li>
        </vue-draggable>
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
                list: [{ name: '' }],
                dragOptions: {
                    animation: 300,
                    disabled: false,
                    filter: '.list-btn, .list-name',
                    preventOnFilter: false,
                    ghostClass: 'ghost'
                },
                order: 1
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
            addList (index) {
                this.list.push({ name: '' })
                this.$nextTick(() => {
                    this.$refs[`name${index + 1}`] && this.$refs[`name${index + 1}`][0].focus()
                })
            },
            deleteList (index) {
                this.list.splice(index, 1)
                this.handleInput()
            },
            validate () {
                return this.$validator.validateAll()
            },
            handleDragEnd () {
                const list = this.list.map(item => item.name)
                this.$emit('input', list)
            },
            handleSort () {
                this.order = this.order * -1
                this.list.sort((A, B) => {
                    return A.name.localeCompare(B.name, 'zh-Hans-CN', { sensitivity: 'accent' }) * this.order
                })

                const list = this.list.map(item => item.name)
                this.$emit('input', list)
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
            position: relative;
            padding: 2px 2px 2px 28px;
            font-size: 0;
            cursor: move;

            &:not(:first-child) {
                margin-top: 16px;
            }

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
