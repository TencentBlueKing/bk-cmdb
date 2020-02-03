<template>
    <bk-form form-type="vertical" :label-width="300">
        <bk-form-item :label="$t('资源录入到')" required>
            <cloud-resource-folder-selector class="selector"
                :data-vv-as="$t('资源目录')"
                data-vv-name="folder"
                v-validate="'required'"
                v-model="selected">
            </cloud-resource-folder-selector>
            <p class="form-error" v-if="errors.has('folder')">{{errors.first('folder')}}</p>
        </bk-form-item>
    </bk-form>
</template>

<script>
    import CloudResourceFolderSelector from './resource-folder-selector.vue'
    export default {
        name: 'cloud-resource-form-all',
        components: {
            CloudResourceFolderSelector
        },
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
