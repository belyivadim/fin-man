FROM cosmtrek/air

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

EXPOSE 8080

CMD [ "air" ]

