// Code generated by "stringer -type=NodeType"; DO NOT EDIT.

package parser

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TypeRoot-0]
	_ = x[TypeAny-1]
	_ = x[TypeText-2]
	_ = x[TypeLeaf-3]
}

const _NodeType_name = "TypeRootTypeAnyTypeTextTypeLeaf"

var _NodeType_index = [...]uint8{0, 8, 15, 23, 31}

func (i NodeType) String() string {
	if i < 0 || i >= NodeType(len(_NodeType_index)-1) {
		return "NodeType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _NodeType_name[_NodeType_index[i]:_NodeType_index[i+1]]
}
