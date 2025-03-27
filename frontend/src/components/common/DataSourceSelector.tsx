'use client';

import React from 'react';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from '@/components/ui/select';

interface DataSourceSelectorProps {
  sources: string[];
  value: string;
  onChange: (value: string) => void;
}

const DataSourceSelector: React.FC<DataSourceSelectorProps> = ({ sources, value, onChange }) => {
  // 将技术数据源名称转换为更友好的显示名称
  const getDisplayName = (source: string) => {
    switch (source) {
      case 'sat_oneprep':
        return 'SAT Oneprep';
      case 'sat_ixl':
        return 'SAT IXL';
      default:
        return source;
    }
  };

  return (
    <div className="flex flex-col gap-2">
      <label htmlFor="data-source" className="text-sm font-medium">
        数据源
      </label>
      <Select
        value={value}
        onValueChange={onChange}
      >
        <SelectTrigger className="w-[200px]" id="data-source">
          <SelectValue placeholder="选择数据源" />
        </SelectTrigger>
        <SelectContent>
          {sources.map((source) => (
            <SelectItem key={source} value={source}>
              {getDisplayName(source)}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
};

export default DataSourceSelector; 