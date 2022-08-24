package commands

import (
	"bufio"
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

const (
	ToolName                      = "consumer-chain-tool"
	PrepareProposalCmdName        = "prepare-proposal"
	VerifyProposalCmdName         = "verify-proposal"
	FinalizeGenesisCmdName        = "finalize-genesis"
	PrepareProposalCmdParamsCount = 10
	VerifyProposalCmdParamsCount  = 7
	FinalizeGenesisCmdParamsCount = 7
	ConsumerBinary                = "wasmd_consumer"
	CosmWasmBinary                = "wasmd"
	ContractBinariesLocation      = "contract-binaries-location"
	ConsumerChainId               = "consumer-chain-id"
	MultisigAddress               = "multisig-address"
	ProposalId                    = "proposal-id"
	ProposalTitle                 = "proposal-title"
	ProposalDescription           = "proposal-description"
	ProposalRevisionHeight        = "proposal-revision-height"
	ProposalRevisionNumber        = "proposal-revision-number"
	ProposalSpawnTime             = "proposal-spawn-time"
	ProposalDeposit               = "proposal-deposit"
	ProposalGenesisHash           = "proposal-genesis-hash"
	ProposalBinaryHash            = "proposal-binary-hash"
	ProviderNodeAddress           = "provider-node-address"
	ToolOutputLocation            = "tool-output-location"
	ProviderBinaryPath            = "provider-binary-path"
)

func init() {
	cobra.EnableCommandSorting = false
}

func Execute() {
	var rootCmd = &cobra.Command{
		Use:   ToolName,
		Short: fmt.Sprintf(ToolShortDesc, ToolName),
		Long:  ToolLongDesc,
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(
		NewPrepareProposalCommand(),
		NewVerifyProposalCommand(),
		NewFinalizeGenesisCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "An error occured while executing command: '%s'", err)
		os.Exit(1)
	}
}

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

const (
	ToolShortDesc = "%s - prepare and verify proposal and genesis file for a new Interchain Security enabled CosmWasm consumer chain"
	ToolLongDesc  = `The purpose of the tool is to produce an output in the form of proposal and genesis files. In that way, the process of starting the CosmWasm consumer chain with the pre-deployed smart contract codes is simplified. The process of creating the proposal and genesis data should be done in the following steps:
	1. The proposer runs a prepare-proposal tool command which generates the genesis.json and proposal.json file. All the sections within the genesis file are populated with the final data, except for the ccvconsumer section, which represents the consumer module. The consumer section will be finalized in a later step. The proposal file contains several fields, among which there are the hashes of the genesis file and consumer binary file which will be used to run the consumer chain. The description field of the proposal file should contain a link to the location from where the tool output and the source code of the wasm contract can be downloaded. 
	2. After the first step, when the proposal.json is created, the proposer manually submits the 'create consumer chain' proposal to the provider chain.
	3. After the proposal is submitted, validators and all the interested parties can optionally run a verify-proposal command of the tool to verify the genesis data. This is done mostly to check if the pre-deployed smart contract codes match the source which is uploaded by the proposer. To do so, a user running this command will first download the contract source codes, review them and build them. The built smart contracts will be input for the verify-proposal command. Afterwards, the tool will check if the hash of the regenerated genesis matches the one from the proposal and the user can decide whether to vote for the proposal or not.
	4. Finally, validators run a finalize-genesis command, which will generate the consumer binary and the final genesis file by adding proper data in ccvconsumer section. Validators can then use this genesis and binary to run the consumer chain. This step also requires for smart contracts to be built and given as a command input.`
)
