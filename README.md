# Teleport Plus

Teleport Plus enhances [Teleport][] to be able to register teleport resources as Custom Resource on Kubernetes.

Supported environments
----------------------

- Kubernetes
    - 1.15
- Teleport
    - 4.0.x

Features
--------

- Create/update/delete [Teleport Resources][] via custom resource on the Kubernetes cluster.

Usage
-----

### Deployment

Teleport-plus needs to be deployed as a sidecar container of the teleport-auth container.

In addition, the teleport-plus container must share teleport configuration and teleport storage with the teleport-auth container.

If you already have a manifest of teleport, you just add the following container next to the teleport-auth container.

```yaml
      containers:
      - name: teleport-plus
        image: teleport-plus:v1
        volumeMounts:
        - mountPath: /etc/teleport
          name: teleport-config
          readOnly: true
        - mountPath: /var/lib/teleport
          name: teleport-storage
```

See the [sample manifest](./e2e/teleport.yaml) for details.

### Apply custom resource

You can use a `TeleportResource` which is a custom resource on Kubernetes cluster to register resources of teleport.

An example of `TeleportResource` is shown below.

```yaml
apiVersion: teleport.gravitational.com/v1
kind: TeleportResource
metadata:
  name: github-integration
  namespace: teleport
spec:
  data: |
    kind: github
    version: v3
    metadata:
      name: github
    spec:
      client_id: <client-id>
      client_secret: <client-secret>
      display: Github
      redirect_url: https://<proxy-address>/v1/webapi/github/callback
      teams_to_logins:
        - organization: octocats
          team: admins
          logins:
            - root
          kubernetes_groups: ["system:masters"]
```

The contents specified in `.spec.data` will be registered to teleport.

The namespace must be the same namespace as teleport-plus container.

Getting Started
---------------

You can try teleport-plus on [kind][] by running the following command.

```
make e2e
```

[Teleport]: https://github.com/gravitational/teleport
[Teleport Resources]: https://gravitational.com/teleport/docs/admin-guide/#resources
[kind]: https://github.com/kubernetes-sigs/kind
