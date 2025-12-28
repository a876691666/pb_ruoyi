import type { BaseCollectionModel } from 'pocketbase';

export interface ProcessDefinition extends BaseCollectionModel {
  id: string;
  createTime: string;
  updateTime: string;
  tenantId: string;
  delFlag: string;
  flowCode: string;
  flowName: string;
  category: string;
  categoryName: string;
  version: string;
  isPublish: number;
  formCustom: string;
  formPath: string;
  activityStatus: number;
  listenerType?: any;
  listenerPath?: any;
  ext?: any;
}
