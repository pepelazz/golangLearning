<template>
  <div class="row q-gutter-md">
    <div class="col-auto">
      <q-uploader
        ref="uploader"
        label="Выберите файл для загрузки"
        auto-upload
        accept="jpeg, jpg, png, gif"
        :url="uploadUrl"
        @uploaded='uploaded'
        @failed='failed'
        :form-fields="formFields"
      />
    </div>
    <div class="col-6" v-if="imgHref">
      <q-img :src="imgHref">
        <div class="absolute-top text-subtitle1 text-center">
          {{imgHref}}
        </div>
      </q-img>
    </div>
  </div>
</template>

<script>
    import config from '../config'

    export default {
        props: ['resize'],
        computed: {
            uploadUrl() {
                return this.resize ? `${config.apiUrl()}/upload_image_resize` : `${config.apiUrl()}/upload_image`
            }
        },
        data() {
            return {
                formFields: [{name: 'product_id', value: 12}],
                imgHref: null,
            }
        },
        methods: {
            uploaded({xhr: {response}}) {
                const res = JSON.parse(response)
                if (!res.ok) {
                    this.$q.notify({
                        color: 'negative',
                        position: 'bottom',
                        message: res.message,
                    })
                } else {
                    this.$refs.uploader.reset()
                    this.imgHref = `${config.apiUrl()}${res.result.file}`
                    this.$emit('update', this.imgHref)
                }
            },
            failed(msg) {
                this.$q.notify({
                    color: 'negative',
                    position: 'bottom',
                    message: 'ошибка загрузки',
                })
            }
        }
    }
</script>
