<template>
    <bk-dropdown-menu trigger="click" :disabled="disabled" ref="dropdownMenu">
        <bk-button class="table-selector-trigger" type="default" slot="dropdown-trigger" :disabled="disabled">
            <i class="table-selector-checkbox" :class="type" @click.stop="handleCurrentClick"></i>
            <i class="bk-icon icon-angle-down"></i>
        </bk-button>
        <ul class="table-selector-list" slot="dropdown-content">
            <li class="table-selector-item" @click.stop="handleClick('current')">{{$t("Common['全选本页']")}}</li>
            <li class="table-selector-item" @click.stop="handleClick('all')">{{$t("Common['跨页全选']")}}</li>
        </ul>
    </bk-dropdown-menu>
</template>

<script>
    export default {
        name: 'cmdb-table-selector',
        props: {
            layout: Object
        },
        computed: {
            table () {
                return this.layout.table
            },
            total () {
                return this.table.pagination.count
            },
            size () {
                return this.table.pagination.size
            },
            selectedCount () {
                return this.table.checked.length
            },
            disabled () {
                return this.table.loading
            },
            checkboxColumn () {
                return this.layout.columns.find(column => column.type === 'checkbox')
            },
            type () {
                if (this.total && this.selectedCount) {
                    return this.selectedCount === this.total ? 'all' : 'part'
                }
                return null
            }
        },
        methods: {
            handleClick (type) {
                if (type === 'current' && this.table.hasCheckbox && !this.table.$scopedSlots[this.checkboxColumn.id]) {
                    const checked = this.table.list.map(item => item[this.checkboxColumn.id])
                    this.table.$emit('update:checked', checked)
                }
                this.table.$emit('handleCheckAll', type)
                this.$refs.dropdownMenu.hide()
            },
            handleCurrentClick () {
                if (this.disabled) return
                if (this.selectedCount) {
                    this.handleCancel()
                } else {
                    this.handleClick('current')
                }
            },
            handleCancel () {
                this.table.$emit('update:checked', [])
                this.table.$emit('handleCancelCheck')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .table-selector-trigger{
        position: relative;
        width: 100%;
        height: 38px;
        border: none;
        background-color: #fafbfd;
        .icon-angle-down{
            position: absolute;
            top: 11px;
            right: 11px;
        }
    }
    .table-selector-checkbox {
        display: inline-block;
        vertical-align: middle;
        width: 14px;
        height: 14px;
        margin: 0 5px 0 0;
        background: url('../../../assets/images/checkbox-sprite.png') no-repeat;
        background-position: 0 -95px;
        &.part {
            background-position: 0 -122px;
        }
        &.all {
            background-position: -33px -95px;
        }
    }
    .table-selector-list{
        width: 100%;
        font-size: 14px;
        line-height: 40px;
        .table-selector-item{
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