import type { BaseCollectionModel } from 'pocketbase';

export interface OssFile extends BaseCollectionModel {
  id: string;
  file_name: string;
  original_name: string;
  file_suffix: string;
  url: string;
  create_time: string;
  create_by: number;
  create_by_name: string;
}
