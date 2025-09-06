/*
 * Copyright (C) distroy
 */

package ldbyte

func IsDigit(c byte) bool { return c >= '0' && c <= '9' }
func IsLower(c byte) bool { return c >= 'a' && c <= 'z' }
func IsUpper(c byte) bool { return c >= 'A' && c <= 'Z' }
func IsPrint(c byte) bool { return c >= 32 && c <= 126 }

func IsCtrl(c byte) bool {
	switch c {
	// in octal
	case 0000: // NUL
	case 0001: // SOH
	case 0002: // STX
	case 0003: // ETX
	case 0004: // EOT
	case 0005: // ENQ
	case 0006: // ACK
	case 0007: // BEL
	case 0010: // BS
	case 0011: // HT
	case 0012: // NL
	case 0013: // VT
	case 0014: // NP
	case 0015: // CR
	case 0016: // SO
	case 0017: // SI
	case 0020: // DLE
	case 0021: // DC1
	case 0022: // DC2
	case 0023: // DC3
	case 0024: // DC4
	case 0025: // NAK
	case 0026: // SYN
	case 0027: // ETB
	case 0030: // CAN
	case 0031: // EM
	case 0032: // SUB
	case 0033: // ESC
	case 0034: // FS
	case 0035: // GS
	case 0036: // RS
	case 0037: // US
	case 0177: // DEL
	default:
		return false
	}
	return true
}

func IsBlank(c byte) bool {
	switch c {
	case ' ':
	case '\t':
	default:
		return false
	}
	return true
}

func IsSpace(c byte) bool {
	switch c {
	case ' ':
	case '\t':
	case '\v':
	case '\f':
	case '\r':
	case '\n':
	default:
		return false
	}
	return true
}

func IsXDigit(c byte) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

func IsPunct(c byte) bool {
	switch c {
	// in octal
	case 0041: // ``!''
	case 0042: // ``"''
	case 0043: // ``#''
	case 0044: // ``$''
	case 0045: // ``%''
	case 0046: // ``&''
	case 0047: // ``'''
	case 0050: // ``(''
	case 0051: // ``)''
	case 0052: // ``*''
	case 0053: // ``+''
	case 0054: // ``,''
	case 0055: // ``-''
	case 0056: // ``.''
	case 0057: // ``/''
	case 0072: // ``:''
	case 0073: // ``;''
	case 0074: // ``<''
	case 0075: // ``=''
	case 0076: // ``>''
	case 0077: // ``?''
	case 0100: // ``@''
	case 0133: // ``[''
	case 0134: // ``\''
	case 0135: // ``]''
	case 0136: // ``^''
	case 0137: // ``_''
	case 0140: // ```''
	case 0173: // ``{''
	case 0174: // ``|''
	case 0175: // ``}''
	case 0176: // ``~''
	default:
		return false
	}
	return true
}

// IsAlpha = IsUpper || IsLower
func IsAlpha(c byte) bool { return IsUpper(c) || IsLower(c) }

// IsAlNum = IsAlpha || IsDigit
func IsAlNum(c byte) bool { return IsDigit(c) || IsAlpha(c) }
