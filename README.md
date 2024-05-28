# Archived Warning
A lot of known security vulnerabilities of outdated packages detected.

If you are interested in this project, please upgrade the dependencies.

# Linux 集群管理及监控系统
1. Start components: `docker-compose up -d`
2. Build: `make all`
3. Start server: `./bin/server`
4. Start client: `./bin/client`

Client 提供一个 REPL，可以进行命令补全、解析、执行：

![client](./res/启动client.png)

监控图表：

![stat_graph](./res/动态图表.png)

## TODO
- 基于snmp的监控
- 远程桌面 参考：[rust实现远程桌面](https://www.bilibili.com/video/BV1tZ4y1X7hT)
