<template>
  <v-dialog persistent v-model="createDialog" max-width="1000">
    <template v-slot:activator="{ props }">
      <v-btn class="action-btn"  size="small" color="primary" v-bind="props">Create</v-btn>
    </template>
    <v-card>
      <v-card-title class="headline">Create AWS Account</v-card-title>
      <v-card-text>
        <v-row no-gutters style='margin-top:25px;'>
          <v-col cols="12">
            <v-text-field
              v-model="name"
              autocomplete="off"
              label="Account Name"
              variant="outlined"
            ></v-text-field>
          </v-col>
          <v-col cols="12">
            <v-text-field
              v-model="accountId"
              autocomplete="off"
              label="AWS Account ID"
              variant="outlined"
            ></v-text-field>
          </v-col>
          <v-col cols="12">
            <v-text-field
              v-model="roleArn"
              autocomplete="off"
              label="Role ARN"
              variant="outlined"
            ></v-text-field>
          </v-col>
          <v-col cols="12">
            <v-text-field
              v-model="externalId"
              autocomplete="off"
              label="External ID"
              variant="outlined"
            ></v-text-field>
          </v-col>
          <v-col cols="12">
            <v-switch v-model="active" label="Active"></v-switch>
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="green-darken-1"  variant="text" @click="createDialog = false">Cancel</v-btn>
        <v-btn color="green-darken-1"  variant="text" @click="createAccount">Create</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
export default {
  name: 'CreateCredentialsDialog',
  computed: {
    createDialog: {
      get() {
        return this.$store.state.credentials.createDialog
      },
      set(value) {
        this.$store.commit('credentials/updateCreateDialog', value)
      }
    },
    name: {
      get() {
        return this.$store.state.credentials.accountEdit.name
      },
      set(value) {
        this.$store.commit('credentials/updateAccountEdit', {name: value})
      }
    },
    accountId: {
      get() {
        return this.$store.state.credentials.accountEdit.accountId
      },
      set(value) {
        this.$store.commit('credentials/updateAccountEdit', {accountId: value})
      }
    },
    roleArn: {
      get() {
        return this.$store.state.credentials.accountEdit.roleArn
      },
      set(value) {
        this.$store.commit('credentials/updateAccountEdit', {roleArn: value})
      }
    },
    externalId: {
      get() {
        return this.$store.state.credentials.accountEdit.externalId
      },
      set(value) {
        this.$store.commit('credentials/updateAccountEdit', {externalId: value})
      }
    },
    active: {
      get() {
        return this.$store.state.credentials.accountEdit.active
      },
      set(value) {
        this.$store.commit('credentials/updateAccountEdit', {active: value})
      }
    }
  },
  methods: {
    createAccount(){
      this.$store.dispatch('credentials/createAccount')
    }
  }
}
</script>
