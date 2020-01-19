<template>
    <bk-select v-model="selected"
        ref="selector"
        searchable
        size="small"
        font-size="small"
        :placeholder="$t('请选择xx', { name: $t('资源目录') })"
        @toggle="handleSelectToggle">
        <bk-option v-for="folder in folders"
            :key="folder.id"
            :id="folder.id"
            :name="folder.name">
        </bk-option>
        <bk-option class="create-option"
            ref="createOptionComponent"
            :id="createFolderId"
            :name="$t('新增目录')"
            :disabled="true"
            @click.native.stop="handleCreateClick">
            <template v-if="!createMode">{{$t('新增目录')}}</template>
            <template v-else>
                <bk-input ref="input"
                    size="small"
                    font-size="small"
                    :placeholder="$t('请输入目录名称，回车结束')"
                    v-model.trim="newFolder"
                    @enter="handleConfirmCreate">
                </bk-input>
            </template>
        </bk-option>
        <a href="javascript:void(0)" class="extension-link" slot="extension">
            <i class="bk-icon icon-plus-circle"></i>
            {{$t('申请其他目录权限')}}
        </a>
    </bk-select>
</template>

<script>
    export default {
        name: 'cloud-resource-folder-selector',
        props: {
            value: {
                type: [String, Number],
                default: ''
            }
        },
        data () {
            return {
                selected: this.value,
                createMode: false,
                createFolderId: 'createFolderId',
                newFolder: '',
                folders: [{
                    id: 'ali',
                    name: '阿里云'
                }, {
                    id: 'tcloud',
                    name: '腾讯云'
                }]
            }
        },
        watch: {
            value (value) {
                this.selected = value
            },
            selected (value, oldValue) {
                this.$emit('input', value)
                this.$emit('change', value, oldValue)
            }
        },
        methods: {
            handleSelectToggle (isVisible) {
                if (!isVisible) {
                    this.toggleCreate(false)
                }
            },
            handleCreateClick () {
                this.toggleCreate(true)
            },
            toggleCreate (isCreateMode) {
                this.createMode = isCreateMode
                if (isCreateMode) {
                    this.$nextTick(() => {
                        this.$refs.input.focus()
                    })
                } else {
                    this.newFolder = ''
                }
            },
            async handleConfirmCreate () {
                if (!this.newFolder.length) {
                    return false
                }
                try {
                    const id = await Promise.resolve(Date.now())
                    this.folders.push({
                        id: id,
                        name: this.newFolder
                    })
                    this.selected = id
                    this.toggleCreate(false)
                } catch (e) {
                    console.error(e)
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .create-option.is-disabled {
        cursor: pointer;
        color: $textColor;
    }
</style>
