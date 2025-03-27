import React from 'react';
import { ReviewableItem, QueryParams } from '../models/reviewable-item';

export interface TableColumn {
  title: string;
  dataIndex: string;
  key: string;
  render?: (value: any, record: any) => React.ReactNode;
  sorter?: boolean;
  filters?: { text: string; value: string }[];
  width?: number | string;
}

export interface ReviewStrategy<T extends ReviewableItem = ReviewableItem> {
  // 获取表格列定义
  getTableColumns(): TableColumn[];
  
  // 渲染列表项
  renderListItem(item: T): React.ReactNode;
  
  // 渲染详情视图
  renderDetailView(item: T): React.ReactNode;
  
  // 渲染过滤选项
  renderFilterOptions(filters: any, options: any, onChange: (filters: any) => void): React.ReactNode;
  
  // 格式化过滤参数为API请求格式
  formatFilterParams(filters: any): Partial<QueryParams>;
  
  // 判断是否可以处理该类型的数据
  canHandle(dataSource: string): boolean;
} 