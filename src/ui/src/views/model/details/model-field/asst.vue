<template>
    <div class="form-item">
        <div class="form-item">
            <label class="form-label">{{$t('ModelManagement["关联模型"]')}}</label>
            <div class="input-box">
                <bk-selector
                    class="form-selector"
                    :list="asstList"
                    :has-children='true'
                    :selected.sync="selected"
                    @item-selected="handleChange"
                    :disabled="isReadOnly || isEditField"
                    :content-max-height="200"
                ></bk-selector>
                <input type="text" hidden 
                    v-model="selected"
                    v-validate="'required'"
                    name="asst">
                <span v-show="errors.has('asst')" class="error-msg color-danger">{{ errors.first('asst') }}</span>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        props: {
            value: {
                default: ''
            },
            isEditField: {
                type: Boolean,
                default: false
            },
            isReadOnly: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                selected: ''
            }
        },
        computed: {
            ...mapGetters('objectModelClassify', [
                'classifications'
            ]),
            ...mapGetters('objectModel', [
                'activeModel'
            ]),
            asstList () {
                let asstList = []
                this.classifications.map(classify => {
                    if (classify['bk_objects'].length && classify['bk_classification_id'] !== 'bk_biz_topo') {
                        let objects = []
                        classify['bk_objects'].map(({bk_obj_id: objId, bk_obj_name: objName}) => {
                            if (['plat', 'process', 'set', 'module', this.activeModel['bk_obj_id']].indexOf(objId) === -1) {
                                objects.push({
                                    id: objId,
                                    name: objName
                                })
                            }
                        })
                        if (objects.length) {
                            asstList.push({
                                name: classify['bk_classification_name'],
                                children: objects
                            })
                        }
                    }
                })
                return asstList
            }
        },
        watch: {
            value () {
                this.initValue()
            }
        },
        created () {
            this.initValue()
        },
        methods: {
            handleChange () {
                this.$emit('input', this.selected)
            },
            initValue () {
                this.selected = this.value
            },
            validate () {
                return this.$validator.validateAll()
            }
        }
    }
</script>
