<template>
    <bk-select
        v-bind="$attrs"
        v-model="localValue"
        ext-popover-cls="audit-model-selector">
        <bk-option-group
            v-for="(group, index) in displayModelList"
            :name="group.bk_classification_name"
            :key="index">
            <bk-option v-for="option in group.bk_objects"
                :key="option.bk_obj_id"
                :id="option.bk_obj_id"
                :name="option.bk_obj_name">
            </bk-option>
        </bk-option-group>
    </bk-select>
</template>

<script>
    export default {
        props: {
            value: {
                type: String,
                default: ''
            }
        },
        data () {
            return {
                classifications: []
            }
        },
        computed: {
            localValue: {
                get () {
                    return this.value
                },
                set (values) {
                    this.$emit('input', values)
                    this.$emit('change', values)
                }
            },
            displayModelList () {
                const displayModelList = []
                this.classifications.forEach(classification => {
                    displayModelList.push({
                        ...classification,
                        'bk_objects': classification['bk_objects'].filter(model => {
                            return !model.bk_ispaused && !model.bk_ishidden
                        })
                    })
                })
                return displayModelList.filter(item => item.bk_objects.length > 0)
            }
        },
        created () {
            this.getModelList()
        },
        methods: {
            async getModelList () {
                try {
                    this.classifications = await this.$store.dispatch('objectModelClassify/searchClassificationsObjects', {
                        fromCache: true
                    })
                } catch (error) {
                    this.classifications = []
                }
            }
        }
    }
</script>

<style lang="scss">
    .bk-select-dropdown-content{
        &.audit-model-selector {
            .bk-options {
                overflow-y: auto;
                overflow-x: hidden;

                .bk-option-group-name {
                    overflow: hidden;
                    text-overflow: ellipsis;
                    white-space: nowrap;
                }
            }
        }
    }
</style>
