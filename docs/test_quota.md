### Documentação de Teste: Verificação de Quota de Recursos no Kubernetes

Este documento descreve o processo de teste para verificar se a quota de recursos definida para um namespace no Kubernetes está funcionando corretamente. O teste envolve a criação de um Deployment que solicita mais recursos do que os permitidos pela quota.

#### Passos para Definir a Quota de Recursos

1. **Definir a Quota de Recursos para o Namespace**

   Execute o comando abaixo para definir a quota de recursos para o namespace `test`:

   ```sh
   ./kubeon namespace-quotas test --action=set --limits-cpu=1 --limits-memory=1Gi --requests-cpu=1 --requests-memory=1Gi
   ```

2. **Verificar a Quota de Recursos Definida**

   Verifique se a quota foi definida corretamente:

   ```sh
   kubectl get resourcequota -n test
   ```

#### Passos para Criar um Deployment que Excede a Quota

1. **Criar o Manifesto do Deployment**

   Crie um arquivo chamado `deployment-exceed-quota.yaml` com o seguinte conteúdo:

   ```yaml
   apiVersion: apps/v1
   kind: Deployment
   metadata:
     name: exceed-quota-deployment
     namespace: test
   spec:
     replicas: 1
     selector:
       matchLabels:
         app: exceed-quota
     template:
       metadata:
         labels:
           app: exceed-quota
       spec:
         containers:
         - name: exceed-quota-container
           image: nginx
           resources:
             requests:
               cpu: "5"         # Excede a quota de requests.cpu
               memory: "5Gi"    # Excede a quota de requests.memory
             limits:
               cpu: "5"         # Excede a quota de limits.cpu
               memory: "5Gi"    # Excede a quota de limits.memory
   ```

2. **Aplicar o Manifesto do Deployment**

   Execute o comando abaixo para aplicar o manifesto do Deployment:

   ```sh
   kubectl apply -f deployment-exceed-quota.yaml
   ```

3. **Verificar o Status do Deployment**

   Verifique o status do Deployment para confirmar que ele não foi criado devido à quota de recursos excedida:

   ```sh
   kubectl get deploy -n test
   ```

   Você deve ver que o Deployment não possui réplicas disponíveis.

4. **Descrever o Deployment para Verificar os Eventos**

   Descreva o Deployment para verificar os eventos e confirmar que a criação falhou devido à quota de recursos:

   ```sh
   kubectl describe deploy exceed-quota-deployment -n test
   ```

   A saída deve mostrar eventos indicando que a criação do Deployment falhou devido à quota de recursos excedida.

#### Exemplo de Saída

```sh
kubectl apply -f deployment-exceed-quota.yaml
deployment.apps/exceed-quota-deployment created

kubectl get po -n test
No resources found in test namespace.

kubectl get deploy -n test
NAME                      READY   UP-TO-DATE   AVAILABLE   AGE
exceed-quota-deployment   0/1     0            0           119s

kubectl describe deploy exceed-quota-deployment -n test
Name:                   exceed-quota-deployment
Namespace:              test
CreationTimestamp:      Mon, 27 Jan 2025 16:26:35 -0300
Labels:                 <none>
Annotations:            deployment.kubernetes.io/revision: 2
Selector:               app=exceed-quota
Replicas:               1 desired | 0 updated | 0 total | 0 available | 2 unavailable
StrategyType:           RollingUpdate
MinReadySeconds:        0
RollingUpdateStrategy:  25% max unavailable, 25% max surge
Pod Template:
  Labels:  app=exceed-quota
  Containers:
   exceed-quota-container:
    Image:      nginx
    Port:       <none>
    Host Port:  <none>
    Limits:
      cpu:     5
      memory:  5Gi
    Requests:
      cpu:         5
      memory:      5Gi
    Environment:   <none>
    Mounts:        <none>
  Volumes:         <none>
  Node-Selectors:  <none>
  Tolerations:     <none>
Conditions:
  Type             Status  Reason
  ----             ------  ------
  Progressing      True    NewReplicaSetCreated
  Available        False   MinimumReplicasUnavailable
  ReplicaFailure   True    FailedCreate
OldReplicaSets:    exceed-quota-deployment-77485f6cb (0/1 replicas created)
NewReplicaSet:     exceed-quota-deployment-669cfddd4 (0/1 replicas created)
Events:
  Type    Reason             Age    From                   Message
  ----    ------             ----   ----                   -------
  Normal  ScalingReplicaSet  2m22s  deployment-controller  Scaled up replica set exceed-quota-deployment-77485f6cb to 1
  Normal  ScalingReplicaSet  63s    deployment-controller  Scaled up replica set exceed-quota-deployment-669cfddd4 to 1
```

#### Conclusão

Este teste confirma que a quota de recursos definida para o namespace `test` está funcionando corretamente, impedindo a criação de um Deployment que solicita mais recursos do que os permitidos.