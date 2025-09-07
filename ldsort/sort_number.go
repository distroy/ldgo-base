/*
 * Copyright (C) distroy
 */

package ldsort

type (
	ints   = slice[int]
	int64s = slice[int64]
	int32s = slice[int32]
	int16s = slice[int16]
	int8s  = slice[int8]

	uints    = slice[uint]
	uint64s  = slice[uint64]
	uint32s  = slice[uint32]
	uint16s  = slice[uint16]
	uint8s   = slice[uint8]
	uintptrs = slice[uintptr]

	float32s = slice[float32]
	float64s = slice[float64]
)

func SortInts(a []int)              { internalSort(ints(a)) }
func UniqInts(a []int) []int        { return a[:uniq(ints(a))] }
func IsSortedInts(a []int) bool     { return internalIsSorted(ints(a)) }
func SearchInts(a []int, x int) int { return templateSearch(a, x) }
func IndexInts(a []int, x int) int  { return templateIndex(a, x) }

func SortInt64s(a []int64)                { internalSort(int64s(a)) }
func UniqInt64s(a []int64) []int64        { return a[:uniq(int64s(a))] }
func IsSortedInt64s(a []int64) bool       { return internalIsSorted(int64s(a)) }
func SearchInt64s(a []int64, x int64) int { return templateSearch(a, x) }
func IndexInt64s(a []int64, x int64) int  { return templateIndex(a, x) }

func SortInt32s(a []int32)                { internalSort(int32s(a)) }
func UniqInt32s(a []int32) []int32        { return a[:uniq(int32s(a))] }
func IsSortedInt32s(a []int32) bool       { return internalIsSorted(int32s(a)) }
func SearchInt32s(a []int32, x int32) int { return templateSearch(a, x) }
func IndexInt32s(a []int32, x int32) int  { return templateIndex(a, x) }

func SortInt16s(a []int16)                { internalSort(int16s(a)) }
func UniqInt16s(a []int16) []int16        { return a[:uniq(int16s(a))] }
func IsSortedInt16s(a []int16) bool       { return internalIsSorted(int16s(a)) }
func SearchInt16s(a []int16, x int16) int { return templateSearch(a, x) }
func IndexInt16s(a []int16, x int16) int  { return templateIndex(a, x) }

func SortInt8s(a []int8)               { internalSort(int8s(a)) }
func UniqInt8s(a []int8) []int8        { return a[:uniq(int8s(a))] }
func IsSortedInt8s(a []int8) bool      { return internalIsSorted(int8s(a)) }
func SearchInt8s(a []int8, x int8) int { return templateSearch(a, x) }
func IndexInt8s(a []int8, x int8) int  { return templateIndex(a, x) }

func SortUints(a []uint)               { internalSort(uints(a)) }
func UniqUints(a []uint) []uint        { return a[:uniq(uints(a))] }
func IsSortedUints(a []uint) bool      { return internalIsSorted(uints(a)) }
func SearchUints(a []uint, x uint) int { return templateSearch(a, x) }
func IndexUints(a []uint, x uint) int  { return templateIndex(a, x) }

func SortUint64s(a []uint64)                 { internalSort(uint64s(a)) }
func UniqUint64s(a []uint64) []uint64        { return a[:uniq(uint64s(a))] }
func IsSortedUint64s(a []uint64) bool        { return internalIsSorted(uint64s(a)) }
func SearchUint64s(a []uint64, x uint64) int { return templateSearch(a, x) }
func IndexUint64s(a []uint64, x uint64) int  { return templateIndex(a, x) }

func SortUint32s(a []uint32)                 { internalSort(uint32s(a)) }
func UniqUint32s(a []uint32) []uint32        { return a[:uniq(uint32s(a))] }
func IsSortedUint32s(a []uint32) bool        { return internalIsSorted(uint32s(a)) }
func SearchUint32s(a []uint32, x uint32) int { return templateSearch(a, x) }
func IndexUint32s(a []uint32, x uint32) int  { return templateIndex(a, x) }

func SortUint16s(a []uint16)                 { internalSort(uint16s(a)) }
func UniqUint16s(a []uint16) []uint16        { return a[:uniq(uint16s(a))] }
func IsSortedUint16s(a []uint16) bool        { return internalIsSorted(uint16s(a)) }
func SearchUint16s(a []uint16, x uint16) int { return templateSearch(a, x) }
func IndexUint16s(a []uint16, x uint16) int  { return templateIndex(a, x) }

func SortUint8s(a []uint8)                { internalSort(uint8s(a)) }
func UniqUint8s(a []uint8) []uint8        { return a[:uniq(uint8s(a))] }
func IsSortedUint8s(a []uint8) bool       { return internalIsSorted(uint8s(a)) }
func SearchUint8s(a []uint8, x uint8) int { return templateSearch(a, x) }
func IndexUint8s(a []uint8, x uint8) int  { return templateIndex(a, x) }

func SortUintptrs(a []uintptr)                  { internalSort(uintptrs(a)) }
func UniqUintptrs(a []uintptr) []uintptr        { return a[:uniq(uintptrs(a))] }
func IsSortedUintptrs(a []uintptr) bool         { return internalIsSorted(uintptrs(a)) }
func SearchUintptrs(a []uintptr, x uintptr) int { return templateSearch(a, x) }
func IndexUintptrs(a []uintptr, x uintptr) int  { return templateIndex(a, x) }

func SortFloat64s(a []float64)                  { internalSort(float64s(a)) }
func UniqFloat64s(a []float64) []float64        { return a[:uniq(float64s(a))] }
func IsSortedFloat64s(a []float64) bool         { return internalIsSorted(float64s(a)) }
func SearchFloat64s(a []float64, x float64) int { return templateSearch(a, x) }
func IndexFloat64s(a []float64, x float64) int  { return templateIndex(a, x) }

func SortFloat32s(a []float32)                  { internalSort(float32s(a)) }
func UniqFloat32s(a []float32) []float32        { return a[:uniq(float32s(a))] }
func IsSortedFloat32s(a []float32) bool         { return internalIsSorted(float32s(a)) }
func SearchFloat32s(a []float32, x float32) int { return templateSearch(a, x) }
func IndexFloat32s(a []float32, x float32) int  { return templateIndex(a, x) }
