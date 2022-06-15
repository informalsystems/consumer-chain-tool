#!/bin/bash
set -eu

#bash prepare_proposal.sh $HOME/wasm_contracts wasm wasm1243cuuy98lxaf7ufgav0w76xt5es93afr8a3ya $HOME/tool_output_step1 "Create a chain" "Gonna be a great chain" 1 2022-06-01T09:10:00.000000000-00:00 10000001stake

TOOL_INPUT="$1"
CONSUMER_CHAIN_ID="$2"
CONSUMER_CHAIN_MULTISIG_ADDRESS="$3"
TOOL_OUTPUT_DIR="$4"
PROPOSAL_TITLE="$5"
PROPOSAL_DESCRIPTION="$6" #TODO: add link with output to description
PROPOSAL_REVISION_HEIGHT="$7"
PROPOSAL_SPAWN_TIME="$8"
PROPOSAL_DEPOSIT="$9"
CONSUMER_CHAIN_BINARY="wasmd_consumer"
WASM_BINARY="wasmd"
WASM_CONTRACTS_SOURCES="$TOOL_INPUT/contracts_source_code" #TODO: temporary separated, in the end we will have only one folder that contains just the source code
TOOL_OUTPUT="$TOOL_OUTPUT_DIR/$(date +"%Y-%m-%d_%H-%M-%S")"

# Delete all generated data.
clean_up () {
    rm -f $TOOL_OUTPUT/sha256hashes.json
} 
trap clean_up EXIT

if ! bash prepare_proposal_inputs.sh $TOOL_INPUT $CONSUMER_CHAIN_ID $CONSUMER_CHAIN_MULTISIG_ADDRESS $CONSUMER_CHAIN_BINARY $WASM_BINARY $TOOL_OUTPUT $PROPOSAL_SPAWN_TIME;then
    echo "Error while preparing proposal data!"
    exit 1
fi

#################################### COPY SMART CONTRACTS #############################

#TODO filter only *.rs files, they might be in different folders/subfolders, some malicious files can be added, etc.
# not sure if it really helps, since they can add malicious things afterwards
cp -r $WASM_CONTRACTS_SOURCES $TOOL_OUTPUT/contracts

#################################### CREATE PROPOSAL JSON #############################

#TODO: check all params
tee $TOOL_OUTPUT/proposal.json<<EOF
{
    "title": "$PROPOSAL_TITLE",
    "description": "$PROPOSAL_DESCRIPTION",
    "chain_id": "$CONSUMER_CHAIN_ID",
    "initial_height": {
        "revision_height": $PROPOSAL_REVISION_HEIGHT 
    },
    "genesis_hash": "$(jq -r ".genesis_hash" $TOOL_OUTPUT/sha256hashes.json)",
    "binary_hash": "$(jq -r ".binary_hash" $TOOL_OUTPUT/sha256hashes.json)",
    "spawn_time": "$PROPOSAL_SPAWN_TIME",
    "deposit": "$PROPOSAL_DEPOSIT"
}
EOF
