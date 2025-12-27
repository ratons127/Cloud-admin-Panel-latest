<template>
  <v-menu
    v-model="timeMenu"
    :close-on-content-click="false"
    transition="scale-transition"
    offset="8"
    max-width="290px"
    min-width="290px"
  >
    <template v-slot:activator="{ props }">
      <v-text-field
      v-model="startTimeStr"
      label="Time filter"
      prepend-inner-icon="mdi-clock-outline"
      style='margin-left:25px; margin-top:0; padding-top:0;max-width:350px;'
      readonly
      v-bind="props"
      ></v-text-field>
    </template>
    <v-time-picker
      v-if="timeMenu"
      v-model="startTimeStr"
      format="24hr"
      full-width
      @update:modelValue="startTimeStrChange"
    ></v-time-picker>
  </v-menu>
</template>

<script>
export default {
  name: 'TimeMenu',
  computed: {
    timeMenu: {
      get() {
        return this.$store.state.cloudwatchlogs.timeMenu
      },
      set(value) {
        this.$store.commit('cloudwatchlogs/updateTimeMenu', value)
      }
    },
    startTimeStr: {
      get() {
        return this.$store.state.cloudwatchlogs.startTimeStr
      },
      set(value) {
        this.$store.commit('cloudwatchlogs/updateStartTimeStr', value)
      }
    }
  },
  methods: {
    startTimeStrChange(){
      this.$store.commit('cloudwatchlogs/updateTimeMenu', false)
      this.$store.dispatch('cloudwatchlogs/onStartTimeChange')
    }
  }
}
</script>

