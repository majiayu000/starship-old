'use client';

import React, { useState, useEffect } from 'react';
import { SATOneprepItem } from '@/models/sat-oneprep-item';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Button } from '@/components/ui/button';
import QuestionRenderer from './QuestionRenderer';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import ReviewForm from '@/components/common/ReviewForm';
import { reviewApi } from '@/api/review';
import { toast } from 'sonner';

interface DetailViewProps {
  item: SATOneprepItem;
  onItemUpdated?: (updatedItem: SATOneprepItem) => void;
}

// 定义错误边界组件
class ErrorBoundary extends React.Component<
  { children: React.ReactNode, fallback?: React.ReactNode },
  { hasError: boolean, error: Error | null }
> {
  constructor(props: { children: React.ReactNode, fallback?: React.ReactNode }) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Component Error:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback || (
        <div className="p-4 bg-red-50 rounded-md">
          <h3 className="text-lg font-medium text-red-800 mb-2">渲染错误</h3>
          <details>
            <summary className="text-sm cursor-pointer text-red-600">查看错误详情</summary>
            <pre className="mt-2 text-xs bg-gray-100 p-2 rounded overflow-auto">
              {this.state.error?.toString()}
            </pre>
          </details>
        </div>
      );
    }

    return this.props.children;
  }
}

