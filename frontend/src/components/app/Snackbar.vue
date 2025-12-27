<template>
  <v-snackbar v-model="snackbar" location="right" :color="snackbarColor">{{ snackbarMessage }}
    <v-btn  variant="text" @click="snackbar = false">Close</v-btn>
  </v-snackbar>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'Snackbar',
  computed: {
    ...mapState({
      snackbarColor: state => state.core.snackbarColor,
      snackbarMessage: state => state.core.snackbarMessage
    }),
    snackbar: {
      get() {
        return this.$store.state.core.snackbar
      },
      set(value) {
        this.$store.commit('core/updateSnackbar', value)
      }
    }
  }
}
</script>
