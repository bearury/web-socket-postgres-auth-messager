-- 1. Добавляем колонку как NULLable
ALTER TABLE chats ADD COLUMN creator UUID REFERENCES users(id) ON DELETE CASCADE;

-- 2. Заполняем пустые значения дефолтным юзером (ВСТАВЬТЕ СВОЙ UUID вместо примера!)
UPDATE chats SET creator = 'c5f68de1-5086-4d90-b36e-777e5b4dbe44' WHERE creator IS NULL;

-- 3. Теперь спокойно делаем колонку NOT NULL
ALTER TABLE chats ALTER COLUMN creator SET NOT NULL;
