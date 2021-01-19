<template>
  <button :title="$t('buttons.info')" :aria-label="$t('buttons.info')" class="action" @click="$emit('rotate', rotate_left)">
    <i class="material-icons">{{ this.icon }}</i>
    <span>{{ $t('buttons.info') }}</span>
  </button>
</template>

<script>
export default {
  name: 'rotate-button',
  props: [ 'rotate_left' ],
  computed: {
    icon () {
      if (this.rotate_left) {
        return 'rotate_left'
      }

      return 'rotate_right'
    }
  }
}
</script>
