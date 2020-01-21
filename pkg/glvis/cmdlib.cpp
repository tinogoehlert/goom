//**************************************************************************
//**
//**	##   ##    ##    ##   ##   ####     ####   ###     ###
//**	##   ##  ##  ##  ##   ##  ##  ##   ##  ##  ####   ####
//**	 ## ##  ##    ##  ## ##  ##    ## ##    ## ## ## ## ##
//**	 ## ##  ########  ## ##  ##    ## ##    ## ##  ###  ##
//**	  ###   ##    ##   ###    ##  ##   ##  ##  ##       ##
//**	   #    ##    ##    #      ####     ####   ##       ##
//**
//**	$Id: cmdlib.cpp 1599 2006-07-06 17:20:16Z dj_jl $
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

// HEADER FILES ------------------------------------------------------------

#include "cmdlib.h"

namespace VavoomUtils {

// MACROS ------------------------------------------------------------------

// TYPES -------------------------------------------------------------------

// EXTERNAL FUNCTION PROTOTYPES --------------------------------------------

// PUBLIC FUNCTION PROTOTYPES ----------------------------------------------

// PRIVATE FUNCTION PROTOTYPES ---------------------------------------------

static void *DefaultMalloc(size_t size);
static void *DefaultRealloc(void *data, size_t size);
static void DefaultFree(void *ptr);

// EXTERNAL DATA DECLARATIONS ----------------------------------------------

// PUBLIC DATA DEFINITIONS -------------------------------------------------

void *(*Malloc)(size_t size) = DefaultMalloc;
void *(*Realloc)(void *data, size_t size) = DefaultRealloc;
void (*Free)(void *ptr) = DefaultFree;

// PRIVATE DATA DEFINITIONS ------------------------------------------------

static int		empty_ptr;

// CODE --------------------------------------------------------------------

//==========================================================================
//
//	DefaultMalloc
//
//==========================================================================

static void *DefaultMalloc(size_t size)
{
	if (!size)
	{
		return &empty_ptr;
	}

	void *ptr = malloc(size);
	if (!ptr)
	{
		Error("Couldn't alloc %d bytes", (int)size);
	}
	memset(ptr, 0, size);
	return ptr;
}

//==========================================================================
//
//	DefaultRealloc
//
//==========================================================================

static void *DefaultRealloc(void *data, size_t size)
{
	void *ptr = realloc(data, size);
	if (!ptr)
	{
		Error("Couldn't realloc %d bytes", (int)size);
	}
	return ptr;
}

//==========================================================================
//
//	DefaultFree
//
//==========================================================================

static void DefaultFree(void *ptr)
{
	if (ptr != &empty_ptr)
	{
		free(ptr);
	}
}

//==========================================================================
//
//	Error
//
//==========================================================================

void Error(const char *error, ...)
{
	va_list		argptr;

	va_start(argptr, error);
	vfprintf(stderr, error, argptr);
	va_end(argptr);

	exit(1);
}

//==========================================================================
//
//	va
//
//	Very usefull function from QUAKE
//
//==========================================================================

char *va(const char *text, ...)
{
	va_list		args;
	static char	va_buffer[1024];

	va_start(args, text);
	vsprintf(va_buffer, text, args);
	va_end(args);

	return va_buffer;
}

//==========================================================================
//
//	LittleShort
//
//==========================================================================

short LittleShort(short val)
{
	vuint8* bv = (vuint8*)&val;
	return bv[0] + (bv[1] << 8);
}

//==========================================================================
//
//	LittleLong
//
//==========================================================================

int LittleLong(int val)
{
	vuint8* bv = (vuint8*)&val;
	return bv[0] + (bv[1] << 8) + (bv[2] << 16) + (bv[3] << 24);
}

//==========================================================================
//
//	DefaultPath
//
//==========================================================================

void DefaultPath(char *path, const char *basepath)
{
	char	temp[128];

	if (path[0] == '/')
		return;					// absolute path location
	strcpy(temp, path);
	strcpy(path, basepath);
	strcat(path, temp);
}

//==========================================================================
//
//	DefaultExtension
//
//==========================================================================

void DefaultExtension(char *path, const char *extension)
{
	char	*src;

	//
	// if path doesn't have a .EXT, append extension
	// (extension should include the .)
	//
	src = path + strlen(path) - 1;

	while (*src != '/' && src != path)
	{
		if (*src == '.')
		{
			return;			// it has an extension
		}
		src--;
	}

	strcat(path, extension);
}

//==========================================================================
//
//	StripFilename
//
//==========================================================================

void StripFilename(char *path)
{
	int			 length;

	length = strlen(path)-1;
	while (length > 0 && path[length] != '/')
		length--;
	path[length] = 0;
}

//==========================================================================
//
//	StripExtension
//
//==========================================================================

void StripExtension(char *path)
{
	char	*search;

	search = path + strlen(path) - 1;
	while (*search != '/' && search != path)
	{
		if (*search == '.')
		{
			*search = 0;
			return;
		}
		search--;
	}
}

//==========================================================================
//
//	ExtractFilePath
//
//==========================================================================

void ExtractFilePath(const char *path, char *dest)
{
	const char	*src;

	src = path + strlen(path) - 1;

	//
	// back up until a \ or the start
	//
	while (src != path && *(src-1) != '/')
		src--;

	memcpy(dest, path, src - path);
	dest[src - path] = 0;
}

//==========================================================================
//
//	ExtractFileBase
//
//==========================================================================

void ExtractFileBase(const char *path, char *dest)
{
	const char	*src;

	src = path + strlen(path) - 1;

	//
	// back up until a \ or the start
	//
	while (src != path && *(src-1) != '/')
		src--;

	while (*src && *src != '.')
	{
		*dest++ = *src++;
	}
	*dest = 0;
}

//==========================================================================
//
//	ExtractFileExtension
//
//==========================================================================

void ExtractFileExtension(const char *path, char *dest)
{
	const char	*src;

	src = path + strlen(path) - 1;

	//
	// back up until a . or the start
	//
	while (src != path && *(src-1) != '.')
		src--;
	if (src == path)
	{
		*dest = 0;		// no extension
		return;
	}

	strcpy(dest, src);
}

//==========================================================================
//
//	FixFileSlashes
//
//==========================================================================

void FixFileSlashes(char *path)
{
	while (*path)
	{
		if (*path == '\\')
			*path = '/';
		path++;
	}
}

//==========================================================================
//
//	LoadFile
//
//==========================================================================

int LoadFile(const char *name, void **bufferptr)
{
	FILE		*f;
	int			length;
	int			count;
	void		*buffer;

	f = fopen(name, "rb");
	if (!f)
	{
		Error("Couldn't open file \"%s\".", name);
	}

	fseek(f, 0, SEEK_END);
	length = ftell(f);
	fseek(f, 0, SEEK_SET);

	buffer = Malloc(length);

	count = fread(buffer, 1, length, f);
	fclose (f);

	if (count != length)
	{
		Error("Couldn't read file \"%s\".", name);
	}

	*bufferptr = buffer;
	return length;
}

} // namespace VavoomUtils
