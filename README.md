### Cadence + Descript Overdub API
This repo use the Descript api to genereate a long audio file using cadence as workflow orchestrator

#### Requirements
```
golang
sox
docker-compose
```

#### Install
Start cadence and dependency
```bash
$ docker-compose up
```
Create domain
```bash
$ docker run --rm ubercadence/cli:master --address host.docker.internal:7933 --domain samples-domain domain register
```


#### Workflows
Overdub workflow allow you to generate audio book using overdub API

spliting text ---> started overdub request for each chunk -----> waiting overdub to complete in parallel ---> download audio and merge them in one wav file.  

#### Test
- Change voice_id and token in services/overdub/main.go
- run the backend
- go run main.go
- sh testOverdub.sh
- Check in temp folder
#### UI
localhost:8088