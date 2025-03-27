'use client';

import React, { useState, useEffect } from 'react';
import { SATIXLItem } from '@/models/sat-ixl-item';
import JsonContent from '@/components/common/JsonContent';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Button } from '@/components/ui/button';
import ReviewForm from '@/components/common/ReviewForm';
import { reviewApi } from '@/api/review';
import { toast } from 'sonner';
import QuestionRenderer from '../../SATOneprepStrategy/components/QuestionRenderer';
import { Card, CardContent, CardHeader } from '@/components/ui/card';

interface DetailViewProps {
  item: SATIXLItem;
  onItemUpdated?: (updatedItem: SATIXLItem) => void;
}

const DetailView: React.FC<DetailViewProps> = ({ item, onItemUpdated }) => {
  const [submitting, setSubmitting] = useState(false);
  const [showAnswer, setShowAnswer] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // 状态变量来存储解析后的数据
  const [parsedData, setParsedData] = useState<{
    content: any;
    answer: any;
    explanation: any;
  }>({
    content: null,
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
        content: parseJSON(item.content, 'content'),
        answer: parseJSON(item.answer, 'answer'),
        explanation: item.explanation ? parseJSON(item.explanation, 'explanation') : null
      };
      
      // 打印解析结果的数据类型
      console.log('DetailView: 数据解析完成', {
        content: typeof newParsedData.content + (Array.isArray(newParsedData.content) ? `[${newParsedData.content.length}]` : ''),
        answer: typeof newParsedData.answer + (Array.isArray(newParsedData.answer) ? `[${newParsedData.answer.length}]` : ''),
        explanation: item.explanation ? (
          typeof newParsedData.explanation + (Array.isArray(newParsedData.explanation) ? `[${newParsedData.explanation.length}]` : '')
        ) : '不存在'
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
      
      try {
        const parsed = JSON.parse(cleanedData);
        console.log(`${fieldName}成功解析为`, Array.isArray(parsed) ? `数组[${parsed.length}]` : '对象');
        return parsed;
      } catch (err) {
        console.error(`${fieldName}解析失败:`, err);
        return data; // 返回原始数据
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
        dataSource: 'sat_ixl',
        id: item.originalId,
        statusUpdate
      });
      
      await reviewApi.reviewItem(
        'sat_ixl', 
        item.originalId, 
        {
          ...statusUpdate,
          reviewedBy: localStorage.getItem('username') || 'unknown_user'
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
        <div className="grid grid-cols-2 gap-4">
          <div>
            <h3 className="text-sm font-medium">题目ID</h3>
            <p>{item.questionId}</p>
          </div>
          <div>
            <h3 className="text-sm font-medium">技能</h3>
            <p>{item.skill}</p>
          </div>
          <div>
            <h3 className="text-sm font-medium">知识点</h3>
            <p>{item.knowledgePoint || '-'}</p>
          </div>
        </div>
        
        <Separator />
        
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
              
              <div className="bg-white p-4 rounded-md border">
                <QuestionRenderer
                  content={parsedData.content}
                  answer={parsedData.answer}
                  explanation={parsedData.explanation}
                  showAnswer={showAnswer}
                />
              </div>
            </div>
          </>
        )}
        
        {item.reviewedBy && (
          <div>
            <h3 className="text-sm font-medium mb-1">审核人</h3>
            <p>{item.reviewedBy}</p>
          </div>
        )}
        
        <Separator />
        
        <ReviewForm item={item} onSubmit={handleReviewSubmit} />
      </CardContent>
    </Card>
  );
};

export default DetailView; 