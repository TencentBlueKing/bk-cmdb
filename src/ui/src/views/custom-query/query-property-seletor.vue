<template>
    <div class="query-property-seletor">
        <bk-input v-model="searchProperty"
            :clearable="true"
            right-icon="bk-icon icon-search">
        </bk-input>
        <div class="property-box">
            <div class="group"
                v-for="group in groups"
                :key="group.id">
                <h3 class="group-title">{{group.name}}</h3>
                <div class="property-list">
                    <div class="property-item"
                        v-for="property in group.children"
                        :key="property.bk_property_id">
                        <bk-checkbox v-model="property.__selected__">{{property.bk_property_name}}</bk-checkbox>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        props: {
            properties: {
                type: Object,
                default: () => {}
            },
            selectedProperties: {
                type: Array,
                default: () => []
            }
        },
        data () {
            return {
                searchProperty: '',
                groups: [],
                originGroups: []
            }
        },
        computed: {
            selectedPropertyList () {
                const list = []
                this.groups.forEach(item => {
                    list.push(...item.children)
                })
                return list.filter(property => property.__selected__)
            },
            removePropertyList () {
                const selectedFilterIds = this.selectedPropertyList.map(property => property.filter_id)
                return this.selectedProperties.filter(filterId => !selectedFilterIds.includes(filterId))
            },
            addPropertyList () {
                return this.selectedPropertyList.filter(property => !this.selectedProperties.includes(property.filter_id))
            },
            hasChanged () {
                let res = false
                for (let i = 0; i < this.groups.length; i++) {
                    const originGroup = this.originGroups.find(item => item.id === this.groups[i].id) || {}
                    if (JSON.stringify(originGroup.children) !== JSON.stringify(this.groups[i].children)) {
                        res = true
                        break
                    }
                }
                return res
            }
        },
        watch: {
            properties () {
                this.init()
            }
        },
        created () {
            this.init()
        },
        methods: {
            init () {
                this.groups = this.getGroups()
                this.originGroups = this.$tools.clone(this.groups)
            },
            getGroups () {
                const properties = this.$tools.clone(this.properties)
                return Object.keys(properties).map(modelId => {
                    const model = this.$store.getters['objectModelClassify/getModelById'](modelId) || {}
                    return {
                        id: modelId,
                        name: model.bk_obj_name,
                        children: properties[modelId].map(property => {
                            if (this.selectedProperties.includes(`${property.bk_obj_id}-${property.bk_property_id}`)) {
                                property.__selected__ = true
                            }
                            return property
                        })
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .query-property-seletor {
        padding: 12px 20px;
        .property-box {
            padding: 15px 0 0 0;
            .group-title {
                font-size: 14px;
                color: #313237;
                padding: 0 0 20px 0;
            }
            .property-list {
                display: flex;
                flex-wrap: wrap;
                .property-item {
                    flex: 0 0 50%;
                    font-size: 14px;
                    color: #63656e;
                    padding: 0 0 20px 0;
                }
            }
        }
    }
</style>
