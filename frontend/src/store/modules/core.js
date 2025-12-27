import vuetify from '../../plugins/vuetify'
import configuration from '../../api/configuration'
import sessionReset from '../../api/sessionReset'

const menuItems = [
  {text:"CloudFormation", value:"CloudFormation",image:"code-json", query:"cloudformation"},
  {text:"System Manager",value: "SSM",image:"console", query: "ssm"},
  {text:"EC2", value:"EC2", image:"server", query:"ec2"},
  {text:"Lambda", value:"Lambda", image: "memory", query:"lambda"},
  {text:"S3", value:"S3", image: "database", query:"s3"},
  {text:"VPCs", value: "VPC", image:"lan", query:"vpc"},
  {text:"Subnets", value: "Subnets", image:"lan", query:"subnets"},
  {text:"Load Balancers", value:"LoadBalancers", image:"lan", query:"loadbalancer"},
  {text: "Target Groups", value: "TargetGroups", image:"server-network", query:"targetgroup"},
  {text:"Security Groups", value:"SecurityGroups",image:"security", query:"securitygroup"},
  {text:"Route 53", value:"Route53", image:"dns", query: "route53"},
  {text:"CloudWatch", value:"CloudWatch",image: "math-log", query:"cloudwatch"},
  {text:"ECR", value:"ECR",image: "docker", query:"ecr"},
  {text:"ECS", value:"ECS",image: "docker", query:"ecs"},
  {text:"Certificate Manager", value:"ACM", image:"certificate-outline", query:"acm"},
  {text:"Keypair", value:"Keypair",image: "key-variant", query:"keypair"},
  {text:"SecretsManager", value:"Secret Manager",image: "key-variant", query:"secretsmanager"},
  {text:"AWS Accounts", value:"Credentials", image:"account",query: "credentials"},
  {text:"Admin", value:"Admin", image:"shield-account", query:"admin", adminOnly:true}
];

const defaultRegion = "eu-west-3";

const filterMenuItems = (items, isAdmin) => {
  return (items || []).filter(item => !item.adminOnly || isAdmin)
}

const state = {
  region: defaultRegion,
  theme: "dark",
  drawer: true,
  title: "AWS Admin Dashboard",
  navbarItems: menuItems,
  remainingNavItems: [],
  regionList: [
    {name: "us-east-2"},
    {name: "us-east-1"},
    {name: "us-west-1"},
    {name: "us-west-2"},
    {name: "ap-east-1"},
    {name: "ap-south-1"},
    {name: "ap-northeast-3"},
    {name: "ap-northeast-2"},
    {name: "ap-southeast-1"},
    {name: "ap-southeast-2"},
    {name: "ap-northeast-1"},
    {name: "ca-central-1"},
    {name: "cn-north-1"},
    {name: "cn-northwest-1"},
    {name: "eu-central-1"},
    {name: "eu-west-1"},
    {name: "eu-west-2"},
    {name: "eu-west-3"},
    {name: "eu-north-1"},
    {name: "me-south-1"},
    {name: "sa-east-1"},
    {name: "us-gov-east-1"},
    {name: "us-gov-west-1"}
  ],
  settingsDialog: false,
  aboutDialog: false,
  navbarSettingsDialog: false,
  serverPath: "",
  serverPathEdit: "",
  loaded: false,
  snackbar: false,
  snackbarColor: 'error',
  snackbarMessage: '',
  isAdmin: false
}

const getters = {
  remainingNavItems(state){
    var result = []
    const baseItems = filterMenuItems(menuItems, state.isAdmin)
    for (var i = 0;i < baseItems.length;i++){
      if (state.navbarItems.filter(it => it.value === baseItems[i].value).length === 0){
          result.push(baseItems[i])
      }
    }
    return result;
  }
}

const mutations = {
  updateRegion(state, value){
    state.region = value
  },
  toggleDark(state){
    const current = vuetify.theme.global.name.value
    state.theme = current === 'dark' ? 'light' : 'dark'
    vuetify.theme.global.name.value = state.theme
    localStorage.setItem("theme", state.theme)
  },
  toggleDrawer(state) {
    state.drawer = !state.drawer
  },
  updateSettingsDialog(state, value){
    state.settingsDialog = value
  },
  updateAboutDialog(state, value){
    state.aboutDialog = value
  },
  updateServerPath(state, value){
    state.serverPath = value
  },
  updateTheme(state, value){
    state.theme = value
  },
  updateLoaded(state, value){
    state.loaded = value
  },
  updateServerPathEdit(state, value){
    state.serverPathEdit = value
  },
  updateSnackbar(state, value){
    state.snackbar = value
  },
  updateSnackbarColor(state, value){
    state.snackbarColor = value
  },
  updateSnackbarMessage(state, value){
    state.snackbarMessage = value
  },
  updateNavbarSettingsDialog(state, value){
    state.navbarSettingsDialog = value
  },
  updateRemainingNavItems(state, value){
    state.remainingNavItems = value
  },
  updateNavbarItems(state, value){
    state.navbarItems = value
    localStorage.setItem('navbarMenu', JSON.stringify(state.navbarItems))
  },
  updateIsAdmin(state, value){
    state.isAdmin = !!value
  },
  resetNavMenu(state){
    state.navbarItems = filterMenuItems(menuItems, state.isAdmin)
    state.remainingNavItems = []
    localStorage.setItem('navbarMenu', JSON.stringify(state.navbarItems))
  }
}

