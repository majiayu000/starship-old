import React from 'react';
import { ReviewStrategy, TableColumn } from '../ReviewStrategy';
import { SATIXLItem } from '@/models/sat-ixl-item';
import { QueryParams } from '@/models/reviewable-item';
import DetailView from './components/DetailView';
import Filters from './components/Filters';
import { formatSATIXLFilters } from './utils';
import { Badge } from '@/components/ui/badge';

class SATIXLStrategy implements ReviewStrategy<SATIXLItem> {
  private onItemUpdated: ((item: SATIXLItem) => void) | null = null;

  setItemUpdateCallback(callback: (item: SATIXLItem) => void): void {
    this.onItemUpdated = callback;
  }
  
  getTableColumns(): TableColumn[] {
    return [
      { 
        title: '题目ID', 
        dataIndex: 'questionId', 
        key: 'questionId',
        width: 120
      },
      { 
        title: '技能', 
        dataIndex: 'skill', 
        key: 'skill',
        width: 150
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
  
  renderListItem(item: SATIXLItem): React.ReactNode {
    return (
      <div className="flex justify-between items-center p-2">
        <div>
          <div className="font-medium">ID: {item.questionId}</div>
          <div className="text-sm text-gray-500">
            技能: {item.skill}
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
  
  renderDetailView(item: SATIXLItem): React.ReactNode {
    return <DetailView item={item} onItemUpdated={this.onItemUpdated || undefined} />;
  }
  
  renderFilterOptions(filters: any, options: any, onChange: (filters: any) => void): React.ReactNode {
    return <Filters filters={filters} options={options} onChange={onChange} />;
  }
  
  formatFilterParams(filters: any): Partial<QueryParams> {
    return formatSATIXLFilters(filters);
  }
  
  canHandle(dataSource: string): boolean {
    return dataSource === 'sat_ixl';
  }
}

export { SATIXLStrategy };
export default SATIXLStrategy; 