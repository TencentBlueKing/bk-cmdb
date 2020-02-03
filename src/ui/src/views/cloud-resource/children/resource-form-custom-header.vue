<template>
    <div class="custom-header">
        {{data.column.label}}
        <bk-popover class="batch-trigger"
            trigger="manual" theme="light" placement="bottom" ref="popover"
            :tippy-options="{
                hideOnClick: false,
                onHidden: handleHidden
            }">
            <i class="icon-cc-batch-update"
                v-bk-tooltips="{
                    content: $t('请选择xx', { name: 'VPC' }),
                    disabled: !disabled
                }"
                :class="{ 'is-disabled': disabled }"
                @click="handleClick">
            </i>
            <bk-form form-type="vertical" slot="content">
                <bk-form-item :label="$t('批量编辑')">
                    <cloud-resource-folder-selector class="folder-selector"
                        v-model="selected">
                    </cloud-resource-folder-selector>
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
    import CloudResourceFolderSelector from './resource-folder-selector.vue'
    export default {
        name: 'cloud-resource-custom-header',
        components: {
            CloudResourceFolderSelector
        },
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
            },
            disabled: Boolean
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
            &.is-disabled {
                color: #c4c6cc;
            }
        }
    }
    .folder-selector {
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
