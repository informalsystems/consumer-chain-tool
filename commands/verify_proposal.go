package commands

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const verifyProposalScript = `#!/bin/bash
set -eu
LOCATION_OF_SMART_CONTRACTS_BINARIES="%s"
CHAIN_ID="%s"
MULTISIG_ADDRESS="%s"
CONSUMER_BINARY="%s"
COSMWASM_BINARY="%s"
TOOL_OUTPUT_LOCATION="%s"
CREATE_OUTPUT_SUBFOLDER="%s"
GENESIS_HASH="%s"
BINARY_HASH="%s"
SPAWN_TIME="%s"

docker run --rm \
-v "$LOCATION_OF_SMART_CONTRACTS_BINARIES":/contract_binaries \
-v "$TOOL_OUTPUT_LOCATION":/tool_output \
consumer-chain-tool:latest sh ./verify_proposal.sh "/contract_binaries" "$CHAIN_ID" "$MULTISIG_ADDRESS" "$CONSUMER_BINARY" "$COSMWASM_BINARY" "/tool_output" "$CREATE_OUTPUT_SUBFOLDER" "$GENESIS_HASH" "$BINARY_HASH" "$SPAWN_TIME"
`

func NewVerifyProposalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     getVerifyCommandUsage(),
		Example: getVerifyCommandExample(),
		Short:   VerifyProposalShortDesc,
		Long:    getVerifyProposalLongDesc(),
		Args:    cobra.ExactArgs(VerifyProposalCmdParamsCount),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputs, err := NewVerifyProposalArgs(args)
			if err != nil {
				return err
			}

			bashCmd := exec.Command("/bin/bash", "-c", fmt.Sprintf(verifyProposalScript, inputs.contractBinariesLocation, inputs.consumerChainId,
				inputs.multisigAddress, ConsumerBinary, CosmWasmBinary, inputs.toolOutputLocation, "true", inputs.proposalGenesisHash,
				inputs.proposalBinaryHash, inputs.proposalSpawnTime))

			RunCmdAndPrintOutput(bashCmd)

			return nil
		},
	}

	return cmd
}

func getVerifyCommandUsage() string {
	return fmt.Sprintf("%s [%s] [%s] [%s] [%s] [%s] [%s] [%s]",
		VerifyProposalCmdName, ContractBinariesLocation, ConsumerChainId,
		MultisigAddress, ToolOutputLocation, ProposalGenesisHash, ProposalBinaryHash, ProposalSpawnTime)
}

func getVerifyCommandExample() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s %s %s",
		ToolName, VerifyProposalCmdName, "$HOME/contract_binaries", "wasm", "wasm1ykqt29d4ekemh5pc0d2wdayxye8yqupttf6vyz",
		"$HOME/cli_tool_output", "519df96a862c30f53e67b1277e6834ab4bd59dfdd08c781d1b7cf3813080fb28", "09184916f3e85aa6fa24d3c12f1e5465af2214f13db265a52fa9f4617146dea5", "2022-06-01T09:10:00Z")
}

func getVerifyProposalLongDesc() string {
	return fmt.Sprintf(VerifyProposalLongDesc, ContractBinariesLocation, ConsumerChainId,
		MultisigAddress, ToolOutputLocation, ProposalGenesisHash, ProposalBinaryHash, ProposalSpawnTime)
}

type VerifyProposalArgs struct {
	contractBinariesLocation string
	consumerChainId          string
	multisigAddress          string
	toolOutputLocation       string
	proposalGenesisHash      string
	proposalBinaryHash       string
	proposalSpawnTime        string
}

func NewVerifyProposalArgs(args []string) (*VerifyProposalArgs, error) {
	if len(args) != VerifyProposalCmdParamsCount {
		return nil, fmt.Errorf("unexpected number of arguments. Expected: %d, received: %d", VerifyProposalCmdParamsCount, len(args))
	}

	commandArgs := new(VerifyProposalArgs)
	var errors []string

	contractBinariesLocation := strings.TrimSpace(args[0])
	if IsValidInputPath(contractBinariesLocation) {
		commandArgs.contractBinariesLocation = contractBinariesLocation
	} else {
		errors = append(errors, fmt.Sprintf("Provided input path '%s' is not a valid directory.", contractBinariesLocation))
	}

	consumerChainId := strings.TrimSpace(args[1])
	if IsValidString(consumerChainId) {
		commandArgs.consumerChainId = consumerChainId
	} else {
		errors = append(errors, fmt.Sprintf("Provided chain-id '%s' is not valid.", consumerChainId))
	}

	multisigAddress := strings.TrimSpace(args[2])
	if IsValidString(multisigAddress) {
		commandArgs.multisigAddress = multisigAddress
	} else {
		errors = append(errors, fmt.Sprintf("Provided multisig address '%s' is not valid.", multisigAddress))
	}

	toolOutputLocation := strings.TrimSpace(args[3])
	if IsValidOutputPath(toolOutputLocation) {
		commandArgs.toolOutputLocation = toolOutputLocation
	} else {
		errors = append(errors, fmt.Sprintf("Provided output path '%s' is not a valid directory.", toolOutputLocation))
	}

	proposalGenesisHash := strings.TrimSpace(args[4])
	if IsValidString(proposalGenesisHash) {
		commandArgs.proposalGenesisHash = proposalGenesisHash
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal genesis hash '%s' is not valid.", proposalGenesisHash))
	}

	proposalBinaryHash := strings.TrimSpace(args[5])
	if IsValidString(proposalBinaryHash) {
		commandArgs.proposalBinaryHash = proposalBinaryHash
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal binary hash '%s' is not valid.", proposalBinaryHash))
	}

	proposalSpawnTime := strings.TrimSpace(args[6])
	if spawnTime, isValid := IsValidDateTime(proposalSpawnTime); isValid {
		commandArgs.proposalSpawnTime = spawnTime.Format(time.RFC3339Nano)
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal spawn time '%s' is not valid.", proposalSpawnTime))
	}
	if len(errors) > 0 {
		return nil, fmt.Errorf(strings.Join(errors, "\n"))
	}

	return commandArgs, nil
}

const (
	VerifyProposalShortDesc = "Verify that genesis and binary hash inputs match the hashes of the consumer binary and the regenerated genesis file."
	VerifyProposalLongDesc  = `This command takes the same inputs and goes through the same process as 'prepare-proposal' command to create the genesis.json file and calculate its hash.
  If the input hashes from the command match the recalculated ones, then the resulting genesis.json file contains the smart contracts provided to the input of this command.

Command arguments:
    %s - The location of the directory that contains the compiled smart contracts .wasm binaries.
    %s - The proposed chain-id of the new consumer chain must be different from all the other consumer chain ids of the executing provider chain.
    %s - The multi-signature address that will have the permission to instantiate the contracts from the set of pre-deployed codes.
    %s - The location of the directory where the verification data will be saved.
    %s - The hash of the genesis file can be obtained by querying the proposal previously submitted to the provider chain.
    %s - The hash of the consumer binary can be obtained by querying the proposal previously submitted to the provider chain.
    %s - The time on the provider chain at which the consumer chain genesis is finalized and all the validators will be responsible for starting their consumer chain validator node. `
)
