FROM golang:latest

WORKDIR /app
   
COPY go.* ./

RUN go mod download
   
COPY . .

RUN go build -o cc.exe
   
ENTRYPOINT ["./cc.exe"]
