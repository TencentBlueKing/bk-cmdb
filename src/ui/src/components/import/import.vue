/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template lang="html">
    <div>
        <div class="up-file upload-file" v-bkloading="{isLoading: isLoading}">
            <img src="../../common/images/up_file.png">
            <input ref="fileInput" type="file" class="fullARea" @change.prevent="handleFile"/>
            <i18n path="Inst['导入提示']" tag="p" :places="{allowType: allowType.join(','), maxSize: maxSize}">
                <b place="clickUpload">{{$t("Inst['点击上传']")}}</b>
                <br place="breakRow">
            </i18n>
        </div>
        <div :class="['upload-file-info', {'success': uploaded}, {'fail': failed}]">
            <div class="upload-file-name">{{fileInfo.name}}</div>
            <div class="upload-file-size fr">{{fileInfo.size}}</div>
            <div class="upload-file-status" hidden>{{fileInfo.status}}</div>
            <div class="upload-file-status-icon" hidden>
                <i :class="['bk-icon ',{'icon-check-circle-shape':uploaded,'icon-close-circle-shape':failed}]"></i>
            </div>
        </div>
        <div class="upload-details" v-if="(uploadResult.success && uploadResult.success.length) || (uploadResult.error && uploadResult.error.length) || (uploadResult.update_error && uploadResult.update_error.length)">
            <div class="upload-details-success" v-if="uploadResult.success && uploadResult.success.length">
                <i class="bk-icon icon-check-circle-shape"></i>
                <span>{{$t("Inst['成功上传N条数据']", {N: uploadResult.success.length})}}</span>
            </div>
            <!-- 上传失败列表  -->
            <div class="upload-details-fail" v-if="uploadResult.error && uploadResult.error.length">
                <div class="upload-details-fail-title">
                    <i class="bk-icon icon-close-circle-shape"></i>
                    <span>{{$t("Inst['上传失败列表']")}}({{uploadResult.error.length}})</span>
                </div>
                <ul ref="failList" class="upload-details-fail-list">
                    <li v-for="(errorMsg, index) in uploadResult.error" :title="errorMsg">{{errorMsg}}</li>
                </ul>
            </div>
            <div class="upload-details-fail" v-if="uploadResult.update_error && uploadResult.update_error.length">
                <div class="upload-details-fail-title">
                    <i class="bk-icon icon-close-circle-shape"></i>
                    <span>{{$t("Inst['更新失败列表']")}}({{uploadResult.update_error.length}})</span>
                </div>
                <ul ref="failList" class="upload-details-fail-list">
                    <li v-for="(errorMsg, index) in uploadResult.update_error" :title="errorMsg">{{errorMsg}}</li>
                </ul>
            </div>
        </div>
        <div class="clearfix down-model-content">
            <slot name="download-desc"></slot>
            <a :href="templateUrl" style="text-decoration: none;">
                <img src="../../common/images/icon/down_model_icon.png" alt="">
                <span class="submit-btn">{{$t("Inst['下载模版']")}}</span>
            </a>
        </div>
    </div>
</template>

<script type="text/javascript">
    export default {
        props: {
            templateUrl: {
                type: String,
                required: true
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
                default: 500 // kb
            }
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
                    update_error: null
                }
            }
        },
        computed: {
            allowTypeRegExp () {
                return new RegExp(`^.*?.(${this.allowType.join('|')})$`)
            }
        },
        methods: {
            handleFile (e) {
                this.reset()
                let files = e.target.files
                let fileInfo = files[0]
                if (!this.allowTypeRegExp.test(fileInfo.name)) {
                    this.$refs.fileInput.value = ''
                    this.$alertMsg(this.$t("Inst['文件格式非法']", {allowType: this.allowType.join(',')}))
                    return false
                } else if (fileInfo.size / 1024 > this.maxSize) {
                    this.$refs.fileInput.value = ''
                    this.$alertMsg(this.$t("Inst['文件大小溢出']", {maxSize: this.maxSize}))
                    return false
                } else {
                    this.fileInfo.name = fileInfo.name
                    this.fileInfo.size = `${(fileInfo.size / 1024).toFixed(2)}kb`
                    let formData = new FormData()
                    formData.append('file', files[0])
                    this.isLoading = true
                    this.$axios.post(this.importUrl, formData).then(res => {
                        this.uploadResult = Object.assign(this.uploadResult, res.data || {success: null, error: null, update_error: null})
                        if (res.result) {
                            this.uploaded = true
                            this.fileInfo.status = this.$t("Inst['成功']")
                            this.$emit('success', res)
                        } else if (res.data && res.data.success) {
                            this.failed = true
                            this.fileInfo.status = this.$t("Inst['部分成功']")
                            this.$emit('partialSuccess', res)
                        } else {
                            this.failed = true
                            this.fileInfo.status = this.$t("Inst['失败']")
                            this.$alertMsg(res['bk_error_msg'])
                            this.$emit('error', res)
                        }
                        this.$refs.fileInput.value = ''
                        this.$nextTick(() => {
                            this.calcFailListHeight()
                        })
                        this.isLoading = false
                    }).catch(error => {
                        this.reset()
                        this.isLoading = false
                        this.$emit('error', error)
                    })
                }
            },
            calcFailListHeight () {
                const failListOffsetHeight = 550
                const maxHeight = document.body.getBoundingClientRect().height - failListOffsetHeight
                let failList = this.$refs.failList
                if (failList) {
                    if (Array.isArray(failList)) {
                        failList.map(list => {
                            list.style.maxHeight = `${maxHeight / failList.length}px`
                        })
                    } else {
                        failList.style.maxHeight = `${maxHeight}px`
                    }
                }
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
                    update_error: null
                }
            }
        }
    }
</script>

<style media="screen" lang="scss" scoped>
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

        &.success {
            background:#f9f9f9;
            &:before {
                content: '';
                position: absolute;
                left: 0;
                top: 0;
                width: 100%;
                height: 100%;
                background: #e3f5eb;
                -webkit-transition: all .5s;
                transition: all .5s;
            }

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
            width: 245px;
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
    .upload-details{
        margin: 2px 29px 0 33px;
        &-success{
            padding: 0 21px;
            line-height: 56px;
            background-color: #f9f9f9;
            color: #34d97b;
            .bk-icon{
                vertical-align: -1px;
            }
        }
        &-fail{
            margin: 2px 0 0 0;
            padding: 13px 0 15px;
            line-height: 32px;
            background-color: #f9f9f9;
            color: #ef4c4c;
            &-title{
                padding: 0 21px;
                .bk-icon{
                    vertical-align: -1px;
                }
            }
            &-list{
                line-height: 28px;
                color: #6b7baa;
                font-size: 12px;
                white-space: nowrap;
                overflow: auto;
                &::-webkit-scrollbar{
                    width: 6px;
                }
                &::-webkit-scrollbar-thumb{
                    border-radius: 3px;
                    background: #c7cee3;
                }
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