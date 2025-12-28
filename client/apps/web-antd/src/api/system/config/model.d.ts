import type { BaseCollectionModel } from 'pocketbase';

export interface SysConfig extends BaseCollectionModel {
  id: number;
  name: string;
  key: string;
  value: string;
  type: string;
  remark: string;
  createTime: string;
}
