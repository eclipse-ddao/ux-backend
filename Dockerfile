FROM golang:1.19-alpine as build
WORKDIR /src 
# Copy the two important files
COPY go.mod go.sum ./
COPY . ./
# Download all the deps
# Build the executable
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build  -o /exe 
# # Run the executable
# CMD [ "/exe" ]

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /exe /exe

USER nonroot:nonroot

ENTRYPOINT ["/exe"]


