#include "dht.h"
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <time.h> 
#include <gpiod.h>

#define BUFFER_MAX 90

typedef struct {
	uint16_t * buffer;
	struct gpiod_line_config* inputConfig;
	struct gpiod_line_config *output0LowConfig;
	struct gpiod_line_config *output1HighConfig;
	int i;
	int edgeCount;
	enum gpiod_line_value firstValue;
} dht_buffer ;

struct dht_info {
	char *chip;
	unsigned int line;
	dht_buffer*buffer;
};

struct gpiod_line_config* lineConfig(
	unsigned int offset,
	enum gpiod_line_direction direction,
	enum gpiod_line_value value
){

	struct gpiod_line_settings *settings;
	struct gpiod_line_config *cfg;

	settings = gpiod_line_settings_new();
	if (settings==NULL){
		return NULL;
	}

	gpiod_line_settings_set_direction(settings, direction);
	gpiod_line_settings_set_output_value(settings, value);

	cfg = gpiod_line_config_new();
	if (cfg == NULL){
		gpiod_line_settings_free(settings);
		return NULL;
	}

	int r;
	r = gpiod_line_config_add_line_settings(cfg, &offset, 1, settings);
	gpiod_line_settings_free(settings);
	if (r != 0){
		gpiod_line_config_free(cfg);
		return NULL;
	}
	return cfg;
}


struct gpiod_line_request *newRequest(
	const char *chip_path,
	struct gpiod_line_config *cfg
){

	struct gpiod_line_request *request = NULL;
	struct gpiod_chip *chip;

	chip = gpiod_chip_open(chip_path);
	if (!chip){
		return NULL;
	}

	request = gpiod_chip_request_lines(chip, NULL, cfg);

	gpiod_chip_close(chip);

	return request;
}


int _array2data(const uint16_t*buffer, uint32_t*data,int bit){
	//最大値を抽出して0と1の閾値にする。
	int i=0;
	int border50or70usMax=1;
	int max=bit*2-1;//79=40*2-1
	for(i=0;i<max;i++){
		border50or70usMax=buffer[i]>border50or70usMax?buffer[i]:border50or70usMax;
	}
	uint32_t _data=0;
	for(i=0;i<(max-15);i+=2){//64=(max-15)
		_data <<= 1;
		_data |= buffer[i] *2 > border50or70usMax ? 1 : 0;
	}
	uint8_t checksum=0;
	for(i=(max-15);i<max;i+=2){//64=(max-15)
		checksum <<= 1;
		checksum |= buffer[i] *2 > border50or70usMax ? 1 : 0;
	}
	uint8_t sum = (( ((_data)&0xFF) + ((_data>>8)&0xFF) + ((_data>>16)&0xFF) + ((_data>>24)&0xFF)  )&0xFF);
	if (checksum == sum){
		*data=_data;
		return 1;
	}
	return 0;
}

int _array2data2(const uint16_t*buffer, uint32_t*data,int bit){
	//50usの平均値を作成して0と1の閾値にする。
	int i=0;
	int max=bit*2-1;//79=40*2-1
	//avg50us
	int border50us=0;
	int borderCount=0;
	for(i=1;i<max;i+=2){
		border50us += buffer[i];
		borderCount++;
	}
	if(borderCount>0){
		border50us /= borderCount;
	}
	// >=
	uint32_t _data=0;
	for(i=0;i<(max-15);i+=2){//64=(max-15)
		_data <<= 1;
		_data |= buffer[i] >= border50us ? 1 : 0;
	}
	uint8_t checksum=0;
	for(i=(max-15);i<max;i+=2){//64=(max-15)
		checksum <<= 1;
		checksum |= buffer[i] >= border50us ? 1 : 0;
	}
	uint8_t sum = (( ((_data)&0xFF) + ((_data>>8)&0xFF) + ((_data>>16)&0xFF) + ((_data>>24)&0xFF)  )&0xFF);
	if (checksum == sum){
		*data=_data;
		return 1;
	}
	// >
	_data=0;
	for(i=0;i<(max-15);i+=2){//64=(max-15)
		_data <<= 1;
		_data |= buffer[i] > border50us ? 1 : 0;
	}
	checksum=0;
	for(i=(max-15);i<max;i+=2){//64=(max-15)
		checksum <<= 1;
		checksum |= buffer[i] > border50us ? 1 : 0;
	}
	sum = (( ((_data)&0xFF) + ((_data>>8)&0xFF) + ((_data>>16)&0xFF) + ((_data>>24)&0xFF)  )&0xFF);
	if (checksum == sum){
		*data=_data;
		return 1;
	}

	return 0;
}

