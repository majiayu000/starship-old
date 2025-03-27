declare module 'react-katex' {
  import React from 'react';
  
  export interface InlineMathProps {
    math: string;
    renderError?: (error: Error) => React.ReactNode;
    errorColor?: string;
  }
  
  export interface BlockMathProps {
    math: string;
    renderError?: (error: Error) => React.ReactNode;
    errorColor?: string;
  }
  
  export const InlineMath: React.FC<InlineMathProps>;
  export const BlockMath: React.FC<BlockMathProps>;
} 