<template>
    <div class="bk-combobox"
         :class="[extCls, {'open': open}]"
         v-clickoutside="close">
        <div class="bk-combobox-wrapper">
            <input class="bk-combobox-input" 
                    :class="{'active': open}" 
                    :placeholder="placeholder" 
                    @input="onInput"
                    v-model="localValue"/>
            <div class="bk-combobox-icon-clear" v-if="localValue.length > 0" @click="localValue=''">
                <span title="clear"><i class="bk-icon icon-close"></i></span>
            </div>
            <div class="bk-combobox-icon-box" @click="openFn">
                <i class="bk-icon icon-angle-down bk-combobox-icon"></i>
            </div>

        </div>
        <transition name="toggle-slide">
            <div class="bk-combobox-list" v-show="open">
                <ul v-if="showList.length > 0">
                    <li class="bk-combobox-item" @click.stop="selectItem(item)" v-for="(item, index) in showList"
                        :class="{'bk-combobox-item-target': index===0 && localValue.length > 0}">
                        <div class="text">
                            {{ item }}
                        </div>
                    </li>
                </ul>
                <ul v-else>
                    <li class="bk-combobox-item" disabled>
                        <div class="text">
                            无匹配数据
                        </div>
                    </li>
                </ul>
            </div>
        </transition>
    </div>
</template>

<script>
    import clickoutside from '../../directives/clickoutside'

    export default {
        name: 'bk-combobox',
        props: {
            placeholder: {
                type: String,
                default: ''
            },
            list: {
                type: Array
            },
            value: {
                type: String,
                required: true
            },
            extCls: {
                type: String
            }
        },
        computed: {
            showList () {
                if (this.localValue === '') {
                    return this.list
                } else {
                    let newList = []
                    for (let item of this.list) {
                        if (item.indexOf(this.localValue) !== -1) {
                            newList.push(item)
                        }
                    }
                    return newList
                }
            }
        },
        data () {
            return {
                open: false,
                localValue: this.value
            }
        },
        watch: {
            localValue () {
                this.$emit('update:value', this.localValue)
            },
            value (newVal) {
                this.localValue = newVal
            }
        },
        directives: {
            clickoutside
        },
        methods: {
            selectItem (item) {
                this.localValue = item
                this.close()

                this.$emit('update:value', this.localValue)
                this.$emit('item-selected', item)
            },
            openFn () {
                if (!this.disabled) {
                    this.open = !this.open
                    this.$emit('visible-toggle', this.open)
                }
            },
            close () {
                this.open = false
                this.$emit('visible-toggle', this.open)
            },
            onInput () {
                this.open = true
                this.$emit('update:value', this.localValue)
                this.$emit('input', this.localValue)
            }
        }
    }
</script>

<style lang="scss">
   @import '../../bk-magic-ui/src/combox.scss'
</style>
