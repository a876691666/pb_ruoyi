import type { BaseCollectionModel } from 'pocketbase';

export interface DictData extends BaseCollectionModel {
  create_by: string;
  create_time: string;
  css_class: string;
  default: boolean;
  id: string;
  dict_label: string;
  dict_sort: number;
  dict_type: string;
  dict_value: string;
  is_default: string;
  list_class: string;
  remark: string;
  status: string;
  update_by?: any;
  update_time?: any;
}
