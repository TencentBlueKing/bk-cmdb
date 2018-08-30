<template>
    <div class="details-wrapper">
        <bk-tab :active-name.sync="curTabName">
            <bk-tabpanel name="host" :title="$t('ModelManagement[\'模型配置\']')">
                <v-field :isEdit.sync="isEdit"
                    @confirm="updateModel"
                    @cancel="cancel"
                ></v-field>
            </bk-tabpanel>
            <bk-tabpanel name="layout" :title="$t('ModelManagement[\'字段分组\']')" :show="isEdit">
                <v-layout></v-layout>
            </bk-tabpanel>
            <bk-tabpanel name="other" :title="$t('ModelManagement[\'其他操作\']')" :show="isEdit">
                <v-other
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
                curTabName: ''
            }
        },
        methods: {
            cancel () {
                this.$emit('cancel')
            },
            updateModel () {
                this.$emit('updateModel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .details-wrapper{
        padding: 8px 20px 0;
        height: 100%;
    }
</style>
