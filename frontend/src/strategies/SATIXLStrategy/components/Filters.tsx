'use client';

import React, { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

interface FiltersProps {
  filters: any;
  options: any;
  onChange: (filters: any) => void;
}

const Filters: React.FC<FiltersProps> = ({ filters, options, onChange }) => {
  const [localFilters, setLocalFilters] = useState({
    reviewStatus: filters?.reviewStatus || 'all',
    skill: filters?.skill || 'all',
    searchTerm: filters?.searchTerm || '',
  });

  useEffect(() => {
    setLocalFilters({
      reviewStatus: filters?.reviewStatus || 'all',
      skill: filters?.skill || 'all',
      searchTerm: filters?.searchTerm || '',
    });
  }, [filters]);

  const handleChange = (key: string, value: any) => {
    setLocalFilters((prev: any) => ({
      ...prev,
      [key]: value,
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // 转换回API需要的格式
    const apiFilters = {
      ...localFilters,
      reviewStatus: localFilters.reviewStatus === 'all' ? '' : localFilters.reviewStatus,
      skill: localFilters.skill === 'all' ? '' : localFilters.skill,
    };
    onChange(apiFilters);
  };

  const handleReset = () => {
    const resetFilters = {
      reviewStatus: 'all',
      skill: 'all',
      searchTerm: '',
    };
    setLocalFilters(resetFilters);
    onChange({
      reviewStatus: '',
      skill: '',
      searchTerm: '',
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div>
        <Label htmlFor="reviewStatus">审核状态</Label>
        <Select
          value={localFilters.reviewStatus}
          onValueChange={(value) => handleChange('reviewStatus', value)}
        >
          <SelectTrigger>
            <SelectValue placeholder="选择审核状态" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">全部</SelectItem>
            {options?.reviewStatus?.map((status: string) => (
              <SelectItem key={status} value={status}>
                {status === 'approved'
                  ? '已通过'
                  : status === 'rejected'
                  ? '已拒绝'
                  : '待审核'}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div>
        <Label htmlFor="skill">技能</Label>
        <Select
          value={localFilters.skill}
          onValueChange={(value) => handleChange('skill', value)}
        >
          <SelectTrigger>
            <SelectValue placeholder="选择技能" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">全部</SelectItem>
            {options?.skill?.map((skill: string) => (
              <SelectItem key={skill} value={skill}>
                {skill}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div>
        <Label htmlFor="searchTerm">搜索</Label>
        <Input
          id="searchTerm"
          placeholder="输入关键词"
          value={localFilters.searchTerm || ''}
          onChange={(e) => handleChange('searchTerm', e.target.value)}
        />
      </div>

      <div className="flex space-x-2">
        <Button type="submit">应用筛选</Button>
        <Button
          type="button"
          variant="outline"
          onClick={handleReset}
        >
          重置
        </Button>
      </div>
    </form>
  );
};

export default Filters; 