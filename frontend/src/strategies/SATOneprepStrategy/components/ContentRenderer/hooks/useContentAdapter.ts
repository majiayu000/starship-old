import { ContentItemType, ContentRendererOptions, ContentTypeEnum, ParagraphContainer, StringContent, ContentMetadata, BlankContent } from '../types';

/**
 * 内容适配器钩子 - 处理和标准化内容数据
 * 
 * 实现功能：
 * 1. 字符串合并
 * 2. 根号转换
 * 3. 段落分组
 * 4. 一般数据标准化
 */
export function useContentAdapter(
  content: ContentItemType[] | ContentItemType | string | null,
  options: ContentRendererOptions
): ContentItemType[] {
  if (!content) return [];
  
  // 将单个项目转换为数组
  let contentArray: ContentItemType[] = Array.isArray(content) 
    ? content 
    : (typeof content === 'string' 
      ? [{ type: ContentTypeEnum.STRING, content } as ContentItemType] 
      : [content]);
  
  // 预处理：转换旧式根号
  if (options.convertOldRootRadicand) {
    contentArray = convertOldRootRadicand(contentArray);
  }
  
  // 预处理：合并连续字符串
  if (options.mergeConsecutiveStrings) {
    contentArray = mergeConsecutiveStrings(contentArray);
  }
  
  // 预处理：增强空白填充元素
  contentArray = enhanceBlankElements(contentArray);
  
  // 预处理：规范化figure内容
  contentArray = normalizeFigureContent(contentArray);
  
  // 处理：段落分组
  if (options.enableParagraphGrouping) {
    return createParagraphs(contentArray);
  }
  
  return contentArray;
}

/**
 * 特殊处理：优化blank元素在文本中的位置和样式
 */
function enhanceBlankElements(contentArray: ContentItemType[]): ContentItemType[] {
  return contentArray.map((item, index, array) => {
    if (item.type === ContentTypeEnum.BLANK) {
      // 获取前后元素上下文
      const prevItem = index > 0 ? array[index - 1] : null;
      const nextItem = index < array.length - 1 ? array[index + 1] : null;
      
      // 克隆blank元素以添加增强功能
      const enhancedBlank = Object.assign({}, item) as BlankContent;
      
      // 根据位置设置宽度（如果没有自定义宽度）
      if (!enhancedBlank.content) {
        enhancedBlank.content = { width: '5em' };
        
        // 如果前后都有文本，可能需要更合适的宽度
        if (prevItem?.type === ContentTypeEnum.STRING && nextItem?.type === ContentTypeEnum.STRING) {
          enhancedBlank.content.width = '4em';
        }
        // 如果是段落末尾的blank
        else if (prevItem?.type === ContentTypeEnum.STRING && (!nextItem || nextItem.type === ContentTypeEnum.LINE_BREAK)) {
          enhancedBlank.content.width = '6em';
        }
      }
      
      // 添加标记以便调试
      if (!enhancedBlank.metadata) {
        enhancedBlank.metadata = {};
      }
      enhancedBlank.metadata.enhanced = true;
      enhancedBlank.metadata.position = `${index}`;
      
      return enhancedBlank;
    }
    
    return item;
  });
}

/**
 * 转换旧式根号表示为LaTeX格式
 */
function convertOldRootRadicand(contentArray: ContentItemType[]): ContentItemType[] {
  return contentArray.map(item => {
    // 只处理HTML类型内容
    if (item.type === ContentTypeEnum.HTML && typeof item.content === 'string') {
      const htmlContent = item.content;
      
      // 查找旧式根号模式 <span class="old-root-radicand">16</span>
      const regex = /<span\s+class="old-root-radicand">(.*?)<\/span>/g;
      let match;
      let newContent = htmlContent;
      
      // 如果找到匹配，替换为LaTeX格式
      while ((match = regex.exec(htmlContent)) !== null) {
        const radicand = match[1];
        
        // 创建新的公式内容
        newContent = newContent.replace(match[0], `<span class="latex-formula">\\sqrt{${radicand}}</span>`);
      }
      
      // 如果内容被修改，替换为转换后的内容
      if (newContent !== htmlContent) {
        const newItem = Object.assign({}, item, {
          content: newContent
        });
        return newItem;
      }
    }
    
    return item;
  });
}

