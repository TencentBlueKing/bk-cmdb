<template>
    <div class="custom-header">
        {{data.column.label}}
        <bk-popover class="batch-trigger"
            trigger="manual" theme="light" placement="bottom" ref="popover"
            :tippy-options="{
                hideOnClick: false,
                onHidden: handleHidden
            }">
            <i class="icon-cc-batch-update" @click="handleClick"></i>
            <bk-form form-type="vertical" slot="content">
                <bk-form-item :label="$t('批量编辑')">
                    <bk-select class="form-table-selector"
                        searchable
                        :placeholder="$t('请选择xx', { name: $t('资源目录') })"
                        v-model="selected">
                        <bk-option v-for="folder in folders"
                            :key="folder.id"
                            :id="folder.id"
                            :name="folder.name">
                        </bk-option>
                        <a href="javascript:void(0)" class="extension-link" slot="extension">
                            <i class="bk-icon icon-plus-circle"></i>
                            {{ $t('申请其他目录权限') }}
                        </a>
                    </bk-select>
                    <div class="selector-options">
                        <link-button class="selector-link-button" @click="handleConfirm">{{$t('确定')}}</link-button>
                        <link-button class="selector-link-button ml10" @click="handleCancel">{{$t('取消')}}</link-button>
                    </div>
                </bk-form-item>
            </bk-form>
        </bk-popover>
    </div>
</template>

<script>
    export default {
        props: {
            data: {
                type: Object,
                required: true
            },
            folders: {
                type: Array,
                required: true
            },
            batchSelectHandler: {
                type: Function,
                required: true
            }
        },
        data () {
            return {
                selected: ''
            }
        },
        methods: {
            handleClick () {
                this.$refs.popover.instance.show()
            },
            handleConfirm () {
                this.batchSelectHandler(this.selected)
                this.handleCancel()
            },
            handleCancel () {
                this.$refs.popover.instance.hide()
            },
            handleHidden () {
                this.selected = ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    .batch-trigger {
        line-height: 1;
        display: inline-block;
        vertical-align: baseline;
        [class^="icon"] {
            font-size: 14px;
            color: $primaryColor;
            cursor: pointer;
            &:hover {
                color: #1964E1;
            }
        }
    }
    .form-table-selector {
        width: 235px;
    }
    .selector-options {
        text-align: right;
        font-size: 0;
        .selector-link-button {
            font-size: 12px;
        }
    }
</style>
