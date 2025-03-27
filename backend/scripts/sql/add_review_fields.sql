-- 添加缺少的review字段到t_math_sat_temp_classify_oneprep表
ALTER TABLE t_math_sat_temp_classify_oneprep 
ADD COLUMN IF NOT EXISTS review_comment TEXT DEFAULT '',
ADD COLUMN IF NOT EXISTS reviewed_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS reviewed_by TEXT DEFAULT '';

-- 添加注释以说明字段用途
COMMENT ON COLUMN t_math_sat_temp_classify_oneprep.review_comment IS '评审评论';
COMMENT ON COLUMN t_math_sat_temp_classify_oneprep.reviewed_at IS '评审时间';
COMMENT ON COLUMN t_math_sat_temp_classify_oneprep.reviewed_by IS '评审人';

-- 检查是否成功添加字段
SELECT 
    column_name, 
    data_type, 
    is_nullable, 
    column_default
FROM 
    information_schema.columns 
WHERE 
    table_schema = 'public' 
    AND table_name = 't_math_sat_temp_classify_oneprep'
    AND column_name IN ('review_comment', 'reviewed_at', 'reviewed_by'); 