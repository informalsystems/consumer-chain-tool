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
		Short:   "Create genesis.json and proposal.json for the given set of smart contracts",
		Long:    `TODO: Add a longer description`,
		Args:    cobra.ExactArgs(PrepareProposalCmdParamsCount),
		RunE: func(cmd *cobra.Command, args []string) error {
			inputs, err := NewPrepareProposalArgs(args)
			if err != nil {
				return err
			}

			bashCmd := exec.Command("/bin/bash", "prepare_proposal.sh",
				inputs.smartContractsLocation, inputs.consumerChainId, inputs.multisigAddress,
				inputs.toolOutputLocation, inputs.proposalTitle, inputs.proposalDescription,
				inputs.proposalRevisionHeight, inputs.proposalSpawnTime, inputs.proposalDeposit)

			RunCmdAndPrintOutput(bashCmd)

			return nil
		},
	}

	return prepareProposalCmd
}

func getPrepareCommandUsage() string {
	return fmt.Sprintf("%s [%s] [%s] [%s] [%s] [%s] [%s] [%s] [%s] [%s]",
		PrepareProposalCmdName, SmartContractsLocation, ConsumerChainId, MultisigAddress, ToolOutputLocation,
		ProposalTitle, ProposalDescription, ProposalRevisionHeight, ProposalSpawnTime, ProposalDeposit)
}

func getPrepareCommandExample() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s %s",
		ToolName, PrepareProposalCmdName, "$HOME/wasm_contracts", "wasm", "wasm1243cuuy98lxaf7ufgav0w76xt5es93afr8a3ya", "$HOME/tool_output_step1",
		"\"Create a chain\"", "\"Gonna be a great chain\"", "1", "2022-06-01T09:10:00.000000000-00:00", "10000001stake")
}

// TODO: proposalRevisionHeight and proposalSpawnTime are only passed to shell script so there is no need to make them int and time.Time?
// TODO: expand with RevisionNumber, proposal also!
type PrepareProposalArgs struct {
	smartContractsLocation string
	consumerChainId        string
	multisigAddress        string
	toolOutputLocation     string
	proposalTitle          string
	proposalDescription    string
	proposalRevisionHeight string
	proposalSpawnTime      string
	proposalDeposit        string
}

func NewPrepareProposalArgs(args []string) (*PrepareProposalArgs, error) {
	if len(args) != PrepareProposalCmdParamsCount {
		return nil, fmt.Errorf("Unexpected number of arguments. Expected: %d, received: %d.", PrepareProposalCmdParamsCount, len(args))
	}

	commandArgs := new(PrepareProposalArgs)
	var errors []string

	smartContractsLocation := strings.TrimSpace(args[0])
	if IsValidInputPath(smartContractsLocation) {
		commandArgs.smartContractsLocation = smartContractsLocation
	} else {
		errors = append(errors, fmt.Sprintf("Provided input path '%s' is not a valid directory.", smartContractsLocation))
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
	if IsValidProposalRevisionHeight(proposalRevisionHeight) {
		commandArgs.proposalRevisionHeight = proposalRevisionHeight
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal revision height '%s' is not valid.", proposalRevisionHeight))
	}

	proposalSpawnTime := strings.TrimSpace(args[7])
	if spawnTime, isValid := IsValidDateTime(proposalSpawnTime); isValid {
		commandArgs.proposalSpawnTime = spawnTime.Format(time.RFC3339Nano)
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal spawn time '%s' is not valid.", proposalSpawnTime))
	}

	proposalDeposit := strings.TrimSpace(args[8])
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
