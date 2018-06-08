<template>
    <div class="association-list">
        <div class="association-options">
            <button class="bk-button bk-primary association-btn btn-new">{{$t('Association["新增关联"]')}}</button>
            <button class="bk-button association-btn btn-remove">{{$t('Association["解除选中关联"]')}}</button>
        </div>
        <div class="association-table-layout">
            <v-association class="association-table"
                v-for="(association, index) in associations"
                :key="index"
                :list="association.list"
                :property="association.property"
                @handleCheck="handleCheck">
            </v-association>
        </div>
    </div>
</template>

<script>
    import vAssociation from '@/components/common/association/association'
    import { mapGetters } from 'vuex'
    export default {
        components: {
            vAssociation
        },
        props: {
            objId: {
                type: String,
                required: true
            },
            data: {
                type: Object,
                defaut () {
                    return {}
                }
            }
        },
        data () {
            return {
                checked: {}
            }
        },
        computed: {
            ...mapGetters('object', ['attribute']),
            objAttribute () {
                return this.attribute[this.objId] || []
            },
            associations () {
                const associationAttributes = this.objAttribute.filter(property => ['singleasst', 'multiasst'].includes(property['bk_property_type']))
                const associations = associationAttributes.map(property => {
                    const {bk_property_id: bkPropertyId} = property
                    return {
                        property: property,
                        list: Array.isArray(this.data[bkPropertyId]) ? this.data[bkPropertyId].filter(({id}) => id !== '') : []
                    }
                })
                return associations
            }
        },
        watch: {
            objId (objId) {
                this.$store.dispatch('object/getAttribute', {objId})
            },
            associations (associations) {
                associations.forEach(association => {
                    if (!this.checked.hasOwnProperty(association.property['bk_property_id'])) {
                        this.$set(this.checked, association.property['bk_property_id'], [])
                    }
                })
            }
        },
        created () {
            this.$store.dispatch('object/getAttribute', {objId: this.objId})
        },
        methods: {
            handleCheck (checked, property) {
                this.checked[property['bk_property_id']] = checked
            }
        }
    }
</script>

<style lang="scss" scoped>
    .association-list{
        position: relative;
        padding: 0 30px;
        margin: 30px 0;
        .association-options{
            position: absolute;
            bottom: 100%;
            right: 30px;
            font-size: 0;
            .association-btn{
                height: 24px;
                line-height: 22px;
                font-size: 12px;
                &.btn-new{
                    color: #fff;
                    border: none;
                }
                &.btn-remove {
                    margin: 0 0 0 9px;
                }
            }
        }
    }
    .association-table{
        margin: 10px 0 0 0;
    }
</style>