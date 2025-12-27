<template>
  <v-dialog v-model="editDialog" max-width="1000" v-if="showActions">
    <template v-slot:activator="{ props }">
      <v-btn class="action-btn"  size="small" color="primary" v-bind="props">
        <v-icon left>mdi-pencil</v-icon> Edit
      </v-btn>
    </template>
    <v-card>
      <v-card-title class="headline">Edit AWS Account</v-card-title>
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
        <v-btn color="green-darken-1"  variant="text" @click="editDialog = false">Cancel</v-btn>
        <v-btn color="green-darken-1"  variant="text" @click="saveAccount">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'EditCredentialsDialog',
  computed: {
    ...mapState({
      showActions: state => state.credentials.showActions
    }),
    editDialog: {
      get() {
        return this.$store.state.credentials.editDialog
      },
      set(value) {
        this.$store.commit('credentials/updateEditDialog', value)
      }
    },
    name: {
      get() {
        return this.$store.state.credentials.selectedAccount.name
      },
      set(value) {
        this.$store.commit('credentials/updateSelectedAccount', {name: value})
      }
    },
    accountId: {
      get() {
        return this.$store.state.credentials.selectedAccount.accountId
      },
      set(value) {
        this.$store.commit('credentials/updateSelectedAccount', {accountId: value})
      }
    },
    roleArn: {
      get() {
        return this.$store.state.credentials.selectedAccount.roleArn
      },
      set(value) {
        this.$store.commit('credentials/updateSelectedAccount', {roleArn: value})
      }
    },
    externalId: {
      get() {
        return this.$store.state.credentials.selectedAccount.externalId
      },
      set(value) {
        this.$store.commit('credentials/updateSelectedAccount', {externalId: value})
      }
    },
    active: {
      get() {
        return this.$store.state.credentials.selectedAccount.active
      },
      set(value) {
        this.$store.commit('credentials/updateSelectedAccount', {active: value})
      }
    }
  },
  methods: {
    saveAccount(){
      this.$store.dispatch('credentials/saveAccount')
    }
  }
}
</script>
