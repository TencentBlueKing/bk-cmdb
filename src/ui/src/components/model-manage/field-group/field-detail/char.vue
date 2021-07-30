<template>
  <div class="form-label cmdb-form-item" :class="{ 'is-error': errors.has('option') }">
    <span class="label-text">{{$t('正则校验')}}</span>
    <textarea
      class="raw"
      name="option"
      v-model="localValue"
      :disabled="isReadOnly"
      data-vv-validate-on="blur"
      v-validate="'remoteRegular'"
      @input="handleInput">
    </textarea>
    <p class="form-error">{{errors.first('option')}}</p>
  </div>
</template>

<script>
  export default {
    props: {
      value: {
        type: String,
        default: ''
      },
      isReadOnly: {
        type: Boolean,
        default: false
      }
    },
    data() {
      return {
        localValue: ''
      }
    },
    watch: {
      value() {
        this.localValue = this.value === '' ? '' : this.value
      }
    },
    created() {
      this.localValue = this.value === '' ? '' : this.value
    },
    methods: {
      handleInput() {
        this.$emit('input', this.localValue)
      },
      validate() {
        return this.$validator.validateAll()
      }
    }
  }
</script>
