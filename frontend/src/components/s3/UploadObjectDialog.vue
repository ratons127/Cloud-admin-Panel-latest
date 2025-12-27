<template>
  <v-dialog v-model="uploadDialog" max-width="1000">
    <template v-slot:activator="{ props }">
      <v-btn class="action-btn"  size="small" color="primary" v-bind="props">Upload object</v-btn>
    </template>
    <v-card>
      <v-card-title class="headline">Upload object</v-card-title>
      <v-card-text>
        <v-row no-gutters style='margin-top:25px;'>
          <v-col cols="12">
            <v-file-input v-model="createdObject" @update:modelValue="onUploadFile" label="File"></v-file-input>
          </v-col>
        </v-row>
        <v-row no-gutters style='margin-top:15px;'>
          <v-col cols="12">
            <v-text-field
              v-model="createdObjectKey"
              label="Object path"
              variant="outlined"
            ></v-text-field>
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="green-darken-1"  variant="text" @click="uploadDialog = false">Cancel</v-btn>
        <v-btn color="green-darken-1"  variant="text" @click="uploadObject">Create</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>

export default {
  name: 'UpdateObjectDialog',
  computed: {
    uploadDialog: {
      get() {
        return this.$store.state.s3.uploadDialog
      },
      set(value) {
        this.$store.commit('s3/updateUploadDialog', value)
      }
    },
    createdObject: {
      get() {
        return this.$store.state.s3.createdObject
      },
      set(value) {
        this.$store.commit('s3/updateCreateObject', value)
      }
    },
    createdObjectKey: {
      get() {
        return this.$store.state.s3.createdObjectKey
      },
      set(value) {
        this.$store.commit('s3/updateCreateObjectKey', value)
      }
    }
  },
  methods: {
    onUploadFile(file){
      this.$store.commit('s3/onUploadFile', file)
    },
    uploadObject(){
      this.$store.dispatch('s3/uploadObject')
    }
  }
}
</script>





