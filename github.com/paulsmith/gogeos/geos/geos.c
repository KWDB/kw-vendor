#include "geos.h"
#include <dlfcn.h>
#include <stdio.h>
#include <pthread.h>

static void* geos_library = NULL;
static pthread_mutex_t geosLibraryMutex = PTHREAD_MUTEX_INITIALIZER;

static GEOSContextHandle_t (*initGEOS_r_func)(GEOSMessageHandler ef, GEOSMessageHandler nf) = NULL;
static GEOSGeometry* (*GEOSBuffer_r_func)(GEOSContextHandle_t handle, const GEOSGeometry* g, double width, int quadsegs) = NULL;
static void (*GEOSGeom_destroy_r_func)(GEOSContextHandle_t handle, GEOSGeometry* g) = NULL;
static char (*GEOSTouches_r_func)(GEOSContextHandle_t handle, const GEOSGeometry* g1, const GEOSGeometry* g2) = NULL;
static char (*GEOSIntersects_r_func)(GEOSContextHandle_t handle, const GEOSGeometry* g1, const GEOSGeometry* g2) = NULL;
static char (*GEOSWithin_r_func)(GEOSContextHandle_t handle, const GEOSGeometry* g1, const GEOSGeometry* g2) = NULL;
static char (*GEOSContains_r_func)(GEOSContextHandle_t handle, const GEOSGeometry* g1, const GEOSGeometry* g2) = NULL;
static char (*GEOSEquals_r_func)(GEOSContextHandle_t handle, const GEOSGeometry* g1, const GEOSGeometry* g2) = NULL;
static char (*GEOSCovers_r_func)(GEOSContextHandle_t handle, const GEOSGeometry* g1, const GEOSGeometry* g2) = NULL;
static int (*GEOSArea_r_func)(GEOSContextHandle_t handle, const GEOSGeometry* g1, double *area) = NULL;
static int (*GEOSDistance_r_func)(GEOSContextHandle_t handle, const GEOSGeometry* g1, const GEOSGeometry* g2, double *dist) = NULL;
static GEOSWKTReader* (*GEOSWKTReader_create_r_func)(GEOSContextHandle_t) = NULL;
static void (*GEOSWKTReader_destroy_r_func)(GEOSContextHandle_t handle, GEOSWKTReader* reader) = NULL;
static GEOSGeometry*  (*GEOSWKTReader_read_r_func)(GEOSContextHandle_t handle, GEOSWKTReader* reader, const char *wkt) = NULL;
static GEOSWKTWriter* (*GEOSWKTWriter_create_r_func)(GEOSContextHandle_t handle) = NULL;
static void (*GEOSWKTWriter_destroy_r_func)(GEOSContextHandle_t handle, GEOSWKTWriter* writer) = NULL;
static char* (*GEOSWKTWriter_write_r_func)(GEOSContextHandle_t handle, GEOSWKTWriter* writer, const GEOSGeometry* g) = NULL;

// 初始化 libgeos_c.so
void* initGEOSLibrary() {
    pthread_mutex_lock(&geosLibraryMutex);
    geos_library = dlopen("libgeos_c.so", RTLD_NOW | RTLD_GLOBAL);
    if (geos_library == NULL) {
        // 如果加载 libgeos_c.so 失败，尝试加载 libgeos.so
        geos_library = dlopen("libgeos.so", RTLD_NOW | RTLD_GLOBAL);
    }

    initGEOS_r_func = dlsym(geos_library, "initGEOS_r");
    GEOSBuffer_r_func = dlsym(geos_library, "GEOSBuffer_r");
    GEOSGeom_destroy_r_func = dlsym(geos_library, "GEOSGeom_destroy_r");
    GEOSTouches_r_func = dlsym(geos_library, "GEOSTouches_r");
    GEOSIntersects_r_func = dlsym(geos_library, "GEOSIntersects_r");
    GEOSWithin_r_func = dlsym(geos_library, "GEOSWithin_r");
    GEOSContains_r_func = dlsym(geos_library, "GEOSContains_r");
    GEOSEquals_r_func = dlsym(geos_library, "GEOSEquals_r");
    GEOSCovers_r_func = dlsym(geos_library, "GEOSCovers_r");
    GEOSArea_r_func = dlsym(geos_library, "GEOSArea_r");
    GEOSDistance_r_func = dlsym(geos_library, "GEOSDistance_r");
    GEOSWKTReader_create_r_func = dlsym(geos_library, "GEOSWKTReader_create_r");
    GEOSWKTReader_destroy_r_func = dlsym(geos_library, "GEOSWKTReader_destroy_r");
    GEOSWKTReader_read_r_func = dlsym(geos_library, "GEOSWKTReader_read_r");
    GEOSWKTWriter_create_r_func = dlsym(geos_library, "GEOSWKTWriter_create_r");
    GEOSWKTWriter_destroy_r_func = dlsym(geos_library, "GEOSWKTWriter_destroy_r");
    GEOSWKTWriter_write_r_func = dlsym(geos_library, "GEOSWKTWriter_write_r");

    pthread_mutex_unlock(&geosLibraryMutex);
    return geos_library;
}

