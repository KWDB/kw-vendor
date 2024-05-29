#include <geos_c.h>
#include <stdlib.h>
#include <stdio.h>
#include <stdarg.h>
#include <string.h>

enum FunctionLoad {
    Load_OK = 0,
    Load_ERROR_INIT,
    Load_ERROR_VERSION,
    Load_ERROR_BUFFER,
    Load_ERROR_GEOM_DESTROY,
    Load_ERROR_TOUCHES,
    Load_ERROR_INTERSECT,
    Load_ERROR_WITHIN,
    Load_ERROR_CONTAIN,
    Load_ERROR_EQUAL,
    Load_ERROR_COVER,
    Load_ERROR_AREA,
    Load_ERROR_DISTANCE,
    Load_ERROR_CREATE_READER,
    Load_ERROR_DESTORY_READER,
    Load_ERROR_READ_READER,
    Load_ERROR_CREATE_WRITER,
    Load_ERROR_DESTORY_WRITER,
    Load_ERROR_WRITE_WRITER,
};

void* initGEOSLibrary();
const char* loadGEOSVersion(enum FunctionLoad* status);

GEOSContextHandle_t call_fp_initGEOS_r_func(GEOSMessageHandler notice_function,
                                            GEOSMessageHandler error_function, enum FunctionLoad* status);
GEOSGeometry* call_fp_GEOSBuffer_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g,
                                        double width, int quadsegs, enum FunctionLoad* status);
void call_fp_GEOSGeom_destroy_r_func(GEOSContextHandle_t handle, GEOSGeometry* g, enum FunctionLoad* status);
char call_fp_GEOSTouches_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g1,
                                const GEOSGeometry* g2, enum FunctionLoad* status);
char call_fp_GEOSIntersects_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g1,
                                   const GEOSGeometry* g2, enum FunctionLoad* status);
char call_fp_GEOSWithin_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g1,
                               const GEOSGeometry* g2, enum FunctionLoad* status);
char call_fp_GEOSContains_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g1,
                                 const GEOSGeometry* g2, enum FunctionLoad* status);
char call_fp_GEOSEquals_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g1,
                               const GEOSGeometry* g2, enum FunctionLoad* status);
char call_fp_GEOSCovers_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g1,
                               const GEOSGeometry* g2, enum FunctionLoad* status);
int call_fp_GEOSArea_r_func(GEOSContextHandle_t handle,
                            const GEOSGeometry* g1, double *area, enum FunctionLoad* status);
int call_fp_GEOSDistance_r_func(GEOSContextHandle_t handle,
                                const GEOSGeometry* g1, const GEOSGeometry* g2, double *dist, enum FunctionLoad* status);
GEOSWKTReader* call_fp_createGEOSWKTReader(GEOSContextHandle_t context, enum FunctionLoad* status);
void call_fp_GEOSWKTReader_destroy_r_func(GEOSContextHandle_t handle,
                                          GEOSWKTReader* reader, enum FunctionLoad* status);
GEOSGeometry* call_fp_GEOSWKTReader_read_r(GEOSContextHandle_t handle, GEOSWKTReader* reader,
                                           const char *wkt, enum FunctionLoad* status);
GEOSWKTWriter* call_fp_GEOSWKTWriter_create_r_func(GEOSContextHandle_t handle, enum FunctionLoad* status);
void call_fp_GEOSWKTWriter_destroy_r_func(GEOSContextHandle_t handle, GEOSWKTWriter* writer, enum FunctionLoad* status);
char* call_fp_GEOSWKTWriter_write_r_func(GEOSContextHandle_t handle, GEOSWKTWriter* writer,
                                         const GEOSGeometry* g, enum FunctionLoad* status);

void closeGEOSLibrary();

void gogeos_notice_handler(const char *fmt, ...);
void gogeos_error_handler(const char *fmt, ...);
char *gogeos_get_last_error(void);
GEOSContextHandle_t gogeos_initGEOS(enum FunctionLoad* status);
