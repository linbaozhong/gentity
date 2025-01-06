// Copyright © 2023 Linbaozhong. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"math"
	"time"
)

const (
	Nil = math.MaxUint64
	// NilString  = "NIL"
	// NilInt     = math.MinInt
	// NilInt8    = math.MinInt8
	// NilInt16   = math.MinInt16
	// NilInt32   = math.MinInt32
	// NilInt64   = math.MinInt64
	// NilUint    = math.MaxUint
	// NilUint8   = math.MaxUint8
	// NilUint16  = math.MaxUint16
	// NilUint32  = math.MaxUint32
	// NilUint64  = math.MaxUint64
	// NilFloat64 = math.MaxFloat64
	// NilFloat32 = math.MaxFloat32
)

var (
	NilTime = time.Time{}
)