const envServerPath = (import.meta.env.VITE_SERVER_PATH || "").trim()

const actions = {
  init({ commit, state, dispatch, rootState }){
    var storedRegion = localStorage.getItem("region")
    if (storedRegion){
      commit('updateRegion', storedRegion)
    } else {
      commit('updateRegion', defaultRegion)
      localStorage.setItem("region", defaultRegion)
    }
    commit('updateServerPath', localStorage.getItem('serverPath'))
    if (!state.serverPath && envServerPath){
      commit('updateServerPath', envServerPath.replace(/\/$/, ""))
      localStorage.setItem("serverPath", state.serverPath)
    }

    var navbarMenu = localStorage.getItem('navbarMenu')
    if (navbarMenu){
      commit('updateNavbarItems', JSON.parse(navbarMenu))
    } else {
      commit('updateNavbarItems', menuItems)
      localStorage.setItem('navbarMenu', JSON.stringify(state.navbarItems))
    }

    if (!state.serverPath){
      var path = `${location.protocol}//${location.hostname+(location.port ? ':' + location.port : '')}`;
      commit('updateServerPath', path)
      localStorage.setItem("serverPath", state.serverPath)
    }
    commit('updateTheme', localStorage.getItem('theme') || 'dark')
    vuetify.theme.global.name.value = state.theme

    dispatch('auth/init', null, {root: true})
    dispatch('applyMenuAccess')
    commit('updateLoaded', true)
  },
  applyMenuAccess({ commit, state, rootState }){
    const isAdmin = !!(rootState.auth.isSuperAdmin || rootState.auth.isTenantAdmin)
    commit('updateIsAdmin', isAdmin)
    const baseItems = filterMenuItems(menuItems, isAdmin)
    const storedRaw = localStorage.getItem('navbarMenu')
    let stored = []
    if (storedRaw) {
      try {
        stored = JSON.parse(storedRaw) || []
      } catch (e) {
        stored = []
      }
    }
    const filteredStored = filterMenuItems(stored, isAdmin)
    const nextMenu = filteredStored.length > 0 ? filteredStored : baseItems
    commit('updateNavbarItems', nextMenu)
    const remaining = baseItems.filter(it => !nextMenu.some(item => item.value === it.value))
    commit('updateRemainingNavItems', remaining)
  },
  updateRegion({commit, state}, region){
    commit('updateRegion', region)
    configuration.setConfiguration({
      serverPath: state.serverPath,
      payload: {
        region: state.region
      }
    },
    (response) => {
      localStorage.setItem('region', state.region);
      document.location.href = window.location.href 
    },
    (err) => {
      console.log(err);
    })
  },
  resetSession({state}){
    sessionReset.reset({
      serverPath: state.serverPath
    },
    (response) => {
      document.location.href = window.location.href
    },
    (err) => {
      console.log(err);
    })
  },
  openSettingsDialog({commit, state}){
    commit("updateServerPathEdit", state.serverPath)
    commit("updateSettingsDialog", true)
  },
  saveSettings({commit, state}){
    if (state.serverPathEdit != state.serverPath){
      commit("updateServerPath", state.serverPathEdit.replace(/\/$/, ""))
      localStorage.setItem("serverPath", state.serverPath)
      document.location.href = window.location.href
    }
    commit("updateSettingsDialog", false)
  },
  dispatchMessage({commit, state, dispatch}, data){
    commit('updateSnackbarColor', data.color)
    commit('updateSnackbarMessage', data.message)
    commit('updateSnackbar', true)
  },
  toggleNavbarSettingsDialog({commit, getters, state}, value){
    if (!value) {
      commit('updateRemainingNavItems', [])
      commit('updateNavbarSettingsDialog', false)
    } else {
      commit('updateRemainingNavItems', getters.remainingNavItems)
      commit('updateNavbarSettingsDialog', true)
    }
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
