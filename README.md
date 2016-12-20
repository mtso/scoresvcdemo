# scoresvcdemo

Demo microservice that copies the `github.com/go-kit/go/examples/profilesvc` service architecture.

## Sample curls

    curl -X POST -d '{"id":"kingcandy","score":20}' localhost:8080/
    curl -i localhost:8080/kingcandy


## Process Roadmap

Builds on lessons learned from ulog and addsvcdemo

1. Create Service interface
2. Build endpoints
3. Add transports
4. Add middlewares
5. Combine into cmd/scoresvcdemo/main.go
   ==wip==
0. Turn inmemService into a postgreService (convert in-memory persistence to database)