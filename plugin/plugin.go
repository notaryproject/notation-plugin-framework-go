// Package plugin provides the tooling to use the notation plugin.
//
// includes a CLIManager and a CLIPlugin implementation.
package plugin

import (
	"context"
)

// GenericPlugin is the base requirement to be a plugin.
type GenericPlugin interface {
	// GetMetadata returns the metadata information of the plugin.
	GetMetadata(ctx context.Context, req *GetMetadataRequest) (*GetMetadataResponse, error)
}

// SignPlugin defines the required methods to be a SignPlugin.
type SignPlugin interface {
	GenericPlugin

	// DescribeKey returns the KeySpec of a key.
	DescribeKey(ctx context.Context, req *DescribeKeyRequest) (*DescribeKeyResponse, error)

	// GenerateSignature generates the raw signature based on the request.
	GenerateSignature(ctx context.Context, req *GenerateSignatureRequest) (*GenerateSignatureResponse, error)

	// GenerateEnvelope generates the Envelope with signature based on the
	// request.
	GenerateEnvelope(ctx context.Context, req *GenerateEnvelopeRequest) (*GenerateEnvelopeResponse, error)
}

// VerifyPlugin defines the required method to be a VerifyPlugin.
type VerifyPlugin interface {
	GenericPlugin

	// VerifySignature validates the signature based on the request.
	VerifySignature(ctx context.Context, req *VerifySignatureRequest) (*VerifySignatureResponse, error)
}

// Plugin defines required methods to be a Plugin.
type Plugin interface {
	SignPlugin
	VerifyPlugin
}
