// Code generated by "stringer -type=indexScanOp -linecomment"; DO NOT EDIT.

package analyzer

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[indexScanOpEq-0]
	_ = x[indexScanOpNullSafeEq-1]
	_ = x[indexScanOpInSet-2]
	_ = x[indexScanOpNotInSet-3]
	_ = x[indexScanOpNotEq-4]
	_ = x[indexScanOpGt-5]
	_ = x[indexScanOpGte-6]
	_ = x[indexScanOpLt-7]
	_ = x[indexScanOpLte-8]
	_ = x[indexScanOpAnd-9]
	_ = x[indexScanOpOr-10]
	_ = x[indexScanOpIsNull-11]
	_ = x[indexScanOpIsNotNull-12]
	_ = x[indexScanOpSpatialEq-13]
	_ = x[indexScanOpFulltextEq-14]
}

const _indexScanOp_name = "=<=>=!=!=>>=<<=&&||IS NULLIS NOT NULLSpatialEqFulltextEq"

var _indexScanOp_index = [...]uint8{0, 1, 4, 5, 7, 9, 10, 12, 13, 15, 17, 19, 26, 37, 46, 56}

func (i indexScanOp) String() string {
	if i >= indexScanOp(len(_indexScanOp_index)-1) {
		return "indexScanOp(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _indexScanOp_name[_indexScanOp_index[i]:_indexScanOp_index[i+1]]
}