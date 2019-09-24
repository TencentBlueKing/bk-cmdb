<template>
    <section class="layout">
        <div class="layout-row">
            <label class="row-label inline-block-middle" :title="$t('模板名称')">{{$t('模板名称')}}</label>
            <i class="row-required">*</i>
            <bk-input class="row-content"
                data-vv-name="name"
                v-validate="'required|singlechar|length:20'"
                :placeholder="$t('集群模板名称占位符')">
            </bk-input>
            <p class="row-error" v-if="errors.has('name')">{{errors.first('name')}}</p>
        </div>
    </section>
</template>

<script>
    export default {
        data () {
            return {}
        }
    }
</script>

<style lang="scss" scoped>
    .layout {
        color: #63656E;
    }
    .layout-row {
        position: relative;
        .row-label {
            width: 125px;
            text-align: right;
            font-size: 14px;
            line-height: 32px;
            padding: 0 0 0 20px;
            @include ellipsis;
        }
        .row-required {
            position: relative;
            top: 4px;
            font-style: normal;
            color: #EA3636;
        }
        .row-content {
            display: inline-block;
            vertical-align: top;
            width: 520px;
        }
        .row-error {
            position: absolute;
            color: #EA3636;
            padding-left: 145px;
            font-size: 12px;
            top: 100%;
        }
    }
</style>
