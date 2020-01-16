<template>
    <bk-form form-type="vertical" :label-width="300">
        <bk-form-item :label="$t('资源录入到')" required>
            <bk-select class="selector"
                searchable
                :placeholder="$t('请选择xx', { name: $t('资源目录') })"
                :data-vv-as="$t('资源目录')"
                data-vv-name="folder"
                v-validate="'required'"
                v-model="selected">
                <bk-option v-for="folder in folders"
                    :key="folder.id"
                    :id="folder.id"
                    :name="folder.name">
                </bk-option>
                <a href="javascript:void(0)" class="extension-link" slot="extension">
                    <i class="bk-icon icon-plus-circle"></i>
                    {{$t('申请其他目录权限')}}
                </a>
            </bk-select>
            <p class="form-error" v-if="errors.has('folder')">{{errors.first('folder')}}</p>
        </bk-form-item>
    </bk-form>
</template>

<script>
    export default {
        name: 'cloud-resource-form-all',
        data () {
            return {
                selected: '',
                folders: []
            }
        }
    }
</script>

<style lang="scss" scoped>
    .selector {
        width: 300px;
    }
</style>
