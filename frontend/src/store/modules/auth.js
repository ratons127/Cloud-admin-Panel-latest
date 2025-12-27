import authApi from '../../api/auth'
import tenantsApi from '../../api/tenants'

const state = {
  accessToken: localStorage.getItem('accessToken') || '',
  refreshToken: localStorage.getItem('refreshToken') || '',
  userId: '',
  tenantId: '',
  email: '',
  accounts: [],
  selectedAccountId: localStorage.getItem('awsAccountId') || '',
  tenants: [],
  selectedTenantId: localStorage.getItem('tenantId') || '',
  region: localStorage.getItem('region') || '',
  loggedIn: false,
  authError: '',
  isSuperAdmin: false,
  tenantRole: '',
  isTenantAdmin: false
}

const getters = {
  isAuthenticated: (state) => !!state.accessToken
}

const mutations = {
  setTokens(state, payload) {
    state.accessToken = payload.accessToken
    state.refreshToken = payload.refreshToken
    localStorage.setItem('accessToken', payload.accessToken)
    localStorage.setItem('refreshToken', payload.refreshToken)
  },
  clearTokens(state) {
    state.accessToken = ''
    state.refreshToken = ''
    state.isSuperAdmin = false
    state.tenantRole = ''
    state.isTenantAdmin = false
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
  },
  setProfile(state, payload) {
    state.userId = payload.userId
    state.tenantId = payload.tenantId
    state.email = payload.email
    state.region = payload.region || ''
    state.isSuperAdmin = !!payload.isSuperAdmin
    state.tenantRole = payload.role || ''
    state.isTenantAdmin = ['owner', 'admin'].includes((payload.role || '').toLowerCase()) || state.isSuperAdmin
    if (state.region) {
      localStorage.setItem('region', state.region)
    }
    state.selectedTenantId = payload.tenantId
    localStorage.setItem('tenantId', payload.tenantId)
  },
  setAccounts(state, accounts) {
    state.accounts = accounts || []
  },
  setTenants(state, tenants) {
    state.tenants = tenants || []
  },
  setSelectedAccount(state, accountId) {
    state.selectedAccountId = accountId || ''
    if (accountId) {
      localStorage.setItem('awsAccountId', accountId)
    } else {
      localStorage.removeItem('awsAccountId')
    }
  },
  setSelectedTenant(state, tenantId) {
    state.selectedTenantId = tenantId || ''
    if (tenantId) {
      localStorage.setItem('tenantId', tenantId)
    } else {
      localStorage.removeItem('tenantId')
    }
  },
  setLoggedIn(state, value) {
    state.loggedIn = value
  },
  setAuthError(state, message) {
    state.authError = message || ''
  }
}

const actions = {
  init({ state, dispatch, commit, rootState }) {
    if (!state.accessToken || !state.refreshToken) {
      commit('setLoggedIn', false)
      return
    }
    dispatch('loadMe')
  },
  loadMe({ commit, dispatch, rootState }) {
    return new Promise((resolve) => {
      authApi.me(rootState.core.serverPath, (resp) => {
        commit('setProfile', resp.data)
        commit('setAccounts', resp.data.accounts || [])
        commit('setTenants', resp.data.tenants || [])
        if (resp.data.region) {
          commit('core/updateRegion', resp.data.region, {root:true})
        }
        if (!localStorage.getItem('awsAccountId') && resp.data.accounts && resp.data.accounts.length > 0) {
          commit('setSelectedAccount', resp.data.accounts[0].id)
        }
        commit('setLoggedIn', true)
        dispatch('core/applyMenuAccess', null, {root:true})
        resolve()
      }, (err) => {
        commit('clearTokens')
        commit('setLoggedIn', false)
        resolve(err)
      })
    })
  },
  signup({ commit, rootState }, payload) {
    commit('setAuthError', '')
    return new Promise((resolve, reject) => {
      authApi.signup(rootState.core.serverPath, payload, (resp) => resolve(resp), (err) => {
        commit('setAuthError', 'Signup failed')
        reject(err)
      })
    })
  },
  verifyEmail({ commit, rootState }, token) {
    return new Promise((resolve, reject) => {
      authApi.verify(rootState.core.serverPath, token, (resp) => resolve(resp), (err) => reject(err))
    })
  },
  login({ commit, dispatch, rootState }, payload) {
    commit('setAuthError', '')
    return new Promise((resolve, reject) => {
      authApi.login(rootState.core.serverPath, payload, (resp) => {
        commit('setTokens', resp.data)
        dispatch('loadMe').then(() => resolve(resp))
      }, (err) => {
        commit('setAuthError', 'Login failed')
        reject(err)
      })
    })
  },
  refresh({ commit, rootState, state }) {
    return new Promise((resolve, reject) => {
      authApi.refresh(rootState.core.serverPath, { refreshToken: state.refreshToken }, (resp) => {
        commit('setTokens', resp.data)
        resolve(resp)
      }, (err) => {
        commit('clearTokens')
        reject(err)
      })
    })
  },
  logout({ commit, dispatch, rootState, state }) {
    return new Promise((resolve) => {
      authApi.logout(rootState.core.serverPath, { refreshToken: state.refreshToken }, () => {
        commit('clearTokens')
        commit('setLoggedIn', false)
        commit('setAccounts', [])
        commit('setSelectedAccount', '')
        commit('setTenants', [])
        commit('setSelectedTenant', '')
        commit('setProfile', {userId:'', tenantId:'', email:'', region:'', role:'', isSuperAdmin:false})
        dispatch('core/applyMenuAccess', null, {root:true})
        resolve()
      }, () => {
        commit('clearTokens')
        commit('setLoggedIn', false)
        dispatch('core/applyMenuAccess', null, {root:true})
        resolve()
      })
    })
  },
  selectAccount({ commit }, accountId) {
    commit('setSelectedAccount', accountId)
  },
  switchTenant({ commit, dispatch, rootState }, tenantId) {
    return new Promise((resolve, reject) => {
      tenantsApi.switchTenant(rootState.core.serverPath, tenantId, (resp) => {
        commit('setTokens', resp.data)
        commit('setSelectedTenant', tenantId)
        dispatch('loadMe').then(() => resolve(resp))
      }, (err) => reject(err))
    })
  },
  inviteUser({ rootState }, payload) {
    return new Promise((resolve, reject) => {
      tenantsApi.invite(rootState.core.serverPath, payload.tenantId, {
        email: payload.email,
        role: payload.role
      }, (resp) => resolve(resp), (err) => reject(err))
    })
  },
  acceptInvite({ rootState }, token) {
    return new Promise((resolve, reject) => {
      tenantsApi.accept(rootState.core.serverPath, { token }, (resp) => resolve(resp), (err) => reject(err))
    })
  },
  createTenant({ rootState, dispatch }, payload) {
    return new Promise((resolve, reject) => {
      tenantsApi.create(rootState.core.serverPath, payload, (resp) => {
        dispatch('loadMe').then(() => resolve(resp))
      }, (err) => reject(err))
    })
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}
