<template>
    <div class="bk-upload">
        <div class="file-wrapper" :class="{'isdrag': isdrag}">
            <img class="upload-icon" :src="uploadIcon">
            <span class="drop-upload">{{dragText}}</span>
            <span class="click-upload">{{clickText}}</span>
            <input ref="uploadel" @change="selectFile" :accept="accept" :multiple="multiple" type="file">
        </div>
        <p class="tip" v-if="tip">{{tip}}</p>
        <div class="all-file" v-if="fileList.length" >
            <div v-for="(file, index) in fileList" :key="index">
                <div class="file-item">
                    <div class="file-icon">
                        <img :src="getIcon(file)">
                    </div>
                    <i v-if="!file.done" class="bk-icon icon-close close-upload" @click="deleteFile(index, file)"></i>
                    <div class="file-info">
                        <div class="file-name"><span>{{file.name}}</span></div>
                        <div class="file-message">
                            <span class="upload-speed" v-show="!file.done && file.status === 'running'">{{speed}}{{unit}}</span>
                            <span class="file-size" v-show="!file.done">{{filesize(file.size)}}</span>
                            <span class="file-size done" v-show="file.done">{{t('uploadFile.uploadDone')}}</span>
                        </div>
                        <div class="progress-bar-wrapper">
                            <div :class="{'fail': file.errorMsg}" class="progress-bar" :style="{width: file.progress}"></div>
                        </div>
                    </div>
                </div>
                <p v-if="file.errorMsg" class="error-msg">{{file.errorMsg}}</p>
            </div>
        </div>
    </div>
</template>
<script>
/**
 * bk-upload
 * @module components/upload
 * @desc 文件上传组件
 * @param url（必传） {string}   文件上传到服务器的地址
 * @param name       {string}   - 服务器读取文件的key， 默认为'uplaod_file'
 * @param size       {number}   - 允许上传的文件大小
 * @param multiple   {boolean}  - 是否支持多选
 * @param accept     {string}   - 允许上传的文件类型
 * @param header     {string}   - 请求头
 */
