'use client';

import React, { useState } from 'react';
import { Button } from '@/components/ui/button';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Label } from '@/components/ui/label';
import { ReviewableItem } from '@/models/reviewable-item';

interface ReviewFormProps {
  item: ReviewableItem;
  onSubmit: (statusUpdate: { status: string; comment: string }) => void;
}

const ReviewForm: React.FC<ReviewFormProps> = ({ item, onSubmit }) => {
  const [status, setStatus] = useState<string>(item.reviewStatus || 'pending');
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    
    onSubmit({
      status,
      comment: '' // 保留空字符串以保持接口兼容性
    });
    
    setIsSubmitting(false);
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="review-status">审核状态</Label>
        <RadioGroup 
          id="review-status" 
          value={status} 
          onValueChange={setStatus} 
          className="flex space-x-4"
        >
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="approved" id="approved" />
            <Label htmlFor="approved" className="cursor-pointer">通过</Label>
          </div>
          
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="rejected" id="rejected" />
            <Label htmlFor="rejected" className="cursor-pointer">拒绝</Label>
          </div>
          
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="pending" id="pending" />
            <Label htmlFor="pending" className="cursor-pointer">待审核</Label>
          </div>
        </RadioGroup>
      </div>
      
      <div className="flex justify-end">
        <Button 
          type="submit" 
          disabled={isSubmitting || status === item.reviewStatus}
        >
          {isSubmitting ? '提交中...' : '提交审核结果'}
        </Button>
      </div>
    </form>
  );
};

export default ReviewForm; 