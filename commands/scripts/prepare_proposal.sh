#!/bin/bash
set -eu

TOOL_INPUT="$1"
CONSUMER_CHAIN_ID="$2"
CONSUMER_CHAIN_MULTISIG_ADDRESS="$3"
TOOL_OUTPUT_DIR="$4"
PROPOSAL_TITLE="$5"
PROPOSAL_DESCRIPTION="$6" 
PROPOSAL_REVISION_HEIGHT="$7"
PROPOSAL_REVISION_NUMBER="$8"
PROPOSAL_SPAWN_TIME="$9"
PROPOSAL_DEPOSIT="${10}"
CONSUMER_CHAIN_BINARY="wasmd_consumer"
WASM_BINARY="wasmd"
TOOL_OUTPUT="$TOOL_OUTPUT_DIR"/$(date +"%Y-%m-%d_%H-%M-%S")
LOG="$TOOL_OUTPUT"/log_file.txt

# Delete all generated data.
clean_up () {
    rm -f "$TOOL_OUTPUT"/sha256hashes.json
} 
trap clean_up EXIT

echo "Generating files and hashes..."
if ! bash prepare_proposal_inputs.sh "$TOOL_INPUT" "$CONSUMER_CHAIN_ID" "$CONSUMER_CHAIN_MULTISIG_ADDRESS" "$CONSUMER_CHAIN_BINARY" "$WASM_BINARY" "$TOOL_OUTPUT" "$PROPOSAL_SPAWN_TIME";
then
    echo "Error while preparing proposal data! For more details please check the log file in output directory."
    exit 1
fi

#################################### CREATE PROPOSAL JSON #############################

echo "Generating proposal.json..."
tee "$TOOL_OUTPUT"/proposal.json  &> /dev/null <<EOF
{
    "title": "$PROPOSAL_TITLE",
    "description": "$PROPOSAL_DESCRIPTION",
    "chain_id": "$CONSUMER_CHAIN_ID",
    "initial_height": {
        "revision_number": $PROPOSAL_REVISION_NUMBER,
        "revision_height": $PROPOSAL_REVISION_HEIGHT
    },
    "genesis_hash": "$(jq -r ".genesis_hash" "$TOOL_OUTPUT"/sha256hashes.json)",
    "binary_hash": "$(jq -r ".binary_hash" "$TOOL_OUTPUT"/sha256hashes.json)",
    "spawn_time": "$PROPOSAL_SPAWN_TIME",
    "deposit": "$PROPOSAL_DEPOSIT"
}
EOF
echo "Output data is saved at the specified location"