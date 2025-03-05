# Laboratory work 2 on the discipline "Knowledge bases"

## Get started

1. ```shell
    go mod tidy
    ```

2. Create `.env` file in root directory with `CONFIG_PATH` field
3. Create `config.yaml` file in `/config` directory like:

    ```yaml
    storage_root: './storage/'
    rules: 'rules.txt'
    chart: 'charts.html'
    server_port: '5000'
    ```

## Launch the app

```shell
go run ./cmd/app/main.go
```
