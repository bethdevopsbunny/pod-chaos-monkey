# Kubernetes manifests

Kubernetes manifests for the deployment of the Pod Chaos Monkey and its supporting infrastructure. <br><br>
Also included some nginx Kubernetes manifests for testing purposes.

## pod_chaos_monkey
All 4 manifests must be run against your cluster to deploy the program. <br>

**deployment** <br>
deployment definition for the app.
Requires updating the `spec/template/spec/containers/image` value before deployment to work.

**config map** <br>
Config map contains the `pod_chaos_monkey.yml` used to feed the application with configuration data. 
Using this method you do not need to rebuild the app to update its function. 

**cluster role & cluster role binding**
provide permissions across the cluster to interact with the kubernetes api.<br>


## nginx
For a complete test of the selector mechanism I have included 3 deployment config yaml.

 - destroyable-nginx-deployment-1
 - destroyable-nginx-deployment-2

These both contain the label chaosmonkey set to "true". Using Pod Chaos Monkeys default configuration
defined in the pod_chaos_monkey config map. The pods in these deployments will be targets by the program for deletion. 

- indestructible-nginx-deployment 

This deployment yaml does not contain a reference to the chaosmonkey label and therefor is not targeted by the program. 

Note: you can also deploy a destroyable-nginx deployment into a different namespace to prove that the application honours namespacing. 

