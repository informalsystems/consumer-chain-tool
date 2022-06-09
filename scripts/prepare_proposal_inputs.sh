#!/bin/bash
set -eux

WASM_CONTRACTS="$1"
CHAIN_ID="$2"
MULTISIG_ADDRESS="$3"
CONSUMER_BINARY="$4"
WASM_BINARY="$5"
TOOL_OUTPUT="$6"
MONIKER="moniker"
VALIDATOR="validator"
KEYRING="--keyring-backend test"
STAKE="100000000stake"
TX_FLAGS="--gas-adjustment 100 --gas auto"
NODE_IP="localhost"
WASM_HOME="$HOME/.tool_wasmd"
WASM_RPC_LADDR="$NODE_IP:26638"
WASM_GRPC_ADDR="$NODE_IP:9071"
CONSUMER_HOME="$HOME/.tool_consumer"
CONST_PASPHRASE="torch bargain math dinner van fabric fly crystal answer first crush fan soap moon scene number dial any silk kangaroo clarify empower awake fiscal"

#TODO: check if killall has efect, if not list processes and kill it with kill -9. do this for all kill calls
# Clean start
killall $WASM_BINARY &> /dev/null || true
rm -rf $WASM_HOME
rm -rf $CONSUMER_HOME

# Create directories if they don't exist.
mkdir -p $TOOL_OUTPUT

#################################### CONTRACT COMPILATION #######################
#TODO: add contract compilation, copy source code in tool output

#################################### CONSUMER ###################################
# Generate initial genesis file
./$CONSUMER_BINARY init $MONIKER --chain-id $CHAIN_ID --home $CONSUMER_HOME
sleep 1

#################################### WASMD #####################################
# Init wasm chain
./$WASM_BINARY init $MONIKER --chain-id $CHAIN_ID --home $WASM_HOME
sleep 1

#TODO: Create user account keypair. Must use --recover to always get the same keys since address is stored when wasm stores the code (creator property)
echo $CONST_PASPHRASE | ./$WASM_BINARY keys add $VALIDATOR $KEYRING --home $WASM_HOME --recover --output json > $WASM_HOME/validator_keypair.json 2>&1

# Add stake to user account
./$WASM_BINARY add-genesis-account $(jq -r .address $WASM_HOME/validator_keypair.json)  1000000000stake --home $WASM_HOME

# Generate gentx file
./$WASM_BINARY gentx $VALIDATOR $STAKE --chain-id $CHAIN_ID --home $WASM_HOME $KEYRING
sleep 1

# Add validator
./$WASM_BINARY collect-gentxs --home $WASM_HOME --gentx-dir $WASM_HOME/config/gentx/
sleep 1


sed -i -r "/node =/ s/= .*/= \"tcp:\/\/${WASM_RPC_LADDR}\"/" $WASM_HOME/config/client.toml

# Start the chain; TODO: should ports be passed in as arguments to the script? this is to ensure the chain will start even if the default ports are taken.
./$WASM_BINARY start \
        --rpc.laddr tcp://${WASM_RPC_LADDR} \
        --grpc.address ${WASM_GRPC_ADDR} \
        --address tcp://${NODE_IP}:26635 \
        --p2p.laddr tcp://${NODE_IP}:26636 \
        --grpc-web.enable=false \
        --home $WASM_HOME &> $WASM_HOME/logs &
# TODO: Think about nicer way to make sure chain is up and running (producing block)
sleep 10

#TODO: set permissions for contract instantiation
# Deploy contracts
for CONTRACT in "$WASM_CONTRACTS"/*.wasm; do
  ./$WASM_BINARY tx wasm store $CONTRACT --instantiate-only-address $MULTISIG_ADDRESS --from $VALIDATOR $KEYRING --chain-id $CHAIN_ID --home $WASM_HOME $TX_FLAGS -b block -y
done

#Stop the chain
killall ./$WASM_BINARY &> /dev/null || true
sleep 3

#Export genesis state
./$WASM_BINARY export --home $WASM_HOME > $WASM_HOME/exported_genesis.json
jq -s '.[0].app_state.wasm = .[1].app_state.wasm | .[0]' $CONSUMER_HOME/config/genesis.json $WASM_HOME/exported_genesis.json > $CONSUMER_HOME/genesis_wasm.json

#################################### TIDY GENESIS #####################################
# //TODO: set parameters of each module in CONSUMER_HOME/genesis_wasm.json

#TODO: genesis_time must be some constant value (which one?) each time the tool is run so that we get the same hash
jq '.genesis_time = "2022-01-01T00:00:00.000000000Z" | .' $CONSUMER_HOME/genesis_wasm.json > $CONSUMER_HOME/genesis_1.json

# Copy genesis to the output folder
cp "$CONSUMER_HOME/genesis_1.json" "$TOOL_OUTPUT/genesis.json"

# Calculate binary and genesis hashes
tee $TOOL_OUTPUT/sha256hashes.json<<EOF
{
    "genesis_hash": "$(sha256sum $TOOL_OUTPUT/genesis.json | cut -d " " -f 1)",
    "binary_hash": "$(sha256sum $CONSUMER_BINARY | cut -d " " -f 1)"
}
EOF

#################################### CLEAN UP #########################################
# Delete all generated data. TODO: check if some other directories should be deleted
rm -rf $WASM_HOME
rm -rf $CONSUMER_HOME

