<template>
  <v-container class="login-container" fluid>
    <v-row justify="center">
      <v-col cols="12" sm="8" md="5">
        <v-card>
          <v-card-title class="headline">{{ mode === 'login' ? 'Sign in' : 'Create account' }}</v-card-title>
          <v-card-text>
            <v-text-field
              v-model="email"
              label="Email"
              type="email"
              autocomplete="email"
              variant="outlined"
            ></v-text-field>
            <v-text-field
              v-model="password"
              label="Password"
              type="password"
              autocomplete="current-password"
              variant="outlined"
            ></v-text-field>
            <v-text-field
              v-if="mode === 'signup'"
              v-model="tenantName"
              label="Company / Tenant Name"
              variant="outlined"
            ></v-text-field>
            <v-alert v-if="errorMessage" type="error" density="compact">{{ errorMessage }}</v-alert>
            <v-alert v-if="successMessage" type="success" density="compact">{{ successMessage }}</v-alert>
          </v-card-text>
          <v-card-actions>
            <v-btn variant="text" @click="toggleMode">
              {{ mode === 'login' ? 'Create an account' : 'Back to sign in' }}
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn color="primary" @click="submit">{{ mode === 'login' ? 'Sign in' : 'Sign up' }}</v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
export default {
  name: 'Login',
  data: () => ({
    email: '',
    password: '',
    tenantName: '',
    mode: 'login',
    successMessage: ''
  }),
  computed: {
    errorMessage() {
      return this.$store.state.auth.authError
    }
  },
  methods: {
    toggleMode() {
      this.mode = this.mode === 'login' ? 'signup' : 'login'
      this.successMessage = ''
    },
    submit() {
      this.successMessage = ''
      if (this.mode === 'login') {
        this.$store.dispatch('auth/login', { email: this.email, password: this.password })
          .catch(() => {})
      } else {
        this.$store.dispatch('auth/signup', { email: this.email, password: this.password, tenantName: this.tenantName })
          .then(() => {
            this.successMessage = 'Check your email to verify your account.'
            this.mode = 'login'
            this.password = ''
          })
          .catch(() => {})
      }
    }
  }
}
</script>

<style>
.login-container {
  padding-top: 60px;
}
</style>
