# mackerel-plugin-minecraft

Minecraft custom metrics plugin for mackerel.io agent.

## Install

```bash
$ mkr plugin install natsuneko-laboratory/mackerel-plugin-minecraft
```

## Setting

```toml
[plugin.metrics.minecraft]
command = "/path/to/mackerel-plugin-minecraft"
```

## Example Metrics

```bash
$ mackerel-plugin-minecraft -server=localhost:25565
minecraft.max_players	20	1722694829
minecraft.online_players	2	1722694829
minecraft.latency	4	1722694829
```

## License

MIT by [@6jz](https://twitter.com/6jz)
