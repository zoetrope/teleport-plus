
Deploy teleport-plus container as a sidecar container.

```yaml
      containers:
      - name: teleport-plus
        image: teleport-plus:v1
        imagePullPolicy: Always
        command:
        - /manager
        args:
        - --enable-leader-election
        volumeMounts:
        - mountPath: /etc/teleport
          name: teleport-config
          readOnly: true
        - mountPath: /var/lib/teleport
          name: teleport-storage
```
