<template>
    <bk-dialog
        v-model="show"
        :draggable="false"
        :mask-close="false"
        :width="730"
        header-position="left"
        title="添加字段"
        @value-change="handleVisibleChange"
        @confirm="handleConfirm"
        @cancel="handleCancel"
    >
        <bk-input v-if="propertyList.length"
            class="search"
            type="text"
            :placeholder="$t('请输入字段名称搜索')"
            clearable
            right-icon="bk-icon icon-search"
            v-model.trim="searchName"
            @input="hanldeFilterProperty">
        </bk-input>
        <bk-checkbox-group v-model="localChecked">
            <ul class="property-list">
                <li class="property-item" v-for="property in propertyList" :key="property.bk_property_id" v-show="property.__extra__.visible">
                    <bk-checkbox
                        :disabled="!property.host_apply_enabled"
                        :value="property.id"
                    >
                        <div
                            v-if="!property.host_apply_enabled"
                            v-bk-tooltips.top-start="'该字段不支持配置'"
                            style="outline:none"
                        >
                            {{property.bk_property_name}}
                        </div>
                        <div v-else>
                            {{property.bk_property_name}}
                        </div>
                    </bk-checkbox>
                </li>
            </ul>
        </bk-checkbox-group>
    </bk-dialog>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        props: {
            visible: {
                type: Boolean,
                default: false
            },
            checkedList: {
                type: Array,
                default: () => ([])
            }
        },
        data () {
            return {
                show: this.visible,
                localChecked: [],
                searchName: '',
                propertyList: []
            }
        },
        computed: {
            ...mapGetters('hostApply', ['configPropertyList'])
        },
        watch: {
            visible (val) {
                this.show = val
            },
            checkedList: {
                handler () {
                    this.localChecked = this.checkedList
                },
                immediate: true
            }
        },
        async created () {
            await this.getHostPropertyList()
            this.propertyList = this.configPropertyList.filter(property => property.host_apply_enabled)
        },
        methods: {
            async getHostPropertyList () {
                try {
                    const data = await this.$store.dispatch('hostApply/getProperties', {
                        requestId: 'getHostPropertyList',
                        fromCache: true
                    })
                    this.$store.commit('hostApply/setPropertyList', data)
                } catch (e) {
                    console.error(e)
                }
            },
            handleVisibleChange (val) {
                this.$emit('update:visible', val)
            },
            handleConfirm () {
                this.$emit('update:checkedList', this.localChecked)
            },
            handleCancel () {
                this.localChecked = this.checkedList
            },
            hanldeFilterProperty () {
                // 使用visible方式是为了兼容checkbox-group组件
                this.propertyList.forEach(property => {
                    property.__extra__.visible = property.bk_property_name.indexOf(this.searchName) > -1
                })
                this.propertyList = [...this.propertyList]
            }
        }
    }
</script>

<style lang="scss" scoped>
    .search {
        width: 240px;
        margin-bottom: 10px;
    }
    .property-list {
        display: flex;
        flex-wrap: wrap;
        align-content: flex-start;
        height: 264px;
        @include scrollbar-y;

        .property-item {
            flex: 0 0 33.3333%;
            margin: 8px 0;
        }
    }
</style>
