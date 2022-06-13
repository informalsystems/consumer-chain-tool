package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const (
	ToolName               = "consumer-chain-tool"
	PrepareProposalCmdName = "prepare-proposal"
	VerifyProposalCmdName  = "verify-proposal"
	FinalizeGenesisCmdName = "finalize-genesis"
	ProviderBinary         = "providerd"
	ConsumerBinary         = "wasmd_consumer"
	CosmWasmBinary         = "wasmd"
	SmartContractsLocation = "smart-contracts-location"
	ConsumerChainId        = "consumer-chain-id"
	MultisigAddress        = "multisig-address"
	ToolOutputLocation     = "tool-output-location"
	ProposalId             = "proposal-id"
	ProposalTitle          = "proposal-title"
	ProposalDescription    = "proposal-description"
	ProposalRevisionHeight = "proposal-revision-height"
	ProposalSpawnTime      = "proposal-spawn-time"
	ProposalDeposit        = "proposal-deposit"
	ProviderNodeId         = "provider-node-id"
	PrepareCmdUsageExample = `consumer-chain-tool prepare-proposal \
	--smart-contracts-location $HOME/wasm_contracts \
	--consumer-chain-id wasm \
	--multisig-address wasm1243cuuy98lxaf7ufgav0w76xt5es93afr8a3ya \
	--tool-output-location $HOME/tool_output_step1 \
	--proposal-title "Create a chain" \
	--proposal-description "Gonna be a great chain" \
	--proposal-revision-height 1 \
	--proposal-spawn-time 2022-03-11T17:02:14.718477Z \
	--proposal-deposit 10000001stake`
	VerifyCmdUsageExample = `consumer-chain-tool verify-proposal \
	--smart-contracts-location $HOME/wasm_contracts \
	--consumer-chain-id wasm \
	--multisig-address wasm1243cuuy98lxaf7ufgav0w76xt5es93afr8a3ya \
	--tool-output-location $HOME/tool_output_step2 \
	--proposal-id 1 \
	--provider-node-id "tcp://localhost:26657"`
	FinalizeGenesisCmdUsageExample = `consumer-chain-tool finalize-genesis \
	--smart-contracts-location $HOME/wasm_contracts \
	--consumer-chain-id wasm \
	--multisig-address wasm1243cuuy98lxaf7ufgav0w76xt5es93afr8a3ya \
	--tool-output-location $HOME/tool_output_step2 \
	--proposal-id 1 \
	--provider-node-id "tcp://localhost:26657"`
)

func init() {
	cobra.EnableCommandSorting = false
}

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   ToolName,
		Short: fmt.Sprintf("%s - prepare and verify proposals for a new Interchain Security enabled CosmWasm consumer chain", ToolName),
		Long:  `TODO: Add a longer description`,
	}

	rootCmd.AddCommand(
		NewPrepareProposalCommand(),
		NewVerifyProposalCommand(),
		NewFinalizeGenesisCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}

//TODO: output only "echo" commands
func RunCmdAndPrintOutput(bashCmd *exec.Cmd) {
	cmdReader, err := bashCmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	bashCmd.Stderr = bashCmd.Stdout

	if err := bashCmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(cmdReader)

	for scanner.Scan() {
		out := scanner.Text()
		fmt.Println(out)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