// 获取 GEOS 版本
const char* loadGEOSVersion(enum FunctionLoad* status) {
    typedef const char* (*GEOSVersionFunc)();
    GEOSVersionFunc versionFunc = (GEOSVersionFunc)dlsym(geos_library, "GEOSversion");
    if (versionFunc == NULL) {
        *status = Load_ERROR_VERSION;
        return '\0';
    }
    *status = Load_OK;
    return versionFunc();
}

GEOSContextHandle_t call_fp_initGEOS_r_func(GEOSMessageHandler notice_function, GEOSMessageHandler error_function, enum FunctionLoad* status)
{
    if (initGEOS_r_func == NULL) {
        *status = Load_ERROR_INIT;
        return NULL;
    }
    *status = Load_OK;
    return initGEOS_r_func(notice_function, error_function);
}

GEOSGeometry* call_fp_GEOSBuffer_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g, double width, int quadsegs, enum FunctionLoad* status)
{
    if (GEOSBuffer_r_func == NULL) {
        *status = Load_ERROR_BUFFER;
        return NULL;
    }
    *status = Load_OK;
    return GEOSBuffer_r_func(handle, g, width, quadsegs);
}

void call_fp_GEOSGeom_destroy_r_func(GEOSContextHandle_t handle, GEOSGeometry* g, enum FunctionLoad* status)
{
    if (GEOSGeom_destroy_r_func == NULL) {
        *status = Load_ERROR_GEOM_DESTROY;
        return;
    }
    *status = Load_OK;
    return GEOSGeom_destroy_r_func(handle, g);
}

char call_fp_GEOSTouches_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g1, const GEOSGeometry* g2, enum FunctionLoad* status)
{
    if (GEOSTouches_r_func == NULL) {
        *status = Load_ERROR_TOUCHES;
        return '\0';
    }
    *status = Load_OK;
    return GEOSTouches_r_func(handle, g1, g2);
}

char call_fp_GEOSIntersects_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g1, const GEOSGeometry* g2, enum FunctionLoad* status)
{
    if (GEOSIntersects_r_func == NULL) {
        *status = Load_ERROR_INTERSECT;
        return '\0';
    }
    *status = Load_OK;
    return GEOSIntersects_r_func(handle, g1, g2);
}

char call_fp_GEOSWithin_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g1, const GEOSGeometry* g2, enum FunctionLoad* status)
{
    if (GEOSWithin_r_func == NULL) {
        *status = Load_ERROR_WITHIN;
        return '\0';
    }
    *status = Load_OK;
    return GEOSWithin_r_func(handle, g1, g2);
}

char call_fp_GEOSContains_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g1, const GEOSGeometry* g2, enum FunctionLoad* status)
{
    if (GEOSContains_r_func == NULL) {
        *status = Load_ERROR_CONTAIN;
        return '\0';
    }
    *status = Load_OK;
    return GEOSContains_r_func(handle, g1, g2);
}

char call_fp_GEOSEquals_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g1, const GEOSGeometry* g2, enum FunctionLoad* status)
{
    if (GEOSEquals_r_func == NULL) {
        *status = Load_ERROR_EQUAL;
        return '\0';
    }
    *status = Load_OK;
    return GEOSEquals_r_func(handle, g1, g2);
}

char call_fp_GEOSCovers_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g1, const GEOSGeometry* g2, enum FunctionLoad* status)
{
    if (GEOSCovers_r_func == NULL) {
        *status = Load_ERROR_COVER;
        return '\0';
    }
    *status = Load_OK;
    return GEOSCovers_r_func(handle, g1, g2);
}

