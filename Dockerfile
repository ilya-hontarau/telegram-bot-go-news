FROM golang:1.12.6 as build
ADD . /telegram-bot-go-news/
WORKDIR /telegram-bot-go-news/
RUN make

FROM scratch
COPY --from=build /telegram-bot-go-news/bin/gonews .
CMD ["./gonews"]
