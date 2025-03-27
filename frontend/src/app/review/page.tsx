'use client';

import React, { useState, useEffect } from 'react';
import DataSourceSelector from '@/components/common/DataSourceSelector';
import { StrategyFactory } from '@/strategies/factory';
import { reviewApi } from '@/api/review';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { ReviewableItem, PaginatedResult } from '@/models/reviewable-item';
import ReviewForm from '@/components/common/ReviewForm';
import { Button } from '@/components/ui/button';
import { DataTable } from '@/components/common/DataTable';
import { useSearchParams, useRouter } from 'next/navigation';

export default function ReviewPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [dataSource, setDataSource] = useState<string>('');
  const [availableSources, setAvailableSources] = useState<string[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [items, setItems] = useState<PaginatedResult<ReviewableItem>>({
    items: [],
    total: 0,
    page: 1,
    pageSize: 10,
  });
  const [selectedItem, setSelectedItem] = useState<ReviewableItem | null>(null);
  const [filterOptions, setFilterOptions] = useState<any>({});
  const [filters, setFilters] = useState<any>({});
  const [viewMode, setViewMode] = useState<'table' | 'list'>('table');
  const [updatingItem, setUpdatingItem] = useState<boolean>(false);

  useEffect(() => {
    // 获取可用的数据源
    const fetchDataSources = async () => {
      try {
        const sources = await reviewApi.getDataSources();
        setAvailableSources(sources);
        
        // 检查URL中是否有dataSource参数
        const dataSourceParam = searchParams.get('dataSource');
        if (dataSourceParam && StrategyFactory.isDataSourceSupported(dataSourceParam)) {
          setDataSource(dataSourceParam);
        } else if (sources.length > 0) {
          setDataSource(sources[0]);
        }
        
        setLoading(false);
      } catch (err) {
        setError('获取数据源失败');
        setLoading(false);
      }
    };
    
    fetchDataSources();
  }, [searchParams]);

  useEffect(() => {
    if (!dataSource) return;
    
    const fetchFilterOptions = async () => {
      try {
        const options = await reviewApi.getFilterOptions(dataSource);
        setFilterOptions(options);
      } catch (err) {
        setError('获取筛选选项失败');
      }
    };
    
    fetchFilterOptions();
    
    // 更新URL
    const params = new URLSearchParams(searchParams.toString());
    params.set('dataSource', dataSource);
    router.push(`/review?${params.toString()}`);
    
    // 重置筛选条件和已选择的项目
    setFilters({});
    setSelectedItem(null);
    
    // 加载项目
    fetchItems({
      page: 1,
      pageSize: 10,
    });
  }, [dataSource, router, searchParams]);

  const fetchItems = async (params: any) => {
    if (!dataSource) return;
    
    setLoading(true);
    try {
      const strategy = StrategyFactory.getStrategy(dataSource);
      const filterParams = strategy?.formatFilterParams(filters) || {};
      
      const result = await reviewApi.getItems(dataSource, {
        ...filterParams,
        ...params,
      });
      
      setItems(result);
      setLoading(false);
    } catch (err) {
      setError('获取项目列表失败');
      setLoading(false);
    }
  };

  const handleDataSourceChange = (source: string) => {
    setDataSource(source);
  };

  const handleItemSelect = async (item: ReviewableItem) => {
    setSelectedItem(null); // 先清空，显示加载状态
    
    try {
      const detailedItem = await reviewApi.getItemById(dataSource, item.originalId);
      setSelectedItem(detailedItem);
    } catch (err) {
      setError('获取项目详情失败');
    }
  };

  const handleFilterChange = (newFilters: any) => {
    setFilters(newFilters);
    fetchItems({
      page: 1,
      pageSize: items.pageSize,
    });
  };

  const handlePageChange = (page: number, pageSize: number) => {
    fetchItems({
      page,
      pageSize,
    });
  };

  const handleReviewSubmit = async (statusUpdate: any) => {
    if (!selectedItem) return;
    
    setUpdatingItem(true);
    try {
      await reviewApi.reviewItem(dataSource, selectedItem.originalId, statusUpdate);
      
      // 刷新项目列表
      fetchItems({
        page: items.page,
        pageSize: items.pageSize,
      });
      
      // 重新获取当前项目详情
      const updatedItem = await reviewApi.getItemById(dataSource, selectedItem.originalId);
      setSelectedItem(updatedItem);
      
      setUpdatingItem(false);
    } catch (err) {
      setError('提交审核结果失败');
      setUpdatingItem(false);
    }
  };

  const strategy = dataSource ? StrategyFactory.getStrategy(dataSource) : null;

  if (loading && !dataSource) {
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

  if (error && !dataSource) {
    return (
      <div className="container mx-auto p-6">
        <Alert variant="destructive">
          <AlertTitle>错误</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      </div>
    );
  }

  return (
    <div className="container mx-auto p-6">
      <h1 className="text-2xl font-bold mb-6">题目审核系统</h1>
      
      {error && (
        <Alert variant="destructive" className="mb-4">
          <AlertTitle>错误</AlertTitle>
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
      
      <div className="mb-6">
        <DataSourceSelector
          sources={availableSources}
          value={dataSource}
          onChange={handleDataSourceChange}
        />
      </div>
      
      {!strategy && dataSource && (
        <Alert variant="destructive">
          <AlertTitle>不支持的数据源</AlertTitle>
          <AlertDescription>选择的数据源 {dataSource} 不受支持</AlertDescription>
        </Alert>
      )}
      
      {strategy && (
        <div className="grid grid-cols-1 md:grid-cols-12 gap-6">
          <div className="md:col-span-4">
            <Card>
              <CardContent className="pt-6">
                {strategy.renderFilterOptions(filters, filterOptions, handleFilterChange)}
              </CardContent>
            </Card>
          </div>
          
          <div className="md:col-span-8">
            <Tabs defaultValue="table" onValueChange={(v) => setViewMode(v as 'table' | 'list')}>
              <div className="flex justify-between items-center mb-4">
                <TabsList>
                  <TabsTrigger value="table">表格视图</TabsTrigger>
                  <TabsTrigger value="list">列表视图</TabsTrigger>
                </TabsList>
              </div>
              
              <TabsContent value="table" className="mt-0">
                <Card>
                  <CardContent className="pt-6">
                    <DataTable
                      columns={strategy.getTableColumns()}
                      data={items.items}
                      loading={loading}
                      pagination={{
                        current: items.page,
                        pageSize: items.pageSize,
                        total: items.total,
                        onChange: handlePageChange,
                      }}
                      onRowClick={handleItemSelect}
                    />
                  </CardContent>
                </Card>
              </TabsContent>
              
              <TabsContent value="list" className="mt-0">
                <Card>
                  <CardContent className="pt-6">
                    {loading ? (
                      <div className="space-y-2">
                        {Array.from({ length: 5 }).map((_, idx) => (
                          <Skeleton key={idx} className="h-20 w-full" />
                        ))}
                      </div>
                    ) : (
                      <div className="space-y-2">
                        {items.items.map((item) => (
                          <div
                            key={item.originalId}
                            className="border rounded-md p-4 cursor-pointer hover:bg-gray-50"
                            onClick={() => handleItemSelect(item)}
                          >
                            {strategy.renderListItem(item)}
                          </div>
                        ))}
                      </div>
                    )}
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>
          </div>
        </div>
      )}
      
      {selectedItem && strategy && (
        <div className="mt-8">
          <h2 className="text-xl font-bold mb-4">项目详情</h2>
          <div className="grid grid-cols-1 md:grid-cols-12 gap-6">
            <div className="md:col-span-8">
              <Card>
                <CardContent className="pt-6">
                  {strategy.renderDetailView(selectedItem)}
                </CardContent>
              </Card>
            </div>
            
            <div className="md:col-span-4">
              <Card>
                <CardContent className="pt-6">
                  <h3 className="text-lg font-medium mb-4">提交审核</h3>
                  <ReviewForm 
                    item={selectedItem} 
                    onSubmit={handleReviewSubmit} 
                  />
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      )}
    </div>
  );
} 