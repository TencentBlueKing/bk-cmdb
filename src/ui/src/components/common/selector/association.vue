<template>
    <bk-select :selected.sync="localSelected" :disabled="disabled || localDisabled" :multiple="multiple" 
    :filterable="true" 
    @on-selected="handleSelected">
        <bk-select-option v-for="(option, index) in associateInst"
            :key="index"
            :label="asstObjId === 'biz'? option['bk_biz_name']: option['bk_inst_name']"
            :value="asstObjId === 'biz'? option['bk_biz_id']: option['bk_inst_id']">
        </bk-select-option>
    </bk-select>
</template>
<script>
    import { mapGetters } from 'vuex'
    export default {
        props: {
            disabled: {
                type: Boolean,
                default: false
            },
            multiple: {
                type: Boolean,
                default: false
            },
            asstObjId: {
                type: String,
                required: true
            },
            selected: {
                required: true,
                validator: (selected) => {
                    return typeof selected === 'string' || Array.isArray(selected) || typeof selected === 'undefined' || !isNaN(selected)
                }
            }
        },
        data () {
            return {
                localSelected: '',
                localDisabled: false,
                associateInst: [],
                isLoaded: false
            }
        },
        computed: {
            ...mapGetters(['bkSupplierAccount'])
        },
        watch: {
            selected (selected) {
                if (this.isLoaded && selected !== this.localSelected) {
                    this.setLocalSelected()
                }
            },
            localSelected (localSelected) {
                this.$emit('update:selected', localSelected)
            }
        },
        async created () {
            await this.getAssociateInst()
            this.isLoaded = true
            this.setLocalSelected()
        },
        methods: {
            getAssociateInst () {
                let url = this.asstObjId === 'biz' ? `biz/search/${this.bkSupplierAccount}` : `inst/search/${this.bkSupplierAccount}/${this.asstObjId}`
                return this.$axios.post(url, {}, {globalError: false}).then(res => {
                    if (res.result) {
                        this.associateInst = res.data.info || []
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                }).catch(error => {
                    if (error.response.status === 403) {
                        this.localDisabled = true
                    }
                })
            },
            setLocalSelected () {
                if (Array.isArray(this.selected)) {
                    this.localSelected = this.selected.map(({id, bk_inst_id: bkInstId}) => {
                        return id === '' ? '' : bkInstId
                    }).join(',')
                } else if (typeof this.selected === 'undefined') {
                    this.localSelected = ''
                } else {
                    this.localSelected = this.selected
                }
                this.$emit('update:selected', this.localSelected)
            },
            handleSelected () {
                this.$emit('on-selected', ...arguments)
            }
        }
    }
</script>