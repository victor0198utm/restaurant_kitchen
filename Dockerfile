FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /restaurant_kitchen

EXPOSE 8081

CMD [ "/restaurant_kitchen" ]


