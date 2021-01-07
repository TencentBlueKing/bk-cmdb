<template>
    <bk-select class="service-category-selector"
        v-model="selected"
        searchable
        :multiple="multiple"
        :clearable="allowClear"
        :disabled="disabled"
        :placeholder="placeholder"
        :font-size="fontSize"
        :popover-options="{
            boundary: 'window'
        }"
        ref="selector">
        <bk-option-group
            v-for="(group, groupIndex) in firstClassList"
            :name="group.name"
            :key="groupIndex">
            <bk-option v-for="(option, optionIndex) in group.secondCategory"
                :key="optionIndex"
                :id="option.id"
                :name="option.name">
            </bk-option>
        </bk-option-group>
    </bk-select>
</template>

<script>
    import { mapState, mapGetters } from 'vuex'
    export default {
        name: 'cmdb-service-category',
        props: {
            value: {
                type: [Array, String],
                default: () => ([])
            },
            disabled: {
                type: Boolean,
                default: false
            },
            multiple: {
                type: Boolean,
                default: true
            },
            allowClear: {
                type: Boolean,
                default: false
            },
            autoSelect: {
                type: Boolean,
                default: true
            },
            placeholder: {
                type: String,
                default: ''
            },
            fontSize: {
                type: [String, Number],
                default: 'medium'
            }
        },
        data () {
            return {
                selected: this.value || [],
                firstClassList: []
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapState('businessHost', [
                'categoryMap'
            ])
        },
        watch: {
            value (value) {
                this.selected = value || []
            },
            selected (selected) {
                this.$emit('input', selected)
                this.$emit('on-selected', selected)
            }
        },
        created () {
            this.getServiceCategories()
        },
        methods: {
            async getServiceCategories () {
                if (this.categoryMap.hasOwnProperty(this.bizId)) {
                    this.firstClassList = this.categoryMap[this.bizId]
                } else {
                    try {
                        const data = await this.$store.dispatch('serviceClassification/searchServiceCategory', {
                            params: { bk_biz_id: this.bizId }
                        })
                        const categories = this.collectServiceCategories(data.info)
                        this.firstClassList = categories
                        this.$store.commit('businessHost/setCategories', {
                            id: this.bizId,
                            categories: categories
                        })
                    } catch (e) {
                        console.error(e)
                        this.firstClassList = []
                    }
                }
            },
            collectServiceCategories (data) {
                const categories = []
                data.forEach(item => {
                    if (!item.category.bk_parent_id) {
                        categories.push(item.category)
                    }
                })
                categories.forEach(category => {
                    category.secondCategory = data.filter(item => item.category.bk_parent_id === category.id).map(item => item.category)
                })
                return categories
            }
        }
    }
</script>

<style lang="scss" scoped>
    .service-category-selector {
        width: 100%;
    }
</style>
