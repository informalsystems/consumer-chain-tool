#!/bin/bash
set -eux 

TOTAL_COINS=100000000000stake
STAKE_COINS=100000000stake
PROVIDER_BINARY=providerd
PROVIDER_HOME=$HOME/.tool_provider
PROVIDER_CHAIN_ID=provider
PROVIDER_MONIKER=provider
VALIDATOR=validator

# Clean start
killall $PROVIDER_BINARY &> /dev/null || true
rm -rf $PROVIDER_HOME

./$PROVIDER_BINARY init $PROVIDER_MONIKER --home $PROVIDER_HOME --chain-id $PROVIDER_CHAIN_ID
jq ".app_state.gov.voting_params.voting_period = \"3s\" | .app_state.staking.params.unbonding_time = \"600s\"" \
   $PROVIDER_HOME/config/genesis.json > \
   $PROVIDER_HOME/edited_genesis.json && mv $PROVIDER_HOME/edited_genesis.json $PROVIDER_HOME/config/genesis.json
sleep 1

# Create account keypair
./$PROVIDER_BINARY keys add $VALIDATOR --home $PROVIDER_HOME --keyring-backend test --output json > $PROVIDER_HOME/keypair.json 2>&1
sleep 1

# Add stake to user
./$PROVIDER_BINARY add-genesis-account $(jq -r .address $PROVIDER_HOME/keypair.json) $TOTAL_COINS --home $PROVIDER_HOME --keyring-backend test
sleep 1

# Stake 1/1000 user's coins
./$PROVIDER_BINARY gentx $VALIDATOR $STAKE_COINS --chain-id $PROVIDER_CHAIN_ID --home $PROVIDER_HOME --keyring-backend test --moniker $VALIDATOR
sleep 1

./$PROVIDER_BINARY collect-gentxs --home $PROVIDER_HOME --gentx-dir $PROVIDER_HOME/config/gentx/
sleep 1

# Start the chain
./$PROVIDER_BINARY start --home $PROVIDER_HOME &> $PROVIDER_HOME/logs &
# TODO: Think about nicer way to make sure chain is up and running (producing block)
sleep 10

# Build consumer chain proposal file
tee $PROVIDER_HOME/consumer-proposal.json<<EOF
{
    "title": "Create a chain",
    "description": "Gonna be a great chain",
    "chain_id": "wasm",
    "initial_height": {
        "revision_height": 1,
        "revision_number": 1
    },
    "genesis_hash": "5c5a82f958621228e704c0a00bb591386c9f891f8bfadb1a34b4c15114174d99",
    "binary_hash": "bbe2de71aacd5af0d4a98118ede4911b7993f447a07c773f9c7c6fe7d2d005ca",
    "spawn_time": "2022-06-01T09:10:00.000000000-00:00", 
    "deposit": "10000001stake"
}
EOF

./$PROVIDER_BINARY tx gov submit-proposal create-consumer-chain $PROVIDER_HOME/consumer-proposal.json \
	--chain-id $PROVIDER_CHAIN_ID --from $VALIDATOR --home $PROVIDER_HOME --keyring-backend test -b block -y
sleep 1

# Vote yes to proposal
./$PROVIDER_BINARY tx gov vote 1 yes --from $VALIDATOR --chain-id $PROVIDER_CHAIN_ID --home $PROVIDER_HOME -b block -y --keyring-backend test
sleep 5
