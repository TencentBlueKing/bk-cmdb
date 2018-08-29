<template>
    <div class="create-field-wrapper">
        <div class="title" @click="closeBox">
            <img src="../../../assets/images/icon/down_icon.png" alt="">
        </div>
        <div class="field-content">
            <v-model-field
                :isEditField="false"
            ></v-model-field>
        </div>
    </div>
</template>

<script>
    import vModelField from '@/components/model-field'
    import { mapActions } from 'vuex'
    export default {
        components: {
            vModelField
        },
        data () {
            return {
                
            }
        },
        methods: {
            ...mapActions('objectModelProperty', [
                'createObjectAttribute'
            ]),
            closeBox () {
                this.$emit('closeBox')
            },
            save () {

            }
        }
    }
</script>


<style lang="scss" scoped>
    .create-field-wrapper {
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: #fff;
        .title {
            width: 100%;
            height: 40px;
            background: linear-gradient(to bottom, #f9f9f9, #fff);
            cursor: pointer;
            text-align: center;
            img {
                margin-top: 14px;
            }
        }
        .field-content {
            padding: 0 40px;
        }
    }
</style>