int array2data(dht_buffer *buffer,uint32_t*data,int bit){
	//二つの関数でデータ変換を試みる
	enum gpiod_line_value firstValue;
	int offset=buffer->edgeCount - (bit*2-1);//79
	int ret=-4;
	for(;offset>=0;offset--){
		firstValue = (offset&1)==0? buffer->firstValue : (buffer->firstValue^1) ;
		if(firstValue == GPIOD_LINE_VALUE_ACTIVE){
			ret=_array2data2(buffer->buffer +offset,data,bit);
			if(ret == 0){
				ret=_array2data(buffer->buffer+offset,data,bit);
			}
			break;
		}
	}
	return ret;
}

dht_info * newDHT(const char * chip,unsigned int line){
	dht_info *dht;
	dht_buffer*buffer;
	dht = malloc(sizeof(*dht));
	if (dht==NULL){
		return NULL;
	}
	memset(dht, 0, sizeof(*dht));
	buffer = malloc(sizeof(*buffer));
	if (buffer==NULL){
		free(dht);
		return NULL;
	}
	memset(buffer, 0, sizeof(*buffer));

	dht->chip = strdup(chip);
	dht->line = line;
	dht->buffer=buffer;
	buffer->buffer = malloc(sizeof(uint16_t)*BUFFER_MAX);
	buffer->inputConfig = lineConfig(line,GPIOD_LINE_DIRECTION_INPUT,GPIOD_LINE_VALUE_INACTIVE);
	buffer->output0LowConfig = lineConfig(line,GPIOD_LINE_DIRECTION_OUTPUT,GPIOD_LINE_VALUE_INACTIVE);
	buffer->output1HighConfig = lineConfig(line,GPIOD_LINE_DIRECTION_OUTPUT,GPIOD_LINE_VALUE_ACTIVE);
	if(dht->chip == NULL||buffer->buffer==NULL||buffer->inputConfig==NULL||buffer->output0LowConfig==NULL||buffer->output1HighConfig==NULL){
		freeDHT(dht);
		return NULL;
	}
	return dht;
}
void freeDHT(dht_info * dht){
	if(dht==NULL)
		return;
	dht_buffer* buffer;
	buffer=dht->buffer;

	if(buffer->buffer != NULL){
		free(buffer->buffer);
		buffer->buffer=NULL;
	}
	if(buffer->inputConfig != NULL){
		gpiod_line_config_free(buffer->inputConfig);
		buffer->inputConfig=NULL;
	}
	if(buffer->output0LowConfig != NULL){
		gpiod_line_config_free(buffer->output0LowConfig);
		buffer->output0LowConfig=NULL;
	}
	if(buffer->output1HighConfig != NULL){
		gpiod_line_config_free(buffer->output1HighConfig);
		buffer->output1HighConfig=NULL;
	}
	if(dht->chip!= NULL){
		free(dht->chip);
		dht->chip=NULL;
	}
	free(dht);
	return;
}
int readDHT(
	dht_info* dht,
	dht_config *dhtcfg,
	uint32_t *data
	){
	dht_buffer *buffer=dht->buffer;
	buffer->i=0;
	buffer->edgeCount=0;
	struct gpiod_line_request *request;
	enum gpiod_line_value value;
	enum gpiod_line_value lastvalue = GPIOD_LINE_VALUE_ERROR;
	int ret=0;
	int looptime5ms=dhtcfg->limit;//5ms
	int edgeCount=-1;
	int maxNegative160us=(looptime5ms/5)/(1000/160);
	int i=0;

	request = newRequest(dht->chip, buffer->inputConfig);
	if (!request) {
		return -1;
	}
	buffer->firstValue=-1;
	//pre
	if(dhtcfg->pre_sleep>0){
		usleep(dhtcfg->pre_sleep);
	}
	if( gpiod_line_request_reconfigure_lines(request, buffer->output0LowConfig) != 0){
		ret=-2;
		goto end_dht_read;
	}
	//low
	usleep(dhtcfg->low_sleep);

	if( gpiod_line_request_reconfigure_lines(request, buffer->output1HighConfig) != 0){
		ret=-3;
		goto end_dht_read;
	}
	//high
	if(dhtcfg->high_sleep>0){
		usleep(dhtcfg->high_sleep);
	}
	if( gpiod_line_request_reconfigure_lines(request, buffer->inputConfig) != 0){
		ret=-4;
		goto end_dht_read;
	}

	for(;i<looptime5ms;i++){
		value = gpiod_line_request_get_value(request, dht->line);
		if(value == GPIOD_LINE_VALUE_ERROR){
			break;
		}
		if(lastvalue != value){
			lastvalue = value;
			edgeCount++;
			if(edgeCount == BUFFER_MAX){
				ret=-5;
				goto end_dht_read;
			}
			buffer->buffer[edgeCount]=1;
			if(edgeCount==0){
				buffer->firstValue=value;
			}
		}else{
			buffer->buffer[edgeCount]++;
			if( buffer->buffer[edgeCount] > maxNegative160us ){
				edgeCount--;
				break;
			}
		}
	}
	buffer->i=i;
	buffer->edgeCount=edgeCount;
	
	uint32_t _data;
	for(i=0; i <= dhtcfg->vague ;i++){
		ret=array2data(buffer,&_data,40-i);
		if(ret>0){
			ret+=i;
			*data = _data;
			break;
		}
	}
end_dht_read:
	gpiod_line_request_release(request);
	return ret;
}

