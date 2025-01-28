CREATE TYPE file_type_enum AS ENUM('PHOTO', 'TEXT', 'VIDEO') 

CREATE TABLE files (
    id SERIAL PRIMARY KEY,          
    file_name VARCHAR(255) NOT NULL, 
    file_type file_type_enum NOT NULL,
    file_size BIGINT NOT NULL,       
    upload_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Время загрузки
    path_to_file TEXT NOT NULL       -- Путь к файлу в MinIO
);

CREATE INDEX idx_files_file_name ON files (file_name);