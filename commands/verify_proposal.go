package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func NewVerifyProposalCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     getVerifyCommandUsage(),
		Example: getVerifyCommandExample(),
		Short:   "Verify that genesis and binary hashes created from the provided inputs match the hashes from the 'create consumer chain' proposal",
		Long:    `TODO: Add a longer description`,
		Args:    cobra.ExactArgs(VerifyProposalCmdParamsCount),
		RunE: func(cmd *cobra.Command, args []string) error {
			inpus, err := NewVerifyProposalArgs(args)
			if err != nil {
				return err
			}

			bashCmd := exec.Command("/bin/bash", "verify_proposal.sh",
				inpus.smartContractsLocation, inpus.consumerChainId, inpus.multisigAddress,
				ConsumerBinary, CosmWasmBinary, inpus.toolOutputLocation, "true", // true for create output subdirectory
				inpus.proposalId, inpus.providerNodeId, ProviderBinary)

			RunCmdAndPrintOutput(bashCmd)

			return nil
		},
	}

	return cmd
}

func getVerifyCommandUsage() string {
	return fmt.Sprintf("%s [%s] [%s] [%s] [%s] [%s] [%s]",
		VerifyProposalCmdName, SmartContractsLocation, ConsumerChainId,
		MultisigAddress, ToolOutputLocation, ProposalId, ProviderNodeId)
}

func getVerifyCommandExample() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s %s",
		ToolName, VerifyProposalCmdName, "$HOME/wasm_contracts", "wasm", "wasm1243cuuy98lxaf7ufgav0w76xt5es93afr8a3ya",
		"$HOME/tool_output_step2", "1", "tcp://localhost:26657")
}

type VerifyProposalArgs struct {
	smartContractsLocation string
	consumerChainId        string
	multisigAddress        string
	toolOutputLocation     string
	proposalId             string
	providerNodeId         string
}

func NewVerifyProposalArgs(args []string) (*VerifyProposalArgs, error) {
	if len(args) != VerifyProposalCmdParamsCount {
		return nil, fmt.Errorf("Unexpected number of arguments. Expected: %d, received: %d.", VerifyProposalCmdParamsCount, len(args))
	}

	commandArgs := new(VerifyProposalArgs)
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

	proposalId := strings.TrimSpace(args[4])
	if IsValidProposalId(proposalId) {
		commandArgs.proposalId = proposalId
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal id '%s' is not valid.", proposalId))
	}

	// TODO: not sure if we should validate node id with regex
	providerNodeId := strings.TrimSpace(args[5])
	if IsValidString(providerNodeId) {
		commandArgs.providerNodeId = providerNodeId
	} else {
		errors = append(errors, fmt.Sprintf("Provided provider node id '%s' is not valid.", providerNodeId))
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf(strings.Join(errors, "\n"))
	}

	return commandArgs, nil
}
