//go:build ignore
// +build ignore

// plus hand editing about timeval

/*
Input to cgo -godefs.
*/

package host

/*
#include <sys/time.h>
#include <utmpx.h>
*/
import "C"

type (
	Utmpx   C.struct_utmpx
	Timeval C.struct_timeval
)
