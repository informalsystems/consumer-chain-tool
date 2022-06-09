#!/bin/bash
set -eux

#bash verify_proposal.sh $HOME/wasm_contracts wasm wasm1243cuuy98lxaf7ufgav0w76xt5es93afr8a3ya wasmd_consumer wasmd $HOME/tool_output_step2/$(date +"%Y-%m-%d_%H-%M-%S") 1 "tcp://localhost:26657" providerd

TOOL_INPUT="$1"
CONSUMER_CHAIN_ID="$2"
CONSUMER_CHAIN_MULTISIG_ADDRESS="$3"
CONSUMER_CHAIN_BINARY="$4"
WASM_BINARY="$5"
TOOL_OUTPUT="$6"
PROPOSAL_ID="$7"
PROVIDER_NODE_ID="$8"
PROVIDER_BINARY="$9"

# Create directories if they don't exist.
mkdir -p $TOOL_OUTPUT

bash prepare_proposal_inputs.sh $TOOL_INPUT $CONSUMER_CHAIN_ID $CONSUMER_CHAIN_MULTISIG_ADDRESS $CONSUMER_CHAIN_BINARY $WASM_BINARY $TOOL_OUTPUT

GENESIS_HASH=$(jq -r ".genesis_hash" $TOOL_OUTPUT/sha256hashes.json)
BINARY_HASH=$(jq -r ".binary_hash" $TOOL_OUTPUT/sha256hashes.json)

# Query the proposal to get the hashes from the chain
./$PROVIDER_BINARY q gov proposal $PROPOSAL_ID --node $PROVIDER_NODE_ID --output json > $TOOL_OUTPUT/proposal_info.json 2>&1

GENESIS_HASH_ON_CHAIN=$(jq -r ".content.genesis_hash" $TOOL_OUTPUT/proposal_info.json)
BINARY_HASH_ON_CHAIN=$(jq -r ".content.binary_hash" $TOOL_OUTPUT/proposal_info.json)

rm -f $TOOL_OUTPUT/proposal_info.json

if [ "$GENESIS_HASH" != "$GENESIS_HASH_ON_CHAIN" ] || [ "$BINARY_HASH" != "$BINARY_HASH_ON_CHAIN" ]
then
  echo "Recalculated genesis and binary hashes don't match the ones from the proposal!"
  exit 1
else
  echo "Genesis and binary hashes are correct!"
fi

