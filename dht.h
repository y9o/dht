#ifndef DHT_H
#define DHT_H
#include <stdint.h>

typedef struct {
	int pre_sleep;
	int low_sleep;
	int high_sleep;
	int limit;
	int retry_delay;
	int vague;
} dht_config;
typedef struct dht_info dht_info;

dht_info * newDHT(const char *const chip,unsigned int line);
void  freeDHT(dht_info * dht);
int   readDHT(dht_info * dht, dht_config *dhtcfg, uint32_t *data);
int copyBufDHT(dht_info* dht, uint16_t *buf, int len);
long getReadTimeDHT(dht_info* dht, int readcount);

#endif