<template>
    <div class="edit-label-list">
        <div class="scrollbar-box">
            <div class="label-item" v-for="(label, index) in list" :key="index">
                <div class="label-key">
                    <input class="cmdb-form-input" type="text"
                        :name="'label' + label.key"
                        v-validate="'required'"
                        v-model="label.key"
                        :placeholder="$t('BusinessTopology[\'添加标签键\']')">
                </div>
                <input class="cmdb-form-input label-value" type="text"
                    :name="'label' + label.value"
                    v-validate="'required'"
                    v-model="label.value"
                    :placeholder="$t('BusinessTopology[\'标签值\']')">
                <i class="bk-icon icon-plus-circle-shape icon-btn"
                    v-show="list.length - 1 === index"
                    @click="handleAddLabel(index)">
                </i>
                <i :class="['bk-icon', 'icon-minus-circle-shape', 'icon-btn', { 'disabled': list.length === 1 }]"
                    @click="handleRemoveLabel(index)">
                </i>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        props: {
            defaultList: {
                type: Array,
                default: () => []
            }
        },
        data () {
            return {
                list: this.defaultList,
                removeList: []
            }
        },
        created () {
            this.list = this.list.concat({
                id: -1,
                key: '',
                value: ''
            })
        },
        methods: {
            handleAddLabel (index) {
                const currentTag = this.list[index]
                if (!currentTag.key || !currentTag.value) {
                    this.$bkMessage({
                        message: this.$t("BusinessTopology['请填写完整标签键/值']"),
                        theme: 'warning'
                    })
                    return
                }
                this.list.push({
                    id: -1,
                    key: '',
                    value: ''
                })
            },
            handleRemoveLabel (index) {
                if (index === this.list.length - 1) return
                const currentTag = this.list[index]
                if (currentTag.id !== -1) {
                    this.removeList.push(currentTag)
                }
                this.list.splice(index, 1)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .edit-label-list {
        padding: 8px 12px 18px 26px;
        .scrollbar-box {
            @include scrollbar-y;
            height: 315px;
        }
        .label-item {
            display: flex;
            align-items: center;
            margin-bottom: 20px;
            &:last-child {
                margin-bottom: 0;
            }
        }
        .label-key {
            width: 172px;
            margin-right: 10px;
        }
        .label-value {
            width: 292px;
            margin-right: 10px;
        }
        .icon-btn {
            color: #c4c6cc;
            font-size: 18px;
            margin-right: 8px;
            cursor: pointer;
            &.disabled {
                color: #dcdee5;
                cursor: not-allowed;
                &:hover {
                    color: #dcdee5 !important;
                }
            }
            &:hover {
                color: #979ba5;
            }
        }
    }
</style>
