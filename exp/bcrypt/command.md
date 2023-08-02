```bash
# Produces binary named "bcrypt.exe"
go build bcrypt.go

# Hash a password
bcrypt.exe hash "some password here"

# Compare a hash and a password
bcrypt.exe bcrypt.go compare

```
