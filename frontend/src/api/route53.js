import http from './http'

export default {
  getHostZones(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/route53/hostedzones`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  getRecordSets(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/route53/recordsets?zone=${data.zone}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  deleteRecordSet(data, cb, errorCb) {
    http.delete(`${data.serverPath}/service/route53/recordsets?zone=${data.zone}&recordName=${data.name}&recordType=${data.type}&recordValue=${data.value}&recordTTL=${data.ttl}&alias=${data.alias}&hostzoneID=${data.hostzoneID}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  createRecordSet(data, cb, errorCb) {
    http.post(`${data.serverPath}/service/route53/recordsets?zone=${data.zone}&recordName=${data.name}&recordType=${data.type}&recordValue=${data.value}&recordTTL=${data.ttl}&alias=${data.alias}&hostzoneID=${data.hostzoneID}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
}
