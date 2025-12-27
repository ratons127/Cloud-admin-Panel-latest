<template>
  <v-container class="login-container" fluid>
    <v-row justify="center">
<<<<<<< HEAD
      <v-col cols="12" sm="10" md="6" lg="4">
        <v-card class="login-card">
          <v-card-title class="login-header">
            <div class="login-logo">
              <span class="logo-cloud">cloud</span><span class="logo-sys">sys24</span>
            </div>
            <div class="headline login-title">{{ mode === 'login' ? 'Sign in' : 'Create account' }}</div>
          </v-card-title>
=======
      <v-col cols="12" sm="8" md="5">
        <v-card>
          <v-card-title class="headline">{{ mode === 'login' ? 'Sign in' : 'Create account' }}</v-card-title>
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
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
<<<<<<< HEAD
            <div v-if="mode === 'login'" class="oauth-row">
              <v-btn class="login-google" variant="outlined" @click="googleLogin">
                Continue with Google
              </v-btn>
            </div>
=======
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
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
<<<<<<< HEAD
            <v-btn variant="text" class="login-link" @click="toggleMode">
              {{ mode === 'login' ? 'Create an account' : 'Back to sign in' }}
            </v-btn>
            <v-btn v-if="mode === 'login'" variant="text" class="login-link" @click="forgotPassword">
              Forgot password?
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn class="login-primary" @click="submit">{{ mode === 'login' ? 'Sign in' : 'Sign up' }}</v-btn>
=======
            <v-btn variant="text" @click="toggleMode">
              {{ mode === 'login' ? 'Create an account' : 'Back to sign in' }}
            </v-btn>
            <v-spacer></v-spacer>
            <v-btn color="primary" @click="submit">{{ mode === 'login' ? 'Sign in' : 'Sign up' }}</v-btn>
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
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
<<<<<<< HEAD
    },
    serverPath() {
      return this.$store.state.core.serverPath || `${location.protocol}//${location.hostname+(location.port ? ':' + location.port : '')}`
=======
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
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
<<<<<<< HEAD
    },
    forgotPassword() {
      this.$router.push('/forgot')
    },
    googleLogin() {
      window.location.href = `${this.serverPath}/auth/google`
=======
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
    }
  }
}
</script>

<style>
.login-container {
<<<<<<< HEAD
  min-height: 100vh;
  display: flex;
  align-items: center;
  background: radial-gradient(1200px circle at 20% 20%, #1b2a41 0%, #0b1020 55%, #06080f 100%);
}
.login-container::before {
  content: "";
  position: absolute;
  inset: 0;
  background-image: linear-gradient(120deg, rgba(247, 149, 73, 0.08), transparent 55%),
    radial-gradient(600px circle at 80% 10%, rgba(255, 255, 255, 0.06), transparent 60%);
  pointer-events: none;
}
.login-container > .v-row {
  width: 100%;
  position: relative;
  z-index: 1;
  padding: 24px 12px;
}
.login-card {
  border-top: 4px solid #f79549;
}
.login-header {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.login-logo {
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  font-size: 0.85rem;
}
.logo-cloud {
  color: #ffffff;
}
.logo-sys {
  color: #f79549;
}
.login-title {
  color: #f79549;
}
.login-primary {
  background: #f79549 !important;
  color: #ffffff !important;
}
.login-link {
  color: #f79549 !important;
}
.oauth-row {
  margin-top: 8px;
  margin-bottom: 8px;
}
.login-google {
  border-color: #f79549 !important;
  color: #f79549 !important;
  text-transform: none;
}
@media (max-width: 600px) {
  .login-container > .v-row {
    padding: 16px 8px;
  }
  .v-card-actions {
    flex-wrap: wrap;
    row-gap: 8px;
  }
  .login-primary {
    width: 100%;
  }
=======
  padding-top: 60px;
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
}
</style>
