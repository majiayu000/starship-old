-- 检查表是否存在
SELECT EXISTS (
   SELECT FROM information_schema.tables 
   WHERE table_schema = 'public' 
   AND table_name = 't_math_sat_temp_classify_oneprep'
);

-- 获取表的确切字段信息（按顺序）
SELECT 
    column_name, 
    data_type, 
    character_maximum_length,
    is_nullable, 
    column_default,
    ordinal_position
FROM 
    information_schema.columns 
WHERE 
    table_schema = 'public' 
    AND table_name = 't_math_sat_temp_classify_oneprep'
ORDER BY 
    ordinal_position;

-- 获取主键信息
SELECT 
    tc.constraint_name, 
    kcu.column_name
FROM 
    information_schema.table_constraints tc
    JOIN information_schema.key_column_usage kcu
        ON tc.constraint_catalog = kcu.constraint_catalog
        AND tc.constraint_schema = kcu.constraint_schema
        AND tc.constraint_name = kcu.constraint_name
WHERE 
    tc.table_schema = 'public'
    AND tc.table_name = 't_math_sat_temp_classify_oneprep'
    AND tc.constraint_type = 'PRIMARY KEY';

-- 检查review相关字段
SELECT 
    column_name
FROM 
    information_schema.columns 
WHERE 
    table_schema = 'public' 
    AND table_name = 't_math_sat_temp_classify_oneprep'
    AND column_name IN ('review_status', 'review_comment', 'reviewed_at', 'reviewed_by');

-- 查看表的一行数据来了解实际内容
SELECT * FROM t_math_sat_temp_classify_oneprep LIMIT 1; 