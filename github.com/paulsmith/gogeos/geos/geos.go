// Copyright (c) 2013 Paul Smith
// Package geos provides support for creating and manipulating spatial data.
// At its core, it relies on the GEOS C library for the implementation of
// spatial operations and geometric algorithms.
package geos

/*
#include "geos.h"
*/
import "C"

import (
	"fmt"
)

// XXX: store last error message from handler in a global var (chan?)

func LoadMessage() string {
	switch FunctionLoadStatus {
	case C.Load_ERROR_INIT:
		return "geos: function initGEOS_r load error"
	case C.Load_ERROR_VERSION:
		return "geos: function GEOSversion load error"
	case C.Load_ERROR_BUFFER:
		return "geos: function GEOSBuffer_r load error"
	case C.Load_ERROR_GEOM_DESTROY:
		return "geos: function GEOSGeom_destroy_r load error"
	case C.Load_ERROR_TOUCHES:
		return "geos: function GEOSTouches_r load error"
	case C.Load_ERROR_INTERSECT:
		return "geos: function GEOSIntersects_r load error"
	case C.Load_ERROR_WITHIN:
		return "geos: function GEOSWithin_r load error"
	case C.Load_ERROR_CONTAIN:
		return "geos: function GEOSContains_r load error"
	case C.Load_ERROR_EQUAL:
		return "geos: function GEOSEquals_r load error"
	case C.Load_ERROR_COVER:
		return "geos: function GEOSCovers_r load error"
	case C.Load_ERROR_AREA:
		return "geos: function GEOSArea_r load error"
	case C.Load_ERROR_DISTANCE:
		return "geos: function GEOSDistance_r load error"
	case C.Load_ERROR_CREATE_READER:
		return "geos: function GEOSWKTReader_create_r load error"
	case C.Load_ERROR_DESTORY_READER:
		return "geos: function GEOSWKTReader_destroy_r load error"
	case C.Load_ERROR_READ_READER:
		return "geos: function GEOSWKTReader_read_r load error"
	case C.Load_ERROR_CREATE_WRITER:
		return "geos: function GEOSWKTWriter_create_r load error"
	case C.Load_ERROR_DESTORY_WRITER:
		return "geos: function GEOSWKTWriter_destroy_r load error"
	case C.Load_ERROR_WRITE_WRITER:
		return "geos: function GEOSWKTWriter_write_r load error"
	default:
		return ""
	}
}

// Error gets the last error that occured in the GEOS C API as a Go error type.
func Error() error {
	if FunctionLoadStatus != C.Load_OK {
		return fmt.Errorf("geos: function load error")
	}
	return fmt.Errorf("geos: %s", C.GoString(C.gogeos_get_last_error()))
}
