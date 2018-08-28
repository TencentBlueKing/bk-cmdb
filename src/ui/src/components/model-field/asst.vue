<template>
    <div class="form-item">
        <div class="form-item">
            <label class="form-label">{{$t('ModelManagement["关联模型"]')}}</label>
            <bk-selector
                class="form-selector bk-selector-small"
                :list="list"
                :selected.sync="selected"
            ></bk-selector>
        </div>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                selected: '',
                list: [{
                    id: 1,
                    name: 2
                }, {
                    id: 11,
                    name: 22
                }]
            }
        }
    }
</script>
