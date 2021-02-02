<template>
    <div class="confirm-wrapper">
        <h2 class="title">{{$t('确认归还主机池')}}</h2>
        <i18n tag="p" path="确认归还主机池主机数量" class="content">
            <span class="count" place="count">{{count}}</span>
        </i18n>
        <p class="content">{{$t('归还主机池提示')}}</p>
        <div class="directory">
            {{$t('归还至目录')}}
            <bk-select class="directory-selector ml10"
                v-model="target"
                size="small"
                searchable
                :clearable="false"
                :loading="$loading(request.list)"
                :popover-options="{
                    boundary: 'window'
                }">
                <cmdb-auth-option v-for="directory in directories"
                    :key="directory.bk_module_id"
                    :id="directory.bk_module_id"
                    :name="directory.bk_module_name"
                    :auth="{ type: $OPERATION.HOST_TO_RESOURCE, relation: [[[bizId], [directory.bk_module_id]]] }">
                </cmdb-auth-option>
            </bk-select>
        </div>
        <div class="options">
            <bk-button class="mr10" theme="primary"
                :disabled="!target"
                @click="handleConfirm">{{$t('确定')}}</bk-button>
            <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-move-to-resource-confirm',
        props: {
            count: {
                type: Number,
                default: 0
            },
            bizId: Number
        },
        data () {
            return {
                target: '',
                directories: [],
                request: {
                    list: Symbol('list')
                }
            }
        },
        created () {
            this.getDirectories()
        },
        methods: {
            async getDirectories () {
                try {
                    const { info } = await this.$store.dispatch('resourceDirectory/getDirectoryList', {
                        params: {
                            page: {
                                sort: 'bk_module_name'
                            }
                        },
                        config: {
                            requestId: this.request.list
                        }
                    })
                    this.directories = info
                } catch (error) {
                    console.error(error)
                }
            },
            handleConfirm () {
                this.$emit('confirm', this.target)
            },
            handleCancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .confirm-wrapper {
        text-align: center;
    }
    .title {
        margin: 45px 0 17px;
        line-height: 32px;
        font-size:24px;
        font-weight: normal;
        color: #313238;
    }
    .content {
        line-height:20px;
        font-size:14px;
        color: $textColor;
        .count {
            font-weight: bold;
            color: #EA3536;
            padding: 0 4px;
        }
    }
    .directory {
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        margin-top: 10px;
        .directory-selector {
            margin-left: 10px;
            width: 150px;
            text-align: left;
        }
    }
    .options {
        margin: 20px 0;
        font-size: 0;
    }
</style>
