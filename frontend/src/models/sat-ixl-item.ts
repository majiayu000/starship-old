import { ReviewableItem } from './reviewable-item';

export interface SATIXLItem extends ReviewableItem {
  questionId: string;
  content: any; // JSON content
  answer: any; // JSON content
  explanation?: any; // JSON content for explanation
  skill: string;
  knowledgePoint: string;
  // 根据实际IXL结构添加其他字段
} 