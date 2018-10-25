<template>
    <div class="relation-type-content">
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["唯一标识"]')}}
                <span class="color-danger">*</span>
            </span>
            <input type="text" class="cmdb-form-input">
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('Hosts["名称"]')}}
                <span class="color-danger">*</span>
            </span>
            <input type="text" class="cmdb-form-input">
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["源->目标描述"]')}}
                <span class="color-danger">*</span>
            </span>
            <input type="text" class="cmdb-form-input">
        </label>
        <label class="form-label">
            <span class="label-text">{{$t('ModelManagement["目标描述->源"]')}}</span>
            <span class="color-danger">*</span>
            <input type="text" class="cmdb-form-input">
        </label>
        <div class="radio-box overflow">
            <label class="label-text">
                {{$t('ModelManagement["是否有方向"]')}}<span class="text-desc">({{$t('ModelManagement["仅视图"]')}})</span>
            </label>
            <label class="cmdb-form-radio cmdb-radio-small">
                <input type="radio" name="direction" value="src_to_dest" v-model="relationInfo.direction">
                <span class="cmdb-radio-text">{{$t('ModelManagement["是，源指向目标"]')}}</span>
            </label>
            <label class="cmdb-form-radio cmdb-radio-small">
                <input type="radio" name="direction" value="none" v-model="relationInfo.direction">
                <span class="cmdb-radio-text">{{$t('ModelManagement["否"]')}}</span>
            </label>
        </div>
        <div class="btn-group">
            <bk-button type="primary" :loading="$loading(['updateObjectAttribute', 'createObjectAttribute'])" @click="saveRelation">
                {{$t('ModelManagement["确定"]')}}
            </bk-button>
            <bk-button type="default" @click="cancel">
                {{$t('ModelManagement["取消"]')}}
            </bk-button>
        </div>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        props: {
            relation: {
                type: Object
            },
            isReadOnly: {
                type: Boolean,
                default: false
            },
            isEdit: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                relationInfo: {
                    bk_asst_id: '',
                    bk_asst_name: '',
                    src_des: '',
                    dest_des: '',
                    direction: 'none' // none, src_to_dest
                }
            }
        },
        created () {
            if (this.isEdit) {
                this.relationInfo = this.relation
            }
        },
        methods: {
            ...mapActions('objectAssociation', [
                'createAssociationType',
                'updateAssociationType'
            ]),
            async saveRelation () {
                if (this.isEdit) {
                    await this.updateAssociationType({
                        params: this.relationInfo,
                        config: {
                            requestId: 'updateAssociationType'
                        }
                    })
                } else {
                    await this.createAssociationType({
                        params: this.relationInfo,
                        config: {
                            requestId: 'createAssociationType'
                        }
                    })
                }
            },
            cancel () {

            }
        }
    }
</script>


<style lang="scss" scoped>
    .relation-type-content {
        .radio-box {
            .label-text {
                padding-right: 2px;
            }
        }
    }
</style>
