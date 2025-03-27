import { ReviewableItem } from './reviewable-item';

export interface SATOneprepItem extends ReviewableItem {
  url: string;
  questionSet: string;
  subject: string;
  difficulty: string;
  domain: string;
  skill: string;
  questionType: string;
  questionContent: any; // JSON content
  explanation: any; // JSON content
  answer: any; // JSON content
  options: any; // JSON content
  questionId: string;
  knowledgePoint: string;
} 