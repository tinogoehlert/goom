#ifdef __cplusplus
extern "C" {
#endif

#ifndef uint8_t 
  typedef unsigned char uint8_t; 
  typedef unsigned int uint32_t; 
#endif

  void* GLVisInit(void);
  void GLVisFree(void*);
  void BuildVis(void* f, uint8_t* wad, uint8_t *gwa);
  uint32_t GetVisSize(void* f);
  uint8_t* GetVisData(void* f);
#ifdef __cplusplus
}
#endif
