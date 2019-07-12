<template>
    <div class="edit-label-list">
        <div class="scrollbar-box">
            <div class="label-item" v-for="(label, index) in list" :key="index">
                <div class="label-key" :class="{ 'is-error': errors.has('key-' + index) }">
                    <input class="cmdb-form-input" type="text"
                        :data-vv-name="'key-' + index"
                        v-validate="'required|instanceTag'"
                        v-model="label.key"
                        :placeholder="$t('BusinessTopology[\'添加标签键\']')">
                    <p class="input-error">{{errors.first('key-' + index)}}</p>
                </div>
                <div class="label-value" :class="{ 'is-error': errors.has('value-' + index) }">
                    <input class="cmdb-form-input" type="text"
                        :data-vv-name="'value-' + index"
                        v-validate="'required|instanceTag'"
                        v-model="label.value"
                        :placeholder="$t('BusinessTopology[\'标签值\']')">
                    <p class="input-error">{{errors.first('value-' + index)}}</p>
                </div>
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
                originList: [],
                list: this.defaultList,
                removeKeysList: []
            }
        },
        created () {
            this.initValue()
        },
        methods: {
            initValue () {
                this.originList = this.$tools.clone(this.defaultList)
                if (!this.list.length) {
                    this.list = this.list.concat([{
                        id: -1,
                        key: '',
                        value: ''
                    }])
                } else {
                    this.list = this.list.concat([])
                }
            },
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
                if (this.list.length === 1) return
                const currentTag = this.list[index]
                if (currentTag.id !== -1) {
                    this.removeKeysList.push(currentTag.key)
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
                margin-bottom: 10px;
            }
        }
        .label-key {
            position: relative;
            width: 172px;
            margin-right: 10px;
        }
        .label-value {
            position: relative;
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
        .input-error {
            position: absolute;
            top: 100%;
            left: 0;
            line-height: 14px;
            font-size: 12px;
            color: #ff5656;
        }
        .is-error input.cmdb-form-input {
            border-color: #ff5656;
        }
    }
</style>
