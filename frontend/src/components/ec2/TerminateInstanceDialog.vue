<template>
  <v-dialog v-model="terminateDialog" max-width="900" v-if="showActions">
    <template v-slot:activator="{ props }">
      <v-btn class="action-btn" color="error"  size="small" v-bind="props">Terminate</v-btn>
    </template>
    <v-card>
      <v-card-title class="headline">Terminate instance {{selectedInstance.Name}}</v-card-title>
      <v-card-text>
        Confirm terminating instance {{selectedInstance.InstanceId}} ?
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="green-darken-1"  variant="text" @click="terminateDialog = false">Cancel</v-btn>
        <v-btn color="green-darken-1"  variant="text" v-on:click="terminateInstance">Terminate</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'TerminateInstanceDialog',
  computed: {
    ...mapState({
      selectedInstance: state => state.ec2.selectedInstance,
      showActions: state => state.ec2.showActions,
    }),
    terminateDialog: {
      get() {
        return this.$store.state.ec2.terminateDialog
      },
      set(value) {
        this.$store.commit('ec2/updateTerminateDialog', value)
      }
    }
  },
  methods: {
    terminateInstance(){
      this.$store.dispatch('ec2/terminateInstance')
    }
  }
}
</script>



