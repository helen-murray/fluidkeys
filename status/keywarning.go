package status

import (
	"fmt"
	"time"

	"github.com/fluidkeys/fluidkeys/colour"
)

type WarningType int

const (
	// If you add a type, remember to handle it in all the switch statements.
	UnsetType WarningType = 0

	PrimaryKeyDueForRotation     = 1
	PrimaryKeyOverdueForRotation = 2
	PrimaryKeyExpired            = 3
	PrimaryKeyNoExpiry           = 4
	PrimaryKeyLongExpiry         = 5

	NoValidEncryptionSubkey  = 6
	SubkeyDueForRotation     = 7
	SubkeyOverdueForRotation = 8
	SubkeyNoExpiry           = 9
	SubkeyLongExpiry         = 10

	MissingPreferredSymmetricAlgorithms    = 11
	WeakPreferredSymmetricAlgorithms       = 12
	UnsupportedPreferredSymmetricAlgorithm = 13

	MissingPreferredHashAlgorithms    = 14
	WeakPreferredHashAlgorithms       = 15
	UnsupportedPreferredHashAlgorithm = 16

	MissingPreferredCompressionAlgorithms    = 17
	UnsupportedPreferredCompressionAlgorithm = 18
	MissingUncompressedPreference            = 19 // Implementations MUST implement uncompressed data.

	WeakSelfSignatureHash          = 20
	WeakSubkeyBindingSignatureHash = 21
)

type KeyWarning struct {
	Type WarningType

	SubkeyId          uint64
	DaysUntilExpiry   uint
	DaysSinceExpiry   uint
	CurrentValidUntil *time.Time
	Detail            string
}

func (w KeyWarning) String() string {
	switch w.Type {
	case UnsetType:
		return ""

	case PrimaryKeyDueForRotation:
		return "Primary key needs extending"

	case PrimaryKeyOverdueForRotation:
		return colour.Danger("Primary key needs extending now (" + countdownUntilExpiry(w.DaysUntilExpiry) + ")")

	case PrimaryKeyExpired:
		return colour.Danger("Primary key " + relativeExpiryDate(w.DaysSinceExpiry))

	case PrimaryKeyNoExpiry:
		return "Primary key never expires"

	case PrimaryKeyLongExpiry:
		return "Primary key expires too far in the future"

	case NoValidEncryptionSubkey:
		return colour.Danger("Missing encryption subkey")

	case SubkeyDueForRotation:
		return "Encryption subkey needs rotating"

	case SubkeyOverdueForRotation:
		return colour.Danger("Encryption subkey needs rotating now (" + countdownUntilExpiry(w.DaysUntilExpiry) + ")")

	case SubkeyNoExpiry:
		return "Encryption subkey never expires"

	case SubkeyLongExpiry:
		return "Encryption subkey expires too far in the future"

	case MissingPreferredSymmetricAlgorithms:
		return "Missing cipher preferences"

	case WeakPreferredSymmetricAlgorithms:
		return fmt.Sprintf("Cipher preferences could be stronger (currently: %s)", w.Detail)

	case UnsupportedPreferredSymmetricAlgorithm:
		return fmt.Sprintf("Fluidkeys doesn't support %s cipher", w.Detail)

	case MissingPreferredHashAlgorithms:
		return "Missing hash preferences"

	case WeakPreferredHashAlgorithms:
		return fmt.Sprintf("Hash preferences could be stronger (currently: %s)", w.Detail)

	case UnsupportedPreferredHashAlgorithm:
		return fmt.Sprintf("Fluidkeys doesn't support %s hash", w.Detail)

	case MissingPreferredCompressionAlgorithms:
		return "Missing compression preferences"

	case MissingUncompressedPreference:
		return "Key does not support uncompressed data"

	case UnsupportedPreferredCompressionAlgorithm:
		return fmt.Sprintf("Fluidkeys doesn't support %s compression", w.Detail)

	case WeakSelfSignatureHash:
		return fmt.Sprintf("Weak hash %s used for self signature", w.Detail)

	case WeakSubkeyBindingSignatureHash:
		return fmt.Sprintf("Weak hash %s used for subkey binding signature", w.Detail)
	}

	return fmt.Sprintf("KeyWarning{Type=%d}", w.Type)
}

func (w KeyWarning) IsAboutSubkey() bool {
	switch w.Type {
	case
		SubkeyDueForRotation,
		SubkeyOverdueForRotation,
		SubkeyNoExpiry,
		SubkeyLongExpiry,
		WeakSubkeyBindingSignatureHash:
		return true
	}
	return false
}

func (w KeyWarning) IsAboutPrimaryKey() bool {
	switch w.Type {
	case
		PrimaryKeyDueForRotation,
		PrimaryKeyOverdueForRotation,
		PrimaryKeyExpired,
		PrimaryKeyNoExpiry,
		PrimaryKeyLongExpiry,
		MissingPreferredSymmetricAlgorithms,
		WeakPreferredSymmetricAlgorithms,
		UnsupportedPreferredSymmetricAlgorithm,
		MissingPreferredHashAlgorithms,
		WeakPreferredHashAlgorithms,
		UnsupportedPreferredHashAlgorithm,
		MissingPreferredCompressionAlgorithms,
		UnsupportedPreferredCompressionAlgorithm,
		WeakSelfSignatureHash:
		return true
	}
	return false
}

func ContainsWarningAboutPrimaryKey(warnings []KeyWarning) bool {
	for _, warning := range warnings {
		if warning.IsAboutPrimaryKey() {
			return true
		}
	}
	return false
}

func countdownUntilExpiry(days uint) string {
	switch days {
	case 0:
		return "expires today!"
	case 1:
		return "expires tomorrow!"
	default:
		return fmt.Sprintf("expires in %d days", days)
	}
}

func relativeExpiryDate(days uint) string {
	switch days {
	case 0:
		return "expired today"
	case 1:
		return "expired yesterday"
	case 2, 3, 4, 5, 6, 7, 8, 9:
		return fmt.Sprintf("expired %d days ago", days)
	default:
		return "has expired"
	}
}
