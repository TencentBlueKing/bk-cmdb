<template>
    <div class="host-search-layout">
        <div class="search-bar">
            <bk-input class="search-input"
                ref="searchInput"
                type="textarea"
                :placeholder="$t('请输入IP，多个IP请使用换行分隔')"
                :rows="rows"
                v-model="searchContent"
                @focus="handleFocus"
                @blur="handleBlur">
            </bk-input>
            <bk-button theme="primary" class="search-btn" @click="handleSearch">
                <i class="bk-icon icon-search"></i>
                {{$t('搜索')}}
            </bk-button>
            <span v-if="showEllipsis" class="search-text" @click="handleSearchInput">{{searchText}}</span>
        </div>
    </div>
</template>

<script>
    import { MENU_RESOURCE_HOST } from '@/dictionary/menu-symbol'
    export default {
        data () {
            return {
                rows: 1,
                searchText: '',
                searchContent: '',
                textarea: '',
                showEllipsis: false,
                textareaDom: null
            }
        },
        watch: {
            searchContent () {
                this.$nextTick(this.setRows)
            }
        },
        mounted () {
            this.textareaDom = this.$refs.searchInput && this.$refs.searchInput.$refs.textarea
        },
        methods: {
            setRows () {
                const rows = this.searchContent.split('\n').length || 1
                this.rows = Math.min(10, rows)
            },
            handleFocus () {
                this.$emit('focus', true)
                this.setRows()
            },
            handleBlur () {
                this.textareaDom && this.textareaDom.blur()
                this.$emit('focus', false)
                const data = this.searchContent.split('\n').map(text => text.trim()).filter(text => text)
                if (data.length) {
                    this.showEllipsis = true
                    this.searchText = data.join(',')
                } else {
                    this.searchContent = ''
                }
                this.$nextTick(() => {
                    this.rows = 1
                    this.textareaDom && (this.textareaDom.scrollTop = 0)
                })
            },
            handleSearchInput () {
                this.showEllipsis = false
                this.textareaDom && this.textareaDom.focus()
            },
            async handleSearch () {
                const searchList = this.searchContent.split('\n').map(text => text.trim()).filter(text => text)
                if (searchList.length > 500) {
                    this.$warn(this.$t('目前最多支持搜索500个IP'))
                } else if (searchList.length) {
                    try {
                        const validateQueue = searchList.map(text => this.$validator.verify(text, 'ip'))
                        const results = await Promise.all(validateQueue)
                        const isPassValidate = results.every(res => res.valid)
                        if (!isPassValidate) {
                            this.$warn(this.$t('请输入完整IP进行搜索，多个IP用换行分割'))
                            return
                        }
                        this.$store.commit('hosts/setFilterIP', {
                            text: searchList.join('\n'),
                            exact: true
                        })
                        this.$store.commit('hosts/setIsHostSearch', true)
                        this.$router.push({ name: MENU_RESOURCE_HOST })
                    } catch (e) {
                        console.error(e)
                    }
                } else {
                    this.searchContent = ''
                    this.textareaDom && this.textareaDom.focus()
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .host-search-layout {
        position: relative;
        width: 100%;
        max-width: 726px;
        height: 42px;
        margin: 0 auto;
    }
    .search-bar {
        position: absolute;
        width: 100%;
        height: 42px;
        z-index: 999;
        display: flex;
    }
    .search-input {
        flex: 1;
        /deep/ {
            .bk-textarea-wrapper {
                border: 0;
                border-radius: 2px 0 0 2px;
            }
            .bk-form-textarea {
                min-height: 42px;
                line-height: 30px;
                font-size: 14px;
                border: 1px solid #C4C6CC;
                padding: 5px 16px;
                border-radius: 2px 0 0 2px;
            }
        }
    }
    .search-btn {
        width: 86px;
        height: 42px;
        line-height: 42px;
        padding: 0;
        border-radius: 0 2px 2px 0;
        .icon-search {
            width: 18px;
            height: 18px;
            font-size: 18px;
            margin: -2px 4px 0 0;
        }
    }
    .search-text {
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        max-width: 640px;
        height: 42px;
        line-height: 30px;
        font-size: 14px;
        color: #63656E;
        border: 1px solid #C4C6CC;
        background-color: #FFFFFF;
        padding: 5px 16px;
        z-index: 1;
        cursor: text;
        @include ellipsis;
    }
</style>
