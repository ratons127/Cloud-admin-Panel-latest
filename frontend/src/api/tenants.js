import http from './http'

export default {
  list(serverPath, cb, errorCb) {
    http.get(`${serverPath}/tenants`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  create(serverPath, payload, cb, errorCb) {
    http.post(`${serverPath}/tenants`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  switchTenant(serverPath, id, cb, errorCb) {
    http.post(`${serverPath}/tenants/${id}/switch`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  invite(serverPath, id, payload, cb, errorCb) {
    http.post(`${serverPath}/tenants/${id}/invite`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  accept(serverPath, payload, cb, errorCb) {
    http.post(`${serverPath}/tenants/accept`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  }
}
