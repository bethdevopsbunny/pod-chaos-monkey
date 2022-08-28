# PodChaosMonkey

Welcome to PodChaosMonkey an agent of chaos working from within your Kubernetes cluster to test its resilience.

![](asset/royalty-free-monkey.jpg)


## Build
### Docker
- Install Go if you haven't already. <br><br>

- Clone the program locally. <br>
```git clone [repo]``` <br><br>

- Enter the application directory <br>
```cd PodChaosMonkey```<br><br>

- Build and tag the program <br>
```docker build --tag [yourcontainerregistery/podchaosmonkey] .``` <br><br>

- Push your built container to your registry <br>
```docker image push [yourcontainerregistery/podchaosmonkey] ```<br><br>

### Kubernetes manifests
Note: details on these files can be found [here](kubernetes-yml/README.md)

- Update the [deployment.yml](kubernetes-yml/pod_chaos_monkey/deployment.yml) file, with the name of your public container registry. <br><br>

- Check over the [configmap.yml](kubernetes-yml/pod_chaos_monkey/configmap.yml) defaults to make sure you are happy. <br>
  - the default namespace is... default go figure. <br>
  - The defaults for the cron will delete a pod every minute. <br>
  - The defaults for the pod selector are configured for the nginx demo, feel free to update them. <br><br>
  - the defaults for the concurrent deletes is 2

## Deployment
Note: all development and testing has been done against [googles kubernetes engine](https://cloud.google.com/kubernetes-engine)
1. Deploy the nginx deployment manifests <br>
    ```kubectl apply -f kubernetes-yml/nginx``` <br> <br>
2. Deploy the pod_chaos_monkey manifests<br>
   ```kubectl apply -f kubernetes-yml/pod_chaos_monkey```<br><br>
3. I suggest using [k9s](https://k9scli.io/) to view in realtime Pod Chaos Monkey removing pods. Basic logs and error handling is viewable from stdout (using k9s you can view this by pressing enter on the Pod Chaos Monkey pod).


  

