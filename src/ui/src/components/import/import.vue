<template>
    <div class="import-wrapper">
        <slot name="prepend"></slot>
        <div class="up-file upload-file" v-bkloading="{ isLoading: isLoading }">
            <img src="../../assets/images/up_file.png">
            <input ref="fileInput" type="file" class="fullARea" @change.prevent="handleFile" />
            <i18n path="导入提示" tag="p" :places="{ allowType: allowType.join(','), maxSize: maxSizeLocal }">
                <b place="clickUpload">{{$t('点击上传')}}</b>
                <br place="breakRow">
            </i18n>
        </div>
        <div :class="['upload-file-info', { 'uploading': isLoading }, { 'fail': failed }, { 'uploaded': uploaded }]">
            <div class="upload-file-name" :title="fileInfo.name">{{fileInfo.name}}</div>
            <div class="upload-file-size fr">{{fileInfo.size}}</div>
            <div class="upload-file-status" hidden>{{fileInfo.status}}</div>
            <div class="upload-file-status-icon" hidden>
                <i :class="['bk-icon ', { 'icon-check-circle-shape': uploaded,'icon-close-circle-shape': failed }]"></i>
            </div>
        </div>
        <div class="upload-details">
            <slot name="uploadErrorMessage"></slot>
            <template v-if="$slots.uploadResult">
                <slot name="uploadResult"></slot>
            </template>
            <div v-else-if="hasUploadError()">
                <div class="upload-details-success" v-if="uploadResult.success && uploadResult.success.length">
                    <i class="bk-icon icon-check-circle-shape"></i>
                    <slot name="successTips" v-bind="uploadResult">
                        <span>{{$t(successTips, { N: uploadResult.success.length })}}</span>
                    </slot>
                </div>
                <!-- 上传失败列表  -->
                <div class="upload-details-fail" v-if="uploadResult.error && uploadResult.error.length">
                    <div class="upload-details-fail-title">
                        <i class="bk-icon icon-close-circle-shape"></i>
                        <slot name="errorTips" v-bind="uploadResult">
                            <span>{{$t(errorTips)}}({{uploadResult.error.length}})</span>
                        </slot>
                    </div>
                    <ul ref="failList" class="upload-details-fail-list">
                        <li v-for="(errorMsg, index) in uploadResult.error" :title="errorMsg" :key="index">{{errorMsg}}</li>
                    </ul>
                </div>
                <div class="upload-details-fail" v-if="uploadResult.update_error && uploadResult.update_error.length">
                    <div class="upload-details-fail-title">
                        <i class="bk-icon icon-close-circle-shape"></i>
                        <slot name="updateErrorTips" v-bind="uploadResult">
                            <span>{{$t(updateErrorTips)}}({{uploadResult.update_error.length}})</span>
                        </slot>
                    </div>
                    <ul ref="failList" class="upload-details-fail-list">
                        <li v-for="(errorMsg, index) in uploadResult.update_error" :title="errorMsg" :key="index">{{errorMsg}}</li>
                    </ul>
                </div>
                <div class="upload-details-fail" v-if="uploadResult.asst_error && uploadResult.asst_error.length">
                    <div class="upload-details-fail-title">
                        <i class="bk-icon icon-close-circle-shape"></i>
                        <span>关联关系导入失败列表({{uploadResult.asst_error.length}})</span>
                    </div>
                    <ul ref="failList" class="upload-details-fail-list">
                        <li v-for="(errorMsg, index) in uploadResult.asst_error" :title="errorMsg" :key="index">{{errorMsg}}</li>
                    </ul>
                </div>
            </div>
        </div>
        <div class="clearfix down-model-content" v-if="templdateAvailable">
            <slot name="download-desc"></slot>
            <a href="javascript:void(0);" style="text-decoration: none;" v-if="templateUrl" @click="handleDownloadTemplate">
                <img src="../../assets/images/icon/down_model_icon.png">
                <span class="submit-btn">{{$t('下载模板')}}</span>
            </a>
        </div>
    </div>
</template>

