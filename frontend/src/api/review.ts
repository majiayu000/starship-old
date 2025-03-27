import axios from 'axios';
import { PaginatedResult, QueryParams, ReviewStatusUpdate } from '../models/reviewable-item';
import { DataSourceItemType, FilterOptions } from './types';

// 从环境变量读取API地址，如果未设置则使用默认值（开发环境中指向8080端口）
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const reviewApi = {
  // 获取所有数据源
  getDataSources: async (): Promise<string[]> => {
    try {
      const response = await axios.get(`${API_BASE_URL}/review/sources`);
      return response.data;
    } catch (error) {
      console.error('获取数据源失败:', error);
      // 返回模拟数据，以防API不可用
      return ['sat_oneprep', 'sat_ixl'];
    }
  },

  // 获取过滤选项
  getFilterOptions: async (dataSource: string): Promise<FilterOptions> => {
    try {
      const response = await axios.get(`${API_BASE_URL}/review/${dataSource}/filters`);
      return response.data;
    } catch (error) {
      console.error('获取过滤选项失败:', error);
      // 返回默认的过滤选项
      return {
        reviewStatus: ['pending', 'approved', 'rejected'],
        domain: [],
        skill: []
      };
    }
  },

  // 获取项目列表（带分页和过滤）
  getItems: async <T extends keyof DataSourceItemType>(
    dataSource: T, 
    params: QueryParams
  ): Promise<PaginatedResult<DataSourceItemType[T]>> => {
    try {
      const response = await axios.get(`${API_BASE_URL}/review/${dataSource}/items`, { params });
      return response.data;
    } catch (error) {
      console.error('获取项目列表失败:', error);
      // 返回空列表作为模拟数据
      return {
        items: [],
        totalItems: 0,
        totalPages: 0,
        page: params.page || 1,
        pageSize: params.pageSize || 10,
        dataSource: dataSource
      };
    }
  },

  // 获取单个项目
  getItemById: async <T extends keyof DataSourceItemType>(
    dataSource: T, 
    id: string
  ): Promise<DataSourceItemType[T]> => {
    try {
      const response = await axios.get(`${API_BASE_URL}/review/${dataSource}/items/${id}`);
      return response.data;
    } catch (error) {
      console.error('获取项目详情失败:', error);
      throw new Error('获取项目详情失败');
    }
  },

  // 提交审核
  reviewItem: async (
    dataSource: string, 
    id: string, 
    update: ReviewStatusUpdate
  ): Promise<void> => {
    try {
      // 只发送状态字段
      const formattedUpdate = {
        status: update.status
      };
      
      console.log('提交审核请求:', {
        url: `${API_BASE_URL}/review/${dataSource}/items/${id}/review`,
        data: formattedUpdate
      });
      
      await axios.post(
        `${API_BASE_URL}/review/${dataSource}/items/${id}/review`, 
        formattedUpdate,
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      );
    } catch (error: any) {
      console.error('提交审核失败:', error);
      if (error.response) {
        console.error('响应数据:', error.response.data);
        console.error('响应状态:', error.response.status);
      }
      throw new Error(error.response?.data?.error || '提交审核失败');
    }
  }
}; 