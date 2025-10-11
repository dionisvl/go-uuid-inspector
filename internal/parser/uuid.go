package parser

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	"github.com/google/uuid"
)

// UUIDInfo contains parsed UUID information
type UUIDInfo struct {
	Input     string
	Version   string
	Variant   string
	Canonical string
	Base58    string
	Base32    string
	Hex       string
	URN       string
	Time      string
	Error     string
}

// Parse tries to parse UUID from various formats and returns detailed info
func Parse(input string) UUIDInfo {
	info := UUIDInfo{Input: strings.TrimSpace(input)}

	u, err := parseUUID(info.Input)
	if err != nil {
		info.Error = fmt.Sprintf("Invalid UUID: %v", err)
		return info
	}

	version := u.Version()
	variant := u.Variant()

	info.Version = fmt.Sprintf("%d", version)
	info.Variant = variantString(variant)
	info.Canonical = u.String()
	info.Hex = strings.ToUpper(hex.EncodeToString(u[:]))
	info.URN = u.URN()
	info.Base58 = base58.Encode(u[:])
	info.Base32 = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(u[:])

	// Extract time for time-based UUIDs
	if version == 1 || version == 6 || version == 7 {
		if t, err := extractTime(u, version); err == nil {
			info.Time = t.Format(time.RFC3339Nano)
		}
	}

	return info
}

// parseUUID tries to parse UUID from various formats
func parseUUID(input string) (uuid.UUID, error) {
	// Try standard formats first (canonical, without dashes, URN, braces)
	u, err := uuid.Parse(input)
	if err == nil {
		return u, nil
	}

	// Try Base58
	decoded := base58.Decode(input)
	if len(decoded) == 16 {
		var u uuid.UUID
		copy(u[:], decoded)
		if isValidUUID(u) {
			return u, nil
		}
	}

	// Try Base32 (with padding)
	decoded, err = base32.StdEncoding.DecodeString(strings.ToUpper(input))
	if err == nil && len(decoded) == 16 {
		var u uuid.UUID
		copy(u[:], decoded)
		if isValidUUID(u) {
			return u, nil
		}
	}

	// Try Base32 (without padding)
	decoded, err = base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(input))
	if err == nil && len(decoded) == 16 {
		var u uuid.UUID
		copy(u[:], decoded)
		if isValidUUID(u) {
			return u, nil
		}
	}

	return uuid.UUID{}, fmt.Errorf("unable to parse UUID from input")
}

// isValidUUID checks if the UUID has valid version and variant
func isValidUUID(u uuid.UUID) bool {
	version := u.Version()
	variant := u.Variant()

	// Valid versions are 1-8
	if version < 1 || version > 8 {
		return false
	}

	// Must be RFC4122 variant
	if variant != uuid.RFC4122 {
		return false
	}

	return true
}

func variantString(v uuid.Variant) string {
	switch v {
	case uuid.Invalid:
		return "Invalid"
	case uuid.RFC4122:
		return "RFC 4122"
	case uuid.Reserved:
		return "Reserved (Microsoft)"
	case uuid.Future:
		return "Reserved (Future)"
	default:
		return "Unknown"
	}
}

func extractTime(u uuid.UUID, version uuid.Version) (time.Time, error) {
	if version == 7 {
		// UUID v7: first 48 bits are Unix timestamp in milliseconds
		timestamp := uint64(u[0])<<40 | uint64(u[1])<<32 | uint64(u[2])<<24 |
			uint64(u[3])<<16 | uint64(u[4])<<8 | uint64(u[5])
		sec := int64(timestamp / 1000)
		nsec := int64((timestamp % 1000) * 1000000)
		return time.Unix(sec, nsec), nil
	}

	if version == 1 || version == 6 {
		// UUID v1 and v6: 60-bit timestamp (100-nanosecond intervals since Oct 15, 1582)
		// v6 reorders the timestamp fields for better sorting
		var timestamp uint64

		if version == 6 {
			// v6: time_hi is in bytes 0-3, time_mid in bytes 4-5, time_low in bytes 6-7
			timeHi := uint64(u[0])<<24 | uint64(u[1])<<16 | uint64(u[2])<<8 | uint64(u[3])
			timeMid := uint64(u[4])<<8 | uint64(u[5])
			timeLow := uint64(u[6]&0x0F)<<8 | uint64(u[7])
			timestamp = (timeHi << 28) | (timeMid << 12) | timeLow
		} else {
			// v1: time_low in bytes 0-3, time_mid in bytes 4-5, time_hi in bytes 6-7
			timeLow := uint64(u[0])<<24 | uint64(u[1])<<16 | uint64(u[2])<<8 | uint64(u[3])
			timeMid := uint64(u[4])<<8 | uint64(u[5])
			timeHi := uint64(u[6]&0x0F)<<8 | uint64(u[7])
			timestamp = (timeHi << 48) | (timeMid << 32) | timeLow
		}

		// Convert from UUID timestamp to Unix timestamp
		// UUID epoch: 1582-10-15 00:00:00
		// Unix epoch: 1970-01-01 00:00:00
		const uuidEpochOffset = 122192928000000000 // 100-nanosecond intervals
		unixTime := int64(timestamp - uuidEpochOffset)
		sec := unixTime / 10000000
		nsec := (unixTime % 10000000) * 100
		return time.Unix(sec, nsec).UTC(), nil
	}

	return time.Time{}, fmt.Errorf("time extraction not supported for version %d", version)
}
