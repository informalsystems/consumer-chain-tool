# consumer-chain-tool
The purpose of the tool is to produce an output in the form of proposal and genesis files. In that way, the process of starting the CosmWasm consumer chain with the pre-deployed smart contract codes is simplified. The process of creating the proposal and genesis data should be done in the following steps:
1. The proposer runs a prepare-proposal tool command which generates the genesis.json and proposal.json file. All the sections within the genesis file are populated with the final data, except for the ccvconsumer section, which represents the consumer module. The consumer section will be finalized in a later step. The proposal file contains several fields, among which there are the hashes of the genesis file and consumer binary file which will be used to run the consumer chain. The description field of the proposal file should contain a link to the location from where the tool output and the source code of the wasm contract can be downloaded. 
2. After the first step, when the proposal.json is created, the proposer manually submits the 'create consumer chain' proposal to the provider chain.
3. After the proposal is submitted, validators and all the interested parties can optionally run a verify-proposal command of the tool to verify the genesis data. This is done mostly to check if the pre-deployed smart contract codes match the source which is uploaded by the proposer. To do so, a user running this command will first download the contract source codes, review them and build them. The built smart contracts will be input for the verify-proposal command. Afterwards, the tool will check if the hash of the regenerated genesis matches the one from the proposal and the user can decide whether to vote for the proposal or not.
4. Finally, validators run a finalize-genesis command, which will generate the consumer binary and the final genesis file by adding proper data in ccvconsumer section. Validators can then use this genesis and binary to run the consumer chain. This step also requires for smart contracts to be built and given as a command input.

Note: parameters of the verify-proposal and finalize-genesis commands must have the same values as the corresponding parameters of the prepare-proposal command, e.g. if the CHAIN_ID in the prepare-proposal is 'wasm' then the CHAIN_ID in the verify-proposal must be 'wasm' as well. 

## prepare-proposal command
```
docker run --rm \
    -v <LOCATION_OF_SMART_CONTRACTS_BINARIES>:/contract_binaries \
    -v <TOOL_OUTPUT_LOCATION>:/tool_output \
    ics/cli consumer-chain-tool prepare-proposal <CHAIN_ID> <MULTISIG_ADDRESS> <PROPOSAL_TITLE> <PROPOSAL_DESCRIPTION> <REVISION_NUMBER> <REVISION_HEIGHT> <SPAWN_TIME> <DEPOSIT>
```

Input parameters:
- LOCATION_OF_SMART_CONTRACTS_BINARIES - The location of the directory that contains the compiled smart contracts .wasm binaries.
- TOOL_OUTPUT_LOCATION - The location of the directory where the resulting genesis.json and proposal.json will be saved.
- CHAIN_ID - The proposed chain-id of the new consumer chain must be different from all the other consumer chain ids of the executing provider chain.
- MULTISIG_ADDRESS - The multi-signature address that will have the permission to instantiate the contracts from the set of pre-deployed codes.
- PROPOSAL_TITLE - The title of the proposal.
- PROPOSAL_DESCRIPTION - The proposal description should contain the publicly available link where the contract's source code and the output of this command are uploaded by the proposer.
- REVISION_NUMBER - The revision that the client is currently on.
- REVISION_HEIGHT - The height within the given revision.
- SPAWN_TIME - The time on the provider chain at which the consumer chain genesis is finalized and all the validators will be responsible for starting heir consumer chain validator node. 
- DEPOSIT - The amount of tokens for the initial proposal deposit.

Examlpe: 
```
docker run --rm \
	-v $HOME/contract_binaries:/contract_binaries \
	-v $HOME/cli_tool_output:/tool_output \
	ics/cli consumer-chain-tool prepare-proposal wasm wasm1ykqt29d4ekemh5pc0d2wdayxye8yqupttf6vyz "CosmWasm consumer" "Contracts code location: https://mysharedlocation/proposal_data" 4 0 2022-06-01T09:10:00Z 10000001stake
```

