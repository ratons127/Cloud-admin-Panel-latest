<template>
  <v-main>
    <v-container class="oauth-container" fluid>
      <v-row justify="center">
        <v-col cols="12" sm="8" md="6">
          <v-card>
            <v-card-title class="headline">Signing you in</v-card-title>
            <v-card-text>
              <v-alert v-if="message" :type="alertType" density="compact">{{ message }}</v-alert>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn v-if="status === 'error'" color="primary" @click="goLogin">Back to sign in</v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>

<script>
export default {
  name: 'OAuthCallback',
  data: () => ({
    status: 'pending',
    message: 'Completing Google sign-in...'
  }),
  computed: {
    alertType() {
      return this.status === 'error' ? 'error' : this.status === 'success' ? 'success' : 'info'
    }
  },
  created() {
    const { accessToken, refreshToken, error } = this.$route.query
    if (error) {
      this.status = 'error'
      this.message = 'Google sign-in failed. Please try again.'
      return
    }
    if (!accessToken || !refreshToken) {
      this.status = 'error'
      this.message = 'Missing login tokens.'
      return
    }
    this.$store.commit('auth/setTokens', { accessToken, refreshToken })
    this.$store.dispatch('auth/loadMe').then(() => {
      this.status = 'success'
      this.message = 'Signed in successfully.'
      this.$router.push('/')
    })
  },
  methods: {
    goLogin() {
      this.$router.push('/')
    }
  }
}
</script>

<style>
.oauth-container {
  padding-top: 60px;
}
</style>
