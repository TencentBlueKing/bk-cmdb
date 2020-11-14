<template>
    <section class="filter-wrapper" v-if="selected.length || showIPTag">
        <label class="filter-label">
            <i class="label-icon icon-cc-funnel"></i>
            <span class="label-text">{{$t('检索项')}}</span>
            <span class="label-colon">:</span>
        </label>
        <div class="filter-list">
            <filter-tag-ip v-if="showIPTag"></filter-tag-ip>
            <filter-tag-item class="filter-item"
                v-for="property in selected"
                :key="property.id"
                :property="property"
                v-bind="condition[property.id]">
            </filter-tag-item>
        </div>
    </section>
</template>

<script>
    import FilterTagIp from './filter-tag-ip'
    import FilterTagItem from './filter-tag-item'
    import FilterStore from './store'
    import Utils from './utils'
    export default {
        components: {
            FilterTagIp,
            FilterTagItem
        },
        computed: {
            condition () {
                return FilterStore.condition
            },
            showIPTag () {
                const list = Utils.splitIP(FilterStore.IP.text)
                return !!list.length
            },
            selected () {
                return FilterStore.selected.filter(property => {
                    const value = this.condition[property.id].value
                    return value !== null && !!value.toString().length
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .filter-wrapper {
        display: flex;
        margin: 10px 0 0 0;
        .filter-label {
            display: flex;
            font-size: 12px;
            align-items: center;
            align-self: flex-start;
            line-height: 22px;
            .label-icon {
                color: #979BA5;
            }
            .label-text {
                margin-left: 4px;
            }
            .label-colon {
                margin: 0 5px;
            }
        }
        .filter-list {
            display: flex;
            flex-wrap: wrap;
            flex: 1;
        }
    }
</style>
