<template>
  <div class="3d-preview">
    <div class="loading" v-if="loadingPreview">
      <div class="spinner">
        <div class="bounce1"></div>
        <div class="bounce2"></div>
        <div class="bounce3"></div>
      </div>
    </div>
    <model-obj v-if="req.extension.match(/\.obj$/i)" :src="raw" :backgroundAlpha="0" :rotation="rotation" @on-mousedown="onMousedown" @on-load="onLoad" @on-error="onError"></model-obj>
    <model-stl v-else-if="req.extension.match(/\.stl$/i)" :src="raw" :backgroundAlpha="0" :rotation="rotation" @on-mousedown="onMousedown" @on-load="onLoad" @on-error="onError"></model-stl>
    <model-ply v-else-if="req.extension.match(/\.ply$/i)" :src="raw" :backgroundAlpha="0" :rotation="rotation" @on-mousedown="onMousedown" @on-load="onLoad" @on-error="onError"></model-ply>
    <model-fbx v-else-if="req.extension.match(/\.fbx$/i)" :src="raw" :backgroundAlpha="0" :rotation="rotation" @on-mousedown="onMousedown" @on-load="onLoad" @on-error="onError"></model-fbx>
    <model-gltf v-else-if="req.extension.match(/\.gltf$/i)" :src="raw" :backgroundAlpha="0" :rotation="rotation" @on-mousedown="onMousedown" @on-load="onLoad" @on-error="onError"></model-gltf>
    <model-collada v-else-if="req.extension.match(/\.dae$/i)" :src="raw" :backgroundAlpha="0" :rotation="rotation" @on-mousedown="onMousedown" @on-load="onLoad" @on-error="onError"></model-collada>
  </div>
</template>
<script>

import { ModelCollada, ModelFbx, ModelGltf, ModelObj, ModelPly, ModelStl } from 'vue-3d-model';
import { mapState } from 'vuex'
import url from '@/utils/url'
import { baseURL } from '@/utils/constants'

export default {
  components: {
    ModelCollada,
    ModelFbx,
    ModelGltf,
    ModelObj,
    ModelPly,
    ModelStl
  },
  data() {
    return {
      loadingPreview: false,
      rotating: false,
      rotation: {
        x: -Math.PI / 2,
        y: 0,
        z: 0
      }
    }
  },
  computed: {
    ...mapState(['req', 'jwt']),
    downloadUrl () {
      return `${baseURL}/api/raw${url.encodePath(this.req.path)}?auth=${this.jwt}`
    },
    raw () {
      return `${this.downloadUrl}&inline=true`
    },
  },
  mounted() {
    this.loadingPreview = true
    this.rotating = true
  },
  methods: {
    onLoad () {
      this.loadingPreview = false
      this.rotate();
    },
    onError () {
      this.loadingPreview = false
    },
    onMousedown () {
      this.rotating = false
    },
    rotate () {
      this.rotation.z += 0.01;
      if (this.rotating) requestAnimationFrame( this.rotate );
    }
  }
}
</script>
