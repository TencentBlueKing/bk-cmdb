<template>
    <bk-dropdown-menu trigger="click" :disabled="disabled">
        <bk-button class="clipboard-trigger" type="default" slot="dropdown-trigger"
            :disabled="disabled">
            {{$t('Common["复制"]')}}
            <i class="bk-icon icon-angle-down"></i>
        </bk-button>
        <ul class="clipboard-list" slot="dropdown-content">
            <li v-for="(item, index) in list"
                class="clipboard-item"
                :key="index"
                @click="handleClick(item)">
                {{item[labelKey]}}
            </li>
        </ul>
    </bk-dropdown-menu>
</template>

<script>
    export default {
        name: 'cmdb-clipboard-selector',
        props: {
            disabled: {
                type: Boolean,
                default: false
            },
            list: {
                type: Array,
                default () {
                    return []
                }
            },
            idKey: {
                type: String,
                default: 'id'
            },
            labelKey: {
                type: String,
                default: 'name'
            }
        },
        methods: {
            handleClick (item) {
                this.$emit('on-copy', item)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .clipboard-trigger{
        padding: 0 16px;
        .icon-angle-down {
            font-size: 12px;
            top: 0;
        }
    }
    .clipboard-list{
        width: 100%;
        font-size: 14px;
        line-height: 40px;
        max-height: 160px;
        @include scrollbar-y;
        &::-webkit-scrollbar{
            width: 3px;
            height: 3px;
        }
        .clipboard-item{
            padding: 0 15px;
            cursor: pointer;
            @include ellipsis;
            &:hover{
                background-color: #ebf4ff;
                color: #3c96ff;
            }
        }
    }
</style>