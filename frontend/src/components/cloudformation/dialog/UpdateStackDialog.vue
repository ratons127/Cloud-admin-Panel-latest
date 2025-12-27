<template>
  <v-dialog persistent v-model="updateDialog" max-width="1000">
    <template v-slot:activator="{ props }">
      <v-btn class="action-btn"  size="small" color="primary" v-bind="props">Update</v-btn>
    </template>
    <v-card>
      <v-card-title class="headline">Update stack {{selectedStack.StackName}}</v-card-title>
      <v-card-text>
        <v-row no-gutters style='margin-top:25px;'>
          <v-file-input v-model="updateTemplateFile" @update:modelValue="updateTemplateFileForUpdate" label="Template file"></v-file-input>
        </v-row>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="green-darken-1"  variant="text" @click="updateDialog=false">Cancel</v-btn>
        <v-btn color="green-darken-1"  variant="text" @click="updateStack">Update</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>

import { mapState } from 'vuex'

export default {
  name: 'CloudFormationUpdateStackDialog',
  computed: {
    ...mapState({
      selectedStack: state => state.cloudformation.selectedStack,
    }),
    updateDialog: {
      get() {
        return this.$store.state.cloudformation.updateDialog
      },
      set(value) {
        this.$store.commit('cloudformation/updateUpdateDialog', value)
      }
    },
    updateTemplateFile: {
      get() {
        return this.$store.state.cloudformation.updateTemplateFile
      },
      set(value) {
        this.$store.commit('cloudformation/updateTemplateFileForUpdate', value)
      }
    }
  },
  methods: {
    updateTemplateFileForUpdate(file){
      this.$store.commit('cloudformation/updateTemplateFileForUpdate', file)
    },
    updateStack(){
      this.$store.dispatch('cloudformation/updateStack')
    }
  }
}
</script>




