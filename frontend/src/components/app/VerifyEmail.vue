<template>
  <v-main>
    <v-container class="verify-container" fluid>
      <v-row justify="center">
        <v-col cols="12" sm="8" md="6">
          <v-card>
            <v-card-title class="headline">Email Verification</v-card-title>
            <v-card-text>
              <v-alert v-if="message" :type="alertType" density="compact">{{ message }}</v-alert>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="primary" @click="goHome">Back to sign in</v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>

<script>
export default {
  name: 'VerifyEmail',
  data: () => ({
    status: 'pending',
    message: 'Verifying your email...'
  }),
  computed: {
    alertType() {
      return this.status === 'success' ? 'success' : this.status === 'error' ? 'error' : 'info'
    }
  },
  created() {
    const token = this.$route.query.token
    if (!token) {
      this.status = 'error'
      this.message = 'Missing verification token.'
      return
    }
    this.$store.dispatch('auth/verifyEmail', token)
      .then(() => {
        this.status = 'success'
        this.message = 'Your email is verified. Redirecting...'
        this.$router.push({ path: '/welcome', query: { verified: '1' } })
      })
      .catch(() => {
        this.status = 'error'
        this.message = 'Verification failed or link expired.'
      })
  },
  methods: {
    goHome() {
      this.$router.push('/')
    }
  }
}
</script>

<style>
.verify-container {
  padding-top: 60px;
}
</style>
