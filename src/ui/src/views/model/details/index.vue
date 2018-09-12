<template>
    <div class="details-wrapper">
        <bk-tab :active-name.sync="curTabName">
            <bk-tabpanel name="host" :title="$t('ModelManagement[\'模型配置\']')">
                <v-field 
                    v-if="curTabName === 'host'"
                    :isEdit.sync="isEdit"
                    @createObject="createObject"
                    @confirm="updateModel"
                    @cancel="cancel"
                    ref="tab"
                ></v-field>
            </bk-tabpanel>
            <bk-tabpanel name="layout" :title="$t('ModelManagement[\'字段分组\']')" :show="isEdit">
                <v-layout v-if="curTabName === 'layout'"
                    @cancel="cancel"
                    ref="tab"
                ></v-layout>
            </bk-tabpanel>
            <bk-tabpanel name="other" :title="$t('ModelManagement[\'其他操作\']')" :show="isEdit">
                <v-other
                    v-if="curTabName === 'other'"
                    @closeSlider="updateModel"
                ></v-other>
            </bk-tabpanel>
        </bk-tab>
    </div>
</template>

<script>
    import vField from './field'
    import vLayout from './layout'
    import vOther from './other'
    export default {
        components: {
            vField,
            vLayout,
            vOther
        },
        props: {
            isEdit: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                curTabName: 'host'
            }
        },
        methods: {
            cancel () {
                this.$emit('cancel')
            },
            createObject () {
                this.$emit('createModel')
                this.$emit('update:isEdit', true)
            },
            updateModel () {
                this.$emit('updateModel')
            },
            isCloseConfirmShow () {
                if (this.curTabName === 'other') {
                    return false
                }
                return this.$refs.tab.isCloseConfirmShow()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .details-wrapper{
        padding: 8px 0 0;
        height: 100%;
    }
</style>
