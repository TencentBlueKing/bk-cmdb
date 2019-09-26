<template>
    <bk-dialog :draggable="false" v-model="show">
        <h2 class="title" slot="header">{{$t('添加服务模板')}}</h2>
        <template>
            <cmdb-tips class="tips">{{$t('服务模板选择器提示语')}}</cmdb-tips>
            <ul class="template-list clearfix">
                <li class="template-item fl clearfix"
                    v-for="template in templates"
                    :class="{ 'is-selected': selected.includes(template.id) }"
                    :key="template.id">
                    <i class="bk-icon icon-check-circle-shape fr"></i>
                    <span class="template-name" :title="template.name">{{template.name}}</span>
                </li>
            </ul>
        </template>
    </bk-dialog>
</template>

<script>
    export default {
        data () {
            return {
                show: true,
                templates: [],
                selected: []
            }
        }
    }
</script>

<style lang="scss" scoped>
    .title {
        font-size: 24px;
        font-weight: normal;
    }
    .template-list {
        @include scrollbar-y;
        .template-item {
            width: 231px;
            height: 32px;
            line-height: 32px;
            border-radius: 2px;
            border: 1px solid #DCDEE5;
            color: #63656E;
            &.is-selected {
                background-color: #E1ECFF;
            }
        }
    }
</style>
