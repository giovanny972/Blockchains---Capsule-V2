package types

import (
	"encoding/json"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"
)

var _ depinject.OnePerModuleType = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

func init() {
	appmodule.Register(
		&Module{},
		appmodule.Provide(ProvideModule),
	)
}

type Module struct {
	// Authority defines the custom module authority. If not set, defaults to the governance module.
	Authority string `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
}

// IsAppModule implements the appmodule.AppModule interface.
func (am Module) IsAppModule() {}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am Module) IsOnePerModuleType() {}

func (x *Module) Reset() {
	*x = Module{}
}

func (x *Module) String() string {
	return "timecapsule module"
}

func (x *Module) ProtoMessage() {}

func (x *Module) Descriptor() ([]byte, []int) {
	return nil, nil
}

// DefaultGenesis returns default genesis state as raw bytes for the module.
func (Module) DefaultGenesis() json.RawMessage {
	return nil
}