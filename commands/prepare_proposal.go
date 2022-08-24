package commands

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const prepareProposalScript = `#!/bin/bash
set -eu
LOCATION_OF_SMART_CONTRACTS_BINARIES="%s"
CHAIN_ID="%s"
MULTISIG_ADDRESS="%s"
TOOL_OUTPUT_LOCATION="%s"
PROPOSAL_TITLE="%s"
PROPOSAL_DESCRIPTION="%s"
REVISION_HEIGHT="%s"
REVISION_NUMBER="%s"
SPAWN_TIME="%s"
DEPOSIT="%s"

docker run --rm \
-v "$LOCATION_OF_SMART_CONTRACTS_BINARIES":/contract_binaries \
-v "$TOOL_OUTPUT_LOCATION":/tool_output \
dusanmaksimovic/consumer-chain-tool:latest sh ./prepare_proposal.sh "/contract_binaries" "$CHAIN_ID" "$MULTISIG_ADDRESS" "/tool_output" "$PROPOSAL_TITLE" "$PROPOSAL_DESCRIPTION" "$REVISION_HEIGHT" "$REVISION_NUMBER" "$SPAWN_TIME" "$DEPOSIT"
`

//TODO: change the image name once it gets published

func NewPrepareProposalCommand() *cobra.Command {
	prepareProposalCmd := &cobra.Command{
		Use:     getPrepareCommandUsage(),
		Example: getPrepareCommandExample(),
		Short:   PrepareProposalShortDesc,
		Long:    getPrepareProposalLongDesc(),
		Args:    cobra.ExactArgs(PrepareProposalCmdParamsCount),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputs, err := NewPrepareProposalArgs(args)
			if err != nil {
				return err
			}

			bashCmd := exec.Command("/bin/bash", "-c", fmt.Sprintf(prepareProposalScript, inputs.contractBinariesLocation, inputs.consumerChainId, inputs.multisigAddress, inputs.toolOutputLocation,
				inputs.proposalTitle, inputs.proposalDescription, inputs.proposalRevisionHeight, inputs.proposalRevisionNumber, inputs.proposalSpawnTime, inputs.proposalDeposit))

			RunCmdAndPrintOutput(bashCmd)

			return nil
		},
	}

	return prepareProposalCmd
}

func getPrepareCommandUsage() string {
	return fmt.Sprintf("%s [%s] [%s] [%s] [%s] [%s] [%s] [%s] [%s] [%s] [%s]",
		PrepareProposalCmdName, ContractBinariesLocation, ConsumerChainId, MultisigAddress, ToolOutputLocation,
		ProposalTitle, ProposalDescription, ProposalRevisionHeight, ProposalRevisionNumber, ProposalSpawnTime, ProposalDeposit)
}

func getPrepareCommandExample() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s %s %s",
		ToolName, PrepareProposalCmdName, "$HOME/contract_binaries", "wasm", "wasm1ykqt29d4ekemh5pc0d2wdayxye8yqupttf6vyz", "$HOME/cli_tool_output",
		"\"CosmWasm consumer\"", "\"Contracts code location: https://mysharedlocation/proposal_data\"", "4", "0", "2022-06-01T09:10:00Z", "10000001stake")
}

func getPrepareProposalLongDesc() string {
	return fmt.Sprintf(PrepareProposalLongDesc, ContractBinariesLocation, ConsumerChainId, MultisigAddress, ToolOutputLocation,
		ProposalTitle, ProposalDescription, ProposalRevisionHeight, ProposalRevisionNumber, ProposalSpawnTime, ProposalDeposit)
}

type PrepareProposalArgs struct {
	contractBinariesLocation string
	consumerChainId          string
	multisigAddress          string
	toolOutputLocation       string
	proposalTitle            string
	proposalDescription      string
	proposalRevisionHeight   string
	proposalRevisionNumber   string
	proposalSpawnTime        string
	proposalDeposit          string
}

func NewPrepareProposalArgs(args []string) (*PrepareProposalArgs, error) {
	if len(args) != PrepareProposalCmdParamsCount {
		return nil, fmt.Errorf("unexpected number of arguments. Expected: %d, received: %d", PrepareProposalCmdParamsCount, len(args))
	}

	commandArgs := new(PrepareProposalArgs)
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

	proposalTitle := strings.TrimSpace(args[4])
	if IsValidString(proposalTitle) {
		commandArgs.proposalTitle = proposalTitle
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal title '%s' is not valid.", proposalTitle))
	}

	proposalDescription := strings.TrimSpace(args[5])
	if IsValidString(proposalDescription) {
		commandArgs.proposalDescription = proposalDescription
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal description '%s' is not valid.", proposalDescription))
	}

	proposalRevisionHeight := strings.TrimSpace(args[6])
	if isPositiveInt(proposalRevisionHeight) {
		commandArgs.proposalRevisionHeight = proposalRevisionHeight
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal revision height '%s' is not valid.", proposalRevisionHeight))
	}

	proposalRevisionNumber := strings.TrimSpace(args[7])
	if isPositiveInt(proposalRevisionNumber) {
		commandArgs.proposalRevisionNumber = proposalRevisionNumber
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal revision number '%s' is not valid.", proposalRevisionNumber))
	}

	proposalSpawnTime := strings.TrimSpace(args[8])
	if spawnTime, isValid := IsValidDateTime(proposalSpawnTime); isValid {
		commandArgs.proposalSpawnTime = spawnTime.Format(time.RFC3339Nano)
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal spawn time '%s' is not valid.", proposalSpawnTime))
	}

	proposalDeposit := strings.TrimSpace(args[9])
	if IsValidDeposit(proposalDeposit) {
		commandArgs.proposalDeposit = proposalDeposit
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal deposit '%s' is not valid.", proposalDeposit))
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf(strings.Join(errors, "\n"))
	}

	return commandArgs, nil
}

const (
	PrepareProposalShortDesc = "Create genesis.json and proposal.json for the given set of CosmWasm smart contracts"
	PrepareProposalLongDesc  = `This command uses the provided set of compiled CosmWasm smart contracts, together with some other input arguments, to create genesis.json that will have those smart contracts deployed.
Then it calculates the SHA256 hashes of this genesis.json file and the binary that should be used to start a new blockchain.
It then uses those SHA256 hashes, together with some other input arguments, to create the proposal.json file that can be submitted as a proposal to the Interchain Security enabled provider blockchain.
	
Command arguments:
    %s - The location of the directory that contains the compiled smart contracts .wasm binaries.
    %s - The proposed chain-id of the new consumer chain must be different from all the other consumer chain ids of the executing provider chain.
    %s - The multi-signature address that will have the permission to instantiate the contracts from the set of pre-deployed codes.
    %s - The location of the directory where the resulting genesis.json and proposal.json will be saved.
    %s - The title of the proposal.
    %s - The proposal description should contain the publicly available link where the contract's source code and the output of this command are uploaded by the proposer.
    %s - The height within the given revision.
    %s - The revision that the client is currently on.
    %s - The time on the provider chain at which the consumer chain genesis is finalized and all the validators will be responsible for starting heir consumer chain validator node. 
    %s - The amount of tokens for the initial proposal deposit.`
)