int call_fp_GEOSArea_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g1, double *area, enum FunctionLoad* status)
{
    if (GEOSArea_r_func == NULL) {
        *status = Load_ERROR_AREA;
        return 0;
    }
    *status = Load_OK;
    return GEOSArea_r_func(handle, g1, area);
}

int call_fp_GEOSDistance_r_func(GEOSContextHandle_t handle, const GEOSGeometry* g1, const GEOSGeometry* g2, double *dist, enum FunctionLoad* status)
{
    if (GEOSDistance_r_func == NULL) {
        *status = Load_ERROR_DISTANCE;
        return 0;
    }
    *status = Load_OK;
    return GEOSDistance_r_func(handle, g1, g2, dist);
}

GEOSWKTReader* call_fp_createGEOSWKTReader(GEOSContextHandle_t context, enum FunctionLoad* status)
{
    if (GEOSWKTReader_create_r_func == NULL) {
        *status = Load_ERROR_CREATE_READER;
        return NULL;
    }
    *status = Load_OK;
    return GEOSWKTReader_create_r_func(context);
}

void call_fp_GEOSWKTReader_destroy_r_func(GEOSContextHandle_t handle, GEOSWKTReader* reader, enum FunctionLoad* status)
{
    if (GEOSWKTReader_destroy_r_func == NULL) {
        *status = Load_ERROR_DESTORY_READER;
        return;
    }
    *status = Load_OK;
    return GEOSWKTReader_destroy_r_func(handle, reader);
}

GEOSGeometry* call_fp_GEOSWKTReader_read_r(GEOSContextHandle_t handle, GEOSWKTReader* reader, const char *wkt, enum FunctionLoad* status)
{
    if (GEOSWKTReader_read_r_func == NULL) {
        *status = Load_ERROR_READ_READER;
        return NULL;
    }
    *status = Load_OK;
    return GEOSWKTReader_read_r_func(handle, reader, wkt);
}

GEOSWKTWriter* call_fp_GEOSWKTWriter_create_r_func(GEOSContextHandle_t handle, enum FunctionLoad* status)
{
    if (GEOSWKTWriter_create_r_func == NULL) {
        *status = Load_ERROR_CREATE_WRITER;
        return NULL;
    }
    *status = Load_OK;
    return GEOSWKTWriter_create_r_func(handle);
}

void call_fp_GEOSWKTWriter_destroy_r_func(GEOSContextHandle_t handle, GEOSWKTWriter* writer, enum FunctionLoad* status)
{
    if (GEOSWKTWriter_destroy_r_func == NULL) {
        *status = Load_ERROR_DESTORY_WRITER;
        return;
    }
    *status = Load_OK;
    return GEOSWKTWriter_destroy_r_func(handle, writer);
}

char* call_fp_GEOSWKTWriter_write_r_func(GEOSContextHandle_t handle, GEOSWKTWriter* writer, const GEOSGeometry* g, enum FunctionLoad* status)
{
    if (GEOSWKTWriter_write_r_func == NULL) {
        *status = Load_ERROR_WRITE_WRITER;
        return '\0';
    }
    *status = Load_OK;
    return GEOSWKTWriter_write_r_func(handle, writer, g);
}

void closeGEOSLibrary() {
    if (geos_library != NULL) {
        dlclose(geos_library);
        geos_library = NULL;
    }
}

void gogeos_notice_handler(const char *fmt, ...) {
    va_list ap;
    va_start(ap, fmt);
    fprintf(stderr, "NOTICE: ");
    vfprintf(stderr, fmt, ap);
    va_end(ap);
}

#define ERRLEN 256

char gogeos_last_err[ERRLEN];

void gogeos_error_handler(const char *fmt, ...) {
    va_list ap;
    va_start(ap, fmt);
    vsnprintf(gogeos_last_err, (size_t) ERRLEN, fmt, ap);
    va_end(ap);
}

char *gogeos_get_last_error(void) {
    return gogeos_last_err;
}

GEOSContextHandle_t gogeos_initGEOS(enum FunctionLoad* status) {
    return call_fp_initGEOS_r_func(gogeos_notice_handler, gogeos_error_handler, status);
}