<script type="text/javascript">
    export default {
        name: 'cmdb-import',
        props: {
            templateUrl: {
                type: String,
                required: true
            },
            downloadPayload: {
                type: Object,
                default () {
                    return {}
                }
            },
            importPayload: {
                type: Object,
                default () {
                    return {}
                }
            },
            importUrl: {
                type: String,
                required: true
            },
            allowType: {
                type: Array,
                default () {
                    return ['xlsx']
                }
            },
            maxSize: {
                type: Number,
                default: 20 * 1024 // kb
            },
            uploadDone: {
                type: Function,
                default: null
            },
            templdateAvailable: {
                type: Boolean,
                default: true
            },
            globalError: {
                type: Boolean,
                default: true
            },
            successTips: {
                type: String,
                default: '成功上传N条数据'
            },
            errorTips: {
                type: String,
                default: '上传失败列表'
            },
            updateErrorTips: {
                type: String,
                default: '更新失败列表'
            },
            beforeUpload: Function
        },
        data () {
            return {
                isLoading: false,
                uploaded: false,
                failed: false,
                fileInfo: {
                    name: '',
                    size: '',
                    status: ''
                },
                uploadResult: {
                    success: null,
                    error: null,
                    update_error: null,
                    asst_error: null
                }
            }
        },
        computed: {
            allowTypeRegExp () {
                return new RegExp(`^.*?.(${this.allowType.join('|')})$`)
            },
            maxSizeLocal () {
                const maxSize = this.maxSize * 1024
                return this.formatSize(maxSize)
            }
        },
        methods: {
            handleFile (e) {
                this.reset()

                if (this.beforeUpload && this.beforeUpload() === false) {
                    this.$refs.fileInput.value = ''
                    return
                }

                const files = e.target.files
                const fileInfo = files[0]
                if (!this.allowTypeRegExp.test(fileInfo.name)) {
                    this.$refs.fileInput.value = ''
                    this.$error(this.$t('文件格式非法', { allowType: this.allowType.join(',') }))
                    return false
                } else if (fileInfo.size / 1024 > this.maxSize) {
                    this.$refs.fileInput.value = ''
                    this.$error(this.$t('文件大小溢出', { maxSize: this.maxSizeLocal }))
                    return false
                } else {
                    this.fileInfo.name = fileInfo.name
                    this.fileInfo.size = this.formatSize(fileInfo.size, 2)
                    const formData = new FormData()
                    formData.append('file', files[0])
                    for (const [key, value] of Object.entries(this.importPayload)) {
                        formData.append(key, value)
                    }
                    this.isLoading = true
                    this.$http.post(this.importUrl, formData, { transformData: false, globalError: false }).then(res => {
                        const defaultResult = {
                            success: null,
                            error: null,
                            update_error: null,
                            asst_error: null
                        }
                        this.uploadResult = Object.assign(this.uploadResult, res.data || defaultResult)
                        if (res.result) {
                            this.uploaded = true
                            this.fileInfo.status = this.$t('成功')
                            this.$emit('success', res)
                        } else if (res.data && res.data.success) {
                            this.failed = true
                            this.fileInfo.status = this.$t('部分成功')
                            this.$emit('partialSuccess', res)
                        } else {
                            this.failed = true
                            this.fileInfo.status = this.$t('失败')
                            this.globalError && this.$error(res['bk_error_msg'])
                            this.$emit('error', res)
                        }
                        this.$refs.fileInput.value = ''
                        this.isLoading = false

                        this.$emit('upload-done', res)
                    }).catch(error => {
                        this.reset()
                        this.isLoading = false
                        this.$emit('error', error)
                    })
                }
            },
            hasUploadError () {
                const uploadResult = this.uploadResult
                return (uploadResult.success && uploadResult.success.length)
                    || (uploadResult.error && uploadResult.error.length)
                    || (uploadResult.update_error && uploadResult.update_error.length)
                    || (uploadResult.asst_error && uploadResult.asst_error.length)
            },
            reset () {
                this.uploaded = false
                this.failed = false
                this.fileInfo = {
                    name: '',
                    size: '',
                    status: ''
                }
                this.uploadResult = {
                    success: null,
                    error: null,
                    update_error: null,
                    asst_error: null
                }
            },
            formatSize (value, digits = 0) {
                const uints = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
                const index = Math.floor(Math.log(value) / Math.log(1024))
                let size = value / Math.pow(1024, index)
                size = `${size.toFixed(digits)}${uints[index]}`
                return size
            },
            async handleDownloadTemplate () {
                try {
                    let data = this.downloadPayload
                    if (!(data instanceof FormData)) {
                        data = new FormData()
                        Object.keys(this.downloadPayload).forEach(key => {
                            const value = this.downloadPayload[key]
                            if (typeof value === 'object') {
                                data.append(key, JSON.stringify(value))
                            } else {
                                data.append(key, value)
                            }
                        })
                    }
                    this.$http.download({
                        url: this.templateUrl,
                        data: data
                    })
                } catch (e) {
                    console.log(e)
                }
            }
        }
    }