/**
 * 合并连续的字符串内容
 */
function mergeConsecutiveStrings(contentArray: ContentItemType[]): ContentItemType[] {
  if (contentArray.length <= 1) return contentArray;
  
  const result: ContentItemType[] = [];
  let currentString: StringContent | null = null;
  
  for (let i = 0; i < contentArray.length; i++) {
    const item = contentArray[i];
    
    // 处理字符串类型内容
    if (item.type === ContentTypeEnum.STRING) {
      const stringItem = item as StringContent;
      
      // 如果当前已有字符串，且内容类型相同，则合并
      if (
        currentString && 
        currentString.type === ContentTypeEnum.STRING &&
        currentString.metadata?.content_type === stringItem.metadata?.content_type
      ) {
        // 合并字符串内容，使用空格连接
        const content: string = `${currentString.content} ${stringItem.content}`;
        
        // 创建新的合并样式
        const combinedStyles: string[] = [];
        
        // 添加当前字符串的样式
        if (currentString.metadata?.style) {
          currentString.metadata.style.forEach(style => combinedStyles.push(style));
        }
        
        // 添加新字符串的样式
        if (stringItem.metadata?.style) {
          stringItem.metadata.style.forEach(style => {
            if (!combinedStyles.includes(style)) {
              combinedStyles.push(style);
            }
          });
        }
        
        // 更新当前字符串
        const newMetadata: ContentMetadata = currentString.metadata 
          ? Object.assign({}, currentString.metadata, { style: combinedStyles })
          : { style: combinedStyles };
          
        currentString = {
          type: ContentTypeEnum.STRING,
          content,
          metadata: newMetadata
        };
      } else {
        // 如果之前有其他内容，添加到结果
        if (currentString) {
          result.push(currentString);
        }
        
        // 开始一个新的字符串
        currentString = Object.assign({}, stringItem);
      }
    } else {
      // 如果有待处理的字符串，先添加到结果
      if (currentString) {
        result.push(currentString);
        currentString = null;
      }
      
      // 添加非字符串内容
      result.push(item);
    }
  }
  
  // 处理最后剩余的字符串
  if (currentString) {
    result.push(currentString);
  }
  
  return result;
}

/**
 * 创建段落结构
 */
function createParagraphs(contentArray: ContentItemType[]): ContentItemType[] {
  if (contentArray.length === 0) return [];
  
  const result: ContentItemType[] = [];
  let currentParagraph: ParagraphContainer | null = null;
  let paragraphId = 0;
  
  // 处理每个内容项
  for (let i = 0; i < contentArray.length; i++) {
    const item = contentArray[i];
    const hasBreak = 
      item.type === ContentTypeEnum.LINE_BREAK || 
      item.metadata?.style?.includes('p');
    
    // 如果当前没有段落，创建一个新段落
    if (!currentParagraph) {
      currentParagraph = {
        type: ContentTypeEnum.PARAGRAPH,
        id: `p-${++paragraphId}`,
        content: [],
        metadata: {}
      };
    }
    
    // 先添加当前内容到段落
    if (currentParagraph) {
      currentParagraph.content.push(item);
    }
    
    // 如果当前内容是段落分隔符，或者是最后一个内容
    if (hasBreak || i === contentArray.length - 1) {
      // 添加当前段落到结果，确保不为空
      if (currentParagraph) {
        result.push(currentParagraph);
      }
      
      // 重置当前段落
      currentParagraph = null;
    }
  }
  
  return result;
}

/**
 * 规范化图表内容，确保figure类型的内容格式一致
 */
function normalizeFigureContent(contentArray: ContentItemType[]): ContentItemType[] {
  return contentArray.map(item => {
    if (item.type === ContentTypeEnum.FIGURE) {
      // 克隆figure项以进行修改
      const figureItem = Object.assign({}, item) as any;
      
      // 如果content是字符串，将其转换为标准格式
      if (typeof figureItem.content === 'string') {
        const content = figureItem.content;
        figureItem.content = {
          type: content.includes('<svg') ? 'svg' : 'html',
          content: content
        };
      }
      
      return figureItem;
    }
    
    return item;
  });
} 