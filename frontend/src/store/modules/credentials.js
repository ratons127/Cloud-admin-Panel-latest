import awsAccounts from '../../api/awsAccounts'

const emptyAccount = {
  name: '',
  accountId: '',
  roleArn: '',
  externalId: '',
  active: true
}

const state = {
  createDialog: false,
  deleteDialog: false,
  editDialog: false,
  showActions: false,
  selectedItems: [],
  accounts: [],
  accountEdit: { ...emptyAccount },
  selectedAccount: { ...emptyAccount },
  accountsHeaders: [{
    title: 'Name',
    key: 'name'
  },{
    title: 'Account ID',
    key: 'accountId'
  },{
    title: 'Role ARN',
    key: 'roleArn'
  },{
    title: 'Active',
    key: 'active'
  }]
}

const mutations = {
  updateCreateDialog(state, value){
    state.createDialog = value
  },
  updateDeleteDialog(state, value){
    state.deleteDialog = value
  },
  updateEditDialog(state, value){
    state.editDialog = value
  },
  updateSelected(state, value){
    state.selectedItems = value
  },
  setSelection(state, value){
    state.showActions = true
    state.selectedAccount = value
    state.selectedItems = [value]
  },
  updateAccountEdit(state, value){
    state.accountEdit = { ...state.accountEdit, ...value }
  },
  updateSelectedAccount(state, value){
    state.selectedAccount = { ...state.selectedAccount, ...value }
  },
  resetAccountEdit(state){
    state.accountEdit = { ...emptyAccount }
  },
  resetSelectedAccount(state){
    state.selectedAccount = { ...emptyAccount }
  },
  updateAccounts(state, value){
    state.accounts = value || []
  }
}

const actions = {
  fetchAccounts({commit, rootState}){
    awsAccounts.list(rootState.core.serverPath,
      (resp) => commit('updateAccounts', resp.data),
      (err) => console.log(err)
    )
  },
  createAccount({commit, dispatch, state, rootState}){
    if (!state.accountEdit.name || !state.accountEdit.accountId || !state.accountEdit.roleArn || !state.accountEdit.externalId){
      dispatch('core/dispatchMessage', {color: 'error', message: 'name, accountId, roleArn, externalId are required'}, {root:true})
      return
    }
    awsAccounts.create(rootState.core.serverPath, state.accountEdit,
      () => {
        commit('resetAccountEdit')
        commit('updateCreateDialog', false)
        dispatch('fetchAccounts')
      },
      (err) => console.log(err)
    )
  },
  saveAccount({commit, dispatch, state, rootState}){
    awsAccounts.update(rootState.core.serverPath, state.selectedAccount.id, state.selectedAccount,
      () => {
        commit('updateEditDialog', false)
        dispatch('fetchAccounts')
      },
      (err) => console.log(err)
    )
  },
  deleteAccount({commit, dispatch, state, rootState}){
    awsAccounts.remove(rootState.core.serverPath, state.selectedAccount.id,
      () => {
        commit('resetSelectedAccount')
        commit('updateSelected', [])
        commit('updateDeleteDialog', false)
        dispatch('fetchAccounts')
      },
      (err) => console.log(err)
    )
  },
  onClickRow({commit}, item){
    commit('setSelection', item)
  },
  onItemSelected({commit}, selected){
    if (Array.isArray(selected) && selected.length > 0) {
      commit('setSelection', selected[0])
    } else {
      commit('updateSelected', [])
    }
  }
}

export default {
  namespaced: true,
  state,
  actions,
  mutations
}
