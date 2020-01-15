<template>
    <div class="create-layout" slot="content">
        <bk-form form-type="vertical">
            <bk-form-item class="create-form-item" :label="$t('账户名称')" required>
                <bk-input class="create-form-meta"
                    :placeholder="$t('请输入xx', { name: $t('账户名称') })"
                    :data-vv-as="$t('账户名称')"
                    data-vv-name="name"
                    v-validate="'required|length:256'"
                    v-model.trim="form.name">
                </bk-input>
                <p class="create-form-error" v-if="errors.has('name')">{{errors.first('name')}}</p>
            </bk-form-item>
            <bk-form-item class="create-form-item" :label="$t('账户类型')" required>
                <bk-select class="create-form-meta"
                    :placeholder="$t('请选择xx', { name: $t('账户类型') })"
                    :data-vv-as="$t('账户类型')"
                    data-vv-name="type"
                    v-model="form.type"
                    v-validate="'required'">
                    <bk-option v-for="type in typeList"
                        :key="type.id"
                        :name="type.name"
                        :id="type.id">
                    </bk-option>
                </bk-select>
                <p class="create-form-error" v-if="errors.has('type')">{{errors.first('type')}}</p>
            </bk-form-item>
            <bk-form-item class="create-form-item" label="ID" required>
                <bk-input class="create-form-meta"
                    :placeholder="$t('请输入xx', { name: 'ID' })"
                    data-vv-as="ID"
                    data-vv-name="id"
                    v-model.trim="form.id"
                    v-validate="'required|length:256'">
                </bk-input>
                <link-button class="create-form-link">如何获取ID和Key?</link-button>
                <p class="create-form-error" v-if="errors.has('id')">{{errors.first('id')}}</p>
            </bk-form-item>
            <bk-form-item class="create-form-item clearfix" label="Key" required>
                <bk-button class="create-form-button fr" @click="handleTest">{{$t('连通测试')}}</bk-button>
                <bk-input class="create-form-meta key"
                    :placeholder="$t('请输入xx', { name: 'Key' })"
                    data-vv-as="Key"
                    data-vv-name="key"
                    v-model.trim="form.key"
                    v-validate="'required|length:256'">
                </bk-input>
                <p class="create-test-success" v-if="testState === 'success'">
                    <i class="bk-icon icon-check-circle-shape"></i>
                    {{$t('账户连通成功')}}
                </p>
                <p class="create-test-fail" v-else-if="testState === 'fail'">
                    <i class="bk-icon icon-close-circle-shape"></i>
                    {{$t('账户连通失败')}}
                </p>
                <p class="create-form-error" v-else-if="errors.has('key')">{{errors.first('key')}}</p>
            </bk-form-item>
            <bk-form-item class="create-form-item" :label="$t('备注')">
                <bk-input class="create-form-meta" type="textarea"
                    :placeholder="$t('请输入xx', { name: $t('备注') })"
                    :data-vv-as="$t('备注')"
                    data-vv-name="remarks"
                    v-model.trim="form.remarks"
                    v-validate="'length:2000'">
                </bk-input>
                <p class="create-form-error" v-if="errors.has('remarks')">{{errors.first('remarks')}}</p>
            </bk-form-item>
            <bk-form-item class="create-form-options">
                <bk-button class="mr10" theme="primary" @click.stop.prevent="handleSubmit">{{$t('提交')}}</bk-button>
                <bk-button theme="default" @click.stop.prevent="handleCancel">{{$t('取消')}}</bk-button>
            </bk-form-item>
        </bk-form>
    </div>
</template>

<script>
    const DEFAULT_FORM = {
        name: '',
        type: '',
        id: '',
        key: '',
        remarks: ''
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
                typeList: [{
                    id: 'ali',
                    name: '阿里云'
                }],
                form: {
                    ...DEFAULT_FORM
                },
                testState: null
            }
        },
        created () {
            if (this.mode === 'edit') {
                this.form = Object.assign({}, DEFAULT_FORM, this.account)
            }
        },
        methods: {
            handleTest () {
                this.testState = 'fail'
            },
            async handleSubmit () {
                const valid = await this.$validator.validateAll()
                if (!valid) {
                    return false
                }
                if (this.mode === 'create') {
                    this.doCreate()
                } else {
                    this.doUpdate()
                }
            },
            async doCreate () {
                try {
                    await Promise.resolve()
                    this.$success('新建成功')
                    this.handleHide('create-success')
                } catch (e) {
                    console.error(e)
                }
            },
            async doUpdate () {
                try {
                    await Promise.resolve({})
                    this.$success('修改成功')
                    this.handleCancel()
                } catch (e) {
                    console.error(e)
                }
            },
            handleCancel () {
                if (this.mode === 'edit') {
                    this.container.show({
                        type: 'details',
                        title: `${this.$t('账户详情')} 【${this.account.name}】`,
                        props: {
                            id: 1
                        }
                    })
                } else {
                    this.handleHide(null)
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
        .create-test-success {
            font-size: 12px;
            line-height: 24px;
            .bk-icon {
                color: $successColor;
            }
        }
        .create-test-fail {
            font-size: 12px;
            line-height: 24px;
            .bk-icon {
                color: $dangerColor;
            }
        }
        .create-form-options.bk-form-item {
            margin-top: 20px;
        }
    }
</style>
