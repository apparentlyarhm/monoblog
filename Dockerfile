FROM node:22.21.1-alpine3.21 as fe
WORKDIR /app

# first we set up dependencies
COPY package*.json ./
RUN npm install

# then we copy the rest of the files and build
COPY . .
RUN npm run build

# at this point the app is built which has the generated static html for blogs
# in astro, we will get a dist/ folder as output

# STAGE 2: go
# go is the choice here because its lightweight and can easily be expanded
# for more stuff later on
# if in the future we dont want complex use cases, we can switch to something like nginx
FROM golang:1.25.5-alpine as be
WORKDIR /app

COPY --from=fe /app/dist ./dist
COPY go.mod go.sum main.go ./

RUN go mod download

RUN go build -o server main.go

# since we did ember.FS in the go code, the binary has the static files embedded

# STAGE 3: final minimal image
FROM alpine:latest
WORKDIR /root/

# only the binary is needed
COPY --from=be /app/server .

EXPOSE 8080

CMD ["./server"]