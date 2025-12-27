import http from './http'

export default {
  listUsers(serverPath, params, cb, errorCb) {
    http.get(`${serverPath}/admin/users`, { params })
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  createUser(serverPath, payload, cb, errorCb) {
    http.post(`${serverPath}/admin/users`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  updateUser(serverPath, userId, payload, cb, errorCb) {
    http.put(`${serverPath}/admin/users/${userId}`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  deleteUser(serverPath, userId, cb, errorCb) {
    http.delete(`${serverPath}/admin/users/${userId}`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  resetPassword(serverPath, userId, payload, cb, errorCb) {
    http.post(`${serverPath}/admin/users/${userId}/reset_password`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  verifyEmail(serverPath, userId, cb, errorCb) {
    http.post(`${serverPath}/admin/users/${userId}/verify_email`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  listTenants(serverPath, cb, errorCb) {
    http.get(`${serverPath}/admin/tenants`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  createTenant(serverPath, payload, cb, errorCb) {
    http.post(`${serverPath}/admin/tenants`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  updateTenant(serverPath, tenantId, payload, cb, errorCb) {
    http.put(`${serverPath}/admin/tenants/${tenantId}`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  deleteTenant(serverPath, tenantId, cb, errorCb) {
    http.delete(`${serverPath}/admin/tenants/${tenantId}`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  listTenantMembers(serverPath, tenantId, cb, errorCb) {
    http.get(`${serverPath}/admin/tenants/${tenantId}/users`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  addTenantMember(serverPath, tenantId, payload, cb, errorCb) {
    http.post(`${serverPath}/admin/tenants/${tenantId}/users`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  updateTenantMember(serverPath, tenantId, userId, payload, cb, errorCb) {
    http.put(`${serverPath}/admin/tenants/${tenantId}/users/${userId}`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  removeTenantMember(serverPath, tenantId, userId, cb, errorCb) {
    http.delete(`${serverPath}/admin/tenants/${tenantId}/users/${userId}`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  listTenantMembersSelf(serverPath, cb, errorCb) {
    http.get(`${serverPath}/admin/tenant/users`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  addTenantMemberSelf(serverPath, payload, cb, errorCb) {
    http.post(`${serverPath}/admin/tenant/users`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  updateTenantMemberSelf(serverPath, userId, payload, cb, errorCb) {
    http.put(`${serverPath}/admin/tenant/users/${userId}`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  removeTenantMemberSelf(serverPath, userId, cb, errorCb) {
    http.delete(`${serverPath}/admin/tenant/users/${userId}`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  listAwsAccounts(serverPath, cb, errorCb) {
    http.get(`${serverPath}/admin/aws/accounts`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  listAwsAccountsForTenant(serverPath, tenantId, cb, errorCb) {
    http.get(`${serverPath}/admin/tenants/${tenantId}/aws/accounts`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  listAwsAccountsSelf(serverPath, cb, errorCb) {
    http.get(`${serverPath}/admin/tenant/aws/accounts`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  listAuditLogs(serverPath, params, cb, errorCb) {
    http.get(`${serverPath}/admin/audit`, { params })
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  }
}
