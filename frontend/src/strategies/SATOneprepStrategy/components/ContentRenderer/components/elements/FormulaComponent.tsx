'use client';

import React from 'react';
import BaseComponent from '../base/BaseComponent';
import { ContentMetadata, ContentRendererOptions } from '../../types';
import { InlineMath, BlockMath } from 'react-katex';
import 'katex/dist/katex.min.css';

interface FormulaContent {
  content: string;
  display?: 'inline' | 'block';
}

interface FormulaComponentProps {
  content: FormulaContent;
  metadata?: ContentMetadata;
  options?: ContentRendererOptions;
}

/**
 * Component for rendering mathematical formulas using KaTeX
 */
const FormulaComponent: React.FC<FormulaComponentProps> = ({
  content,
  metadata = {},
  options = {}
}) => {
  if (!content || !content.content) return null;
  
  const { content: formulaContent, display = 'inline' } = content;
  const isBlock = display === 'block' || (metadata.style || []).includes('p');
  
  return (
    <BaseComponent metadata={metadata}>
      {isBlock ? (
        <BlockMath math={formulaContent} />
      ) : (
        <InlineMath math={formulaContent} />
      )}
    </BaseComponent>
  );
};

export default FormulaComponent; 