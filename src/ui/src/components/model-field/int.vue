<template>
    <div>
        <div class="form-item has-right-content">
            <label class="form-label">{{$t('ModelManagement["最小值"]')}}</label>
            <div class="input-box">
                <input type="text" class="cmdb-form-input" placeholder=""
                    v-model="min"
                    v-validate="`number`"
                    maxlength="11"
                    :name="'min'">
                <span v-show="errors.has('min')" class="error-msg color-danger">{{ errors.first('min') }}</span>
            </div>
        </div>
        <div class="form-item">
            <label class="form-label">{{$t('ModelManagement["最大值"]')}}</label>
            <div class="input-box">
                <input type="text" class="cmdb-form-input" placeholder="" v-model="max"
                    name="max"
                    v-validate="`number|isBigger:${min}`">
                <span v-show="errors.has('max')" class="error-msg color-danger">{{ errors.first('max') }}</span>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                max: '',
                min: ''
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-label {
        vertical-align: top;
        line-height: 30px;
    }
</style>
