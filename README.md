# cosmos-transactions-bot

![Latest release](https://img.shields.io/github/v/release/QuokkaStake/cosmos-transactions-bot)
[![Actions Status](https://github.com/QuokkaStake/cosmos-transactions-bot/workflows/test/badge.svg)](https://github.com/QuokkaStake/cosmos-transactions-bot/actions)
[![codecov](https://codecov.io/gh/QuokkaStake/cosmos-transactions-bot/graph/badge.svg?token=NDKDV02PC1)](https://codecov.io/gh/QuokkaStake/cosmos-transactions-bot)

cosmos-transactions-bot is a tool that listens to transactions with a specific filter on multiple chains
and reports them to a Telegram channel.

Here's how it may look like:

![Telegram](https://raw.githubusercontent.com/QuokkaStake/cosmos-transactions-bot/main/images/telegram.png)

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

You can use this template (change the user to whatever user you want this to be executed from. It's advised
to create a separate user for that instead of running it from root):

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

There are multiple nodes this app is connecting to via Websockets (see [this](https://docs.tendermint.com/master/rpc/#/Websocket/subscribe) for more details) and subscribing
to the queries that are set in config. When a new transaction matching the filters is found, it's put through
a deduplication filter first, to make sure we don't send the same transaction twice. Then each message
in transaction is enriched (for example, if someone claims rewards, the app fetches Coingecko price
and validator rewards are claimed from). Lastly, each of these transactions are sent to a reporter
(currently Telegram only) to notify those who need it.

## How can I configure it?

All configuration is done with a `.toml` file, which is passed to an app through a `--config` flag.
See `config.example.toml` for reference.

### Chains, subscriptions, chain subscriptions and reporters

This app's design is quite complex to allow it to be as flexible as possible.
There are the main objects that this app has:

- reporter - something that acts as a destination point (e.g. Telegram bot) and maybe allows you as a user
to interact with it in a special way (like, setting aliases etc.)
- chain - info about chain itself, its denoms, queries (see below), nodes used to receive data from, etc.
- subscription - info about which set of chains and their events to send to which reporter,
has many chain subscriptions
- chain subscription - info about which chain to receive data from, filters on which events to match
(see below) and how to process errors/unparsed/unsupported messages, if any.

Each chain has many chain subscriptions, each subscription has one reporter, each chain subscription
has one chain and many filters.

Generally speaking, the workflow of the app looks something like this:

![Schema](https://raw.githubusercontent.com/QuokkaStake/cosmos-transactions-bot/main/images/schema.png)

This allows to build very flexible setups. Here's the example of the easy and the more difficult setup.

1) "I want to receive all transactions sent from my wallet on chain A, B and C to my Telegram channel"

You can do it the following way:
- have 1 reporter, a Telegram channel
- have 3 chains, A, B and C, and their configs
- have 1 subscription, with Telegram reporter and 3 chain subscriptions inside (one for chain A, B and C
with 1 filter each matching transfers from wallets on these chains)

2) "I want to receive all transactions sent from my wallet on chains A, B and C to one Telegram chat,
all transactions that are votes on chains A and B to another Telegram chat, and all transactions that are delegations
with amount more than 10M $TOKEN on chain C to another Telegram chat"

That's also manageable. You can do the following:
- reporter 1, "first", a bot that sends messages to Telegram channel 1
- reporter 2, "second", a bot that sends messages to Telegram channel 2
- reporter 3, "third", a bot that sends messages to Telegram channel 3
- chain A and its config
- chain B and its config
- chain C and its config
- subscription 1, let's call it "my-wallet-sends", with reporter "first" and the following chain subscriptions
- - chain subscription 1, chain A, 1 filter matching transfers from my wallet on chain A
- - chain subscription 2, chain B, 1 filter matching transfers from my wallet on chain B
- - chain subscription 3, chain C, 1 filter matching transfers from my wallet on chain C
- subscription 2, let's call it "all-wallet-votes", with reporter "second" and the following chain subscriptions
- - chain subscription 1, chain A, 1 filter matching any vote on chain A
- - chain subscription 2, chain B, 1 filter matching any vote on chain B
- subscription 3, let's call it "whale-votes", with reporter "third" and the following chain subscription
- - chain subscription 1, chain C, 1 filter matching any delegations with amount more than 10M $TOKEN on chain C

See config.example.toml for real-life examples.

### Queries and filters

This is another quite complex topic and deserves a special explanation.

When a node starts, it connects to a Websocket of the fullnode and subscribes to queries (`queries` in `.toml` config).
If there's a transaction that does not match these filters, a fullnode won't emit the event for it
and this transaction won't reach the app.

If using filters (`filters` in `.toml` config), when a transaction is received, all messages in the transaction
are checked whether they match these filters, and can be filtered out (and the transaction itself would be filtered out
if there are 0 non-filtered messages left).

Using filters can be useful is you have transactions with multiple messages, where you only need to know about one
(for example, someone claiming rewards from your validator and other ones, when you need to know only about claiming
from your validator).

Keep in mind that queries is set on the app level, while filters are set on a chain subscription level,
so you can have some generic query on a chain, and more granular filter on each of your chain subscriptions.

Filters should follow the same pattern as queries, but they can only match the following pattern (so no AND/OR support):
- `xxx = yyy` (which would filter the transaction if key doesn't match value)
- `xxx! = yyy` (which would filter the transaction if key does match value)

Please note that the message would not be filtered out if it matches at least one filter.
Example: you have a message that has `xxx = yyy` as events, and if using `xxx != yyy` and `xxx != zzz` as filters,
it won't get filtered out (as it would not match the first filter but would match the second one).

You can always use `tx.height > 0`, which will send you the information on all transactions in chain,
or check out something we have:


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

One important thing to keep in mind: by default, Tendermint RPC now only allows 5 connections per client,
so if you have more than 5 filters specified, this will fail when subscribing to 6th one.
If you own the node you are subscribing to, o fix this, change this parameter to something that suits your needs
in `<fullnode folder>/config/config.toml`:

```
max_subscriptions_per_client = 5
```

### Denoms fetching

The app fetches denoms and their prices in the following order:
1. Local chain denoms
2. If it's IBC denom (`ibc/xxxxx`):
- it traverses IBC path,
- it fetches all intermediate chains, if we have them in local config
- when getting a final chain, it tries to get its local config denom
- if there's no local config, or denom in it, it takes data from https://cosmos.directory by chain-id
and denom from there, if found
3. If it's not an IBC denom:
- it fetches the https://cosmos.directory chain by chain-id
- it takes the denom from there, if found.

Consider this config:
```
[chain]
name = "osmosis"
chain-id = "osmosis-1"
denoms = [
    { denom = "uosmo", display-denom = "osmo", coingecko-currency = "osmosis" },
    { denom = "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2", display-denom = "atom", coingecko-currency = "cosmos" }
]

[chain]
name = "cosmoshub"
chain-id = "cosmoshub-4"
denoms = [
    { denom = "uatom", display-denom = "atom", coingecko-currency = "cosmos" },
]

[chain]
name = "akash"
chain-id = "akashnet-2"
denoms = []
```

and the following transactions:
1. IBC transfer from Osmosis (chain-id `osmosis-1`) to Cosmos Hub (chain-id `cosmoshub-4`) of `100uosmo` - it 
will take the denom from local config of `osmosis-1` (as `uosmo` denom is declared there).
2. IBC transfer from Osmosis (chain-id `osmosis-1`) to Cosmos Hub (chain-id `cosmoshub-4`)
of `100ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2` (`uatom` on osmosis-1 chain) - it will take
the denom from local config of `osmosis-1` (`ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2`).
3. IBC transfer from Cosmos Hub (chain-id `cosmoshub-4`) to Akash (chain-id `akashnet-2`) of
`100ibc/1480B8FD20AD5FCAE81EA87584D269547DD4D436843C1D20F15E00EB64743EF4` (`uakt` transferred from Akash to Osmosis)
- it will take the denom from cosmos.directory (as there's no `uakt` denom declared in local config of `akashnet-2`,
or `ibc/1480B8FD20AD5FCAE81EA87584D269547DD4D436843C1D20F15E00EB64743EF4` denom declared on `osmosis-1` chain).
4. Withdraw delegator rewards on Cosmos Hub (chain-id `cosmoshub-4`) in `uatom` - it will take data from local config,
as `uatom` is declared as denom there.
5. Withdraw delegator rewards on Akash (chain-id `akashnet-2`) in `uakt` - it will take data from cosmos.directory,
as `uakt` is not declared as denom in `akashnet-2` config.

Generally, the rule is the following:
- if you want to override how some tokens are displayed - override them in your local config.
- if you do not want to deal with it - just omit specifying them, it should do the job for you.


## Notifications channels

Go to [@BotFather](https://t.me/BotFather) in Telegram and create a bot. After that, there are two options:
- you want to send messages to a user. This user should write a message to [@getmyid_bot](https://t.me/getmyid_bot),
then copy the `Your user ID` number. Also keep in mind that the bot won't be able to send messages
unless you contact it first, so write a message to a bot before proceeding.
- you want to send messages to a channel. Write something to a channel, then forward it to [@getmyid_bot](https://t.me/getmyid_bot)
and copy the `Forwarded from chat` number. Then add the bot as an admin.

Then run a program with Telegram config (see `config.example.toml` as example).

You would likely want to also put only the IDs of trusted people to admins list in Telegram config, so the bot
won't react to anyone writing messages to it except these users.

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
