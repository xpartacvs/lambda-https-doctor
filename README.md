# Lambda HTTPS Doctor
A https health doctor which run on top of AWS Lambda

## How to

### Get the things ready

```bash
# Clone the repository
git clone git@github.com:xpartacvs/lambda-https-doctor.git

# Change directory
cd lambda-https-doctor

# Download the dependencies
go mod tidy

# Compile to binary
GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./https-doctor main.go

# Compress the binary
zip https-doctor.zip https-doctor
```

## Environment Variables

I was designed to get configured by environment variables.

| **Environment Variable** | **Type**  | **Req** | **Default**                           | **Description**                                                                                                                                                    |
| :---                     | :---      | :---:   | :---                                  | :---                                                                                                                                                               |
| `DISHOOK_URL`            | `string`  | √       |                                       | Discord webhook's URL.                                                                                                                                             |
| `DISHOOK_BOT_NAME`       | `string`  |         | `HTTPS Doctor`                        | Discord webhook bot's display name.                                                                                                                                |
| `DISHOOK_BOT_AVATAR`     | `string`  |         | default URL                           | Discord webhook bot's avatar URL.                                                                                                                                  |
| `DISHOOK_BOT_MESSAGE`    | `string`  |         | `Your HTTPS health monitoring result` | Discord webhook bot's alert message.                                                                                                                               |
| `GRACEPERIOD`            | `integer` |         | `14`                                  | Number of days before the host's SSL certificate get expired. When the current time was in the range, alert will get triggered.                                    |
| `HOSTS`                  | `string`  | √       | empty string                          | Coma separated list of hosts to check on.                                                                                                                          |
| `LOGLEVEL`               | `string`  |         | `disabled`                            | The logging mode: `debug`, `info`, `warn`, `error`, and `disabled`.                                                                                                |
| `TZ`                     | `string`  |         | local system                          | The timezone. Must contain one of [IANA Time Zone database](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) associate to your preferred time format. |
