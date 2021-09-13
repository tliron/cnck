CNCK
====

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Latest Release](https://img.shields.io/github/release/tliron/cnck.svg)](https://github.com/tliron/cnck/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/cnck)](https://goreportcard.com/report/github.com/tliron/cnck)

CNCK = Cloud Native Configurations for Kubernetes

Make your Kubernetes applications more cloud native by injecting runtime cluster information into your
ConfigMaps.

CNCK is a Kubernetes operator that renders text templates with simple JavaScript scriptlets that can query
Kubernetes resources to pull data and generate contextual configuration text. It can continuously keep your
configurations up to date. Couple it with the [Reloader operator](https://github.com/stakater/Reloader) that
will make sure to restart your components when the configurations change.


How It Works
------------

An example of a ConfigMap with an embedded scriptlet to gather all running database pods and configure a
loadbalancing connection to them:

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: myapp
  annotations:
    cnck.github.com/render: myapp.yaml.template
    cnck.github.com/refresh: 1m
data:
  myapp.yaml.template: |
    log-file: /var/log/myapp.log
    database-url: <%
        let addresses = [];
        let databases = k8s.select({kind: 'Pod', labels: {app: 'mariadb'}});
        for (let d = 0; d < databases.length; d++)
            addresses.push(databases[d].status.podIP);
        write('jdbc:mariadb:loadbalance://' + addresses.join(',') + '/mydb');
    %>
```

1. CNCK will detect the annotation and render the "myapp.yaml.template" data.
2. The scriptlet embedded in the `<%` and `%>` delimiters will be run.
3. The scriptlet selects pods according to labels, aggregates their IP addresses, and writes a JDBC
   connection string.
4. CNCK will set the result in a new data key, without the `.template` extension, so the final
   result could look something like this:

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: myapp
  annotations:
    cnck.github.com/render: myapp.yaml.template
    cnck.github.com/refresh: 1m
data:
  myapp.yaml: |
    log-file: /var/log/myapp.log
    database-url: jdbc:mariadb:loadbalance://10.244.0.54,10.244.0.76/mydb
  myapp.yaml.template: |
    ...
```

5. The "refresh" annotation is set to 1 minute, so CNCK will re-render the template at that interval.
   It will only update the ConfigMap if the data has changed.
6. An example of a Deployment that mounts this ConfigMap and has a
   [Reloader](https://github.com/stakater/Reloader) annotation:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  annotations:
    reloader.stakater.com/auto: "true"
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: myapp
  template:
    metadata:
      labels:
        app.kubernetes.io/name: myapp
    spec:
      containers:
      - name: main
        image: myimage
        volumeMounts:
        - mountPath: /etc/myapp # will have a "myapp.yaml" file
          name: config
      volumes:
      - name: config
        configMap:
          name: myapp # the name of the ConfigMap to mount
```

See the [examples directory](examples/) for more information.


Installation
------------

You can install the operator using [this manifest](assets/kubernetes/cnck-operator.yaml) as a template.
For example:

    curl -s https://raw.githubusercontent.com/tliron/cnck/main/assets/kubernetes/cnck-operator.yaml | NAMESPACE=default VERSION=1.0 envsubst | kubectl apply -f -
