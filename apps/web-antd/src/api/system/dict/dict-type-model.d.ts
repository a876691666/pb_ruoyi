import type { BaseCollectionModel } from 'pocketbase';

export interface DictType extends BaseCollectionModel {
  createTime: string;
  id: string;
  dict_name: string;
  dict_type: string;
  remark: string;
  status: string;
}
