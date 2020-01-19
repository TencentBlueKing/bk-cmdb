<template>
    <cmdb-sticky-layout>
        <bk-form class="form-layout" form-type="inline" :label-width="300">
            <bk-form-item class="form-item" :label="$t('任务名称')" required>
                <bk-input class="form-meta"
                    :placeholder="$t('请输入xx', { name: $t('任务名称') })"
                    :data-vv-as="$t('任务名称')"
                    data-vv-name="mission_name"
                    v-validate="'required|length:256'"
                    v-model.trim="form.mission_name">
                </bk-input>
                <p class="form-error" v-if="errors.has('mission_name')">{{errors.first('mission_name')}}</p>
            </bk-form-item>
            <bk-form-item class="form-item" :label="$t('账户名称')" required>
                <bk-select class="form-meta"
                    :searchable="true"
                    :readonly="!isCreateMode"
                    :placeholder="$t('请选择xx', { name: $t('账户名称') })"
                    :data-vv-as="$t('账户名称')"
                    data-vv-name="account_name"
                    v-model="form.account_name"
                    v-validate="'required'">
                    <bk-option v-for="account in accountList"
                        :key="account.id"
                        :name="account.name"
                        :id="account.id">
                    </bk-option>
                </bk-select>
                <link-button class="form-account-link" @click="handleLinkToCloudAccount">
                    <i class="icon-cc-share"></i>
                    {{$t('跳转账户管理')}}
                </link-button>
                <p class="form-tips" v-if="form.account_name">
                    <i class="icon-cc-exclamation-tips"></i>
                    {{$t('同步频率提示语')}}
                </p>
                <p class="form-error" v-if="errors.has('account_name')">{{errors.first('account_name')}}</p>
            </bk-form-item>
            <bk-form-item class="form-item" :label="$t('资源类型')" required>
                <bk-select class="form-meta"
                    :readonly="true"
                    :placeholder="$t('请选择xx', { name: $t('资源类型') })"
                    :data-vv-as="$t('资源类型')"
                    data-vv-name="resource_type"
                    v-model="form.resource_type"
                    v-validate="'required'">
                    <bk-option id="host" :name="$t('主机')"></bk-option>
                </bk-select>
                <p class="form-error" v-if="errors.has('resource_type')">{{errors.first('resource_type')}}</p>
            </bk-form-item>
        </bk-form>
        <bk-form class="form-layout" form-type="inline" :label-width="300" v-if="form.account_name">
            <bk-form-item class="form-item" :label="$t('云区域设定')" required>
                <bk-radio-group class="form-meta" v-model="form.setting">
                    <bk-radio :value="settingComponents.custom">{{$t('自定义')}}</bk-radio>
                    <bk-radio style="margin-left: 85px;" :value="settingComponents.all">{{$t('全部及后续新增')}}</bk-radio>
                </bk-radio-group>
            </bk-form-item>
        </bk-form>
        <div class="form-setting-component">
            <component :is="form.setting"></component>
        </div>
        <div class="form-options"
            slot="footer"
            slot-scope="{ sticky }"
            :class="{ 'is-sticky': sticky }">
            <bk-button theme="primary" @click="handleSumbit">{{$t('提交')}}</bk-button>
            <bk-button class="ml10" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
    </cmdb-sticky-layout>
</template>

<script>
    import { MENU_RESOURCE_CLOUD_ACCOUNT } from '@/dictionary/menu-symbol'
    import ResourceFormCustom from './resource-form-custom.vue'
    import ResourceFormAll from './resource-form-all.vue'
    import CloudResourceDetailsInfo from './resource-details-info.vue'
    export default {
        name: 'cloud-resource-form',
        components: {
            [ResourceFormCustom.name]: ResourceFormCustom,
            [ResourceFormAll.name]: ResourceFormAll
        },
        props: {
            mission: {
                type: Object,
                default: null
            },
            mode: {
                type: String,
                default: 'create'
            },
            container: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                accountList: [{ id: 1, name: '阿里云账户' }],
                form: {
                    mission_name: '',
                    account_name: '',
                    resource_type: 'host',
                    setting: null
                },
                settingComponents: {
                    custom: ResourceFormCustom.name,
                    all: ResourceFormAll.name
                }
            }
        },
        computed: {
            isCreateMode () {
                return this.mode === 'create'
            }
        },
        watch: {
            'form.account_name' (value) {
                if (value) {
                    this.form.setting = this.settingComponents.custom
                } else {
                    this.form.setting = null
                }
            }
        },
        methods: {
            handleLinkToCloudAccount () {
                this.$router.push({
                    name: MENU_RESOURCE_CLOUD_ACCOUNT
                })
            },
            handleAccountSelected (value) {
                this.form.setting = this.settingComponents.custom
            },
            handleSumbit () {
                if (this.isCreateMode) {
                    this.doCreate()
                } else {
                    this.doUpdate()
                }
            },
            async doCreate () {
                try {
                    await Promise.resolve()
                    this.container.hide('request-refresh')
                } catch (e) {
                    console.error(e)
                }
            },
            async doUpdate () {
                try {
                    await Promise.resolve()
                    this.container.show({
                        detailsComponent: CloudResourceDetailsInfo.name,
                        props: {
                            mission: this.mission
                        }
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            handleCancel () {
                if (this.isCreateMode) {
                    this.container.hide()
                } else {
                    this.container.show({
                        detailsComponent: CloudResourceDetailsInfo.name,
                        props: {
                            mission: this.mission
                        }
                    })
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-layout {
        padding: 10px 25px;
        display: block;
        .form-item:nth-child(2n) {
            margin-left: 32px;
        }
        .form-item:nth-child(n + 3) {
            margin-top: 10px;
        }
        .form-item {
            width: 300px;
            margin-left: 0;
            /deep/ {
                .bk-label,
                .bk-form-content {
                    display: block;
                    float: none;
                    text-align: left;
                }
            }
            .form-meta {
                display: block;
            }
            .form-error {
                position: absolute;
                top: 100%;
                left: 0;
                line-height: 1.5;
                font-size: 12px;
                color: $dangerColor;
            }
        }
        .form-account-link {
            position: absolute;
            bottom: 100%;
            right: 0;
            font-size: 12px;
        }
        .form-tips {
            position: absolute;
            top: 100%;
            left: 0;
            font-size: 12px;
            line-height: 1.5;
            color: #979BA5;
            [class^=icon] {
                vertical-align: 1px;
            }
        }
    }
    .form-options {
        font-size: 0;
        margin-top: 10px;
        padding: 0 25px;
        &.is-sticky {
            margin-top: 0;
            padding: 15px 25px;
            border-top: 1px solid $borderColor;
            background-color: #FAFBFD;
        }
    }
    .form-setting-component {
        margin: 10px 25px 20px;
    }
</style>
