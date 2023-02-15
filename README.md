# cosmos-transactions-bot

![Latest release](https://img.shields.io/github/v/release/QuokkaStake/cosmos-transactions-bot)
[![Actions Status](https://github.com/QuokkaStake/cosmos-transactions-bot/workflows/test/badge.svg)](https://github.com/QuokkaStake/cosmos-transactions-bot/actions)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FQuokkaStake%2Fcosmos-transactions-bot.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FQuokkaStake%2Fcosmos-transactions-bot?ref=badge_shield)

cosmos-transactions-bot is a tool that listens to transactions with a specific filter on multiple chains and reports them to a Telegram channel.

Here's how it may look like:

![Telegram](https://raw.githubusercontent.com/QuokkaStake/cosmos-transactions-bot/master/images/telegram.png)

## How can I set it up?

Download the latest release from [the releases page](https://github.com/QuokkaStake/cosmos-transactions-bot/releases/). After that, you should unzip it and you are ready to go:

```sh
wget <the link from the releases page>
tar xvfz <filename you have just downloaded>
./cosmos-transactions-bot <params>
```

To have it running in the background, first, we have to copy the file to the system apps folder:

```sh
sudo cp ./cosmos-transactions-bot /usr/bin
```

Then we need to create a systemd service for our app:

```sh
sudo nano /etc/systemd/system/cosmos-transactions-bot.service
```

You can use this template (change the user to whatever user you want this to be executed from. It's advised to create a separate user for that instead of running it from root):

```
[Unit]
Description=Cosmos Transactions Bot
After=network-online.target

[Service]
User=<username>
TimeoutStartSec=0
CPUWeight=95
IOWeight=95
ExecStart=cosmos-transactions-bot --config /path/to/config.toml
Restart=always
RestartSec=2
LimitNOFILE=800000
KillSignal=SIGTERM

[Install]
WantedBy=multi-user.target
```

Then we'll add this service to the autostart and run it:

```sh
sudo systemctl enable cosmos-transactions-bot
sudo systemctl start cosmos-transactions-bot
sudo systemctl status cosmos-transactions-bot # validate it's running
```

If you need to, you can also see the logs of the process:

```sh
sudo journalctl -u cosmos-transactions-bot -f --output cat
```

## How does it work?

There are multiple nodes this app is connecting to via Websockets (see [this](https://docs.tendermint.com/master/rpc/#/Websocket/subscribe) for more details) and subscribing to the queries that are set in config. When a new transaction matching the filters is found, it's put through a deduplication filter first, to make sure we don't send the same transaction twice. Then each message in transaction is enriched (for example, if someone claims rewards, the app fetches Coingecko price and validator rewards are claimed from). Lastly, each of these transactions are sent to a reporter (currently Telegram only) to notify those who need it.

## How can I configure it?

All configuration is done with a `.toml` file, which is passed to an app through a `--config` flag. See `config.example.toml` for reference.

### Queries and filters

This is quite complex and deserves a special explanation.

When a node starts, it connects to a Websocket of the fullnode and subscribes to queries (`queries` in `.toml` config). If there's a transaction that does not match these filters, a fullnode won't emit the event for it and this transaction won't reach the app.

If using filters (`filters` in `.toml` config), when a transaction is received, all messages in the transaction are checked whether they match these filters, and can be filtered out (and the transaction itself would be filtered out if there are 0 non filtered messages left).

Using filters can be useful is you have transactions with multiple messages, where you only need to know about one (for example, someone claiming rewards from your validator and other ones, when you need to know only about claiming from your validator).

Filters should follow the same pattern as queries, but they can only match the following pattern (so no AND/OR support):
- `xxx = yyy` (which would filter the transaction if key doesn't match value)
- `xxx! = yyy` (which would filter the transaction if key does match value)

Please note that the message would not be filtered out if it matches at least one filter. Example: you have a message that has `xxx = yyy` as events, and if using `xxx != yyy` and `xxx != zzz` as filters, it won't get filtered out (as it would not match the first filter but would match the second one).

You can always use `tx.height > 0`, which will send you the information on all transactions in chain, or check out something we have:


```
queries = [
    # claiming rewards from validator's wallet
    "withdraw_rewards.validator = 'sentvaloper1rw9wtyhsus7jvx55v3qv5nzun054ma6kz4237k'",
    # incoming delegations from validator
    "delegate.validator = 'sentvaloper1rw9wtyhsus7jvx55v3qv5nzun054ma6kz4237k'",
    # redelegations from and to validator
    "redelegate.source_validator = 'sentvaloper1rw9wtyhsus7jvx55v3qv5nzun054ma6kz4237k'",
    "redelegate.destination_validator = 'sentvaloper1rw9wtyhsus7jvx55v3qv5nzun054ma6kz4237k'",
    # unbonding from validator
    "unbond.validator = 'sentvaloper1rw9wtyhsus7jvx55v3qv5nzun054ma6kz4237k'",
    # tokens sent from validator's wallet
    "transfer.sender = 'sent1rw9wtyhsus7jvx55v3qv5nzun054ma6kas4u3l'",
    # tokens sent to validator's wallet
    "transfer.recipient = 'sent1rw9wtyhsus7jvx55v3qv5nzun054ma6kas4u3l'",
    # IBC token transferred from validator's wallet
    "ibc_transfer.sender = 'sent1rw9wtyhsus7jvx55v3qv5nzun054ma6kas4u3l'",
    # IBC token received at validator's wallet
    "fungible_token_packet.receiver = 'sent1rw9wtyhsus7jvx55v3qv5nzun054ma6kas4u3l'",
]
```

Or the similar, with filters (then we'll receive all transactions but would filter them on app level):

```
queries = ['tx.height > 1']
filters = [
    "withdraw_rewards.validator = 'sentvaloper1rw9wtyhsus7jvx55v3qv5nzun054ma6kz4237k'",
    "delegate.validator = 'sentvaloper1rw9wtyhsus7jvx55v3qv5nzun054ma6kz4237k'",
    "redelegate.source_validator = 'sentvaloper1rw9wtyhsus7jvx55v3qv5nzun054ma6kz4237k'",
    "redelegate.destination_validator = 'sentvaloper1rw9wtyhsus7jvx55v3qv5nzun054ma6kz4237k'",
    "unbond.validator = 'sentvaloper1rw9wtyhsus7jvx55v3qv5nzun054ma6kz4237k'",
    "transfer.sender = 'sent1rw9wtyhsus7jvx55v3qv5nzun054ma6kas4u3l'",
    "transfer.recipient = 'sent1rw9wtyhsus7jvx55v3qv5nzun054ma6kas4u3l'",
    "ibc_transfer.sender = 'sent1rw9wtyhsus7jvx55v3qv5nzun054ma6kas4u3l'",
    "fungible_token_packet.receiver = 'sent1rw9wtyhsus7jvx55v3qv5nzun054ma6kas4u3l'",
]
```

See [the documentation](https://docs.tendermint.com/master/rpc/#/Websocket/subscribe) for more information on queries.

One important thing to keep in mind: by default, Tendermint RPC now only allows 5 connections per client, so if you have more than 5 filters specified, this will fail when subscribing to 6th one. If you own the node you are subscribing to, o fix this, change this parameter to something that suits your needs in `<fullnode folder>/config/config.toml`:

```
max_subscriptions_per_client = 5
```

## Notifications channels

Go to [@BotFather](https://t.me/BotFather) in Telegram and create a bot. After that, there are two options:
- you want to send messages to a user. This user should write a message to [@getmyid_bot](https://t.me/getmyid_bot), then copy the `Your user ID` number. Also keep in mind that the bot won't be able to send messages unless you contact it first, so write a message to a bot before proceeding.
- you want to send messages to a channel. Write something to a channel, then forward it to [@getmyid_bot](https://t.me/getmyid_bot) and copy the `Forwarded from chat` number. Then add the bot as an admin.

Then run a program with Telegram config (see `config.example.toml` as example).

You would likely want to also put only the IDs of trusted people to admins list in Telegram config, so the bot won't react to anyone writing messages to it except these users.

Additionally, for the ease of using commands, you can put the following list as bot commands in @BotFather settings:

```
help - Display help message
status - Display nodes status
config - Display .toml config
alias - Add a wallet alias
aliases - List wallet aliases
```

## Which networks this is guaranteed to work?

In theory, it should work on a Cosmos-based blockchains that expose a Tendermint RPC endpoint.

## How can I contribute?

Bug reports and feature requests are always welcome! If you want to contribute, feel free to open issues or PRs.


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FQuokkaStake%2Fcosmos-transactions-bot.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FQuokkaStake%2Fcosmos-transactions-bot?ref=badge_large)