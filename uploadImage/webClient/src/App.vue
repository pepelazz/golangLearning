<template>
  <div id="q-app">
    <q-layout view="hHh Lpr lFf">

      <q-header reveal elevated class="bg-white text-grey-8 q-py-xs">
        <q-toolbar>
          <q-btn dense flat round icon="menu"/>
          <q-btn flat no-caps no-wrap class="q-ml-xs" @click="$router.push('/')">
            <q-toolbar-title shrink class="text-weight-bold">
              Upload Image Demo
            </q-toolbar-title>
          </q-btn>
        </q-toolbar>
      </q-header>

      <q-page-container>
        <q-page padding>
          <div class="row q-mb-lg">
            <div class="col">
              <upload-img :resize="true" @update="url => images.unshift({url})"/>
            </div>
          </div>
          <div class="row q-col-gutter-md">
            <div v-for="img in images" :key="img.url" class="col-2">
              <q-img :src="img.url">
                <div class="absolute-top text-subtitle1 text-center">
                  {{img.size}}
                </div>
              </q-img>
            </div>
          </div>
        </q-page>
      </q-page-container>
    </q-layout>
  </div>
</template>

<script>
    import uploadImg from './components/uploadImg'
    import config from './config'

    export default {
        components: {uploadImg},
        data() {
            return {
                images: [],
            }
        },
        mounted() {
            // при открытии страницы загружаем список фотографий с сервера
            fetch(`${config.apiUrl()}/get_all_image`, {method: 'POST'})
                .then(res => res.json())
                .then(res => {
                    this.images = res.result.map(img => {
                        img.url = `${config.apiUrl()}/${img.url}`
                        img.size = (img.size / 1000).toFixed(0) + ' Kb'
                        return img
                    })
                })
        }
    }
</script>
