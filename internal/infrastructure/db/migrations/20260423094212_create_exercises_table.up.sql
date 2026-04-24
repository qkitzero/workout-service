CREATE TABLE exercises (
  id VARCHAR(36) PRIMARY KEY,
  code VARCHAR(64) NOT NULL UNIQUE,
  category VARCHAR(32) NOT NULL,
  created_at TIMESTAMP NOT NULL
);

CREATE TABLE exercise_translations (
  exercise_id VARCHAR(36) NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
  lang VARCHAR(8) NOT NULL,
  name VARCHAR(255) NOT NULL,
  PRIMARY KEY (exercise_id, lang)
);

INSERT INTO exercises (id, code, category, created_at) VALUES
  ('f1f538e5-4a37-409c-be99-09ee7bfefc50', 'bench_press',     'compound',  NOW()),
  ('8b55202d-10aa-4007-a540-ca79e602d9ad', 'squat',           'compound',  NOW()),
  ('e4af570d-163f-41ef-8bac-c6cc773606f4', 'deadlift',        'compound',  NOW()),
  ('43e18e58-c2fb-4a7e-97b2-8b66181ab45c', 'overhead_press',  'compound',  NOW()),
  ('0fadf73b-7a65-4cad-acc2-55d6eb076461', 'bent_over_row',   'compound',  NOW()),
  ('f408191d-22db-4998-87f6-ebe9e76c5c60', 'pull_up',         'compound',  NOW()),
  ('660e3f0b-fc78-45f0-bb36-5533028e0a6b', 'lat_pulldown',    'compound',  NOW()),
  ('6b1e867b-6744-411d-9576-e9f51dc84703', 'leg_press',       'compound',  NOW()),
  ('82906b9f-1503-4e28-800d-678a2d0f8606', 'dumbbell_curl',   'isolation', NOW()),
  ('fa4a0ee1-4544-4964-896e-7b45edc0046f', 'tricep_pushdown', 'isolation', NOW());

INSERT INTO exercise_translations (exercise_id, lang, name) VALUES
  ('f1f538e5-4a37-409c-be99-09ee7bfefc50', 'ja', 'ベンチプレス'),
  ('8b55202d-10aa-4007-a540-ca79e602d9ad', 'ja', 'スクワット'),
  ('e4af570d-163f-41ef-8bac-c6cc773606f4', 'ja', 'デッドリフト'),
  ('43e18e58-c2fb-4a7e-97b2-8b66181ab45c', 'ja', 'オーバーヘッドプレス'),
  ('0fadf73b-7a65-4cad-acc2-55d6eb076461', 'ja', 'ベントオーバーロウ'),
  ('f408191d-22db-4998-87f6-ebe9e76c5c60', 'ja', '懸垂'),
  ('660e3f0b-fc78-45f0-bb36-5533028e0a6b', 'ja', 'ラットプルダウン'),
  ('6b1e867b-6744-411d-9576-e9f51dc84703', 'ja', 'レッグプレス'),
  ('82906b9f-1503-4e28-800d-678a2d0f8606', 'ja', 'ダンベルカール'),
  ('fa4a0ee1-4544-4964-896e-7b45edc0046f', 'ja', 'トライセプスプッシュダウン');
