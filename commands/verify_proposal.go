package commands

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

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

			bashCmd := exec.Command("/bin/bash", "-c", verifyProposalScript, prepareProposalInputsScript,
				inputs.contractBinariesLocation, inputs.consumerChainId, inputs.multisigAddress,
				ConsumerBinary, CosmWasmBinary, inputs.toolOutputLocation, "true", // true for create output subdirectory
				inputs.proposalGenesisHash, inputs.proposalBinaryHash, inputs.proposalSpawnTime)

			RunCmdAndPrintOutput(bashCmd)

			return nil
		},
	}

	return cmd
}

func getVerifyCommandUsage() string {
	return fmt.Sprintf("%s [%s] [%s] [%s] [%s] [%s]",
		VerifyProposalCmdName, ConsumerChainId, MultisigAddress,
		ProposalGenesisHash, ProposalBinaryHash, ProposalSpawnTime)
}

func getVerifyCommandExample() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s",
		ToolName, VerifyProposalCmdName, "wasm", "wasm1ykqt29d4ekemh5pc0d2wdayxye8yqupttf6vyz",
		"8beb03cf0d59d5c77f0521eaf169311f7ea442ca55894c9c9b8bc58d52806e7a",
		"f3414a11bf4ef5dbd1e65fa341d1ece5d8b7b139f648edd0d2513e4c168a859d", "2022-06-01T09:10:00Z")
}

func getVerifyProposalLongDesc() string {
	return fmt.Sprintf(VerifyProposalLongDesc, ConsumerChainId, MultisigAddress,
		ProposalGenesisHash, ProposalBinaryHash, ProposalSpawnTime)
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

	proposalGenesisHash := strings.TrimSpace(args[2])
	if IsValidString(proposalGenesisHash) {
		commandArgs.proposalGenesisHash = proposalGenesisHash
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal genesis hash '%s' is not valid.", proposalGenesisHash))
	}

	proposalBinaryHash := strings.TrimSpace(args[3])
	if IsValidString(proposalBinaryHash) {
		commandArgs.proposalBinaryHash = proposalBinaryHash
	} else {
		errors = append(errors, fmt.Sprintf("Provided proposal binary hash '%s' is not valid.", proposalBinaryHash))
	}

	proposalSpawnTime := strings.TrimSpace(args[4])
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
	VerifyProposalShortDesc = "Verify that genesis and binary hashes created from the provided inputs match the hashes from the 'create consumer chain' proposal with the given proposal ID"
	VerifyProposalLongDesc  = `This command takes the same inputs and goes through the same process as 'prepare-proposal' command to create the genesis.json file and calculate its hash.
It then queries the 'create consumer chain' proposal from the provider chain to obtain the hashes. If the hashes from the proposal match the recalculated ones, then the resulting genesis.json file contains the smart contracts provided to the input of this command.

Command arguments:
    %s - The chain ID of the consumer chain.
    %s - The multi-signature address that will have the permission to instantiate contracts from the set of predeployed codes.
    %s - The proposal's hash of the genesis file. It can be retrieved by quering a provider chain. 
    %s - The proposal's hash of the consumer binary. It can be retrieved by quering a provider chain.
    %s - The proposal's spawn time. It can be retrieved by quering a provider chain.`
)
