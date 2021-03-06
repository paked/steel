BEGIN TRANSACTION;
CREATE TABLE "users" (
	`username`	TEXT NOT NULL,
	`password_hash`	TEXT NOT NULL,
	`salt`	TEXT,
	`id`	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	`email`	TEXT NOT NULL,
	`permission_level`	TEXT DEFAULT 0
);

CREATE TABLE "team_members" (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`student`	INTEGER,
	`submission`	INTEGER,
	`assignment`	INTEGER
);

CREATE TABLE "submissions" (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`team_name`	TEXT,
	`thoughts`	TEXT,
	`assignment`	TEXT
);

CREATE TABLE "assignments" (
	`id`	INTEGER PRIMARY KEY AUTOINCREMENT,
	`name`	TEXT,
	`description`	TEXT,
	`explanation`	TEXT,
	`due`	INTEGER,
	`created_by`	INTEGER,
    `class` INTEGER
);

CREATE TABLE "classes" (
    `id`    INTEGER PRIMARY KEY AUTOINCREMENT,
    `name`  TEXT,
    `description` TEXT,
    `image_url` TEXT
);

CREATE TABLE "students" (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `user` INTEGER,
    `permission_level` INTEGER,
    `class` INTEGER
);

CREATE TABLE  "workshop_pages" (
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `workshop` INTEGER,
    `title` TEXT,
    `contents` TEXT,
    `created` INTEGER,
    `updated` INTEGER,
    `sequence` INTEGER
);
;
COMMIT;
