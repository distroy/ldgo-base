/*
 * Copyright (C) distroy
 */

package ldatomic

import "testing"

func TestUint(t *testing.T)    { testInt(t, NewUint) }
func TestUint8(t *testing.T)   { testInt(t, NewUint8) }
func TestUint16(t *testing.T)  { testInt(t, NewUint16) }
func TestUint32(t *testing.T)  { testInt(t, NewUint32) }
func TestUint64(t *testing.T)  { testInt(t, NewUint64) }
func TestUintptr(t *testing.T) { testInt(t, NewUintptr) }
