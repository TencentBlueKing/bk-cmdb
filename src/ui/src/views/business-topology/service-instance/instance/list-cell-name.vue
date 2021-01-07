<template>
    <div :class="{ 'instance-name': true, disabled }">
        <template v-if="!editing">
            <span class="name-text" v-bk-overflow-tips>{{row.name}}</span>
            <cmdb-auth tag="i" class="name-edit icon-cc-edit-shape"
                :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }"
                @click.native.stop
                @click="handleEdit">
            </cmdb-auth>
        </template>
        <service-instance-name-edit-form v-else ref="nameEditForm"
            :value="row.name"
            @click.native.stop
            @confirm="handleConfirm"
            @cancel="handleCancel" />
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import ServiceInstanceNameEditForm from '@/components/service/instance-name-edit-form'
    export default {
        name: 'list-cell-name',
        components: {
            ServiceInstanceNameEditForm
        },
        props: {
            row: Object
        },
        data () {
            return {
                request: {
                    update: Symbol('update')
                }
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            disabled () {
                return !this.row.process_count
            },
            editing () {
                return this.row.editing.name
            }
        },
        methods: {
            handleEdit () {
                this.$emit('edit')
                this.$nextTick(() => {
                    this.$refs.nameEditForm.focus()
                })
            },
            async handleConfirm (value) {
                try {
                    await this.$store.dispatch('serviceInstance/updateServiceInstance', {
                        bizId: this.bizId,
                        params: {
                            data: [{
                                service_instance_id: this.row.id,
                                update: {
                                    name: value
                                }
                            }]
                        },
                        config: { requestId: this.request.update }
                    })
                    this.$emit('success', value)
                } catch (error) {
                    console.error(error)
                }
            },
            handleCancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .instance-name {
        display: flex;
        width: 100%;
        align-items: center;
        &.disabled {
            color: $textDisabledColor;
        }
        .name-text {
            max-width: calc(100% - 25px);
            font-weight: bold;
            @include ellipsis;
            cursor: default;
        }
        .name-edit {
            visibility: hidden;
            font-size: 14px;
            height: 26px;
            width: 26px;
            text-align: center;
            line-height: 26px;
            color: $primaryColor;
            cursor: pointer;
            &:hover {
                opacity: .8;
            }
            &.disabled {
                color: $textDisabledColor;
            }
        }
        .instance-name-edit-form {
            flex: auto;
        }
    }
</style>
