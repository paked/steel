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
	`user`	INTEGER,
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
	`created_by`	INTEGER
);
;
COMMIT;
