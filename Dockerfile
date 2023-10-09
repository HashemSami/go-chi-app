FROM golang

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

# run the server build
CMD ./server