import locale from '../../mixins/locale'
import uploadZip from '../../bk-magic-ui/src/images/uploadzip.svg'
import uploadFile from '../../bk-magic-ui/src/images/uploadfile.svg'
import uploadIcon from '../../bk-magic-ui/src/images/upload.svg'
export default {
    name: 'bk-upload',
    mixins: [locale],
    props: {
        name: {
            type: String,
            default: 'upload_file'
        },
        multiple: {
            type: Boolean,
            default: true
        },
        accept: {
            type: String,
            default: '*'
        },
        delayTime: {
            type: Number,
            default: 0
        },
        url: {
            required: true,
            type: String
        },
        size: {
            type: [Number, Object],
            default: function () {
                return {
                    maxFileSize: 5,
                    maxImgSize: 1
                }
            }
        },
        handleResCode: {
            type: Function,
            default: function (res) {
                if (res.code === 0) {
                    return true
                } else {
                    return false
                }
            }
        },
        header: [Array, Object],
        tip: {
            type: String,
            default: ''
        },
        validateName: {
            type: RegExp
        },
        withCredentials: {
            type: Boolean,
            default: false
        }
    },
    data () {
        return {
            dragText: this.t('uploadFile.drag'),
            clickText: this.t('uploadFile.click'),
            showDialog: true,
            fileList: [],
            width: 0,
            barEl: null,
            fileIndex: null,
            speed: 0,
            total: 0,
            unit: 'kb/s',
            isdrag: false,
            progress: 0,
            uploadIcon: uploadIcon
        }
    },
    watch: {
        'fileIndex' (val) {
            if (val !== null && val < this.fileList.length) {
                this.uploadFile(this.fileList[val], val)
            }
        }
    },
    methods: {
        filesize (val) {
            let size = val / 1000
            if (size < 1) {
                return `${val.toFixed(3)} KB`
            } else {
                let index = size.toString().indexOf('.')
                return `${size.toString().slice(0, index + 2)} MB`
            }
        },
        selectFile (e) {
            let file = null
            let files = Array.from(e.target.files)
            if (!files.length) return
            files.forEach((file, i) => {
                let fileObj = {
                    name: file.name,
                    originSize: file.size,
                    size: file.size / 1000,
                    maxFileSize: null,
                    maxImgSize: null,
                    type: file.type,
                    fileHeader: '',
                    origin: file,
                    base64: '',
                    status: '',
                    done: false,
                    responseData: '',
                    speed: null,
                    errorMsg: '',
                    progress: ''
                }
                let index = fileObj.type.indexOf('/')
                let type = fileObj.type.slice(0, index)
                let safariImageType = fileObj.type.indexOf('application/x-photoshop') > -1
                fileObj.fileHeader = type
                if (typeof this.size === 'number') {
                    fileObj.maxFileSize = this.size
                    fileObj.maxImgSize = this.size
                } else {
                    fileObj.maxFileSize = this.size.maxFileSize
                    fileObj.maxImgSize = this.size.maxImgSize
                }
                if (type === 'image' || safariImageType) {
                    this.handleImage(fileObj, file)
                }
                if ((type !== 'image' || !safariImageType) && fileObj.size > (fileObj.maxFileSize * 1000)) {
                    fileObj.errorMsg = `${fileObj.name}文件不能超过${fileObj.maxFileSize}MB`
                }
                if (this.validateName) {
                    if (!this.validateName.test(fileObj.name)) {
                        fileObj.errorMsg = '文件名不合法'
                    }
                }
                this.fileList.push(fileObj)
            })
            let len = this.fileList.length
            let fileIndex = this.fileIndex
            if (len - 1 === fileIndex) {
                this.uploadFile(this.fileList[fileIndex], fileIndex)
            } else {
                this.fileIndex = 0
            }
            e.target.value = ''
        },
        hideFileList () {
            if (this.delayTime) {
                setTimeout(() => {
                    this.fileList = []
                }, this.delayTime)
            }
        },
        uploadFile (fileObj) {
            if (fileObj.errorMsg) {
                this.fileIndex += 1
                fileObj.progress = 100 + '%'
                return
            }
            let formData = new FormData()
            let xhr = new XMLHttpRequest()
            formData.append(this.name, fileObj.origin)
            this.isdrag = false
            fileObj.xhr = xhr
            xhr.onreadystatechange = () => {
                if (xhr.readyState === 4) {
                    if (xhr.status === 200) {
                        try {
                            let response = JSON.parse(xhr.responseText)
                            if (this.handleResCode(response)) {
                                fileObj.done = true
                                fileObj.responseData = response
                                this.$emit('on-success', fileObj, this.fileList)
                            } else {
                                fileObj.errorMsg = response.message
                                this.$emit('on-error', fileObj, this.fileList)
                            }
                        } catch (error) {
                            fileObj.progress = 100 + '%'
                            fileObj.errorMsg = error.message
                        }
                    }
                    this.fileIndex += 1
                    this.unit = 'kb/s'
                    this.total = 0
                    fileObj.status = 'done'
                    if (this.fileIndex === this.fileList.length) {
                        this.$emit('on-done', this.fileList)
                        this.hideFileList()
                    }
                }
            }
            let uploadProgress = e => {
                if (e.lengthComputable) {
                    let percentComplete = Math.round(e.loaded * 100 / e.total)
                    let kb = Math.round(e.loaded / 1000)
                    fileObj.progress = percentComplete + '%'
                    this.speed = kb - this.total
                    this.total = kb
                    this.unit = 'kb/s'
                    if (this.speed > 1000) {
                        this.speed = Math.round(this.speed / 1000)
                        this.unit = 'mb/s'
                    }
                    this.$emit('on-progress', e, fileObj, this.fileList)
                }
                fileObj.status = 'running'
            }
            xhr.upload.addEventListener('progress', uploadProgress, false)
            xhr.withCredentials = this.withCredentials
            xhr.open('POST', this.url, true)
            if (this.header) {
                if (Array.isArray(this.header)) {
                    this.header.forEach(head => {
                        let headerKey = this.header.name
                        let headerVal = this.header.value
                        xhr.setRequestHeader(headerKey, headerVal)
                    })
                } else {
                    let headerKey = this.header.name
                    let headerVal = this.header.value
                    xhr.setRequestHeader(headerKey, headerVal)
                }
            }
            xhr.send(formData)
        },
        handleImage (fileObj, file) {
            let isJPGPNG = /image\/(jpg|png|jpeg)$/.test(fileObj.type)
            if (!isJPGPNG) {
                fileObj.errorMsg = '只允许上传JPG|PNG|JPEG格式的图片'
                return false
            }
            if (fileObj.size > (fileObj.maxImgSize * 1000)) {
                fileObj.errorMsg = `图片大小不能超过${fileObj.maxImgSize}MB`
                return false
            }
            let reader = new FileReader()
            reader.onload = (e) => {
                this.smallImage(reader.result, fileObj)
            }
            reader.readAsDataURL(file)
            return true
        },
        smallImage (result, fileObj) {
            let img = new Image()
            let canvas = document.createElement('canvas')
            let context = canvas.getContext('2d')
            img.onload = () => {
                let originWidth = img.width
                let originHeight = img.height
                let maxWidth = 42
                let maxHeight = 42
                let targetWidth = originWidth
                let targetHeight = originHeight
                if (originWidth > maxWidth || originHeight > maxHeight) {
                    if (originWidth / originHeight > maxWidth / maxHeight) {
                        targetWidth = maxWidth
                        targetHeight = Math.round(maxWidth * (originHeight / originWidth))
                    } else {
                        targetWidth = maxWidth
                        targetHeight = maxHeight
                    }
                }
                canvas.width = targetWidth
                canvas.height = targetHeight
                context.clearRect(0, 0, targetWidth, targetHeight)
                context.drawImage(img, 0, 0, targetWidth, targetHeight)
                fileObj['base64'] = canvas.toDataURL()
            }
            img.src = result
        },
        getIcon (file) {
            if (file.base64) {
                return file.base64
            }
            let isZip = false
            let zipType = ['zip', 'rar', 'tar', 'gz']
            for (let i = 0; i < zipType.length; i++) {
                if (file.type.indexOf(zipType[i]) > -1) {
                    isZip = true
                    break
                }
            }
            if (isZip) {
                return uploadZip
            } else {
                return uploadFile
            }
        },
        deleteFile (index, file) {
            if (file.xhr) {
                file.xhr.abort()
            }
            this.fileList.splice(index, 1)
            let len = this.fileList.length
            if (!len) {
                this.fileIndex = null
            }
            if (index === 0 && len) {
                this.fileIndex = 0
                this.uploadFile(this.fileList[0])
            }
        }
    },
    mounted () {
        let uploadEl = this.$refs.uploadel
        uploadEl.addEventListener('dragenter', e => {
            this.isdrag = true
        })
        uploadEl.addEventListener('dragleave', e => {
            this.isdrag = false
        })
        uploadEl.addEventListener('dragend', e => {
            this.isdrag = false
        })
    }
}
</script>
<style lang="scss">
    @import '../../bk-magic-ui/src/upload.scss'
</style>