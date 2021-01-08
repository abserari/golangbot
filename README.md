# golangbot
telegram or other platform to reply the golang program's process result.
https://www.yuque.com/abser/solutions

![Picture](https://cdn.nlark.com/yuque/0/2020/png/176280/1585055803036-c05e1d2c-9195-460f-b4f1-e2854d7b60d2.png)

![Picture](https://cdn.nlark.com/yuque/0/2021/png/176280/1610074773928-a471b068-823d-4ea6-88fd-6ae0312f6824.png)

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
