<template>
    <div>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["唯一标识"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('asstId')}">
                <input type="text" class="cmdb-form-input"
                name="asstId"
                :disabled="isEdit"
                v-model.trim="relationInfo['bk_asst_id']"
                v-validate="'required|singlechar'"
                :placeholder="$t('ModelManagement[\'请输入英文标识\']')">
                <p class="form-error">{{errors.first('asstId')}}</p>
            </div>
        </label>
    </div>
</template>
