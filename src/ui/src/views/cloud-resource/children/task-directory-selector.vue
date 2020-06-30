<template>
    <bk-select
        v-if="display === 'selector'"
        v-model="selected"
        ref="selector"
        searchable
        size="small"
        font-size="small"
        :clearable="false"
        :placeholder="$t('请选择xx', { name: $t('资源目录') })"
        @toggle="handleSelectToggle">
        <bk-option v-for="directory in directories"
            :key="directory.bk_module_id"
            :id="directory.bk_module_id"
            :name="directory.bk_module_name">
        </bk-option>
        <bk-option class="create-option"
            ref="createOptionComponent"
            :id="createDirectoryId"
            :name="$t('新增目录')"
            :disabled="true"
            @click.native.stop="handleCreateClick">
            <template v-if="!createMode">{{$t('新增目录')}}</template>
            <template v-else>
                <bk-input ref="input"
                    size="small"
                    font-size="small"
                    :placeholder="$t('请输入目录名称，回车结束')"
                    v-model.trim="newDirectory"
                    @enter="handleConfirmCreate">
                </bk-input>
            </template>
        </bk-option>
        <a href="javascript:void(0)" class="extension-link" slot="extension">
            <i class="bk-icon icon-plus-circle"></i>
            {{$t('申请其他目录权限')}}
        </a>
    </bk-select>
    <span v-else>{{getInfo()}}</span>
</template>

<script>
    import symbols from '../common/symbol'
    export default {
        name: 'task-directory-selector',
        props: {
            value: {
                type: [String, Number],
                default: ''
            },
            display: {
                type: String,
                default: 'selector'
            }
        },
        data () {
            return {
                createMode: false,
                createDirectoryId: 'createDirectoryId',
                newDirectory: '',
                directories: [],
                request: {
                    findMany: symbols.get('directory'),
                    create: Symbol('create')
                }
            }
        },
        computed: {
            selected: {
                get () {
                    return this.value
                },
                set (value, oldValue) {
                    this.$emit('input', value)
                    this.$emit('change', value, oldValue)
                }
            }
        },
        created () {
            this.getDirectories()
        },
        methods: {
            async getDirectories () {
                try {
                    const { info } = await this.$store.dispatch('resource/directory/findMany', {
                        params: {
                            sort: 'bk_module_id'
                        },
                        config: {
                            requestId: this.request.findMany,
                            fromCache: true
                        }
                    })
                    // 直接进行赋值，后面新增目录后，其他的地方也能获得相同的数据
                    this.directories = info
                    if (!this.selected && info.length) {
                        this.selected = info[0].bk_module_id
                    }
                } catch (e) {
                    this.directories = []
                    console.error(e)
                }
            },
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
                    this.newDirectory = ''
                    this.$refs.selector.close()
                }
            },
            async handleConfirmCreate () {
                if (!this.newDirectory.length) {
                    return false
                }
                try {
                    const data = await this.$store.dispatch('resource/directory/create', {
                        params: {
                            bk_module_name: this.newDirectory
                        },
                        config: {
                            requestId: this.request.create
                        }
                    })
                    this.directories.push({
                        bk_module_id: data.created.id,
                        bk_module_name: this.newDirectory
                    })
                    this.selected = data.created.id
                    this.toggleCreate(false)
                } catch (e) {
                    console.error(e)
                }
            },
            getInfo () {
                const directory = this.directories.find(directory => directory.bk_module_id === this.value)
                if (directory) {
                    return directory.bk_module_name
                }
                return null
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
