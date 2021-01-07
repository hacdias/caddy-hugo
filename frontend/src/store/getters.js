const getters = {
  isLogged: state => state.user !== null,
  isFiles: state => !state.loading && state.route.name === 'Files',
  isListing: (state, getters) => getters.isFiles && state.req.isDir,
  isEditor: (state, getters) => getters.isFiles && (state.req.type === 'text' || state.req.type === 'textImmutable'),
  isPreview: state => state.previewMode,
  isSharing: state =>  !state.loading && state.route.name === 'Share',
  selectedCount: state => state.selected.length,
  progress : state => {
    if (state.upload.progress.length == 0) {
      return 0;
    }

    let sum = state.upload.progress.reduce((acc, val) => acc + val)
    return Math.ceil(sum / state.upload.size * 100);
  },
  getLastViewedDetail: state => (path) => {
    if (state.lastViewed.details.has(path)) {
      let i = state.lastViewed.paths.indexOf(path)
      state.lastViewed.paths.splice(i, 1)
      state.lastViewed.paths.push(path)
      return state.lastViewed.details.get(path)
    }
    return null
  }
}

export default getters
