# MSN Server (MSNP2)

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

## :memo: Description

This project is an implementation written in Golang of Microsoft's MSN Messenger server in its initial release. Due to the complexicity of the project, the server only supports (for now) MSN Messenger 1.0 client based on MSN Protocol 2 (MNSP2).

## :sparkles: Features

The project aims to support the entire feature set of MSNP2, including:

- [x] User authentication
- [x] Contact list management
- [x] Status updates
- [x] Real-time messaging

## :rocket: Runing the server in release mode

> [!IMPORTANT]
> Before running the server, you must have 3 network interfaces available on your machine and reachable by the clients. Those interfaces may be virtual. The example configuration in the repository uses the following IP ranges:
>
> - `192.168.101.0/24`
> - `192.168.102.0/24`
> - `192.168.103.0/24`

1. Add the configuration:

```bash
cp config.yaml.template config.yaml
```

2. Update the configuration file with the correct values.

3. Start the docker containers:

```bash
docker compose up
```

## :hammer: Runing the server in development mode

> [!IMPORTANT]
> The same requirements regarding the network interfaces as in the release mode apply to the development mode.

1. Add the configuration:

```bash
cp config-dev.yaml.template config.yaml
```

2. Update the configuration file with the correct values.

3. Start the database and the Redis server using Docker:

```bash
docker compose -f docker-compose-dev.yaml up
```

4. Run the dispatch server:

```bash
make dispatch
```

5. Run the notification server:

```bash
make notification
```

6. Run the switchboard server:

```bash
make switchboard
```

## :rotating_light: Disclaimer

> [!WARNING]
> This project is a personal hobby project and is not affiliated with, endorsed by, or officially supported by Microsoft or any of its subsidiaries. The content, views, and contributions in this repository are solely those of the project author(s) and do not represent the opinions or positions of Microsoft.
