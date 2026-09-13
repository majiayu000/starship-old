'use client';

import React, { useState, useEffect } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Skeleton } from '@/components/ui/skeleton';
import DataSourceSelector from '@/components/common/DataSourceSelector';
import { StrategyFactory } from '@/strategies/factory';
import { reviewApi } from '@/api/review';
import { DataSourceItemType } from '@/api/types';
import { ChevronDown, ChevronUp } from 'lucide-react';
import { toast } from 'sonner';
import { useAuthSession } from '@/hooks/useAuthSession';

export default function QuestionsPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { authEpoch, isAuthenticated } = useAuthSession();
  const [dataSource, setDataSource] = useState<keyof DataSourceItemType>('sat_oneprep');
  const [availableSources, setAvailableSources] = useState<string[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [items, setItems] = useState<any[]>([]);
  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);
  const [totalItems, setTotalItems] = useState<number>(0);
  const [totalPages, setTotalPages] = useState<number>(1);
  const [filterOptions, setFilterOptions] = useState<any>({});
  // 默认过滤器设置为只显示待审核的题目
  const [filters, setFilters] = useState<any>({ status: 'pending' });
  // 为元数据折叠状态创建一个映射，以originalId为键
  const [metadataOpenMap, setMetadataOpenMap] = useState<Record<string, boolean>>({});
  
  // 处理元数据折叠/展开
  const toggleMetadata = (originalId: string) => {
    setMetadataOpenMap(prev => ({
      ...prev,
      [originalId]: !prev[originalId]
    }));
  };
  
  // 获取数据源
  useEffect(() => {
    if (!isAuthenticated) {
      setAvailableSources([]);
      setItems([]);
      setTotalItems(0);
      setTotalPages(1);
      setFilterOptions({});
      setError(null);
      setLoading(false);
      return;
    }

    const fetchDataSources = async () => {
      try {
        setLoading(true);
        setError(null);
        const sources = await reviewApi.getDataSources();
        setAvailableSources(sources);
        
        const dataSourceParam = searchParams.get('dataSource');
        
        // 设置数据源
        if (dataSourceParam && StrategyFactory.isDataSourceSupported(dataSourceParam)) {
          setDataSource(dataSourceParam as keyof DataSourceItemType);
        } else if (sources.length > 0) {
          setDataSource(sources[0] as keyof DataSourceItemType);
        }
        
        // 获取页码
        const pageParam = searchParams.get('page');
        if (pageParam) {
          setPage(parseInt(pageParam, 10));
        }
      } catch (err) {
        setError('获取数据源失败');
      } finally {
        setLoading(false);
      }
    };
    
    fetchDataSources();
  }, [searchParams, authEpoch, isAuthenticated]);

  // 数据源变化时获取筛选选项
  useEffect(() => {
    if (!isAuthenticated || !dataSource) return;
    
    const fetchFilterOptions = async () => {
      try {
        const options = await reviewApi.getFilterOptions(dataSource);
        setFilterOptions(options);
      } catch (err) {
        setError('获取筛选选项失败');
      }
    };
    
    fetchFilterOptions();
    // 重置页码
    setPage(1);
    // 更新URL
    updateUrl();
    
    // 获取题目列表
    fetchItems();
  }, [dataSource, authEpoch, isAuthenticated]);

  // 页码或筛选条件变化时获取题目列表
  useEffect(() => {
    if (!isAuthenticated || !dataSource) return;
    fetchItems();
    updateUrl();
  }, [page, filters, dataSource, authEpoch, isAuthenticated]);

  // 更新URL
  const updateUrl = () => {
    const params = new URLSearchParams(searchParams.toString());
    params.set('dataSource', dataSource as string);
    params.set('page', page.toString());
    router.push(`/questions?${params.toString()}`);
  };
  
  // 获取题目列表
  const fetchItems = async () => {
    setLoading(true);
    try {
      const strategy = StrategyFactory.getStrategy(dataSource as string);
      const filterParams = strategy?.formatFilterParams(filters) || {};
      
      const result = await reviewApi.getItems(dataSource, {
        ...filterParams,
        page: page,
        pageSize: pageSize,
      });
      
      setItems(result.items);
      setTotalItems(result.totalItems);
      setTotalPages(result.totalPages);
    } catch (err) {
      setError('获取题目列表失败');
      setItems([]);
    } finally {
      setLoading(false);
    }
  };

  // 处理数据源变更
  const handleDataSourceChange = (source: string) => {
    setDataSource(source as keyof DataSourceItemType);
  };

  // 处理筛选条件变更
  const handleFilterChange = (newFilters: any) => {
    setFilters(newFilters);
    setPage(1); // 重置页码
  };

  // 处理页码变更
  const handlePageChange = (newPage: number) => {
    setPage(newPage);
  };

  // 处理审核提交
  const handleReviewSubmit = async (itemId: string, statusUpdate: { status: string; comment: string }) => {
    try {
      await reviewApi.reviewItem(dataSource, itemId, statusUpdate);
      
      // 更新列表中的项目状态
      const updatedItems = items.map(item => 
        item.originalId === itemId
          ? { ...item, reviewStatus: statusUpdate.status, reviewComment: statusUpdate.comment }
          : item
      );
      setItems(updatedItems);
      
      // 成功提示 - 根据状态显示不同消息
      const statusText = statusUpdate.status === 'approved' ? '通过' : 
                         statusUpdate.status === 'rejected' ? '拒绝' : '待审核';
      toast.success(`题目已${statusText}，正在加载下一个待审核题目...`);

      // 在filters.status是pending时，可能需要从列表中移除此项，所以延迟一点时间再刷新
      if (filters.status === 'pending') {
        setTimeout(() => fetchItems(), 1000);
      }
    } catch (err) {
      console.error('提交审核失败:', err);
      toast.error('提交审核结果失败');
    }
  };

  const handleItemUpdated = (updatedItem: any) => {
    // 更新列表中的项目状态
    const updatedItems = items.map(item => 
      item.originalId === updatedItem.originalId
        ? updatedItem
        : item
    );
    setItems(updatedItems);
    
    // 如果审核后状态不再是pending，则可能需要刷新列表获取新的待审核题目
    if (updatedItem.reviewStatus !== 'pending' && filters.status === 'pending') {
      // 使用短暂延迟以确保UI有时间更新
      setTimeout(() => {
        fetchItems();
      }, 500);
    }
  };

  const strategy = dataSource ? StrategyFactory.getStrategy(dataSource as string) : null;
  
  // 设置策略的回调函数
  useEffect(() => {
    if (strategy) {
      strategy.setItemUpdateCallback(handleItemUpdated);
    }
  }, [strategy]);

  // 初始加载显示骨架屏
  if (loading && items.length === 0) {
    return (
      <div className="container mx-auto p-6">
        <div className="space-y-2">
          <Skeleton className="h-8 w-full" />
          <Skeleton className="h-8 w-2/3" />
          <Skeleton className="h-32 w-full" />
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-6 pb-20">
      <h1 className="text-2xl font-bold mb-6">题目审核</h1>
      
      {error && (
        <Alert variant="destructive" className="mb-4">
          <AlertTitle>错误</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
          <Button 
            variant="outline" 
            size="sm" 
            className="mt-2" 
            onClick={() => setError(null)}
          >
            关闭
          </Button>
        </Alert>
      )}
      
      {/* 顶部控制区 */}
      <div className="bg-white shadow-sm mb-6 p-4 rounded-lg">
        <div className="flex flex-wrap items-center justify-between gap-4 mb-4">
          <DataSourceSelector
            sources={availableSources}
            value={dataSource as string}
            onChange={handleDataSourceChange}
          />
          
          <div className="flex items-center text-sm text-gray-500">
            共 {totalItems} 个题目，第 {page}/{totalPages} 页
          </div>
        </div>
        
        {/* 筛选条件 */}
        {strategy && (
          <div className="mt-2">
            <Card>
              <CardContent className="pt-4 pb-4">
                {strategy.renderFilterOptions(filters, filterOptions, handleFilterChange)}
              </CardContent>
            </Card>
          </div>
        )}
      </div>
      
      {/* 分页控件 - 顶部 */}
      <div className="flex justify-between items-center mb-6">
        <Button
          variant="outline"
          onClick={() => handlePageChange(page - 1)}
          disabled={page <= 1 || loading}
        >
          上一页
        </Button>
        
        <span className="text-sm text-gray-500">
          第 {page} / {totalPages} 页
        </span>
        
        <Button
          variant="outline"
          onClick={() => handlePageChange(page + 1)}
          disabled={page >= totalPages || loading}
        >
          下一页
        </Button>
      </div>
      
      {/* 题目内容区 */}
      <div className="space-y-8">
        {loading ? (
          // 加载中占位
          Array.from({ length: 3 }).map((_, idx) => (
            <div key={idx} className="space-y-2">
              <Skeleton className="h-8 w-1/3" />
              <Skeleton className="h-24 w-full" />
              <Skeleton className="h-20 w-full" />
            </div>
          ))
        ) : items.length === 0 ? (
          // 无数据
          <div className="p-12 text-center text-gray-500 bg-gray-50 rounded-lg">
            {filters.status === 'pending' ? (
              <div className="space-y-3">
                <p>🎉 太棒了！所有题目都已审核完毕</p>
                <p className="text-sm">可以尝试切换其他数据源或修改筛选条件查看更多题目</p>
                <Button 
                  variant="outline" 
                  onClick={() => setFilters({...filters, status: 'all'})}
                  className="mt-2"
                >
                  查看所有题目
                </Button>
              </div>
            ) : (
              <div>没有找到匹配的题目</div>
            )}
          </div>
        ) : (
          // 题目内容
          items.map((item, index) => (
            <Card key={item.originalId} className="mb-8">
              <CardHeader className="flex flex-row items-center justify-between bg-gray-50">
                <h2 className="text-xl font-semibold">题目 {page * 10 - 10 + index + 1}</h2>
                <div className="text-sm text-gray-500">
                  ID: {item.originalId}
                </div>
              </CardHeader>
              <CardContent className="pt-6">
                <div className="space-y-6">
                  {/* 题目元数据（可折叠） */}
                  <div className="border rounded-md overflow-hidden mb-4">
                    <button
                      onClick={() => toggleMetadata(item.originalId)}
                      className="w-full flex justify-between items-center p-3 text-left bg-gray-50 hover:bg-gray-100 transition-colors"
                    >
                      <span className="font-medium">题目详情</span>
                      {metadataOpenMap[item.originalId] ? (
                        <ChevronUp className="h-5 w-5 text-gray-500" />
                      ) : (
                        <ChevronDown className="h-5 w-5 text-gray-500" />
                      )}
                    </button>
                    
                    {metadataOpenMap[item.originalId] && (
                      <div className="p-4 bg-white">
                        <div className="grid grid-cols-2 gap-4">
                          {item.questionSet && (
                            <div>
                              <h3 className="text-sm font-medium text-gray-500">题集</h3>
                              <p>{item.questionSet}</p>
                            </div>
                          )}
                          
                          {item.subject && (
                            <div>
                              <h3 className="text-sm font-medium text-gray-500">科目</h3>
                              <p>{item.subject}</p>
                            </div>
                          )}
                          
                          {item.difficulty && (
                            <div>
                              <h3 className="text-sm font-medium text-gray-500">难度</h3>
                              <p>{item.difficulty}</p>
                            </div>
                          )}
                          
                          {item.domain && (
                            <div>
                              <h3 className="text-sm font-medium text-gray-500">领域</h3>
                              <p>{item.domain}</p>
                            </div>
                          )}
                          
                          {item.skill && (
                            <div>
                              <h3 className="text-sm font-medium text-gray-500">技能</h3>
                              <p>{item.skill}</p>
                            </div>
                          )}
                          
                          {item.questionType && (
                            <div>
                              <h3 className="text-sm font-medium text-gray-500">题型</h3>
                              <p>{item.questionType}</p>
                            </div>
                          )}
                        </div>
                      </div>
                    )}
                  </div>
                  
                  {/* 知识点信息 */}
                  {item.knowledgePoint && (
                    <div className="mb-4 p-3 bg-blue-50 rounded-md">
                      <h3 className="text-sm font-medium text-blue-800 mb-1">知识点</h3>
                      <p className="text-blue-700">{item.knowledgePoint}</p>
                    </div>
                  )}
                  
                  {/* 题目详情内容 */}
                  <div className="border-b pb-4">
                    {strategy?.renderDetailView(item)}
                  </div>
                  
                  
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </div>
      
      {/* 分页控件 - 底部 */}
      {items.length > 0 && (
        <div className="flex justify-between items-center mt-8">
          <Button
            variant="outline"
            onClick={() => handlePageChange(page - 1)}
            disabled={page <= 1 || loading}
          >
            上一页
          </Button>
          
          <span className="text-sm text-gray-500">
            第 {page} / {totalPages} 页
          </span>
          
          <Button
            variant="outline"
            onClick={() => handlePageChange(page + 1)}
            disabled={page >= totalPages || loading}
          >
            下一页
          </Button>
        </div>
      )}
    </div>
  );
} 