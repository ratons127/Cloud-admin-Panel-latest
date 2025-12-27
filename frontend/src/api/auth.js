import http from './http'

export default {
  signup(serverPath, payload, cb, errorCb) {
    http.post(`${serverPath}/auth/signup`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  verify(serverPath, token, cb, errorCb) {
    http.get(`${serverPath}/auth/verify?token=${encodeURIComponent(token)}`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  login(serverPath, payload, cb, errorCb) {
    http.post(`${serverPath}/auth/login`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  refresh(serverPath, payload, cb, errorCb) {
    http.post(`${serverPath}/auth/refresh`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  logout(serverPath, payload, cb, errorCb) {
    http.post(`${serverPath}/auth/logout`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  me(serverPath, cb, errorCb) {
    http.get(`${serverPath}/auth/me`)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  }
}
