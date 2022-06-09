#!/bin/bash
set -eux

#bash finalize_genesis.sh $HOME/wasm_contracts wasm wasm1243cuuy98lxaf7ufgav0w76xt5es93afr8a3ya $HOME/tool_output_step2 1 "tcp://localhost:26657"

WASM_CONTRACTS="$1" #TODO: once SC compiling is solved, it should point to the output of step1.sh where source code of SCs are stored
TOOL_INPUT="$HOME/tool_output_step1"
CONSUMER_CHAIN_ID="$2"
CONSUMER_CHAIN_MULTISIG_ADDRESS="$3"
TOOL_OUTPUT_DIR="$4"
PROPOSAL_ID="$5"
PROVIDER_NODE_ID="$6"
PROVIDER_BINARY="providerd"
CONSUMER_CHAIN_BINARY="wasmd_consumer"
WASM_BINARY="wasmd"
TOOL_OUTPUT="$TOOL_OUTPUT_DIR/$(date +"%Y-%m-%d_%H-%M-%S")"

#TODO: check error code
bash verify_proposal.sh $WASM_CONTRACTS $CONSUMER_CHAIN_ID $CONSUMER_CHAIN_MULTISIG_ADDRESS $CONSUMER_CHAIN_BINARY $WASM_BINARY $TOOL_OUTPUT $PROPOSAL_ID $PROVIDER_NODE_ID $PROVIDER_BINARY

./$PROVIDER_BINARY q provider consumer-genesis $CONSUMER_CHAIN_ID --node $PROVIDER_NODE_ID --output json > $TOOL_OUTPUT/consumer_section.json
jq -s '.[0].app_state.ccvconsumer = .[1] | .[0]' $TOOL_OUTPUT/genesis.json $TOOL_OUTPUT/consumer_section.json > $TOOL_OUTPUT/genesis_consumer.json && \
	mv $TOOL_OUTPUT/genesis_consumer.json $TOOL_OUTPUT/genesis.json

rm -f $TOOL_OUTPUT/consumer_section.json
rm -f $TOOL_OUTPUT/sha256hashes.json

