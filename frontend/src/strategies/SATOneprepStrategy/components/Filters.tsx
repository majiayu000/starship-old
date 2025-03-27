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
  SelectValue
} from '@/components/ui/select';
import { FilterOptions } from '@/api/types';

interface FiltersProps {
  filters: any;
  options: FilterOptions;
  onChange: (filters: any) => void;
}

const Filters: React.FC<FiltersProps> = ({ filters, options, onChange }) => {
  const [localFilters, setLocalFilters] = useState({
    status: filters.status || 'all',
    domain: filters.domain || 'all',
    skill: filters.skill || 'all',
    searchTerm: filters.searchTerm || ''
  });

  // 当props中的filters变化时更新本地状态
  useEffect(() => {
    setLocalFilters({
      status: filters.status || 'all',
      domain: filters.domain || 'all',
      skill: filters.skill || 'all',
      searchTerm: filters.searchTerm || ''
    });
  }, [filters]);

  const handleChange = (field: string, value: string) => {
    setLocalFilters(prev => ({
      ...prev,
      [field]: value
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // 转换回API需要的格式
    const apiFilters = {
      ...localFilters,
      status: localFilters.status === 'all' ? '' : localFilters.status,
      domain: localFilters.domain === 'all' ? '' : localFilters.domain,
      skill: localFilters.skill === 'all' ? '' : localFilters.skill
    };
    onChange(apiFilters);
  };

  const handleReset = () => {
    const resetFilters = {
      status: 'all',
      domain: 'all',
      skill: 'all',
      searchTerm: ''
    };
    setLocalFilters(resetFilters);
    onChange({
      status: '',
      domain: '',
      skill: '',
      searchTerm: ''
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <h2 className="text-lg font-semibold">筛选条件</h2>
      
      <div className="space-y-2">
        <Label htmlFor="status">审核状态</Label>
        <Select
          value={localFilters.status}
          onValueChange={(value) => handleChange('status', value)}
        >
          <SelectTrigger id="status">
            <SelectValue placeholder="选择审核状态" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">全部</SelectItem>
            {options.reviewStatus?.map(status => (
              <SelectItem key={status} value={status}>
                {status === 'pending' ? '待审核' : 
                 status === 'approved' ? '已通过' : 
                 status === 'rejected' ? '已拒绝' : status}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-2">
        <Label htmlFor="domain">领域</Label>
        <Select
          value={localFilters.domain}
          onValueChange={(value) => handleChange('domain', value)}
        >
          <SelectTrigger id="domain">
            <SelectValue placeholder="选择领域" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">全部</SelectItem>
            {options.domain?.map(domain => (
              <SelectItem key={domain} value={domain}>{domain}</SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-2">
        <Label htmlFor="skill">技能</Label>
        <Select
          value={localFilters.skill}
          onValueChange={(value) => handleChange('skill', value)}
        >
          <SelectTrigger id="skill">
            <SelectValue placeholder="选择技能" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">全部</SelectItem>
            {options.skill?.map(skill => (
              <SelectItem key={skill} value={skill}>{skill}</SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      <div className="space-y-2">
        <Label htmlFor="searchTerm">搜索</Label>
        <Input
          id="searchTerm"
          placeholder="搜索内容"
          value={localFilters.searchTerm}
          onChange={(e) => handleChange('searchTerm', e.target.value)}
        />
      </div>

      <div className="flex space-x-2 pt-2">
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