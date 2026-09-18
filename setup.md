# Setup

### Project dependencies

- Go
- Asynq
- GORM
- Redis
- sqlite3 (dev)
- Postgresql (prod)
- Node

```bash
# Defines the std path for go and it's installs
go GOROOT
# Tells Go about you cwd
go GOPATH
# Where binaries are install after running go build
go GOBIN
# to ensure 'go install' works well at Go to $PATH
export PATH="$PATH:$(go env GOPATH)/bin"
# Very go appears in $PATH
echo $PATH

#verify redis starts and works
# on mac use brew services start redis
redis-cli ping
#you should see
PONG
```

On MacOS27> bind redis to on ip version (IPV4 OR IPV6)
Instead of bind 127.0.0.1 ::1
use bind 127.0.0.1

Down the file, set tcp_keepalive to 0 instead of the default 300

```bash
#run redis manually
redis-server
```

### Migrate:

Used for easy database migrations and health

```bash
brew install golang-migrate
```

> [!NOTE]
> You should know how to install Node.js come now

### Using blue to create the file structure

Install go-blueprint by using this command

```bash
go install github.com/melkeydev/go-blueprint@latesto

# run blueprint
go-blueprint create \
  --name booking-platform \
  --framework standard-library \
  --driver sqlite \
  --git commit \
  --advanced \
  --feature react \
  --feature tailwind

# Blueprint isn't going to do everything
mkdir -p \
 cmd/worker \
  internal/user \
  internal/booking \
  internal/calendar \
  internal/guest \
  internal/services \
  internal/payments \
  internal/notifications \
  internal/reporting \
  internal/audit \
  infrastructure \
  migrations \
  storage \
  docs/adr \
  docs/architecture
```
