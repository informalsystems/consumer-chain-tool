package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

const finalizeGenesisScript = `#!/bin/bash
set -eu
LOCATION_OF_SMART_CONTRACTS_BINARIES="%s"
CHAIN_ID="%s"
MULTISIG_ADDRESS="%s"
TOOL_OUTPUT_LOCATION="%s"
PROPOSAL_ID="%s"
PROVIDER_NODE_ADDRESS="%s"
PROVIDER_BINARY_PATH="%s"

docker run --rm --network="host" \
    -v "$LOCATION_OF_SMART_CONTRACTS_BINARIES":/contract_binaries \
    -v "$PROVIDER_BINARY_PATH":/go/bin/interchain-security-pd \
    -v "$TOOL_OUTPUT_LOCATION":/tool_output \
    consumer-chain-tool:latest sh ./finalize_genesis.sh "/contract_binaries" "$CHAIN_ID" "$MULTISIG_ADDRESS" "/tool_output" "$PROPOSAL_ID" "$PROVIDER_NODE_ADDRESS" "/go/bin/interchain-security-pd"
`

func NewFinalizeGenesisCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     getFinalizeCommandUsage(),
		Example: getFinalizeCommandExample(),
		Short:   FinalizeGenesisShortDesc,
		Long:    getFinalizeGenesisLongDesc(),
		Args:    cobra.ExactArgs(FinalizeGenesisCmdParamsCount),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputs, err := NewFinalizeGenesisArgs(args)
			if err != nil {
				return err
			}

			bashCmd := exec.Command("/bin/bash", "-c", fmt.Sprintf(finalizeGenesisScript, inputs.contractBinariesLocation, inputs.consumerChainId,
				inputs.multisigAddress, inputs.toolOutputLocation, inputs.proposalId, inputs.providerNodeAddress, inputs.providerBinaryPath))

			RunCmdAndPrintOutput(bashCmd)

			return nil
		},
	}

	return cmd
}

func getFinalizeCommandUsage() string {
	return fmt.Sprintf("%s [%s] [%s] [%s] [%s] [%s] [%s] [%s]",
		FinalizeGenesisCmdName, ContractBinariesLocation, ConsumerChainId,
		MultisigAddress, ToolOutputLocation, ProposalId, ProviderNodeAddress, ProviderBinaryPath)
}

func getFinalizeCommandExample() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s %s %s",
		ToolName, FinalizeGenesisCmdName, "$HOME/contract_binaries", "wasm", "wasm1ykqt29d4ekemh5pc0d2wdayxye8yqupttf6vyz",
		"$HOME/cli_tool_output", "1", "tcp://localhost:26657", "$HOME/gaiad")
}

func getFinalizeGenesisLongDesc() string {
	return fmt.Sprintf(FinalizeGenesisLongDesc, ContractBinariesLocation, ConsumerChainId,
		MultisigAddress, ToolOutputLocation, ProposalId, ProviderNodeAddress, ProviderBinaryPath)
}

type FinalizeGenesisArgs struct {
	contractBinariesLocation string
	consumerChainId          string
	multisigAddress          string
	toolOutputLocation       string
	proposalId               string
	providerNodeAddress      string
	providerBinaryPath       string
}

func NewFinalizeGenesisArgs(args []string) (*FinalizeGenesisArgs, error) {
	if len(args) != FinalizeGenesisCmdParamsCount {
		return nil, fmt.Errorf("unexpected number of arguments. Expected: %d, received: %d", FinalizeGenesisCmdParamsCount, len(args))
	}

	commandArgs := new(FinalizeGenesisArgs)
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

	proposalId := strings.TrimSpace(args[4])
	if isPositiveInt(proposalId) {
		commandArgs.proposalId = proposalId
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal id '%s' is not valid.", proposalId))
	}

	providerNodeAddress := strings.TrimSpace(args[5])
	if IsValidString(providerNodeAddress) {
		commandArgs.providerNodeAddress = providerNodeAddress
	} else {
		errors = append(errors, fmt.Sprintf("Provided provider node address '%s' is not valid.", providerNodeAddress))
	}

	providerBinaryPath := strings.TrimSpace(args[6])
	if IsValidFilePath(providerBinaryPath) {
		commandArgs.providerBinaryPath = providerBinaryPath
	} else {
		errors = append(errors, fmt.Sprintf("Provided provider binary path '%s' is not valid.", providerBinaryPath))
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf(strings.Join(errors, "\n"))
	}

	return commandArgs, nil
}

const (
	FinalizeGenesisShortDesc = "Build the final genesis.json for Interchain Security consumer chain with CosmWasm smart contracts deployed"
	FinalizeGenesisLongDesc  = `This command takes the same inputs and goes through the same process as 'verify-proposal' command to verify the command inputs against the provided proposal.
It then queries the provider chain to obtain the consumer section for the chain ID and appends this data to the initial genesis.json, which results in Interchain Secuirty consumer-enabled genesis with CosmWasm smart contracts deployed.

Command arguments:
    %s - The location of the directory that contains the compiled smart contracts .wasm binaries.
    %s - The proposed chain-id of the new consumer chain must be different from all the other consumer chain ids of the executing provider chain.
    %s - The multi-signature address that will have the permission to instantiate the contracts from the set of pre-deployed codes.
    %s - The location of the directory where the final genesis.json and consumer binary will be saved. The validators will use the outputs to start the consumer chain.
    %s - The ID of the proposal submitted to the provider chain whose data will be used to verify if the inputs of this command match the ones from the proposal.
    %s - This represents the address of the provider chain node in the following format: tcp://IP_ADDRESS:PORT_NUMBER. This address is used to query the provider chain to obtain the consumer section for the genesis file.
    %s - The path to the provider binary.`
)
