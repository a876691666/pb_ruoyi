import type { BaseCollectionModel } from 'pocketbase';

export interface Spel extends BaseCollectionModel {
  id: number;
  componentName: string;
  methodName: string;
  methodParams: string;
  viewSpel: string;
  status: string;
  remark: string;
  createTime: string;
}
