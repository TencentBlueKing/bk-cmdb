<template>
    <bk-dialog
        v-model="show"
        :draggable="false"
        :mask-close="false"
        :width="730"
        header-position="left"
        title="添加字段"
        @value-change="handleVisiableChange"
        @confirm="handleConfirm"
        @cancel="handleCancel"
    >
        <bk-checkbox-group v-model="localChecked">
            <ul class="property-list">
                <li class="property-item" v-for="property in configPropertyList" :key="property.bk_property_id">
                    <bk-checkbox
                        :disabled="property.__extra__.disabled"
                        :value="property.id"
                    >
                        <div
                            v-if="property.__extra__.disabled"
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
    import { mapGetters, mapActions } from 'vuex'
    export default {
        props: {
            visiable: {
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
                show: this.visiable,
                localChecked: []
            }
        },
        computed: {
            ...mapGetters('hosts', ['configPropertyList'])
        },
        watch: {
            visiable (val) {
                this.show = val
            },
            checkedList: {
                handler () {
                    this.localChecked = this.checkedList
                },
                immediate: true
            }
        },
        created () {
            this.getHostPropertyList()
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'searchObjectAttribute'
            ]),
            async getHostPropertyList () {
                try {
                    const data = await this.searchObjectAttribute({
                        params: this.$injectMetadata({
                            bk_obj_id: 'host',
                            bk_supplier_account: this.supplierAccount
                        }),
                        config: {
                            requestId: 'getHostPropertyList',
                            fromCache: true
                        }
                    })

                    this.$store.commit('hosts/setPropertyList', data)
                } catch (e) {
                    console.error(e)
                }
            },
            handleVisiableChange (val) {
                this.$emit('update:visiable', val)
            },
            handleConfirm () {
                this.$emit('update:checkedList', this.localChecked)
            },
            handleCancel () {
                this.localChecked = this.checkedList
            }
        }
    }
</script>

<style lang="scss" scoped>
    .property-list {
        display: flex;
        flex-wrap: wrap;
        max-height: 360px;
        @include scrollbar-y;

        .property-item {
            flex: 0 0 33.3333%;
            margin: 8px 0;
        }
    }
</style>