int copyBufDHT(dht_info* dht,uint16_t *buf,int len){
	int min = len > dht->buffer->edgeCount+1 ? dht->buffer->edgeCount+1 : len ;
	int i;
	for(i=0;i<min;i++){
		buf[i] = dht->buffer->buffer[i];
	}
	for(;i<len;i++){
		buf[i] = 0;
	}
	return dht->buffer->firstValue;
}
long getReadTimeDHT(dht_info* dht, int readcount){
	struct gpiod_line_request *request;
	enum gpiod_line_value value=GPIOD_LINE_VALUE_ERROR;
	struct timespec start, end,since;
	long ret=-1;
	request = newRequest(dht->chip, dht->buffer->inputConfig);
	if (!request) {
		return -1;
	}
	clock_gettime(CLOCK_MONOTONIC, &start);
	for(int i=0;i<readcount;i++){
		value = gpiod_line_request_get_value(request, dht->line);
		if(value == GPIOD_LINE_VALUE_ERROR){
			break;
		}
	}
	clock_gettime(CLOCK_MONOTONIC, &end);
	if(value == GPIOD_LINE_VALUE_ERROR){
		goto end_dht_read;
	}
	since.tv_sec = end.tv_sec - start.tv_sec;
	since.tv_nsec = end.tv_nsec - start.tv_nsec;
	if (since.tv_nsec < 0){
		since.tv_sec--;
		since.tv_nsec+=1000000000L;
	}
	while( since.tv_sec > 0 && since.tv_nsec < 1147483647L ){
		since.tv_sec--;
		since.tv_nsec+=1000000000L;
	}
	if (since.tv_sec > 0) {
		ret=-1;
	}else{
		ret = since.tv_nsec;
	}

end_dht_read:
	gpiod_line_request_release(request);
	return ret;
}