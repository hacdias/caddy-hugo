<template>
  <button @click="getPlaylist" :aria-label="$t('buttons.playlist')" :title="$t('buttons.playlist')" class="action">
    <i class="material-icons">playlist_play</i>
    <span>{{ $t('buttons.playlist') }}</span>
  </button>
</template>

<script>
import { saveAs } from 'file-saver'
import { mapState } from 'vuex'
import { baseURL } from '@/utils/constants'
import url from '@/utils/url'

export default {
  name: 'playlist-button',
  computed: {
    ...mapState([ 'req', 'selected', 'selectedCount', 'jwt', ]),
  },
  methods: {
    raw (fileUrl) {
      return `${baseURL}/api/raw${url.encodePath(fileUrl)}?auth=${this.jwt}&inline=true`
    },
    getPlaylist () {
      const { protocol, hostname, port } = window.location
      let host = `${protocol}//${hostname}`
      if (port) host = `${host}:${port}`

      const validTypes = ['video', 'audio']
      console.log(this.req.items)
      // constraint to video & audio files remove any folders
      const files = this.selected.map((i) => this.req.items[i]).filter((v) => validTypes.indexOf(v.type) !== -1 && !v.isDir)
      console.log(files)
      if (files.length === 0) return

      const urls = files.map((v) => `${host}${this.raw(v.path)}\n`)
      var blob = new Blob(["#EXTM3U\n", ...urls], {type: "text/plain;charset=utf-8"})
      saveAs(blob, "playlist.m3u")
    }
  }
}
</script>
