<template>
  <header v-if="!isEditor && !isPreview">
    <div>
      <button @click="openSidebar" :aria-label="$t('buttons.toggleSidebar')" :title="$t('buttons.toggleSidebar')" class="action">
        <i class="material-icons">menu</i>
      </button>
      <img :src="logoURL" alt="File Browser">
      <search v-if="isLogged"></search>
    </div>
    <div>
      <template v-if="isLogged || isSharing">
        <button v-show="!isSharing" @click="openSearch" :aria-label="$t('buttons.search')" :title="$t('buttons.search')" class="search-button action">
          <i class="material-icons">search</i>
        </button>

        <button @click="openMore" id="more" :aria-label="$t('buttons.more')" :title="$t('buttons.more')" class="action">
          <i class="material-icons">more_vert</i>
        </button>

        <!-- Menu that shows on listing AND mobile when there are files selected -->
        <div id="file-selection" v-if="isMobile && isListing && !isSharing">
          <span v-if="selectedCount > 0">{{ selectedCount }} selected</span>
          <playlist-button v-show="showPlaylistButton"></playlist-button>
          <share-button v-show="showShareButton"></share-button>
          <rename-button v-show="showRenameButton"></rename-button>
          <copy-button v-show="showCopyButton"></copy-button>
          <move-button v-show="showMoveButton"></move-button>
          <delete-button v-show="showDeleteButton"></delete-button>
        </div>

        <!-- This buttons are shown on a dropdown on mobile phones -->
        <div id="dropdown" :class="{ active: showMore }">
          <div v-if="!isListing || !isMobile">
            <playlist-button v-show="showPlaylistButton"></playlist-button>
            <share-button v-show="showShareButton"></share-button>
            <rename-button v-show="showRenameButton"></rename-button>
            <copy-button v-show="showCopyButton"></copy-button>
            <move-button v-show="showMoveButton"></move-button>
            <delete-button v-show="showDeleteButton"></delete-button>
          </div>

          <shell-button v-if="isExecEnabled && !isSharing && user.perm.execute" />
          <switch-button v-show="isListing"></switch-button>
          <download-button v-show="showDownloadButton"></download-button>
          <upload-button v-show="showUpload"></upload-button>
          <info-button v-show="isFiles"></info-button>

          <button v-show="isListing || (isSharing && req.isDir)" @click="toggleMultipleSelection" :aria-label="$t('buttons.selectMultiple')" :title="$t('buttons.selectMultiple')" class="action" >
            <i class="material-icons">check_circle</i>
            <span>{{ $t('buttons.select') }}</span>
          </button>
        </div>
      </template>

      <div v-show="showOverlay" @click="resetPrompts" class="overlay"></div>
    </div>
  </header>
</template>

<script>
import Search from './Search'
import InfoButton from './buttons/Info'
import DeleteButton from './buttons/Delete'
import RenameButton from './buttons/Rename'
import UploadButton from './buttons/Upload'
import DownloadButton from './buttons/Download'
import SwitchButton from './buttons/SwitchView'
import MoveButton from './buttons/Move'
import CopyButton from './buttons/Copy'
import ShareButton from './buttons/Share'
import ShellButton from './buttons/Shell'
import PlaylistButton from './buttons/Playlist'

import {mapGetters, mapState} from 'vuex'
import { logoURL, enableExec } from '@/utils/constants'
import * as api from '@/api'
import buttons from '@/utils/buttons'

export default {
  name: 'header-layout',
  components: {
    Search,
    InfoButton,
    DeleteButton,
    ShareButton,
    RenameButton,
    DownloadButton,
    CopyButton,
    UploadButton,
    SwitchButton,
    MoveButton,
    ShellButton,
    PlaylistButton
  },
  data: function () {
    return {
      width: window.innerWidth,
      pluginData: {
        api,
        buttons,
        'store': this.$store,
        'router': this.$router
      }
    }
  },
  created () {
    window.addEventListener('resize', () => {
      this.width = window.innerWidth
    })
  },
  computed: {
    ...mapGetters([
      'selectedCount',
      'isFiles',
      'isEditor',
      'isPreview',
      'isListing',
      'isLogged',
      'isSharing'
    ]),
    ...mapState([
      'req',
      'user',
      'loading',
      'reload',
      'multiple'
    ]),
    logoURL: () => logoURL,
    isExecEnabled: () => enableExec,
    isMobile () {
      return this.width <= 736
    },
    showUpload () {
      return this.isListing && this.user.perm.create
    },
    showDownloadButton () {
      return (this.isFiles && this.user.perm.download) || (this.isSharing && this.selectedCount > 0)
    },
    showDeleteButton () {
      return this.isFiles && (this.isListing
        ? (this.selectedCount !== 0 && this.user.perm.delete)
        : this.user.perm.delete)
    },
    showRenameButton () {
      return this.isFiles && (this.isListing
        ? (this.selectedCount === 1 && this.user.perm.rename)
        : this.user.perm.rename)
    },
    showShareButton () {
      return this.isFiles && (this.isListing
        ? (this.selectedCount === 1 && this.user.perm.share)
        : this.user.perm.share)
    },
    showMoveButton () {
      return this.isFiles && (this.isListing
        ? (this.selectedCount > 0 && this.user.perm.rename)
        : this.user.perm.rename)
    },
    showPlaylistButton () {
      return this.selectedCount > 0
    },
    showCopyButton () {
      return this.isFiles && (this.isListing
        ? (this.selectedCount > 0 && this.user.perm.create)
        : this.user.perm.create)
    },
    showMore () {
      return (this.isFiles || this.isSharing) && this.$store.state.show === 'more'
    },
    showOverlay () {
      return this.showMore
    }
  },
  methods: {
    openSidebar () {
      this.$store.commit('showHover', 'sidebar')
    },
    openMore () {
      this.$store.commit('showHover', 'more')
    },
    openSearch () {
      this.$store.commit('showHover', 'search')
    },
    toggleMultipleSelection () {
      this.$store.commit('multiple', !this.multiple)
      this.resetPrompts()
    },
    resetPrompts () {
      this.$store.commit('closeHovers')
    }
  }
}
</script>
