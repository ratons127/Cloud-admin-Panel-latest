<template>
  <v-main>
    <v-container class="reset-container" fluid>
      <v-row justify="center">
        <v-col cols="12" sm="8" md="5">
          <v-card>
            <v-card-title class="headline">Reset Password</v-card-title>
            <v-card-text>
              <v-text-field
                v-model="password"
                label="New password"
                type="password"
                autocomplete="new-password"
                variant="outlined"
              ></v-text-field>
              <v-text-field
                v-model="confirmPassword"
                label="Confirm password"
                type="password"
                autocomplete="new-password"
                variant="outlined"
              ></v-text-field>
              <v-alert v-if="message" :type="alertType" density="compact">{{ message }}</v-alert>
            </v-card-text>
            <v-card-actions>
              <v-btn variant="text" @click="goLogin">Back to sign in</v-btn>
              <v-spacer></v-spacer>
              <v-btn color="primary" @click="submit" :disabled="loading">Reset</v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>

<script>
import authApi from '../../api/auth'

export default {
  name: 'ResetPassword',
  data: () => ({
    password: '',
    confirmPassword: '',
    status: '',
    message: '',
    loading: false
  }),
  computed: {
    alertType() {
      return this.status === 'error' ? 'error' : 'success'
    },
    serverPath() {
      return this.$store.state.core.serverPath
    }
  },
  methods: {
    submit() {
      this.status = ''
      this.message = ''
      const token = this.$route.query.token
      if (!token) {
        this.status = 'error'
        this.message = 'Missing reset token.'
        return
      }
      if (!this.password || this.password.length < 8) {
        this.status = 'error'
        this.message = 'Password must be at least 8 characters.'
        return
      }
      if (this.password !== this.confirmPassword) {
        this.status = 'error'
        this.message = 'Passwords do not match.'
        return
      }
      this.loading = true
      authApi.reset(this.serverPath, { token, password: this.password }, () => {
        this.status = 'success'
        this.message = 'Password reset. You can sign in now.'
        this.loading = false
        this.$router.push({ path: '/welcome', query: { reset: '1' } })
      }, () => {
        this.status = 'error'
        this.message = 'Reset failed or link expired.'
        this.loading = false
      })
    },
    goLogin() {
      this.$router.push('/')
    }
  }
}
</script>

<style>
.reset-container {
  padding-top: 60px;
}
</style>
