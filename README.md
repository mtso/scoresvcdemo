# scoresvcdemo

Demo microservice that copies the `github.com/go-kit/go/examples/profilesvc` service architecture.

## Process Roadmap
Builds on lessons learned from ulog and addsvcdemo

1. Create Service interface
2. Build endpoints
3. Add transports
   ==wip==
4. Add middlewares
5. Combine into cmd/scoresvcdemo/main.go

6. Turn inmemService into a postgreService (convert in-memory persistence to database)