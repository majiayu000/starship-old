import React from 'react';
import { ReviewStrategy, TableColumn } from '../ReviewStrategy';
import { SATOneprepItem } from '@/models/sat-oneprep-item';
import { QueryParams } from '@/models/reviewable-item';
import DetailView from './components/DetailView';
import Filters from './components/Filters';
import { formatSATOneprepFilters } from './utils';
import { Badge } from '@/components/ui/badge';

// 修改 ReviewStrategy 接口以支持 onItemUpdated 回调
declare module '../ReviewStrategy' {
  interface ReviewStrategy<T extends ReviewableItem = ReviewableItem> {
    setItemUpdateCallback(callback: (item: T) => void): void;
  }
}

class SATOneprepStrategy implements ReviewStrategy<SATOneprepItem> {
  private onItemUpdated: ((item: SATOneprepItem) => void) | null = null;

  setItemUpdateCallback(callback: (item: SATOneprepItem) => void): void {
    this.onItemUpdated = callback;
  }
  
  getTableColumns(): TableColumn[] {
    return [
      { 
        title: '题目ID', 
        dataIndex: 'originalId', 
        key: 'originalId',
        width: 120
      },
      { 
        title: '题集', 
        dataIndex: 'questionSet', 
        key: 'questionSet',
        width: 150
      },
      { 
        title: '领域', 
        dataIndex: 'domain', 
        key: 'domain',
        width: 120
      },
      { 
        title: '技能', 
        dataIndex: 'skill', 
        key: 'skill',
        width: 120
      },
      { 
        title: '状态', 
        dataIndex: 'reviewStatus', 
        key: 'reviewStatus',
        width: 100,
        render: (status: string) => {
          const variant = 
            status === 'approved' ? 'default' : 
            status === 'rejected' ? 'destructive' : 
            'outline';
          
          const text = 
            status === 'approved' ? '已通过' : 
            status === 'rejected' ? '已拒绝' : 
            '待审核';
          
          return (
            <Badge variant={variant}>{text}</Badge>
          );
        }
      },
      { 
        title: '审核人', 
        dataIndex: 'reviewedBy', 
        key: 'reviewedBy',
        width: 120,
        render: (reviewedBy: string) => reviewedBy || '-'
      }
    ];
  }
  
  renderListItem(item: SATOneprepItem): React.ReactNode {
    return (
      <div className="flex justify-between items-center p-2">
        <div>
          <div className="font-medium">{item.questionSet}</div>
          <div className="text-sm text-gray-500">
            {item.domain} / {item.skill}
          </div>
        </div>
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
              : '待审核'
          }
        </Badge>
      </div>
    );
  }
  
  renderDetailView(item: SATOneprepItem): React.ReactNode {
    return <DetailView item={item} onItemUpdated={this.onItemUpdated || undefined} />;
  }
  
  renderFilterOptions(filters: any, options: any, onChange: (filters: any) => void): React.ReactNode {
    return <Filters filters={filters} options={options} onChange={onChange} />;
  }
  
  formatFilterParams(filters: any): Partial<QueryParams> {
    return formatSATOneprepFilters(filters);
  }
  
  canHandle(dataSource: string): boolean {
    return dataSource === 'sat_oneprep';
  }
}

export { SATOneprepStrategy };
export default SATOneprepStrategy; 