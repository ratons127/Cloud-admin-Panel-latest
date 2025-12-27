import http from './http'

export default {
  list(serverPath, cb, errorCb) {
    http.get(`${serverPath}/aws/accounts`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  create(serverPath, payload, cb, errorCb) {
    http.post(`${serverPath}/aws/accounts`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  update(serverPath, id, payload, cb, errorCb) {
    http.put(`${serverPath}/aws/accounts/${id}`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  remove(serverPath, id, cb, errorCb) {
    http.delete(`${serverPath}/aws/accounts/${id}`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  }
}