## verify-proposal command
```
docker run --rm \
    -v <LOCATION_OF_SMART_CONTRACTS_BINARIES>:/contract_binaries \
    -v <TOOL_OUTPUT_LOCATION>:/tool_output \
    ics/cli consumer-chain-tool verify-proposal <CHAIN_ID> <MULTISIG_ADDRESS> <GENESIS_HASH> <CONSUMER_BINARY_HASH> <SPAWN_TIME> 
```

Input parameters:
- LOCATION_OF_SMART_CONTRACTS_BINARIES - The location of the directory that contains the compiled smart contracts .wasm binaries.
- TOOL_OUTPUT_LOCATION - TThe location of the directory where the verification data will be saved.
- CHAIN_ID - The proposed chain-id of the new consumer chain must be different from all the other consumer chain ids of the executing provider chain.
- MULTISIG_ADDRESS - The multi-signature address that will have the permission to instantiate the contracts from the set of pre-deployed codes.
- GENESIS_HASH - The hash of the genesis file can be obtained by querying the proposal previously submitted to the provider chain.
- CONSUMER_BINARY_HASH - The hash of the consumer binary can be obtained by querying the proposal previously submitted to the provider chain.
- SPAWN_TIME - The time on the provider chain at which the consumer chain genesis is finalized and all the validators will be responsible for starting their consumer chain validator node. 

Example:
```
docker run --rm \
	-v $HOME/contract_binaries:/contract_binaries \
	-v $HOME/cli_tool_output:/tool_output \
	ics/cli consumer-chain-tool verify-proposal wasm wasm1ykqt29d4ekemh5pc0d2wdayxye8yqupttf6vyz 5e637f4dbc6d6fb4b950ee259b13594deebfd7f92c68644d1b2264f2daa1b9df 09184916f3e85aa6fa24d3c12f1e5465af2214f13db265a52fa9f4617146dea5 2022-06-01T09:10:00Z
```

## finalize-genesis command
```
docker run --rm --network="host" \
    -v <LOCATION_OF_SMART_CONTRACTS_BINARIES>:/contract_binaries \
    -v <PROVIDER_BINARY_PATH>:/go/bin/interchain-security-pd \
    -v <TOOL_OUTPUT_LOCATION>:/tool_output \
    ics/cli consumer-chain-tool finalize-genesis <CHAIN_ID> <MULTISIG_ADDRESS> <PROPOSAL_ID> <PROVIDER_NODE_ADDRESS>
```

Input parameters:
- LOCATION_OF_SMART_CONTRACTS_BINARIES - The location of the directory that contains the compiled smart contracts .wasm binaries.
- PROVIDER_BINARY_PATH - The path to the provider binary.
- TOOL_OUTPUT_LOCATION - The location of the directory where the final genesis.json and consumer binary will be saved. The validators will use the outputs to start the consumer chain.
- CHAIN_ID - The proposed chain-id of the new consumer chain must be different from all the other consumer chain ids of the executing provider chain.
- MULTISIG_ADDRESS - The multi-signature address that will have the permission to instantiate the contracts from the set of pre-deployed codes.
- PROPOSAL_ID - The ID of the proposal submitted to the provider chain whose data will be used to verify if the inputs of this command match the ones from the proposal.
- PROVIDER_NODE_ADDRESS - This represents the address of the provider chain node in the following format: tcp://IP_ADDRESS:PORT_NUMBER. This address is used to query the provider chain to obtain the consumer section for the genesis file.

Example:
```
docker run --rm --network="host" \
	-v $HOME/contract_binaries:/contract_binaries \
	-v $HOME/go/src/consumer-chain-tool/gaiad:/go/bin/interchain-security-pd \
	-v $HOME/cli_tool_output:/tool_output \
	ics/cli consumer-chain-tool finalize-genesis wasm wasm1ykqt29d4ekemh5pc0d2wdayxye8yqupttf6vyz 1 tcp://ec2-13-40-188-91.eu-west-2.compute.amazonaws.com:1478
```