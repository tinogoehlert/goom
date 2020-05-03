//**************************************************************************
//**
//**	##   ##    ##    ##   ##   ####     ####   ###     ###
//**	##   ##  ##  ##  ##   ##  ##  ##   ##  ##  ####   ####
//**	 ## ##  ##    ##  ## ##  ##    ## ##    ## ## ## ## ##
//**	 ## ##  ########  ## ##  ##    ## ##    ## ##  ###  ##
//**	  ###   ##    ##   ###    ##  ##   ##  ##  ##       ##
//**	   #    ##    ##    #      ####     ####   ##       ##
//**
//**	$Id: wadlib.cpp 1599 2006-07-06 17:20:16Z dj_jl $
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
#include "wadlib.h"

namespace VavoomUtils {

#include "fwaddefs.h"

void TIWadFile::Open(uint8_t* buff)
{
	wadinfo_t		header;
	lumpinfo_t*		lump_p;
	int				i;
	filelump_t*		fileinfo;
	filelump_t*		fi_p;

	this->buffer = buff;
	memcpy(&header, buff, sizeof(header));
	
	if (strncmp(header.identification, "IWAD", 4))
	{
		// Homebrew levels?
		if (strncmp(header.identification, "PWAD", 4))
		{
			throw WadLibError(va("buffer doesn't have IWAD or PWAD id"));
		}
	}
	strcpy(wadid, header.identification);
	numlumps = LittleLong(header.numlumps);
	this->lumpinfo = new lumpinfo_t[numlumps];
	fileinfo = new filelump_t[numlumps];
	
	memcpy(fileinfo, this->buffer+header.infotableofs, numlumps * sizeof(filelump_t));

	// Fill in lumpinfo
	lump_p = this->lumpinfo;
	fi_p = fileinfo;
	for (i = 0; i < numlumps; i++, lump_p++, fi_p++)
	{
		CleanupName(fi_p->name, lump_p->name);
		lump_p->position = LittleLong(fi_p->filepos);
		lump_p->size = LittleLong(fi_p->size);
	}
	
	delete fileinfo;
}

//==========================================================================
//
//	TIWadFile::LumpNumForName
//
//==========================================================================

int TIWadFile::LumpNumForName(const char* name)
{
	int			i;
	char		buf[12];

	CleanupName(name, buf);
	for (i = numlumps - 1; i >= 0; i--)
	{
		if (!strcmp(buf, lumpinfo[i].name))
			return i;
	}
	throw WadLibError(va("W_GetNumForName: %s not found!", name));
}

//==========================================================================
//
//	TIWadFile::GetLump
//
//==========================================================================

void* TIWadFile::GetLump(int lump)
{
	void*		ptr;
	lumpinfo_t*	l;
	
	l = &lumpinfo[lump];
	ptr = Malloc(l->size);
	memcpy(ptr, buffer+(l->position), l->size);
	return ptr;
}

//==========================================================================
//
//	TIWadFile::Close
//
//==========================================================================

void TIWadFile::Close()
{
	buffer = NULL;
	delete lumpinfo;
}

//==========================================================================
//
//	TOWadFile::Open
//
//==========================================================================

void TOWadFile::Open(const char *Awadid)
{
	wadinfo_t	header;
	memset(&header, 0, sizeof(header));
	buffer = (uint8_t*)Malloc(sizeof(header));
	memcpy(buffer, &header, sizeof(header));
	this->size = sizeof(header);
	lumpinfo = new lumpinfo_t[8 * 1024];
	strncpy(wadid, Awadid, 4);
	numlumps = 0;
}

//==========================================================================
//
//	TOWadFile::AddLump
//
//==========================================================================

void TOWadFile::AddLump(const char *name, const void *data, int size)
{
	CleanupName(name, lumpinfo[numlumps].name);
	lumpinfo[numlumps].size = size;
	lumpinfo[numlumps].position = this->size;
	if (size)
	{
		buffer = (uint8_t*)Realloc(buffer, this->size + size);
		memcpy(buffer+this->size, data, size);
		this->size += size;
	}
	numlumps++;
}

//==========================================================================
//
//	TOWadFile::Close
//
//==========================================================================

void TOWadFile::Close()
{
	wadinfo_t	header;
	strcpy(header.identification, wadid);
	header.numlumps = LittleLong(numlumps);
	header.infotableofs = this->size;
	for (int i = 0; i < numlumps; i++)
	{
		filelump_t	fileinfo;
		strncpy(fileinfo.name, lumpinfo[i].name, 8);
		fileinfo.size = LittleLong(lumpinfo[i].size);
		fileinfo.filepos = LittleLong(lumpinfo[i].position);

		buffer = (uint8_t*)Realloc(buffer, this->size + sizeof(fileinfo));
		memcpy(buffer+this->size, &fileinfo, sizeof(fileinfo));
		this->size += sizeof(fileinfo);
	}
	memcpy(buffer, &header, sizeof(header));
	delete lumpinfo;
}

} // namespace VavoomUtils
