<template>
  <div>
    <template v-if="isAuthenticated">
      <NavigationDrawer />
      <AppBar />
      <SettingsDialog />
      <AboutDialog />
    </template>
    <v-main v-if="loaded">
      <Login v-if="!isAuthenticated" />
      <template v-else>
        <Root />
        <Snackbar />
      </template>
    </v-main>
    <!--
    <v-footer app>
      <span>&copy; 2020</span>
    </v-footer>
      -->
  </div>
</template>

<script>
import { mapState } from 'vuex'

import AppBar from './AppBar'
import NavigationDrawer from './NavigationDrawer'
import SettingsDialog from './SettingsDialog'
import AboutDialog from './AboutDialog'
import Root from './Root'
import Snackbar from './Snackbar'
import Login from './Login'

export default {
  name: 'AppShell',
  components: {
    NavigationDrawer,
    AppBar,
    SettingsDialog,
    AboutDialog,
    Root,
    Snackbar,
    Login
  },
  computed: {
    ...mapState({
      loaded: state => state.core.loaded
    }),
    isAuthenticated() {
      return this.$store.getters['auth/isAuthenticated']
    }
  }
}
</script>
