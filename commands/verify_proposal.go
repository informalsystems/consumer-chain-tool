package commands

import (
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

func NewVerifyProposalCommand() *cobra.Command {
	var smartContractsLocation string
	var consumerChainId string
	var multisigAddress string
	var toolOutputLocation string
	var proposalId int
	var providerNodeId string

	cmd := &cobra.Command{
		Use:     VerifyProposalCmdName,
		Example: VerifyCmdUsageExample,
		Short:   "TODO: Verify proposal description",
		Long:    `TODO: Add a longer description`,
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: validate all input values

			bashCmd := exec.Command("/bin/bash", "verify_proposal.sh",
				smartContractsLocation, consumerChainId, multisigAddress,
				ConsumerBinary, CosmWasmBinary, toolOutputLocation,
				strconv.Itoa(proposalId), providerNodeId, ProviderBinary)

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
