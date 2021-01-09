# golangbot
telegram or other bot by golang.
https://www.yuque.com/abser/solutions

## Run Golang bot
![Picture](https://cdn.nlark.com/yuque/0/2020/png/176280/1585055803036-c05e1d2c-9195-460f-b4f1-e2854d7b60d2.png)

## View Pixiv Picture
支持翻页

![Picture](https://cdn.nlark.com/yuque/0/2021/png/176280/1610166356617-482bae82-3898-4c32-b68a-af8925ed5aa8.png?x-oss-process=image%2Fresize%2Cw_738)

## Usage

### Pre
- Docker or Go v1.15+
- TelegramBot Token @botfather
- PixivCookies or Username&Password [optional]

### Docker
```bash
docker run -v ~/config.yaml:/config.yaml yhyddr/golangbot
```

### Build
```bash
git clone https://github.com/abserari/golangbot

cd golangbot

# edit the config.yaml
go run main.go
```
