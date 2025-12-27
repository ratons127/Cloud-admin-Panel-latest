<template>
  <v-main>
    <v-container class="invite-container" fluid>
      <v-row justify="center">
        <v-col cols="12" sm="8" md="6">
          <v-card>
            <v-card-title class="headline">Accept Invitation</v-card-title>
            <v-card-text>
              <v-alert v-if="message" :type="alertType" density="compact">{{ message }}</v-alert>
              <div v-if="!isAuthenticated" class="invite-note">
<<<<<<< HEAD
                Set a password to accept the invitation.
              </div>
              <div v-if="!isAuthenticated" class="invite-form">
                <v-checkbox v-model="hasAccount" label="I already have an account" density="compact" />
                <v-text-field
                  v-model="password"
                  label="Password"
                  type="password"
                  autocomplete="new-password"
                  variant="outlined"
                  :disabled="hasAccount"
                ></v-text-field>
                <v-text-field
                  v-model="confirmPassword"
                  label="Confirm password"
                  type="password"
                  autocomplete="new-password"
                  variant="outlined"
                  :disabled="hasAccount"
                ></v-text-field>
              </div>
            </v-card-text>
            <v-card-actions>
              <v-btn v-if="!isAuthenticated" variant="text" @click="goLogin">Back to sign in</v-btn>
              <v-spacer></v-spacer>
              <v-btn v-if="!isAuthenticated" color="primary" @click="acceptInvitePublic" :disabled="processing">Accept invitation</v-btn>
=======
                Please sign in to accept the invitation.
              </div>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
              <v-btn v-if="status === 'success'" color="primary" @click="goHome">Go to dashboard</v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
<<<<<<< HEAD
=======
      <Login v-if="!isAuthenticated" />
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
    </v-container>
  </v-main>
</template>

<script>
<<<<<<< HEAD
import authApi from '../../api/auth'

export default {
  name: 'InviteAccept',
  data: () => ({
    status: 'idle',
    message: 'Waiting for invitation...',
    processing: false,
    completed: false,
    password: '',
    confirmPassword: '',
    hasAccount: false
=======
import Login from './Login'

export default {
  name: 'InviteAccept',
  components: {
    Login
  },
  data: () => ({
    status: 'idle',
    message: 'Waiting for sign in...',
    processing: false,
    completed: false
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
  }),
  computed: {
    isAuthenticated() {
      return this.$store.getters['auth/isAuthenticated']
    },
    alertType() {
      return this.status === 'success' ? 'success' : this.status === 'error' ? 'error' : 'info'
    }
  },
  created() {
    if (this.isAuthenticated) {
      this.acceptInvite()
<<<<<<< HEAD
    } else {
      this.status = 'idle'
      this.message = 'Enter a password to accept the invitation.'
=======
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
    }
  },
  watch: {
    isAuthenticated(value) {
      if (value) {
        this.acceptInvite()
      }
    }
  },
  methods: {
    acceptInvite() {
      if (this.processing || this.completed) {
        return
      }
      const token = this.$route.query.token
      if (!token) {
        this.status = 'error'
        this.message = 'Missing invite token.'
        return
      }
      this.processing = true
      this.status = 'pending'
      this.message = 'Accepting invitation...'
      this.$store.dispatch('auth/acceptInvite', token)
        .then(() => this.$store.dispatch('auth/loadMe'))
        .then(() => {
          this.status = 'success'
          this.message = 'Invite accepted.'
          this.completed = true
        })
        .catch(() => {
          this.status = 'error'
          this.message = 'Invite invalid or expired.'
        })
        .finally(() => {
          this.processing = false
        })
    },
<<<<<<< HEAD
    acceptInvitePublic() {
      if (this.processing || this.completed) {
        return
      }
      const token = this.$route.query.token
      if (!token) {
        this.status = 'error'
        this.message = 'Missing invite token.'
        return
      }
      if (!this.hasAccount) {
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
      }
      this.processing = true
      this.status = 'pending'
      this.message = 'Accepting invitation...'
      const payload = this.hasAccount ? { token } : { token, password: this.password }
      authApi.acceptInvite(this.$store.state.core.serverPath, payload, (resp) => {
        if (resp.data && resp.data.requiresLogin) {
          this.status = 'success'
          this.message = 'Invite accepted. Please sign in.'
        } else {
          this.status = 'success'
          this.message = 'Invite accepted. You can sign in now.'
        }
        this.completed = true
      }, () => {
        this.status = 'error'
        this.message = 'Invite invalid or expired.'
      }).finally(() => {
        this.processing = false
      })
    },
    goHome() {
      this.$router.push('/')
    },
    goLogin() {
      this.$router.push('/')
=======
    goHome() {
      this.$router.push('/')
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
    }
  }
}
</script>

<style>
.invite-container {
  padding-top: 60px;
}
.invite-note {
  margin-top: 10px;
}
<<<<<<< HEAD
.invite-form {
  margin-top: 10px;
}
=======
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
</style>
