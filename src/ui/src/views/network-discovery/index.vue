<template>
    <div class="network-wrapper">
        <div class="filter-wrapper">
            <bk-button type="primary">
                {{$t('NetworkDiscovery["配置网络发现"]')}}
            </bk-button>
            <div class="filter-content fr">
                <input type="text" class="cmdb-form-input fl" :placeholder="$t('NetworkDiscovery[\'请输入云区域名称\']')">
                <bk-button type="default" class="fl" v-tooltip="$t('NetworkDiscovery[\'查看完成历史\']')">
                    <i class="icon-cc-history"></i>
                </bk-button>
            </div>
        </div>
        <cmdb-table></cmdb-table>
    </div>
</template>

<script>
    export default {
        
    }
</script>


<style lang="scss" scoped>
    .network-wrapper {
        .filter-wrapper {
            .filter-content {
                .cmdb-form-input {
                    width: 260px;
                    margin-right: 10px;
                }
            }
        }
    }
</style>
