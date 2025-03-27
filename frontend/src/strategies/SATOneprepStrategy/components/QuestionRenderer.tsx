'use client';

import React from 'react';
import 'katex/dist/katex.min.css';
import ContentRenderer from './ContentRenderer';
import { ContentRendererOptions } from './ContentRenderer/types';

// 定义题目组件的内容类型
type ContentType = 
  | StringContent
  | FormulaContent
  | ImageContent 
  | HTMLContent
  | BlankContent
  | FigureContent
  | TableContent
  | string;

// 字符串内容
interface StringContent {
  type: 'string';
  content: string;
  metadata?: {
    style?: string[];
    id?: string;
    order?: number;
    [key: string]: any;
  };
}

// 公式内容
interface FormulaContent {
  type: 'formula';
  content: {
    type: 'latex' | 'mathml' | 'html';
    content: string;
    display?: 'block' | 'inline';
    alternatives?: {
      mathml?: string;
      ascii?: string;
    };
    variables?: Record<string, any>;
  };
  metadata?: {
    content_type?: string;
    [key: string]: any;
  };
}

// 图片内容
interface ImageContent {
  type: 'image';
  content: {
    type: 'url' | 'base64' | 'svg' | 'html';
    content: string;
    fallback?: string;
    thumbnail?: string;
    dimensions?: {
      width: number;
      height: number;
    };
    altText?: string;
    caption?: string;
  };
}

// HTML内容
interface HTMLContent {
  type: 'html';
  content: string;
}

// 空白填充
interface BlankContent {
  type: 'blank';
}

// 添加Figure内容类型
interface FigureContent {
  type: 'figure';
  content: {
    type: 'svg' | 'html';
    content: string;
    altText?: string;
    caption?: string;
  };
}

// 添加Table内容类型
interface TableContent {
  type: 'table';
  content: {
    type: 'html';
    content: string;
    caption?: string;
  };
  metadata?: {
    style?: string[];
    id?: string;
    [key: string]: any;
  };
}

// 选项
interface Option {
  identifier?: string;
  letter?: string;
  is_correct: boolean | string;
  content?: ContentType[];
  text?: ContentType[] | string;
}

// 单选答案
interface SingleChoiceAnswer {
  type: 'choice';
  format: {
    single_choice: {
      identifier: string;
    }
  };
  content?: ContentType[];
}

// 多选答案
interface MultipleChoiceAnswer {
  type: 'multiple_choice';
  format: {
    multiple_choice: {
      identifiers: string[];
    }
  };
  content?: ContentType[];
}

// 填空答案
interface TextAnswer {
  type: 'text';
  format: {
    text: {
      value: string;
      alternatives?: string[];
      case_sensitive?: boolean;
    }
  };
  content?: ContentType[];
}

// 答案类型
type AnswerType = SingleChoiceAnswer | MultipleChoiceAnswer | TextAnswer;

// 导出必要的类型供外部使用
export type { ContentType, Option, AnswerType };

// 组件属性
interface QuestionRendererProps {
  content: any; // 使用any类型以兼容现有数据结构
  options?: any;
  answer?: any;
  explanation?: any;
  className?: string;
  showAnswer?: boolean;
}

/**
 * SAT Oneprep题目渲染组件
 * 专门用于渲染SAT Oneprep格式的问题
 * 使用组件化架构和适配器模式
 */
