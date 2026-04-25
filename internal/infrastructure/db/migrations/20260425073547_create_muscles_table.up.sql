CREATE TABLE muscles (
  id VARCHAR(36) PRIMARY KEY,
  code VARCHAR(32) NOT NULL UNIQUE,
  created_at TIMESTAMP NOT NULL
);

CREATE TABLE muscle_translations (
  muscle_id VARCHAR(36) NOT NULL REFERENCES muscles(id) ON DELETE CASCADE,
  lang VARCHAR(8) NOT NULL,
  name VARCHAR(255) NOT NULL,
  PRIMARY KEY (muscle_id, lang)
);

CREATE TABLE exercise_muscle (
  exercise_id VARCHAR(36) NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
  muscle_id VARCHAR(36) NOT NULL REFERENCES muscles(id),
  PRIMARY KEY (exercise_id, muscle_id)
);

INSERT INTO muscles (id, code, created_at) VALUES
  ('4b5a784a-3333-4721-a071-2e3fbd570c7f', 'chest',     NOW()),
  ('ea60868e-b2ee-428b-9b16-90dff700de4d', 'back',      NOW()),
  ('b127e524-2301-4bc6-b075-ad1f13e25115', 'legs',      NOW()),
  ('938c4201-0ac7-46c0-a74f-9ed682bbc339', 'shoulders', NOW()),
  ('24edcb38-ee90-48ef-925d-c12de371d20b', 'arms',      NOW()),
  ('a306f225-dcd2-43be-adf1-5b1bfbc6ea8a', 'core',      NOW());

INSERT INTO muscle_translations (muscle_id, lang, name) VALUES
  ('4b5a784a-3333-4721-a071-2e3fbd570c7f', 'ja', '胸'),
  ('ea60868e-b2ee-428b-9b16-90dff700de4d', 'ja', '背中'),
  ('b127e524-2301-4bc6-b075-ad1f13e25115', 'ja', '脚'),
  ('938c4201-0ac7-46c0-a74f-9ed682bbc339', 'ja', '肩'),
  ('24edcb38-ee90-48ef-925d-c12de371d20b', 'ja', '腕'),
  ('a306f225-dcd2-43be-adf1-5b1bfbc6ea8a', 'ja', '体幹'),
  ('4b5a784a-3333-4721-a071-2e3fbd570c7f', 'en', 'Chest'),
  ('ea60868e-b2ee-428b-9b16-90dff700de4d', 'en', 'Back'),
  ('b127e524-2301-4bc6-b075-ad1f13e25115', 'en', 'Legs'),
  ('938c4201-0ac7-46c0-a74f-9ed682bbc339', 'en', 'Shoulders'),
  ('24edcb38-ee90-48ef-925d-c12de371d20b', 'en', 'Arms'),
  ('a306f225-dcd2-43be-adf1-5b1bfbc6ea8a', 'en', 'Core');

INSERT INTO exercise_muscle (exercise_id, muscle_id) VALUES
  ('f1f538e5-4a37-409c-be99-09ee7bfefc50', '4b5a784a-3333-4721-a071-2e3fbd570c7f'),
  ('f1f538e5-4a37-409c-be99-09ee7bfefc50', '938c4201-0ac7-46c0-a74f-9ed682bbc339'),
  ('f1f538e5-4a37-409c-be99-09ee7bfefc50', '24edcb38-ee90-48ef-925d-c12de371d20b'),
  ('8b55202d-10aa-4007-a540-ca79e602d9ad', 'b127e524-2301-4bc6-b075-ad1f13e25115'),
  ('8b55202d-10aa-4007-a540-ca79e602d9ad', 'a306f225-dcd2-43be-adf1-5b1bfbc6ea8a'),
  ('e4af570d-163f-41ef-8bac-c6cc773606f4', 'ea60868e-b2ee-428b-9b16-90dff700de4d'),
  ('e4af570d-163f-41ef-8bac-c6cc773606f4', 'b127e524-2301-4bc6-b075-ad1f13e25115'),
  ('e4af570d-163f-41ef-8bac-c6cc773606f4', 'a306f225-dcd2-43be-adf1-5b1bfbc6ea8a'),
  ('43e18e58-c2fb-4a7e-97b2-8b66181ab45c', '938c4201-0ac7-46c0-a74f-9ed682bbc339'),
  ('43e18e58-c2fb-4a7e-97b2-8b66181ab45c', '24edcb38-ee90-48ef-925d-c12de371d20b'),
  ('43e18e58-c2fb-4a7e-97b2-8b66181ab45c', 'a306f225-dcd2-43be-adf1-5b1bfbc6ea8a'),
  ('0fadf73b-7a65-4cad-acc2-55d6eb076461', 'ea60868e-b2ee-428b-9b16-90dff700de4d'),
  ('0fadf73b-7a65-4cad-acc2-55d6eb076461', '24edcb38-ee90-48ef-925d-c12de371d20b'),
  ('f408191d-22db-4998-87f6-ebe9e76c5c60', 'ea60868e-b2ee-428b-9b16-90dff700de4d'),
  ('f408191d-22db-4998-87f6-ebe9e76c5c60', '24edcb38-ee90-48ef-925d-c12de371d20b'),
  ('660e3f0b-fc78-45f0-bb36-5533028e0a6b', 'ea60868e-b2ee-428b-9b16-90dff700de4d'),
  ('660e3f0b-fc78-45f0-bb36-5533028e0a6b', '24edcb38-ee90-48ef-925d-c12de371d20b'),
  ('6b1e867b-6744-411d-9576-e9f51dc84703', 'b127e524-2301-4bc6-b075-ad1f13e25115'),
  ('82906b9f-1503-4e28-800d-678a2d0f8606', '24edcb38-ee90-48ef-925d-c12de371d20b'),
  ('fa4a0ee1-4544-4964-896e-7b45edc0046f', '24edcb38-ee90-48ef-925d-c12de371d20b');
