# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

The Cosmos SDK is a modular framework for building blockchain applications in Go. It uses CometBFT (formerly Tendermint) for consensus and provides a comprehensive set of modules for building application-specific blockchains.

## Build System and Development Commands

### Essential Commands

Build the simapp binary:
```bash
make build
# Binary will be in build/simd
```

Install the simapp binary:
```bash
make install
# Installs simd to $GOPATH/bin
```

Run all tests:
```bash
make test
# Runs unit tests across all modules
```

Run tests for a specific package:
```bash
go test ./x/bank/...
```

Run integration tests:
```bash
make test-integration
```

Run E2E tests:
```bash
make test-e2e
```

### Linting and Code Quality

Run linting:
```bash
make lint
# Uses golangci-lint v1.54.2
```

Fix linting issues:
```bash
make lint-fix
```

Check for security vulnerabilities:
```bash
make vulncheck
```

### Protobuf Development

Generate protobuf files:
```bash
make proto-gen
```

Format protobuf files:
```bash
make proto-format
```

Lint protobuf files:
```bash
make proto-lint
```

Check for breaking proto changes:
```bash
make proto-check-breaking
```

### Simulation Testing

Run determinism test:
```bash
make test-sim-nondeterminism
```

Run multi-seed simulation:
```bash
make test-sim-multi-seed-short
make test-sim-multi-seed-long
```

### Local Development

Initialize a local single-node network:
```bash
make init-simapp
simd start
```

Start a 4-node local testnet:
```bash
make localnet-start
```

## Architecture

### Module Structure

The Cosmos SDK follows a modular architecture where each module is self-contained:

- **baseapp/**: Core ABCI application logic and state machine
- **x/**: Standard modules (auth, bank, staking, gov, etc.)
- **simapp/**: Example application demonstrating module integration
- **runtime/**: Application wiring and dependency injection
- **client/**: Client libraries and CLI tooling
- **codec/**: Encoding/decoding utilities (Amino, Protobuf)
- **crypto/**: Cryptographic primitives and key management
- **store/**: State storage interfaces and implementations
- **types/**: Core SDK types and utilities

### Key Components

- **BaseApp**: Core application that handles ABCI interactions
- **Keeper Pattern**: Each module has a keeper that manages state
- **AnteHandler**: Middleware for transaction preprocessing
- **PostHandler**: Middleware for transaction postprocessing
- **Module Manager**: Orchestrates module initialization and execution
- **Msg Service Router**: Routes messages to appropriate handlers

### State Management

- Uses **IAVL trees** for authenticated state storage
- **KVStore interface** abstracts storage operations
- **Collections** package provides type-safe state management
- **Prefix stores** for module state isolation

### Module Communication

- **Inter-module communication** via keeper interfaces
- **Events** for cross-module notifications
- **Capabilities** for object-capability security model

## Go Module Structure

This is a multi-module Go workspace:
- Root module: `github.com/cosmos/cosmos-sdk`
- Submodules: `collections/`, `core/`, `depinject/`, `errors/`, `log/`, `math/`, `orm/`, `store/`, `tools/`, `x/`

Each module has its own `go.mod` file and can be versioned independently.

## Testing Guidelines

### Unit Tests
- Each package should have comprehensive unit tests
- Use table-driven tests where appropriate
- Mock external dependencies using interfaces
- Test files should end with `_test.go`

### Integration Tests
- Located in `tests/integration/`
- Test inter-module interactions
- Use the test network utilities in `testutil/network/`

### Simulation Tests
- Randomized testing of application behavior
- Located in module `simulation/` directories
- Use `make test-sim-*` commands

## Module Development

### Creating a New Module

1. Create directory structure in `x/yourmodule/`
2. Implement required interfaces: `AppModule`, `AppModuleBasic`
3. Define types in `types/` directory
4. Implement keeper in `keeper/` directory
5. Add CLI commands in `client/cli/`
6. Write comprehensive tests

### Module Requirements

- **Genesis handling**: Import/export state
- **Query handlers**: gRPC and REST endpoints
- **Message handlers**: Transaction processing
- **Events**: Emit relevant events
- **Parameters**: Module configuration
- **Migrations**: State migration between versions

## Important Notes

- Always run `make lint` before submitting changes
- Use `make proto-gen` after modifying `.proto` files
- Follow the existing code patterns and conventions
- Write comprehensive tests for new functionality
- Update relevant documentation when adding features
- Check `CONTRIBUTING.md` for detailed contribution guidelines

## Debugging

- Use `make build-linux-amd64` or `make build-linux-arm64` for cross-compilation
- Set `COSMOS_BUILD_OPTIONS=debug` for debug builds
- Use the `simd debug` command for debugging utilities
- Enable trace logging with appropriate log levels

## Module Locations

Key modules are located in `x/`:
- `auth/`: Account management and authentication
- `bank/`: Token transfers and balances  
- `staking/`: Proof-of-stake functionality
- `gov/`: On-chain governance
- `distribution/`: Fee and reward distribution
- `slashing/`: Validator punishment
- `evidence/`: Byzantine behavior handling
- `crisis/`: Invariant checking
- `upgrade/`: On-chain upgrades
- `params/`: Parameter management