<template>
    <div class="search-box">
        <input id="search"
            class="search-keywords"
            type="text"
            maxlength="32"
            :placeholder="$t('请输入搜索内容')"
            v-model.trim="queryString"
            @keypress.enter="handleSearch">
        <div class="search-btn">
            <label class="bk-icon icon-search" for="search" @click="handleSearch"></label>
        </div>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                queryString: ''
            }
        },
        methods: {
            handleSearch () {
                if (!this.queryString) return
                this.$router.push({
                    path: '/index/search',
                    query: {
                        keywords: this.queryString,
                        from: this.$route.fullPath
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .search-box {
        position: relative;
        width: 50%;
        max-width: 640px;
        line-height: 40px;
        margin: 50px auto 38px;
        .search-keywords {
            width: 100%;
            height: 42px;
            border-radius: 2px;
            border: 1px solid #c3cdd7;
            padding: 0 56px 0 16px;
            font-size: 14px;
        }
        .search-btn {
            position: absolute;
            right: 0;
            top: 0;
            width: 50px;
            height: 42px;
            display: flex;
            justify-content: center;
            align-items: center;
            cursor: pointer;
            .icon-search {
                font-size: 18px;
                cursor: pointer;
            }
        }
    }
</style>