const DetailView: React.FC<DetailViewProps> = ({ item, onItemUpdated }) => {
  const [showAnswer, setShowAnswer] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  
  // 状态变量来存储解析后的数据
  const [parsedData, setParsedData] = useState<{
    questionContent: any;
    options: any;
    answer: any;
    explanation: any;
  }>({
    questionContent: null,
    options: null,
    answer: null,
    explanation: null
  });

  // 使用useEffect解析JSON数据
  useEffect(() => {
    console.log('DetailView: 开始解析数据');
    setLoading(true);
    setError(null);
    
    try {
      // 创建新的解析数据对象
      const newParsedData = {
        questionContent: parseJSON(item.questionContent, 'questionContent'),
        options: parseJSON(item.options, 'options'),
        answer: parseJSON(item.answer, 'answer'),
        explanation: parseJSON(item.explanation, 'explanation')
      };
      
      // 打印解析结果的数据类型
      console.log('DetailView: 数据解析完成', {
        questionContent: typeof newParsedData.questionContent + (Array.isArray(newParsedData.questionContent) ? `[${newParsedData.questionContent.length}]` : ''),
        options: typeof newParsedData.options + (Array.isArray(newParsedData.options) ? `[${newParsedData.options.length}]` : ''),
        answer: typeof newParsedData.answer,
        explanation: typeof newParsedData.explanation + (Array.isArray(newParsedData.explanation) ? `[${newParsedData.explanation.length}]` : '')
      });
      
      setParsedData(newParsedData);
      setLoading(false);
    } catch (err) {
      console.error('DetailView: 数据解析错误', err);
      setError(`数据解析错误: ${err instanceof Error ? err.message : String(err)}`);
      setLoading(false);
    }
  }, [item]);
  
  // 安全地解析JSON
  const parseJSON = (data: any, fieldName: string): any => {
    // 如果已经是对象，直接返回
    if (data && typeof data === 'object') {
      console.log(`${fieldName}已经是对象类型`);
      return data;
    }
    
    // 如果是空值或非字符串，返回原始值
    if (!data || typeof data !== 'string') {
      console.log(`${fieldName}是空值或非字符串:`, data);
      return data;
    }
    
    try {
      // 尝试清理和处理数据
      let cleanedData = data.trim();
      
      // 检测是否以引号开始，如果不是可能不是JSON
      if (!(cleanedData.startsWith('{') || cleanedData.startsWith('[') || 
            cleanedData.startsWith('"') || cleanedData.startsWith("'"))) {
        console.log(`${fieldName}不是标准JSON格式，尝试返回原始值:`, cleanedData);
        return data;
      }
      
      // 特别处理：修复类似 .1764 这样的无效数字格式
      cleanedData = cleanedData.replace(/("\s*:\s*)\./g, '$10.');
      
      // 移除潜在的 BOM 标记
      if (cleanedData.charCodeAt(0) === 0xFEFF) {
        cleanedData = cleanedData.slice(1);
      }
      
      try {
        const parsed = JSON.parse(cleanedData);
        console.log(`${fieldName}成功解析为`, Array.isArray(parsed) ? `数组[${parsed.length}]` : '对象');
        return parsed;
      } catch (err) {
        // 如果解析失败，记录详细错误并尝试其他处理方式
        console.error(`${fieldName}解析失败:`, err);
        console.log('问题数据片段:', cleanedData.slice(0, 100) + '...');
        
        // 尝试使用更宽松的方式解析
        try {
          // 对于特殊格式的数字进行全局替换
          cleanedData = cleanedData.replace(/:\s*\.(\d+)/g, ': 0.$1');
          const parsed = JSON.parse(cleanedData);
          console.log(`${fieldName}在特殊处理后成功解析`);
          return parsed;
        } catch (err2) {
          console.error(`${fieldName}在特殊处理后仍解析失败:`, err2);
          // 返回原始数据，防止页面崩溃
          return data;
        }
      }
    } catch (err) {
      console.error(`${fieldName}处理过程中出错:`, err);
      return data; // 返回原始数据
    }
  };

  // 处理表单提交
  const handleReviewSubmit = async (statusUpdate: { status: string; comment: string }) => {
    if (submitting) return;
    
    try {
      setSubmitting(true);
      
      console.log('提交审核', {
        dataSource: 'sat_oneprep',
        id: item.originalId,
        statusUpdate
      });
      
      await reviewApi.reviewItem(
        'sat_oneprep', 
        item.originalId, 
        {
          ...statusUpdate,
          reviewedBy: localStorage.getItem('username') || 'unknown_user' // 可以从本地存储获取用户名
        }
      );
      
      // 更新本地界面状态
      const now = new Date().toISOString();
      const updatedItem = {
        ...item,
        reviewStatus: statusUpdate.status,
        reviewComment: statusUpdate.comment,
        reviewedAt: now,
        reviewedBy: localStorage.getItem('username') || 'unknown_user'
      };
      
      // 如果有回调，可以更新父组件状态
      if (onItemUpdated) {
        onItemUpdated(updatedItem);
      }
      
      toast.success('审核提交成功');
    } catch (error) {
      console.error('审核提交失败:', error);
      toast.error('审核提交失败，请重试');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Card className="w-full">
      <CardHeader className="flex flex-row items-center justify-between">
        
        <Badge
          variant={
            item.reviewStatus === 'approved'
              ? 'default'
              : item.reviewStatus === 'rejected'
                ? 'destructive'
                : 'outline'
          }
        >
          {item.reviewStatus === 'approved'
            ? '已通过'
            : item.reviewStatus === 'rejected'
              ? '已拒绝'
              : '待审核'}
        </Badge>
      </CardHeader>
      
      <CardContent className="space-y-4">
        {loading ? (
          <div className="flex items-center justify-center h-40">
            <p>加载中...</p>
          </div>
        ) : error ? (
          <div className="bg-red-50 p-4 rounded-md text-red-800">
            <p>{error}</p>
          </div>
        ) : (
          <>
            <Separator />
            
            <div>
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg font-medium">题目内容</h3>
                <Button 
                  variant="outline" 
                  onClick={() => setShowAnswer(!showAnswer)}
                >
                  {showAnswer ? '隐藏答案' : '显示答案'}
                </Button>
              </div>
              
              <ErrorBoundary>
                <div className="bg-white p-4 rounded-md border">
                  <QuestionRenderer
                    content={parsedData.questionContent}
                    options={parsedData.options}
                    answer={parsedData.answer}
                    explanation={parsedData.explanation}
                    showAnswer={showAnswer}
                  />
                </div>
              </ErrorBoundary>
            </div>
            
            <Separator />
            
            <ReviewForm item={item} onSubmit={handleReviewSubmit} />
          </>
        )}
      </CardContent>
    </Card>
  );
};

export default DetailView; 