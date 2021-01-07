<template>
    <div class="create-layout" v-bkloading="{ isLoading: $loading([request.create, request.update]) }">
        <bk-form form-type="vertical">
            <bk-form-item class="create-form-item" :label="$t('账户名称')" required>
                <bk-input class="create-form-meta"
                    :placeholder="$t('请输入xx', { name: $t('账户名称') })"
                    :data-vv-as="$t('账户名称')"
                    data-vv-name="bk_account_name"
                    v-validate="'required|length:256'"
                    v-model.trim="form.bk_account_name">
                </bk-input>
                <p class="create-form-error" v-if="errors.has('bk_account_name')">{{errors.first('bk_account_name')}}</p>
            </bk-form-item>
            <bk-form-item class="create-form-item" :label="$t('账户类型')" required>
                <bk-select class="create-form-meta"
                    :placeholder="$t('请选择xx', { name: $t('账户类型') })"
                    :readonly="!isCreateMode"
                    :data-vv-as="$t('账户类型')"
                    :loading="$loading(request.property)"
                    data-vv-name="bk_cloud_vendor"
                    v-model="form.bk_cloud_vendor"
                    v-validate="'required'">
                    <bk-option v-for="vendor in vendors"
                        :key="vendor.id"
                        :name="vendor.name"
                        :id="vendor.id">
                        <cmdb-vendor :type="vendor.id"></cmdb-vendor>
                    </bk-option>
                    <cmdb-vendor class="vendor-selector-trigger" slot="trigger" :type="form.bk_cloud_vendor" empty-text=""></cmdb-vendor>
                </bk-select>
                <p class="create-form-error" v-if="errors.has('bk_cloud_vendor')">
                    {{errors.first('bk_cloud_vendor')}}
                </p>
            </bk-form-item>
            <bk-form-item class="create-form-item" label="ID" required>
                <bk-input class="create-form-meta"
                    :placeholder="$t('请输入xx', { name: 'ID' })"
                    data-vv-as="ID"
                    data-vv-name="bk_secret_id"
                    v-model.trim="form.bk_secret_id"
                    v-validate="'required|length:256'">
                </bk-input>
                <bk-link theme="primary" class="create-form-link"
                    v-show="form.bk_cloud_vendor"
                    @click="handleLinkToFAQ">
                    {{$t('如何获取ID和Key')}}
                </bk-link>
                <p class="create-form-error" v-if="errors.has('bk_secret_id')">
                    {{errors.first('bk_secret_id')}}
                </p>
            </bk-form-item>
            <bk-form-item class="create-form-item clearfix" label="Key" required>
                <bk-button class="create-form-button fr"
                    :loading="$loading(request.verify)"
                    @click="handleTest">
                    {{$t('连通测试')}}
                </bk-button>
                <bk-input class="create-form-meta key"
                    :placeholder="$t('请输入xx', { name: 'Key' })"
                    data-vv-as="Key"
                    data-vv-name="bk_secret_key"
                    v-model.trim="form.bk_secret_key"
                    v-validate="'required|length:256'"
                    @focus="handleSecretKeyFocus"
                    @blur="handleSecretKeyBlur">
                </bk-input>
                <template v-if="verifyResult.hasOwnProperty('done')">
                    <p class="create-verify-success" v-if="verifyResult.connected">
                        <i class="bk-icon icon-check-circle-shape"></i>
                        {{$t('账户连通成功')}}
                    </p>
                    <p class="create-verify-fail" v-else>
                        <i class="bk-icon icon-close-circle-shape"></i>
                        {{verifyResult.msg}}
                    </p>
                </template>
                <p class="create-form-error" v-else-if="errors.has('bk_secret_key')">
                    {{errors.first('bk_secret_key')}}
                </p>
            </bk-form-item>
            <bk-form-item class="create-form-item" :label="$t('备注')">
                <bk-input class="create-form-meta" type="textarea"
                    :placeholder="$t('请输入xx', { name: $t('备注') })"
                    :data-vv-as="$t('备注')"
                    data-vv-name="bk_description"
                    v-model.trim="form.bk_description"
                    v-validate="'length:2000'">
                </bk-input>
                <p class="create-form-error" v-if="errors.has('bk_description')">{{errors.first('bk_description')}}</p>
            </bk-form-item>
            <div class="create-form-options">
                <cmdb-auth class="mr10" :auth="auth">
                    <bk-button theme="primary" slot-scope="{ disabled }"
                        :disabled="disabled || !hasChange"
                        @click.stop.prevent="handleSubmit">
                        {{isCreateMode ? $t('提交') : $t('保存')}}
                    </bk-button>
                </cmdb-auth>
                <bk-button theme="default" @click.stop.prevent="handleCancel(null)">{{$t('取消')}}</bk-button>
            </div>
        </bk-form>
    </div>
