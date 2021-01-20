<template>
    <div class="button-group">
        <template v-if="expand">
            <template v-for="button in available">
                <span class="button-item"
                    v-if="button.auth"
                    v-cursor="button.auth"
                    :key="button.id">
                    <bk-button class="button-item"
                        size="normal"
                        :theme="button.theme || 'default'"
                        :disabled="button.disabled || false"
                        @click="handleClick(button)">
                        <slot :name="button.id">
                            {{button.text}}
                        </slot>
                    </bk-button>
                </span>
                <bk-button
                    v-else
                    size="normal"
                    :theme="button.theme || 'default'"
                    :disabled="button.disabled || false"
                    :key="button.id"
                    @click="handleClick(button)">
                    <slot :name="button.id">
                        {{button.text}}
                    </slot>
                </bk-button>
            </template>
        </template>
        <bk-dropdown-menu
            v-else
            trigger="click"
            font-size="medium"
            @show="toggleDropdownState(true)"
            @hide="toggleDropdownState(false)">
            <bk-button slot="dropdown-trigger">
                <span>{{triggerText || $t('更多')}}</span>
                <i :class="['bk-icon icon-angle-down', { 'icon-flip': isDropdownShow }]"></i>
            </bk-button>
            <ul class="dropdown-list" slot="dropdown-content">
                <template v-for="button in available">
                    <li class="dropdown-item" v-if="button.auth" :key="button.id" v-bk-tooltips="getTooltips(button.tooltips)">
                        <slot :name="button.id">
                            <cmdb-auth style="display: block;" :auth="button.auth">
                                <bk-button slot-scope="{ disabled }"
                                    class="dropdown-item-btn"
                                    text
                                    theme="primary"
                                    :disabled="button.disabled || disabled"
                                    @click="handleClick(button)">
                                    {{button.text}}
                                </bk-button>
                            </cmdb-auth>
                        </slot>
                    </li>
                    <li class="dropdown-item" v-else :key="button.id" v-bk-tooltips="getTooltips(button.tooltips)">
                        <slot :name="button.id">
                            <bk-button text theme="primary"
                                class="dropdown-item-btn"
                                :disabled="button.disabled"
                                @click="handleClick(button)">
                                {{button.text}}
                            </bk-button>
                        </slot>
                    </li>
                </template>
            </ul>
        </bk-dropdown-menu>
    </div>
</template>

<script>
    export default {
        props: {
            buttons: {
                type: Array,
                required: true
            },
            expand: Boolean,
            triggerText: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                isDropdownShow: false
            }
        },
        computed: {
            available () {
                return this.buttons.filter(button => {
                    if (button.hasOwnProperty('available')) {
                        return button.available
                    }
                    return true
                })
            }
        },
        methods: {
            handleClick (button) {
                if (button.disabled) {
                    return false
                }
                if (typeof button.handler === 'function') {
                    button.handler.call(null)
                }
            },
            toggleDropdownState (state) {
                this.isDropdownShow = state
            },
            getTooltips (tooltips) {
                let tooltipsSettings = { disabled: true }
                if (tooltips) {
                    tooltipsSettings.disabled = false
                    tooltipsSettings.interactive = true
                    tooltipsSettings.boundary = 'window'
                    const type = Object.prototype.toString.call(tooltips)
                    if (type === '[object String]') {
                        tooltipsSettings.content = tooltips
                    } else if (type === '[object Object]') {
                        tooltipsSettings = { ...tooltipsSettings, ...tooltips }
                    }
                }
                return tooltipsSettings
            }
        }
    }
</script>

<style lang="scss" scoped>
    .button-group {
        display: inline-block;
    }
    .button-item {
        display: inline-block;
        margin-right: 10px;
        &:last-child {
            margin-right: 0;
        }
    }
    .dropdown-list {
        .dropdown-item {
            font-size: 14px;
            height: 36px;
            line-height: 36px;
            cursor: pointer;
            @include ellipsis;
            .dropdown-item-btn {
                width: 100%;
                padding: 0 15px;
                height: auto;
                text-align: left;
                color: #63656e;
                &:disabled {
                    color: #dcdee5 !important;
                    background-color: transparent !important;
                }
            }
            &:hover {
                .dropdown-item-btn {
                    background-color: #ebf4ff;
                    color: #3c96ff;
                }
            }
        }
    }
</style>
