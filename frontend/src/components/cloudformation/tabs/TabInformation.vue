<template>
  <div v-show="selectedStack.StackName != null">
    <v-row no-gutters style='padding-left:25px;padding-top:25px;padding-left:25px;'>
      <v-col cols="6">
        <v-text-field
          :value="description.StackName"
          label="Name"
          variant="outlined"
          readonly
        ></v-text-field>
      </v-col>
      <v-col cols="6" style='padding-left:25px;'>
        <v-text-field
          :value="description.StackStatus"
          label="Status"
          variant="outlined"
          readonly
        ></v-text-field>
      </v-col>
    </v-row>
    <v-row no-gutters style='padding-left:25px;'>
      <v-col cols="12">
        <v-textarea
          variant="outlined"
          readonly
          :value="description.Description">
          <template v-slot:label>
            <div>
              Description
            </div>
          </template>
        </v-textarea>
      </v-col>
    </v-row>
    <v-row no-gutters style='padding-left:25px;'>
      <v-col cols="12">
        <v-text-field
          :value="description.StackId"
          label="StackID"
          variant="outlined"
          readonly
        ></v-text-field>
      </v-col>
    </v-row>
    <v-row no-gutters style='padding-left:25px;padding-left:25px;'>
      <v-col cols="6">
        <v-text-field
          :value="description.CreationTime"
          label="CreationTime"
          variant="outlined"
          readonly
        ></v-text-field>
      </v-col>
      <v-col cols="6" style='padding-left:25px;' v-if="selectedStack.DeletionTime != null">
        <v-text-field
          :value="description.DeletionTime"
          label="DeletionTime"
          variant="outlined"
          readonly
        ></v-text-field>
      </v-col>
    </v-row>
    <v-row no-gutters style='padding-left:25px;padding-left:25px;'>
      <v-col cols="6">
        <v-text-field
          :value="terminationProtection ? 'protected' : 'unprotected'"
          label="Termination protection"
          variant="outlined"
          readonly
        ></v-text-field>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'TabInformation',
  computed: {
    ...mapState({
      description: state => state.cloudformation.description,
      terminationProtection: state => state.cloudformation.terminationProtection,
      selectedStack: state => state.cloudformation.selectedStack,
    })
  }
}
</script>
