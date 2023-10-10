
# will use the golang image just to build up the binary file
# then will send the binary file to the live server to run
FROM golang AS builder

WORKDIR /app

# adding the mod and the sum files first and download them
# so the next build can know that
# nothing has changed from the last
# build and all will be built faster.
# getting advantage from the caching
COPY go.mod go.sum ./

RUN go mod download

# then move everything to the app direcrory
COPY . .

# build the app and call the build "server"
# -v for listing the build process in the terminal
RUN go build -v -o ./server ./cmd/server

# ================================================
# build the server
FROM ubuntu

WORKDIR /app

# copy only the required files the binary will use
COPY ./assets ./assets
COPY .env .env

# copy the binary file
COPY --from=builder /app/server ./server

# run the server build
CMD ./server