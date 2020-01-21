//**************************************************************************
//**
//**	##   ##    ##    ##   ##   ####     ####   ###     ###
//**	##   ##  ##  ##  ##   ##  ##  ##   ##  ##  ####   ####
//**	 ## ##  ##    ##  ## ##  ##    ## ##    ## ## ## ## ##
//**	 ## ##  ########  ## ##  ##    ## ##    ## ##  ###  ##
//**	  ###   ##    ##   ###    ##  ##   ##  ##  ##       ##
//**	   #    ##    ##    #      ####     ####   ##       ##
//**
//**	$Id: cmdlib.h 1610 2006-07-09 14:11:23Z dj_jl $
//**
//**	Copyright (C) 1999-2006 Jānis Legzdiņš
//**
//**	This program is free software; you can redistribute it and/or
//**  modify it under the terms of the GNU General Public License
//**  as published by the Free Software Foundation; either version 2
//**  of the License, or (at your option) any later version.
//**
//**	This program is distributed in the hope that it will be useful,
//**  but WITHOUT ANY WARRANTY; without even the implied warranty of
//**  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//**  GNU General Public License for more details.
//**
//**************************************************************************

#ifndef CMDLIB_H
#define CMDLIB_H

// HEADER FILES ------------------------------------------------------------

//	C headers
#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>
#include <string.h>
#include <ctype.h>

#ifdef _MSC_VER
// Disable some unwanted warnings
#pragma warning(disable : 4097) // typedef-name '' used as synonym for class-name ''
#pragma warning(disable : 4127) // Conditional expression is constant
#pragma warning(disable : 4244) // conversion to float, possible loss of data
#pragma warning(disable : 4284) // return type is not a UDT or reference to a UDT
#pragma warning(disable : 4291) // no matching operator delete found
#pragma warning(disable : 4305) // truncation from 'const double' to 'float'
#pragma warning(disable : 4512) // assignment operator could not be generated
#pragma warning(disable : 4514)	// unreferenced inline function has been removed
#pragma warning(disable : 4702) // unreachable code in inline expanded function
#pragma warning(disable : 4710) // inline function not expanded

// Workaround of for variable scoping
#define for		if (0) ; else for
#endif

#ifdef __BORLANDC__
// Disable some unwanted warnings
#pragma warn -8027		// inline function not expanded
#pragma warn -8071		// conversation may loose significant digits
#endif

#ifndef __GNUC__
#define __attribute__(whatever)
#endif

#if !defined _WIN32 && !defined DJGPP
#undef stricmp	//	Allegro defines them
#undef strnicmp
#define stricmp		strcasecmp
#define strnicmp	strncasecmp
#endif

namespace VavoomUtils {

// MACROS ------------------------------------------------------------------

#define MIN_VINT8	((vint8)-128)
#define MIN_VINT16	((vint16)-32768)
#define MIN_VINT32	((vint32)-2147483648)

#define MAX_VINT8	((vint8)0x7f)
#define MAX_VINT16	((vint16)0x7fff)
#define MAX_VINT32	((vint32)0x7fffffff)

// TYPES -------------------------------------------------------------------

#ifdef HAVE_INTTYPES_H
#include <inttypes.h>
typedef int8_t				vint8;
typedef uint8_t				vuint8;
typedef int16_t				vint16;
typedef uint16_t			vuint16;
typedef int32_t				vint32;
typedef uint32_t			vuint32;
/* Needs more changes to compile with MSVC
#elif defined _WIN32
typedef __int8				vint8;
typedef unsigned __int8		vuint8;
typedef __int16				vint16;
typedef unsigned __int16	vuint16;
typedef __int32				vint32;
typedef unsigned __int32	vuint32;
*/
#else
typedef char				vint8;
typedef unsigned char		vuint8;
typedef short				vint16;
typedef unsigned short		vuint16;
typedef int					vint32;
typedef unsigned int		vuint32;
#endif

// PUBLIC FUNCTION PROTOTYPES ----------------------------------------------

void Error(const char *error, ...) __attribute__ ((noreturn))
	__attribute__ ((format(printf, 1, 2)));

char *va(const char *text, ...) __attribute__ ((format(printf, 1, 2)));

short LittleShort(short val);
int LittleLong(int val);

void DefaultPath(char *path, const char *basepath);
void DefaultExtension(char *path, const char *extension);
void StripFilename(char *path);
void StripExtension(char *path);
void ExtractFilePath(const char *path, char *dest);
void ExtractFileBase(const char *path, char *dest);
void ExtractFileExtension(const char *path, char *dest);
void FixFileSlashes(char *path);
int LoadFile(const char *name, void **bufferptr);

// PUBLIC DATA DECLARATIONS ------------------------------------------------

extern void *(*Malloc)(size_t size);
extern void *(*Realloc)(void *data, size_t size);
extern void (*Free)(void *ptr);

} // namespace VavoomUtils

#endif
