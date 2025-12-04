# Contracts Directory

This folder contains Solidity smart contracts and related deployment tooling. It is
structured to make it easy to plug your existing contracts into the Go backend
without tightly coupling the codebase to any particular implementation.

* `src/` – place your Solidity source files here. Organise them by
  functionality or domain as you see fit.
* `deploy/` – optional helper scripts for deploying contracts. You may choose
  to use a framework such as Hardhat or Foundry; if so, initialise the
  project here. Scripts in this folder can compile and deploy contracts to
  Ethereum or zkSync networks.

## Integrating Contracts with Go

1. Compile your Solidity contracts and generate Go bindings using a tool such
   as `abigen` from the go-ethereum suite. The generated Go code should
   live outside of `internal` to keep it reusable; `internal/bridge/contracts`
   is a good place.
2. Define interfaces in Go that describe the functions you intend to call. This
   allows the rest of the application to depend on abstractions rather than
   concrete implementations.
3. Load contract addresses and ABI file paths from configuration (see
   `internal/config`). This enables switching between local, testnet and
   mainnet deployments without code changes.

If you prefer not to use Hardhat or any JavaScript tooling, you can skip
initialising a project in this folder. The key requirement is that the Go
project remains prepared to accept ABI/binding packages when you are ready
to integrate your contracts.