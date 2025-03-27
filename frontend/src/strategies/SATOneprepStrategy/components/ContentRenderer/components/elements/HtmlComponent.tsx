'use client';

import React from 'react';
import BaseComponent from '../base/BaseComponent';
import { ContentMetadata, ContentRendererOptions } from '../../types';

interface HtmlComponentProps {
  content: string;
  metadata?: ContentMetadata;
  options?: ContentRendererOptions;
}

/**
 * Component for safely rendering HTML content
 */
const HtmlComponent: React.FC<HtmlComponentProps> = ({
  content,
  metadata = {},
  options = {}
}) => {
  if (!content) return null;
  
  // Sanitize HTML content to remove potential security risks
  const sanitizedContent = sanitizeHtml(content);
  
  return (
    <BaseComponent metadata={metadata}>
      <div 
        className="html-content" 
        dangerouslySetInnerHTML={{ __html: sanitizedContent }}
      />
    </BaseComponent>
  );
};

/**
 * Basic HTML sanitization to prevent script injection
 */
function sanitizeHtml(html: string): string {
  return html
    .replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')
    .replace(/on\w+="[^"]*"/gi, '')
    .replace(/javascript:/gi, '');
}

export default HtmlComponent; 