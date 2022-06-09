package commands

import (
	"os/exec"
	"strconv"

	"github.com/spf13/cobra"
)

func NewPrepareProposalCommand() *cobra.Command {
	var smartContractsLocation string
	var consumerChainId string
	var multisigAddress string
	var toolOutputLocation string
	var proposalTitle string
	var proposalDescription string
	var proposalRevisionHeight int
	var proposalSpawnTime string
	var proposalDeposit string

	prepareProposalCmd := &cobra.Command{
		Use:     PrepareProposalCmdName,
		Example: PrepareCmdUsageExample,
		Short:   "Create genesis.json and proposal.json for the given smart contracts",
		Long:    `TODO: Add a longer description`,
		//TODO: it would be better to use arguments but then we need to parse them by ourself
		//Args:  cobra.ExactArgs(9),
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: validate all input values

			bashCmd := exec.Command("/bin/bash", "prepare_proposal.sh",
				smartContractsLocation, consumerChainId, multisigAddress,
				toolOutputLocation, proposalTitle, proposalDescription,
				strconv.Itoa(proposalRevisionHeight), proposalSpawnTime, proposalDeposit)

			RunCmdAndPrintOutput(bashCmd)
		},
	}

	createProposalFlags := prepareProposalCmd.Flags()
	createProposalFlags.StringVar(&smartContractsLocation, SmartContractsLocation, "", "Path to the smart contracts source code folder")
	createProposalFlags.StringVar(&consumerChainId, ConsumerChainId, "", "TODO")
	createProposalFlags.StringVar(&multisigAddress, MultisigAddress, "", "TODO")
	createProposalFlags.StringVar(&toolOutputLocation, ToolOutputLocation, "", "TODO")
	createProposalFlags.StringVar(&proposalTitle, ProposalTitle, "", "TODO")
	createProposalFlags.StringVar(&proposalDescription, ProposalDescription, "", "TODO")
	createProposalFlags.IntVar(&proposalRevisionHeight, ProposalRevisionHeight, 1, "TODO")
	createProposalFlags.StringVar(&proposalSpawnTime, ProposalSpawnTime, "", "TODO")
	createProposalFlags.StringVar(&proposalDeposit, ProposalDeposit, "", "TODO")

	return prepareProposalCmd
}

// func getCommandUsage() string {
// 	return fmt.Sprintf("%s \n\t[%s] \n\t[%s] \n\t[%s] \n\t[%s] \n\t[%s] \n\t[%s] [%s] [%s] [%s]",
// 		PrepareProposalCmdName, SmartContractsLocation, ConsumerChainId, MultisigAddress, ToolOutputLocation,
// 		ProposalTitle, ProposalDescription, ProposalRevisionHeight, ProposalSpawnTime, ProposalDeposit)
// }
