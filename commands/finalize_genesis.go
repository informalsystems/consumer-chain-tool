package commands

import (
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

func NewFinalizeGenesisCommand() *cobra.Command {
	var smartContractsLocation string
	var consumerChainId string
	var multisigAddress string
	var toolOutputLocation string
	var proposalId int
	var providerNodeId string

	cmd := &cobra.Command{
		Use:     FinalizeGenesisCmdName,
		Example: FinalizeGenesisCmdUsageExample,
		Short:   "TODO: Finalize genesis description",
		Long:    `TODO: Add a longer description`,
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: validate all input values

			bashCmd := exec.Command("/bin/bash", "finalize_genesis.sh",
				smartContractsLocation, consumerChainId, multisigAddress,
				toolOutputLocation, strconv.Itoa(proposalId), providerNodeId)

			RunCmdAndPrintOutput(bashCmd)
		},
	}

	cmdFlags := cmd.Flags()
	cmdFlags.StringVar(&smartContractsLocation, SmartContractsLocation, "", "Path to the smart contracts source code folder")
	cmdFlags.StringVar(&consumerChainId, ConsumerChainId, "", "TODO")
	cmdFlags.StringVar(&multisigAddress, MultisigAddress, "", "TODO")
	cmdFlags.StringVar(&toolOutputLocation, ToolOutputLocation, "", "TODO")
	cmdFlags.IntVar(&proposalId, ProposalId, 1, "TODO")
	cmdFlags.StringVar(&providerNodeId, ProviderNodeId, "", "TODO")

	return cmd
}
