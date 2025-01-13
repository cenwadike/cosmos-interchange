# interchange
**interchange** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

## Get started

### Run mars chain

```
ignite chain serve -c mars.yml
```

### Run venus chain

```
ignite chain serve -c venus.yml
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

### Run Hermes Relayer
```
ignite relayer hermes configure "mars" "http://localhost:26657" "http://localhost:9090" "venus" "http://localhost:26659" "http://localhost:9092" --chain-a-faucet "http://0.0.0.0:4500" --chain-b-faucet "http://0.0.0.0:4501" --chain-a-port-id "dex" --chain-b-port-id "dex" --channel-version "dex-1"
```

```
cp ./config.toml ~/.hermes/config.toml
```

```
hermes create channel --a-chain mars --b-chain venus --a-port dex --b-port dex --order ordered --chan-version "dex-1" --new-client-connection
```

```
hermes query channels --chain mars
```

```
hermes query channels --chain venus
```

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend

Additionally, Ignite CLI offers both Vue and React options for frontend scaffolding:

For a Vue frontend, use: `ignite scaffold vue`
For a React frontend, use: `ignite scaffold react`
These commands can be run within your scaffolded blockchain project. 


For more information see the [monorepo for Ignite front-end development](https://github.com/ignite/web).

## Release
To release a new version of your blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

### Install
To install the latest version of your blockchain node's binary, execute the following command on your machine:

```
curl https://get.ignite.com/username/interchange@latest! | sudo bash
```
`username/interchange` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/allinbits/starport-installer).

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.gg/ignite)
