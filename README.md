# telegram-bot-go-news 

[![CircleCI](https://circleci.com/gh/illfate/telegram-bot-go-news.svg?style=svg)](https://circleci.com/gh/illfate/telegram-bot-go-news)
[![codecov](https://codecov.io/gh/illfate/telegram-bot-go-news/branch/master/graph/badge.svg)](https://codecov.io/gh/illfate/telegram-bot-go-news)

It's a telegram bot that sends latest news from RSS channel according to passed tag.

## Getting Started

### Prerequisites

You need to have:

- git
- go 1.12+

## Installing

Next commands will build app:
```bash
git clone https://github.com/illfate/telegram-bot-go-news && \
cd telegram-bot-go-news && \
make
```

Binary will be compiled in current repository in `bin` directory.

## Running the tests

You need to be in root of this repository and type:
```bash
make test
```
