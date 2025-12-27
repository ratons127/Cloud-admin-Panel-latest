import { createStore } from 'vuex'
import core from './modules/core'
 
import credentials from './modules/credentials'
import cloudformation from './modules/cloudformation'
import ssm from './modules/ssm'
import ec2 from './modules/ec2'
import lambda from './modules/lambda'
import s3 from './modules/s3'
import vpc from './modules/vpc'
import subnets from './modules/subnets'
import loadbalancers from './modules/loadbalancers'
import targetgroups from './modules/targetgroups'
import securitygroups from './modules/securitygroups'
import route53 from './modules/route53'
import cloudwatchlogs from './modules/cloudwatchlogs'
import ecr from './modules/ecr'
import ecs from './modules/ecs'
import acm from './modules/acm'
import keypair from './modules/keypair'
import secretsmanager from './modules/secretsmanager'
import auth from './modules/auth'

export default createStore({
  modules: {
    core,
    credentials,
    cloudformation,
    ssm,
    ec2,
    lambda,
    s3,
    vpc,
    subnets,
    loadbalancers,
    targetgroups,
    securitygroups,
    route53,
    cloudwatchlogs,
    ecr,
    ecs,
    acm,
    keypair,
    secretsmanager,
    auth
  },
  /*
  strict: debug,
  plugins: debug ? [createLogger()] : []
  */
})
