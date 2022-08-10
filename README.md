# consumer-chain-tool
The purpose of the tool is to produce an output in the form of proposal and genesis files. In that way, the process of starting the CosmWasm consumer chain with the pre-deployed smart contract codes is simplified. The process of creating the proposal and genesis data should be done in the following steps:
1. The proposer runs a prepare-proposal tool command which generates the genesis.json and proposal.json file. All the sections within the genesis file are populated with the final data, except for the ccvconsumer section, which represents the consumer module. The consumer section will be finalized in a later step. The proposal file contains several fields, among which there are the hashes of the genesis file and consumer binary file which will be used to run the consumer chain. The description field of the proposal file should contain a link to the location from where the tool output and the source code of the wasm contract can be downloaded. 
2. After the first step, when the proposal.json is created, the proposer manually submits the 'create consumer chain' proposal to the provider chain.
3. After the proposal is submitted, validators and all the interested parties can optionally run a verify-proposal command of the tool to verify the genesis data. This is done mostly to check if the pre-deployed smart contract codes match the source which is uploaded by the proposer. To do so, a user running this command will first download the contract source codes, review them and build them. The built smart contracts will be input for the verify-proposal command. Afterwards, the tool will check if the hash of the regenerated genesis matches the one from the proposal and the user can decide whether to vote for the proposal or not.
4. Finally, validators run a finalize-genesis command, which will generate the consumer binary and the final genesis file by adding proper data in ccvconsumer section. Validators can then use this genesis and binary to run the consumer chain. This step also requires for smart contracts to be built and given as a command input.

Note: parameters of the verify-proposal and finalize-genesis commands must have the same values as the corresponding parameters of the prepare-proposal command, e.g. if the CHAIN_ID in the prepare-proposal is 'wasm' then the CHAIN_ID in the verify-proposal must be 'wasm' as well. 

## Prerequisits
Running the tool requires Docker to be installed. 

## prepare-proposal command
```
consumer-chain-tool prepare-proposal [CONTRACT-BINARIES-LOCATION] [CONSUMER-CHAIN-ID] [MULTISIG-ADDRESS] [TOOL-OUTPUT-LOCATION] [PROPOSAL-TITLE] [PROPOSAL-DESCRIPTION] [PROPOSAL-REVISION-HEIGHT] [PROPOSAL-REVISION-NUMBER] [PROPOSAL-SPAWN-TIME] [PROPOSAL-DEPOSIT]
```

Input parameters:
- CONTRACT-BINARIES-LOCATION - The location of the directory that contains the compiled smart contracts .wasm binaries.
- CONSUMER-CHAIN-ID - The proposed chain-id of the new consumer chain must be different from all the other consumer chain ids of the executing provider chain.
- MULTISIG-ADDRESS - The multi-signature address that will have the permission to instantiate the contracts from the set of pre-deployed codes.
- TOOL-OUTPUT-LOCATION - The location of the directory where the resulting genesis.json and proposal.json will be saved.
- PROPOSAL-TITLE - The title of the proposal.
- PROPOSAL-DESCRIPTION - The proposal description should contain the publicly available link where the contract's source code and the output of this command are uploaded by the proposer.
- PROPOSAL-REVISION-HEIGHT - The height within the given revision.
- PROPOSAL-REVISION-NUMBER - The revision that the client is currently on.
- PROPOSAL-SPAWN-TIME - The time on the provider chain at which the consumer chain genesis is finalized and all the validators will be responsible for starting heir consumer chain validator node. 
- PROPOSAL-DEPOSIT - The amount of tokens for the initial proposal deposit.

Examlpe: 
```
consumer-chain-tool prepare-proposal $HOME/contract_binaries wasm wasm1ykqt29d4ekemh5pc0d2wdayxye8yqupttf6vyz $HOME/cli_tool_output "CosmWasm consumer" "Contracts code location: https://mysharedlocation/proposal_data" 4 0 2022-06-01T09:10:00Z 10000001stake
```

## verify-proposal command
```
consumer-chain-tool verify-proposal [CONTRACT-BINARIES-LOCATION] [CONSUMER-CHAIN-ID] [MULTISIG-ADDRESS] [TOOL-OUTPUT-LOCATION] [PROPOSAL-GENESIS-HASH] [PROPOSAL-BINARY-HASH] [PROPOSAL-SPAWN-TIME]
```

Input parameters:
- CONTRACT-BINARIES-LOCATION - The location of the directory that contains the compiled smart contracts .wasm binaries.
- CONSUMER-CHAIN-ID - The proposed chain-id of the new consumer chain must be different from all the other consumer chain ids of the executing provider chain.
- MULTISIG-ADDRESS - The multi-signature address that will have the permission to instantiate the contracts from the set of pre-deployed codes.
- TOOL-OUTPUT-LOCATION - The location of the directory where the verification data will be saved.
- PROPOSAL-GENESIS-HASH - The hash of the genesis file can be obtained by querying the proposal previously submitted to the provider chain.
- PROPOSAL-BINARY-HASH - The hash of the consumer binary can be obtained by querying the proposal previously submitted to the provider chain.
- PROPOSAL-SPAWN-TIME - The time on the provider chain at which the consumer chain genesis is finalized and all the validators will be responsible for starting their consumer chain validator node. 

Example:
```
consumer-chain-tool verify-proposal $HOME/contract_binaries wasm wasm1ykqt29d4ekemh5pc0d2wdayxye8yqupttf6vyz $HOME/cli_tool_output 519df96a862c30f53e67b1277e6834ab4bd59dfdd08c781d1b7cf3813080fb28 09184916f3e85aa6fa24d3c12f1e5465af2214f13db265a52fa9f4617146dea5 2022-06-01T09:10:00Z
```

## finalize-genesis command
```
consumer-chain-tool finalize-genesis [CONTRACT-BINARIES-LOCATION] [CONSUMER-CHAIN-ID] [MULTISIG-ADDRESS] [TOOL-OUTPUT-LOCATION] [PROPOSAL-ID] [PROVIDER-NODE-ADDRESS] [PROVIDER-BINARY-PATH]
```

Input parameters:
- CONTRACT-BINARIES-LOCATION - The location of the directory that contains the compiled smart contracts .wasm binaries.
- CONSUMER-CHAIN-ID - The proposed chain-id of the new consumer chain must be different from all the other consumer chain ids of the executing provider chain.
- MULTISIG-ADDRESS - The multi-signature address that will have the permission to instantiate the contracts from the set of pre-deployed codes.
- TOOL-OUTPUT-LOCATION - The location of the directory where the final genesis.json and consumer binary will be saved. The validators will use the outputs to start the consumer chain.
- PROPOSAL-ID - The ID of the proposal submitted to the provider chain whose data will be used to verify if the inputs of this command match the ones from the proposal.
- PROVIDER-NODE-ADDRESS - This represents the address of the provider chain node in the following format: tcp://IP_ADDRESS:PORT_NUMBER. This address is used to query the provider chain to obtain the consumer section for the genesis file.
- PROVIDER_BINARY_PATH - The path to the provider binary.

Example:
```
consumer-chain-tool finalize-genesis $HOME/contract_binaries wasm wasm1ykqt29d4ekemh5pc0d2wdayxye8yqupttf6vyz $HOME/cli_tool_output 1 tcp://localhost:26657 $HOME/gaiad
```