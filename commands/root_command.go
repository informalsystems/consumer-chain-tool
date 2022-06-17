package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

const (
	ToolName                      = "consumer-chain-tool"
	PrepareProposalCmdName        = "prepare-proposal"
	VerifyProposalCmdName         = "verify-proposal"
	FinalizeGenesisCmdName        = "finalize-genesis"
	PrepareProposalCmdParamsCount = 9
	VerifyProposalCmdParamsCount  = 6
	FinalizeGenesisCmdParamsCount = 6
	ProviderBinary                = "providerd"
	ConsumerBinary                = "wasmd_consumer"
	CosmWasmBinary                = "wasmd"
	SmartContractsLocation        = "smart-contracts-location"
	ConsumerChainId               = "consumer-chain-id"
	MultisigAddress               = "multisig-address"
	ToolOutputLocation            = "tool-output-location"
	ProposalId                    = "proposal-id"
	ProposalTitle                 = "proposal-title"
	ProposalDescription           = "proposal-description"
	ProposalRevisionHeight        = "proposal-revision-height"
	ProposalSpawnTime             = "proposal-spawn-time"
	ProposalDeposit               = "proposal-deposit"
	ProviderNodeId                = "provider-node-id"
	VerifyCmdUsageExample         = `consumer-chain-tool verify-proposal \
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

var (
	reDnmString = `[a-zA-Z][a-zA-Z0-9/-]{2,127}`
	reDecAmt    = `[[:digit:]]+(?:\.[[:digit:]]+)?|\.[[:digit:]]+`
	reSpc       = `[[:space:]]*`
	reDecCoin   = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, reDecAmt, reSpc, reDnmString))
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

func IsValidInputPath(pathStr string) bool {
	fileInfo, err := os.Stat(pathStr)

	return err == nil && fileInfo.IsDir()
}

func IsValidOutputPath(pathStr string) bool {
	return os.MkdirAll(pathStr, os.ModePerm) == nil
}

// TODO: should we use regular expressions in this check?
func IsValidString(input string) bool {
	return input != ""
}

// TODO: this is ugly :)
func IsValidProposalRevisionHeight(input string) bool {
	return isPositiveInt(input)
}

// TODO: this is ugly :)
func IsValidProposalId(input string) bool {
	return isPositiveInt(input)
}

func isPositiveInt(input string) bool {
	revHeight, err := strconv.Atoi(input)
	if err != nil || revHeight < 0 {
		return false
	}

	return true
}

func IsValidDateTime(input string) (time.Time, bool) {
	t, err := time.Parse(time.RFC3339Nano, input)
	if err != nil {
		return time.Now().UTC(), false
	}

	return t, true
}

// TODO: basic validation, expects only one coin and its amount
func IsValidDeposit(input string) bool {
	matches := reDecCoin.FindStringSubmatch(input)
	if matches == nil || len(matches) != 3 {
		return false
	}

	return true
}
