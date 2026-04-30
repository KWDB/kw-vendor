// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build !noasm

package memory

import "unsafe"

//go:noescape
func _memset_avx2(buf, len, c unsafe.Pointer)

/* douzt: the below code fails checkptr check, the main reason is uintptr should not
   be applied to integer numbers.  (This does NOT look good to me as well.)
   After a very shallow look into the asm code, it seems that arrow package has its
   way to handle memory, so we add ptr-checking free comment here.
   It's said that uintptr and unsafe.Pointer should be always in one same line, so we remove the vars.
*/

//go:nocheckptr
func memory_memset_avx2(buf []byte, c byte) {
	if len(buf) == 0 {
		return
	}

	if len(buf) > 2000 || isMultipleOfPowerOf2(len(buf), 256) {
		_memset_avx2(unsafe.Pointer(&buf[0]), unsafe.Pointer(uintptr(len(buf))), unsafe.Pointer(uintptr(c)))
	} else {
		_memset_sse4(unsafe.Pointer(&buf[0]), unsafe.Pointer(uintptr(len(buf))), unsafe.Pointer(uintptr(c)))
	}
}
