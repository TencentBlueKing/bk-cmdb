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
                    data-vv-name="bk_cloud_vendor"
                    v-model="form.bk_cloud_vendor"
                    v-validate="'required'">
                    <bk-option v-for="vendor in vendors"
                        :key="vendor.id"
                        :name="vendor.name"
                        :id="vendor.id">
                    </bk-option>
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
                <link-button class="create-form-link">如何获取ID和Key?</link-button>
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
                    v-validate="'required|length:256'">
                </bk-input>
                <template v-if="verifyResult.hasOwnProperty('done')">
                    <p class="create-verify-success" v-if="verifyResult.connected">
                        <i class="bk-icon icon-check-circle-shape"></i>
                        {{$t('账户连通成功')}}
                    </p>
                    <p class="create-verify-fail" v-else>
                        <i class="bk-icon icon-close-circle-shape"></i>
                        {{$t('账户连通成功')}}
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
                <bk-button class="mr10" theme="primary" @click.stop.prevent="handleSubmit">{{isCreateMode ? $t('提交') : $t('保存')}}</bk-button>
                <bk-button theme="default" @click.stop.prevent="handleCancel(null)">{{$t('取消')}}</bk-button>
            </div>
        </bk-form>
    </div>
</template>

<script>
    import vendors from '@/dictionary/cloud-vendor'
    const DEFAULT_FORM = {
        bk_account_name: '',
        bk_cloud_vendor: '',
        bk_secret_id: '',
        bk_secret_key: '',
        bk_description: ''
    }
    export default {
        name: 'cloud-account-form',
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
                vendors: vendors,
                form: {
                    ...DEFAULT_FORM
                },
                verifyResult: {},
                request: {
                    create: Symbol('create'),
                    update: Symbol('update'),
                    verify: Symbol('verify')
                }
            }
        },
        computed: {
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
            }
        },
        created () {
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
            async handleTest () {
                const isValid = await this.$validator.validateAll()
                if (!isValid) {
                    return false
                }
                const params = { ...this.verifyParams }
                try {
                    this.verifyResult = {}
                    const result = await this.$store.dispatch('cloud/account/verify', {
                        params: params,
                        config: {
                            requestId: this.request.verify
                        }
                    })
                    this.verifyResult = {
                        ...result,
                        ...params,
                        done: true
                    }
                } catch (e) {
                    this.verifyResult = {
                        ...params,
                        done: true,
                        connected: false,
                        error_msg: e.message
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
                    this.$error(this.$t('账户连通测试提示'))
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
                    this.handleHide('request-refresh')
                } catch (e) {
                    console.error(e)
                }
            },
            async doUpdate () {
                try {
                    await this.$store.dispatch('cloud/account/update', {
                        id: this.account.bk_account_id,
                        params: this.form,
                        config: {
                            requestId: this.request.update
                        }
                    })
                    this.$success('修改成功')
                    this.handleCancel('request-refresh')
                } catch (e) {
                    console.error(e)
                }
            },
            handleCancel (eventType) {
                if (!this.isCreateMode) {
                    this.container.show({
                        type: 'details',
                        title: `${this.$t('账户详情')} 【${this.account.bk_account_name}】`,
                        props: {
                            id: 1
                        }
                    })
                } else {
                    this.handleHide(eventType)
                }
            },
            handleHide (eventType) {
                this.container.hide(eventType)
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
        }
        .create-form-button {
            margin-left: 10px;
        }
        .create-form-link {
            position: absolute;
            bottom: 100%;
            right: 0;
            font-size: 12px;
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
