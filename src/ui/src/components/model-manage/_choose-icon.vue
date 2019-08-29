<template>
    <div class="choose-icon">
        <bk-input type="text" class="system-icon-search" v-show="activeTab === 'system'"
            right-icon="bk-icon icon-search"
            :placeholder="$t('请输入关键词')"
            v-model.trim="searchText">
        </bk-input>
        <bk-tab :active.sync="activeTab" type="unborder-card" class="icon-tab">
            <bk-tab-panel name="system" :label="$t('系统图标')">
                <icon-set></icon-set>
            </bk-tab-panel>
            <bk-tab-panel name="custom" :label="$t('自定义图标')"></bk-tab-panel>
        </bk-tab>
        <div class="footer">
            <bk-button theme="primary">{{$t('确定')}}</bk-button>
            <bk-button>{{$t('取消')}}</bk-button>
        </div>
    </div>
</template>

<script>
    import iconSet from './icon-set'
    export default {
        components: {
            iconSet
        },
        data () {
            return {
                activeTab: 'system',
                searchText: ''
            }
        },
        methods: {
        }
    }
</script>

<style lang="scss" scoped>
    .choose-icon {
        position: relative;
        height: 460px;
        overflow: hidden;
        .system-icon-search {
            position: absolute;
            top: 12px;
            right: 20px;
            width: 240px;
            z-index: 2;
        }
        .icon-tab {
            width: 100%;
            height: calc(100% - 58px);
            /deep/ .bk-tab-section {
                height: calc(100% - 57px);
                @include scrollbar-y;
            }
            /deep/ .bk-tab-header {
                padding: 0;
                margin: 0 20px;
            }
        }
        .footer {
            height: 57px;
            line-height: 56px;
            text-align: right;
            font-size: 0;
            padding-right: 24px;
            background-color: #fafbfd;
            border-top: 1px solid #dcdee5;
            .bk-button {
                margin-left: 10px;
            }
        }
    }
</style>
