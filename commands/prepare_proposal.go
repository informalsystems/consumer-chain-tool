package commands

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

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

			bashCmd := exec.Command("/bin/bash", "-c", prepareProposalScript, prepareProposalInputsScript,
				inputs.contractBinariesLocation, inputs.consumerChainId, inputs.multisigAddress,
				inputs.toolOutputLocation, inputs.proposalTitle, inputs.proposalDescription,
				inputs.proposalRevisionHeight, inputs.proposalRevisionNumber, inputs.proposalSpawnTime, inputs.proposalDeposit)

			RunCmdAndPrintOutput(bashCmd)

			return nil
		},
	}

	return prepareProposalCmd
}

func getPrepareCommandUsage() string {
	return fmt.Sprintf("%s [%s] [%s] [%s] [%s] [%s] [%s] [%s] [%s]",
		PrepareProposalCmdName, ConsumerChainId, MultisigAddress, ProposalTitle, ProposalDescription,
		ProposalRevisionHeight, ProposalRevisionNumber, ProposalSpawnTime, ProposalDeposit)
}

func getPrepareCommandExample() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s",
		ToolName, PrepareProposalCmdName, "wasm", "wasm1ykqt29d4ekemh5pc0d2wdayxye8yqupttf6vyz", "\"Create a chain\"",
		"\"Gonna be a great chain\"", "4", "0", "2022-06-01T09:10:00.000000000-00:00", "10000001stake")
}

func getPrepareProposalLongDesc() string {
	return fmt.Sprintf(PrepareProposalLongDesc, ConsumerChainId, MultisigAddress, ProposalTitle,
		ProposalDescription, ProposalRevisionHeight, ProposalRevisionNumber, ProposalSpawnTime, ProposalDeposit)
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

	commandArgs.contractBinariesLocation = ContractBinariesLocation

	consumerChainId := strings.TrimSpace(args[0])
	if IsValidString(consumerChainId) {
		commandArgs.consumerChainId = consumerChainId
	} else {
		errors = append(errors, fmt.Sprintf("Provided chain-id '%s' is not valid.", consumerChainId))
	}

	multisigAddress := strings.TrimSpace(args[1])
	if IsValidString(multisigAddress) {
		commandArgs.multisigAddress = multisigAddress
	} else {
		errors = append(errors, fmt.Sprintf("Provided multisig address '%s' is not valid.", multisigAddress))
	}

	commandArgs.toolOutputLocation = ToolOutputLocation

	proposalTitle := strings.TrimSpace(args[2])
	if IsValidString(proposalTitle) {
		commandArgs.proposalTitle = proposalTitle
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal title '%s' is not valid.", proposalTitle))
	}

	proposalDescription := strings.TrimSpace(args[3])
	if IsValidString(proposalDescription) {
		commandArgs.proposalDescription = proposalDescription
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal description '%s' is not valid.", proposalDescription))
	}

	proposalRevisionHeight := strings.TrimSpace(args[4])
	if isPositiveInt(proposalRevisionHeight) {
		commandArgs.proposalRevisionHeight = proposalRevisionHeight
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal revision height '%s' is not valid.", proposalRevisionHeight))
	}

	proposalRevisionNumber := strings.TrimSpace(args[5])
	if isPositiveInt(proposalRevisionNumber) {
		commandArgs.proposalRevisionNumber = proposalRevisionNumber
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal revision number '%s' is not valid.", proposalRevisionNumber))
	}

	proposalSpawnTime := strings.TrimSpace(args[6])
	if spawnTime, isValid := IsValidDateTime(proposalSpawnTime); isValid {
		commandArgs.proposalSpawnTime = spawnTime.Format(time.RFC3339Nano)
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal spawn time '%s' is not valid.", proposalSpawnTime))
	}

	proposalDeposit := strings.TrimSpace(args[7])
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
Then it calculates the SHA256 hashes of this genesis.json file and the binary that should be used to start a new blockchain with this genesis.json file.
It then uses those SHA256 hashes, together with some other input arguments, to create the proposal.json file that can be submitted as a proposal to the Interchain Security enabled provider blockchain.
	
Command arguments:
    %s - The desired chain ID of the consumer chain.
    %s - The multi-signature address that will have the permission to instantiate contracts from the set of pre-deployed codes.
    %s - Proposal title.
    %s - Proposal description. It should contain the publicly available link where the results of this command will be placed.
    %s - The proposal revision height
    %s - The proposal revision number
    %s - The desired time of consumer chain start in the yyyy-MM-ddTHH:mm:ss.fffffffff-zz:zz format (e.g. 2022-06-01T09:10:00.000000000-00:00).
    %s - The amount of tokens for the initial proposal deposit.`
)
