<template>
  <v-main>
    <v-container class="welcome-container" fluid>
      <v-row justify="center">
        <v-col cols="12" sm="8" md="6">
          <v-card>
            <v-card-title class="headline">Welcome</v-card-title>
            <v-card-text>
              <v-alert type="success" density="compact">{{ message }}</v-alert>
            </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="primary" @click="goLogin">Sign in</v-btn>
            </v-card-actions>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>

<script>
export default {
  name: 'Welcome',
  computed: {
    message() {
      if (this.$route.query.reset === '1') {
        return 'Your password has been reset. You can sign in now.'
      }
      return 'Your email is verified. You can sign in now.'
    }
  },
  methods: {
    goLogin() {
      this.$router.push('/')
    }
  }
}
</script>

<style>
.welcome-container {
  padding-top: 60px;
}
</style>
