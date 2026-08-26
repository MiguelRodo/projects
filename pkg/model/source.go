package model

import (
	"errors"
	"strings"
)

var (
	// ErrEmptySourceID indicates the source identifier is empty.
	ErrEmptySourceID = errors.New("source id cannot be empty")
	// ErrEmptySourceLocation indicates the source location is empty.
	ErrEmptySourceLocation = errors.New("source location cannot be empty")
)

// SourceType indicates the category of an external source.
type SourceType string

const (
	// SourceTypeDoc represents a public documentation or specification source.
	SourceTypeDoc SourceType = "doc"
	// SourceTypeCatalogue represents a reusable repository catalogue or schema list.
	SourceTypeCatalogue SourceType = "catalogue"
	// SourceTypeAsset represents static templates or asset collections.
	SourceTypeAsset SourceType = "asset"
)

// SourceRef references an external repository-safe catalogue, document, or template source.
type SourceRef struct {
	ID          string     `json:"id" yaml:"id"`
	Name        string     `json:"name,omitempty" yaml:"name,omitempty"`
	Location    string     `json:"location" yaml:"location"`
	Type        SourceType `json:"type,omitempty" yaml:"type,omitempty"`
	Description string     `json:"description,omitempty" yaml:"description,omitempty"`
}

// Validate checks source reference fields.
func (s *SourceRef) Validate() error {
	s.ID = strings.TrimSpace(s.ID)
	s.Location = strings.TrimSpace(s.Location)
	if s.ID == "" {
		return ErrEmptySourceID
	}
	if s.Location == "" {
		return ErrEmptySourceLocation
	}
	if s.Type == "" {
		s.Type = SourceTypeDoc
	}
	return nil
}

// SourceCatalogue is a collection of versioned source references.
type SourceCatalogue struct {
	Version string      `json:"version" yaml:"version"`
	Sources []SourceRef `json:"sources" yaml:"sources"`
}

// Validate checks catalogue validity.
func (c *SourceCatalogue) Validate() error {
	for i := range c.Sources {
		if err := c.Sources[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}
