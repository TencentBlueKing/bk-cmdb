<template>
    <div class="field-wrapper">
       <v-base-info></v-base-info>
       <div class="field-content clearfix">
           <div class="create-field">
               <bk-button type="primary" :title="$t('ModelManagement[\'新增字段\']')" @click="createField">
                   {{$t('ModelManagement["新增字段"]')}}
               </bk-button>
               <div class="btn-group">
                   <bk-button type="default" :title="$t('ModelManagement[\'导入\']')" :disabled="isReadOnly" class="btn mr10">
                        {{$t('ModelManagement["导入"]')}}
                        <input ref="fileInput" type="file" @change.prevent="">
                    </bk-button>
                    <form :action="exportUrl" method="POST" class="form">
                        <bk-button type="default" btnType="submit" :title="$t('ModelManagement[\'导出\']')" class="btn">
                            {{$t('ModelManagement["导出"]')}}
                        </bk-button>
                    </form>
               </div>
           </div>
           <div class="field-table">
                <ul class="table-title">
                    <li>{{$t('ModelManagement["唯一"]')}}</li>
                    <li>{{$t('ModelManagement["必填字段"]')}}</li>
                    <li>{{$t('ModelManagement["类型"]')}}</li>
                    <li>{{$t('ModelManagement["字段名"]')}}</li>
                    <li>{{$t('ModelManagement["操作"]')}}</li>
                </ul>
                <ul class="table-content">
                    <li>
                        <div></div>
                        <div>
                            <v-model-field></v-model-field>
                        </div>
                    </li>
                </ul>
           </div>
       </div>
    </div>
</template>

<script>
    import vBaseInfo from './base-info'
    import vModelField from '@/components/model-field'
    export default {
        components: {
            vBaseInfo,
            vModelField
        },
        methods: {
            createField () {
                
            }
        }
    }
</script>

<style lang="scss" scoped>
    .field-wrapper {
        padding: 0 30px;
        .field-content {
            padding-top: 20px;
            border-top: 1px solid $cmdbBorderLightColor;
        }
        .create-field {
            .btn-group {
                float: right;
                font-size: 0;
                .btn {
                    position: relative;
                    overflow: hidden;
                    input[type="file"] {
                        position: absolute;
                        left: -70px;
                        top: 0;
                        opacity: 0;
                        width: calc(100% + 70px);
                        height: 100%;
                        cursor: pointer;
                    }
                }
                .form {
                    display: inline-block;
                }
            }
        }
        .field-table {
            margin-top: 10px;
            .table-title {
                width: 100%;
                height: 40px;
                >li {
                    float: left;
                    height: 40px;
                    line-height: 40px;
                    background: $cmdbPrimaryColor;
                    text-align: center;
                    border: 1px solid $cmdbBorderLightColor;
                    border-right: none;
                    &:nth-child(1) {
                        width:80px;
                    }
                    &:nth-child(2) {
                        width:80px;
                    }
                    &:nth-child(3) {
                        width:145px;
                    }
                    &:nth-child(4) {
                        width:290px;
                    }
                    &:nth-child(5) {
                        width: 105px;
                        border-right: 1px solid $cmdbBorderLightColor;
                    }
                }
            }
        }
    }
</style>
