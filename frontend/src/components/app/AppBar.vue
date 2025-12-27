<template>
  <v-app-bar app clipped-left>
    <v-app-bar-nav-icon @click.stop="toggleDrawer" />
    <router-link class="titleLink" to="/"><v-toolbar-title class="title">{{title}}</v-toolbar-title></router-link>
    <v-spacer></v-spacer>
    <v-col cols="auto" class='credentialSelect'>
      <v-select
        v-model="selectedAccountId"
        :items="accounts"
        label="AWS Account"
        item-title="name"
        item-value="id"
      ></v-select>
    </v-col>
    <v-col cols="auto" class='credentialSelect'>
      <v-select
        v-model="selectedTenantId"
        :items="tenants"
        label="Tenant"
        item-title="name"
        item-value="id"
      ></v-select>
    </v-col>
    <v-col cols="auto" class='regionSelect'>
      <v-select
        v-model="region"
        :items="regionList"
        label="Region"
        item-title="name"
        item-value="name"
      ></v-select>
    </v-col>
    <v-btn icon style='margin-right:15px;' @click="openSettingsDialog">
      <v-icon>mdi-cog-outline</v-icon>
    </v-btn>
    <v-btn icon style='margin-right:15px;' @click="aboutDialog = true">
      <v-icon>mdi-information-outline</v-icon>
    </v-btn>
    <InviteDialog />
    <v-btn icon style='margin-right:15px;' @click="logout">
      <v-icon>mdi-logout</v-icon>
    </v-btn>
  </v-app-bar>
</template>

<script>
import { mapState } from 'vuex'
import InviteDialog from './InviteDialog'

  export default {
    name: 'AppBar',
    components: {
      InviteDialog
    },
    computed: {
      ...mapState({
        title: state => state.core.title,
        regionList: state => state.core.regionList,
        accounts: state => state.auth.accounts,
        tenants: state => state.auth.tenants
      }),
      settingsDialog: {
        get() {
          return this.$store.state.core.settingsDialog
        },
        set(value) {
          this.$store.commit('core/updateSettingsDialog', value)
        }
      },
      aboutDialog: {
        get() {
          return this.$store.state.core.aboutDialog
        },
        set(value) {
          this.$store.commit('core/updateAboutDialog', value)
        }
      },
      region: {
        get() {
          return this.$store.state.core.region
        },
        set(value) {
          this.$store.dispatch('core/updateRegion', value);
        }
      },
      selectedAccountId: {
        get() {
          return this.$store.state.auth.selectedAccountId
        },
        set(value) {
          this.$store.dispatch('auth/selectAccount', value)
        }
      },
      selectedTenantId: {
        get() {
          return this.$store.state.auth.selectedTenantId
        },
        set(value) {
          if (value) {
            this.$store.dispatch('auth/switchTenant', value)
          }
        }
      }
    },
    methods: {
      toggleDrawer(){
        this.$store.commit('core/toggleDrawer')
      },
      openSettingsDialog(){
        this.$store.dispatch('core/openSettingsDialog')
      },
      logout(){
        this.$store.dispatch('auth/logout')
      }
    }
  }
</script>

<style>
.regionSelect {
  margin-top: 30px;
  width: 200px;
  margin-right:10px;
}
.credentialSelect {
  margin-top: 30px;
  width: 200px;
  margin-right:10px;
}
</style>
