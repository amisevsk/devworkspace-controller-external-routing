# Example external workspacerouting controller for DevWorkspace operator

This project is a small controller sample that reuses the workspacerouting controller from the DevWorkspace operator to provide a controller for workspaceRoutings with the routingClass `external-sample`.

Bootstrapped by Operator SDK v1.1, with minimal updates from the raw templates.

## Testing
1. Deploy the [DevWorkspace operator](https://github.com/devfile/devworkspace-operator) as normal
2. Deploy this controller alongside it
    ```bash
    export IMG=<my-image>
    export NAMESPACE=<namespace> # Default: external-routing-example
    make docker-build docker-push deploy
    ```
3. Create DevWorkspace with external routingclass:
    ```bash
    kubectl apply -f samples/theia-external.yaml
    ```
4. Check for created service:
    ```
    $ kubectl get svc
    NAME              TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
    example-service   ClusterIP   10.96.203.255   <none>        9999/TCP   12s
    ```

## Makefile rules
- `make manager`: build binary locally
- `make deploy`: deploy controller to cluster
- `make logs`: get logs for controller
- `make restart`: restart controller deployment
- `make fmt vet`: run `go fmt` and `go vet`
- `make docker-build docker-push`: build and push `${IMG}`