</script>

<style media="screen" lang="scss" scoped>
    .import-wrapper {
        height: 100%;
        @include scrollbar-y;
    }
    .up-file{
        .up-file-text{
            p{
                font-size:14px;
                font-weight: bold;
                line-height:1;
                color:#bec6de;
            }
            .click-text{
                color: #3c96ff;
                cursor:pointer;
                position:relative;
            }
        }
    }
    .input-file{
        left: 0;
        top: 0;
        color: #3c96ff;
        cursor: pointer;
        position: relative;
        border: none !important;
        text-decoration: none;
        &:hover{
            background: none !important;
        }
    }
    .submit-btn{
        display: inline-block;
        vertical-align: 2px;
        border: none;
        background: #fff;
        padding: 0;
        color: #3c96ff;
        outline: none;
        &:hover{
            color:#3c96ff;
        }
    }

    .upload-file {
        position: relative;
        height:182px;
        margin: 30px 29px 0 33px;
        padding: 33px 0;
        text-align: center;
        overflow: hidden;
        background-color: #f9f9f9;
        cursor: pointer;
        border: 1px solid transparent;
        -webkit-transition: all .5s ease;
        transition: all .5s ease;
        &:hover{
            background-color: #fff;
            border-color: #3c96ff;
            box-shadow: 0 4px 6px rgba(0,0,0,.1);
            p {
                color: #6b7baa;
            }
        }
        p {
            margin: 23px 0 0 0;
            line-height: 18px;
            font-size: 14px;
            font-weight: bold;
            color: #bec6de;
            b {
                color: #3c96ff;
            }
        }
    }

    .down-model-content {
        padding: 10px 30px;
        font-size: 14px;
    }

    .upload-file{
        &-name,
        &-size,
        &-status,
        &-status-icon{
            float: left;
            position: relative;
            height: 100%;
        }
    }
    .upload-file-info {
        overflow: hidden;
        line-height: 50px;
        position: relative;
        margin: 2px 29px 0 33px;
        .bk-icon{
            vertical-align: -1px;
        }
        .icon-check-circle-shape{
            color: #4cd084;
        }
        .icon-close-circle-shape{
            color: red;
        }
        &:before {
            content: '';
            position: absolute;
            left: 0;
            top: 0;
            width: 0%;
            height: 100%;
        }

        &.uploading {
            &:before {
                background: #e3f5eb;
                width: 90%;
                transition: width 20s;
            }
        }

        &.uploaded {
            &:before {
                background: #e3f5eb;
                width: 100%;
                transition: width 1s;
            }
            background: #e3f5eb;
        }

        &.fail {
            background:#f9f9f9;
            &:before {
                content: '';
                position: absolute;
                left: 0;
                top: 0;
                width: 0;
                height: 0;
                background: #e3f5eb;
                -webkit-transition: all .5s;
                transition: all .5s;
            }

        }

        .upload-file-name {
            width: calc(100% - 200px);
            overflow: hidden;
            white-space: nowrap;
            text-overflow: ellipsis;
            padding: 0 24px;
        }

        .upload-file-size {
            /* width: 174px; */
            padding: 0 24px;
            span{
                color: #c7cee3;
            }
        }

        .upload-file-meta {
            width: 8%;
        }

        .upload-file-status {
            width: 212px;
        }

    }
    .upload-details {
        margin: 2px 29px 0 33px;
        font-size: 14px;
        &-success{
            padding: 0 21px;
            line-height: 56px;
            background-color: #f9f9f9;
            color: #34d97b;
        }
        &-fail{
            margin: 2px 0 0 0;
            padding: 13px 0 15px;
            line-height: 32px;
            background-color: #f9f9f9;
            color: #ef4c4c;
            &-title{
                padding: 0 21px;
            }
            &-list{
                line-height: 28px;
                color: #6b7baa;
                font-size: 12px;
                white-space: nowrap;
                overflow: auto;
                li{
                    padding: 0 43px;
                    overflow: hidden;
                    text-overflow: ellipsis;
                }
                li:hover{
                    background-color: #edf5ff;
                }
            }
        }
    }
    .fullARea {
        position: absolute;
        cursor: pointer;
        left: 0;
        top: 0;
        width: 100%;
        height: 100%;
        opacity: 0;
        filter: alpha(opacity=0);
        cursor: pointer;
        -webkit-transition: all .5s ease;
        transition: all .5s ease;
    }
</style>
