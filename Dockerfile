FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY /appData ./appData
COPY /models ./models

RUN go build -o /restaurant_kitchen

CMD [ "/restaurant_kitchen" ]


