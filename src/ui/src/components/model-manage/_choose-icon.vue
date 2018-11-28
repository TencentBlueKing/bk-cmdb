<template>
    <div class="model-icon-list">
        <div class="page clearfix">
            <input type="text" class="cmdb-form-input" :placeholder="$t('ModelManagement[\'请输入关键词\']')" v-model.trim="searchText">
            <div class="page-btn">
                <bk-button type="default" :disabled="!page.current" @click="pageTurning(--page.current)">
                    <i class="bk-icon icon-angle-left"></i>
                </bk-button>
                <bk-button type="default" :disabled="page.current === page.totalPage - 1" @click="pageTurning(++page.current)">
                    <i class="bk-icon icon-angle-right"></i>
                </bk-button>
            </div>
        </div>
        <ul class="icon-box clearfix">
            <li class="icon" 
                :class="{'create': type === 'create', 'active': icon.value === localValue}"
                v-tooltip="{content: language === 'zh_CN' ? icon.nameZh : icon.nameEn}"
                v-for="(icon, index) in curIconList" 
                :key="index" @click="chooseIcon(icon.value)">
                <i :class="icon.value"></i>
            </li>
        </ul>
    </div>
</template>

<script>
    import iconList from '@/assets/json/model-icon.json'
    import { mapGetters } from 'vuex'
    export default {
        props: {
            value: {
                default: 'icon-cc-default'
            },
            type: {
                default: 'create'
            }
        },
        data () {
            return {
                iconList,
                localValue: this.value,
                searchText: '',
                page: {
                    current: 0,
                    size: 28,
                    totalPage: Math.ceil(iconList.length / 28)
                }
            }
        },
        computed: {
            ...mapGetters([
                'language'
            ]),
            curIconList () {
                let {
                    searchText,
                    page
                } = this
                let curIconList = this.iconList
                if (searchText.length) {
                    curIconList = this.iconList.filter(icon => {
                        return icon.nameZh.toLowerCase().indexOf(searchText.toLowerCase()) > -1 || icon.nameEn.toLowerCase().indexOf(searchText.toLowerCase()) > -1
                    })
                }
                this.page.totalPage = Math.ceil(curIconList.length / page.size)
                return curIconList.slice(page.size * page.current, page.size * (page.current + 1))
            }
        },
        watch: {
            searchText () {
                this.page.current = 0
            }
        },
        methods: {
            chooseIcon (value) {
                this.localValue = value
                this.$emit('input', value)
                this.$emit('chooseIcon')
            },
            pageTurning (page) {
                this.page.current = page
            }
        }
    }
</script>

<style lang="scss" scoped>
    .model-icon-list {
        display: block;
        height: 100%;
    }
    .page {
        padding: 15px;
        .cmdb-form-input {
            float: left;
            width: 220px;
            height: 30px;
            line-height: 28px;
        }
        .page-btn {
            float: right;
            .bk-button {
                padding: 0;
                width: 30px;
                height: 30px;
                line-height: 1;
                vertical-align: middle;
            }
        }
    }
    .icon-box {
        padding: 0 15px 10px;
        width: 100%;
        height: calc(100% - 60px);
        .icon {
            float: left;
            width: calc(100% / 7);
            height: 46px;
            padding: 5px;
            font-size: 24px;
            text-align: center;
            cursor: pointer;
            &.create {
                font-size: 30px;
                margin-bottom: 10px;
                height: 50px;
            }
            &:hover,
            &.active {
                background: #e2efff;
                color: #3c96ff;
            }
        }
        .page {
            height: 52px;
            padding: 10px 20px;
            .cmdb-form-input {
                float: left;
                width: 200px;
                height: 30px;
            }
            .page-btn {
                float: right;
                .bk-button {
                    padding: 0;
                    width: 30px;
                    height: 30px;
                    line-height: 1;
                    vertical-align: top;
                }
            }
        }
    }
</style>
