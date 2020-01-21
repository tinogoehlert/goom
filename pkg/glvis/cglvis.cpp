#include "glvis.h"
#include "glvis.hpp"
#include "glvisint.h"
#include <stdarg.h>

using namespace VavoomUtils;

class TCGOGLVis:public TGLVis
{
 public:
  TCGOGLVis();
	void DisplayMessage(const char *text, ...)
		__attribute__((format(printf, 2, 3)));
	void DisplayStartMap(const char *levelname);
	void DisplayBaseVisProgress(int count, int total);
	void DisplayPortalVisProgress(int count, int total);
	void DisplayMapDone(int accepts, int total);
  void RunBuilder(uint8_t* wadBuff, uint8_t* gwaBuff);
  TOWadFile* GetOutput();
 private:
  TVisBuilder* visBuilder;
};

TCGOGLVis::TCGOGLVis() {
  this->visBuilder = new TVisBuilder(*this);
}

void TCGOGLVis::DisplayMessage(const char *text, ...)
{
    va_list args;
    va_start(args, text);
		vfprintf(stderr, text, args);
		va_end(args);
}

void TCGOGLVis::DisplayStartMap(const char *)
{
}

void TCGOGLVis::DisplayBaseVisProgress(int count, int)
{
}

void TCGOGLVis::DisplayPortalVisProgress(int count, int total)
{
}

void TCGOGLVis::DisplayMapDone(int accepts, int total)
{
}

void TCGOGLVis::RunBuilder(uint8_t* wadBuff, uint8_t* gwaBuff) {
  this->visBuilder->Run(wadBuff,gwaBuff);
}

TOWadFile* TCGOGLVis::GetOutput() {
  return this->visBuilder->GetOutput();
}

void* GLVisInit()
{
  TCGOGLVis *ret = new TCGOGLVis();
  return (void*)ret;
}
void GLVisFree(void* f)
{
  TCGOGLVis *glvis = (TCGOGLVis*)f;
  delete glvis;
}

void BuildVis(void* f, uint8_t* wad, uint8_t *gwa)
{
  TCGOGLVis *glvis = (TCGOGLVis*)f;
  try {
    glvis->fastvis = true;
    glvis->RunBuilder(wad, gwa);
  } catch(GLVisError err) {
    printf("ERR: %s", err.message);
  }
}

uint32_t GetVisSize(void* f) {
  TCGOGLVis *glvis = (TCGOGLVis*)f;
  TOWadFile *out = glvis->GetOutput();
  return glvis->GetOutput()->size;
}

uint8_t* GetVisData(void* f) {
  TCGOGLVis *glvis = (TCGOGLVis*)f;
  return (uint8_t*)glvis->GetOutput()->buffer;
}
