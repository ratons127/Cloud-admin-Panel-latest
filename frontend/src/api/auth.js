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
<<<<<<< HEAD
  acceptInvite(serverPath, payload, cb, errorCb) {
    http.post(`${serverPath}/auth/invite/accept`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
=======
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
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
<<<<<<< HEAD
  forgot(serverPath, payload, cb, errorCb) {
    http.post(`${serverPath}/auth/forgot`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
  reset(serverPath, payload, cb, errorCb) {
    http.post(`${serverPath}/auth/reset`, payload)
      .then((response) => cb(response))
      .catch((err) => errorCb(err))
  },
=======
>>>>>>> e0ac5ea8763b5bbbe5af1dffd73ebd9de417e8af
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
