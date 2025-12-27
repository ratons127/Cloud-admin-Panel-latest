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
                Please sign in to accept the invitation.
              </div>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn v-if="status === 'success'" color="primary" @click="goHome">Go to dashboard</v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
      <Login v-if="!isAuthenticated" />
    </v-container>
  </v-main>
</template>

<script>
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
    goHome() {
      this.$router.push('/')
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
</style>
