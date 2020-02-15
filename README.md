# GO Ping
Use GAE to ping GCP Internal DNS Server

## How to use
### Local Test
```
$ INTERNAL_DNS=8.8.8.8 go run ./main.go
```

```
$ curl http://localhost:8080/ping
```

### GAE
- Create a GCP project
- Deploy
```
$ gcloud config set project <YOUR_GCP_PROJECT_ID>
$ gcloud app deploy app.yaml
```
- After deploy
```
$ curl https://<YOUR_GCP_PROJECT_ID>.appspot.com/ping
```
