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
- [ ] Real-time messaging

## :rocket: Runing the server in development mode

1. Start the database and the Redis server using Docker:

```bash
docker compose up
```

2. Run the dispatch server:

```bash
make dispatch
```

3. Run the notification server:

```bash
make notification
```

4. Run the switchboard server:

```bash
make switchboard
```

## :rotating_light: Disclaimer

> [!WARNING]
> This project is a personal hobby project and is not affiliated with, endorsed by, or officially supported by Microsoft or any of its subsidiaries. The content, views, and contributions in this repository are solely those of the project author(s) and do not represent the opinions or positions of Microsoft.
