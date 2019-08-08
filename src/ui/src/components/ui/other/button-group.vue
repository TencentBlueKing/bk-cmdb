<template>
    <div class="button-group">
        <template v-if="expand">
            <template v-for="button in available">
                <span class="button-item"
                    v-if="button.auth"
                    v-cursor="button.auth"
                    :key="button.id">
                    <bk-button class="button-item"
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
            @show="toggleDropdownState(true)"
            @hide="toggleDropdownState(false)">
            <bk-button slot="dropdown-trigger">
                <span>{{$t('更多')}}</span>
                <i :class="['bk-icon icon-angle-down', { 'icon-flip': isDropdownShow }]"></i>
            </bk-button>
            <ul class="dropdown-list" slot="dropdown-content">
                <template v-for="button in available">
                    <li class="dropdown-item"
                        v-if="button.auth"
                        v-cursor="button.auth"
                        :key="button.id"
                        :class="{
                            'is-disabled': button.disabled
                        }"
                        @click="handleClick(button)">
                        <slot :name="button.id">
                            {{button.text}}
                        </slot>
                    </li>
                    <li class="dropdown-item"
                        v-else
                        :key="button.id"
                        :class="{
                            'is-disabled': button.disabled
                        }"
                        @click="handleClick(button)">
                        <slot :name="button.id">
                            {{button.text}}
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
            expand: Boolean
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
            padding: 0 15px;
            cursor: pointer;
            @include ellipsis;
            &:hover {
                background-color: #ebf4ff;
                color: #3c96ff;
            }
            &.is-disabled {
                background-color: #fff;
                color: #c4c6cc;
                cursor: not-allowed;
            }
        }
    }
</style>
