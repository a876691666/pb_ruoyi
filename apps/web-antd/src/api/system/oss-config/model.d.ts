import type { BaseCollectionModel } from 'pocketbase';

export interface OssConfig extends BaseCollectionModel {
  ossConfigId: number;
  configKey: string;
  accessKey: string;
  secretKey: string;
  bucketName: string;
  prefix: string;
  endpoint: string;
  domain: string;
  isHttps: string;
  region: string;
  status: string;
  ext1: string;
  remark: string;
  accessPolicy: string;
}
