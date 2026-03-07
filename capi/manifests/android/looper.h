// Stub looper.h for c-for-go parser — provides forward declarations only.
#ifndef _ANDROID_LOOPER_H
#define _ANDROID_LOOPER_H

typedef struct ALooper ALooper;
typedef int (*ALooper_callbackFunc)(int fd, int events, void* data);

#define ALOOPER_POLL_WAKE 1
#define ALOOPER_POLL_CALLBACK 2
#define ALOOPER_POLL_TIMEOUT 3
#define ALOOPER_POLL_ERROR 4
#define ALOOPER_EVENT_INPUT 1
#define ALOOPER_EVENT_OUTPUT 2
#define ALOOPER_EVENT_ERROR 4
#define ALOOPER_EVENT_HANGUP 8
#define ALOOPER_EVENT_INVALID 16

ALooper* ALooper_forThread(void);
ALooper* ALooper_prepare(int opts);
void ALooper_acquire(ALooper* looper);
void ALooper_release(ALooper* looper);
int ALooper_pollOnce(int timeoutMillis, int* outFd, int* outEvents, void** outData);
int ALooper_pollAll(int timeoutMillis, int* outFd, int* outEvents, void** outData);
void ALooper_wake(ALooper* looper);
int ALooper_addFd(ALooper* looper, int fd, int ident, int events, ALooper_callbackFunc callback, void* data);
int ALooper_removeFd(ALooper* looper, int fd);

#endif
