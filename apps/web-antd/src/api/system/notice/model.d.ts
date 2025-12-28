import type { BaseCollectionModel } from 'pocketbase';

export interface Notice extends BaseCollectionModel {
  id: number;
  notice_title: string;
  notice_type: string;
  notice_content: string;
  status: string;
  remark: string;
  create_by: number;
  create_by_name: string;
  create_time: string;
}
