<template>
  <div>
    <div id="breadcrumbs" v-if="isListing || error">
      <router-link to="/files/" :aria-label="$t('files.home')" :title="$t('files.home')">
        <i class="material-icons">home</i>
      </router-link>

      <span v-for="(link, index) in breadcrumbs" :key="index">
        <span class="chevron"><i class="material-icons">keyboard_arrow_right</i></span>
        <router-link :to="link.url">{{ link.name }}</router-link>
      </span>
    </div>

    <div v-if="error">
      <not-found v-if="error.message === '404'"></not-found>
      <forbidden v-else-if="error.message === '403'"></forbidden>
      <internal-error v-else></internal-error>
    </div>
    <preview v-else-if="isPreview"></preview>
    <editor v-else-if="isEditor"></editor>
    <listing :class="{ multiple }" v-else-if="isListing"></listing>
    <div v-else>
      <h2 class="message">
        <span>{{ $t('files.loading') }}</span>
      </h2>
    </div>
  </div>
</template>

<script>
import Forbidden from './errors/403'
import NotFound from './errors/404'
import InternalError from './errors/500'
import Preview from '@/components/files/Preview'
import Listing from '@/components/files/Listing'
import { files as api } from '@/api'
import { mapGetters, mapState, mapMutations } from 'vuex'

function clean (path) {
  return path.endsWith('/') ? path.slice(0, -1) : path
}

export default {
  name: 'files',
  components: {
    Forbidden,
    NotFound,
    InternalError,
    Preview,
    Listing,
    Editor: () => import('@/components/files/Editor')
  },
  computed: {
    ...mapGetters([
      'selectedCount',
      'isListing',
      'isEditor',
      'isFiles',
      'getLastViewedDetail'
    ]),
    ...mapState([
      'req',
      'user',
      'reload',
      'multiple',
      'loading',
      'show'
    ]),
    isPreview () {
      return !this.loading && !this.isListing && !this.isEditor || this.loading && this.$store.state.previewMode
    },
    breadcrumbs () {
      let parts = this.$route.path.split('/')

      if (parts[0] === '') {
        parts.shift()
      }

      if (parts[parts.length - 1] === '') {
        parts.pop()
      }

      let breadcrumbs = []

      for (let i = 0; i < parts.length; i++) {
        if (i === 0) {
          breadcrumbs.push({ name: decodeURIComponent(parts[i]), url: '/' + parts[i] + '/' })
        } else {
          breadcrumbs.push({ name: decodeURIComponent(parts[i]), url: breadcrumbs[i - 1].url + parts[i] + '/' })
        }
      }

      breadcrumbs.shift()

      if (breadcrumbs.length > 3) {
        while (breadcrumbs.length !== 4) {
          breadcrumbs.shift()
        }

        breadcrumbs[0].name = '...'
      }

      return breadcrumbs
    }
  },
  data: function () {
    return {
      error: null
    }
  },
  created () {
    this.fetchData()
  },
  watch: {
    '$route': 'nav',
    'reload': function () {
      this.fetchData()
    }
  },
  mounted () {
    window.addEventListener('keydown', this.keyEvent)
    window.addEventListener('scroll', this.scroll)
  },
  beforeDestroy () {
    window.removeEventListener('keydown', this.keyEvent)
    window.removeEventListener('scroll', this.scroll)
  },
  destroyed () {
    this.$store.commit('updateRequest', {})
  },
  methods: {
    ...mapMutations([ 'setLoading' ]),
    async nav() {
      if (!this.isPreview && !this.isEditor && clean(`/${this.$route.params.pathMatch}`).startsWith(clean(this.req.path))) {
        let dirs = clean(this.$route.fullPath).split("/")

        this.$store.commit('addLastViewed', {
          path: clean(this.req.path),
          detail: {
            clicked: decodeURIComponent(dirs[dirs.length - 1]),
            pageOffset: window.pageYOffset
          }
        })
      }
      await this.fetchData()
      let detail = this.getLastViewedDetail(clean(`/${this.$route.params.pathMatch}`))
      if (detail !== null) {
        let offsetTarget = Math.min(1000, detail.pageOffset), oldPageOffset = 0, pageOffset = window.pageYOffset
        let interval = setInterval(function () {
          window.scrollTo(0, offsetTarget)
          if (offsetTarget >= detail.pageOffset || oldPageOffset === pageOffset) clearInterval(interval);
          oldPageOffset = pageOffset
          pageOffset = window.pageYOffset
          offsetTarget += Math.min(1000, detail.pageOffset - offsetTarget)
        }, 20);

        for (let i = 0; i < this.req.items.length; i++) {
          if (this.req.items[i].name === detail.clicked) {
            this.$store.commit('addSelected', i)
            break
          }
        }
      }
    },
    async fetchData () {
      // Reset view information.
      this.$store.commit('setReload', false)
      this.$store.commit('resetSelected')
      this.$store.commit('multiple', false)
      this.$store.commit('closeHovers')

      // Set loading to true and reset the error.
      this.setLoading(true)
      this.error = null

      let url = this.$route.path
      if (url === '') url = '/'
      if (url[0] !== '/') url = '/' + url

      try {
        const res = await api.fetch(url)

        if (clean(res.path) !== clean(`/${this.$route.params.pathMatch}`)) {
          return
        }

        this.$store.commit('updateRequest', res)
        document.title = res.name
      } catch (e) {
        this.error = e
      } finally {
        this.setLoading(false)
      }
    },
    keyEvent (event) {
      if (this.show !== null) {
        // Esc!
        if (event.keyCode === 27) {
          this.$store.commit('closeHovers')
        }

        return
      }

      // Esc!
      if (event.keyCode === 27) {        
        // If we're on a listing, unselect all
        // files and folders.
        if (this.isListing) {
          this.$store.commit('resetSelected')
        }
      }

      // Del!
      if (event.keyCode === 46) {
        if (this.isEditor ||
          !this.isFiles ||
          this.loading ||
          !this.user.perm.delete ||
          (this.isListing && this.selectedCount === 0) ||
          this.$store.state.show != null) return

        this.$store.commit('showHover', 'delete')
      }

      // F1!
      if (event.keyCode === 112) {
        event.preventDefault()
        this.$store.commit('showHover', 'help')
      }

      // F2!
      if (event.keyCode === 113) {
        if (this.isEditor ||
          !this.isFiles ||
          this.loading ||
          !this.user.perm.rename ||
          (this.isListing && this.selectedCount === 0) ||
          (this.isListing && this.selectedCount > 1)) return

        this.$store.commit('showHover', 'rename')
      }

      // CTRL + S
      if (event.ctrlKey || event.metaKey) {
        if (this.isEditor) return

        if (String.fromCharCode(event.which).toLowerCase() === 's') {
          event.preventDefault()

          if (this.req.kind !== 'editor') {
            document.getElementById('download-button').click()
          }
        }
      }
    },
    scroll () {
      if (this.req.kind !== 'listing' || this.$store.state.user.viewMode === 'mosaic') return

      let top = 112 - window.scrollY

      if (top < 64) {
        top = 64
      }

      document.querySelector('#listing.list .item.header').style.top = top + 'px'
    },
    openSidebar () {
      this.$store.commit('showHover', 'sidebar')
    },
    openSearch () {
      this.$store.commit('showHover', 'search')
    }
  }
}
</script>