const QuestionRenderer: React.FC<QuestionRendererProps> = ({
  content,
  options,
  answer,
  explanation,
  className = '',
  showAnswer = false
}) => {
  // 创建全局渲染选项
  const rendererOptions: ContentRendererOptions = {
    mergeConsecutiveStrings: true,  // 合并连续字符串
    convertOldRootRadicand: true,   // 转换旧式根号
    showSpecialSymbols: true        // 显示特殊符号
  };

  // 渲染选项
  const renderOptions = () => {
    // 安全检查：如果options为空或格式不正确，返回null
    if (!options) {
      console.log('选项为空');
      return null;
    }
    
    try {
      // 确保options是数组
      const optionsArray = Array.isArray(options) ? options : 
                         (typeof options === 'object' ? [options] : []);
      
      console.log('渲染选项，数量:', optionsArray.length);
      
      if (optionsArray.length === 0) return null;
      
      return (
        <div className="sat-options mt-4">
          <div className="text-lg font-medium mb-3">选项</div>
          <div className="space-y-4">
            {optionsArray.map((option: any, index: number) => {
              // 安全检查：如果选项为null或undefined，跳过
              if (!option) return null;

              // 安全获取选项标识符
              let identifier;
              try {
                identifier = option.identifier || option.letter || String.fromCharCode(65 + index);
              } catch (e) {
                identifier = `选项${index + 1}`;
                console.error('选项标识符处理错误:', e);
              }
              
              // 安全判断是否正确选项
              const isCorrect = showAnswer && 
                (option.is_correct === true || 
                 option.is_correct === 'True' || 
                 option.is_correct === 'true');
              
              return (
                <div 
                  key={index} 
                  className={`sat-option p-4 rounded-md ${isCorrect ? 'bg-green-50 border border-green-200' : 'bg-gray-50'}`}
                >
                  <span className="font-bold mr-3">{identifier}.</span>
                  
                  {/* 使用ContentRenderer渲染选项内容 */}
                  <span className="option-content">
                    {option.content ? (
                      <ContentRenderer content={option.content} options={rendererOptions} />
                    ) : option.text ? (
                      typeof option.text === 'string' ? (
                        option.text
                      ) : (
                        <ContentRenderer content={option.text} options={rendererOptions} />
                      )
                    ) : (
                      <span className="text-gray-400">无内容</span>
                    )}
                  </span>
                </div>
              );
            })}
          </div>
        </div>
      );
    } catch (error) {
      console.error('选项整体渲染错误:', error);
      return (
        <div className="sat-options mt-4">
          <div className="text-lg font-medium mb-2">选项</div>
          <div className="p-2 bg-red-50 text-red-500 rounded">
            选项数据格式不正确，无法显示
          </div>
        </div>
      );
    }
  };

  // 渲染答案
  const renderAnswer = () => {
    if (!showAnswer || !answer) {
      console.log('不显示答案 - showAnswer:', showAnswer, '答案是否存在:', !!answer);
      return null;
    }
    
    try {
      console.log('渲染答案，类型:', typeof answer);
      
      // 答案标题
      const title = <div className="text-lg font-medium mb-3">答案</div>;
      
      // 处理不同类型的答案
      let answerContent;
      
      // 如果answer是字符串，尝试解析成JSON
      if (typeof answer === 'string') {
        try {
          const parsedAnswer = JSON.parse(answer);
          console.log('答案从字符串解析为对象，解析成功');
          if (parsedAnswer && typeof parsedAnswer === 'object') {
            // 使用ContentRenderer渲染解析后的答案
            answerContent = renderParsedAnswer(parsedAnswer);
          } else {
            // 如果解析后不是对象，直接显示字符串
            console.log('解析后的答案不是对象');
            answerContent = <div className="p-4 bg-green-50 rounded">{answer}</div>;
          }
        } catch (e) {
          // 解析失败，直接显示字符串
          console.error('答案JSON解析错误:', e);
          answerContent = <div className="p-4 bg-green-50 rounded">{answer}</div>;
        }
      } else if (answer && typeof answer === 'object') {
        // 如果answer已经是对象
        console.log('答案已经是对象，类型:', answer.type || '无类型');
        answerContent = renderParsedAnswer(answer);
      } else {
        // 不支持的答案格式
        console.log('不支持的答案格式');
        answerContent = <div className="p-4 bg-red-50 text-red-500 rounded">无法识别的答案格式</div>;
      }
      
      return (
        <div className="sat-answer mt-6">
          {title}
          {answerContent}
        </div>
      );
    } catch (error) {
      console.error('答案渲染错误:', error);
      return (
        <div className="sat-answer mt-6">
          <div className="text-lg font-medium mb-3">答案</div>
          <div className="p-4 bg-red-50 text-red-500 rounded">
            答案渲染过程中发生错误
          </div>
        </div>
      );
    }
  };
  
  // 处理已解析的答案对象
  const renderParsedAnswer = (parsedAnswer: any): React.ReactNode => {
    try {
      // 处理选择题情况 - 如果答案中包含identifier和content
      if (parsedAnswer.identifier && parsedAnswer.content) {
        return (
          <div className="p-2 bg-green-50 rounded">
            <div className="flex items-center mb-2">
              <span className="font-medium mr-2">正确选项:</span> 
              <span className="font-bold text-lg">{parsedAnswer.identifier}</span>
            </div>
            
            {/* 使用ContentRenderer渲染内容 */}
            <div className="mt-1">
              <ContentRenderer content={parsedAnswer.content} options={rendererOptions} />
            </div>
          </div>
        );
      }

      // 单选答案
      if (
        parsedAnswer.type === 'choice' && 
        parsedAnswer.format && 
        parsedAnswer.format.single_choice && 
        parsedAnswer.format.single_choice.identifier
      ) {
        const identifier = parsedAnswer.format.single_choice.identifier;
        return (
          <div className="p-2 bg-green-50 rounded">
            <div className="flex items-center mb-2">
              <span className="font-medium mr-2">正确选项:</span> 
              <span className="font-bold text-lg">{identifier}</span>
            </div>
            
            {parsedAnswer.content && (
              <div className="mt-1">
                <ContentRenderer content={parsedAnswer.content} options={rendererOptions} />
              </div>
            )}
          </div>
        );
      }
      
      // 多选答案
      if (
        parsedAnswer.type === 'multiple_choice' && 
        parsedAnswer.format && 
        parsedAnswer.format.multiple_choice && 
        parsedAnswer.format.multiple_choice.identifiers
      ) {
        const identifiers = parsedAnswer.format.multiple_choice.identifiers;
        return (
          <div className="p-2 bg-green-50 rounded">
            <div className="flex items-center mb-2">
              <span className="font-medium mr-2">正确选项:</span> 
              {Array.isArray(identifiers) 
                ? identifiers.map((id: string, index: number) => (
                    <span key={index} className="font-bold text-lg mx-1">{id}{index < identifiers.length - 1 ? ',' : ''}</span>
                  ))
                : <span className="font-bold text-lg">{identifiers}</span>
              }
            </div>
            
            {parsedAnswer.content && (
              <div className="mt-1">
                <ContentRenderer content={parsedAnswer.content} options={rendererOptions} />
              </div>
            )}
          </div>
        );
      }
      
      // 文本答案
      if (
        parsedAnswer.type === 'text' && 
        parsedAnswer.format && 
        parsedAnswer.format.text
      ) {
        const textValue = parsedAnswer.format.text.value;
        const alternatives = parsedAnswer.format.text.alternatives;
        return (
          <div className="p-2 bg-green-50 rounded">
            <div className="mb-2">
              <span className="font-medium">答案:</span> {textValue}
            </div>
            {alternatives && alternatives.length > 0 && (
              <div className="mb-2">
                <span className="font-medium">可接受的答案:</span> 
                {alternatives.map((alt: string, i: number) => (
                  <span key={i} className="ml-1">{alt}{i < alternatives.length - 1 ? ',' : ''}</span>
                ))}
              </div>
            )}
            {parsedAnswer.content && (
              <div className="mt-1">
                <ContentRenderer content={parsedAnswer.content} options={rendererOptions} />
              </div>
            )}
          </div>
        );
      }
      
      // 未知类型的答案
      return (
        <div className="p-2 bg-green-50 rounded">
          <pre className="text-sm overflow-auto">{JSON.stringify(parsedAnswer, null, 2)}</pre>
        </div>
      );
    } catch (error) {
      console.error('解析后的答案渲染错误:', error);
      return (
        <div className="p-2 bg-red-50 text-red-500 rounded">
          答案格式无效
        </div>
      );
    }
  };

  // 渲染解析
  const renderExplanation = () => {
    if (!explanation || !showAnswer) {
      console.log('不显示解析 - showAnswer:', showAnswer, '解析是否存在:', !!explanation);
      return null;
    }
    
    try {
      console.log('渲染解析，类型:', typeof explanation, Array.isArray(explanation) ? '数组' : '非数组');
      
      return (
        <div className="sat-explanation mt-6 p-4 bg-gray-50 rounded-md">
          <div className="text-lg font-medium mb-3">解析</div>
          <div>
            <ContentRenderer content={explanation} options={rendererOptions} />
          </div>
        </div>
      );
    } catch (error) {
      console.error('解析渲染错误:', error);
      return (
        <div className="sat-explanation mt-6 p-4 bg-red-50 rounded-md">
          <div className="text-lg font-medium mb-2">解析</div>
          <div className="text-red-500">解析渲染过程中发生错误</div>
        </div>
      );
    }
  };

  return (
    <div className={`sat-oneprep-question ${className}`}>
      <div className="sat-question-content">
        <ContentRenderer content={content} options={rendererOptions} />
      </div>
      {renderOptions()}
      {renderAnswer()}
      {renderExplanation()}
    </div>
  );
};

export default QuestionRenderer; 