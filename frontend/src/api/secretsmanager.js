import http from './http'

export default {
  getSecrets(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/secretsmanager`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  createSecret(data, cb, errorCb) {
    http.post(`${data.serverPath}/service/secretsmanager`, data.formData)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  deleteSecret(data, cb, errorCb) {
    http.delete(`${data.serverPath}/service/secretsmanager?arn=${data.arn}&forceDelete=${data.forceDelete}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  getSecretValue(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/secretsmanager/value?arn=${data.arn}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
}