</template>

<script>
    import CmdbVendor from '@/components/ui/other/vendor'
    import { mapGetters } from 'vuex'
    import { CLOUD_AREA_PROPERTIES } from '@/dictionary/request-symbol'
    import RouterQuery from '@/router/query'
    const DEFAULT_FORM = {
        bk_account_name: '',
        bk_cloud_vendor: '',
        bk_secret_id: '',
        bk_secret_key: '',
        bk_description: ''
    }
    export default {
        name: 'cloud-account-form',
        components: {
            CmdbVendor
        },
        props: {
            container: {
                type: Object,
                default: () => ({})
            },
            mode: {
                type: String,
                default: 'create'
            },
            account: {
                type: Object,
                default: () => ({})
            }
        },
        data () {
            return {
                form: {
                    ...DEFAULT_FORM
                },
                fakeSecret: '******',
                verifyResult: {},
                properties: [],
                request: {
                    create: Symbol('create'),
                    update: Symbol('update'),
                    verify: Symbol('verify'),
                    property: CLOUD_AREA_PROPERTIES
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            isCreateMode () {
                return this.mode === 'create'
            },
            verifyParams () {
                return {
                    bk_cloud_vendor: this.form.bk_cloud_vendor,
                    bk_secret_id: this.form.bk_secret_id,
                    bk_secret_key: this.form.bk_secret_key
                }
            },
            verifyRequired () {
                const changed = ['bk_cloud_vendor', 'bk_secret_id', 'bk_secret_key'].some(key => this.form[key] !== this.verifyResult[key])
                return changed || !this.verifyResult.connected
            },
            vendors () {
                const vendorProperty = this.properties.find(property => property.bk_property_id === 'bk_cloud_vendor')
                if (vendorProperty) {
                    return vendorProperty.option || []
                }
                return []
            },
            hasChange () {
                if (this.isCreateMode) {
                    return true
                }
                return Object.keys(this.form).some(key => this.form[key] !== this.account[key])
            },
            auth () {
                if (this.isCreateMode) {
                    return {
                        type: this.$OPERATION.C_CLOUD_ACCOUNT
                    }
                }
                return {
                    type: this.$OPERATION.U_CLOUD_ACCOUNT,
                    relation: [this.account.bk_account_id]
                }
            }
        },
        created () {
            this.getCloudAreaProperties()
            if (!this.isCreateMode) {
                Object.keys(DEFAULT_FORM).forEach(key => {
                    this.form[key] = this.account[key]
                })
                this.verifyResult = {
                    bk_cloud_vendor: this.account.bk_cloud_vendor,
                    bk_secret_id: this.account.bk_secret_id,
                    bk_secret_key: this.account.bk_secret_key,
                    connected: true
                }
            }
        },
        methods: {
            async getCloudAreaProperties () {
                try {
                    this.properties = await this.$store.dispatch('objectModelProperty/searchObjectAttribute', {
                        params: {
                            bk_obj_id: 'plat',
                            bk_supplier_account: this.supplierAccount
                        },
                        config: {
                            requestId: CLOUD_AREA_PROPERTIES,
                            fromCache: true
                        }
                    })
                } catch (error) {
                    console.error(error)
                }
            },
            handleSecretKeyFocus () {
                if (this.isCreateMode || this.form.bk_secret_key !== this.fakeSecret) {
                    return
                }
                this.form.bk_secret_key = ''
                this.$nextTick(async () => {
                    await this.$validator.validateAll()
                    this.errors.clear()
                })
            },
            handleSecretKeyBlur () {
                if (this.isCreateMode || this.form.bk_secret_key) {
                    return
                }
                this.form.bk_secret_key = this.fakeSecret
            },
            async handleTest () {
                if (!this.isCreateMode && this.form.bk_secret_key === this.fakeSecret) {
                    this.useBackendTest()
                } else {
                    this.useFrontendTest()
                }
            },
            async useBackendTest () {
                try {
                    const [result] = await this.$store.dispatch('cloud/account/getStatus', {
                        params: {
                            account_ids: [this.account.bk_account_id]
                        },
                        config: {
                            cancelPrevious: true,
                            requestId: this.request.verify
                        }
                    })
                    this.verifyResult = {
                        done: true,
                        connected: !result.err_msg,
                        msg: result.err_msg
                    }
                } catch (error) {
                    console.error(error)
                }
            },
            async useFrontendTest () {
                const results = await Promise.all([
                    this.$validator.validate('bk_cloud_vendor'),
                    this.$validator.validate('bk_secret_id'),
                    this.$validator.validate('bk_secret_key')
                ])
                if (results.some(valid => !valid)) {
                    return false
                }
                const params = { ...this.verifyParams }
                try {
                    this.verifyResult = {}
                    const response = await this.$store.dispatch('cloud/account/verify', {
                        params: params,
                        config: {
                            requestId: this.request.verify,
                            transformData: false
                        }
                    })
                    this.verifyResult = {
                        ...params,
                        connected: response.result,
                        msg: response.bk_error_msg,
                        done: true
                    }
                } catch (e) {
                    this.verifyResult = {
                        ...params,
                        done: true,
                        connected: false,
                        msg: e.message
                    }
                    console.error(e.message)
                }
            },
            async handleSubmit () {
                const valid = await this.$validator.validateAll()
                if (!valid) {
                    return false
                }
                if (this.verifyRequired) {
                    this.$error(this.$t('请先进行账户连通测试'))
                    return false
                }
                if (this.isCreateMode) {
                    this.doCreate()
                } else {
                    this.doUpdate()
                }
            },
            async doCreate () {
                try {
                    await this.$store.dispatch('cloud/account/create', {
                        params: this.form,
                        config: {
                            requestId: this.request.create
                        }
                    })
                    this.$success('新建成功')
                    this.container.hide()
                    RouterQuery.set({
                        _t: Date.now(),
                        page: RouterQuery.get('page', 1)
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            async doUpdate () {
                try {
                    const params = { ...this.form }
                    if (this.form.bk_secret_key === this.fakeSecret) {
                        delete params.bk_secret_key
                    }
                    await this.$store.dispatch('cloud/account/update', {
                        id: this.account.bk_account_id,
                        params: params,
                        config: {
                            requestId: this.request.update
                        }
                    })
                    this.$success('修改成功')
                    this.handleCancel()
                    RouterQuery.set({
                        _t: Date.now(),
                        page: RouterQuery.get('page', 1)
                    })
                } catch (e) {
                    console.error(e)
                }
            },
            handleCancel () {
                if (!this.isCreateMode) {
                    this.container.show({
                        type: 'details',
                        title: `${this.$t('账户详情')} 【${this.account.bk_account_name}】`,
                        props: {
                            id: this.account.bk_account_id
                        }
                    })
                } else {
                    this.container.hide()
                }
            },
            handleLinkToFAQ () {
                const FAQLink = {
                    '1': 'https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_enable-console-custom-url.html',
                    '2': 'https://cloud.tencent.com/document/product/598/37140'
                }
                const link = FAQLink[this.form.bk_cloud_vendor]
                if (link) {
                    window.open(FAQLink[this.form.bk_cloud_vendor])
                } else {
                    this.$warn('No link provided')
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-layout {
        padding: 10px 27px;
        .create-form-item + .create-form-item {
            margin-top: 13px;
            font-size: 0;
        }
        .create-form-meta {
            width: 100%;
            &.key {
                display: block;
                width: auto;
                overflow: hidden;
            }
            .vendor-selector-trigger {
                display: flex;
                height: 30px;
                padding-left: 10px;
            }
        }
        .create-form-button {
            margin-left: 10px;
        }
        .create-form-link {
            position: absolute;
            bottom: 100%;
            right: 0;
            /deep/ {
                .bk-link-text {
                    font-size: 12px;
                }
            }
        }
        .create-form-error {
            font-size: 12px;
            line-height: 1.2;
            color: $dangerColor;
        }
        .create-verify-success {
            font-size: 12px;
            line-height: 24px;
            .bk-icon {
                color: $successColor;
            }
        }
        .create-verify-fail {
            font-size: 12px;
            line-height: 24px;
            .bk-icon {
                color: $dangerColor;
            }
        }
        .create-form-options {
            margin-top: 20px;
        }
    }
</style>
