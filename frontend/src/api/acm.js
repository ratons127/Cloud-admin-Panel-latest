import http from './http'

export default {
  getCertificates(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/acm/list`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  deleteCertificate(data, cb, errorCb) {
    http.delete(`${data.serverPath}/service/acm?arn=${data.arn}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  importCertificate(data, cb, errorCb) {
    http.put(`${data.serverPath}/service/acm/import`, data.formData)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
}
