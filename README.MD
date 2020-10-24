Порядок первичного старта приложения:
```bash
go mod download
prisma-client-go migrate --experimental up
prisma-client-go generate
go build
./addreality
```