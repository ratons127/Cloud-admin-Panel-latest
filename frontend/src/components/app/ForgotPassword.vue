<template>
  <v-main>
    <v-container class="forgot-container" fluid>
      <v-row justify="center">
        <v-col cols="12" sm="8" md="5">
          <v-card>
            <v-card-title class="headline">Forgot Password</v-card-title>
            <v-card-text>
              <v-text-field
                v-model="email"
                label="Email"
                type="email"
                autocomplete="email"
                variant="outlined"
              ></v-text-field>
              <v-alert v-if="message" :type="alertType" density="compact">{{ message }}</v-alert>
            </v-card-text>
            <v-card-actions>
              <v-btn variant="text" @click="goLogin">Back to sign in</v-btn>
              <v-spacer></v-spacer>
              <v-btn color="primary" @click="submit" :disabled="loading">Send reset link</v-btn>
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
  name: 'ForgotPassword',
  data: () => ({
    email: '',
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
      this.loading = true
      authApi.forgot(this.serverPath, { email: this.email }, () => {
        this.status = 'success'
        this.message = 'If that email exists, a reset link has been sent.'
        this.loading = false
      }, () => {
        this.status = 'error'
        this.message = 'Failed to request password reset.'
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
.forgot-container {
  padding-top: 60px;
}
</style>
