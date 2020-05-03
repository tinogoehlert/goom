//**************************************************************************
//**
//**	##   ##    ##    ##   ##   ####     ####   ###     ###
//**	##   ##  ##  ##  ##   ##  ##  ##   ##  ##  ####   ####
//**	 ## ##  ##    ##  ## ##  ##    ## ##    ## ## ## ## ##
//**	 ## ##  ########  ## ##  ##    ## ##    ## ##  ###  ##
//**	  ###   ##    ##   ###    ##  ##   ##  ##  ##       ##
//**	   #    ##    ##    #      ####     ####   ##       ##
//**
//**	$Id: wadlib.h 1599 2006-07-06 17:20:16Z dj_jl $
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

#ifndef WADLIB_H
#define WADLIB_H

namespace VavoomUtils {

// HEADER FILES ------------------------------------------------------------

// MACROS ------------------------------------------------------------------

// TYPES -------------------------------------------------------------------

struct lumpinfo_t
{
	char	name[12];
	int		position;
	int		size;
};

class WadLibError
{
 public:
	WadLibError(const char *Amessage)
	{
		strcpy(message, Amessage);
	}

	char message[256];
};

class TIWadFile
{
 public:
	TIWadFile()
	{
		buffer = NULL;
	}
	~TIWadFile()
	{
		
	}
	void Open(uint8_t* buffer);
	int LumpNumForName(const char* name);
	const char* LumpName(int lump)
	{
		return lump >= numlumps ? "" : lumpinfo[lump].name;
	}
	int LumpSize(int lump)
	{
		return lump >= numlumps ? 0 : lumpinfo[lump].size;
	}
	void* GetLump(int lump);
	void* GetLumpName(const char* name)
	{
		return GetLump(LumpNumForName(name));
	}
	void Close();

	uint8_t*		buffer;
	char			wadid[4];
	lumpinfo_t*		lumpinfo;
	int				numlumps;
};

class TOWadFile
{
 public:
	TOWadFile()
	{
		buffer = NULL;
	}
	~TOWadFile()
	{
		if (buffer)
		{
			Free(buffer);
		}
	}
	void Open(const char *Awadid);
	void AddLump(const char *name, const void *data, int size);
	void Close();

	uint8_t*		buffer;			
	uint32_t		size;
	char			wadid[4];
	lumpinfo_t*		lumpinfo;
	int				numlumps;
};

// PUBLIC FUNCTION PROTOTYPES ----------------------------------------------

// PUBLIC DATA DECLARATIONS ------------------------------------------------

//==========================================================================
//
//	CleanupName
//
//==========================================================================

inline void CleanupName(const char *src, char *dst)
{
	int i;
	for (i = 0; i < 8 && src[i]; i++)
	{
		dst[i] = toupper(src[i]);
	}
	for (; i < 12; i++)
	{
		dst[i] = 0;
	}
}

} // namespace VavoomUtils

#endif
