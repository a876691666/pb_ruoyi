import type { BaseCollectionModel } from 'pocketbase';

export interface Dept extends BaseCollectionModel {
  create_by: string;
  id: string;
  create_time: string;
  update_by?: string;
  update_time?: string;
  remark?: string;
  id: number;
  parent_id: number;
  ancestors: string;
  dept_name: string;
  order_num: number;
  leader: string;
  phone: string;
  email: string;
  status: string;
  del_flag: string;
  parent_name?: string;
  children?: Dept[];
}
