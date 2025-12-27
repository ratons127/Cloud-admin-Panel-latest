import http from './http'

export default {
  getInstances(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/ec2/instances`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  rebootInstance(data, cb, errorCb) {
    http.post(`${data.serverPath}/service/ec2/instance?id=${data.instanceId}&status=reboot`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  startInstance(data, cb, errorCb) {
    http.post(`${data.serverPath}/service/ec2/instance?id=${data.instanceId}&status=start`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  stopInstance(data, cb, errorCb) {
    http.post(`${data.serverPath}/service/ec2/instance?id=${data.instanceId}&status=stop`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  terminateInstance(data, cb, errorCb) {
    http.post(`${data.serverPath}/service/ec2/instance?id=${data.instanceId}&status=terminate`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  getVPC(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/ec2/vpc`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  getSubnets(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/ec2/subnets`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  getSecurityGroups(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/ec2/sg`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  getKeypairs(data, cb, errorCb) {
    http.get(`${data.serverPath}/service/ec2/keypair`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  createKeypair(data, cb, errorCb) {
    http.put(`${data.serverPath}/service/ec2/keypair?name=${data.name}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
  deleteKeypair(data, cb, errorCb) {
    http.delete(`${data.serverPath}/service/ec2/keypair?name=${data.name}`)
      .then(function(response){
        cb(response)
      })
      .catch(function(err){
        errorCb(err)
      });
  },
}
