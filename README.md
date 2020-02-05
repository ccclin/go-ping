# GO Ping
Use GAE to ping GCP Internal DNS Server

## How to use
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
